package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	db "github.com/socialpay/socialpay/src/pkg/webhook/adapter/gateway/repository/generated"
	"github.com/socialpay/socialpay/src/pkg/webhook/core/entity"
)

type CallbackRepositoryImpl struct {
	queries *db.Queries
}

func NewCallbackRepository(dbConn *sql.DB) CallbackRepository {
	return &CallbackRepositoryImpl{
		queries: db.New(dbConn),
	}
}

func (r *CallbackRepositoryImpl) Create(ctx context.Context, log *entity.CallbackLog) error {
	// Log the merchant_id before insertion
	fmt.Printf("Attempting to create callback log with merchant_id: %s\n", log.MerchantID)

	params := db.CreateCallbackLogParams{
		ID:           log.ID,
		UserID:       log.UserID,
		TxnID:        log.TxnID,
		MerchantID:   log.MerchantID,
		Status:       int32(log.Status),
		RequestBody:  log.RequestBody,
		ResponseBody: sql.NullString{String: log.ResponseBody, Valid: log.ResponseBody != ""},
		RetryCount:   int32(log.RetryCount),
	}
	return r.queries.CreateCallbackLog(ctx, params)
}

func (r *CallbackRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.CallbackLog, error) {
	dbLog, err := r.queries.GetCallbackLogByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toEntityCallbackLog(&dbLog), nil
}

func (r *CallbackRepositoryImpl) Update(ctx context.Context, log *entity.CallbackLog) error {
	params := db.UpdateCallbackLogParams{
		ID:           log.ID,
		Status:       int32(log.Status),
		ResponseBody: sql.NullString{String: log.ResponseBody, Valid: log.ResponseBody != ""},
		RetryCount:   int32(log.RetryCount),
	}
	return r.queries.UpdateCallbackLog(ctx, params)
}

func (r *CallbackRepositoryImpl) GetByTransactionID(ctx context.Context, txnID uuid.UUID) ([]*entity.CallbackLog, error) {
	dbLogs, err := r.queries.GetCallbackLogsByTransactionID(ctx, txnID)
	if err != nil {
		return nil, err
	}
	return toEntityCallbackLogs(dbLogs), nil
}

func (r *CallbackRepositoryImpl) GetByStatus(ctx context.Context, status int) ([]*entity.CallbackLog, error) {
	dbLogs, err := r.queries.GetCallbackLogsByStatus(ctx, int32(status))
	if err != nil {
		return nil, err
	}
	return toEntityCallbackLogs(dbLogs), nil
}

func (r *CallbackRepositoryImpl) GetByMerchantID(ctx context.Context, merchantID uuid.UUID, pagination *txEntity.Pagination) ([]*entity.CallbackLog, error) {
	// Calculate limit and offset
	limit := int32(pagination.PageSize)
	offset := int32((pagination.Page - 1) * pagination.PageSize)

	dbLogs, err := r.queries.GetCallbackLogsByMerchantID(ctx, db.GetCallbackLogsByMerchantIDParams{
		MerchantID: merchantID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, err
	}
	return toEntityCallbackLogs(dbLogs), nil
}

func (r *CallbackRepositoryImpl) GetAll(ctx context.Context, pagination *txEntity.Pagination) ([]*entity.CallbackLog, error) {
	// Calculate limit and offset
	limit := int32(pagination.PageSize)
	offset := int32((pagination.Page - 1) * pagination.PageSize)

	dbLogs, err := r.queries.GetAllCallbackLogs(ctx, db.GetAllCallbackLogsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toEntityCallbackLogs(dbLogs), nil
}

// Helper functions to convert between database and entity types
func toEntityCallbackLog(dbLog *db.WebhookCallbackLog) *entity.CallbackLog {
	return &entity.CallbackLog{
		ID:           dbLog.ID,
		UserID:       dbLog.UserID,
		TxnID:        dbLog.TxnID,
		MerchantID:   dbLog.MerchantID,
		Status:       int(dbLog.Status),
		RequestBody:  dbLog.RequestBody,
		ResponseBody: dbLog.ResponseBody.String,
		RetryCount:   int(dbLog.RetryCount),
		CreatedAt:    dbLog.CreatedAt,
		UpdatedAt:    dbLog.UpdatedAt,
	}
}

func toEntityCallbackLogs(dbLogs []db.WebhookCallbackLog) []*entity.CallbackLog {
	logs := make([]*entity.CallbackLog, len(dbLogs))
	for i, dbLog := range dbLogs {
		logs[i] = toEntityCallbackLog(&dbLog)
	}
	return logs
}
