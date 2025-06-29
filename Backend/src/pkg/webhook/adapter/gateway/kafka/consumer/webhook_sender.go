package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/socialpay/socialpay/src/pkg/config"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
	webhookUsecase "github.com/socialpay/socialpay/src/pkg/webhook/usecase"
)

type WebhookSenderWorker struct {
	cfg     *config.Config
	reader  *kafka.Reader
	client  *http.Client
	logger  logging.Logger
	usecase webhookUsecase.WebhookUseCase
}

func NewWebhookSenderWorker(cfg *config.Config, usecase webhookUsecase.WebhookUseCase) *WebhookSenderWorker {
	logger := logging.NewStdLogger("[WEBHOOK-SENDER]")

	logger.Info("Initializing WebhookSenderWorker", map[string]interface{}{
		"brokers":         cfg.Kafka.Brokers,
		"topic":           cfg.Kafka.Topics.WebhookSend,
		"group_id":        cfg.Kafka.GroupID,
		"min_bytes":       "10KB",
		"max_bytes":       "10MB",
		"request_timeout": cfg.Webhook.RequestTimeout.String(),
		"max_retries":     cfg.Webhook.MaxRetries,
		"retry_intervals": cfg.Webhook.RetryIntervals,
	})

	worker := &WebhookSenderWorker{
		cfg: cfg,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  cfg.Kafka.Brokers,
			Topic:    cfg.Kafka.Topics.WebhookSend,
			GroupID:  cfg.Kafka.GroupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		client:  &http.Client{Timeout: cfg.Webhook.RequestTimeout},
		logger:  logger,
		usecase: usecase,
	}

	logger.Info("WebhookSenderWorker initialized successfully", map[string]interface{}{
		"kafka_topic": cfg.Kafka.Topics.WebhookSend,
		"group_id":    cfg.Kafka.GroupID,
	})

	return worker
}

func (w *WebhookSenderWorker) Start(ctx context.Context) {
	w.logger.Info("Starting WebhookSenderWorker", map[string]interface{}{
		"kafka_topic": w.cfg.Kafka.Topics.WebhookSend,
		"group_id":    w.cfg.Kafka.GroupID,
	})

	defer func() {
		w.logger.Info("Closing Kafka reader", map[string]interface{}{
			"topic": w.cfg.Kafka.Topics.WebhookSend,
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
					"topic":    w.cfg.Kafka.Topics.WebhookSend,
					"brokers":  w.cfg.Kafka.Brokers,
					"group_id": w.cfg.Kafka.GroupID,
				})
				continue
			}

			w.logger.Debug("Received webhook send message", map[string]interface{}{
				"topic":      msg.Topic,
				"partition":  msg.Partition,
				"offset":     msg.Offset,
				"key":        string(msg.Key),
				"value_size": len(msg.Value),
				"timestamp":  msg.Time,
			})

			var webhookMsg dto.WebhookEventMerchant
			if err := json.Unmarshal(msg.Value, &webhookMsg); err != nil {
				w.logger.Error("Failed to unmarshal webhook message", map[string]interface{}{
					"error":     err.Error(),
					"raw_value": string(msg.Value),
					"key":       string(msg.Key),
				})
				continue
			}

			w.logger.Info("Processing webhook send message", map[string]interface{}{
				"event":          webhookMsg.Event,
				"transaction_id": webhookMsg.SocialPayTxnID,
				"status":         webhookMsg.Status,
				"message":        webhookMsg.Message,
				"provider_txid":  webhookMsg.ProviderTxID,
				"timestamp":      webhookMsg.Timestamp,
				"callback_url":   webhookMsg.CallbackURL,
			})

			if err := w.processMessage(ctx, webhookMsg); err != nil {
				w.logger.Error("Failed to process webhook message", map[string]interface{}{
					"error":          err.Error(),
					"transaction_id": webhookMsg.SocialPayTxnID,
				})
			}
		}
	}
}

func (w *WebhookSenderWorker) processMessage(ctx context.Context, msg dto.WebhookEventMerchant) error {
	var retries int
	var lastErr error

	for retries < w.cfg.Webhook.MaxRetries {
		w.logger.Info("Attempting to send webhook", map[string]interface{}{
			"transaction_id": msg.SocialPayTxnID,
			"attempt":        retries + 1,
			"max_attempts":   w.cfg.Webhook.MaxRetries,
			"callback_url":   msg.CallbackURL,
		})

		responseStatus, responseBody, err := w.sendWebhook(msg.CallbackURL, msg)
		if err != nil {
			lastErr = err
			retries++
			retryInterval := w.cfg.Webhook.RetryIntervals[retries-1]

			w.logger.Warn("Webhook send failed, will retry", map[string]interface{}{
				"error":          err.Error(),
				"transaction_id": msg.SocialPayTxnID,
				"retry_attempt":  retries,
				"retry_interval": retryInterval.String(),
				"next_retry_at":  time.Now().Add(retryInterval),
			})

			time.Sleep(retryInterval)
			continue
		}

		// Log the webhook response
		requestBody, _ := json.Marshal(msg)
		txnID, err := uuid.Parse(msg.SocialPayTxnID)
		if err != nil {
			w.logger.Error("Failed to parse transaction ID", map[string]interface{}{
				"error":          err.Error(),
				"transaction_id": msg.SocialPayTxnID,
			})
			return fmt.Errorf("failed to parse transaction ID: %w", err)
		}

		if err := w.usecase.CreateCallbackLog(ctx, txnID, responseStatus, string(requestBody), responseBody, msg.MerchantID, msg.UserID); err != nil {
			w.logger.Error("Failed to create callback log", map[string]interface{}{
				"error":          err.Error(),
				"transaction_id": msg.SocialPayTxnID,
			})
			return fmt.Errorf("failed to create callback log: %w", err)
		}

		w.logger.Info("Webhook sent successfully", map[string]interface{}{
			"transaction_id":  msg.SocialPayTxnID,
			"attempts_needed": retries + 1,
			"status":          msg.Status,
			"message":         msg.Message,
			"response_status": responseStatus,
		})
		return nil
	}

	if lastErr != nil {
		w.logger.Error("Failed to send webhook after all retries", map[string]interface{}{
			"error":          lastErr.Error(),
			"transaction_id": msg.SocialPayTxnID,
			"attempts_made":  retries,
			"max_retries":    w.cfg.Webhook.MaxRetries,
		})
		return fmt.Errorf("failed to send webhook after %d retries: %w", retries, lastErr)
	}

	return nil
}

func (w *WebhookSenderWorker) sendWebhook(url string, payload dto.WebhookEventMerchant) (int, string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Social Pay")

	resp, err := w.client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", fmt.Errorf("failed to read response body: %w", err)
	}

	return resp.StatusCode, string(body), nil
}
