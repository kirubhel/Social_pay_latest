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
	StatusInactive            MerchantStatus = "inactive"
	StatusSuspended           MerchantStatus = "suspended"
	StatusTerminated          MerchantStatus = "terminated"
)

// SupportedFileType represents file types of exported merchants data
type SupportedFileType string

const (
	CSV  SupportedFileType = "csv"
	XLSX SupportedFileType = "xlsx"
	JSON SupportedFileType = "json"
)

// Merchant represents a merchant in the system
type Merchant struct {
	ID                         uuid.UUID      `json:"id"`
	UserID                     uuid.UUID      `json:"userId"`
	LegalName                  string         `json:"legalName"`
	TradingName                *string        `json:"tradingName,omitempty"`
	BusinessRegistrationNumber string         `json:"businessRegistrationNumber"`
	TaxIdentificationNumber    string         `json:"taxIdentificationNumber"`
	BusinessType               string         `json:"businessType"`
	IndustryCategory           *string        `json:"industryCategory,omitempty"`
	IsBettingCompany           bool           `json:"isBettingCompany"`
	LotteryCertificateNumber   *string        `json:"lotteryCertificateNumber,omitempty"`
	WebsiteURL                 *string        `json:"websiteUrl,omitempty"`
	EstablishedDate            *time.Time     `json:"establishedDate,omitempty"`
	CreatedAt                  time.Time      `json:"createdAt"`
	UpdatedAt                  time.Time      `json:"updatedAt"`
	Status                     MerchantStatus `json:"status"`
}

