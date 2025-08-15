package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductCatalog struct {
	Id           uuid.UUID `json:"id"`
	MerchantID   uuid.UUID `json:"merchant_id"`
	CatalogID    uuid.UUID `json:"catalog_id"`
	ProductID    uuid.UUID `json:"product_id"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    uuid.UUID `json:"created_by"`
	UpdatedBy    uuid.UUID `json:"updated_by"`
}
