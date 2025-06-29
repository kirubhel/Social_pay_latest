package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/config"

	"github.com/google/uuid"
	tipService "github.com/socialpay/socialpay/src/pkg/socialpayapi/usecase"
	notificationUsecase "github.com/socialpay/socialpay/src/pkg/notifications/usecase"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	transactionRepo "github.com/socialpay/socialpay/src/pkg/transaction/core/repository"
	walletUsecase "github.com/socialpay/socialpay/src/pkg/wallet/usecase"
	webhookDto "github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
	"github.com/socialpay/socialpay/src/pkg/webhook/adapter/gateway/kafka/producer"
	webhookRepo "github.com/socialpay/socialpay/src/pkg/webhook/adapter/gateway/repository"
	webhook "github.com/socialpay/socialpay/src/pkg/webhook/core/entity"
)

type WebhookUseCaseImpl struct {
	transactionRepo     transactionRepo.TransactionRepository
	callbackRepo        webhookRepo.CallbackRepository
	walletUsecase       walletUsecase.WalletUseCase
	log                 logging.Logger
	producer            *producer.GroupedProducer
	sendProducer        *producer.GroupedProducer
	tipService          tipService.TipProcessingService
	transactionNotifier *notificationUsecase.TransactionNotifier
}

func NewWebhookUseCase(
	cfg *config.Config,
	transactionRepo transactionRepo.TransactionRepository,
	callbackRepo webhookRepo.CallbackRepository,
	walletUsecase walletUsecase.WalletUseCase,
	tipService tipService.TipProcessingService,
	transactionNotifier *notificationUsecase.TransactionNotifier,
) WebhookUseCase {
	log := logging.NewStdLogger("[webhook]")
	log.Info("initializing webhook use case", map[string]interface{}{
		"kafka_brokers": cfg.Kafka.Brokers,
		"kafka_topic":   cfg.Kafka.Topics.WebhookDispatch,
	})

	dispatchProducer := producer.NewGroupedProducer(cfg.Kafka.Brokers, cfg.Kafka.Topics.WebhookDispatch, 5, logging.NewStdLogger("[webhook][WebhookDispatch]"))
	dispatchProducer.Start() // Start the producer workers

	sendProducer := producer.NewGroupedProducer(cfg.Kafka.Brokers, cfg.Kafka.Topics.WebhookSend, 5, logging.NewStdLogger("[webhook][WebhookSend]"))
	sendProducer.Start() // Start the send producer workers

	log.Info("webhook producers and notification service started", nil)

	return &WebhookUseCaseImpl{
		transactionRepo:     transactionRepo,
		callbackRepo:        callbackRepo,
		walletUsecase:       walletUsecase,
		log:                 log,
		producer:            dispatchProducer,
		sendProducer:        sendProducer,
		tipService:          tipService,
		transactionNotifier: transactionNotifier,
	}
}

func (uc *WebhookUseCaseImpl) GetProducer() *producer.GroupedProducer {
	return uc.producer
}

func (uc *WebhookUseCaseImpl) GetSendProducer() *producer.GroupedProducer {
	return uc.sendProducer
}

