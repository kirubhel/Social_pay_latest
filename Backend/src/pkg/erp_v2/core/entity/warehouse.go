package entity

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseStatus string

const (
	WarehouseStatusActive   WarehouseStatus = "active"
	WarehouseStatusInactive WarehouseStatus = "inactive"
)

type Warehouse struct {
	ID          uuid.UUID       `json:"id"`
	MerchantID  uuid.UUID       `json:"merchant_id"`
	Name        string          `json:"name"`
	Location    string          `json:"location"`
	Capacity    int             `json:"capacity"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	UpdatedBy   uuid.UUID       `json:"updated_by"`
	IsActive    bool            `json:"is_active"`
	Description string          `json:"description"`
	Status      WarehouseStatus `json:"status"`
}

type WarehouseResponse struct {
	ID          uuid.UUID       `json:"id"`
	MerchantID  uuid.UUID       `json:"merchant_id"`
	Name        string          `json:"name"`
	Location    string          `json:"location"`
	Capacity    int             `json:"capacity"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	UpdatedBy   uuid.UUID       `json:"updated_by"`
	IsActive    bool            `json:"is_active"`
	Description string          `json:"description"`
	Status      WarehouseStatus `json:"status"`
}

type CreateWarehouseRequest struct {
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Capacity    int    `json:"capacity"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateWarehouseRequest struct {
	Name        *string `json:"name,omitempty"`
	Location    *string `json:"location,omitempty"`
	Capacity    *int    `json:"capacity,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type GetWarehousesParams struct {
	Text       string    `json:"text"`
	Skip       int       `json:"skip"`
	Take       int       `json:"take"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Status     string    `json:"status"`
	MerchantID uuid.UUID `json:"merchant_id"`
}

type WarehousesResponse struct {
	Count      int                 `json:"count"`
	Warehouses []WarehouseResponse `json:"warehouses"`
}


