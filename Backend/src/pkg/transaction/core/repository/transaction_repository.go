package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/filter"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

type TransactionRepository interface {
	GetTransactions(c context.Context, user_id uuid.UUID, limit, offset int32) ([]entity.Transaction, int, error)
	GetTransactionByParamenters(ctx context.Context, parameters *entity.FilterParameters) ([]entity.Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetByIDWithMerchant(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetByReferenceID(ctx context.Context, referenceID string) (*entity.Transaction, error)
	GetByUserIdAndReferenceID(ctx context.Context, userID uuid.UUID, referenceID string) (*entity.Transaction, error)
	GetByMerchantIdAndReferenceID(ctx context.Context, merchantID uuid.UUID, referenceID string) (*entity.Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TransactionStatus) error

	// UpdateTransactionWithProviderData updates transaction with provider information
	UpdateTransactionWithProviderData(ctx context.Context, id uuid.UUID, updateParams map[string]interface{}) error

	// Create creates a new transaction
	Create(ctx context.Context, tx *entity.Transaction) error

	// Update updates an existing transaction
	Update(ctx context.Context, tx *entity.Transaction) error

	// GetTransactionsByParameters retrieves transactions filtered by parameters
	// GetTransactionsByParameters(ctx context.Context, params entity.FilterParameters, userID uuid.UUID, limit, offset int32) ([]entity.Transaction, error)
	GetTransactionsByParameters(ctx context.Context, filterParam filter.Filter, userID uuid.UUID) ([]entity.Transaction, error)
	// GetTransactionsByStatus retrieves transactions by status
	GetTransactionsByStatus(ctx context.Context, status entity.TransactionStatus, limit, offset int32) ([]entity.Transaction, error)

	// GetTransactionsByType retrieves transactions by type
	GetTransactionsByType(ctx context.Context, txType entity.TransactionType, limit, offset int32) ([]entity.Transaction, error)

	// Get Transaction history
	GetMerchantTransactions(ctx context.Context, merchantID uuid.UUID, limit, offset int32) ([]entity.Transaction, error)

	GetFilteredMerchantTransactions(ctx context.Context, params *entity.FilterParameters, merchantID uuid.UUID, limit, offset int32) ([]entity.Transaction, error)
	OverrideTransactionStatus(ctx context.Context, txnID uuid.UUID, newStatus entity.TransactionStatus, reason string, adminID string) error

	//

	GetTransactionByParametersCount(ctx context.Context,
		filterParam filter.Filter, userID uuid.UUID) (int, error)
	// Count
	// CountTransactionWithParameter(ctx context.Context,clause string,
	// 	args []interface{}) (int, error)

	// CountWithClause(ctx context.Context,
	// 	baseTable string, clause string, args ...interface{}) (int, error)

	// QR Payment and Tip Processing methods
	CreateWithContext(ctx context.Context, tx *entity.Transaction) error
	UpdateTipProcessing(ctx context.Context, transactionID, tipTransactionID uuid.UUID) error
	GetTransactionsWithPendingTips(ctx context.Context) ([]entity.Transaction, error)
	GetTransactionsByQRLink(ctx context.Context, qrLinkID uuid.UUID, limit, offset int32) ([]entity.Transaction, error)

	// Analytics methods
	GetTransactionAnalytics(ctx context.Context, filter *entity.AnalyticsFilter, userID uuid.UUID) (*entity.TransactionAnalytics, error)
	GetChartData(ctx context.Context, filter *entity.ChartFilter, userID uuid.UUID) (*entity.ChartData, error)

	// Admin analytics methods
	GetAdminTransactionAnalytics(ctx context.Context, filter *entity.AnalyticsFilter) (*entity.AdminTransactionAnalytics, error)
	GetAdminChartData(ctx context.Context, filter *entity.ChartFilter) (*entity.ChartData, error)
	GetMerchantGrowthAnalytics(ctx context.Context, startDate, endDate time.Time, dateUnit entity.DateUnit) (*entity.MerchantGrowthAnalytics, error)
}
