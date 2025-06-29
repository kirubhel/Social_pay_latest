package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
)

// TransactionType represents the type of transaction
// @Description Type of transaction (e.g., SALE, WITHDRAWAL)
type TransactionType string

const (
	REPLENISHMENT TransactionType = "REPLENISHMENT"
	P2P           TransactionType = "P2P"
	SALE          TransactionType = "SALE"
	BILL_PAYMENT  TransactionType = "BILL_PAYMENT"
	SETTLEMENT    TransactionType = "SETTLEMENT"
	BILL          TransactionType = "BILL_AIRTIME"
	DEPOSIT       TransactionType = "DEPOSIT"
	WITHDRAWAL    TransactionType = "WITHDRAWAL"
	REFUND        TransactionType = "REFUND"
)

// TransactionMedium represents the payment medium/provider
// @Description Payment medium or provider (e.g., MPESA, TELEBIRR)
type TransactionMedium string

const (
	SOCIALPAY   TransactionMedium = "SOCIALPAY"
	CYBERSOURCE TransactionMedium = "CYBERSOURCE"
	ETHSWITCH   TransactionMedium = "ETHSWITCH"
	MPESA       TransactionMedium = "MPESA"
	TELEBIRR    TransactionMedium = "TELEBIRR"
	CBE         TransactionMedium = "CBE"
	AWASH       TransactionMedium = "AWASH"
)

// TransactionSource represents the source/origin of the transaction
// @Description Source of the transaction (e.g., DIRECT, QR_PAYMENT)
type TransactionSource string

const (
	DIRECT          TransactionSource = "DIRECT"          // Direct API calls
	HOSTED_CHECKOUT TransactionSource = "HOSTED_CHECKOUT" // Hosted checkout page
	QR_PAYMENT      TransactionSource = "QR_PAYMENT"      // QR payment link
	WITHDRAWAL_TIP  TransactionSource = "WITHDRAWAL"      // Withdrawal/tip transactions
)

// QR transaction tags
const (
	QR_SHOP_PAYMENT       = "QR_SHOP_PAYMENT"
	QR_RESTAURANT_PAYMENT = "QR_RESTAURANT_PAYMENT"
	QR_DONATION_PAYMENT   = "QR_DONATION_PAYMENT"
)

// TransactionDetails contains additional transaction details
// @Description Additional details specific to the transaction
type TransactionDetails struct {

	// Item name
	ItemName string `json:"item_name" example:"Item name"`

	// Item description
	ItemDescription string `json:"item_description" example:"Item description"`

	// Item quantity
	ItemQuantity int `json:"item_quantity" example:"1"`

	// Item price
	ItemId string `json:"item_id" example:"1234567890"`
}

// TransactionRedirects contains URLs for payment redirects
// @Description URLs for redirecting users after payment
type TransactionRedirects struct {
	// URL to redirect on successful payment
	Success string `json:"success" example:"https://example.com/success"`
	// URL to redirect on failed payment
	Failed string `json:"failed" example:"https://example.com/failed"`
}

// TransactionStatus represents the current status of a transaction
// @Description Current status of the transaction
type TransactionStatus string

const (
	INITIATED TransactionStatus = "INITIATED"
	PENDING   TransactionStatus = "PENDING"
	SUCCESS   TransactionStatus = "SUCCESS"
	FAILED    TransactionStatus = "FAILED"
	REFUNDED  TransactionStatus = "REFUNDED"
	EXPIRED   TransactionStatus = "EXPIRED"
	CANCELED  TransactionStatus = "CANCELED"
)

// Transaction represents a payment transaction
// @Description Complete transaction details including payment information
type Transaction struct {
	// Unique identifier for the transaction
	Id uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Merchant ID
	MerchantId uuid.UUID `json:"merchant_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Phone number of the payer
	PhoneNumber string `json:"phone_number" example:"+251911234567"`
	// User ID associated with the transaction
	UserId uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Type of transaction
	Type TransactionType `json:"type" example:"SALE"`
	// Payment medium used
	Medium TransactionMedium `json:"medium" example:"MPESA"`
	// Reference code
	Reference string `json:"reference" example:"REF123456"`
	// Additional comment
	Comment string `json:"comment" example:"Payment for order #123"`
	// Whether the transaction is verified
	Verified bool `json:"verified" example:"true"`
	// Time-to-live in seconds
	TTL int64 `json:"ttl" example:"3600"`
	// Additional transaction details
	Details interface{} `json:"details"`
	// Creation timestamp
	CreatedAt time.Time `json:"created_at" example:"2024-05-03T10:00:00Z"`
	// Last update timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2024-05-03T10:01:00Z"`
	// Confirmation timestamp
	Confirm_Timestamp time.Time `json:"confirm_timestamp" example:"2024-05-03T10:02:00Z"`
	// Reference number
	ReferenceNumber string `json:"reference_number" example:"TXN123456"`
	// Whether this is a test transaction
	Test bool `json:"test" example:"false"`
	// Current status
	Status TransactionStatus `json:"status" example:"SUCCESS"`
	// Transaction description
	Description string `json:"description" example:"Payment for service"`
	// Security token
	Token string `json:"token" example:"tok_123456"`
	// Transaction amount
	Amount float64 `json:"amount" example:"1000.50"`
	// Whether challenge is required
	HasChallenge bool `json:"has_challenge" example:"false"`

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

	// QR Payment Context
	TransactionSource TransactionSource `json:"transaction_source" db:"transaction_source"`
	QRLinkID          *uuid.UUID        `json:"qr_link_id,omitempty" db:"qr_link_id"`
	HostedCheckoutID  *uuid.UUID        `json:"hosted_checkout_id,omitempty" db:"hosted_checkout_id"`
	QRTag             *string           `json:"qr_tag,omitempty" db:"qr_tag"`

	// Tip Information
	HasTip           bool       `json:"has_tip" db:"has_tip"`
	TipAmount        *float64   `json:"tip_amount,omitempty" db:"tip_amount"`
	TipeePhone       *string    `json:"tipee_phone,omitempty" db:"tipee_phone"`
	TipMedium        *string    `json:"tip_medium,omitempty" db:"tip_medium"`
	TipTransactionID *uuid.UUID `json:"tip_transaction_id,omitempty" db:"tip_transaction_id"`
	TipProcessed     bool       `json:"tip_processed" db:"tip_processed"`

	// Merchant information (populated when fetched with merchant details)
	Merchant *entity.Merchant `json:"merchant,omitempty"`
}

// GetTransactionResponse represents the response for transaction lookup
// @Description Response containing transaction details
type GetTransactionResponse struct {
	// Status of the request
	Status string `json:"status" example:"success"`
	// Transaction data
	Data interface{} `json:"data"`
}
