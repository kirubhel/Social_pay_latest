package entity

import (
	"time"

	"github.com/google/uuid"
)

// TwoFactorCode represents a 2FA verification code
type TwoFactorCode struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Code      string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TwoFactorStatus represents the 2FA status for a user
type TwoFactorStatus struct {
	Enabled    bool       `json:"enabled"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}
