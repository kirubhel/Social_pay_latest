package entity

import (
	"time"

	"github.com/google/uuid"
)

// an ENUM type for transaction status, it will be stored in the database
type TransactionStatus string

const (
	Pending   TransactionStatus = "Pending"
	Completed TransactionStatus = "Completed"
	Failed    TransactionStatus = "Failed"
	Reversed  TransactionStatus = "Reversed"
)

type AirtimeTransaction struct {
	ID                   uuid.UUID         `json:"id" db:"id"`
	MerchantID           uuid.UUID         `json:"merchant_id" db:"merchant_id"`
	Amount               float64           `json:"amount" db:"amount"`
	MSISDN               string            `json:"msisdn" db:"phone_number"`
	TransactionRef       string            `json:"transaction_ref" db:"transaction_ref"`
	PaymentMethod        string            `json:"payment_method" db:"payment_method"`
	Currency             string            `json:"currency" db:"currency"`
	Status               TransactionStatus `json:"status" db:"status"`
	CreatedAt            time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at" db:"updated_at"`
	YimuluTransactionRef string            `json:"yimulu_transaction_ref" db:"yimulu_transaction_ref"`
	TransactionDetails   []byte            `json:"transaction_details" db:"transaction_details"`
	WebhookNotified      bool              `json:"webhook_notified" db:"webhook_notified"` // Flag to track if the webhook has been notified
}
