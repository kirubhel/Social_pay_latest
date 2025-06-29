package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	Kafka struct {
		Brokers []string
		Topics  struct {
			PaymentStatus  string
			WebhookDispatch string
		}
		GroupID string
	}
	Webhook struct {
		MaxRetries      int
		RetryIntervals  []time.Duration
		RequestTimeout  time.Duration
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Database configuration
	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port, _ = strconv.Atoi(getEnv("DB_PORT", "5432"))
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.DB.DBName = getEnv("DB_NAME", "webhook")
	cfg.DB.SSLMode = getEnv("DB_SSL_MODE", "disable")

	// Kafka configuration
	cfg.Kafka.Brokers = []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
	cfg.Kafka.Topics.PaymentStatus = getEnv("KAFKA_TOPIC_PAYMENT_STATUS", "payment_status")
	cfg.Kafka.Topics.WebhookDispatch = getEnv("KAFKA_TOPIC_WEBHOOK_DISPATCH", "webhook_dispatch")
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