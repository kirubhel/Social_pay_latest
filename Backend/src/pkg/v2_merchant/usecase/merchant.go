package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/repository"

	"github.com/google/uuid"
)

// MerchantUseCase defines the merchant management use case interface
type MerchantUseCase interface {
	GetMerchant(ctx context.Context, id uuid.UUID) (*entity.MerchantResponse, error)
	GetMerchantDetails(ctx context.Context, id uuid.UUID) (*entity.MerchantDetails, error)
	GetMerchantByUserID(ctx context.Context, userID uuid.UUID) (*entity.MerchantResponse, error)
}

type merchantUseCase struct {
	log  logging.Logger
	repo repository.Repository
}

// NewMerchantUseCase creates a new merchant management use case
func NewMerchantUseCase(repo repository.Repository) MerchantUseCase {
	return &merchantUseCase{
		log:  logging.NewStdLogger("[V2_MERCHANT]"),
		repo: repo,
	}
}

// GetMerchant gets a merchant by ID
func (u *merchantUseCase) GetMerchant(ctx context.Context, id uuid.UUID) (*entity.MerchantResponse, error) {
	merchant, err := u.repo.GetMerchant(ctx, id)
	if err != nil {
		u.log.Error("Failed to get merchant", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to get merchant: %w", err)
	}

	if merchant == nil {
		return nil, errors.New("merchant not found")
	}

	response := &entity.MerchantResponse{
		ID:                         merchant.ID,
		UserID:                     merchant.UserID,
		LegalName:                  merchant.LegalName,
		TradingName:                merchant.TradingName,
		BusinessRegistrationNumber: merchant.BusinessRegistrationNumber,
		TaxIdentificationNumber:    merchant.TaxIdentificationNumber,
		BusinessType:               merchant.BusinessType,
		IndustryCategory:           merchant.IndustryCategory,
		IsBettingCompany:           merchant.IsBettingCompany,
		LotteryCertificateNumber:   merchant.LotteryCertificateNumber,
		WebsiteURL:                 merchant.WebsiteURL,
		EstablishedDate:            merchant.EstablishedDate,
		CreatedAt:                  merchant.CreatedAt,
		UpdatedAt:                  merchant.UpdatedAt,
		Status:                     merchant.Status,
	}

	return response, nil
}

// GetMerchantDetails gets complete merchant information with related data
func (u *merchantUseCase) GetMerchantDetails(ctx context.Context, id uuid.UUID) (*entity.MerchantDetails, error) {
	details, err := u.repo.GetMerchantDetails(ctx, id)
	if err != nil {
		u.log.Error("Failed to get merchant details", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to get merchant details: %w", err)
	}

	if details == nil {
		return nil, errors.New("merchant not found")
	}

	return details, nil
}

// GetMerchantByUserID gets a merchant by user ID
func (u *merchantUseCase) GetMerchantByUserID(ctx context.Context, userID uuid.UUID) (*entity.MerchantResponse, error) {
	merchant, err := u.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		u.log.Error("Failed to get merchant by user ID", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		return nil, fmt.Errorf("failed to get merchant by user ID: %w", err)
	}

	if merchant == nil {
		return nil, errors.New("merchant not found")
	}

	response := &entity.MerchantResponse{
		ID:                         merchant.ID,
		UserID:                     merchant.UserID,
		LegalName:                  merchant.LegalName,
		TradingName:                merchant.TradingName,
		BusinessRegistrationNumber: merchant.BusinessRegistrationNumber,
		TaxIdentificationNumber:    merchant.TaxIdentificationNumber,
		BusinessType:               merchant.BusinessType,
		IndustryCategory:           merchant.IndustryCategory,
		IsBettingCompany:           merchant.IsBettingCompany,
		LotteryCertificateNumber:   merchant.LotteryCertificateNumber,
		WebsiteURL:                 merchant.WebsiteURL,
		EstablishedDate:            merchant.EstablishedDate,
		CreatedAt:                  merchant.CreatedAt,
		UpdatedAt:                  merchant.UpdatedAt,
		Status:                     merchant.Status,
	}

	return response, nil
}
