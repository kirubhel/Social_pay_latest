package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/commission/core/entity"
)

type CommissionUseCase interface {
	// Calculate commission for a transaction
	CalculateCommission(ctx context.Context, amount float64, merchantID uuid.UUID) (*entity.CommissionSettings, error)

	// Get merchant commission settings
	GetMerchantCommission(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantCommission, error)

	// Update merchant commission settings
	UpdateMerchantCommission(ctx context.Context, merchantID uuid.UUID, commission *entity.MerchantCommission) error

	// Get default commission settings
	GetDefaultCommission(ctx context.Context) (*entity.CommissionSettings, error)

	// Update default commission settings
	UpdateDefaultCommission(ctx context.Context, settings *entity.CommissionSettings) error
}
