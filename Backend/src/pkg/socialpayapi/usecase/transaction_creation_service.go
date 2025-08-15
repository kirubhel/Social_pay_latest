package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	commission_usecase "github.com/socialpay/socialpay/src/pkg/commission/usecase"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

// TransactionCreationService provides unified transaction creation logic
type TransactionCreationService struct {
	commissionUsecase commission_usecase.CommissionUseCase
	logger            logging.Logger
}

func NewTransactionCreationService(
	commissionUsecase commission_usecase.CommissionUseCase,
	logger logging.Logger,
) *TransactionCreationService {
	return &TransactionCreationService{
		commissionUsecase: commissionUsecase,
		logger:            logger,
	}
}

// TransactionCreationRequest contains all data needed to create a transaction
type TransactionCreationRequest struct {
	// Basic transaction data
	UserID      uuid.UUID
	MerchantID  uuid.UUID
	BaseAmount  float64
	Description string
	Medium      txEntity.TransactionMedium
	Type        txEntity.TransactionType
	PaymentType string // "direct", "checkout", "qr"

	Status txEntity.TransactionStatus

	// Optional tip information
	TipAmount  *float64                    `json:"tip_amount,omitempty"`
	TipeePhone *string                     `json:"tipee_phone,omitempty"`
	TipMedium  *txEntity.TransactionMedium `json:"tip_medium,omitempty"`

	// Fee configuration
	MerchantPaysFee bool

	// URLs and metadata
	CallbackURL string
	SuccessURL  string
	FailedURL   string
	QRTag       *string
	Details     map[string]interface{}
}

// TransactionCreationResponse contains the created transaction with calculated amounts
type TransactionCreationResponse struct {
	Transaction    *txEntity.Transaction
	BaseAmount     float64
	TotalAmount    float64
	CustomerNet    float64
	MerchantNet    float64
	AdminNet       float64
	CommissionRate float64
	TipAmount      float64
}

// CreateTransaction creates a transaction with proper amount calculations and commission
func (s *TransactionCreationService) CreateTransaction(ctx context.Context, req TransactionCreationRequest) (*TransactionCreationResponse, error) {
	// Get merchant-specific commission rate
	commissionResult, err := s.commissionUsecase.CalculateCommission(ctx, req.BaseAmount, req.MerchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate commission: %w", err)
	}

	s.logger.Info("Commission calculated", map[string]interface{}{
		"merchant_id":     req.MerchantID,
		"base_amount":     req.BaseAmount,
		"commission_rate": commissionResult.Percent,
		"commission":      commissionResult.Percent * req.BaseAmount / 100,
		"payment_type":    req.PaymentType,
	})

	// Handle tip amount
	tipAmount := 0.0
	if req.TipAmount != nil && *req.TipAmount > 0 {
		tipAmount = *req.TipAmount
		// Validate tip requirements
		if req.TipeePhone == nil || *req.TipeePhone == "" {
			return nil, fmt.Errorf("tipee_phone is required when tip_amount > 0")
		}
		if req.TipMedium == nil {
			return nil, fmt.Errorf("tip_medium is required when tip_amount > 0")
		}
	}

	// Calculate amounts based on transaction type
	var baseAmount, totalAmount, customerNet, merchantNet, adminNet float64
	var feeAmount, vatAmount float64

	baseAmount = req.BaseAmount
	feeAmount = (baseAmount * commissionResult.Percent) / 100.0
	feeAmount += commissionResult.Cent // Add fixed cent amount
	vatAmount = feeAmount * 0.15       // 15% VAT
	adminNet = feeAmount               // Admin gets fee minus VAT

	switch req.Type {
	case txEntity.DEPOSIT:
		// DEPOSIT: Customer pays base_amount + commission + tip
		if req.MerchantPaysFee {
			// Merchant pays fee: customer pays only base amount + tip
			totalAmount = baseAmount + tipAmount + vatAmount + feeAmount
			customerNet = baseAmount
			merchantNet = baseAmount - feeAmount - vatAmount // Merchant gets less
		} else {
			// Customer pays fee: customer pays base + fee + tip
			totalAmount = baseAmount + feeAmount + vatAmount + tipAmount
			customerNet = totalAmount
			merchantNet = baseAmount // Merchant gets full base amount
		}

	case txEntity.WITHDRAWAL:
		// WITHDRAWAL: Merchant pays base_amount + commission (positive values)
		if req.MerchantPaysFee {
			// Merchant pays fee: total includes base + commission
			totalAmount = baseAmount + feeAmount + vatAmount
			merchantNet = totalAmount // Merchant pays this amount (POSITIVE)
			customerNet = baseAmount  // Customer receives base amount
		} else {
			// Customer pays fee: customer receives less
			totalAmount = baseAmount + vatAmount + feeAmount
			merchantNet = baseAmount                         // Merchant pays base amount (POSITIVE)
			customerNet = baseAmount - feeAmount - vatAmount // Customer receives amount minus fee
		}
	default:
		return nil, fmt.Errorf("unsupported transaction type: %s", req.Type)
	}

	// Create transaction entity
	transaction := &txEntity.Transaction{
		Id:              uuid.New(),
		UserId:          req.UserID,
		MerchantId:      req.MerchantID,
		BaseAmount:      baseAmount,
		TotalAmount:     totalAmount,
		CustomerNet:     customerNet,
		MerchantNet:     merchantNet,
		AdminNet:        adminNet,
		FeeAmount:       feeAmount,
		VatAmount:       vatAmount,
		Description:     req.Description,
		Medium:          req.Medium,
		Type:            req.Type,
		Status:          req.Status,
		CallbackURL:     req.CallbackURL,
		SuccessURL:      req.SuccessURL,
		FailedURL:       req.FailedURL,
		Details:         req.Details,
		MerchantPaysFee: req.MerchantPaysFee,
	}

	// Add QR tag if provided
	if req.QRTag != nil {
		transaction.QRTag = req.QRTag
	}

	// Add tip information if provided
	if tipAmount > 0 {
		transaction.HasTip = true
		transaction.TipAmount = &tipAmount
		transaction.TipeePhone = req.TipeePhone
		if req.TipMedium != nil {
			tipMediumStr := string(*req.TipMedium)
			transaction.TipMedium = &tipMediumStr
		}
	}

	response := &TransactionCreationResponse{
		Transaction:    transaction,
		BaseAmount:     baseAmount,
		TotalAmount:    totalAmount,
		CustomerNet:    customerNet,
		MerchantNet:    merchantNet,
		AdminNet:       adminNet,
		CommissionRate: commissionResult.Percent,
		TipAmount:      tipAmount,
	}

	s.logger.Info("Transaction amounts", map[string]interface{}{
		"Amounts": response,
	})

	return response, nil
}
