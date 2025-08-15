package entity

import (
	"time"

	"github.com/google/uuid"
	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	rbac_entity "github.com/socialpay/socialpay/src/pkg/rbac/core/entity"
)

// Group represents a user group/role
type Group struct {
	ID          uuid.UUID                `json:"id" db:"id"`
	Title       string                   `json:"title" db:"title"`
	Description *string                  `json:"description" db:"description"`
	MerchantID  *uuid.UUID               `json:"merchant_id" db:"merchant_id"`
	Permissions []rbac_entity.Permission `json:"permissions,omitempty"`
	CreatedAt   time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at" db:"updated_at"`
}

// TeamMember represents a member in the system
type TeamMember struct {
	ID       uuid.UUID        `json:"id" db:"id"`
	UserData auth_entity.User `json:"user_data" db:"user_data"`
	Group    string           `json:"group" db:"group"`
}

// UserGroup represents relationship between users and groups tables
type UserGroup struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserId     uuid.UUID  `json:"user_id" db:"user_id"`
	GroupId    uuid.UUID  `json:"group_id" db:"group_id"`
	UserType   string     `json:"user_type" db:"user_type"`
	MerchantID *uuid.UUID `json:"merchant_id" db:"merchant_id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateGroupRequest represents a request to create a new group
type CreateGroupRequest struct {
	Title       string                                `json:"title" validate:"required"`
	Description string                                `json:"description" validate:"required"`
	MerchantID  *uuid.UUID                            `json:"merchant_id,omitempty"`
	Permissions []rbac_entity.CreatePermissionRequest `json:"permissions" validate:"required"`
}

// UpdateGroupRequest represents a request to update an existing group
type UpdateGroupRequest struct {
	ID          uuid.UUID                              `json:"id" validate:"required"`
	Title       *string                                `json:"title,omitempty"`
	Description *string                                `json:"description,omitempty"`
	MerchantID  *uuid.UUID                             `json:"merchant_id" validate:"required"`
	Permissions *[]rbac_entity.CreatePermissionRequest `json:"permissions,omitempty"`
}

// CreateTeamMemberRequest represents a request to create a new user
type CreateTeamMemberRequest struct {
	GroupId  uuid.UUID                     `json:"group_id" validate:"required"`
	UserData auth_entity.CreateUserRequest `json:"user_data" validate:"required"`
}

// UpdateTeamMemberRequest represents a request to update an existing team member
type UpdateTeamMemberRequest struct {
	GroupId  uuid.UUID                     `json:"group_id" validate:"required"`
	UserData auth_entity.UpdateUserRequest `json:"user_data" validate:"required"`
}
