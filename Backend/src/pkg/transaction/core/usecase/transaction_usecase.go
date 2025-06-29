package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	"github.com/socialpay/socialpay/src/pkg/webhook/adapter/gateway/kafka/producer"
)

type TransactionRepository interface {
	GetByID(ctx context.Context, txnID uuid.UUID) (*entity.Transaction, error)
	UpdateStatus(ctx context.Context, txnID uuid.UUID, status string) error
}

type TransactionUseCase interface {
	OverrideTransactionStatus(ctx context.Context, txnID uuid.UUID, newStatus string, reason string, adminID string) error
}

type TransactionUseCaseImpl struct {
	transactionRepo TransactionRepository
	sendProducer    *producer.GroupedProducer
	log             logging.Logger
}

func NewTransactionUseCase(transactionRepo TransactionRepository, sendProducer *producer.GroupedProducer) TransactionUseCase {
	return &TransactionUseCaseImpl{
		transactionRepo: transactionRepo,
		sendProducer:    sendProducer,
		log:             logging.NewStdLogger("[transaction]"),
	}
}

func (uc *TransactionUseCaseImpl) OverrideTransactionStatus(ctx context.Context, txnID uuid.UUID, newStatus string, reason string, adminID string) error {
	uc.log.Info("overriding transaction status", map[string]interface{}{
		"txnID":     txnID,
		"newStatus": newStatus,
		"reason":    reason,
		"adminID":   adminID,
	})

	txn, err := uc.transactionRepo.GetByID(ctx, txnID)
	if err != nil {
		uc.log.Error("failed to get transaction", map[string]interface{}{
			"error": err,
			"txnID": txnID,
		})
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if txn.Status == "SUCCESS" || txn.Status == "FAILED" {
		uc.log.Error("transaction is already finalized", map[string]interface{}{
			"txnID":  txnID,
			"status": txn.Status,
		})
		return fmt.Errorf("transaction is already finalized")
	}

	oldStatus := txn.Status
	if err := uc.transactionRepo.UpdateStatus(ctx, txnID, newStatus); err != nil {
		uc.log.Error("failed to update transaction status", map[string]interface{}{
			"error":  err,
			"txnID":  txnID,
			"status": newStatus,
		})
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	uc.log.Info("transaction status updated successfully", map[string]interface{}{
		"txnID":     txnID,
		"oldStatus": oldStatus,
		"newStatus": newStatus,
	})

	// Trigger webhook with updated status
	event := map[string]interface{}{
		"transaction_id": txnID.String(),
		"status":         newStatus,
		"reason":         reason,
		"admin_id":       adminID,
		"timestamp":      time.Now(),
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		uc.log.Error("failed to marshal webhook event", map[string]interface{}{
			"error": err,
			"event": event,
		})
		return fmt.Errorf("failed to marshal webhook event: %w", err)
	}

	uc.sendProducer.Produce(txn.MerchantId.String(), bytes)

	uc.log.Info("webhook event produced to Kafka", map[string]interface{}{
		"merchantID": txn.MerchantId.String(),
		"eventSize":  len(bytes),
	})

	return nil
}
