package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	rbac_entity "github.com/socialpay/socialpay/src/pkg/rbac/core/entity"
	rbac_repository "github.com/socialpay/socialpay/src/pkg/rbac/core/repository"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/entity"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/repository"
	"github.com/lib/pq"
)

// TeamMemberRepositoryImpl implements the TeamMemberRepository interface
type TeamMemberRepositoryImpl struct {
	db             *sql.DB
	logger         *log.Logger
	rbacRepository rbac_repository.RBACRepository
}

// NewTeamMemberRepository creates a new PostgreSQL auth repository
func NewTeamMemberRepository(db *sql.DB, logger *log.Logger, rbacRepository rbac_repository.RBACRepository) repository.TeamMemberRepository {
	return &TeamMemberRepositoryImpl{
		db:             db,
		logger:         logger,
		rbacRepository: rbacRepository,
	}
}

// CreateGroup creates new group
func (r *TeamMemberRepositoryImpl) CreateGroup(ctx context.Context, req *entity.CreateGroupRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Creating new permissions and storing their id
	permissionsStr := make([]string, len(req.Permissions))
	for i, perm := range req.Permissions {
		var uuidPerms []uuid.UUID

		for _, operation := range perm.Operations {
			uuidPerms = append(uuidPerms, operation)
		}
		newPermission, err := r.rbacRepository.CreatePermission(ctx, &rbac_entity.CreatePermissionRequest{
			ResourceId: perm.ResourceId,
			Effect:     perm.Effect,
			Operations: uuidPerms,
		})

		if err != nil {
			return fmt.Errorf("failed to getting new permission: %w", err)
		}

		permissionsStr[i] = newPermission.String()
	}

	groupId := uuid.New()
	_, err = tx.ExecContext(ctx, `
			INSERT INTO auth."groups"(id, title, description, permissions, merchant_id, created_at, updated_at)
			VALUES($1, $2, $3, $4, $5, $6, $7);
			`, groupId, req.Title, req.Description, pq.Array(permissionsStr), req.MerchantID, time.Now(), time.Now())

	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// create new group permission
	for i := range len(req.Permissions) {
		permissionUUID, err := uuid.Parse(permissionsStr[i])
		if err != nil {
			return fmt.Errorf("failed to parse permission id: %w", err)
		}

		err = r.rbacRepository.CreateGroupPermission(ctx, groupId, permissionUUID)
		if err != nil {
			return fmt.Errorf("failed to create group permission: %w", err)
		}
	}

	return nil
}

