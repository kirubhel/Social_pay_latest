package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/rbac/core/entity"
	"github.com/socialpay/socialpay/src/pkg/rbac/core/repository"
	"github.com/lib/pq"
)

// RBACRepositoryImpl implements the RBACRepository interface
type RBACRepositoryImpl struct {
	db     *sql.DB
	logger *log.Logger
}

// NewRBACRepository creates a new PostgreSQL rbac repository
func NewRBACRepository(db *sql.DB, logger *log.Logger) repository.RBACRepository {
	return &RBACRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// GetOperationById queries operation from the db by given id
func (r *RBACRepositoryImpl) GetOperationById(ctx context.Context, id uuid.UUID) (*entity.Operation, error) {
	var operation entity.Operation

	query := `SELECT id, name, description, created_at, updated_at FROM auth.operations WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&operation.Id,
		&operation.Name,
		&operation.Description,
		&operation.CreatedAt,
		&operation.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query operation by id: %w", err)
	}

	return &operation, nil
}

// GetResources queries all the resources from the db
func (r *RBACRepositoryImpl) GetResources(ctx context.Context, isAdmin bool) ([]entity.Resource, error) {
	var resources []entity.Resource

	query := `SELECT * FROM auth.resources;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query resources: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var resource entity.Resource
		var uuidOps []uuid.UUID
		var operations []entity.Operation

		err := rows.Scan(
			&resource.Id,
			&resource.Name,
			pq.Array(&uuidOps),
			&resource.Description,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan resource: %w", err)
		}

		if resource.Name == "ALL" {
			continue
		}

		if resource.Name == "ADMIN_ALL" && !isAdmin {
			continue
		}

		for _, id := range uuidOps {
			operation, err := r.GetOperationById(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("row iteration error: %w", err)
			}

			if strings.Contains(operation.Name, "ADMIN") {
				if !isAdmin {
					continue
				}
			}
			operations = append(operations, *operation)
		}

		resource.Operations = operations
		resources = append(resources, resource)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return resources, nil
}

// CreatePermission creates new permission
func (r *RBACRepositoryImpl) CreatePermission(ctx context.Context, req *entity.CreatePermissionRequest) (*uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	permissionId := uuid.New()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO auth.permissions
			(id, operations, effect, created_at, updated_at, resource_id)
			VALUES($1, $2, $3, $4, $5, $6);
		`, permissionId, pq.Array(req.Operations), req.Effect, time.Now(), time.Now(), req.ResourceId)

	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &permissionId, nil
}

// GetPermissionById fetches specific permission by id from db
func (r *RBACRepositoryImpl) GetPermissionById(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	var permission entity.Permission
	var uuidOps []uuid.UUID
	var operations []entity.Operation

	query := `SELECT id, resource_id, operations, effect, created_at, updated_at FROM auth.permissions WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&permission.ID,
		&permission.ResourceId,
		pq.Array(&uuidOps),
		&permission.Effect,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query permission by id: %w", err)
	}

	for _, id := range uuidOps {
		operation, err := r.GetOperationById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("row iteration error: %w", err)
		}
		operations = append(operations, *operation)
	}

	permission.Operations = operations

	return &permission, nil
}

// DeletePermissionById deletes specific permission with the given id from the db
func (r *RBACRepositoryImpl) DeletePermissionById(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, `DELETE FROM auth.permissions WHERE id = $1;`, id)

	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CreateGroupPermission create new group permission
func (r *RBACRepositoryImpl) CreateGroupPermission(ctx context.Context, groupID uuid.UUID, permissionID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO auth.group_permissions(id, group_id, permission_id, created_at, updated_at) VALUES($1, $2, $3, $4, $5);`

	id := uuid.New()
	_, err = tx.ExecContext(ctx, query, id, groupID, permissionID, time.Now(), time.Now())

	if err != nil {
		return fmt.Errorf("failed to create group permission: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
