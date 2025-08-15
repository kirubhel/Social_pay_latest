package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusArchived ProductStatus = "archived"
)

type Product struct {
	ID          uuid.UUID     `json:"id"`
	MerchantID  uuid.UUID     `json:"merchant_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Currency    string        `json:"currency"`
	SKU         string        `json:"sku"`
	Weight      float64       `json:"weight"`
	Dimensions  string        `json:"dimensions"`
	ImageURL    string        `json:"image_url"`
	Status      ProductStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	CreatedBy   uuid.UUID     `json:"created_by"`
	UpdatedBy   uuid.UUID     `json:"updated_by"`
}

type ProductResponse struct {
	ID          uuid.UUID     `json:"id"`
	MerchantID  uuid.UUID     `json:"merchant_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Currency    string        `json:"currency"`
	SKU         string        `json:"sku"`
	Weight      float64       `json:"weight"`
	Dimensions  string        `json:"dimensions"`
	ImageURL    string        `json:"image_url"`
	Status      ProductStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	CreatedBy   uuid.UUID     `json:"created_by"`
	UpdatedBy   uuid.UUID     `json:"updated_by"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Currency    string  `json:"currency" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
	Weight      float64 `json:"weight"`
	Dimensions  string  `json:"dimensions"`
	ImageURL    string  `json:"image_url"`
	Status      string  `json:"status"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Currency    *string  `json:"currency,omitempty"`
	SKU         *string  `json:"sku,omitempty"`
	Weight      *float64 `json:"weight,omitempty"`
	Dimensions  *string  `json:"dimensions,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

type GetProductsParams struct {
	Text       string    `json:"text"`
	Skip       int       `json:"skip"`
	Take       int       `json:"take"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Status     string    `json:"status"`
	MerchantID uuid.UUID `json:"merchant_id"`
}

type ProductsResponse struct {
	Count    int               `json:"count"`
	Products []ProductResponse `json:"products"`
}


