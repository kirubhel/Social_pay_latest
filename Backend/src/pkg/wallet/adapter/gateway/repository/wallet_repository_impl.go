package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	db "github.com/socialpay/socialpay/src/pkg/wallet/adapter/gateway/repository/generated"
	"github.com/socialpay/socialpay/src/pkg/wallet/core/entity"
)

type merchantWalletRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewWalletRepository(dbConn *sql.DB) WalletRepository {
	return &merchantWalletRepository{
		queries: db.New(dbConn),
		db:      dbConn,
	}
}

func (r *merchantWalletRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (r *merchantWalletRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (r *merchantWalletRepository) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

func (r *merchantWalletRepository) CreateMerchantWallet(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID, amount float64, lockedAmount float64, currency string) error {
	err := r.queries.CreateMerchantWallet(ctx, db.CreateMerchantWalletParams{
		ID:           uuid.New(),
		UserID:       userID,
		MerchantID:   merchantID,
		Amount:       amount,
		LockedAmount: lockedAmount,
		Currency:     currency,
		WalletType:   string(entity.WalletTypeMerchant),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *merchantWalletRepository) GetMerchantWalletByUserID(ctx context.Context, userID uuid.UUID) (*entity.MerchantWallet, error) {
	wallet, err := r.queries.GetMerchantWalletByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant wallet not found")
		}
		return nil, err
	}

	return &entity.MerchantWallet{
		ID:           wallet.ID,
		UserID:       wallet.UserID,
		MerchantID:   wallet.MerchantID,
		Amount:       wallet.Amount,
		LockedAmount: wallet.LockedAmount,
		Currency:     entity.Currency(wallet.Currency),
		WalletType:   entity.WalletType(wallet.WalletType),
		CreatedAt:    wallet.CreatedAt,
		UpdatedAt:    wallet.UpdatedAt,
	}, nil
}

func (r *merchantWalletRepository) GetMerchantWalletByMerchantID(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantWallet, error) {
	wallet, err := r.queries.GetMerchantWalletByMerchantID(ctx, merchantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant wallet not found")
		}
		return nil, err
	}

	return &entity.MerchantWallet{
		ID:           wallet.ID,
		UserID:       wallet.UserID,
		MerchantID:   wallet.MerchantID,
		Amount:       wallet.Amount,
		LockedAmount: wallet.LockedAmount,
		Currency:     entity.Currency(wallet.Currency),
		WalletType:   entity.WalletType(wallet.WalletType),
		CreatedAt:    wallet.CreatedAt,
		UpdatedAt:    wallet.UpdatedAt,
	}, nil
}

func (r *merchantWalletRepository) GetMerchantWalletByMerchantIDForUpdate(ctx context.Context, tx *sql.Tx, merchantID uuid.UUID) (*entity.MerchantWallet, error) {
	q := r.queries.WithTx(tx)
	wallet, err := q.GetMerchantWalletByMerchantIDForUpdate(ctx, merchantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant wallet not found")
		}
		return nil, err
	}

	return &entity.MerchantWallet{
		ID:           wallet.ID,
		UserID:       wallet.UserID,
		MerchantID:   wallet.MerchantID,
		Amount:       wallet.Amount,
		LockedAmount: wallet.LockedAmount,
		Currency:     entity.Currency(wallet.Currency),
		WalletType:   entity.WalletType(wallet.WalletType),
		CreatedAt:    wallet.CreatedAt,
		UpdatedAt:    wallet.UpdatedAt,
	}, nil
}

