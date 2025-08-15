package entity

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrGatewayNotFound = errors.New("gateway not found")
var ErrGatewayInUse = errors.New("gateway is currently in use")
var ErrInvalidConfig = errors.New("invalid gateway configuration")
var ErrMerchantNotFound = errors.New("merchant not found")
var ErrGatewayAlreadyExists = errors.New("gateway already exists")
var ErrAlreadyLinked = errors.New("gateway is already linked to the merchant")
var ErrNotLinked = errors.New("gateway is not linked to the merchant")
var ErrGatewayDisabled = errors.New("gateway is disabled")
var ErrAlreadyDisabled = errors.New("gateway is already disabled")
var ErrAlreadyEnabled = errors.New("gateway is already enabled")

type GatewayMerchantsResult struct {
	Merchants   []GatewayMerchant
	TotalCount  int
	ActiveCount int
}

type PaymentGateway struct {
	ID             uuid.UUID         `json:"id" db:"id"`
	Name           string            `json:"name" db:"name"`
	Description    string            `json:"description" db:"description"`
	Type           string            `json:"type" db:"type"` // e.g., "bank", "mobile_money", "card_processor"
	IsActive       bool              `json:"is_active" db:"is_active"`
	Config         GatewayConfig     `json:"config" db:"config"` // JSON configuration
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`
	LinkedAt       time.Time         `json:"linked_at,omitempty" db:"linked_at"`
	DisabledReason string            `json:"disabled_reason,omitempty" db:"disabled_reason"`
	DisabledAt     *time.Time        `json:"disabled_at,omitempty" db:"disabled_at"`
	Merchant       Merchant          `json:"merchant,omitempty" db:"merchant"`   // Optional, if linked to a merchant
	Merchants      []GatewayMerchant `json:"merchants,omitempty" db:"merchants"` // List of merchants linked to this gateway
}

type ListPaymentGateway struct {
	ID             uuid.UUID         `json:"id" db:"id"`
	Name           string            `json:"name" db:"name"`
	Description    string            `json:"description" db:"description"`
	Type           string            `json:"type" db:"type"` // e.g., "bank", "mobile_money", "card_processor"
	IsActive       bool              `json:"is_active" db:"is_active"`
	Config         GatewayConfig     `json:"config" db:"config"` // JSON configuration
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`
}

type GatewayConfig struct {
	APIKey         string  `json:"api_key,omitempty"`
	SecretKey      string  `json:"secret_key,omitempty"`
	BaseURL        string  `json:"base_url"`
	WebhookURL     string  `json:"webhook_url"`
	TransactionFee float64 `json:"transaction_fee"`
	IsTest         bool    `json:"is_test"`
}

type GatewayMerchant struct {
	MerchantID     uuid.UUID  `json:"merchant_id"`
	Name           string     `json:"name"`
	BusinessID     string     `json:"business_id"`
	Type           string     `json:"type"` // e.g., "bank", "mobile_money", "card_processor"
	IsActive       bool       `json:"is_active"`
	LinkedAt       time.Time  `json:"linked_at"`
	DisabledAt     *time.Time `json:"disabled_at,omitempty"`
	DisabledReason *string    `json:"disabled_reason,omitempty"`
	Merchants      Merchant
}

var ethiopianPhoneRegex = regexp.MustCompile(`^\+251(9|7|1)\d{8}$`)

type MerchantStatus string

const (
	StatusActive    MerchantStatus = "active"
	StatusPending   MerchantStatus = "pending"
	StatusSuspended MerchantStatus = "suspended"
)

type FieldErrors map[string]string

type Merchant struct {
	MerchantID         string    `json:"merchant_id"`         // UUID
	UserID             uuid.UUID `json:"user_id"`             // UUID
	LegalName          string    `json:"legal_name"`          // Varying character
	TradingName        string    `json:"trading_name"`        // Varying character
	BusinessRegNumber  string    `json:"business_reg_number"` // Varying character
	TaxIdentifier      string    `json:"tax_identifier"`      // Varying character
	IndustryType       string    `json:"industry_type"`       // Varying character
	BusinessType       string    `json:"business_type"`       // Varying character
	IsBettingClient    bool      `json:"is_betting_client"`   // Boolean
	LoyaltyCertificate string    `json:"loyalty_certificate"` // Varying character
	WebsiteURL         string    `json:"website_url"`         // Varying character
	EstablishedDate    time.Time `json:"established_date"`    // Date
	CreatedAt          time.Time `json:"created_at"`          // Timestamp
	UpdatedAt          time.Time `json:"updated_at"`          // Timestamp
	Status             string    `json:"status"`              // Varying character
}

type MerchantDetails struct {
	MerchantID         uuid.UUID      `json:"merchant_id"`
	UserID             uuid.UUID      `json:"user_id"`
	LegalName          string         `json:"legal_name"`
	TradingName        string         `json:"trading_name"`
	BusinessRegNumber  string         `json:"business_registration_number"`
	TaxIdentifier      string         `json:"tax_identification_number"`
	IndustryType       string         `json:"industry_category"`
	BusinessType       string         `json:"business_type"`
	IsBettingClient    bool           `json:"is_betting_company"`
	LoyaltyCertificate string         `json:"lottery_certificate_number"`
	WebsiteURL         string         `json:"website_url"`
	EstablishedDate    time.Time      `json:"established_date"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	Status             MerchantStatus `json:"status"`
}

func (fe FieldErrors) Error() string {

	return "validation failed "

}

func (m Merchant) Validate() error {

	validationErrors := FieldErrors{}

	if strings.TrimSpace(m.LegalName) == "" {
		validationErrors["legal_name"] = "legal_name is required"
	}

	if strings.TrimSpace(m.TradingName) == "" {
		validationErrors["trading_name"] = "trading_name is required"
	}

	if strings.TrimSpace(m.BusinessRegNumber) == "" {
		validationErrors["business_reg_number"] = "business_reg_number is required"
	}

	if strings.TrimSpace(m.TaxIdentifier) == "" {
		validationErrors["tex_indentifier"] = "tax_identifier is required"
	}

	if strings.TrimSpace(m.IndustryType) == "" {
		validationErrors["industry_type"] = "industry_type is required"
	}

	if strings.TrimSpace(m.BusinessType) == "" {
		validationErrors["business_type"] = "business_type is required"
	}

	if strings.TrimSpace(m.LoyaltyCertificate) == "" {
		validationErrors["loyalty_certificate"] = "loyalty_certificate is required"
	}

	if !m.EstablishedDate.Before(time.Now()) {
		validationErrors["established_date"] = "established_date must be in the past"
	}

	if m.WebsiteURL != "" {
		if _, err := url.ParseRequestURI(m.WebsiteURL); err != nil {
			validationErrors["website_url"] = "website_url must be a valid URL"
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation failed: %v", validationErrors)

	}

	return nil
}
