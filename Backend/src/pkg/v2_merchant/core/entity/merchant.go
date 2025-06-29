package entity

import (
	"time"

	"github.com/google/uuid"
)

// MerchantStatus represents the status of a merchant
type MerchantStatus string

const (
	StatusPendingVerification MerchantStatus = "pending_verification"
	StatusActive              MerchantStatus = "active"
	StatusSuspended           MerchantStatus = "suspended"
	StatusTerminated          MerchantStatus = "terminated"
)

// Merchant represents a merchant in the system
type Merchant struct {
	ID                         uuid.UUID      `json:"id"`
	UserID                     uuid.UUID      `json:"user_id"`
	LegalName                  string         `json:"legal_name"`
	TradingName                *string        `json:"trading_name,omitempty"`
	BusinessRegistrationNumber string         `json:"business_registration_number"`
	TaxIdentificationNumber    string         `json:"tax_identification_number"`
	BusinessType               string         `json:"business_type"`
	IndustryCategory           *string        `json:"industry_category,omitempty"`
	IsBettingCompany           bool           `json:"is_betting_company"`
	LotteryCertificateNumber   *string        `json:"lottery_certificate_number,omitempty"`
	WebsiteURL                 *string        `json:"website_url,omitempty"`
	EstablishedDate            *time.Time     `json:"established_date,omitempty"`
	CreatedAt                  time.Time      `json:"created_at"`
	UpdatedAt                  time.Time      `json:"updated_at"`
	Status                     MerchantStatus `json:"status"`
}

// MerchantAddress represents a merchant address
type MerchantAddress struct {
	ID             uuid.UUID `json:"id"`
	MerchantID     uuid.UUID `json:"merchant_id"`
	AddressType    string    `json:"address_type"`
	StreetAddress1 string    `json:"street_address_1"`
	StreetAddress2 *string   `json:"street_address_2,omitempty"`
	City           string    `json:"city"`
	Region         string    `json:"region"`
	PostalCode     *string   `json:"postal_code,omitempty"`
	Country        string    `json:"country"`
	IsPrimary      bool      `json:"is_primary"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// MerchantContact represents a merchant contact
type MerchantContact struct {
	ID          uuid.UUID `json:"id"`
	MerchantID  uuid.UUID `json:"merchant_id"`
	ContactType string    `json:"contact_type"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Position    *string   `json:"position,omitempty"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MerchantDocument represents a merchant document
type MerchantDocument struct {
	ID              uuid.UUID  `json:"id"`
	MerchantID      uuid.UUID  `json:"merchant_id"`
	DocumentType    string     `json:"document_type"`
	DocumentNumber  *string    `json:"document_number,omitempty"`
	FileURL         string     `json:"file_url"`
	FileHash        *string    `json:"file_hash,omitempty"`
	VerifiedBy      *uuid.UUID `json:"verified_by,omitempty"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	Status          string     `json:"status"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// MerchantBankAccount represents a merchant bank account
type MerchantBankAccount struct {
	ID                     uuid.UUID  `json:"id"`
	MerchantID             uuid.UUID  `json:"merchant_id"`
	AccountHolderName      string     `json:"account_holder_name"`
	BankName               string     `json:"bank_name"`
	BankCode               string     `json:"bank_code"`
	BranchCode             *string    `json:"branch_code,omitempty"`
	AccountNumber          string     `json:"account_number"`
	AccountType            string     `json:"account_type"`
	Currency               string     `json:"currency"`
	IsPrimary              bool       `json:"is_primary"`
	IsVerified             bool       `json:"is_verified"`
	VerificationDocumentID *uuid.UUID `json:"verification_document_id,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

// MerchantSettings represents merchant settings
type MerchantSettings struct {
	MerchantID          uuid.UUID `json:"merchant_id"`
	DefaultCurrency     string    `json:"default_currency"`
	DefaultLanguage     string    `json:"default_language"`
	CheckoutTheme       *string   `json:"checkout_theme,omitempty"`
	EnableWebhooks      bool      `json:"enable_webhooks"`
	WebhookURL          *string   `json:"webhook_url,omitempty"`
	WebhookSecret       *string   `json:"webhook_secret,omitempty"`
	AutoSettlement      bool      `json:"auto_settlement"`
	SettlementFrequency string    `json:"settlement_frequency"`
	RiskSettings        *string   `json:"risk_settings,omitempty"` // JSON string
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// MerchantDetails represents complete merchant information with related data
type MerchantDetails struct {
	Merchant     Merchant              `json:"merchant"`
	Addresses    []MerchantAddress     `json:"addresses,omitempty"`
	Contacts     []MerchantContact     `json:"contacts,omitempty"`
	Documents    []MerchantDocument    `json:"documents,omitempty"`
	BankAccounts []MerchantBankAccount `json:"bank_accounts,omitempty"`
	Settings     *MerchantSettings     `json:"settings,omitempty"`
}

// MerchantResponse represents a merchant response for API
type MerchantResponse struct {
	ID                         uuid.UUID      `json:"id"`
	UserID                     uuid.UUID      `json:"user_id"`
	LegalName                  string         `json:"legal_name"`
	TradingName                *string        `json:"trading_name,omitempty"`
	BusinessRegistrationNumber string         `json:"business_registration_number"`
	TaxIdentificationNumber    string         `json:"tax_identification_number"`
	BusinessType               string         `json:"business_type"`
	IndustryCategory           *string        `json:"industry_category,omitempty"`
	IsBettingCompany           bool           `json:"is_betting_company"`
	LotteryCertificateNumber   *string        `json:"lottery_certificate_number,omitempty"`
	WebsiteURL                 *string        `json:"website_url,omitempty"`
	EstablishedDate            *time.Time     `json:"established_date,omitempty"`
	CreatedAt                  time.Time      `json:"created_at"`
	UpdatedAt                  time.Time      `json:"updated_at"`
	Status                     MerchantStatus `json:"status"`
}
