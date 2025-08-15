package service

import (
	"context"

	"github.com/socialpay/socialpay/src/pkg/rbac/core/entity"
)

// RBACService defines the interface for rbac business logic
type RBACService interface {
	// Resource
	GetResources(ctx context.Context, isAdmin bool) ([]entity.Resource, error)
}
