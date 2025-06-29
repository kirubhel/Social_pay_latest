package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	db "github.com/socialpay/socialpay/src/pkg/wallet/adapter/gateway/repository/generated"
	"github.com/socialpay/socialpay/src/pkg/wallet/core/entity"
)

type walletRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewWalletRepository(dbConn *sql.DB) WalletRepository {
	return &walletRepository{
		queries: db.New(dbConn),
		db:      dbConn,
	}
}

func (r *walletRepository) UpdateMerchantWallet(ctx context.Context, walletID uuid.UUID, amount float64, lockedAmount float64) error {
	err := r.queries.UpdateMerchantWallet(ctx, db.UpdateMerchantWalletParams{
		ID:           walletID,
		Amount:       amount,
		LockedAmount: lockedAmount,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *walletRepository) CreateMerchantWallet(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID, amount float64, lockedAmount float64, currency string) error {
	err := r.queries.CreateMerchantWallet(ctx, db.CreateMerchantWalletParams{
		ID:           uuid.New(),
		UserID:       userID,
		MerchantID:   merchantID,
		Amount:       amount,
		LockedAmount: lockedAmount,
		Currency:     currency,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *walletRepository) GetMerchantWallet(ctx context.Context, merchantID uuid.UUID) (*entity.Wallet, error) {
	wallet, err := r.queries.GetMerchantWalletByMerchantID(ctx, merchantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant wallet not found")
		}
		return nil, err
	}

	return &entity.Wallet{
		ID:           wallet.ID,
		UserID:       wallet.UserID,
		Amount:       wallet.Amount,
		LockedAmount: wallet.LockedAmount,
		Currency:     wallet.Currency,
		CreatedAt:    wallet.CreatedAt,
		UpdatedAt:    wallet.UpdatedAt,
	}, nil
}

// Begin a database transaction
func (r *walletRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable, // Highest isolation level to prevent race conditions
	})
}

// Commit a database transaction
func (r *walletRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

// Rollback a database transaction
func (r *walletRepository) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

// Get merchant wallet with FOR UPDATE lock to prevent race conditions
func (r *walletRepository) GetMerchantWalletForUpdate(ctx context.Context, tx *sql.Tx, merchantID uuid.UUID) (*entity.Wallet, error) {
	// Create a queries instance that uses the transaction
	q := r.queries.WithTx(tx)

	wallet, err := q.GetMerchantWalletByMerchantID(ctx, merchantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant wallet not found")
		}
		return nil, err
	}

	return &entity.Wallet{
		ID:           wallet.ID,
		UserID:       wallet.UserID,
		Amount:       wallet.Amount,
		LockedAmount: wallet.LockedAmount,
		Currency:     wallet.Currency,
		CreatedAt:    wallet.CreatedAt,
		UpdatedAt:    wallet.UpdatedAt,
	}, nil
}

// Update merchant wallet within a transaction
func (r *walletRepository) UpdateMerchantWalletWithTx(ctx context.Context, tx *sql.Tx, walletID uuid.UUID, amount float64, lockedAmount float64) error {
	// Create a queries instance that uses the transaction
	q := r.queries.WithTx(tx)

	err := q.UpdateMerchantWallet(ctx, db.UpdateMerchantWalletParams{
		ID:           walletID,
		Amount:       amount,
		LockedAmount: lockedAmount,
	})
	if err != nil {
		return err
	}
	return nil
}
