package utils

import (
	"strings"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/entity"
)

// ValidateCreateGroupRequest validates create group request
func ValidateCreateGroupRequest(req *entity.CreateGroupRequest) *entity.TeamManagementError {
	// Required fields
	if strings.TrimSpace(req.Title) == "" {
		return entity.NewTeamManagementError(entity.ErrMissingRequiredData, "Title is required")
	}
	if strings.TrimSpace(req.Description) == "" {
		return entity.NewTeamManagementError(entity.ErrMissingRequiredData, "Description is required")
	}
	if len(req.Permissions) == 0 {
		return entity.NewTeamManagementError(entity.ErrMissingRequiredData, "Permissions are empty")
	}

	return nil
}

// ValidateUpdateGroupRequest validates update group request
func ValidateUpdateGroupRequest(req *entity.UpdateGroupRequest) *entity.TeamManagementError {
	// Required fields
	if req.ID == uuid.Nil {
		return entity.NewTeamManagementError(entity.ErrMissingRequiredData, "Group ID is required")
	}

	// Validate optional fields if provided
	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		return entity.NewTeamManagementError(entity.ErrInvalidRequest, "Title cannot be empty")
	}
	if req.Description != nil && strings.TrimSpace(*req.Description) == "" {
		return entity.NewTeamManagementError(entity.ErrInvalidRequest, "Description cannot be empty")
	}
	if req.Permissions != nil && len(*req.Permissions) == 0 {
		return entity.NewTeamManagementError(entity.ErrInvalidRequest, "Permissions cannot be empty if provided")
	}

	// At least one field must be provided for update
	if req.Title == nil && req.Description == nil && req.Permissions == nil {
		return entity.NewTeamManagementError(entity.ErrInvalidRequest, "At least one field must be provided for update")
	}

	return nil
}

// SanitizeInput sanitizes string input
func SanitizeInput(input string) string {
	// Trim whitespace and remove any potentially harmful characters
	cleaned := strings.TrimSpace(input)

	// Remove null bytes and other control characters
	cleaned = strings.ReplaceAll(cleaned, "\x00", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "\r", "")
	cleaned = strings.ReplaceAll(cleaned, "\t", "")

	return cleaned
}
