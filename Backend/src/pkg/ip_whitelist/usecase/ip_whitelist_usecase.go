package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/ip_whitelist/core/entity"
	"github.com/socialpay/socialpay/src/pkg/ip_whitelist/core/repository"
	"github.com/socialpay/socialpay/src/pkg/ip_whitelist/utils"
)

// IPWhitelistUseCase defines the interface for IP whitelist operations
type IPWhitelistUseCase interface {
	// GetIPWhitelist retrieves the IP whitelist for a merchant
	GetIPWhitelist(ctx context.Context, merchantID uuid.UUID) ([]entity.IPWhitelist, error)

	// CreateIPWhitelist whitelists new IP address
	CreateIPWhitelist(ctx context.Context, merchantID uuid.UUID, req entity.CreateIPWhitelistRequest) error

	// UpdateIPWhitelist updates the IP whitelist for a merchant
	UpdateIPWhitelist(ctx context.Context, merchantID uuid.UUID, request entity.UpdateIPWhitelistRequest) error

	// DeleteIPWhitelist deletes the IP whitelist for a merchant
	DeleteIPWhitelist(ctx context.Context, id uuid.UUID, merchantID uuid.UUID) error
}

type ipWhitelistUseCase struct {
	repo repository.IPWhitelistRepository
}

// NewIPWhitelistUseCase creates a new instance of IPWhitelistUseCase
func NewIPWhitelistUseCase(repo repository.IPWhitelistRepository) IPWhitelistUseCase {
	return &ipWhitelistUseCase{
		repo: repo,
	}
}

func (u *ipWhitelistUseCase) GetIPWhitelist(ctx context.Context, merchantID uuid.UUID) ([]entity.IPWhitelist, error) {
	return u.repo.GetIPWhitelist(ctx, merchantID)
}

func (u *ipWhitelistUseCase) CreateIPWhitelist(ctx context.Context, merchantID uuid.UUID, req entity.CreateIPWhitelistRequest) error {
	err := utils.ValidateIPAddress(req.IPAddress)
	if err != nil {
		return err
	}
	return u.repo.CreateIPWhitelist(ctx, merchantID, req)
}

func (u *ipWhitelistUseCase) UpdateIPWhitelist(ctx context.Context, merchantID uuid.UUID, request entity.UpdateIPWhitelistRequest) error {
	if err := utils.ValidateIPAddress(request.IPAddress); err != nil {
		return err
	}

	if err := u.repo.UpdateIPWhitelist(ctx, merchantID, request); err != nil {
		return err
	}

	return nil
}

func (u *ipWhitelistUseCase) DeleteIPWhitelist(ctx context.Context, id uuid.UUID, merchantID uuid.UUID) error {
	return u.repo.DeleteIPWhitelist(ctx, id, merchantID)
}