// GroupExists checks whether group exists or not based on the given title and merchant id
func (r *TeamMemberRepositoryImpl) GroupExists(ctx context.Context, title string, merchantId *uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth.groups WHERE title = $1 AND merchant_id IS NOT DISTINCT FROM $2)`
	err := r.db.QueryRowContext(ctx, query, title, merchantId).Scan(&exists)
	return exists, err
}

// GetGroups queries list of groups from db
func (r *TeamMemberRepositoryImpl) GetGroups(ctx context.Context, merchantId *uuid.UUID) ([]entity.Group, error) {
	var groups []entity.Group
	var query string
	var rows *sql.Rows
	var err error

	if merchantId != nil {
		query = `SELECT id, title, permissions, merchant_id, description, created_at, updated_at FROM auth.groups WHERE merchant_id = $1;`
		rows, err = r.db.QueryContext(ctx, query, merchantId)
	} else {
		query = `SELECT id, title, permissions, merchant_id, description, created_at, updated_at FROM auth.groups WHERE merchant_id IS NULL;`
		rows, err = r.db.QueryContext(ctx, query)
	}

	if err != nil {
		fmt.Println("Error querying groups:", err)
		return nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group entity.Group
		var uuidPerms []uuid.UUID
		var permissions []rbac_entity.Permission

		err := rows.Scan(
			&group.ID,
			&group.Title,
			pq.Array(&uuidPerms),
			&group.MerchantID,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}

		for _, id := range uuidPerms {
			permission, err := r.rbacRepository.GetPermissionById(ctx, id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					// If permission not found, skip it
					continue
				}
				return nil, fmt.Errorf("failed to query permission by id: %w", err)
			}
			permissions = append(permissions, *permission)
		}

		group.Permissions = permissions
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return groups, nil
}

// GetGroupById fetches specific group by id and merchant id
func (r *TeamMemberRepositoryImpl) GetGroupById(ctx context.Context, merchantId *uuid.UUID, id uuid.UUID) (*entity.Group, error) {
	var group entity.Group
	var uuidPerms []uuid.UUID
	var permsissions []rbac_entity.Permission

	query := `SELECT id, title, permissions, merchant_id, description, created_at, updated_at FROM auth.groups WHERE id = $1 AND merchant_id IS NOT DISTINCT FROM $2;`

	err := r.db.QueryRowContext(ctx, query, id, merchantId).Scan(
		&group.ID,
		&group.Title,
		pq.Array(&uuidPerms),
		&group.MerchantID,
		&group.Description,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewTeamManagementError(entity.ErrGroupNotFound, entity.MsgGroupNotFound)
		}
		return nil, fmt.Errorf("failed to query group  by id: %w", err)
	}

	for _, id := range uuidPerms {
		permission, err := r.rbacRepository.GetPermissionById(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// If permission not found, skip it
				continue
			}
			return nil, fmt.Errorf("failed to query permission by id: %w", err)
		}

		permsissions = append(permsissions, *permission)
	}

	group.Permissions = permsissions

	return &group, nil
}

// UpdateGroup updates an existing group efficiently
func (r *TeamMemberRepositoryImpl) UpdateGroup(ctx context.Context, req *entity.UpdateGroupRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get current group to compare permissions
	currentGroup, err := r.GetGroupById(ctx, req.MerchantID, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get current group: %w", err)
	}

	// Handle permission updates efficiently
	var newPermissionIds []string
	if req.Permissions != nil {
		// Create map of existing permissions for efficient comparison
		existingPermMap := make(map[string]rbac_entity.Permission)
		for _, perm := range currentGroup.Permissions {
			// Create a key based on permission content for comparison
			key := fmt.Sprintf("%s_%s_%v", perm.ResourceId, perm.Effect, len(perm.Operations))
			existingPermMap[key] = perm
		}

		// Process new permissions
		for _, newPerm := range *req.Permissions {
			key := fmt.Sprintf("%s_%s_%v", newPerm.ResourceId, newPerm.Effect, len(newPerm.Operations))

			// Check if this permission configuration already exists
			if existingPerm, exists := existingPermMap[key]; exists {
				// Permission exists, check if operations are identical
				if r.areOperationsIdentical(existingPerm.Operations, newPerm.Operations) {
					// Reuse existing permission
					newPermissionIds = append(newPermissionIds, existingPerm.ID.String())
					delete(existingPermMap, key) // Remove from deletion list
					continue
				}
			}

			// Create new permission
			var uuidPerms []uuid.UUID
			uuidPerms = append(uuidPerms, newPerm.Operations...)

			newPermission, err := r.rbacRepository.CreatePermission(ctx, &rbac_entity.CreatePermissionRequest{
				ResourceId: newPerm.ResourceId,
				Effect:     newPerm.Effect,
				Operations: uuidPerms,
			})
			if err != nil {
				return fmt.Errorf("failed to create new permission: %w", err)
			}
			newPermissionIds = append(newPermissionIds, newPermission.String())
		}

		// Delete unused existing permissions
		for _, unusedPerm := range existingPermMap {
			err = r.rbacRepository.DeletePermissionById(ctx, unusedPerm.ID)
			if err != nil {
				return fmt.Errorf("failed to delete unused permission: %w", err)
			}
		}
	} else {
		// Keep existing permissions if not updating
		for _, perm := range currentGroup.Permissions {
			newPermissionIds = append(newPermissionIds, perm.ID.String())
		}
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *req.Title)
		argIndex++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if req.Permissions != nil {
		setParts = append(setParts, fmt.Sprintf("permissions = $%d", argIndex))
		args = append(args, pq.Array(newPermissionIds))
		argIndex++
	}

	// Always update updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE conditions
	args = append(args, req.ID, req.MerchantID)

	query := fmt.Sprintf(`
		UPDATE auth."groups" 
		SET %s 
		WHERE id = $%d AND merchant_id IS NOT DISTINCT FROM $%d
	`, strings.Join(setParts, ", "), argIndex, argIndex+1)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// areOperationsIdentical compares two operation slices for equality
func (r *TeamMemberRepositoryImpl) areOperationsIdentical(existing []rbac_entity.Operation, new []uuid.UUID) bool {
	if len(existing) != len(new) {
		return false
	}

	// Create maps for efficient comparison
	existingIds := make(map[uuid.UUID]bool)
	for _, op := range existing {
		existingIds[op.Id] = true
	}

	for _, newId := range new {
		if !existingIds[newId] {
			return false
		}
	}

	return true
}

// DeleteGroupById deletes specific group by id
func (r *TeamMemberRepositoryImpl) DeleteGroupById(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE from auth.groups WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("failed to delete group by id: %w", err)
	}

	return nil
}

// DeleteGroupPermissionById deletes permission in the group
func (r *TeamMemberRepositoryImpl) DeleteGroupPermissionById(ctx context.Context, id uuid.UUID) error {
	err := r.rbacRepository.DeletePermissionById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete permission by id: %w", err)
	}

	return nil
}

// GetUserGroups fetches user groups based on merchant id and group id
func (r *TeamMemberRepositoryImpl) GetUserGroups(ctx context.Context, merchantId *uuid.UUID, groupId uuid.UUID) ([]entity.UserGroup, error) {
	var userGroups []entity.UserGroup
	query := `SELECT id, user_id, group_id, user_type, merchant_id, created_at, updated_at FROM auth.user_groups WHERE group_id = $1 AND merchant_id IS NOT DISTINCT FROM $2`
	rows, err := r.db.QueryContext(ctx, query, groupId, merchantId)

	if err != nil {
		return nil, fmt.Errorf("failed to delete permission by id: %w", err)
	}

	for rows.Next() {
		var userGroup entity.UserGroup
		err := rows.Scan(
			&userGroup.ID,
			&userGroup.UserId,
			&userGroup.GroupId,
			&userGroup.UserType,
			&userGroup.MerchantID,
			&userGroup.CreatedAt,
			&userGroup.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}

		userGroups = append(userGroups, userGroup)
	}

	return userGroups, err
}

// GetUserGroupById fetches a specific user group by user group ID and merchant ID
func (r *TeamMemberRepositoryImpl) GetUserGroupById(ctx context.Context, merchantId *uuid.UUID, userGroupId uuid.UUID) (*entity.UserGroup, error) {
	var userGroup entity.UserGroup
	query := `SELECT id, user_id, group_id, user_type, merchant_id, created_at, updated_at FROM auth.user_groups WHERE id = $1 AND merchant_id IS NOT DISTINCT FROM $2`

	err := r.db.QueryRowContext(ctx, query, userGroupId, merchantId).Scan(
		&userGroup.ID,
		&userGroup.UserId,
		&userGroup.GroupId,
		&userGroup.UserType,
		&userGroup.MerchantID,
		&userGroup.CreatedAt,
		&userGroup.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil, nil when user group not found
		}
		return nil, fmt.Errorf("failed to get user group by id: %w", err)
	}

	return &userGroup, nil
}
