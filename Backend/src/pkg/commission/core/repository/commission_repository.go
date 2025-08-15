package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/commission/core/entity"
)

// TODO setup caching
type CommissionRepository interface {
	// Get default commission settings
	GetDefaultCommission(ctx context.Context) (*entity.CommissionSettings, error)

	// Get merchant commission settings
	GetMerchantCommission(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantCommission, error)

	// Update merchant commission settings
	UpdateMerchantCommission(ctx context.Context, merchantID uuid.UUID, commission *entity.MerchantCommission) error

	// Update default commission settings
	UpdateDefaultCommission(ctx context.Context, setstings *entity.CommissionSettings) error
}