func (uc *WebhookUseCaseImpl) ProcessTransactionStatus(ctx context.Context, txnID uuid.UUID, status txEntity.TransactionStatus) error {
	uc.log.Info("processing transaction status", map[string]interface{}{
		"txnID":  txnID,
		"status": status,
	})

	txn, err := uc.transactionRepo.GetByID(ctx, txnID)
	if err != nil {
		uc.log.Error("failed to get transaction", map[string]interface{}{
			"error": err,
			"txnID": txnID,
		})
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	uc.log.Info("transaction found", map[string]interface{}{
		"txnID":        txnID,
		"currentState": txn.Status,
		"targetState":  status,
	})

	if !isValidStatusTransition(txn.Status, status) {
		uc.log.Error("invalid status transition", map[string]interface{}{
			"from":  txn.Status,
			"to":    status,
			"txnID": txnID,
		})
		return fmt.Errorf("invalid status transition from %s to %s", txn.Status, status)
	}

	if err := uc.transactionRepo.UpdateStatus(ctx, txnID, status); err != nil {
		uc.log.Error("failed to update transaction status", map[string]interface{}{
			"error":  err,
			"txnID":  txnID,
			"status": status,
		})
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	uc.log.Info("transaction status updated successfully", map[string]interface{}{
		"txnID":     txnID,
		"newStatus": status,
	})

	return nil
}

func (uc *WebhookUseCaseImpl) HandlePaymentStatusUpdate(ctx context.Context, msg webhookDto.WebhookMessage) error {
	uc.log.Info("handling payment status update", map[string]interface{}{
		"type":          msg.Type,
		"transactionID": msg.TransactionID,
		"status":        msg.Status,
		"amount":        msg.TotalAmount,
		"merchantID":    msg.MerchantID,
	})

	parsedTxnID, err := uuid.Parse(msg.TransactionID)
	if err != nil {
		uc.log.Error("failed to parse transaction ID", map[string]interface{}{
			"error":         err,
			"transactionID": msg.TransactionID,
		})
		return fmt.Errorf("failed to parse transaction ID: %w", err)
	}

	txn, err := uc.transactionRepo.GetByID(ctx, parsedTxnID)
	if err != nil {
		uc.log.Error("failed to get transaction", map[string]interface{}{
			"type":          msg.Type,
			"error":         err,
			"transactionID": msg.TransactionID,
		})
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	uc.log.Info("transaction found", map[string]interface{}{
		"type":          txn.Type,
		"currentStatus": txn.Status,
		"transactionID": msg.TransactionID,
	})

	// Update transaction status
	txnStatus := txEntity.TransactionStatus(msg.Status)
	if err := uc.transactionRepo.UpdateStatus(ctx, parsedTxnID, txnStatus); err != nil {
		uc.log.Error("failed to update transaction status", map[string]interface{}{
			"error":         err,
			"transactionID": msg.TransactionID,
			"status":        msg.Status,
		})
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	uc.log.Info("transaction status updated successfully", map[string]interface{}{
		"type":      txn.Type,
		"newStatus": txnStatus,
		"txnID":     msg.TransactionID,
	})

	// Send SMS notifications for transaction status updates
	uc.log.Info("sending transaction status notifications", map[string]interface{}{
		"transactionID": msg.TransactionID,
		"status":        txnStatus,
	})

	if uc.transactionNotifier != nil {
		if err := uc.transactionNotifier.NotifyTransactionStatus(ctx, txn, string(txnStatus)); err != nil {
			uc.log.Error("failed to send transaction notifications", map[string]interface{}{
				"error":         err,
				"transactionID": msg.TransactionID,
			})
			// Don't return error here - notification failure shouldn't block transaction processing
		} else {
			uc.log.Info("transaction notifications sent successfully", map[string]interface{}{
				"transactionID": msg.TransactionID,
				"status":        txnStatus,
			})
		}
	} else {
		uc.log.Info("notification service not available, skipping SMS notifications", map[string]interface{}{
			"transactionID": msg.TransactionID,
		})
	}

	// Parse merchant ID
	merchantID, err := uuid.Parse(msg.MerchantID)
	if err != nil {
		uc.log.Error("invalid merchant ID", map[string]interface{}{
			"error":      err,
			"merchantID": msg.MerchantID,
		})
		return fmt.Errorf("invalid merchant ID: %w", err)
	}
	uc.log.Info("Processing Wallet", map[string]interface{}{
		"merchantID": merchantID,
		"type":       txn.Type,
		"status":     txnStatus,
	})
	if txn.Type == txEntity.WITHDRAWAL {
		// Process withdrawal transaction using transaction-safe methods
		isSuccess := txnStatus == txEntity.SUCCESS
		if err := uc.walletUsecase.ProcessTransactionStatus(ctx, merchantID, msg.TotalAmount, isSuccess, true); err != nil {
			uc.log.Error("failed to process withdrawal status", map[string]interface{}{
				"error":      err,
				"merchantID": msg.MerchantID,
				"amount":     msg.TotalAmount,
				"status":     txnStatus,
			})
			return fmt.Errorf("failed to process withdrawal status: %w", err)
		}
	} else if txnStatus == txEntity.SUCCESS {
		// Process deposit transaction using transaction-safe methods
		if err := uc.walletUsecase.ProcessTransactionStatus(ctx, merchantID, txn.MerchantNet, true, false); err != nil {
			uc.log.Error("failed to process deposit status", map[string]interface{}{
				"error":      err,
				"merchantID": msg.MerchantID,
				"amount":     txn.MerchantNet,
				"status":     txnStatus,
			})
			return fmt.Errorf("failed to process deposit status: %w", err)
		}
		uc.log.Info("Checking if transaction has tip", map[string]interface{}{
			"transactionID": txn.Id,
			"hasTip":        txn.HasTip,
		})
		if txn.HasTip && !txn.TipProcessed {
			uc.log.Info("Processing tip", map[string]interface{}{
				"transactionID": txn.Id,
			})
			uc.tipService.ProcessTipForTransaction(ctx, txn.Id)
		}
	}

	uc.log.Info("preparing webhook message for Kafka", map[string]interface{}{
		"callbackURL":   txn.CallbackURL,
		"transactionID": txn.Id,
		"merchantID":    txn.MerchantId,
	})

	// Create event for Kafka
	event := webhookDto.WebhookEventMerchant{
		Event:          txn.Type,
		SocialPayTxnID: txn.Id.String(),
		ReferenceId:    txn.Reference,
		Status:         string(txnStatus),
		Amount:         fmt.Sprintf("%f", txn.MerchantNet),
		CallbackURL:    txn.CallbackURL,
		Timestamp:      time.Now(),
		ProviderTxID:   msg.ProviderTxID,
		Message:        msg.Message,
		MerchantID:     msg.MerchantID,
		UserID:         msg.UserID,
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		uc.log.Error("failed to marshal webhook event", map[string]interface{}{
			"error": err,
			"event": event,
		})
		return fmt.Errorf("failed to marshal webhook event: %w", err)
	}

	// Send to Kafka using merchant ID as key for sequential processing
	uc.sendProducer.Produce(msg.MerchantID, bytes)

	uc.log.Info("webhook event produced to Kafka", map[string]interface{}{
		"merchantID": msg.MerchantID,
		"eventSize":  len(bytes),
	})

	return nil
}

func (uc *WebhookUseCaseImpl) HandleWebhookDispatch(ctx context.Context, req webhookDto.WebhookRequest) error {
	uc.log.Info("handling webhook dispatch", map[string]interface{}{
		"type":          req.Type,
		"transactionID": req.TransactionID,
		"status":        req.Status,
	})

	parsedTransactionID, err := uuid.Parse(req.TransactionID)
	if err != nil {
		uc.log.Error("failed to parse transaction ID", map[string]interface{}{
			"error":         err,
			"transactionID": req.TransactionID,
		})
		return fmt.Errorf("failed to parse transaction ID: %w", err)
	}

	txn, err := uc.transactionRepo.GetByID(ctx, parsedTransactionID)
	if err != nil {
		uc.log.Error("failed to get transaction", map[string]interface{}{
			"error":         err,
			"transactionID": req.TransactionID,
		})
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	uc.log.Info("transaction found", map[string]interface{}{
		"type":          txn.Type,
		"transactionID": req.TransactionID,
		"merchantID":    txn.MerchantId,
	})

	event := map[string]interface{}{
		"type":             txn.Type,
		"transaction_id":   req.TransactionID,
		"total_amount":     txn.TotalAmount,
		"status":           req.Status,
		"message":          req.Message,
		"provider_txid":    req.ProviderTxID,
		"provider_data":    req.ProviderData,
		"timestamp":        req.Timestamp,
		"user_id":          txn.UserId,
		"merchant_id":      txn.MerchantId,
		"callback_url":     txn.CallbackURL,
		"isHostedCheckout": req.IsHostedCheckout,
	}

	// check status
	if txn.Status != txEntity.PENDING && txn.Status != txEntity.INITIATED {
		uc.log.Info("transaction is already finalized: CLOSED", map[string]interface{}{
			"transactionID":   req.TransactionID,
			"EXISTING_STATUS": txn.Status,
			"NEW_STATUS":      req.Status,
		})
		return nil
	}

	uc.log.Info("preparing webhook event", map[string]interface{}{
		"event": event,
	})

	bytes, err := json.Marshal(event)
	if err != nil {
		uc.log.Error("failed to marshal event", map[string]interface{}{
			"error": err,
			"event": event,
		})
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	uc.log.Info("producing webhook message", map[string]interface{}{
		"merchantID": txn.MerchantId.String(),
		"eventSize":  len(bytes),
	})

	uc.producer.Produce(txn.MerchantId.String(), bytes)

	uc.log.Info("webhook dispatch completed successfully", map[string]interface{}{
		"transactionID": req.TransactionID,
	})

	return nil
}

func (uc *WebhookUseCaseImpl) CreateCallbackLog(ctx context.Context, txnID uuid.UUID, responseStatus int, requestBody string, responseBody string, merchantID string, userID string) error {
	uc.log.Info("creating callback log", map[string]interface{}{
		"txnID":          txnID,
		"responseStatus": responseStatus,
		"merchantID":     merchantID,
		"userID":         userID,
	})

	parsedMerchantID, err := uuid.Parse(merchantID)
	if err != nil {
		uc.log.Error("invalid merchant ID", map[string]interface{}{
			"error":      err,
			"merchantID": merchantID,
		})
		return fmt.Errorf("invalid merchant ID: %w", err)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		uc.log.Error("invalid user ID", map[string]interface{}{
			"error":  err,
			"userID": userID,
		})
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Parse the request body to get the message
	var webhookMsg webhookDto.WebhookRequest
	if err := json.Unmarshal([]byte(requestBody), &webhookMsg); err != nil {
		uc.log.Error("failed to parse webhook message", map[string]interface{}{
			"error": err,
			"body":  requestBody,
		})
		return fmt.Errorf("failed to parse webhook message: %w", err)
	}

	log := &webhook.CallbackLog{
		ID:           uuid.New(),
		TxnID:        txnID,
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		Status:       responseStatus,
		Message:      webhookMsg.Message,
		RetryCount:   0,
		MerchantID:   parsedMerchantID,
		UserID:       parsedUserID,
	}

	if err := uc.callbackRepo.Create(ctx, log); err != nil {
		uc.log.Error("failed to create callback log", map[string]interface{}{
			"error": err,
			"txnID": txnID,
		})
		return fmt.Errorf("failed to create callback log: %w", err)
	}

	uc.log.Info("callback log created successfully", map[string]interface{}{
		"logID": log.ID,
		"txnID": txnID,
	})

	return nil
}

func (uc *WebhookUseCaseImpl) UpdateCallbackLog(ctx context.Context, id uuid.UUID, responseBody string, responseStatus int) error {
	uc.log.Info("updating callback log", map[string]interface{}{
		"logID":          id,
		"responseStatus": responseStatus,
	})

	log, err := uc.callbackRepo.GetByID(ctx, id)
	if err != nil {
		uc.log.Error("failed to get callback log", map[string]interface{}{
			"error": err,
			"logID": id,
		})
		return fmt.Errorf("failed to get callback log: %w", err)
	}

	log.ResponseBody = responseBody
	log.Status = responseStatus

	if err := uc.callbackRepo.Update(ctx, log); err != nil {
		uc.log.Error("failed to update callback log", map[string]interface{}{
			"error": err,
			"logID": id,
		})
		return fmt.Errorf("failed to update callback log: %w", err)
	}

	uc.log.Info("callback log updated successfully", map[string]interface{}{
		"logID": id,
	})

	return nil
}

func (uc *WebhookUseCaseImpl) GetCallbackLogByID(ctx context.Context, id uuid.UUID) (*webhook.CallbackLog, error) {
	uc.log.Info("getting callback log by ID", map[string]interface{}{
		"logID": id,
	})

	log, err := uc.callbackRepo.GetByID(ctx, id)
	if err != nil {
		uc.log.Error("failed to get callback log", map[string]interface{}{
			"error": err,
			"logID": id,
		})
		return nil, fmt.Errorf("failed to get callback log: %w", err)
	}

	uc.log.Info("callback log retrieved successfully", map[string]interface{}{
		"logID": id,
	})

	return log, nil
}

func (uc *WebhookUseCaseImpl) GetCallbackLogsByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*webhook.CallbackLog, error) {
	uc.log.Info("getting callback logs by merchant ID", map[string]interface{}{
		"merchantID": merchantID,
	})

	logs, err := uc.callbackRepo.GetByMerchantID(ctx, merchantID)
	if err != nil {
		uc.log.Error("failed to get callback logs", map[string]interface{}{
			"error":      err,
			"merchantID": merchantID,
		})
		return nil, fmt.Errorf("failed to get callback logs: %w", err)
	}

	uc.log.Info("callback logs retrieved successfully", map[string]interface{}{
		"merchantID": merchantID,
		"count":      len(logs),
	})

	return logs, nil
}

func (uc *WebhookUseCaseImpl) GetAllCallbackLogs(ctx context.Context, pagination *txEntity.Pagination) ([]*webhook.CallbackLog, error) {
	uc.log.Info("getting all callback logs", map[string]interface{}{
		"action":    "get_all_callback_logs",
		"page":      pagination.Page,
		"page_size": pagination.PageSize,
	})

	if pagination == nil {
		uc.log.Error("pagination is nil", nil)
		return nil, fmt.Errorf("pagination parameters are required")
	}

	if err := pagination.Validate(); err != nil {
		uc.log.Error("pagination validation error", map[string]interface{}{
			"error":      err,
			"pagination": pagination,
		})
		return nil, fmt.Errorf("invalid pagination parameters: %w", err)
	}

	logs, err := uc.callbackRepo.GetAll(ctx, pagination)
	if err != nil {
		uc.log.Error("failed to get callback logs", map[string]interface{}{
			"error": err,
		})
		return nil, fmt.Errorf("failed to get callback logs: %w", err)
	}

	uc.log.Info("callback logs retrieved successfully", map[string]interface{}{
		"count":     len(logs),
		"page":      pagination.Page,
		"page_size": pagination.PageSize,
	})

	return logs, nil
}

func (uc *WebhookUseCaseImpl) OverrideTransactionStatus(ctx context.Context, txnID uuid.UUID, newStatus txEntity.TransactionStatus, reason string, adminID string) error {
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

	if txn.Status == txEntity.SUCCESS || txn.Status == txEntity.FAILED {
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

// Helper functions
func isValidStatusTransition(from, to txEntity.TransactionStatus) bool {
	validTransitions := map[txEntity.TransactionStatus][]txEntity.TransactionStatus{
		txEntity.PENDING: {
			txEntity.INITIATED,
			txEntity.EXPIRED,
			txEntity.FAILED,
		},
		txEntity.INITIATED: {
			txEntity.SUCCESS,
			txEntity.FAILED,
		},
		txEntity.SUCCESS: {},
		txEntity.FAILED:  {},
		txEntity.EXPIRED: {},
	}

	validNextStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, validState := range validNextStates {
		if validState == to {
			return true
		}
	}

	return false
}
