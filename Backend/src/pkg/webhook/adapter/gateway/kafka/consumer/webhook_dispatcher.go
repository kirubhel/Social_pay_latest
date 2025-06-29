package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/socialpay/socialpay/src/pkg/config"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/repository"
	transactionUsecase "github.com/socialpay/socialpay/src/pkg/transaction/usecase"
	webhookDto "github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
	webhookUsecase "github.com/socialpay/socialpay/src/pkg/webhook/usecase"
	"github.com/segmentio/kafka-go"
)

type WebhookDispatcherWorker struct {
	cfg                *config.Config
	db                 *sql.DB
	reader             *kafka.Reader
	client             *http.Client
	logger             logging.Logger
	usecase            webhookUsecase.WebhookUseCase
	transactionUsecase transactionUsecase.TransactionUseCase
	hostedPaymentRepo  repository.HostedPaymentRepository
}

func NewWebhookDispatcherWorker(cfg *config.Config, db *sql.DB, usecase webhookUsecase.WebhookUseCase,
	transactionUsecase transactionUsecase.TransactionUseCase,
	hostedPaymentRepo repository.HostedPaymentRepository,
) *WebhookDispatcherWorker {
	logger := logging.NewStdLogger("[WEBHOOK-DISPATCHER]")

	logger.Info("Initializing WebhookDispatcherWorker", map[string]interface{}{
		"brokers":         cfg.Kafka.Brokers,
		"topic":           cfg.Kafka.Topics.WebhookDispatch,
		"group_id":        cfg.Kafka.GroupID,
		"min_bytes":       "10KB",
		"max_bytes":       "10MB",
		"request_timeout": cfg.Webhook.RequestTimeout.String(),
		"max_retries":     cfg.Webhook.MaxRetries,
		"retry_intervals": cfg.Webhook.RetryIntervals,
		"has_db":          db != nil,
	})

	worker := &WebhookDispatcherWorker{
		cfg: cfg,
		db:  db,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  cfg.Kafka.Brokers,
			Topic:    cfg.Kafka.Topics.WebhookDispatch,
			GroupID:  cfg.Kafka.GroupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		client:             &http.Client{Timeout: cfg.Webhook.RequestTimeout},
		logger:             logger,
		usecase:            usecase,
		transactionUsecase: transactionUsecase,
		hostedPaymentRepo:  hostedPaymentRepo,
	}

	logger.Info("WebhookDispatcherWorker initialized successfully", map[string]interface{}{
		"kafka_topic": cfg.Kafka.Topics.WebhookDispatch,
		"group_id":    cfg.Kafka.GroupID,
	})

	return worker
}

func (w *WebhookDispatcherWorker) Start(ctx context.Context) {
	w.logger.Info("Starting WebhookDispatcherWorker", map[string]interface{}{
		"kafka_topic": w.cfg.Kafka.Topics.WebhookDispatch,
		"group_id":    w.cfg.Kafka.GroupID,
	})

	defer func() {
		w.logger.Info("Closing Kafka reader", map[string]interface{}{
			"topic": w.cfg.Kafka.Topics.WebhookDispatch,
		})
		w.reader.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Context cancelled, stopping worker", map[string]interface{}{
				"reason": ctx.Err().Error(),
			})
			return
		default:
			msg, err := w.reader.ReadMessage(ctx)
			if err != nil {
				w.logger.Error("Failed to read message from Kafka", map[string]interface{}{
					"error":    err.Error(),
					"topic":    w.cfg.Kafka.Topics.WebhookDispatch,
					"brokers":  w.cfg.Kafka.Brokers,
					"group_id": w.cfg.Kafka.GroupID,
				})
				continue
			}

			w.logger.Debug("Received webhook dispatch message", map[string]interface{}{
				"topic":      msg.Topic,
				"partition":  msg.Partition,
				"offset":     msg.Offset,
				"key":        string(msg.Key),
				"value_size": len(msg.Value),
				"value":      string(msg.Value),
				"timestamp":  msg.Time,
			})

			var webhookMsg webhookDto.WebhookMessage
			if err := json.Unmarshal(msg.Value, &webhookMsg); err != nil {
				w.logger.Error("Failed to unmarshal webhook message", map[string]interface{}{
					"error":     err.Error(),
					"raw_value": string(msg.Value),
					"key":       string(msg.Key),
				})
				continue
			}

			w.logger.Info("Processing webhook message", map[string]interface{}{
				"type":           webhookMsg.Type,
				"transaction_id": webhookMsg.TransactionID,
				"merchant_id":    webhookMsg.MerchantID,
				"status":         webhookMsg.Status,
				"message":        webhookMsg.Message,
				"total_amount":   webhookMsg.TotalAmount,
				"provider_txid":  webhookMsg.ProviderTxID,
				"message_time":   webhookMsg.Timestamp,
				"user_id":        webhookMsg.UserID,
			})

			if err := w.processMessage(ctx, webhookMsg); err != nil {
				w.logger.Error("Failed to process webhook message", map[string]interface{}{
					"error":          err.Error(),
					"transaction_id": webhookMsg.TransactionID,
					"status":         webhookMsg.Status,
					"user_id":        webhookMsg.UserID,
				})
			} else {
				w.logger.Info("Successfully processed webhook message", map[string]interface{}{
					"transaction_id": webhookMsg.TransactionID,
					"status":         webhookMsg.Status,
					"user_id":        webhookMsg.UserID,
				})
			}
		}
	}
}

