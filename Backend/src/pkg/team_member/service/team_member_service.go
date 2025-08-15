package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	auth_repository "github.com/socialpay/socialpay/src/pkg/authv2/core/repository"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/entity"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/repository"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/service"
	"github.com/socialpay/socialpay/src/pkg/team_member/utils"
)

// TeamMemberServiceImpl implements the TeamMemberService interface
type TeamMemberServiceImpl struct {
	authRepository auth_repository.AuthRepository
	repo           repository.TeamMemberRepository
	jwtSecret      string
	logger         *log.Logger
}

// NewTeamMemberService creates a new team member service
func NewTeamMemberService(authRepository auth_repository.AuthRepository, repo repository.TeamMemberRepository, jwtSecret string, logger *log.Logger) service.TeamMemberService {
	return &TeamMemberServiceImpl{
		authRepository: authRepository,
		repo:           repo,
		jwtSecret:      jwtSecret,
		logger:         logger,
	}
}

// CreateGroup creates new group
func (s *TeamMemberServiceImpl) CreateGroup(ctx context.Context, req *entity.CreateGroupRequest) error {
	// Validate request
	if err := utils.ValidateCreateGroupRequest(req); err != nil {
		return err
	}

	// Sanitize inputs
	req.Title = utils.SanitizeInput(req.Title)
	req.Description = utils.SanitizeInput(req.Description)

	if req.MerchantID != nil {
		// Check if group already exists
		exists, err := s.repo.GroupExists(ctx, req.Title, req.MerchantID)
		if err != nil {
			s.logger.Printf("Error checking group existence: %v", err)
			return entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
		}
		if exists {
			return entity.NewTeamManagementError(entity.ErrGroupAlreadyExists, entity.MsgGroupAlreadyExists)
		}
	}

	// Create group
	err := s.repo.CreateGroup(ctx, req)
	if err != nil {
		s.logger.Printf("Error creating group: %v", err)
		return entity.NewTeamManagementError(entity.ErrGroupCreationFailed, entity.MsgGroupCreationFailed)
	}

	return nil
}

// GetGroups fetches list of groups
func (s *TeamMemberServiceImpl) GetGroups(ctx context.Context, merchantId *uuid.UUID) ([]entity.Group, error) {
	return s.repo.GetGroups(ctx, merchantId)
}

// GetGroupById fetches group by specific id and merchant id
func (s *TeamMemberServiceImpl) GetGroupById(ctx context.Context, merchantId *uuid.UUID, id uuid.UUID) (*entity.Group, error) {
	return s.repo.GetGroupById(ctx, merchantId, id)
}

// UpdateGroup updates an existing group efficiently
func (s *TeamMemberServiceImpl) UpdateGroup(ctx context.Context, req *entity.UpdateGroupRequest) error {
	// Validate request
	if err := utils.ValidateUpdateGroupRequest(req); err != nil {
		return err
	}

	// Sanitize inputs if provided
	if req.Title != nil {
		sanitized := utils.SanitizeInput(*req.Title)
		req.Title = &sanitized
	}
	if req.Description != nil {
		sanitized := utils.SanitizeInput(*req.Description)
		req.Description = &sanitized
	}

	// Get current group to verify ownership and check for changes
	currentGroup, err := s.repo.GetGroupById(ctx, req.MerchantID, req.ID)
	if err != nil {
		s.logger.Printf("Error getting current group: %v", err)
		return entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
	}
	if currentGroup == nil {
		return entity.NewTeamManagementError(entity.ErrGroupNotFound, entity.MsgGroupNotFound)
	}

	// Check if title is changing and if new title already exists
	if req.Title != nil && *req.Title != currentGroup.Title {
		exists, err := s.repo.GroupExists(ctx, *req.Title, req.MerchantID)
		if err != nil {
			s.logger.Printf("Error checking group existence: %v", err)
			return entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
		}
		if exists {
			return entity.NewTeamManagementError(entity.ErrGroupAlreadyExists, entity.MsgGroupAlreadyExists)
		}
	}

	// Update group
	err = s.repo.UpdateGroup(ctx, req)
	if err != nil {
		s.logger.Printf("Error updating group: %v", err)
		return entity.NewTeamManagementError(entity.ErrGroupUpdateFailed, entity.MsgGroupUpdateFailed)
	}

	return nil
}

