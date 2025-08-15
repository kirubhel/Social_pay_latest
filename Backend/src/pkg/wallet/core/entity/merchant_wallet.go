package entity

import (
	"time"

	"github.com/google/uuid"
)

type WalletType string

const (
	WalletTypeMerchant WalletType = "merchant"
	WalletTypeAdmin    WalletType = "admin"
)

type Currency string

const (
	CurrencyETB Currency = "ETB"
	CurrencyUSD Currency = "USD"
)

type MerchantWallet struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	MerchantID   uuid.UUID  `json:"merchant_id"`
	Balance      float64    `json:"balance"`
	Amount       float64    `json:"amount"`
	LockedAmount float64    `json:"locked_amount"`
	Description  string     `json:"description,omitempty"`
	Currency     Currency   `json:"currency"`
	WalletType   WalletType `json:"wallet_type"`
	IsActive     bool       `json:"is_active"`
	LastSyncAt   time.Time  `json:"last_sync_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// WalletHealthCheck represents the health status of wallet balances vs transaction history
type WalletHealthCheck struct {
	IsHealthy        bool                        `json:"is_healthy"`
	CheckedAt        time.Time                   `json:"checked_at"`
	TotalWallets     int                         `json:"total_wallets"`
	HealthyWallets   int                         `json:"healthy_wallets"`
	UnhealthyWallets int                         `json:"unhealthy_wallets"`
	WalletDetails    []WalletBalanceHealthDetail `json:"wallet_details"`
	Summary          WalletHealthSummary         `json:"summary"`
}

// WalletBalanceHealthDetail represents detailed health info for a specific wallet
type WalletBalanceHealthDetail struct {
	WalletID          uuid.UUID `json:"wallet_id"`
	MerchantID        uuid.UUID `json:"merchant_id"`
	WalletType        string    `json:"wallet_type"`
	CurrentBalance    float64   `json:"current_balance"`
	CalculatedBalance float64   `json:"calculated_balance"`
	Difference        float64   `json:"difference"`
	IsHealthy         bool      `json:"is_healthy"`
	TotalDeposits     float64   `json:"total_deposits"`
	TotalWithdrawals  float64   `json:"total_withdrawals"`
	TotalCommissions  float64   `json:"total_commissions,omitempty"` // For admin wallets
	TransactionCount  int       `json:"transaction_count"`
}

// WalletHealthSummary provides aggregated health statistics
type WalletHealthSummary struct {
	TotalMerchantWallets   int     `json:"total_merchant_wallets"`
	HealthyMerchantWallets int     `json:"healthy_merchant_wallets"`
	TotalAdminWallets      int     `json:"total_admin_wallets"`
	HealthyAdminWallets    int     `json:"healthy_admin_wallets"`
	TotalBalanceDifference float64 `json:"total_balance_difference"`
	LargestDiscrepancy     float64 `json:"largest_discrepancy"`
}
