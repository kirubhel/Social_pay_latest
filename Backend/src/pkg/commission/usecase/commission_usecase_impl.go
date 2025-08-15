package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/commission/core/entity"
	"github.com/socialpay/socialpay/src/pkg/commission/core/repository"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
)

type commissionUseCaseImpl struct {
	repo repository.CommissionRepository
	log  logging.Logger
}

func NewCommissionUseCase(repo repository.CommissionRepository) CommissionUseCase {
	return &commissionUseCaseImpl{
		repo: repo,
		log:  logging.NewStdLogger("[commission]"),
	}
}

func (uc *commissionUseCaseImpl) CalculateCommission(ctx context.Context, amount float64, merchantID uuid.UUID) (*entity.CommissionSettings, error) {
	// First check if merchant has custom commission
	merchantCommission, err := uc.repo.GetMerchantCommission(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant commission: %w", err)
	}

	// If merchant has active custom commission, use it
	if merchantCommission.CommissionActive && merchantCommission.CommissionPercent != nil {
		return &entity.CommissionSettings{
			Percent: *merchantCommission.CommissionPercent,
			Cent:    entity.GetFloat64OrDefault(merchantCommission.CommissionCent, 0),
		}, nil
	}

	// Fallback to default commission
	defaultCommission, err := uc.repo.GetDefaultCommission(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get default commission: %w", err)
	}

	return defaultCommission, nil
}

func (uc *commissionUseCaseImpl) GetMerchantCommission(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantCommission, error) {
	return uc.repo.GetMerchantCommission(ctx, merchantID)
}

func (uc *commissionUseCaseImpl) UpdateMerchantCommission(ctx context.Context, merchantID uuid.UUID, commission *entity.MerchantCommission) error {
	return uc.repo.UpdateMerchantCommission(ctx, merchantID, commission)
}

func (uc *commissionUseCaseImpl) GetDefaultCommission(ctx context.Context) (*entity.CommissionSettings, error) {
	return uc.repo.GetDefaultCommission(ctx)
}

func (uc *commissionUseCaseImpl) UpdateDefaultCommission(ctx context.Context, settings *entity.CommissionSettings) error {
	return uc.repo.UpdateDefaultCommission(ctx, settings)
}
