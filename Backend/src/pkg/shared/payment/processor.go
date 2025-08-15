package payment

import (
	"context"

	"github.com/google/uuid"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

// PaymentStatus represents the status of a payment

// PaymentRequest represents a unified payment request structure
type PaymentRequest struct {
	TransactionID uuid.UUID                  `json:"transaction_id"`
	Amount        float64                    `json:"amount"`
	Medium        txEntity.TransactionMedium `json:"medium"`
	Currency      string                     `json:"currency"`
	PhoneNumber   string                     `json:"phone_number,omitempty"`
	Reference     string                     `json:"reference"`
	Description   string                     `json:"description"`
	CallbackURL   string                     `json:"callback_url"`
	SuccessURL    string                     `json:"success_url"`
	FailedURL     string                     `json:"failed_url"`
	Metadata      map[string]interface{}     `json:"metadata,omitempty"`
}

// PaymentResponse represents a unified payment response structure
type PaymentResponse struct {
	TransactionID uuid.UUID                  `json:"transaction_id"`
	Success       bool                       `json:"success"`
	Status        txEntity.TransactionStatus `json:"status"`
	PaymentURL    string                     `json:"payment_url,omitempty"`
	Message       string                     `json:"message,omitempty"`
	ProcessorRef  string                     `json:"processor_ref,omitempty"`
	Metadata      map[string]interface{}     `json:"metadata,omitempty"`
}

type TransactionStatusQueryResponse struct {
	Status       txEntity.TransactionStatus `json:"status"`
	ProviderTxId string                     `json:"provider_tx_id"`
	ProviderData map[string]interface{}     `json:"provider_data"`
}

// CallbackRequest represents the payment callback request
type CallbackRequest struct {
	TransactionID uuid.UUID                  `json:"transaction_id"`
	ProcessorRef  string                     `json:"processor_ref"`
	Status        txEntity.TransactionStatus `json:"status"`
	Metadata      map[string]interface{}     `json:"metadata,omitempty"`
}

// Processor defines the interface that all payment processors must implement
type Processor interface {
	// InitiatePayment starts a payment transaction
	InitiatePayment(ctx context.Context, apikey string, req *PaymentRequest) (*PaymentResponse, error)

	// SettlePayment handles the payment callback/settlement
	SettlePayment(ctx context.Context, req *CallbackRequest) error

	InitiateWithdrawal(ctx context.Context, apikey string, req *PaymentRequest) (*PaymentResponse, error)

	// GetType returns the processor type
	GetType() txEntity.TransactionMedium

	QueryTransactionStatus(ctx context.Context, transactionID string) (*TransactionStatusQueryResponse, error)
}
