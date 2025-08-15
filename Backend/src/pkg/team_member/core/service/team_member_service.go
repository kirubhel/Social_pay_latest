package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/entity"
)

// TeamMemberService defines the interface for team member management business logic
type TeamMemberService interface {
	//Group
	CreateGroup(ctx context.Context, req *entity.CreateGroupRequest) error
	GetGroups(ctx context.Context, merchantId *uuid.UUID) ([]entity.Group, error)
	GetGroupById(ctx context.Context, merchantId *uuid.UUID, id uuid.UUID) (*entity.Group, error)
	UpdateGroup(ctx context.Context, req *entity.UpdateGroupRequest) error
	DeleteGroupById(ctx context.Context, merchantId *uuid.UUID, id uuid.UUID) error

	//Team member operations
	CreateTeamMember(ctx context.Context, merchantId *uuid.UUID, req *entity.CreateTeamMemberRequest) error
	UpdateTeamMember(ctx context.Context, merchantId *uuid.UUID, userGroupId uuid.UUID, req *entity.UpdateTeamMemberRequest) error
	GetTeamMembers(ctx context.Context, merchantId *uuid.UUID, groupId *uuid.UUID) ([]entity.TeamMember, error)
	RemoveTeamMember(ctx context.Context, merchantId *uuid.UUID, userGroupId uuid.UUID) error
}
