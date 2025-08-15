package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
)

// Repository defines the interface for merchant storage operations
type Repository interface {
	// GetMerchant retrieves a merchant by its ID
	GetMerchant(ctx context.Context, id uuid.UUID) (*entity.Merchant, error)

	// GetMerchants retrieves list of merchants
	GetMerchants(ctx context.Context, params entity.GetMerchantsParams) (*entity.MerchantsResponse, error)

	// GetAllMerchants retries all the merchants
	GetAllMerchants(ctx context.Context) ([]entity.Merchant, error)

	// GetMerchantDetails retrieves complete merchant information with related data
	GetMerchantDetails(ctx context.Context, id uuid.UUID) (*entity.MerchantDetails, error)

	// GetMerchantByUserID retrieves a merchant by user ID
	GetMerchantByUserID(ctx context.Context, userID uuid.UUID) (*entity.Merchant, error)

	// GetMerchantAddresses retrieves all addresses for a merchant
	GetMerchantAddresses(ctx context.Context, merchantID uuid.UUID) ([]entity.MerchantAddress, error)

	// GetMerchantContacts retrieves all contacts for a merchant
	GetMerchantContacts(ctx context.Context, merchantID uuid.UUID) ([]entity.MerchantContact, error)

	// GetMerchantDocuments retrieves all documents for a merchant
	GetMerchantDocuments(ctx context.Context, merchantID uuid.UUID) ([]entity.MerchantDocument, error)

	// GetMerchantBankAccounts retrieves all bank accounts for a merchant
	GetMerchantBankAccounts(ctx context.Context, merchantID uuid.UUID) ([]entity.MerchantBankAccount, error)

	// GetMerchantSettings retrieves settings for a merchant
	GetMerchantSettings(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantSettings, error)

	// UpdateMerchant updates merchant info
	UpdateMerchant(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantRequest) error

	// UpdateMerchantStatus updates merchant status
	UpdateMerchantStatus(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantStatusRequest) error

	// UpdateMerchantContact updates merchant contact info
	UpdateMerchantContact(ctx context.Context, id uuid.UUID, req *entity.UpdateMerchantContactRequest) error

	// UpdateMerchantDocument updates merchant document
	UpdateMerchantDocument(ctx context.Context, id uuid.UUID, req *entity.UpdateMerchantDocumentRequest) error

	// UpdateAdminMerchant updates merchant by admin
	UpdateAdminMerchant(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantRequest) error

	// DeleteMerchant soft deletes merchant
	DeleteMerchant(ctx context.Context, merchantID uuid.UUID) error

	// DeleteMerchants soft deletes list of merchants
	DeleteMerchants(ctx context.Context, ids []uuid.UUID) error

	// GetMerchantStats retrieves merchant statistics
	GetMerchantStats(ctx context.Context) (*entity.MerchantStats, error)
}
