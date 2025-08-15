package entity

import (
	"time"

	"github.com/google/uuid"
)

// APIKeyPermission represents a permission that can be granted to an API key
type APIKeyPermission string

const (
	// Define standard permissions
	PermissionRead   APIKeyPermission = "read"
	PermissionWrite  APIKeyPermission = "write"
	PermissionDelete APIKeyPermission = "delete"
	PermissionAdmin  APIKeyPermission = "admin"
)

// CustomTime is a wrapper around time.Time that handles empty string JSON unmarshaling
type CustomTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "\"\"" || s == "null" {
		return nil
	}

	// Try to parse the time
	var err error
	t, err := time.Parse(`"2006-01-02T15:04:05Z07:00"`, s)
	if err != nil {
		return nil // Return nil to make it optional
	}
	ct.Time = t
	return nil
}

// APIKey represents an API key in the system
type APIKey struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"user_id"`
	MerchantID        uuid.UUID  `json:"merchant_id"`
	CreatedBy         uuid.UUID  `json:"created_by"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	PublicKey         string     `json:"public_key"`
	SecretKey         string     `json:"secret_key"`
	CanWithdrawal     bool       `json:"can_withdrawal"`
	CanProcessPayment bool       `json:"can_process_payment"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	LastUsedAt        *time.Time `json:"last_used_at,omitempty"`
	IsActive          bool       `json:"is_active"`
}

// APIKeyResponse represents an API key response without sensitive data
type APIKeyResponse struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"user_id"`
	MerchantID        uuid.UUID  `json:"merchant_id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	PublicKey         string     `json:"public_key"`
	SecretKey         string     `json:"secret_key"`
	CanWithdrawal     bool       `json:"can_withdrawal"`
	CanProcessPayment bool       `json:"can_process_payment"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	LastUsedAt        *time.Time `json:"last_used_at,omitempty"`
	IsActive          bool       `json:"is_active"`
}

// CreateAPIKeyRequest represents the data needed to create a new API key
type CreateAPIKeyRequest struct {
	Name              string      `json:"name" validate:"required"`
	Description       string      `json:"description"`
	CanWithdrawal     *bool       `json:"can_withdrawal"`
	CanProcessPayment *bool       `json:"can_process_payment"`
	ExpiresAt         *CustomTime `json:"expires_at,omitempty"`
}

// GetExpiresAt returns the actual time.Time pointer from CustomTime
func (r *CreateAPIKeyRequest) GetExpiresAt() *time.Time {
	if r.ExpiresAt == nil {
		return nil
	}
	if r.ExpiresAt.IsZero() {
		return nil
	}
	return &r.ExpiresAt.Time
}

// UpdateAPIKeyRequest represents the data that can be updated for an API key
type UpdateAPIKeyRequest struct {
	Name              *string     `json:"name,omitempty"`
	Description       *string     `json:"description,omitempty"`
	CanWithdrawal     *bool       `json:"can_withdrawal,omitempty"`
	CanProcessPayment *bool       `json:"can_process_payment,omitempty"`
	IsActive          *bool       `json:"is_active,omitempty"`
	ExpiresAt         *CustomTime `json:"expires_at,omitempty"`
}

// GetExpiresAt returns the actual time.Time pointer from CustomTime
func (r *UpdateAPIKeyRequest) GetExpiresAt() *time.Time {
	if r.ExpiresAt == nil {
		return nil
	}
	if r.ExpiresAt.IsZero() {
		return nil
	}
	return &r.ExpiresAt.Time
}

// APIKeyRotateResponse represents the response when rotating an API key's secret
type APIKeyRotateResponse struct {
	APIKey    *APIKey `json:"api_key"`
	SecretKey string  `json:"secret_key"`
}

// APIKeyValidateRequest represents the data needed to validate an API key
type APIKeyValidateRequest struct {
	PublicKey string `json:"public_key" validate:"required"`
	SecretKey string `json:"secret_key" validate:"required"`
}
