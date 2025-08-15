package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/google/uuid"
	ip_whitelist "github.com/socialpay/socialpay/src/pkg/ip_whitelist/adapter/gateway/repository/generated"
	"github.com/socialpay/socialpay/src/pkg/ip_whitelist/core/entity"
	"github.com/socialpay/socialpay/src/pkg/ip_whitelist/core/repository"
	"github.com/sqlc-dev/pqtype"
)

type ipWhitelistRepositoryImpl struct {
	db *ip_whitelist.Queries
}

func NewIPWhitelistRepository(db *sql.DB) repository.IPWhitelistRepository {
	return &ipWhitelistRepositoryImpl{
		db: ip_whitelist.New(db),
	}
}

func (r *ipWhitelistRepositoryImpl) GetIPWhitelist(ctx context.Context, merchantID uuid.UUID) ([]entity.IPWhitelist, error) {
	whitelistedIPs, err := r.db.GetWhitelistedIPsByMerchantID(ctx, merchantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]entity.IPWhitelist, 0), nil
		}
		return nil, fmt.Errorf("failed to get IP whitelist: %w", err)
	}

	var ipWhitelists []entity.IPWhitelist
	for _, ip := range whitelistedIPs {
		ipWhitelists = append(ipWhitelists, entity.IPWhitelist{
			ID:        ip.ID,
			IpAddress: ip.IpAddress.IPNet.String(),
			IsActive:  ip.IsActive,
			CreatedAt: ip.CreatedAt,
			UpdatedAt: ip.UpdatedAt,
		})
	}

	return ipWhitelists, nil
}

func (r *ipWhitelistRepositoryImpl) CreateIPWhitelist(ctx context.Context, merchantID uuid.UUID, req entity.CreateIPWhitelistRequest) error {
	// Parse the IP string to ensure it's valid
	_, ipNet, err := net.ParseCIDR(req.IPAddress)
	if err != nil {
		return fmt.Errorf("invalid IP address format %s: %w", req.IPAddress, err)
	}

	// Convert net.IPNet to pqtype.CIDR
	cidr := pqtype.CIDR{
		IPNet: *ipNet,
		Valid: true,
	}

	whitelisted, err := r.db.CheckIPWhitelisted(ctx, ip_whitelist.CheckIPWhitelistedParams{
		MerchantID: merchantID,
		Column2:    pqtype.Inet(cidr),
	})

	if err != nil {
		return fmt.Errorf("Error checking ip address is whitelisted: %w", err)
	}

	if whitelisted {
		return fmt.Errorf("IP address already whitelisted")
	}

	id := uuid.New()
	err = r.db.CreateWhitelistedIP(ctx, ip_whitelist.CreateWhitelistedIPParams{
		ID:          id,
		MerchantID:  merchantID,
		IpAddress:   cidr,
		Description: sql.NullString{String: "Added via API", Valid: true},
		IsActive:    true,
	})

	if err != nil {
		return fmt.Errorf("Failed to create ip whitelist: %w", err)
	}

	return nil
}

func (r *ipWhitelistRepositoryImpl) UpdateIPWhitelist(ctx context.Context, merchantID uuid.UUID, req entity.UpdateIPWhitelistRequest) error {
	ipWhitelist, err := r.db.GetWhitelistedIPByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("IP Whitelist not found")
		}
		return fmt.Errorf("Error getting IP Whitelist: %w", err)
	}

	if ipWhitelist.MerchantID != merchantID {
		return fmt.Errorf("Invalid action!!")
	}

	// Parse the IP string to ensure it's valid
	_, ipNet, err := net.ParseCIDR(req.IPAddress)
	if err != nil {
		return fmt.Errorf("invalid IP address format %s: %w", req.IPAddress, err)
	}

	// Convert net.IPNet to pqtype.CIDR
	cidr := pqtype.CIDR{
		IPNet: *ipNet,
		Valid: true,
	}
	err = r.db.UpdateWhitelistedIP(ctx, ip_whitelist.UpdateWhitelistedIPParams{
		ID:        req.ID,
		IpAddress: cidr,
		IsActive:  req.IsActive,
	})

	if err != nil {
		return fmt.Errorf("failed to update ip whitelist: %w", err)
	}

	return nil
}

func (r *ipWhitelistRepositoryImpl) DeleteIPWhitelist(ctx context.Context, id uuid.UUID, merchantID uuid.UUID) error {
	ipWhitelist, err := r.db.GetWhitelistedIPByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("IP Whitelist not found")
		}
		return fmt.Errorf("Error getting IP Whitelist: %w", err)
	}

	if ipWhitelist.MerchantID != merchantID {
		return fmt.Errorf("Invalid action!!")
	}

	err = r.db.DeleteWhitelistedIP(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete IP %s: %w", ipWhitelist.IpAddress.IPNet.String(), err)
	}

	return nil
}
