// github.com/socialpay/socialpay/src/pkg/key/core/usecase/repository.go
package usecase

import (
	"context"

	"github.com/socialpay/socialpay/src/pkg/merchants/core/entity"

	"github.com/google/uuid"
)

type Repository interface {
	Save(apiKey *entity.Merchant) error
	FindByUserID(userID uuid.UUID) (*entity.Merchant, error)
	GetMerchants() ([]entity.Merchant, error)
	GetMerchantsByUserID(userID uuid.UUID) ([]entity.Merchant, error)
	SaveDocument(ctx context.Context, doc entity.MerchantDocument) error
	UpdateMerchantStatus(merchantID uuid.UUID, status entity.MerchantStatus) error
	UpdateFullMerchant(ctx context.Context, merchantID uuid.UUID, merchant *entity.Merchant,
		address *entity.MerchantAdditionalInfo, documents []entity.MerchantDocument) error
	DeleteMerchant(ctx context.Context, merchantID uuid.UUID) error
	GetMerchantByID(merchantID uuid.UUID) (*entity.Merchant, error)
	GetMerchantDetails(uuid.UUID) (*entity.MerchantDetails, error)
	CreateMerchantAdditionalInfo(ctx context.Context, info entity.MerchantAdditionalInfo) error
}
