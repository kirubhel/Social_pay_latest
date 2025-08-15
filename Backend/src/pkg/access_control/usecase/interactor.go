package usecase

import (
	"github.com/socialpay/socialpay/src/pkg/access_control/core/entity"

	"github.com/google/uuid"
)

type Interactor interface {
	// Permission
	GrantPermissionToUser(userID uuid.UUID, permissionID uuid.UUID) error
	RevokePermissionFromUser(userID uuid.UUID, permissionID uuid.UUID) error
	GrantPermissionToGroup(groupID uuid.UUID, permissionID uuid.UUID) error
	RevokePermissionFromGroup(groupID uuid.UUID, permissionID uuid.UUID) error
	ListUserPermissions(userID uuid.UUID) ([]entity.Permission, error)
	ListGroupPermissions(groupID uuid.UUID) ([]entity.Permission, error)
	CreatePermission(resourceID string, resource uuid.UUID, operation []uuid.UUID, effect string) (*entity.Permission, error)
	DeletePermission(resourceID uuid.UUID, permissionID uuid.UUID) error
	UpdatePermission(permissionID uuid.UUID, resourceID string, resource uuid.UUID, operation []uuid.UUID, effect string) (*entity.Permission, error)
	ListPermissions() ([]entity.Permission, error)
	CheckPermission(userID uuid.UUID, requiredPermission entity.Permission) (bool, error)

	ListUsers() ([]entity.User, error)

	// Group
	CreateGroup(title string) (*entity.Group, error)
	UpdateGroup(groupID uuid.UUID, title string) (*entity.Group, error)
	ListGroups() ([]entity.Group, error)
	DeleteGroup(groupID uuid.UUID) error

	// User Group
	AddUserToGroup(userID uuid.UUID, groupID uuid.UUID) error
	RemoveUserFromGroup(userID, groupID uuid.UUID) error
	ListUserGroups(userID uuid.UUID) ([]entity.Group, error)
	ListGroupUsers(groupID uuid.UUID) ([]entity.User, error)
	// Resources

	CreateResource(name, description string, operations []uuid.UUID) (*entity.Resource, error)
	ListResources() ([]entity.Resource, error)
	UpdateResource(resourceID uuid.UUID, name, description string, operations []uuid.UUID) (*entity.Resource, error)
	DeleteResource(resourceID uuid.UUID) error
	GetResourceByID(resourceID uuid.UUID) (*entity.Resource, error)
	// Operations

	CreateOperations(name string, description string) (*entity.Operation, error)
	UpdateOperations(id uuid.UUID, name string, description string) (*entity.Operation, error)
	DeleteOperations(id uuid.UUID) error
	GetOperationsByID(id uuid.UUID) (*entity.Operation, error)
	ListOperations() ([]*entity.Operation, error)

	// RBAC permission check
	CheckUserPermission(userID uuid.UUID, resourceName string, operationName string) bool
}
