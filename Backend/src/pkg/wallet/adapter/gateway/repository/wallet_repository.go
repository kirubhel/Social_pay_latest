package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/wallet/core/entity"
)

type WalletRepository interface {
	// Common operations
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	RollbackTx(tx *sql.Tx) error

	// Merchant wallet operations
	CreateMerchantWallet(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID, amount float64, lockedAmount float64, currency string) error
	GetMerchantWalletByUserID(ctx context.Context, userID uuid.UUID) (*entity.MerchantWallet, error)
	GetMerchantWalletByMerchantID(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantWallet, error)
	GetMerchantWalletByMerchantIDForUpdate(ctx context.Context, tx *sql.Tx, merchantID uuid.UUID) (*entity.MerchantWallet, error)
	UpdateMerchantWallet(ctx context.Context, walletID uuid.UUID, amount float64, lockedAmount float64) error
	UpdateMerchantWalletAmountByMerchantID(ctx context.Context, merchantID uuid.UUID, amount float64) error

	// Admin wallet operations
	GetAdminWallet(ctx context.Context) (*entity.MerchantWallet, error)
	GetSingleAdminWallet(ctx context.Context) (*entity.MerchantWallet, error) // New method for single admin wallet
	GetAdminWalletForUpdate(ctx context.Context, tx *sql.Tx) (*entity.MerchantWallet, error)
	GetSingleAdminWalletForUpdate(ctx context.Context, tx *sql.Tx) (*entity.MerchantWallet, error) // New method
	GetTotalAdminWalletAmount(ctx context.Context) (map[string]float64, error)

	// Atomic transaction processing methods (high-performance)
	ProcessDepositSuccess(ctx context.Context, merchantID uuid.UUID, merchantAmount float64, adminAmount float64) error
	ProcessWithdrawalSuccess(ctx context.Context, merchantID uuid.UUID, merchantAmount float64, adminAmount float64) error
	ProcessWithdrawalFailure(ctx context.Context, merchantID uuid.UUID, merchantAmount float64) error
	LockWithdrawalAmountAtomic(ctx context.Context, merchantID uuid.UUID, amount float64) error

	// Health check operations
	CheckWalletBalanceHealth(ctx context.Context) (*entity.WalletHealthCheck, error)
}
