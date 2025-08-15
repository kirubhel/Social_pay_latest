package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/ip_whitelist/core/entity"
)

// IPWhitelistRepository defines the interface for IP whitelist operations
type IPWhitelistRepository interface {
	// GetIPWhitelist retrieves the IP whitelist for a merchant
	GetIPWhitelist(ctx context.Context, merchantID uuid.UUID) ([]entity.IPWhitelist, error)

	// CreateIPWhitelist adds new IP address to the whitelist
	CreateIPWhitelist(ctx context.Context, merchantID uuid.UUID, req entity.CreateIPWhitelistRequest) error

	// UpdateIPWhitelist updates the IP whitelist for a merchant
	UpdateIPWhitelist(ctx context.Context, merchantID uuid.UUID, req entity.UpdateIPWhitelistRequest) error

	// DeleteIPWhitelist deletes the IP whitelist for a merchant
	DeleteIPWhitelist(ctx context.Context, id uuid.UUID, merchantID uuid.UUID) error
}
