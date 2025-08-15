package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	STORED AccountType = "STORED"
	BANK   AccountType = "BANK"
	CARD   AccountType = "CARD"
	WALLET AccountType = "WALLET"
)

type Account struct {
	ID                 uuid.UUID   `json:"id"`
	Title              string      `json:"title"`
	Type               AccountType `json:"type"`
	Default            bool        `json:"default"`
	UserID             uuid.UUID   `json:"user_id"`
	VerificationStatus struct {
		Verified   bool `json:"verified"`
		VerifiedBy *struct {
			Method  string      `json:"method"`
			Details interface{} `json:"details"`
		} `json:"verified_by,omitempty"`
	} `json:"verification_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Account_VerificationStatus struct {
	Verified   bool `json:"verified"`
	VerifiedBy *struct {
		Method  string      `json:"method"`
		Details interface{} `json:"details"`
	} `json:"verified_by,omitempty"`
}
