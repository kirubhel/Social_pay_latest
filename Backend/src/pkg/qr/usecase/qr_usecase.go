package usecase

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/google/uuid"

	SocialpayUsecase "github.com/socialpay/socialpay/src/pkg/socialpayapi/usecase"
	"github.com/socialpay/socialpay/src/pkg/qr/core/entity"
	"github.com/socialpay/socialpay/src/pkg/qr/core/repository"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	txRepo "github.com/socialpay/socialpay/src/pkg/transaction/core/repository"
	transaction_usecase "github.com/socialpay/socialpay/src/pkg/transaction/usecase"
	walletUsecase "github.com/socialpay/socialpay/src/pkg/wallet/usecase"
)

// QRUseCase defines the interface for QR link business operations
type QRUseCase interface {
	// CreateQRLink creates a new QR payment link
	CreateQRLink(ctx context.Context, userID, merchantID uuid.UUID, req *entity.CreateQRLinkRequest) (*entity.QRLinkResponse, error)

	// GetQRLink retrieves a QR link by ID
	GetQRLink(ctx context.Context, id uuid.UUID) (*entity.QRLinkResponse, error)

	// GetQRLinksByMerchant retrieves QR links for a merchant with pagination
	GetQRLinksByMerchant(ctx context.Context, merchantID uuid.UUID, pagination *pagination.Pagination) (*entity.QRLinksListResponse, error)

	// GetQRLinksByUser retrieves QR links for a user with pagination
	GetQRLinksByUser(ctx context.Context, userID uuid.UUID, pagination *pagination.Pagination) (*entity.QRLinksListResponse, error)

	// UpdateQRLink updates an existing QR link
	UpdateQRLink(ctx context.Context, id, userID uuid.UUID, req *entity.UpdateQRLinkRequest) (*entity.QRLinkResponse, error)

	// DeleteQRLink soft deletes a QR link
	DeleteQRLink(ctx context.Context, id, userID uuid.UUID) error

	// ProcessQRPayment processes a payment using a QR link
	ProcessQRPayment(ctx context.Context, qrLinkID uuid.UUID, req *entity.QRPaymentRequest) (*entity.QRPaymentResponse, error)
}

type qrUseCase struct {
	qrRepo             repository.QRRepository
	transactionRepo    txRepo.TransactionRepository
	transactionUseCase transaction_usecase.TransactionUseCase
	paymentService     SocialpayUsecase.PaymentProcessor
	walletUseCase      walletUsecase.WalletUseCase
	log                logging.Logger
}

func NewQRUseCase(
	qrRepo repository.QRRepository,
	transactionRepo txRepo.TransactionRepository,
	transactionUseCase transaction_usecase.TransactionUseCase,
	paymentService SocialpayUsecase.PaymentProcessor,
	walletUseCase walletUsecase.WalletUseCase,
) QRUseCase {
	return &qrUseCase{
		qrRepo:             qrRepo,
		transactionRepo:    transactionRepo,
		transactionUseCase: transactionUseCase,
		paymentService:     paymentService,
		walletUseCase:      walletUseCase,
		log:                logging.NewStdLogger("qr_usecase"),
	}
}

