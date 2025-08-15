package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	auth_service "github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	auth_utils "github.com/socialpay/socialpay/src/pkg/authv2/utils"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/repository"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/utils"

	"github.com/google/uuid"
)

// MerchantUseCase defines the merchant management use case interface
type MerchantUseCase interface {
	AddMerchant(ctx context.Context, req *auth_entity.CreateUserRequest) error
	GetMerchant(ctx context.Context, id uuid.UUID) (*entity.MerchantResponse, error)
	GetMerchants(ctx context.Context, params entity.GetMerchantsParams) (*entity.MerchantsResponse, error)
	GetMerchantDetails(ctx context.Context, id uuid.UUID) (*entity.MerchantDetails, error)
	GetMerchantByUserID(ctx context.Context, userID uuid.UUID) (*entity.MerchantResponse, error)
	ExportMerchants(ctx context.Context, req *entity.ExportMerchantsRequest) (*string, error)
	UpdateMerchant(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantRequest) error
	UpdateMerchantStatus(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantStatusRequest) error
	UpdateMerchantContact(ctx context.Context, id uuid.UUID, req *entity.UpdateMerchantContactRequest) error
	UpdateMerchantDocument(ctx context.Context, id uuid.UUID, req *entity.UpdateMerchantDocumentRequest) error
	UpdateAdminMerchant(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantRequest) error
	DeleteMerchant(ctx context.Context, merchantID uuid.UUID) error
	DeleteMerchants(ctx context.Context, req *entity.DeleteMerchantsRequest) error
	GetMerchantStats(ctx context.Context) (*entity.MerchantStats, error)
	ImpersonateMerchant(ctx context.Context, merchantID uuid.UUID) (*auth_entity.AuthResponse, error)
}

type merchantUseCase struct {
	authService auth_service.AuthService
	log         logging.Logger
	repo        repository.Repository
}

// NewMerchantUseCase creates a new merchant management use case
func NewMerchantUseCase(authService auth_service.AuthService, repo repository.Repository) MerchantUseCase {
	return &merchantUseCase{
		authService: authService,
		log:         logging.NewStdLogger("[V2_MERCHANT]"),
		repo:        repo,
	}
}

