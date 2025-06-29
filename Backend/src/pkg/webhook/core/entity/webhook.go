package entity

import (
	"time"

	"github.com/google/uuid"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

// Status codes for webhook callbacks
const (
	StatusPending   = 0
	StatusInitiated = 1
	StatusSuccess   = 2
	StatusFailed    = 3
	StatusExpired   = 4
)

// WebhookMessage represents a message to be sent to a webhook URL
type WebhookMessage struct {
	TransactionID uuid.UUID
	ReferenceID   string
	Status        txEntity.TransactionStatus
	Amount        float64
	Timestamp     time.Time
}

// IsValidStatus checks if a status code is valid
func IsValidStatus(status int) bool {
	return status >= StatusPending && status <= StatusExpired
}

