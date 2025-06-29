package entity

import (
	"time"

	"github.com/google/uuid"
)

// HostedPaymentStatus represents the status of a hosted payment
type HostedPaymentStatus string

const (
	HostedPaymentPending   HostedPaymentStatus = "PENDING"
	HostedPaymentCompleted HostedPaymentStatus = "COMPLETED"
	HostedPaymentExpired   HostedPaymentStatus = "EXPIRED"
	HostedPaymentCanceled  HostedPaymentStatus = "CANCELED"
)

// HostedPayment represents a hosted checkout payment
type HostedPayment struct {
	// Unique identifier for the hosted payment
	ID uuid.UUID `json:"id"`

	// User and merchant information
	UserID     uuid.UUID `json:"user_id"`
	MerchantID uuid.UUID `json:"merchant_id"`

	// Payment details
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	Reference   string  `json:"reference"`

	// Supported payment mediums
	SupportedMediums []TransactionMedium `json:"supported_mediums"`

	// Optional phone number from merchant
	PhoneNumber string `json:"phone_number,omitempty"`

	// URLs for redirects and callbacks
	SuccessURL  string `json:"success_url"`
	FailedURL   string `json:"failed_url"`
	CallbackURL string `json:"callback_url,omitempty"`

	// Status and transaction linking
	Status        HostedPaymentStatus `json:"status"`
	TransactionID *uuid.UUID          `json:"transaction_id,omitempty"`

	// Selected payment details (filled when user makes payment)
	SelectedMedium      string `json:"selected_medium,omitempty"`
	SelectedPhoneNumber string `json:"selected_phone_number,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
