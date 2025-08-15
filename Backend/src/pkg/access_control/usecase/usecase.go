package usecase

import (
	"log"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/access_control/core/entity"
)

type Error struct {
	Type    string
	Message string
}

func (err Error) Error() string {
	return err.Message
}

type Usecase struct {
	log  *log.Logger
	repo Repo
}

func New(log *log.Logger, repo Repo) Interactor {
	return Usecase{log: log, repo: repo}
}

func (uc Usecase) CheckUserPermission(userID uuid.UUID, resourceName string, operationName string) bool {
	// 1. Get all groups for the user
	groups, err := uc.repo.ListUserGroups(userID)
	if err != nil {
		uc.log.Println("CHECK PERMISSION ERROR: Failed to list user groups", err)
		return false
	}

	// 2. Collect all permissions from user groups
	var groupPermissions []entity.Permission
	for _, group := range groups {
		perms, err := uc.repo.ListGroupPermissions(group.ID)
		if err != nil {
			uc.log.Println("CHECK PERMISSION ERROR: Failed to list group permissions", err)
			continue
		}
		groupPermissions = append(groupPermissions, perms...)
	}

	// 3. Collect all direct user permissions
	userPermissions, err := uc.repo.ListUserPermissions(userID, resourceName)
	if err != nil {
		uc.log.Println("CHECK PERMISSION ERROR: Failed to list user permissions", err)
	}

	allPermissions := append(groupPermissions, userPermissions...)

	// 4. Check for allow/deny
	allowed := false
	for _, perm := range allPermissions {
		if perm.ResourceName == resourceName && perm.OperationName == operationName {
			if perm.Effect == "deny" {
				return false
			}
			if perm.Effect == "allow" {
				allowed = true
			}
		}
	}
	return allowed
}

func (uc Usecase) CheckPermission(userID uuid.UUID, requiredPermission entity.Permission) (bool, error) {
	// First check direct user permissions
	permissions, err := uc.repo.ListUserPermissions(userID, requiredPermission.ResourceName)
	if err != nil {
		uc.log.Printf("[ERROR] Failed to fetch user permissions: %v", err)
		return false, err
	}

	// Check if any direct permission matches
	for _, permission := range permissions {
		if permission.ResourceName == requiredPermission.ResourceName &&
			permission.OperationName == requiredPermission.OperationName &&
			permission.Effect == requiredPermission.Effect {
			uc.log.Printf("[DEBUG] User %s has direct permission to perform the operation", userID)
			return true, nil
		}
	}

	// If no direct permission found, check group permissions
	groups, err := uc.repo.ListUserGroups(userID)
	if err != nil {
		uc.log.Printf("[ERROR] Failed to fetch user groups: %v", err)
		return false, err
	}

	for _, group := range groups {
		groupPermissions, err := uc.repo.ListGroupPermissions(group.ID)
		if err != nil {
			uc.log.Printf("[ERROR] Failed to fetch permissions for group %s: %v", group.ID, err)
			continue
		}

		for _, permission := range groupPermissions {
			if permission.ResourceName == requiredPermission.ResourceName &&
				permission.OperationName == requiredPermission.OperationName &&
				permission.Effect == requiredPermission.Effect {
				uc.log.Printf("[DEBUG] User %s has permission via group %s", userID, group.Title)
				return true, nil
			}
		}
	}

	uc.log.Printf("[DEBUG] Permission check failed for user %s: No matching permissions found", userID)
	return false, nil
}