// AddMerchat adds new merchant
func (u *merchantUseCase) AddMerchant(ctx context.Context, req *auth_entity.CreateUserRequest) error {
	_, err := u.authService.Register(ctx, req)

	if err != nil {
		u.log.Error("Failed to add merchant", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to get merchant: %w", err)
	}

	return nil
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

// GetMerchants gets list of merchants
func (u *merchantUseCase) GetMerchants(ctx context.Context, params entity.GetMerchantsParams) (*entity.MerchantsResponse, error) {
	merchants, err := u.repo.GetMerchants(ctx, params)
	if err != nil {
		u.log.Error("Failed to get merchants", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get merchants: %w", err)
	}

	return merchants, nil
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

// ExportMerchants exports merchants data in the form of specific file type
func (u *merchantUseCase) ExportMerchants(ctx context.Context, req *entity.ExportMerchantsRequest) (*string, error) {
	var merchantsDetails []entity.MerchantDetails
	var merchants []uuid.UUID

	if len(req.Merchants) == 0 {
		allMerchants, err := u.repo.GetAllMerchants(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get all merchants details: %w", err)
		}

		for _, merchant := range allMerchants {
			merchants = append(merchants, merchant.ID)
		}
	} else {
		merchants = req.Merchants
	}

	fmt.Println(merchants)
	for _, merchant := range merchants {
		merchantDetails, err := u.GetMerchantDetails(ctx, merchant)
		if err != nil {
			return nil, fmt.Errorf("failed to get merchant details: %w", err)
		}

		//filteredMerchantDetails := utils.FilterMerchantDetails(*merchantDetails, req.Data)
		filteredMerchantDetails, err := utils.FilterMerchantDetailsJSON(*merchantDetails, req.Data)

		if err != nil {
			return nil, fmt.Errorf("failed to filter merchat details: %v", filteredMerchantDetails)
		}

		fmt.Println(filteredMerchantDetails)
		merchantsDetails = append(merchantsDetails, filteredMerchantDetails)
	}

	var filePath string

	switch req.FileType {
	case "csv":
		fileName := "export-" + time.Now().String() + ".zip"
		filePath = "public/" + fileName
		err := utils.WriteMerchantDetailsToCSVZip(merchantsDetails, fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to write data to zip of cv files: %w", err)
		}
	case "xlsx":
		fileName := "export-" + time.Now().String() + ".xlsx"
		filePath = "public/" + fileName
		err := utils.WriteMerchantDetailsToXLSX(merchantsDetails, fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to write data to xlsx file: %w", err)
		}
	case "json":
		fileName := "export-" + time.Now().String() + ".json"
		filePath = "public/" + fileName
		err := utils.WriteMerchantDetailsToJSON(merchantsDetails, fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to write data to json file: %w", err)
		}
	default:
		err := errors.New("INVALID FILE TYPE")
		return nil, fmt.Errorf("invalid file type: %w", err)
	}

	return &filePath, nil
}

// UpdateMerchant updates merchant
func (u *merchantUseCase) UpdateMerchant(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantRequest) error {
	err := u.repo.UpdateMerchant(ctx, merchantID, req)

	if err != nil {
		return err
	}

	return nil
}

// UpdateMerchantStatus updates merchant status
func (u *merchantUseCase) UpdateMerchantStatus(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantStatusRequest) error {
	err := u.repo.UpdateMerchantStatus(ctx, merchantID, req)

	if err != nil {
		return err
	}

	return nil
}

// UpdateMerchantContact updates merchant contact
func (u *merchantUseCase) UpdateMerchantContact(ctx context.Context, id uuid.UUID, req *entity.UpdateMerchantContactRequest) error {
	err := u.repo.UpdateMerchantContact(ctx, id, req)

	if err != nil {
		return err
	}

	return nil
}

// UpdateMerchantDocument updates merchant document
func (u *merchantUseCase) UpdateMerchantDocument(ctx context.Context, id uuid.UUID, req *entity.UpdateMerchantDocumentRequest) error {
	err := u.repo.UpdateMerchantDocument(ctx, id, req)

	if err != nil {
		return err
	}

	return nil
}

// UpdateAdminMerchant updates merchant by admin
func (u *merchantUseCase) UpdateAdminMerchant(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantRequest) error {
	err := u.repo.UpdateAdminMerchant(ctx, merchantID, req)

	if err != nil {
		return err
	}

	return nil
}

// DeleteMerchant deletes merchant
func (u *merchantUseCase) DeleteMerchant(ctx context.Context, merchantID uuid.UUID) error {
	err := u.repo.DeleteMerchant(ctx, merchantID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMerchants deletes list of merchants
func (u *merchantUseCase) DeleteMerchants(ctx context.Context, req *entity.DeleteMerchantsRequest) error {
	err := u.repo.DeleteMerchants(ctx, req.IDs)

	if err != nil {
		return err
	}

	return nil
}

// GetMerchantStats gets merchant statistics
func (u *merchantUseCase) GetMerchantStats(ctx context.Context) (*entity.MerchantStats, error) {
	stats, err := u.repo.GetMerchantStats(ctx)
	if err != nil {
		u.log.Error("Failed to get merchant stats", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get merchant stats: %w", err)
	}

	return stats, nil
}

// ImpersonateMerchant impersonate merchant
func (u *merchantUseCase) ImpersonateMerchant(ctx context.Context, merchantID uuid.UUID) (*auth_entity.AuthResponse, error) {
	merchant, err := u.repo.GetMerchant(ctx, merchantID)
	if err != nil {
		u.log.Error("Failed to get merchant", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get merchant: %w", err)
	}

	user, err := u.authService.GetUserProfile(ctx, merchant.UserID)
	if err != nil {
		return nil, auth_entity.NewAuthError(auth_entity.ErrInternalServer, auth_entity.MsgInternalServer)
	}

	device, err := u.authService.CreateDevice(ctx, auth_entity.CreateDeviceArgs{
		IP:    "unknown",
		Name:  "login",
		Agent: "login",
	})

	if err != nil {
		u.log.Error("Failed to create session", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, auth_entity.NewAuthError(auth_entity.ErrInternalServer, auth_entity.MsgInternalServer)
	}

	// Generate tokens
	sessionID := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := auth_utils.GenerateJWT(user.ID.String(), string(user.UserType), merchantID.String(), sessionID.String(), jwtSecret, 24)
	if err != nil {
		u.log.Error("Failed to generate JWT", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, auth_entity.NewAuthError(auth_entity.ErrInternalServer, auth_entity.MsgInternalServer)
	}

	refreshToken, err := auth_utils.GenerateRefreshToken()
	if err != nil {
		u.log.Error("Failed to generate refresh token", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, auth_entity.NewAuthError(auth_entity.ErrInternalServer, auth_entity.MsgInternalServer)
	}

	// Create session
	_, err = u.authService.CreateSession(ctx, auth_entity.CreateSessionArgs{
		UserID:       user.ID,
		DeviceID:     device.ID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		u.log.Error("Failed to create session", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, auth_entity.NewAuthError(auth_entity.ErrInternalServer, auth_entity.MsgInternalServer)
	}

	var merchants []auth_entity.Merchant
	merchants = append(merchants, auth_entity.Merchant{
		ID:                         merchant.ID,
		LegalName:                  merchant.LegalName,
		TradingName:                *merchant.TradingName,
		BusinessRegistrationNumber: merchant.BusinessRegistrationNumber,
		TaxIdentificationNumber:    merchant.TaxIdentificationNumber,
		IndustryCategory:           *merchant.IndustryCategory,
		UserID:                     merchant.UserID,
		BusinessType:               merchant.BusinessType,
		IsBettingCompany:           merchant.IsBettingCompany,
		LotteryCertificateNumber:   *merchant.LotteryCertificateNumber,
		WebsiteURL:                 *merchant.WebsiteURL,
		EstablishedDate:            *merchant.EstablishedDate,
		Status:                     string(merchant.Status),
		CreatedAt:                  merchant.CreatedAt,
		UpdatedAt:                  merchant.UpdatedAt,
	})

	return &auth_entity.AuthResponse{
		User:         user,
		Merchants:    merchants,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
