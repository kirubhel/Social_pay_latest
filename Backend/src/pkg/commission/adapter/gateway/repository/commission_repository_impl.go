package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	commission "github.com/socialpay/socialpay/src/pkg/commission/adapter/gateway/repository/generated"
	"github.com/socialpay/socialpay/src/pkg/commission/core/entity"
)

type commissionRepositoryImpl struct {
	db *commission.Queries
}

func NewCommissionRepository(db *sql.DB) CommissionRepository {
	return &commissionRepositoryImpl{
		db: commission.New(db),
	}
}

func (r *commissionRepositoryImpl) GetDefaultCommission(ctx context.Context) (*entity.CommissionSettings, error) {
	valueStr, err := r.db.GetDefaultCommission(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get default commission: %w", err)
	}

	var settings entity.CommissionSettings
	if err := json.Unmarshal([]byte(valueStr), &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commission settings: %w", err)
	}

	return &settings, nil
}

func (r *commissionRepositoryImpl) GetMerchantCommission(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantCommission, error) {
	commission, err := r.db.GetMerchantCommission(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant commission: %w", err)
	}

	result := &entity.MerchantCommission{
		MerchantID:       merchantID,
		CommissionActive: commission.CommissionActive,
	}

	if commission.CommissionPercent.Valid {
		percent, _ := strconv.ParseFloat(commission.CommissionPercent.String, 64)
		result.CommissionPercent = &percent
	}
	if commission.CommissionCent.Valid {
		cent, _ := strconv.ParseFloat(commission.CommissionCent.String, 64)
		result.CommissionCent = &cent
	}

	return result, nil
}

func (r *commissionRepositoryImpl) UpdateMerchantCommission(ctx context.Context, merchantID uuid.UUID, merchantCommission *entity.MerchantCommission) error {
	var percent, cent sql.NullString
	if merchantCommission.CommissionPercent != nil {
		percent.String = strconv.FormatFloat(*merchantCommission.CommissionPercent, 'f', 2, 64)
		percent.Valid = true
	}
	if merchantCommission.CommissionCent != nil {
		cent.String = strconv.FormatFloat(*merchantCommission.CommissionCent, 'f', 4, 64)
		cent.Valid = true
	}

	err := r.db.UpdateMerchantCommission(ctx, commission.UpdateMerchantCommissionParams{
		CommissionActive:  merchantCommission.CommissionActive,
		CommissionPercent: percent,
		CommissionCent:    cent,
		ID:                merchantID,
	})
	if err != nil {
		return fmt.Errorf("failed to update merchant commission: %w", err)
	}

	return nil
}

func (r *commissionRepositoryImpl) UpdateDefaultCommission(ctx context.Context, settings *entity.CommissionSettings) error {
	value, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal commission settings: %w", err)
	}

	err = r.db.UpdateDefaultCommission(ctx, value)
	if err != nil {
		return fmt.Errorf("failed to update default commission: %w", err)
	}

	return nil
}
