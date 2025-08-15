package entity

import (
	"github.com/google/uuid"
)

type PaymentMethod struct {
	Id         uuid.UUID `json:"id"`
	MerchantID string    `json:"merchant_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Comission  float64   `json:"comission"`
	Details    string    `json:"details"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
	CreatedBy  uuid.UUID `json:"created_by"`
	UpdatedBy  uuid.UUID `json:"updated_by"`
}
