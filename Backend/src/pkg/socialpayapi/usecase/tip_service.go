package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	txRepository "github.com/socialpay/socialpay/src/pkg/transaction/core/repository"
	walletUseCase "github.com/socialpay/socialpay/src/pkg/wallet/usecase"
)

// TipProcessingService handles automatic tip processing for QR payments
type TipProcessingService interface {
	// ProcessPendingTips processes all transactions with pending tips
	ProcessPendingTips(ctx context.Context) error

	// ProcessTipForTransaction processes a tip for a specific transaction
	ProcessTipForTransaction(ctx context.Context, transactionID uuid.UUID) error

	// CreateTipWithdrawal creates a withdrawal transaction for processed tips
	CreateTipWithdrawal(ctx context.Context, mainTx *txEntity.Transaction, tipeePhone string, tipAmount float64, medium txEntity.TransactionMedium) (*txEntity.Transaction, error)
}

type tipProcessingService struct {
	transactionRepo txRepository.TransactionRepository
	walletUseCase   walletUseCase.WalletUseCase
	paymentService  PaymentProcessor
	log             logging.Logger
}

// NewTipProcessingService creates a new tip processing service
func NewTipProcessingService(
	transactionRepo txRepository.TransactionRepository,
	walletUseCase walletUseCase.WalletUseCase,
	paymentService PaymentProcessor,
) TipProcessingService {
	return &tipProcessingService{
		transactionRepo: transactionRepo,
		walletUseCase:   walletUseCase,
		paymentService:  paymentService,
		log:             logging.NewStdLogger("[TIP-PROCESSING]"),
	}
}

