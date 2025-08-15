package entity

import "fmt"

// TeamManagementError represents an authentication error
type TeamManagementError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"-"`
}

func (e *TeamManagementError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Error types
const (
	ErrInvalidRequest      = "INVALID_REQUEST"
	ErrGroupAlreadyExists  = "GROUP_ALREADY_EXISTS"
	ErrGroupNotFound       = "GROUP_NOT_FOUND"
	ErrUserNotFound        = "USER_NOT_FOUND"
	ErrMissingRequiredData = "MISSING_REQUIRED_DATA"
	ErrGroupCreationFailed = "GROUP_CREATION_FAILED"
	ErrGroupDeletionFailed = "GROUP_DELETION_FAILED"
	ErrGroupUpdateFailed   = "GROUP_UPDATE_FAILED"
	ErrInternalServer      = "INTERNAL_ERROR"
)

// Error messages
const (
	MsgInvalidRequest      = "Invalid request format"
	MsgGroupAlreadyExists  = "This group is already created"
	MsgGroupNotFound       = "Group not found"
	MsgUserNotFound        = "User not found"
	MsgMissingRequiredData = "Required information was not provided"
	MsgGroupCreationFailed = "Failed to create group. Please try again."
	MsgGroupDeletionFailed = "Failed to delete group. Please try again."
	MsgGroupUpdateFailed   = "Failed to update group. Please try again."
	MsgInternalServer      = "An unexpected error occurred"
)

// NewTeamManagementError creates a new team management error
func NewTeamManagementError(errType, message string) *TeamManagementError {
	return &TeamManagementError{
		Type:    errType,
		Message: message,
	}
}

// NewTeamManagementErrorWithDetail creates a new team management error with details
func NewTeamManagementErrorWithDetail(errType, message, detail string) *TeamManagementError {
	return &TeamManagementError{
		Type:    errType,
		Message: message,
		Detail:  detail,
	}
}