func (r *merchantWalletRepository) UpdateMerchantWallet(ctx context.Context, walletID uuid.UUID, amount float64, lockedAmount float64) error {
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

func (r *merchantWalletRepository) UpdateMerchantWalletAmountByMerchantID(ctx context.Context, merchantID uuid.UUID, amount float64) error {
	err := r.queries.UpdateMerchantWalletAmountByMerchantID(ctx, db.UpdateMerchantWalletAmountByMerchantIDParams{
		MerchantID: merchantID,
		Amount:     amount,
	})
	if err != nil {
		return err
	}
	return nil
}

// Admin wallet operations
func (r *merchantWalletRepository) GetAdminWallet(ctx context.Context) (*entity.MerchantWallet, error) {
	wallet, err := r.queries.GetAdminWallet(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin wallet not found")
		}
		return nil, err
	}

	return &entity.MerchantWallet{
		ID:           wallet.ID,
		UserID:       wallet.UserID,
		MerchantID:   wallet.MerchantID,
		Amount:       wallet.Amount,
		LockedAmount: wallet.LockedAmount,
		Currency:     entity.Currency(wallet.Currency),
		WalletType:   entity.WalletType(wallet.WalletType),
		CreatedAt:    wallet.CreatedAt,
		UpdatedAt:    wallet.UpdatedAt,
	}, nil
}

func (r *merchantWalletRepository) GetAdminWalletForUpdate(ctx context.Context, tx *sql.Tx) (*entity.MerchantWallet, error) {
	q := r.queries.WithTx(tx)
	wallet, err := q.GetAdminWalletForUpdate(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin wallet not found")
		}
		return nil, err
	}

	return &entity.MerchantWallet{
		ID:           wallet.ID,
		UserID:       wallet.UserID,
		MerchantID:   wallet.MerchantID,
		Amount:       wallet.Amount,
		LockedAmount: wallet.LockedAmount,
		Currency:     entity.Currency(wallet.Currency),
		WalletType:   entity.WalletType(wallet.WalletType),
		CreatedAt:    wallet.CreatedAt,
		UpdatedAt:    wallet.UpdatedAt,
	}, nil
}

