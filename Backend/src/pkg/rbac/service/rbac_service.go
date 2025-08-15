package service

import (
	"context"
	"log"

	"github.com/socialpay/socialpay/src/pkg/rbac/core/entity"
	"github.com/socialpay/socialpay/src/pkg/rbac/core/repository"
	"github.com/socialpay/socialpay/src/pkg/rbac/core/service"
)

// RBACServiceImpl implements the RBACService interface
type RBACServiceImpl struct {
	repo      repository.RBACRepository
	jwtSecret string
	logger    *log.Logger
}

// NewRBACService creates a new rbac service
func NewRBACService(repo repository.RBACRepository, jwtSecret string, logger *log.Logger) service.RBACService {
	return &RBACServiceImpl{
		repo:      repo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *RBACServiceImpl) GetResources(ctx context.Context, isAdmin bool) ([]entity.Resource, error) {
	return s.repo.GetResources(ctx, isAdmin)
}
