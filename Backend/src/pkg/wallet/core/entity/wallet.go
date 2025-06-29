package entity

import (
	"time"

	"github.com/google/uuid"
)

// Wallet represents a merchant's wallet
// @Description Wallet information for a merchant
type Wallet struct {
	// @Description Unique identifier for the wallet
	ID uuid.UUID `json:"id"`
	// @Description ID of the merchant who owns this wallet
	MerchantID uuid.UUID `json:"merchant_id"`
	// @Description Current balance in the wallet
	Balance float64 `json:"balance"`
	// @Description Currency code (e.g., ETB, USD)
	Currency string `json:"currency"`
	// @Description When the wallet was created
	CreatedAt time.Time `json:"created_at"`
	// @Description When the wallet was last updated
	UpdatedAt time.Time `json:"updated_at"`
	// @Description When the wallet was last synchronized
	LastSyncAt time.Time `json:"last_sync_at"`
	// @Description Whether the wallet is active
	IsActive bool `json:"is_active"`
	// @Description Optional description of the wallet
	Description string `json:"description,omitempty"`
	// @Description Locked amount in the wallet that cannot be used for transactions
	LockedAmount float64 `json:"locked_amount,omitempty"`
	// @Description Amount total amount in the wallet that is available for transactions
	Amount float64 `json:"amount,omitempty"`
	// @Description User ID associated with the wallet
	UserID uuid.UUID `json:"user_id,omitempty"`
}
