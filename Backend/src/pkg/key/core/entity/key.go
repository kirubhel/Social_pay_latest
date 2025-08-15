// models/key.go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type KeyPair struct {
	PrivateKey string
	PublicKey  string
}

type APIKey struct {
	ID         uuid.UUID
	MerchantID string
	PrivateKey string
	PublicKey  string
	APIKey     string
	Service    string
	ExpiryDate time.Time
	Store      string
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Merchant struct {
	ID          string
	CompanyName string
}
