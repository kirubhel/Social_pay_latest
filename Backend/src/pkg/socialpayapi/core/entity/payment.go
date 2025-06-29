package entity

import (
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	merchantEntity "github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
)

// DirectPaymentRequest represents the request for direct payment
// @Description Payment request details for processing a direct payment
type DirectPaymentRequest struct {
	// Payment medium/provider to use
	// @Example MPESA
	Medium entity.TransactionMedium `json:"medium" example:"MPESA"`

	// Description of the payment
	// @Example Payment for order #123
	Description string `json:"description" example:"Payment for order #123"`

	// Phone number to be paid
	// @Example 251911111111
	PhoneNumber string `json:"phone_number" example:"251911111111"`

	// Merchant ID
	// @Example 1234567890
	Reference string `json:"reference" example:"1234567890"`

	// Amount to be paid
	// @Example 1000.50
	Amount float64 `json:"amount" example:"1000.50"`

	// Three-letter currency code
	// @Example ETB
	Currency string `json:"currency" example:"ETB"`

	// Additional payment details
	Details entity.TransactionDetails `json:"details"`

	// URLs for payment redirects
	Redirects entity.TransactionRedirects `json:"redirects"`

	// URL to receive payment status updates
	// @Example https://example.com/callback
	CallbackURL string `json:"callback_url" example:"https://example.com/callback"`
}

func (r DirectPaymentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Medium, validation.Required, validation.In(
			entity.CYBERSOURCE,
			entity.MPESA,
			entity.TELEBIRR,
			entity.CBE,
			entity.AWASH, // Awash
		)),
		validation.Field(&r.Amount, validation.Required, validation.Min(0.01)),
		validation.Field(&r.Currency, validation.Required, validation.Length(3, 3)),
		validation.Field(&r.Details, validation.Required),
		validation.Field(&r.PhoneNumber, validation.Required, validation.Length(12, 12), validation.Match(regexp.MustCompile(`^251\d{9}$`))),
		validation.Field(&r.Redirects, validation.Required),
		validation.Field(&r.CallbackURL, validation.Required, is.URL),
	)
}

// HostedCheckoutRequest represents the request for hosted checkout
// @Description Hosted checkout request details for creating a payment link
type HostedCheckoutRequest struct {
	// Amount to be paid
	// @Example 1000.50
	Amount float64 `json:"amount" example:"1000.50"`

	// Three-letter currency code
	// @Example ETB
	Currency string `json:"currency" example:"ETB"`

	// Description of the payment
	// @Example Payment for order #123
	Description string `json:"description" example:"Payment for order #123"`

	// Merchant reference
	// @Example ORDER123456
	Reference string `json:"reference" example:"ORDER123456"`

	// Supported payment mediums
	// @Example ["MPESA", "TELEBIRR", "CBE"]
	SupportedMediums []entity.TransactionMedium `json:"supported_mediums" example:"[\"MPESA\", \"TELEBIRR\", \"CBE\",\"ETHSWITCH\"]"`

	// Optional phone number (can be pre-filled)
	// @Example 251911111111
	PhoneNumber string `json:"phone_number,omitempty" example:"251911111111"`

	// URLs for payment redirects
	Redirects entity.TransactionRedirects `json:"redirects"`

	// URL to receive payment status updates
	// @Example https://example.com/callback
	CallbackURL string `json:"callback_url,omitempty" example:"https://example.com/callback"`

	// Optional expiry date time in UTC (ISO 8601 format)
	// @Example 2024-12-31T23:59:59Z
	ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

func (r HostedCheckoutRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Amount, validation.Required, validation.Min(0.01)),
		validation.Field(&r.Currency, validation.Required, validation.Length(3, 3)),
		validation.Field(&r.Reference, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.SupportedMediums, validation.Required, validation.Length(1, 10)),
		validation.Field(&r.Redirects, validation.Required),
		validation.Field(&r.CallbackURL, validation.Required, is.URL),
		validation.Field(&r.PhoneNumber, validation.When(
			r.PhoneNumber != "",
			validation.Length(12, 12),
			validation.Match(regexp.MustCompile(`^251\d{9}$`)),
		)),
		validation.Field(&r.ExpiresAt, validation.When(r.ExpiresAt != nil, validation.By(func(value interface{}) error {
			if expiresAt, ok := value.(*time.Time); ok && expiresAt != nil {
				if expiresAt.Before(time.Now().UTC()) {
					return fmt.Errorf("expires_at must be in the future")
				}
			}
			return nil
		}))),
	)
}