// DeleteGroupById deletes specific group by id and merchant id
func (s *TeamMemberServiceImpl) DeleteGroupById(ctx context.Context, merchantId *uuid.UUID, id uuid.UUID) error {
	// Get group
	group, err := s.repo.GetGroupById(ctx, merchantId, id)
	if err != nil {
		s.logger.Printf("Error getting group: %v", err)
		return entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// If group doesn't exisit with the given id and merchat id throw err
	if group == nil {
		return entity.NewTeamManagementError(entity.ErrGroupNotFound, entity.MsgGroupNotFound)
	}

	// Delete group
	err = s.repo.DeleteGroupById(ctx, id)
	if err != nil {
		s.logger.Printf("Error deleting group: %v", err)
		return entity.NewTeamManagementError(entity.ErrGroupDeletionFailed, entity.MsgGroupDeletionFailed)
	}

	// Delete permissions in the group
	for _, permission := range group.Permissions {
		err = s.repo.DeleteGroupPermissionById(ctx, permission.ID)

		if err != nil {
			s.logger.Printf("Error deleting group permission: %v", err)
			return entity.NewTeamManagementError(entity.ErrGroupDeletionFailed, entity.MsgGroupDeletionFailed)
		}
	}

	return nil
}

// CreateTeamMember create new user with member account and assign it to specific group
func (s *TeamMemberServiceImpl) CreateTeamMember(ctx context.Context, merchantId *uuid.UUID, req *entity.CreateTeamMemberRequest) error {
	// Get group
	group, err := s.repo.GetGroupById(ctx, merchantId, req.GroupId)
	if err != nil {
		s.logger.Printf("Error getting group: %v", err)
		return entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// If group doesn't exisist with the given id and merchat id throw err
	if group == nil {
		return entity.NewTeamManagementError(entity.ErrGroupNotFound, entity.MsgGroupNotFound)
	}

	// create new user
	user, err := s.authRepository.CreateUser(ctx, &req.UserData)
	if err != nil {
		s.logger.Panicf("Error creating new user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	// assign user to the group
	err = s.authRepository.AssignUserToGroup(ctx, user.ID, req.GroupId)
	if err != nil {
		s.logger.Panicf("Error assigning group to user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetTeamMembers fetches team members in a group
func (s *TeamMemberServiceImpl) GetTeamMembers(ctx context.Context, merchantId *uuid.UUID, groupId *uuid.UUID) ([]entity.TeamMember, error) {
	var userGroups []entity.UserGroup

	if groupId != nil {
		// Get group
		group, err := s.repo.GetGroupById(ctx, merchantId, *groupId)
		if err != nil {
			s.logger.Printf("Error getting group: %v", err)
			return nil, entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
		}

		// If group doesn't exisist with the given id and merchat id throw err
		if group == nil {
			return nil, entity.NewTeamManagementError(entity.ErrGroupNotFound, entity.MsgGroupNotFound)
		}

		// Get user groups
		userGroups, err = s.repo.GetUserGroups(ctx, merchantId, *groupId)
		if err != nil {
			s.logger.Printf("Error getting user groups: %v", err)
			return nil, entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
		}
	} else {
		// Get groups
		groups, err := s.repo.GetGroups(ctx, merchantId)
		if err != nil {
			s.logger.Printf("Error getting group: %v", err)
			return nil, entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
		}

		for _, group := range groups {
			ug, err := s.repo.GetUserGroups(ctx, merchantId, group.ID)
			if err != nil {
				s.logger.Printf("Error getting user groups: %v", err)
				return nil, entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
			}

			userGroups = append(userGroups, ug...)
		}
	}

	var members []entity.TeamMember

	// Get members data
	for _, userGroup := range userGroups {
		user, err := s.authRepository.GetUserByID(ctx, userGroup.UserId)
		if err != nil {
			s.logger.Printf("Error getting user groups: %v", err)
			continue
		}

		group, err := s.repo.GetGroupById(ctx, merchantId, userGroup.GroupId)
		if err != nil {
			s.logger.Printf("Error getting group: %v", err)
			return nil, entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
		}

		members = append(members, entity.TeamMember{
			ID:       userGroup.ID,
			UserData: *user,
			Group:    group.Title,
		})
	}

	return members, nil
}

// UpdateTeamMember updates an existing team member's user data and group assignment
func (s *TeamMemberServiceImpl) UpdateTeamMember(ctx context.Context, merchantId *uuid.UUID, userGroupId uuid.UUID, req *entity.UpdateTeamMemberRequest) error {
	// Get the specific user group record by ID
	currentUserGroup, err := s.repo.GetUserGroupById(ctx, merchantId, userGroupId)
	if err != nil {
		s.logger.Printf("Error getting user group by id: %v", err)
		return entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	if currentUserGroup == nil {
		return entity.NewTeamManagementError(entity.ErrUserNotFound, "User group not found")
	}

	// Use the user ID from the user group record
	userID := currentUserGroup.UserId

	// Validate the new group exists and belongs to the merchant
	newGroup, err := s.repo.GetGroupById(ctx, merchantId, req.GroupId)
	if err != nil {
		s.logger.Printf("Error getting new group: %v", err)
		return entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	if newGroup == nil {
		return entity.NewTeamManagementError(entity.ErrGroupNotFound, entity.MsgGroupNotFound)
	}

	// Prepare updates map for user data (excluding user_type and device_info)
	updates := make(map[string]interface{})

	if req.UserData.Title != nil {
		updates["title"] = *req.UserData.Title
	}
	if req.UserData.FirstName != nil {
		updates["first_name"] = *req.UserData.FirstName
	}
	if req.UserData.LastName != nil {
		updates["last_name"] = *req.UserData.LastName
	}
	if req.UserData.Email != nil {
		updates["email"] = *req.UserData.Email
	}
	if req.UserData.PhonePrefix != nil {
		updates["phone_prefix"] = *req.UserData.PhonePrefix
	}
	if req.UserData.PhoneNumber != nil {
		updates["phone_number"] = *req.UserData.PhoneNumber
	}
	if req.UserData.Password != nil {
		updates["password"] = *req.UserData.Password
	}
	if req.UserData.PasswordHint != nil {
		updates["password_hint"] = *req.UserData.PasswordHint
	}

	// Update user data if there are changes
	if len(updates) > 0 {
		err = s.authRepository.UpdateUser(ctx, userID, updates)
		if err != nil {
			s.logger.Printf("Error updating user: %v", err)
			return entity.NewTeamManagementError(entity.ErrInternalServer, "Failed to update user data")
		}
	}

	// Update user group assignment if group has changed
	if currentUserGroup.GroupId != req.GroupId {
		err = s.authRepository.UpdateUserGroup(ctx, userID, currentUserGroup.GroupId, req.GroupId)
		if err != nil {
			s.logger.Printf("Error updating user group: %v", err)
			return entity.NewTeamManagementError(entity.ErrInternalServer, "Failed to update user group")
		}
	}

	return nil
}

// RemoveTeamMember removes team member from a group
func (s *TeamMemberServiceImpl) RemoveTeamMember(ctx context.Context, merchantId *uuid.UUID, userGroupId uuid.UUID) error {
	// Get the specific user group record by ID
	currentUserGroup, err := s.repo.GetUserGroupById(ctx, merchantId, userGroupId)
	if err != nil {
		s.logger.Printf("Error getting user group by id: %v", err)
		return entity.NewTeamManagementError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	if currentUserGroup == nil {
		return entity.NewTeamManagementError(entity.ErrUserNotFound, "User group not found")
	}

	// remove user from the group
	err = s.authRepository.RemoveUserFromGroup(ctx, currentUserGroup.UserId, currentUserGroup.GroupId)

	if err != nil {
		s.logger.Printf("Error removing team member: %v", err)
		return entity.NewTeamManagementError(entity.ErrInternalServer, "Failed to remove user from a group")
	}

	// remove user
	err = s.authRepository.RemoveUser(ctx, currentUserGroup.UserId)
	if err != nil {
		s.logger.Printf("Error removing user data: %v", err)
		return entity.NewTeamManagementError(entity.ErrInternalServer, "Failed to remove user")
	}

	return nil
}
