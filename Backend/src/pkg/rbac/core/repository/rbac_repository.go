package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/rbac/core/entity"
)

// RBACRepository defines the interface for resource and operation data operations
type RBACRepository interface {
	// Resource operations
	GetResources(ctx context.Context, isAdmin bool) ([]entity.Resource, error)

	// Operation
	GetOperationById(ctx context.Context, id uuid.UUID) (*entity.Operation, error)

	// Permission operations
	CreatePermission(ctx context.Context, req *entity.CreatePermissionRequest) (*uuid.UUID, error)
	GetPermissionById(ctx context.Context, id uuid.UUID) (*entity.Permission, error)
	DeletePermissionById(ctx context.Context, id uuid.UUID) error

	// Group Permission operations
	CreateGroupPermission(ctx context.Context, groupID uuid.UUID, permissionID uuid.UUID) error
}
