package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CallbackLog struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	TxnID        uuid.UUID `json:"txn_id"`
	MerchantID   uuid.UUID `json:"merchant_id"`
	Status       int       `json:"status"`
	RequestBody  string    `json:"request_body"`
	ResponseBody string    `json:"response_body"`
	RetryCount   int       `json:"retry_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Message      string    `json:"message"`
}

// Validate checks if the callback log is valid
func (c *CallbackLog) Validate() error {
	if !IsValidStatus(c.Status) {
		return fmt.Errorf("invalid status code: %d", c.Status)
	}
	return nil
}