func (w *WebhookDispatcherWorker) processMessage(ctx context.Context, msg webhookDto.WebhookMessage) error {
	w.logger.Debug("Starting webhook message processing", map[string]interface{}{
		"type":            msg.Type,
		"transaction_id":  msg.TransactionID,
		"max_retries":     w.cfg.Webhook.MaxRetries,
		"retry_intervals": w.cfg.Webhook.RetryIntervals,
	})

	var retries int
	for retries < w.cfg.Webhook.MaxRetries {
		w.logger.Info("Attempting to update payment status", map[string]interface{}{
			"type":           msg.Type,
			"transaction_id": msg.TransactionID,
			"attempt":        retries + 1,
			"max_attempts":   w.cfg.Webhook.MaxRetries,
		})

		err := w.usecase.HandlePaymentStatusUpdate(ctx, msg)
		if err != nil {
			retries++
			retryInterval := w.cfg.Webhook.RetryIntervals[retries-1]

			w.logger.Warn("Payment status update failed, will retry", map[string]interface{}{
				"error":          err.Error(),
				"transaction_id": msg.TransactionID,
				"retry_attempt":  retries,
				"retry_interval": retryInterval.String(),
				"next_retry_at":  time.Now().Add(retryInterval),
			})

			time.Sleep(retryInterval)
			continue
		}

		// if msg.IsHostedCheckout {
		// 	// update the hosted checkout
		// 	err := w.hostedPaymentRepo.UpdateStatus(ctx, uuid.MustParse(msg.TransactionID),
		// 		entity.HostedPaymentCompleted)

		// 	if err != nil {
		// 		retries++
		// 		retryInterval := w.cfg.Webhook.RetryIntervals[retries-1]

		// 		w.logger.Warn("Failed to update hosted checkout status ,will retry", map[string]interface{}{
		// 			"error":          err.Error(),
		// 			"operation":      "getHostedCheckout",
		// 			"transaction_id": msg.TransactionID,
		// 		})

		// 		time.Sleep(retryInterval)
		// 		continue
		// 	}

		// }

		w.logger.Info("Payment status updated successfully", map[string]interface{}{
			"transaction_id":  msg.TransactionID,
			"attempts_needed": retries + 1,
			"status":          msg.Status,
		})
		break
	}

	if retries >= w.cfg.Webhook.MaxRetries {
		err := fmt.Errorf("max retries exceeded for payment status update")
		w.logger.Error("Failed to update payment status after all retries", map[string]interface{}{
			"error":          err.Error(),
			"transaction_id": msg.TransactionID,
			"attempts_made":  retries,
			"max_retries":    w.cfg.Webhook.MaxRetries,
		})
		return err
	}

	return nil
}
