package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/types"
)

type CatalogStatus string

const (
	CatalogStatusActive   CatalogStatus = "active"
	CatalogStatusInactive CatalogStatus = "inactive"
	CatalogStatusDraft    CatalogStatus = "draft"
)

type Catalog struct {
	ID          uuid.UUID        `json:"id"`
	MerchantID  uuid.UUID        `json:"merchant_id"`
	Name        string           `json:"name"`
	Description types.NullString `json:"description"`
	Status      CatalogStatus    `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	CreatedBy   uuid.UUID        `json:"created_by"`
	UpdatedBy   uuid.UUID        `json:"updated_by"`
}

type CatalogResponse struct {
	ID          uuid.UUID        `json:"id"`
	MerchantID  uuid.UUID        `json:"merchant_id"`
	Name        string           `json:"name"`
	Description types.NullString `json:"description"`
	Status      CatalogStatus    `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	CreatedBy   uuid.UUID        `json:"created_by"`
	UpdatedBy   uuid.UUID        `json:"updated_by"`
}

type CreateCatalogRequest struct {
	Name        string           `json:"name" binding:"required"`
	Description types.NullString `json:"description"`
	Status      string           `json:"status"`
}

type UpdateCatalogRequest struct {
	Name        *string           `json:"name,omitempty"`
	Description *types.NullString `json:"description,omitempty"`
	Status      *CatalogStatus    `json:"status,omitempty"`
}

type GetCatalogsParams struct {
	Text       string    `json:"text"`
	Skip       int       `json:"skip"`
	Take       int       `json:"take"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Status     string    `json:"status"`
	MerchantID uuid.UUID `json:"merchant_id"`
}

type CatalogsResponse struct {
	Count    int               `json:"count"`
	Catalogs []CatalogResponse `json:"catalogs"`
}
