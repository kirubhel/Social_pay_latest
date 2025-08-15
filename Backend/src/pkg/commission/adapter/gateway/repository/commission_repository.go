package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/commission/core/entity"
)

type CommissionRepository interface {
	GetDefaultCommission(ctx context.Context) (*entity.CommissionSettings, error)
	GetMerchantCommission(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantCommission, error)
	UpdateMerchantCommission(ctx context.Context, merchantID uuid.UUID, commission *entity.MerchantCommission) error
	UpdateDefaultCommission(ctx context.Context, settings *entity.CommissionSettings) error
}
