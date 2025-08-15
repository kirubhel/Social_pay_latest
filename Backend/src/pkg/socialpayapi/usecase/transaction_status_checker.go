package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/filter"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	txRepo "github.com/socialpay/socialpay/src/pkg/transaction/core/repository"
	settlementdto "github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
)

// WebhookDispatcher is an interface to avoid import cycles
type WebhookDispatcher interface {
	HandleWebhookDispatch(ctx context.Context, req settlementdto.WebhookRequest) error
}

type TransactionStatusChecker struct {
	transactionRepo   txRepo.TransactionRepository
	paymentService    PaymentProcessor
	webhookDispatcher WebhookDispatcher
	log               logging.Logger
}

func NewTransactionStatusChecker(
	transactionRepo txRepo.TransactionRepository,
	paymentService PaymentProcessor,
	webhookDispatcher WebhookDispatcher,
) *TransactionStatusChecker {
	return &TransactionStatusChecker{
		transactionRepo:   transactionRepo,
		paymentService:    paymentService,
		webhookDispatcher: webhookDispatcher,
		log:               logging.NewStdLogger("[TRANSACTION-STATUS-CHECKER]"),
	}
}

func (tsc *TransactionStatusChecker) CheckPendingCBETransactions(ctx context.Context) error {
	tsc.log.Info("Starting to check pending CBE transactions", map[string]interface{}{})

	// Create filter for pending CBE transactions older than 2 minutes
	twoMinutesAgo := time.Now().Add(-5 * time.Minute)

	filterParam := filter.Filter{
		Pagination: pagination.Pagination{
			Page:    1,
			PerPage: 100, // Limit to 100 per run to avoid overloading
		},
		Sort: []filter.Sort{
			{
				Field:    "created_at",
				Operator: "ASC",
			},
		},
		Group: filter.FilterGroup{
			Linker: "AND",
			Fields: []filter.FilterItem{
				filter.Field{
					Name:     "status",
					Operator: "=",
					Value:    string(txEntity.PENDING),
				},
				filter.Field{
					Name:     "medium",
					Operator: "=",
					Value:    string(txEntity.CBE),
				},
				filter.Field{
					Name:     "created_at",
					Operator: "<",
					Value:    twoMinutesAgo,
				},
			},
		},
	}

	// Get pending CBE transactions using the filter
	// Pass nil as userID since we want to query for all users
	transactions, err := tsc.transactionRepo.GetTransactionsByParameters(ctx, filterParam, uuid.Nil)
	if err != nil {
		tsc.log.Error("Failed to get pending CBE transactions", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to get pending CBE transactions: %w", err)
	}

	tsc.log.Info("Found pending CBE transactions", map[string]interface{}{
		"count": len(transactions),
	})

	processedCount := 0
	updatedCount := 0

	// Process each transaction
	for _, tx := range transactions {
		processedCount++

		tsc.log.Info("Checking status for transaction", map[string]interface{}{
			"transaction_id": tx.Id,
			"provider_tx_id": tx.ProviderTxId,
			"created_at":     tx.CreatedAt,
		})

		// Determine the transaction ID to use for status query
		queryID := tx.ProviderTxId
		if queryID == "" {
			// If ProviderTxId is empty, use the transaction ID (for older transactions)
			queryID = tx.Id.String()
		}

		// Query transaction status from CBE
		queryResp, err := tsc.paymentService.QueryTransactionStatus(ctx, txEntity.CBE, queryID)
		if err != nil {
			tsc.log.Error("Failed to query transaction status", map[string]interface{}{
				"transaction_id": tx.Id,
				"provider_tx_id": queryID,
				"error":          err.Error(),
			})
			// Continue with next transaction even if one fails
			continue
		}

		if queryResp == nil {
			tsc.log.Warn("Received nil response from status query", map[string]interface{}{
				"transaction_id": tx.Id,
				"provider_tx_id": queryID,
			})
			continue
		}

		// Check if status has changed
		if queryResp.Status != tx.Status {
			tsc.log.Info("Transaction status changed, updating", map[string]interface{}{
				"transaction_id": tx.Id,
				"old_status":     tx.Status,
				"new_status":     queryResp.Status,
				"provider_tx_id": queryResp.ProviderTxId,
			})

			// Dispatch webhook for status change
			if err := tsc.dispatchWebhookSettlement(ctx, tx.Id, queryResp); err != nil {
				tsc.log.Error("Failed to dispatch webhook", map[string]interface{}{
					"transaction_id": tx.Id,
					"error":          err.Error(),
				})
				// Don't fail the whole process if webhook dispatch fails
			}

			updatedCount++
			tsc.log.Info("Successfully updated transaction status", map[string]interface{}{
				"transaction_id": tx.Id,
				"new_status":     queryResp.Status,
			})
		} else {
			tsc.log.Debug("Transaction status unchanged", map[string]interface{}{
				"transaction_id": tx.Id,
				"status":         tx.Status,
			})
		}
	}

	tsc.log.Info("Completed checking pending CBE transactions", map[string]interface{}{
		"total_processed": processedCount,
		"total_updated":   updatedCount,
	})

	return nil
}

func (tsc *TransactionStatusChecker) dispatchWebhookSettlement(ctx context.Context, transactionID uuid.UUID, transactionStatusQueryResponse *payment.TransactionStatusQueryResponse) error {
	providerData, _ := json.Marshal(transactionStatusQueryResponse.ProviderData)

	// Prepare webhook request
	webhookReq := settlementdto.WebhookRequest{
		TransactionID: transactionID.String(),
		Status:        string(transactionStatusQueryResponse.Status),
		Message:       "Transaction status updated by cron job status checker",
		ProviderTxID:  transactionStatusQueryResponse.ProviderTxId,
		ProviderData:  string(providerData),
		Timestamp:     time.Now(),
	}

	tsc.log.Info("Dispatching webhook from cron job", map[string]interface{}{
		"webhook_request": webhookReq,
		"transaction_id":  transactionID.String(),
	})

	// Dispatch webhook
	if err := tsc.webhookDispatcher.HandleWebhookDispatch(ctx, webhookReq); err != nil {
		tsc.log.Error("Failed to dispatch webhook from cron job", map[string]interface{}{
			"error":          err.Error(),
			"transaction_id": transactionID.String(),
		})
		return fmt.Errorf("failed to dispatch webhook: %w", err)
	}

	tsc.log.Info("Successfully dispatched webhook from cron job", map[string]interface{}{
		"transaction_id": transactionID.String(),
		"status":         transactionStatusQueryResponse.Status,
	})

	return nil
}
