package entity

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/socialpay/socialpay/src/pkg/merchants/errors"

	"github.com/google/uuid"
)

type MerchantAdditionalInfo struct {
	MerchantID           uuid.UUID `json:"merchant_id,omitempty"`            // UUID
	PersonalName         string    `json:"personal_name"`                    // Varying character
	PhoneNumber          string    `json:"phone_number"`                     // Varying character
	Region               string    `json:"region"`                           // Varying character
	City                 string    `json:"city"`                             // Varying character
	SubCity              string    `json:"sub_city"`                         // Varying character
	Woreda               string    `json:"woreda"`                           // Varying character
	PostalCode           string    `json:"postal_code"`                      // Varying character
	SecondaryPhoneNumber string    `json:"secondary_phone_number,omitempty"` // Varying character
	Email                string    `json:"email"`                            // Varying character
}

func ValidateMerchantAdditionalInfo(info MerchantAdditionalInfo) error {
	fieldErrors := map[string]string{}

	if strings.TrimSpace(info.PersonalName) == "" {
		fieldErrors["personal_name"] = "personal_name is required"
	}
	if strings.TrimSpace(info.PhoneNumber) == "" {
		fieldErrors["phone_number"] = "phone_number is required"
	}
	if strings.TrimSpace(info.Region) == "" {
		fieldErrors["region"] = "region is required"
	}
	if strings.TrimSpace(info.City) == "" {
		fieldErrors["city"] = "city is required"
	}
	if strings.TrimSpace(info.SubCity) == "" {
		fieldErrors["sub_city"] = "sub_city is required"
	}
	if strings.TrimSpace(info.Woreda) == "" {
		fieldErrors["woreda"] = "woreda is required"
	}
	if strings.TrimSpace(info.Email) == "" {
		fieldErrors["email"] = "email is required"
	} else if _, err := mail.ParseAddress(info.Email); err != nil {
		fieldErrors["email"] = "invalid email format"
	}

	if info.PhoneNumber != "" && !ethiopianPhoneRegex.MatchString(info.PhoneNumber) {
		fieldErrors["phone_number"] = "invalid Ethiopian phone number format (expected +2519XXXXXXXX)"
	}

	if info.SecondaryPhoneNumber != "" && !ethiopianPhoneRegex.MatchString(info.SecondaryPhoneNumber) {
		fieldErrors["secondary_phone_number"] = "invalid Ethiopian phone number format"
	}

	if len(info.PersonalName) > 100 {
		fieldErrors["personal_name"] = "personal_name is too long (max 100 characters)"
	}
	if len(info.PostalCode) > 20 {
		fieldErrors["postal_code"] = "postal_code is too long (max 20 characters)"
	}

	if len(fieldErrors) > 0 {
		return errors.Error{
			Type:    "VALIDATION_ERROR",
			Message: "validation failed",
			Code:    http.StatusBadRequest,
			Data:    fieldErrors,
		}
	}

	return nil
}
