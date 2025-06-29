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
	userPermissions, err := uc.repo.ListUserPermissions(userID)
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
