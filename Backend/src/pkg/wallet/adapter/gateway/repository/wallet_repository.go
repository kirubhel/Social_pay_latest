package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/wallet/core/entity"
)

type WalletRepository interface {
	UpdateMerchantWallet(ctx context.Context, walletID uuid.UUID, amount float64, lockedAmount float64) error
	CreateMerchantWallet(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID, amount float64, lockedAmount float64, currency string) error
	GetMerchantWallet(ctx context.Context, merchantID uuid.UUID) (*entity.Wallet, error)

	GetMerchantWalletForUpdate(ctx context.Context, tx *sql.Tx, merchantID uuid.UUID) (*entity.Wallet, error)
	UpdateMerchantWalletWithTx(ctx context.Context, tx *sql.Tx, walletID uuid.UUID, amount float64, lockedAmount float64) error

	// Begin a database transaction
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// Commit a database transaction
	CommitTx(tx *sql.Tx) error

	// Rollback a database transaction
	RollbackTx(tx *sql.Tx) error
}