// CheckoutPaymentRequest represents the request when user makes payment from hosted checkout
// @Description Payment request from hosted checkout page
type CheckoutPaymentRequest struct {
	// Hosted checkout ID
	// @Example 123e4567-e89b-12d3-a456-426614174000
	HostedCheckoutID uuid.UUID `json:"hosted_checkout_id" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Selected payment medium
	// @Example MPESA
	Medium entity.TransactionMedium `json:"medium" example:"MPESA,ETHSWITCH"`

	// Phone number for payment
	// @Example 251911111111
	PhoneNumber string `json:"phone_number" example:"251911111111"`
}

func (r CheckoutPaymentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.HostedCheckoutID, validation.Required),
		validation.Field(&r.Medium, validation.Required, validation.In(
			entity.CYBERSOURCE,
			entity.MPESA,
			entity.TELEBIRR,
			entity.CBE,
			entity.ETHSWITCH,
		)),
		validation.Field(&r.PhoneNumber, validation.Required, validation.Length(12, 12), validation.Match(regexp.MustCompile(`^251\d{9}$`))),
	)
}

// HostedCheckoutResponseDTO represents the response for hosted checkout details
// @Description Response containing hosted checkout information for the checkout page
type HostedCheckoutResponseDTO struct {
	// Unique identifier for the hosted checkout
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Payment details
	Amount      float64 `json:"amount" example:"1000.50"`
	Currency    string  `json:"currency" example:"ETB"`
	Description string  `json:"description" example:"Payment for order #123"`
	Reference   string  `json:"reference" example:"ORDER123456"`

	// Supported payment mediums
	SupportedMediums []entity.TransactionMedium `json:"supported_mediums" example:"[\"MPESA\", \"TELEBIRR\", \"CBE\"]"`

	// Optional pre-filled phone number
	PhoneNumber string `json:"phone_number,omitempty" example:"251911111111"`

	// URLs for redirects
	SuccessURL string `json:"success_url" example:"https://example.com/success"`
	FailedURL  string `json:"failed_url" example:"https://example.com/failed"`

	// Status and timestamps
	Status    string    `json:"status" example:"PENDING"`
	CreatedAt time.Time `json:"created_at" example:"2024-05-03T10:00:00Z"`
	ExpiresAt time.Time `json:"expires_at" example:"2024-05-04T10:00:00Z"`

	// Merchant information (optional)
	MerchantName string `json:"merchant_name,omitempty" example:"Example Store"`
}

// HostedCheckoutWithMerchantResponseDTO represents the response for hosted checkout details with merchant information
// @Description Response containing hosted checkout information and merchant details for the checkout page
type HostedCheckoutWithMerchantResponseDTO struct {
	// Unique identifier for the hosted checkout
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Payment details
	Amount      float64 `json:"amount" example:"1000.50"`
	Currency    string  `json:"currency" example:"ETB"`
	Description string  `json:"description" example:"Payment for order #123"`
	Reference   string  `json:"reference" example:"ORDER123456"`

	// Supported payment mediums
	SupportedMediums []entity.TransactionMedium `json:"supported_mediums" example:"[\"MPESA\", \"TELEBIRR\", \"CBE\"]"`

	// Optional pre-filled phone number
	PhoneNumber string `json:"phone_number,omitempty" example:"251911111111"`

	// URLs for redirects
	SuccessURL string `json:"success_url" example:"https://example.com/success"`
	FailedURL  string `json:"failed_url" example:"https://example.com/failed"`

	// Status and timestamps
	Status    string    `json:"status" example:"PENDING"`
	CreatedAt time.Time `json:"created_at" example:"2024-05-03T10:00:00Z"`
	ExpiresAt time.Time `json:"expires_at" example:"2024-05-04T10:00:00Z"`

	// Merchant information
	Merchant *merchantEntity.Merchant `json:"merchant,omitempty"`
}

// PaymentResponse represents the response for payment operations
// @Description Response containing payment processing results
type PaymentResponse struct {
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

	// Socialpay transaction ID
	SocialPayTransactionID string `json:"Socialpay_transaction_id,omitempty" example:"1234567890"`
}

// WithdrawalRequest represents the request for withdrawal
// @Description Withdrawal request details
type WithdrawalRequest struct {
	// Amount to withdraw
	Amount float64 `json:"amount" example:"500.00"`

	// Three-letter currency code
	Currency string `json:"currency" example:"ETB"`

	// Payment medium/provider to use
	// @Example MPESA
	Medium entity.TransactionMedium `json:"medium" example:"MPESA"`

	// URL to receive payment status updates
	// @Example https://example.com/callback
	CallbackURL string `json:"callback_url" example:"https://example.com/callback"`

	// Bank account number
	PhoneNumber string `json:"phone_number" example:"1234567890"`

	// Client-provided reference
	Reference string `json:"reference" example:"WD123456789"`
}

func (r WithdrawalRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Amount, validation.Required, validation.Min(0.01)),
		validation.Field(&r.Currency, validation.Required, validation.Length(3, 3)),
		validation.Field(&r.PhoneNumber, validation.Required, validation.Length(5, 20)),
		validation.Field(&r.CallbackURL, validation.Required, is.URL),
		validation.Field(&r.Reference, validation.Required, validation.Length(3, 50)),
	)
}

// TransactionQuery represents the query parameters for transaction lookup
type TransactionQuery struct {
	// Transaction UUID
	ID uuid.UUID `uri:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

func (q TransactionQuery) Validate() error {
	return validation.ValidateStruct(&q,
		validation.Field(&q.ID, validation.Required),
	)
}

// TransactionResponseDTO represents the API response for transaction details
// @Description API response containing transaction information
type TransactionResponseDTO struct {
	// Unique identifier for the transaction
	Id uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Merchant ID
	MerchantId uuid.UUID `json:"merchant_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Phone number of the payer
	PhoneNumber string `json:"phone_number" example:"+251911234567"`
	// User ID associated with the transaction
	UserId uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Type of transaction
	Type entity.TransactionType `json:"type" example:"SALE"`
	// Payment medium used
	Medium entity.TransactionMedium `json:"medium" example:"MPESA"`
	// Reference code
	Reference string `json:"reference" example:"REF123456"`
	// Additional comment
	Comment string `json:"comment" example:"Payment for order #123"`
	// Whether the transaction is verified
	Verified bool `json:"verified" example:"true"`
	// Additional transaction details
	Details interface{} `json:"details"`
	// Creation timestamp
	CreatedAt time.Time `json:"created_at" example:"2024-05-03T10:00:00Z"`
	// Last update timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2024-05-03T10:01:00Z"`
	// Reference number
	ReferenceNumber string `json:"reference_number" example:"TXN123456"`
	// Whether this is a test transaction
	Test bool `json:"test" example:"false"`
	// Current status
	Status entity.TransactionStatus `json:"status" example:"SUCCESS"`
	// Transaction description
	Description string `json:"description" example:"Payment for service"`
	// Security token
	Token string `json:"token" example:"tok_123456"`
	// Transaction amount
	Amount float64 `json:"amount" example:"1000.50"`

	//webhook_received
	WebhookReceived bool `json:"webhook_received" example:"false"`

	// Fee amount
	FeeAmount float64 `json:"fee_amount" example:"10.00"`
	// Admin net amount
	AdminNet float64 `json:"admin_net" example:"5.00"`
	// VAT amount
	VatAmount float64 `json:"vat_amount" example:"15.00"`
	// Merchant net amount
	MerchantNet float64 `json:"merchant_net" example:"970.50"`
	// Total amount
	TotalAmount float64 `json:"total_amount" example:"1000.50"`
	// Currency code
	Currency string `json:"currency" example:"ETB"`
	// Callback URL
	CallbackURL string `json:"callback_url" example:"https://example.com/callback"`
	// Success redirect URL
	SuccessURL string `json:"success_url" example:"https://example.com/success"`
	// Failed redirect URL
	FailedURL string `json:"failed_url" example:"https://example.com/failed"`
	// Merchant information
	Merchant *merchantEntity.Merchant `json:"merchant,omitempty"`
}

type DirectPaymentResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"`
	RedirectURL   string    `json:"redirect_url,omitempty"`
	Message       string    `json:"message,omitempty"`
}

