package entity

import (
	"time"

	"github.com/google/uuid"
)

type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusInactive CustomerStatus = "inactive"
	CustomerStatusBlocked  CustomerStatus = "blocked"
)

type Customer struct {
	ID            uuid.UUID      `json:"id"`
	CustomerID    uuid.UUID      `json:"customer_id"`
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	Phone         string         `json:"phone"`
	Address       string         `json:"address"`
	LoyaltyPoints int            `json:"loyalty_points"`
	DateOfBirth   string         `json:"date_of_birth"`
	Status        CustomerStatus `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CreatedBy     uuid.UUID      `json:"created_by"`
	UpdatedBy     uuid.UUID      `json:"updated_by"`
	MerchantID    uuid.UUID      `json:"merchant_id"`
}

type CustomerResponse struct {
	ID            uuid.UUID      `json:"id"`
	CustomerID    uuid.UUID      `json:"customer_id"`
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	Phone         string         `json:"phone"`
	Address       string         `json:"address"`
	LoyaltyPoints int            `json:"loyalty_points"`
	DateOfBirth   string         `json:"date_of_birth"`
	Status        CustomerStatus `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CreatedBy     uuid.UUID      `json:"created_by"`
	UpdatedBy     uuid.UUID      `json:"updated_by"`
	MerchantID    uuid.UUID      `json:"merchant_id"`
}

type CreateCustomerRequest struct {
	MerchantID    string `json:"merchant_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	PhoneNumber   string `json:"phone_number,omitempty"`
	Address       string `json:"address,omitempty"`
	DateOfBirth   string `json:"date_of_birth,omitempty"`
	Status        string `json:"status,omitempty"`
	LoyaltyPoints int    `json:"loyalty_points,omitempty"`
}

type UpdateCustomerRequest struct {
	Name          *string `json:"name,omitempty"`
	Email         *string `json:"email,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Address       *string `json:"address,omitempty"`
	LoyaltyPoints *int    `json:"loyalty_points,omitempty"`
	DateOfBirth   *string `json:"date_of_birth,omitempty"`
	Status        *string `json:"status,omitempty"`
}

type GetCustomersParams struct {
	Text       string    `json:"text"`
	Skip       int       `json:"skip"`
	Take       int       `json:"take"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Status     string    `json:"status"`
	MerchantID uuid.UUID `json:"merchant_id"`
}

type CustomersResponse struct {
	Count     int                `json:"count"`
	Customers []CustomerResponse `json:"customers"`
}

