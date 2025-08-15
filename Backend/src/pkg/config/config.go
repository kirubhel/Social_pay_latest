package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Kafka struct {
		Brokers []string
		Topics  struct {
			WebhookDispatch string
			PaymentStatus   string
			WebhookSend     string
		}
		GroupID string
	}
	Webhook struct {
		RequestTimeout time.Duration
		MaxRetries     int
		RetryIntervals []time.Duration
	}
}	

func Load() (*Config, error) {
	cfg := &Config{}

	// Kafka configuration
	cfg.Kafka.Brokers = []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
	cfg.Kafka.Topics.WebhookDispatch = getEnv("KAFKA_TOPIC_WEBHOOK_DISPATCH", "webhook_dispatch")
	cfg.Kafka.Topics.PaymentStatus = getEnv("KAFKA_TOPIC_PAYMENT_STATUS", "payment_status")
	cfg.Kafka.Topics.WebhookSend = getEnv("KAFKA_TOPIC_WEBHOOK_SEND", "webhook_send")
	cfg.Kafka.GroupID = getEnv("KAFKA_GROUP_ID", "webhook-service")

	// Webhook configuration
	cfg.Webhook.MaxRetries, _ = strconv.Atoi(getEnv("WEBHOOK_MAX_RETRIES", "3"))
	cfg.Webhook.RequestTimeout, _ = time.ParseDuration(getEnv("WEBHOOK_REQUEST_TIMEOUT", "30s"))

	// Set retry intervals
	cfg.Webhook.RetryIntervals = []time.Duration{
		time.Second,
		5 * time.Second,
		8 * time.Second,
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
