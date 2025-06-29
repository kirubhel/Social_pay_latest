package usecase

import (
	"context"

	"github.com/google/uuid"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	"github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
	webhookDto "github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
	"github.com/socialpay/socialpay/src/pkg/webhook/adapter/gateway/kafka/producer"
	"github.com/socialpay/socialpay/src/pkg/webhook/core/entity"
)

type WebhookUseCase interface {
	ProcessTransactionStatus(ctx context.Context, txnID uuid.UUID, status txEntity.TransactionStatus) error
	CreateCallbackLog(ctx context.Context, txnID uuid.UUID, responseStatus int, requestBody string, responseBody string, merchantID string, userID string) error
	UpdateCallbackLog(ctx context.Context, id uuid.UUID, responseBody string, responseStatus int) error
	HandlePaymentStatusUpdate(ctx context.Context, msg webhookDto.WebhookMessage) error
	HandleWebhookDispatch(ctx context.Context, req dto.WebhookRequest) error
	GetProducer() *producer.GroupedProducer
	GetSendProducer() *producer.GroupedProducer
	GetCallbackLogByID(ctx context.Context, id uuid.UUID) (*entity.CallbackLog, error)
	GetCallbackLogsByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*entity.CallbackLog, error)
	GetAllCallbackLogs(ctx context.Context, pagination *txEntity.Pagination) ([]*entity.CallbackLog, error)
}