// UpdateHostedCheckoutRequest represents the request for updating hosted checkout
// @Description Request for updating hosted checkout details (only allowed when status is PENDING)
type UpdateHostedCheckoutRequest struct {
	// Amount to be paid
	// @Example 1000.50
	Amount *float64 `json:"amount,omitempty" example:"1000.50"`

	// Three-letter currency code
	// @Example ETB
	Currency *string `json:"currency,omitempty" example:"ETB"`

	// Description of the payment
	// @Example Payment for order #123
	Description *string `json:"description,omitempty" example:"Payment for order #123"`

	// Supported payment mediums
	// @Example ["MPESA", "TELEBIRR", "CBE"]
	SupportedMediums []entity.TransactionMedium `json:"supported_mediums,omitempty" example:"[\"MPESA\", \"TELEBIRR\", \"CBE\"]"`

	// Optional phone number (can be pre-filled)
	// @Example 251911111111
	PhoneNumber *string `json:"phone_number,omitempty" example:"251911111111"`

	// URLs for payment redirects
	Redirects *entity.TransactionRedirects `json:"redirects,omitempty"`

	// URL to receive payment status updates
	// @Example https://example.com/callback
	CallbackURL *string `json:"callback_url,omitempty" example:"https://example.com/callback"`

	// Optional expiry date time in UTC (ISO 8601 format)
	// @Example 2024-12-31T23:59:59Z
	ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

func (r UpdateHostedCheckoutRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Amount, validation.When(r.Amount != nil, validation.Min(0.01))),
		validation.Field(&r.Currency, validation.When(r.Currency != nil, validation.Length(3, 3))),
		validation.Field(&r.SupportedMediums, validation.When(len(r.SupportedMediums) > 0, validation.Length(1, 10))),
		validation.Field(&r.CallbackURL, validation.When(r.CallbackURL != nil && *r.CallbackURL != "", is.URL)),
		validation.Field(&r.PhoneNumber, validation.When(
			r.PhoneNumber != nil && *r.PhoneNumber != "",
			validation.Length(12, 12),
			validation.Match(regexp.MustCompile(`^251\d{9}$`)),
		)),
		validation.Field(&r.ExpiresAt, validation.When(r.ExpiresAt != nil, validation.By(func(value interface{}) error {
			if expiresAt, ok := value.(*time.Time); ok && expiresAt != nil {
				if expiresAt.Before(time.Now().UTC()) {
					return fmt.Errorf("expires_at must be in the future")
				}
			}
			return nil
		}))),
	)
}
