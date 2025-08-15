package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/commission/core/entity"
	"github.com/socialpay/socialpay/src/pkg/commission/usecase"
)

type CommissionCalculator struct {
	commissionUseCase usecase.CommissionUseCase
}

func NewCommissionCalculator(commissionUseCase usecase.CommissionUseCase) *CommissionCalculator {
	return &CommissionCalculator{
		commissionUseCase: commissionUseCase,
	}
}

type CommissionResult struct {
	// Base commission amount
	FeeAmount float64
	// VAT amount
	VatAmount float64
	// Total commission (FeeAmount + VatAmount)
	TotalCommission float64
	// Merchant's net amount (TotalAmount - TotalCommission)
	MerchantNet float64
	// Admin's net amount (FeeAmount)
	AdminNet float64
}

// CalculateCommission calculates the commission for a transaction based on the payment processor
func (c *CommissionCalculator) CalculateCommission(ctx context.Context, amount float64, merchantID uuid.UUID) CommissionResult {
	// Get commission settings from the commission usecase
	commissionSettings, err := c.commissionUseCase.CalculateCommission(ctx, amount, merchantID)
	if err != nil {
		// If there's an error getting commission settings, use default values
		commissionSettings = &entity.CommissionSettings{
			Percent:   2.75, // Default commission rate
			Cent:      0.00,
		}
	}

	// Calculate base commission
	feeAmount := amount * (commissionSettings.Percent / 100.0)
	feeAmount += commissionSettings.Cent // Add fixed cent amount


	// Calculate VAT (15% of commission)
	vatAmount := feeAmount * 0.15

	// Calculate total commission
	totalCommission := feeAmount + vatAmount

	// Calculate merchant's net amount
	merchantNet := amount - totalCommission

	return CommissionResult{
		FeeAmount:       feeAmount,
		VatAmount:       vatAmount,
		TotalCommission: totalCommission,
		MerchantNet:     merchantNet,
		AdminNet:        feeAmount,
	}
}
