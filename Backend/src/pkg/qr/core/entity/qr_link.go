package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

// QRLinkType represents the type of QR link
type QRLinkType string

const (
	DYNAMIC QRLinkType = "DYNAMIC"
	STATIC  QRLinkType = "STATIC"
)

// QRLinkTag represents the category/tag of the QR link
type QRLinkTag string

const (
	RESTAURANT QRLinkTag = "RESTAURANT"
	DONATION   QRLinkTag = "DONATION"
	SHOP       QRLinkTag = "SHOP"
)

// QRLink represents a QR payment link
type QRLink struct {
	ID               uuid.UUID                  `json:"id" db:"id"`
	UserID           uuid.UUID                  `json:"user_id" db:"user_id"`
	MerchantID       uuid.UUID                  `json:"merchant_id" db:"merchant_id"`
	Type             QRLinkType                 `json:"type" db:"type"`
	Amount           *float64                   `json:"amount,omitempty" db:"amount"` // Only for STATIC type
	SupportedMethods []entity.TransactionMedium `json:"supported_methods" db:"supported_methods"`
	Tag              QRLinkTag                  `json:"tag" db:"tag"`
	Title            *string                    `json:"title,omitempty" db:"title"`
	Description      *string                    `json:"description,omitempty" db:"description"`
	ImageURL         *string                    `json:"image_url,omitempty" db:"image_url"`
	IsTipEnabled     bool                       `json:"is_tip_enabled" db:"is_tip_enabled"`
	IsActive         bool                       `json:"is_active" db:"is_active"`
	CreatedAt        time.Time                  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at" db:"updated_at"`
}

// CreateQRLinkRequest represents request to create a QR link
// @Description Request to create a new QR payment link
type CreateQRLinkRequest struct {
	// Type of QR link (DYNAMIC or STATIC)
	Type QRLinkType `json:"type" example:"DYNAMIC"`

	// Amount for STATIC type QR links
	Amount *float64 `json:"amount,omitempty" example:"100.50"`

	// Supported payment methods
	SupportedMethods []entity.TransactionMedium `json:"supported_methods" example:"[\"MPESA\",\"TELEBIRR\"]"`

	// QR link category
	Tag QRLinkTag `json:"tag" example:"SHOP"`

	// Optional title
	Title *string `json:"title,omitempty" example:"My Shop Payment"`

	// Optional description
	Description *string `json:"description,omitempty" example:"Payment for shop items"`

	// Optional image URL
	ImageURL *string `json:"image_url,omitempty" example:"https://example.com/image.jpg"`

	// Whether tipping is enabled
	IsTipEnabled bool `json:"is_tip_enabled" example:"true"`
}

func (r CreateQRLinkRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Type, validation.Required, validation.In(DYNAMIC, STATIC)),
		validation.Field(&r.SupportedMethods, validation.Required, validation.Length(1, 10)),
		validation.Field(&r.Tag, validation.Required, validation.In(RESTAURANT, DONATION, SHOP)),
		validation.Field(&r.Amount, validation.By(func(value interface{}) error {
			if r.Type == STATIC {
				if value == nil {
					return validation.NewError("validation_required", "amount is required for static QR links")
				}
				if amount, ok := value.(*float64); ok && amount != nil && *amount < 0.01 {
					return validation.NewError("validation_min", "amount must be at least 0.01")
				}
			}
			return nil
		})),
	)
}

// UpdateQRLinkRequest represents request to update a QR link
// @Description Request to update an existing QR payment link
type UpdateQRLinkRequest struct {
	// Amount for STATIC type QR links
	Amount *float64 `json:"amount,omitempty" example:"150.75"`

	// Supported payment methods
	SupportedMethods []entity.TransactionMedium `json:"supported_methods,omitempty"`

	// QR link category
	Tag *QRLinkTag `json:"tag,omitempty" example:"RESTAURANT"`

	// Title
	Title *string `json:"title,omitempty" example:"Updated Shop Payment"`

	// Description
	Description *string `json:"description,omitempty" example:"Updated payment description"`

	// Image URL
	ImageURL *string `json:"image_url,omitempty" example:"https://example.com/new-image.jpg"`

	// Whether tipping is enabled
	IsTipEnabled *bool `json:"is_tip_enabled,omitempty" example:"false"`

	// Whether QR link is active
	IsActive *bool `json:"is_active,omitempty" example:"true"`
}

// QRPaymentRequest represents a payment request via QR link
// @Description Request to make a payment using a QR link
type QRPaymentRequest struct {
	// Amount for DYNAMIC QR links (ignored for STATIC)
	Amount *float64 `json:"amount,omitempty" example:"50.00"`

	// Payment method/provider
	Medium entity.TransactionMedium `json:"medium" example:"MPESA"`

	// Payer's phone number
	PhoneNumber string `json:"phone_number" example:"251911234567"`

	// Optional tip amount
	TipAmount *float64 `json:"tip_amount,omitempty" example:"5.00"`

	// Tipee phone number (required if tip amount > 0)
	TipeePhone *string `json:"tipee_phone,omitempty" example:"251911234568"`

	// Tip payment method (required if tip amount > 0)
	TipMedium *entity.TransactionMedium `json:"tip_medium,omitempty" example:"TELEBIRR"`

	// MerchantPaysFee indicates if the merchant pays the transaction fee
	MerchantPaysFee bool `json:"merchant_pays_fee" example:"false"`
}

func (r QRPaymentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Medium, validation.Required),
		validation.Field(&r.PhoneNumber, validation.Required),
	)
}

// QRLinkResponse represents QR link details response
// @Description Response containing QR link details
type QRLinkResponse struct {
	*QRLink
	// URL to display QR code
	QRCodeURL string `json:"qr_code_url" example:"https://api.socialpay.co/qr/display/123e4567-e89b-12d3-a456-426614174000"`

	// Payment URL for the QR link
	PaymentURL string `json:"payment_url" example:"https://checkout.socialpay.co/qr/123e4567-e89b-12d3-a456-426614174000"`
}

// QRLinksListResponse represents paginated QR links response
// @Description Paginated list of QR links
type QRLinksListResponse struct {
	QRLinks []QRLinkResponse `json:"qr_links"`
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
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

	// Payment amount processed
	PaymentAmount float64 `json:"payment_amount" example:"100.50"`

	// Tip amount (if any)
	TipAmount *float64 `json:"tip_amount,omitempty" example:"5.00"`

	// SocialPay transaction ID for the main payment
	SocialPayTransactionID string `json:"socialpay_transaction_id" example:"123e4567-e89b-12d3-a456-426614174000"`

	// SocialPay transaction ID for the tip (if any)
	TipTransactionID *string `json:"tip_transaction_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"`

	// Payment URL for the QR link
	PaymentURL string `json:"payment_url" example:"https://checkout.socialpay.co/qr/123e4567-e89b-12d3-a456-426614174000"`
}