func (s *tipProcessingService) ProcessPendingTips(ctx context.Context) error {
	s.log.Info("Starting to process pending tips", nil)

	// Get all transactions with pending tips
	transactions, err := s.transactionRepo.GetTransactionsWithPendingTips(ctx)
	if err != nil {
		s.log.Error("Failed to get transactions with pending tips", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to get pending tips: %w", err)
	}

	s.log.Info("Found transactions with pending tips", map[string]interface{}{
		"count": len(transactions),
	})

	// Process each transaction
	for _, tx := range transactions {
		if err := s.ProcessTipForTransaction(ctx, tx.Id); err != nil {
			s.log.Error("Failed to process tip for transaction", map[string]interface{}{
				"transaction_id": tx.Id,
				"error":          err.Error(),
			})
			// Continue processing other tips even if one fails
			continue
		}
	}

	s.log.Info("Completed processing pending tips", map[string]interface{}{
		"processed_count": len(transactions),
	})

	return nil
}

func (s *tipProcessingService) ProcessTipForTransaction(ctx context.Context, transactionID uuid.UUID) error {
	s.log.Info("Processing tip for transaction", map[string]interface{}{
		"transaction_id": transactionID,
	})

	// Get the transaction details
	tx, err := s.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	// Validate transaction has tip and is not processed
	if !tx.HasTip || tx.TipProcessed {
		s.log.Info("Transaction tip already processed or has no tip", map[string]interface{}{
			"transaction_id": transactionID,
			"has_tip":        tx.HasTip,
			"tip_processed":  tx.TipProcessed,
		})
		return nil
	}

	// Validate tip data
	if tx.TipAmount == nil || tx.TipeePhone == nil || tx.TipMedium == nil {
		return fmt.Errorf("incomplete tip data for transaction %s", transactionID)
	}

	s.log.Info("Creating tip withdrawal", map[string]interface{}{
		"transaction_id": transactionID,
		"tip_amount":     *tx.TipAmount,
		"tipee_phone":    *tx.TipeePhone,
		"tip_medium":     *tx.TipMedium,
	})

	// Create tip withdrawal transaction
	tipTransaction, err := s.CreateTipWithdrawal(ctx, tx, *tx.TipeePhone, *tx.TipAmount, txEntity.TransactionMedium(*tx.TipMedium))
	if err != nil {
		return fmt.Errorf("failed to create tip withdrawal: %w", err)
	}

	// Update the original transaction to mark tip as processed
	if err := s.transactionRepo.UpdateTipProcessing(ctx, transactionID, tipTransaction.Id); err != nil {
		s.log.Error("Failed to update tip processing status", map[string]interface{}{
			"transaction_id":     transactionID,
			"tip_transaction_id": tipTransaction.Id,
			"error":              err.Error(),
		})
		return fmt.Errorf("failed to update tip processing status: %w", err)
	}

	s.log.Info("Successfully processed tip for transaction", map[string]interface{}{
		"transaction_id":     transactionID,
		"tip_transaction_id": tipTransaction.Id,
		"tip_amount":         *tx.TipAmount,
	})

	return nil
}

func (s *tipProcessingService) CreateTipWithdrawal(ctx context.Context, mainTx *txEntity.Transaction, tipeePhone string, tipAmount float64, medium txEntity.TransactionMedium) (*txEntity.Transaction, error) {
	s.log.Info("Creating tip withdrawal transaction", map[string]interface{}{
		"tipee_phone": tipeePhone,
		"tip_amount":  tipAmount,
		"medium":      medium,
		"merchant_id": mainTx.MerchantId,
	})

	// Create withdrawal transaction for the tip
	tipTx := &txEntity.Transaction{
		Id:                uuid.New(),
		PhoneNumber:       tipeePhone,
		UserId:            mainTx.UserId,
		MerchantId:        mainTx.MerchantId,
		Type:              txEntity.WITHDRAWAL,
		CallbackURL:       mainTx.CallbackURL,
		Medium:            medium,
		Currency:          "ETB", // Default currency
		Reference:         fmt.Sprintf("TIP-%s", mainTx.Id.String()),
		Comment:           fmt.Sprintf("Tip withdrawal %s for amount  %.2f", mainTx.Id.String(), tipAmount),
		Description:       fmt.Sprintf("Tip for QR payment %s", mainTx.Id.String()),
		TransactionSource: txEntity.WITHDRAWAL_TIP,
		Status:            txEntity.INITIATED,
		Amount:            tipAmount, // TODO: add fee and vat
		TotalAmount:       tipAmount, // TODO: add fee and vat
		MerchantNet:       tipAmount, // TODO: add fee and vat
		Test:              false,
		Verified:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Store the tip withdrawal transaction
	if err := s.transactionRepo.CreateWithContext(ctx, tipTx); err != nil {
		s.log.Error("Failed to create tip withdrawal transaction", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create tip withdrawal transaction: %w", err)
	}

	s.log.Info("Created tip withdrawal transaction", map[string]interface{}{
		"tip_transaction_id": tipTx.Id,
		"tipee_phone":        tipeePhone,
		"amount":             tipAmount,
	})

	// Process the withdrawal (this would typically be done asynchronously)
	go func() {
		s.processTipWithdrawalAsync(context.Background(), tipTx)
	}()

	return tipTx, nil
}

// processTipWithdrawalAsync processes the tip withdrawal asynchronously
func (s *tipProcessingService) processTipWithdrawalAsync(ctx context.Context, tipTx *txEntity.Transaction) {
	s.log.Info("Processing tip withdrawal asynchronously", map[string]interface{}{
		"tip_transaction_id": tipTx.Id,
	})

	// Create payment request for withdrawal
	paymentReq := &payment.PaymentRequest{
		TransactionID: tipTx.Id,
		Medium:        tipTx.Medium,
		Amount:        tipTx.Amount,
		Currency:      tipTx.Currency,
		PhoneNumber:   tipTx.PhoneNumber,
		Reference:     tipTx.Reference,
		Description:   tipTx.Description,
	}

	// Process withdrawal using payment service
	// Note: In a real implementation, you'd need proper API keys and authentication
	_, err := s.paymentService.ProcessWithdrawal(ctx, tipTx.MerchantId.String(), paymentReq)

	// Update transaction status based on result
	if err != nil {
		tipTx.Status = txEntity.FAILED
		tipTx.Comment = fmt.Sprintf("Tip withdrawal failed: %s", err.Error())
		s.log.Error("Tip withdrawal failed", map[string]interface{}{
			"tip_transaction_id": tipTx.Id,
			"error":              err.Error(),
		})
	} else {
		tipTx.Status = txEntity.PENDING // Will be updated by webhook when completed
		s.log.Info("Tip withdrawal initiated successfully", map[string]interface{}{
			"tip_transaction_id": tipTx.Id,
		})
	}

	// Update transaction status
	if updateErr := s.transactionRepo.Update(ctx, tipTx); updateErr != nil {
		s.log.Error("Failed to update tip withdrawal status", map[string]interface{}{
			"tip_transaction_id": tipTx.Id,
			"error":              updateErr.Error(),
		})
	}
}
