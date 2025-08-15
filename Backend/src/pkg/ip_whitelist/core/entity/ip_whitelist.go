package entity

import (
	"time"

	"github.com/google/uuid"
)

// IPWhitelist represents a merchant's whitelisted IP addresses
type IPWhitelist struct {
	ID        uuid.UUID `json:"id"`
	IpAddress string    `json:"ip_address"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateIPWhitelistRequest contains create IP Whitelist request body
type CreateIPWhitelistRequest struct {
	IPAddress string `json:"ip_address"`
}

// UpdateIPWhitelistRequest contains create IP Whitelist request body
type UpdateIPWhitelistRequest struct {
	ID        uuid.UUID `json:"id"`
	IPAddress string    `json:"ip_address"`
	IsActive  bool      `json:"is_active"`
}