func (uc *qrUseCase) CreateQRLink(ctx context.Context, userID, merchantID uuid.UUID, req *entity.CreateQRLinkRequest) (*entity.QRLinkResponse, error) {
	uc.log.Info("Creating QR link", map[string]interface{}{
		"user_id":     userID,
		"merchant_id": merchantID,
		"type":        req.Type,
		"tag":         req.Tag,
	})

	if err := req.Validate(); err != nil {
		uc.log.Error("QR link creation request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	qrLink := &entity.QRLink{
		ID:               uuid.New(),
		UserID:           userID,
		MerchantID:       merchantID,
		Type:             req.Type,
		Amount:           req.Amount,
		SupportedMethods: req.SupportedMethods,
		Tag:              req.Tag,
		Title:            req.Title,
		Description:      req.Description,
		ImageURL:         req.ImageURL,
		IsTipEnabled:     req.IsTipEnabled,
		IsActive:         true,
	}

	if err := uc.qrRepo.Create(ctx, qrLink); err != nil {
		uc.log.Error("Failed to create QR link", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create QR link: %w", err)
	}

	uc.log.Info("QR link created successfully", map[string]interface{}{
		"qr_link_id": qrLink.ID,
	})

	return uc.buildQRLinkResponse(qrLink), nil
}

func (uc *qrUseCase) GetQRLink(ctx context.Context, id uuid.UUID) (*entity.QRLinkResponse, error) {
	uc.log.Info("Getting QR link", map[string]interface{}{
		"qr_link_id": id,
	})

	qrLink, err := uc.qrRepo.GetByID(ctx, id)
	if err != nil {
		uc.log.Error("Failed to get QR link", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get QR link: %w", err)
	}

	return uc.buildQRLinkResponse(qrLink), nil
}

func (uc *qrUseCase) GetQRLinksByMerchant(ctx context.Context, merchantID uuid.UUID, pag *pagination.Pagination) (*entity.QRLinksListResponse, error) {
	uc.log.Info("Getting QR links by merchant", map[string]interface{}{
		"merchant_id": merchantID,
		"page":        pag.Page,
		"page_size":   pag.PerPage,
	})

	qrLinks, total, err := uc.qrRepo.GetByMerchant(ctx, merchantID, pag)
	if err != nil {
		uc.log.Error("Failed to get QR links by merchant", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get QR links: %w", err)
	}

	responses := make([]entity.QRLinkResponse, len(qrLinks))
	for i, qrLink := range qrLinks {
		responses[i] = *uc.buildQRLinkResponse(&qrLink)
	}

	return &entity.QRLinksListResponse{
		QRLinks: responses,
		Total:   total,
		Page:    pag.Page,
		Limit:   pag.PerPage,
	}, nil
}

func (uc *qrUseCase) GetQRLinksByUser(ctx context.Context, userID uuid.UUID, pag *pagination.Pagination) (*entity.QRLinksListResponse, error) {
	uc.log.Info("Getting QR links by user", map[string]interface{}{
		"user_id":   userID,
		"page":      pag.Page,
		"page_size": pag.PerPage,
	})

	qrLinks, total, err := uc.qrRepo.GetByUser(ctx, userID, pag)
	if err != nil {
		uc.log.Error("Failed to get QR links by user", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get QR links: %w", err)
	}

	responses := make([]entity.QRLinkResponse, len(qrLinks))
	for i, qrLink := range qrLinks {
		responses[i] = *uc.buildQRLinkResponse(&qrLink)
	}

	return &entity.QRLinksListResponse{
		QRLinks: responses,
		Total:   total,
		Page:    pag.Page,
		Limit:   pag.PerPage,
	}, nil
}

func (uc *qrUseCase) UpdateQRLink(ctx context.Context, id, userID uuid.UUID, req *entity.UpdateQRLinkRequest) (*entity.QRLinkResponse, error) {
	uc.log.Info("Updating QR link", map[string]interface{}{
		"qr_link_id": id,
		"user_id":    userID,
	})

	updatedQRLink, err := uc.qrRepo.Update(ctx, id, userID, req)
	if err != nil {
		uc.log.Error("Failed to update QR link", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to update QR link: %w", err)
	}

	uc.log.Info("QR link updated successfully", map[string]interface{}{
		"qr_link_id": id,
	})

	return uc.buildQRLinkResponse(updatedQRLink), nil
}

func (uc *qrUseCase) DeleteQRLink(ctx context.Context, id, userID uuid.UUID) error {
	uc.log.Info("Deleting QR link", map[string]interface{}{
		"qr_link_id": id,
		"user_id":    userID,
	})

	if err := uc.qrRepo.Delete(ctx, id, userID); err != nil {
		uc.log.Error("Failed to delete QR link", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to delete QR link: %w", err)
	}

	uc.log.Info("QR link deleted successfully", map[string]interface{}{
		"qr_link_id": id,
	})

	return nil
}

func (uc *qrUseCase) ProcessQRPayment(ctx context.Context, qrLinkID uuid.UUID, req *entity.QRPaymentRequest) (*entity.QRPaymentResponse, error) {
	uc.log.Info("Processing QR payment", map[string]interface{}{
		"qr_link_id": qrLinkID,
		"medium":     req.Medium,
	})

	if err := req.Validate(); err != nil {
		uc.log.Error("QR payment request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get QR link details
	qrLink, err := uc.qrRepo.GetByID(ctx, qrLinkID)
	if err != nil {
		uc.log.Error("Failed to get QR link for payment", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("QR link not found: %w", err)
	}

	if !qrLink.IsActive {
		return nil, fmt.Errorf("QR link is not active")
	}

	// Determine payment amount
	var paymentAmount float64
	if qrLink.Type == entity.STATIC {
		if qrLink.Amount == nil {
			return nil, fmt.Errorf("static QR link has no amount set")
		}
		paymentAmount = *qrLink.Amount
	} else {
		if req.Amount == nil {
			return nil, fmt.Errorf("amount is required for dynamic QR link")
		}
		paymentAmount = *req.Amount
	}

	// Validate payment medium is supported
	mediumSupported := false
	for _, supportedMedium := range qrLink.SupportedMethods {
		if supportedMedium == req.Medium {
			mediumSupported = true
			break
		}
	}
	if !mediumSupported {
		return nil, fmt.Errorf("payment medium %s is not supported by this QR link", req.Medium)
	}

	// Determine transaction tag based on QR link tag
	var transactionTag string
	switch qrLink.Tag {
	case entity.SHOP:
		transactionTag = txEntity.QR_SHOP_PAYMENT
	case entity.RESTAURANT:
		transactionTag = txEntity.QR_RESTAURANT_PAYMENT
	case entity.DONATION:
		transactionTag = txEntity.QR_DONATION_PAYMENT
	default:
		transactionTag = txEntity.QR_SHOP_PAYMENT // Default fallback
	}

	// Generate callback URL
	baseURL := os.Getenv("APP_URL_V2")
	if baseURL == "" {
		baseURL = "http://localhost:8080" // fallback
	}
	callbackURL := fmt.Sprintf("%s/api/v2/qr/callback", baseURL)

	// Safely handle tip amount - make it optional
	var tipAmount float64
	if req.TipAmount != nil {
		tipAmount = *req.TipAmount
	}

	totalAmountIncludingTip := paymentAmount + tipAmount

	hasTip := qrLink.IsTipEnabled && req.TipAmount != nil && *req.TipAmount > 0
	// Create main payment transaction
	mainTx := &txEntity.Transaction{
		Id:                uuid.New(),
		PhoneNumber:       req.PhoneNumber,
		UserId:            qrLink.UserID,
		MerchantId:        qrLink.MerchantID,
		Type:              txEntity.DEPOSIT,
		Medium:            req.Medium,
		Amount:            totalAmountIncludingTip,
		Currency:          "ETB", // Default currency
		Reference:         fmt.Sprintf("QR_%s", qrLinkID.String()[:8]),
		Description:       uc.buildPaymentDescription(qrLink),
		CallbackURL:       callbackURL,
		Test:              false,
		Status:            txEntity.INITIATED,
		TransactionSource: txEntity.QR_PAYMENT,
		QRTag:             &transactionTag,
		HasTip:            hasTip,
		// QR specific context
		QRLinkID: &qrLinkID,
	}

	// Calculate fees and totals (same as direct payment)
	commission := 2.75
	if qrLink.MerchantID == uuid.MustParse("66e3c8c4-5308-4537-955e-c8d6cd1b4afe") {
		commission = 1
	}
	mainTx.FeeAmount = RoundToTwoDecimals((mainTx.Amount * commission) / 100) // 2.75% transaction fee
	mainTx.VatAmount = RoundToTwoDecimals((mainTx.FeeAmount * 15) / 100)      // 15% VAT on fee
	mainTx.TotalAmount = RoundToTwoDecimals(mainTx.Amount + mainTx.FeeAmount + mainTx.VatAmount)
	mainTx.AdminNet = RoundToTwoDecimals(mainTx.FeeAmount - mainTx.VatAmount)
	mainTx.MerchantNet = RoundToTwoDecimals(mainTx.Amount)

	uc.log.Info("Created QR payment transaction", map[string]interface{}{
		"transaction_id":  mainTx.Id,
		"original_amount": mainTx.Amount,
		"total_amount":    mainTx.TotalAmount,
		"qr_link_id":      qrLinkID,
		"tag":             transactionTag,
	})

	if hasTip {
		mainTx.HasTip = true
		mainTx.TipAmount = req.TipAmount
		// Safely handle optional tipee phone
		if req.TipeePhone != nil {
			mainTx.TipeePhone = req.TipeePhone
		}
		// Safely handle optional tip medium
		if req.TipMedium != nil {
			mainTx.TipMedium = (*string)(req.TipMedium)
		}
	}

	// Store main transaction
	if err := uc.transactionRepo.Create(ctx, mainTx); err != nil {
		uc.log.Error("Failed to create QR payment transaction", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Process main payment
	paymentReq := &payment.PaymentRequest{
		TransactionID: mainTx.Id,
		Medium:        req.Medium,
		Amount:        mainTx.TotalAmount,
		Currency:      mainTx.Currency,
		PhoneNumber:   req.PhoneNumber,
		Reference:     mainTx.Reference,
		Description:   mainTx.Description,
		CallbackURL:   callbackURL,
	}

	uc.log.Info("Processing main QR payment", map[string]interface{}{
		"transaction_id": mainTx.Id,
		"amount":         mainTx.TotalAmount,
	})

	paymentResp, err := uc.paymentService.ProcessPayment(ctx, qrLink.MerchantID.String(), paymentReq)
	if err != nil {
		uc.log.Error("QR payment processing failed", map[string]interface{}{
			"error": err.Error(),
		})
		mainTx.Status = txEntity.FAILED
		_ = uc.transactionRepo.Update(ctx, mainTx)
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Update main transaction status
	mainTx.Status = txEntity.TransactionStatus(paymentResp.Status)
	if err := uc.transactionRepo.Update(ctx, mainTx); err != nil {
		uc.log.Error("Failed to update QR transaction status", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	response := &entity.QRPaymentResponse{
		Success:                paymentResp.Success,
		Status:                 string(paymentResp.Status),
		Message:                paymentResp.Message,
		PaymentAmount:          totalAmountIncludingTip,
		SocialPayTransactionID: mainTx.Id.String(),
	}

	// Add tip information if present
	if hasTip {
		response.TipAmount = &tipAmount
	}

	uc.log.Info("QR payment processed successfully", map[string]interface{}{
		"qr_link_id":     qrLinkID,
		"payment_amount": totalAmountIncludingTip,
		"tip_amount":     tipAmount,
		"has_tip":        hasTip,
		"transaction_id": mainTx.Id,
		"status":         response.Status,
	})

	return response, nil
}

func (uc *qrUseCase) buildQRLinkResponse(qrLink *entity.QRLink) *entity.QRLinkResponse {
	return &entity.QRLinkResponse{
		QRLink:     qrLink,
		QRCodeURL:  fmt.Sprintf("https://api.Socialpay.co/qr/display/%s", qrLink.ID),
		PaymentURL: fmt.Sprintf("https://checkout.Socialpay.co/qr/%s", qrLink.ID),
	}
}

// Helper function to build payment description
func (uc *qrUseCase) buildPaymentDescription(qrLink *entity.QRLink) string {
	if qrLink.Description != nil && *qrLink.Description != "" {
		return *qrLink.Description
	}

	if qrLink.Title != nil && *qrLink.Title != "" {
		return *qrLink.Title
	}

	switch qrLink.Tag {
	case entity.SHOP:
		return "QR Shop Payment"
	case entity.RESTAURANT:
		return "QR Restaurant Payment"
	case entity.DONATION:
		return "QR Donation Payment"
	default:
		return "QR Payment"
	}
}

// Helper function to round to two decimals (same as socialpayapi)
func RoundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
