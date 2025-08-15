package entity

import (
	"time"

	"github.com/google/uuid"
)

// Resource represents resources in the system
type Resource struct {
	Id          uuid.UUID   `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Operations  []Operation `json:"operations" db:"operations"`
	Description *string     `json:"description" db:"description"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// Operation represents operations in the system
type Operation struct {
	Id               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      *string   `json:"description" db:"description"`
	IsAdminOperation bool      `json:"is_admin_operation" db:"is_admin_operation"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Permission represents a permission
type Permission struct {
	ID         uuid.UUID   `json:"id" db:"id"`
	ResourceId uuid.UUID   `json:"resource_id" db:"resource_id"`
	Operations []Operation `json:"operations" db:"operations"`
	Effect     string      `json:"effect" db:"effect"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
}

// CreatePermissionRequest represents a request to create a new permission
type CreatePermissionRequest struct {
	ResourceId uuid.UUID   `json:"resource_id" validate:"required"`
	Operations []uuid.UUID `json:"operations" validate:"required"`
	Effect     string      `json:"effect" validate:"required"`
}
