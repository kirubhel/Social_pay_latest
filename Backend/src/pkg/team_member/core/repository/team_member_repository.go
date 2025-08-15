package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/entity"
)

// TeamMemberRepository defines the interface for team member data operations
type TeamMemberRepository interface {
	// Group operations
	CreateGroup(ctx context.Context, req *entity.CreateGroupRequest) error
	GroupExists(ctx context.Context, title string, merchantId *uuid.UUID) (bool, error)
	GetGroups(ctx context.Context, merchantId *uuid.UUID) ([]entity.Group, error)
	GetGroupById(ctx context.Context, merchantId *uuid.UUID, id uuid.UUID) (*entity.Group, error)
	UpdateGroup(ctx context.Context, req *entity.UpdateGroupRequest) error
	DeleteGroupById(ctx context.Context, id uuid.UUID) error
	DeleteGroupPermissionById(ctx context.Context, id uuid.UUID) error

	// User Group operations
	GetUserGroups(ctx context.Context, merchantId *uuid.UUID, groupId uuid.UUID) ([]entity.UserGroup, error)
	GetUserGroupById(ctx context.Context, merchantId *uuid.UUID, userGroupId uuid.UUID) (*entity.UserGroup, error)
}
