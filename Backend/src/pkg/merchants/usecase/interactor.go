package usecase

import (
	"context"
	"database/sql"
	"mime/multipart"
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/merchants/core/entity"

	"github.com/google/uuid"
)

type Interactor interface {
	CreateMerchant(entity.Merchant) (*entity.Merchant, error)
	GetMerchantByUserID(uuid.UUID) (*entity.Merchant, error)
	GetMerchantDetails(uuid.UUID) (*entity.MerchantDetails, error)

	GetMerchants() ([]entity.Merchant, error)
	UpdateMerchantStatus(merchantID uuid.UUID, status entity.MerchantStatus) error
	UpdateFullMerchant(ctx context.Context, merchantID uuid.UUID, merchant *entity.Merchant,
		address *entity.MerchantAdditionalInfo, documents []entity.MerchantDocument) error
	DeleteMerchant(ctx context.Context, merchantID uuid.UUID) error
	AddDocument(ctx context.Context, file multipart.File,
		fileHeader multipart.FileHeader, doc entity.MerchantDocument) error

	AddMerchantInfo(ctx context.Context, userId uuid.UUID,
		info entity.MerchantAdditionalInfo) error
}

type MerchantRepository interface {
	GetMerchantID(cookie http.Cookie) (string, error)
}

type MySQLKeyRepository struct {
	DB *sql.DB
}

func NewMySQLKeyRepository(db *sql.DB) *MySQLKeyRepository {
	return &MySQLKeyRepository{DB: db}
}

// Implement all KeyRepository methods for MySQLKeyRepository...