// MerchantAddress represents a merchant address
type MerchantAddress struct {
	ID             uuid.UUID `json:"id"`
	MerchantID     uuid.UUID `json:"merchantId"`
	AddressType    string    `json:"addressType"`
	StreetAddress1 string    `json:"streetAddress1"`
	StreetAddress2 *string   `json:"streetAddress2,omitempty"`
	City           string    `json:"city"`
	Region         string    `json:"region"`
	PostalCode     *string   `json:"postalCode,omitempty"`
	Country        string    `json:"country"`
	IsPrimary      bool      `json:"isPrimary"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// MerchantContact represents a merchant contact
type MerchantContact struct {
	ID          uuid.UUID `json:"id"`
	MerchantID  uuid.UUID `json:"merchantId"`
	ContactType string    `json:"contactType"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Position    *string   `json:"position,omitempty"`
	IsVerified  bool      `json:"isVerified"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// MerchantDocument represents a merchant document
type MerchantDocument struct {
	ID              uuid.UUID  `json:"id"`
	MerchantID      uuid.UUID  `json:"merchantId"`
	DocumentType    string     `json:"documentType"`
	DocumentNumber  *string    `json:"documentNumber,omitempty"`
	FileURL         string     `json:"fileUrl"`
	FileHash        *string    `json:"fileHash,omitempty"`
	VerifiedBy      *uuid.UUID `json:"verifiedBy,omitempty"`
	VerifiedAt      *time.Time `json:"verifiedAt,omitempty"`
	Status          string     `json:"status"`
	RejectionReason *string    `json:"rejectionReason,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// MerchantBankAccount represents a merchant bank account
type MerchantBankAccount struct {
	ID                     uuid.UUID  `json:"id"`
	MerchantID             uuid.UUID  `json:"merchantId"`
	AccountHolderName      string     `json:"accountHolderName"`
	BankName               string     `json:"bankName"`
	BankCode               string     `json:"bankCode"`
	BranchCode             *string    `json:"branchCode,omitempty"`
	AccountNumber          string     `json:"accountNumber"`
	AccountType            string     `json:"accountType"`
	Currency               string     `json:"currency"`
	IsPrimary              bool       `json:"isPrimary"`
	IsVerified             bool       `json:"isVerified"`
	VerificationDocumentID *uuid.UUID `json:"verificationDocumentId,omitempty"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}

// MerchantSettings represents merchant settings
type MerchantSettings struct {
	MerchantID          uuid.UUID `json:"merchantId"`
	DefaultCurrency     string    `json:"defaultCurrency"`
	DefaultLanguage     string    `json:"defaultLanguage"`
	CheckoutTheme       *string   `json:"checkoutTheme,omitempty"`
	EnableWebhooks      bool      `json:"enableWebhooks"`
	WebhookURL          *string   `json:"webhookUrl,omitempty"`
	WebhookSecret       *string   `json:"webhookSecret,omitempty"`
	AutoSettlement      bool      `json:"autoSettlement"`
	SettlementFrequency string    `json:"settlementFrequency"`
	RiskSettings        *string   `json:"riskSettings,omitempty"` // JSON string
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
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

// MerchantsResponse contains extra merchant data from other tables
type MerchantsResponse struct {
	Count     int               `json:"count"`
	Merchants []MerchantDetails `json:"merchants"`
}

// GetMerchantsParams contians filters used to fetch list of merchants
type GetMerchantsParams struct {
	Text      string    `json:"text"`
	Skip      int       `json:"skip"`
	Take      int       `json:"take"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Status    string    `json:"status"`
}

// ExportMerchantsRequest contains export request body
type ExportMerchantsRequest struct {
	FileType  SupportedFileType `json:"fileType"`
	Data      []string          `json:"data"`
	Merchants []uuid.UUID       `json:"merchants"`
}

// UpdateMerchantBusinessInformationRequest contains merchant business info
type UpdateMerchantBusinessInformationRequest struct {
	LegalName                  *string         `json:"legal_name"`
	TradingName                *string         `json:"trading_name"`
	BusinessRegistrationNumber *string         `json:"business_registration_number"`
	TaxIdentificationNumber    *string         `json:"tax_identification_number"`
	BusinessType               *string         `json:"business_type"`
	IndustryCategory           *string         `json:"industry_category"`
	IsBettingCompany           *bool           `json:"is_betting_company"`
	LotteryCertificateNumber   *string         `json:"lottery_certificate_number"`
	WebsiteURL                 *string         `json:"website_url"`
	EstablishedDate            *time.Time      `json:"established_date"`
	Status                     *MerchantStatus `json:"status"`
	Mode                       *string         `json:"mode"`
}

// CreateMerchantDocumnetRequest contains create merchant document request
type CreateMerchantDocumnetRequest struct {
	DocumentType string `json:"document_type"`
	Status       string `json:"status"`
	FileUrl      string `json:"file_url"`
}

// UpdateMerchantDocumentRequest contains update merchant document request with ID
type UpdateMerchantDocumentWithIDRequest struct {
	ID           uuid.UUID `json:"id"`
	DocumentType string    `json:"document_type"`
	Status       string    `json:"status"`
	FileUrl      string    `json:"file_url"`
}

// CreateMerchantPersonalInformationRequest contains create merchant personal info
type CreateMerchantPersonalInformationRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// UpdateMerchantRequest contains update merchant request body
type UpdateMerchantRequest struct {
	BusinessInfo UpdateMerchantBusinessInformationRequest `json:"business_info"`
	PersonalInfo CreateMerchantPersonalInformationRequest `json:"personal_info"`
	Documents    []UpdateMerchantDocumentWithIDRequest    `json:"documents"`
}

// UpdateMerchantContactRequest contains update merchant contact request body
type UpdateMerchantContactRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	IsVerified  bool   `json:"is_verified"`
}

// UpdateMerchantDocumentRequest contains update merchant document request body
type UpdateMerchantDocumentRequest struct {
	FileUrl         string     `json:"file_url"`
	VerifiedBy      *uuid.UUID `json:"verified_by"`
	Status          string     `json:"status"`
	RejectionReason *string    `json:"rejection_reason"`
}

// UpdateMerchantStatusRequest contains update merchant status request body
type UpdateMerchantStatusRequest struct {
	Status string `json:"status"`
}

// DeleteMerchantsRequest contains delete merchants request body
type DeleteMerchantsRequest struct {
	IDs []uuid.UUID `json:"ids"`
}

// MerchantStats represents merchant statistics
type MerchantStats struct {
	TotalMerchants  int64 `json:"total_merchants"`
	ActiveMerchants int64 `json:"active_merchants"`
	PendingKyc      int64 `json:"pending_kyc"`
	NewThisMonth    int64 `json:"new_this_month"`
}

// ImpersonateMerchantRequest represents impersonate merchant request
type ImpersonateMerchantRequest struct {
	MerchantID uuid.UUID `json:"merchant_id"`
}
