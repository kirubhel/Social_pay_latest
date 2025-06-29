package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/errorxx"
	"github.com/socialpay/socialpay/src/pkg/shared/filter"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	transaction_repository "github.com/socialpay/socialpay/src/pkg/transaction/core/repository"
)

type transactionUseCase struct {
	repo transaction_repository.TransactionRepository
	log  logging.Logger
}

// OverrideTransactionStatus implements TransactionUseCase.
func (t *transactionUseCase) OverrideTransactionStatus(ctx context.Context, txnID uuid.UUID, newStatus entity.TransactionStatus, reason string, adminID string) error {
	return t.repo.OverrideTransactionStatus(ctx, txnID, newStatus, reason, adminID)
}

func (t *transactionUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	return t.repo.GetByID(ctx, id)
}

func (t *transactionUseCase) GetByReferenceID(ctx context.Context, referenceID string) (*entity.Transaction, error) {
	return t.repo.GetByReferenceID(ctx, referenceID)
}

func (t *transactionUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TransactionStatus) error {
	return t.repo.UpdateStatus(ctx, id, status)
}

func NewTransactionUsecase(
	repo transaction_repository.TransactionRepository) TransactionUseCase {
	return &transactionUseCase{
		repo: repo,
		log:  logging.NewStdLogger("[TRANSACTION] [USECASE]"),
	}
}

func (t *transactionUseCase) GetTransactions(c context.Context, UserId uuid.UUID,
	pagination pagination.Pagination) ([]entity.Transaction, int, error) {

	// calculating limit and offset
	limit := pagination.GetLimit()
	offset := pagination.GetOffset()

	// get transaction by user_id by getting from context
	transactions, count, err := t.repo.GetTransactions(c, UserId,
		int32(limit), int32(offset))

	if err != nil {

		err = errorxx.ErrDBRead.Wrap(err, "Db Read transaction").
			WithProperty(errorxx.ErrorCode, 500)

		t.log.Error("get transactions error",
			map[string]interface{}{
				"error":   err,
				"context": c,
			})
		return nil, 0, err
	}
	return transactions, count, nil
}

func (t *transactionUseCase) GetTransactionByParamenters(c context.Context, UserId uuid.UUID,
	parameter *entity.FilterParameters, pagination pagination.Pagination) ([]entity.Transaction, int, error) {

	// validating parameters
	if err := parameter.Validate(); err != nil {

		err = errorxx.ErrAppBadInput.Wrap(err, "get transaction err").
			WithProperty(errorxx.ErrorCode, 500)

		t.log.Info("parameter validation error",
			map[string]interface{}{
				"error":   err,
				"context": c,
			})
		return nil, 0, err
	}
	// building filter obj
	filterParameter := parameter.ToFilter()

	filterParameter.Group.Fields = append(filterParameter.Group.Fields, filter.Field{
		Name:     "user_id",
		Operator: "=",
		Value:    UserId,
	})

	count, err := t.repo.GetTransactionByParametersCount(c, filterParameter, UserId)

	if err != nil {

		err = errorxx.ErrDBRead.Wrap(err, "count err").
			WithProperty(errorxx.ErrorCode, 500)

		t.log.Error("GET_TRANSACTION::ERR::",
			map[string]interface{}{
				"error":   err,
				"user_id": UserId,
			})

		return nil, 0, err
	}

	filterParameter.Pagination = pagination
	trans, err := t.repo.GetTransactionsByParameters(c, filterParameter, UserId)

	if err != nil {
		err = errorxx.ErrDBRead.Wrap(err, "get transaction with filter parameters err").
			WithProperty(errorxx.ErrorCode, 500)

		t.log.Error("GET_TRANSACTION::ERR::",
			map[string]interface{}{
				"error":   err,
				"user_id": UserId,
			})

		return nil, count, err
	}

	return trans, count, nil
}