// GetSingleAdminWallet gets the single admin wallet without requiring userID
func (r *merchantWalletRepository) GetSingleAdminWallet(ctx context.Context) (*entity.MerchantWallet, error) {
	query := `
		SELECT id, user_id, merchant_id, amount, locked_amount, currency, wallet_type, created_at, updated_at
		FROM merchant.wallet 
		WHERE wallet_type = 'super_admin'
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query)

	var wallet entity.MerchantWallet
	var walletType string
	var merchantID sql.NullString
	err := row.Scan(
		&wallet.ID,
		&wallet.UserID,
		&merchantID,
		&wallet.Amount,
		&wallet.LockedAmount,
		&wallet.Currency,
		&walletType,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin wallet not found")
		}
		return nil, err
	}

	wallet.WalletType = entity.WalletType(walletType)
	if merchantID.Valid {
		wallet.MerchantID, _ = uuid.Parse(merchantID.String)
	}

	return &wallet, nil
}

// GetSingleAdminWalletForUpdate gets the single admin wallet with row-level locking
func (r *merchantWalletRepository) GetSingleAdminWalletForUpdate(ctx context.Context, tx *sql.Tx) (*entity.MerchantWallet, error) {
	query := `
		SELECT id, user_id, merchant_id, amount, locked_amount, currency, wallet_type, created_at, updated_at
		FROM merchant.wallet 
		WHERE wallet_type = 'super_admin'
		FOR UPDATE
		LIMIT 1
	`
	row := tx.QueryRowContext(ctx, query)

	var wallet entity.MerchantWallet
	var walletType string
	var merchantID sql.NullString
	err := row.Scan(
		&wallet.ID,
		&wallet.UserID,
		&merchantID,
		&wallet.Amount,
		&wallet.LockedAmount,
		&wallet.Currency,
		&walletType,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin wallet not found")
		}
		return nil, err
	}

	wallet.WalletType = entity.WalletType(walletType)
	if merchantID.Valid {
		wallet.MerchantID, _ = uuid.Parse(merchantID.String)
	}

	return &wallet, nil
}

func (r *merchantWalletRepository) GetTotalAdminWalletAmount(ctx context.Context) (map[string]float64, error) {
	result, err := r.queries.GetTotalAdminWalletAmount(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"total_amount":        result.TotalAmount,
		"total_locked_amount": result.TotalLockedAmount,
	}, nil
}

// CheckWalletBalanceHealth verifies if wallet balances match transaction history
func (r *merchantWalletRepository) CheckWalletBalanceHealth(ctx context.Context) (*entity.WalletHealthCheck, error) {
	healthCheck := &entity.WalletHealthCheck{
		CheckedAt:     time.Now(),
		WalletDetails: []entity.WalletBalanceHealthDetail{},
		Summary:       entity.WalletHealthSummary{},
	}

	// Check merchant wallet health
	merchantWallets, err := r.getAllMerchantWallets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant wallets: %w", err)
	}

	for _, wallet := range merchantWallets {
		detail, err := r.checkMerchantWalletHealth(ctx, wallet)
		if err != nil {
			return nil, fmt.Errorf("failed to check merchant wallet health: %w", err)
		}
		healthCheck.WalletDetails = append(healthCheck.WalletDetails, *detail)
		healthCheck.Summary.TotalMerchantWallets++
		if detail.IsHealthy {
			healthCheck.Summary.HealthyMerchantWallets++
		}
		healthCheck.Summary.TotalBalanceDifference += math.Abs(detail.Difference)
		if math.Abs(detail.Difference) > healthCheck.Summary.LargestDiscrepancy {
			healthCheck.Summary.LargestDiscrepancy = math.Abs(detail.Difference)
		}
	}

	// Check admin wallet health (single admin wallet)
	adminWallets, err := r.getAllAdminWallets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin wallet: %w", err)
	}

	// Should only have 0 or 1 admin wallet
	if len(adminWallets) > 1 {
		return nil, fmt.Errorf("system error: found %d admin wallets, expected exactly 1", len(adminWallets))
	}

	if len(adminWallets) == 1 {
		wallet := adminWallets[0]
		detail, err := r.checkAdminWalletHealth(ctx, wallet)
		if err != nil {
			return nil, fmt.Errorf("failed to check admin wallet health: %w", err)
		}
		healthCheck.WalletDetails = append(healthCheck.WalletDetails, *detail)
		healthCheck.Summary.TotalAdminWallets = 1
		if detail.IsHealthy {
			healthCheck.Summary.HealthyAdminWallets = 1
		}
		healthCheck.Summary.TotalBalanceDifference += math.Abs(detail.Difference)
		if math.Abs(detail.Difference) > healthCheck.Summary.LargestDiscrepancy {
			healthCheck.Summary.LargestDiscrepancy = math.Abs(detail.Difference)
		}
	}

	// Calculate overall health
	healthCheck.TotalWallets = len(healthCheck.WalletDetails)
	healthCheck.HealthyWallets = healthCheck.Summary.HealthyMerchantWallets + healthCheck.Summary.HealthyAdminWallets
	healthCheck.UnhealthyWallets = healthCheck.TotalWallets - healthCheck.HealthyWallets
	healthCheck.IsHealthy = healthCheck.UnhealthyWallets == 0

	return healthCheck, nil
}

// Helper methods
func (r *merchantWalletRepository) getAllMerchantWallets(ctx context.Context) ([]entity.MerchantWallet, error) {
	query := `
		SELECT id, user_id, merchant_id, balance, currency, wallet_type, created_at, updated_at
		FROM public.wallets 
		WHERE type = 'merchant'
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []entity.MerchantWallet
	for rows.Next() {
		var wallet entity.MerchantWallet
		var walletType string
		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.MerchantID,
			&wallet.Balance,
			&wallet.Currency,
			&walletType,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallet.WalletType = entity.WalletType(walletType)
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func (r *merchantWalletRepository) getAllAdminWallets(ctx context.Context) ([]entity.MerchantWallet, error) {
	query := `
		SELECT id, user_id, merchant_id, balance, currency, wallet_type, created_at, updated_at
		FROM public.wallets 
		WHERE type = 'super_admin'
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query)

	var wallet entity.MerchantWallet
	var walletType string
	var merchantID sql.NullString
	err := row.Scan(
		&wallet.ID,
		&wallet.UserID,
		&merchantID,
		&wallet.Balance,
		&wallet.Currency,
		&walletType,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.MerchantWallet{}, nil // No admin wallet found
		}
		return nil, err
	}

	wallet.WalletType = entity.WalletType(walletType)
	if merchantID.Valid {
		wallet.MerchantID, _ = uuid.Parse(merchantID.String)
	}

	return []entity.MerchantWallet{wallet}, nil
}

func (r *merchantWalletRepository) checkMerchantWalletHealth(ctx context.Context, wallet entity.MerchantWallet) (*entity.WalletBalanceHealthDetail, error) {
	// Calculate deposits and withdrawals from transaction history
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type IN ('deposit', 'payment') AND status = 'SUCCESS' AND merchant_net > 0 THEN merchant_net ELSE 0 END), 0) as total_deposits,
			COALESCE(SUM(CASE WHEN type = 'withdrawal' AND status = 'SUCCESS' THEN ABS(merchant_net) ELSE 0 END), 0) as total_withdrawals,
			COUNT(*) as transaction_count
		FROM public.transactions 
		WHERE merchant_id = $1 AND status = 'SUCCESS' AND merchant_net IS NOT NULL
	`

	var totalDeposits, totalWithdrawals float64
	var transactionCount int

	err := r.db.QueryRowContext(ctx, query, wallet.MerchantID).Scan(&totalDeposits, &totalWithdrawals, &transactionCount)
	if err != nil {
		return nil, err
	}

	calculatedBalance := totalDeposits - totalWithdrawals
	difference := wallet.Balance - calculatedBalance
	isHealthy := math.Abs(difference) < 0.01 // Consider healthy if difference is less than 1 cent

	return &entity.WalletBalanceHealthDetail{
		WalletID:          wallet.ID,
		MerchantID:        wallet.MerchantID,
		WalletType:        string(wallet.WalletType),
		CurrentBalance:    wallet.Balance,
		CalculatedBalance: calculatedBalance,
		Difference:        difference,
		IsHealthy:         isHealthy,
		TotalDeposits:     totalDeposits,
		TotalWithdrawals:  totalWithdrawals,
		TransactionCount:  transactionCount,
	}, nil
}

func (r *merchantWalletRepository) checkAdminWalletHealth(ctx context.Context, wallet entity.MerchantWallet) (*entity.WalletBalanceHealthDetail, error) {
	// For the single admin wallet, calculate total commissions collected from all transactions
	query := `
		SELECT 
			COALESCE(SUM(admin_net), 0) as total_commissions,
			COUNT(*) as transaction_count
		FROM public.transactions 
		WHERE status = 'SUCCESS' AND admin_net IS NOT NULL AND admin_net > 0
	`

	var totalCommissions float64
	var transactionCount int

	err := r.db.QueryRowContext(ctx, query).Scan(&totalCommissions, &transactionCount)
	if err != nil {
		return nil, err
	}

	calculatedBalance := totalCommissions
	difference := wallet.Balance - calculatedBalance
	isHealthy := math.Abs(difference) < 0.01 // Consider healthy if difference is less than 1 cent

	return &entity.WalletBalanceHealthDetail{
		WalletID:          wallet.ID,
		MerchantID:        wallet.MerchantID,
		WalletType:        string(wallet.WalletType),
		CurrentBalance:    wallet.Balance,
		CalculatedBalance: calculatedBalance,
		Difference:        difference,
		IsHealthy:         isHealthy,
		TotalCommissions:  totalCommissions,
		TransactionCount:  transactionCount,
	}, nil
}

// Atomic transaction processing methods (high-performance)
// These methods use single SQL statements to update both merchant and admin wallets atomically

func (r *merchantWalletRepository) ProcessDepositSuccess(ctx context.Context, merchantID uuid.UUID, merchantAmount float64, adminAmount float64) error {
	// Single transaction with CTE for atomic updates of both wallets
	query := `
	WITH merchant_update AS (
		UPDATE merchant.wallet 
		SET amount = amount + $2,
			updated_at = NOW()
		WHERE merchant_id = $1
		RETURNING id
	),
	admin_update AS (
		UPDATE merchant.wallet 
		SET amount = amount + $3,
			updated_at = NOW()
		WHERE wallet_type = 'super_admin'
		RETURNING id
	)
	SELECT 
		(SELECT COUNT(*) FROM merchant_update) as merchant_updated,
		(SELECT COUNT(*) FROM admin_update) as admin_updated
	`

	var merchantUpdated, adminUpdated int
	err := r.db.QueryRowContext(ctx, query, merchantID, merchantAmount, adminAmount).Scan(&merchantUpdated, &adminUpdated)
	if err != nil {
		return fmt.Errorf("failed to process deposit success: %w", err)
	}

	if merchantUpdated == 0 {
		return fmt.Errorf("merchant wallet not found for merchantID: %s", merchantID)
	}
	if adminUpdated == 0 {
		return fmt.Errorf("admin wallet not found")
	}

	return nil
}

func (r *merchantWalletRepository) ProcessWithdrawalSuccess(ctx context.Context, merchantID uuid.UUID, merchantAmount float64, adminAmount float64) error {
	// Single transaction: unlock amount from merchant wallet and add commission to admin
	query := `
	WITH merchant_update AS (
		UPDATE merchant.wallet 
		SET locked_amount = locked_amount - $2,
			updated_at = NOW()
		WHERE merchant_id = $1
		RETURNING id
	),
	admin_update AS (
		UPDATE merchant.wallet 
		SET amount = amount + $3,
			updated_at = NOW()
		WHERE wallet_type = 'super_admin'
		RETURNING id
	)
	SELECT 
		(SELECT COUNT(*) FROM merchant_update) as merchant_updated,
		(SELECT COUNT(*) FROM admin_update) as admin_updated
	`

	var merchantUpdated, adminUpdated int
	err := r.db.QueryRowContext(ctx, query, merchantID, merchantAmount, adminAmount).Scan(&merchantUpdated, &adminUpdated)
	if err != nil {
		return fmt.Errorf("failed to process withdrawal success: %w", err)
	}

	if merchantUpdated == 0 {
		return fmt.Errorf("merchant wallet not found for merchantID: %s", merchantID)
	}
	if adminUpdated == 0 {
		return fmt.Errorf("admin wallet not found")
	}

	return nil
}

func (r *merchantWalletRepository) ProcessWithdrawalFailure(ctx context.Context, merchantID uuid.UUID, merchantAmount float64) error {
	// Single statement: return locked amount to available balance and unlock it
	query := `
	UPDATE merchant.wallet 
	SET amount = amount + $2,
		locked_amount = locked_amount - $2,
		updated_at = NOW()
	WHERE merchant_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, merchantID, merchantAmount)
	if err != nil {
		return fmt.Errorf("failed to process withdrawal failure: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("merchant wallet not found for merchantID: %s", merchantID)
	}

	return nil
}

func (r *merchantWalletRepository) LockWithdrawalAmountAtomic(ctx context.Context, merchantID uuid.UUID, amount float64) error {
	// Single atomic SQL operation to check balance and lock amount
	query := `
	WITH wallet_update AS (
		UPDATE merchant.wallet 
		SET amount = amount - $2,
			locked_amount = locked_amount + $2,
			updated_at = NOW()
		WHERE merchant_id = $1
			AND wallet_type = 'merchant'
			AND amount >= $2
		RETURNING id
	)
	SELECT COUNT(*) FROM wallet_update
	`

	var rowsAffected int64
	err := r.db.QueryRowContext(ctx, query, merchantID, amount).Scan(&rowsAffected)
	if err != nil {
		return fmt.Errorf("failed to lock withdrawal amount: %w", err)
	}

	if rowsAffected == 0 {
		// Get current balance to provide a helpful error message
		wallet, err := r.GetMerchantWalletByMerchantID(ctx, merchantID)
		if err != nil {
			return fmt.Errorf("insufficient funds or wallet not found")
		}
		return fmt.Errorf("insufficient funds: available balance is %.2f %s", wallet.Amount, wallet.Currency)
	}

	return nil
}
