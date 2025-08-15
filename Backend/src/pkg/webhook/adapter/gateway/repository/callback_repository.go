package repository

import (
	"context"

	"github.com/google/uuid"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	webhookEntity "github.com/socialpay/socialpay/src/pkg/webhook/core/entity"
)

type CallbackRepository interface {
	Create(ctx context.Context, log *webhookEntity.CallbackLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*webhookEntity.CallbackLog, error)
	Update(ctx context.Context, log *webhookEntity.CallbackLog) error
	GetByTransactionID(ctx context.Context, txnID uuid.UUID) ([]*webhookEntity.CallbackLog, error)
	GetByStatus(ctx context.Context, status int) ([]*webhookEntity.CallbackLog, error)
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID, pagination *txEntity.Pagination) ([]*webhookEntity.CallbackLog, error)
	GetAll(ctx context.Context, pagination *txEntity.Pagination) ([]*webhookEntity.CallbackLog, error)
}
