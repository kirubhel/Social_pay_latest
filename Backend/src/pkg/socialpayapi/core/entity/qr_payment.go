package entity

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

// QRPaymentRequest represents the request for QR payment
// @Description QR Payment request details for processing a payment from checkout
type QRMerchantPaymentRequest struct {
	// Payment medium/provider to use
	// @Example MPESA
	Medium entity.TransactionMedium `json:"medium" example:"MPESA"`

	// Amount to be paid
	// @Example 1000.50
	Amount float64 `json:"amount" example:"1000.50"`

	// Merchant ID
	// @Example 123e4567-e89b-12d3-a456-426614174000
	MerchantID uuid.UUID `json:"merchant_id" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Optional: Three-letter currency code (defaults to ETB)
	// @Example ETB
	Currency string `json:"currency,omitempty" example:"ETB"`

	// Optional: Description of the payment
	// @Example Payment for order #123
	Description string `json:"description,omitempty" example:"Payment for order #123"`

	// Optional: Phone number to be paid
	// @Example 251911111111
	PhoneNumber string `json:"phone_number,omitempty" example:"251911111111"`

	// Optional: Client-provided reference
	// @Example ORDER123456
	Reference string `json:"reference,omitempty" example:"ORDER123456"`
}

func (r QRMerchantPaymentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Medium, validation.Required, validation.In(
			entity.CYBERSOURCE,
			entity.MPESA,
			entity.TELEBIRR,
			entity.CBE,
			entity.KACHA,
			entity.AWASH,
		)),
		validation.Field(&r.Amount, validation.Required, validation.Min(0.01)),
		validation.Field(&r.MerchantID, validation.Required),
	)
}

// QRPaymentResponse represents the response for QR payment operations
// @Description Response containing QR payment processing results
type QRPaymentResponse struct {
	// Whether the operation was successful
	Success bool `json:"success" example:"true"`

	// Status of the payment
	Status string `json:"status" example:"PENDING"`

	// Human-readable message about the operation
	Message string `json:"message" example:"Payment initiated successfully"`

	// URL to redirect the user for payment (if applicable)
	PaymentURL string `json:"payment_url,omitempty" example:"https://pay.example.com/checkout"`

	// Unique payment reference
	Reference string `json:"reference_id" example:"PAY123456789"`

	// SocialPay transaction ID
	SocialPayTransactionID string `json:"socialpay_transaction_id,omitempty" example:"1234567890"`
}
