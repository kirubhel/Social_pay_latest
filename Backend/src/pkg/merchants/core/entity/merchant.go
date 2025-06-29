// models/key.go
package entity

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/socialpay/socialpay/src/pkg/merchants/errors"

	"github.com/google/uuid"
)

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
    MerchantID               uuid.UUID                `json:"merchant_id"`
    UserID                   uuid.UUID                `json:"user_id"`
    LegalName                string                   `json:"legal_name"`
    TradingName              string                   `json:"trading_name"`
    BusinessRegNumber        string                   `json:"business_registration_number"`
    TaxIdentifier           string                   `json:"tax_identification_number"`
    IndustryType            string                   `json:"industry_category"`
    BusinessType            string                   `json:"business_type"`
    IsBettingClient         bool                     `json:"is_betting_company"`
    LoyaltyCertificate      string                   `json:"lottery_certificate_number"`
    WebsiteURL              string                   `json:"website_url"`
    EstablishedDate         time.Time                `json:"established_date"`
    CreatedAt               time.Time                `json:"created_at"`
    UpdatedAt               time.Time                `json:"updated_at"`
    Status                  MerchantStatus           `json:"status"`
    Address                 *MerchantAdditionalInfo  `json:"address,omitempty"`
    Documents               []MerchantDocument       `json:"documents,omitempty"`
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
		return errors.Error{
			Type:    "BADREQUEST",
			Message: "validation failed",
			Code:    http.StatusBadRequest,
			Data:    validationErrors,
		}

	}

	return nil
}
