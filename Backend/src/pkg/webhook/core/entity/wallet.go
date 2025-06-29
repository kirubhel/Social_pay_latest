package entity

import (
	"time"

	"github.com/google/uuid"
)

type Currency string

const (
	CurrencyETB Currency = "ETB"
	CurrencyUSD Currency = "USD"
)

//TODO: Add merchant_id
type Wallet struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Amount       float64   `json:"amount"`
	LockedAmount float64   `json:"locked_amount"`
	Currency     Currency  `json:"currency"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
} 