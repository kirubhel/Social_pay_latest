package entity

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	Id          uuid.UUID `json:"id"`
	MerchantID  uuid.UUID `json:"merchant_id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uuid.UUID `json:"created_by"`
	UpdatedBy   uuid.UUID `json:"updated_by"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
}