func (t *transactionUseCase) ValidateReferenceId(c context.Context, merchantID uuid.UUID, referenceID string) error {
	// check if transaction already exists
	t.log.Info("validating reference id, merchantID", map[string]interface{}{
		"referenceID": referenceID,
		"merchantID":  merchantID,
	})
	checktx, err := t.repo.GetByMerchantIdAndReferenceID(c, merchantID, referenceID)
	t.log.Info("checktx", map[string]interface{}{
		"checktx": checktx,
	})
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			t.log.Error("failed to get transaction", map[string]interface{}{
				"error": err.Error(),
			})
			return fmt.Errorf("failed to get transaction: %w", err)
		}
		// No existing transaction found, which is what we want
		t.log.Info("No existing transaction found with this reference", map[string]interface{}{
			"reference": referenceID,
		})
	} else if checktx != nil {
		t.log.Error("duplicate transaction reference", map[string]interface{}{
			"reference": referenceID,
		})
		return fmt.Errorf("please use unique reference for each transaction")
	}
	return nil
}

// GetTransactionAnalytics retrieves aggregated transaction analytics
func (t *transactionUseCase) GetTransactionAnalytics(ctx context.Context, filter *entity.AnalyticsFilter) (*entity.TransactionAnalytics, error) {
	// Validate the filter parameters
	if err := filter.Validate(); err != nil {
		err = errorxx.ErrAppBadInput.Wrap(err, "analytics filter validation error").
			WithProperty(errorxx.ErrorCode, 400)

		t.log.Error("analytics filter validation error", map[string]interface{}{
			"error":   err,
			"context": ctx,
		})
		return nil, err
	}

	// Get user ID from context (assuming it's available)
	merchantID, ok := ctx.Value("merchant_id").(uuid.UUID)
	if !ok {
		err := errorxx.ErrAuthUnauthorized.Wrap(nil, "user ID not found in context").
			WithProperty(errorxx.ErrorCode, 401)

		t.log.Error("user ID not found in context", map[string]interface{}{
			"context": ctx,
		})
		return nil, err
	}

	// Get analytics from repository
	analytics, err := t.repo.GetTransactionAnalytics(ctx, filter, merchantID)
	if err != nil {
		err = errorxx.ErrDBRead.Wrap(err, "failed to get transaction analytics").
			WithProperty(errorxx.ErrorCode, 500)

		t.log.Error("failed to get transaction analytics", map[string]interface{}{
			"error":      err,
			"merchantID": merchantID,
			"filter":     filter,
		})
		return nil, err
	}

	t.log.Info("transaction analytics retrieved successfully", map[string]interface{}{
		"merchantID":         merchantID,
		"total_transactions": analytics.TotalTransactions,
		"total_amount":       analytics.TotalAmount,
		"date_range":         fmt.Sprintf("%s to %s", filter.StartDate.Format("2006-01-02"), filter.EndDate.Format("2006-01-02")),
	})

	return analytics, nil
}

// GetChartData retrieves chart data for analytics
func (t *transactionUseCase) GetChartData(ctx context.Context, filter *entity.ChartFilter) (*entity.ChartData, error) {
	// Validate the filter parameters
	if err := filter.Validate(); err != nil {
		err = errorxx.ErrAppBadInput.Wrap(err, "chart filter validation error").
			WithProperty(errorxx.ErrorCode, 400)

		t.log.Error("chart filter validation error", map[string]interface{}{
			"error":   err,
			"context": ctx,
		})
		return nil, err
	}

	// Get user ID from context
	merchantID, ok := ctx.Value("merchant_id").(uuid.UUID)
	if !ok {
		err := errorxx.ErrAuthUnauthorized.Wrap(nil, "user ID not found in context").
			WithProperty(errorxx.ErrorCode, 401)

		t.log.Error("user ID not found in context", map[string]interface{}{
			"context": ctx,
		})
		return nil, err
	}

	// Get chart data from repository
	chartData, err := t.repo.GetChartData(ctx, filter, merchantID)
	if err != nil {
		err = errorxx.ErrDBRead.Wrap(err, "failed to get chart data").
			WithProperty(errorxx.ErrorCode, 500)

		t.log.Error("failed to get chart data", map[string]interface{}{
			"error":      err,
			"merchantID": merchantID,
			"filter":     filter,
		})
		return nil, err
	}

	t.log.Info("chart data retrieved successfully", map[string]interface{}{
		"merchantID":  merchantID,
		"chart_type":  filter.ChartType,
		"date_unit":   filter.DateUnit,
		"data_points": len(chartData.Data),
	})

	return chartData, nil
}
