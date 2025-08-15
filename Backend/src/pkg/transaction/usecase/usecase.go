package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

type TransactionUseCase interface {
	// TODO
	GetTransactions(c context.Context, UserId uuid.UUID,
		pagination pagination.Pagination) ([]entity.Transaction, int, error)

	GetTransactionByParamenters(c context.Context, UserId uuid.UUID,
		parameter *entity.FilterParameters, pagination pagination.Pagination, queryForAllUsers bool) ([]entity.Transaction, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetByReferenceID(ctx context.Context, referenceID string) (*entity.Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TransactionStatus) error
	ValidateReferenceId(ctx context.Context, merchantID uuid.UUID, referenceID string) error

	// Analytics methods
	GetTransactionAnalytics(ctx context.Context, filter *entity.AnalyticsFilter) (*entity.TransactionAnalytics, error)
	GetChartData(ctx context.Context, filter *entity.ChartFilter) (*entity.ChartData, error)

	// Admin analytics methods
	GetAdminTransactionAnalytics(ctx context.Context, filter *entity.AnalyticsFilter) (*entity.AdminTransactionAnalytics, error)
	GetAdminChartData(ctx context.Context, filter *entity.ChartFilter) (*entity.ChartData, error)
	GetMerchantGrowthAnalytics(ctx context.Context, startDate, endDate time.Time, dateUnit entity.DateUnit) (*entity.MerchantGrowthAnalytics, error)
}
