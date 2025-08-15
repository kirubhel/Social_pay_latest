package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod struct {
	ID         uuid.UUID `json:"id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Commission float64   `json:"commission"`
	Details    string    `json:"details"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  uuid.UUID `json:"created_by"`
	UpdatedBy  uuid.UUID `json:"updated_by"`
}

type PaymentMethodResponse struct {
	ID         uuid.UUID `json:"id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Commission float64   `json:"commission"`
	Details    string    `json:"details"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  uuid.UUID `json:"created_by"`
	UpdatedBy  uuid.UUID `json:"updated_by"`
}

type CreatePaymentMethodRequest struct {
	Name       string  `json:"name" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Commission float64 `json:"commission"`
	Details    string  `json:"details"`
	IsActive   bool    `json:"is_active"`
}

type UpdatePaymentMethodRequest struct {
	Name       *string  `json:"name,omitempty"`
	Type       *string  `json:"type,omitempty"`
	Commission *float64 `json:"commission,omitempty"`
	Details    *string  `json:"details,omitempty"`
	IsActive   *bool    `json:"is_active,omitempty"`
}

type GetPaymentMethodsParams struct {
	Text       string    `json:"text"`
	Skip       int       `json:"skip"`
	Take       int       `json:"take"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	MerchantID uuid.UUID `json:"merchant_id"`
}

type PaymentMethodsResponse struct {
	Count         int                     `json:"count"`
	PaymentMethods []PaymentMethodResponse `json:"payment_methods"`
}


