package entity

import "fmt"

// RBACError represents an rbac error
type RBACError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"-"`
}

func (e *RBACError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Error types
const (
	ErrInvalidRequest      = "INVALID_REQUEST"
	ErrUserNotFound        = "USER_NOT_FOUND"
	ErrMissingRequiredData = "MISSING_REQUIRED_DATA"
	ErrInternalServer      = "INTERNAL_ERROR"
	ErrSessionExpired      = "SESSION_EXPIRED"
	ErrInvalidToken        = "INVALID_TOKEN"
	ErrPermissionDenied    = "PERMISSION_DENIED"
	ErrTooManyAttempts     = "TOO_MANY_ATTEMPTS"
	ErrDeviceNotTrusted    = "DEVICE_NOT_TRUSTED"
)

// Error messages
const (
	MsgInvalidRequest      = "Invalid request format"
	MsgResourceNotFound    = "Resource not found"
	MsgMissingRequiredData = "Required information was not provided"
	MsgInternalServer      = "An unexpected error occurred"
	MsgSessionExpired      = "Your session has expired. Please login again."
	MsgInvalidToken        = "Invalid authentication token"
	MsgPermissionDenied    = "You don't have permission to perform this action"
	MsgTooManyAttempts     = "Too many failed attempts. Please try again later."
	MsgDeviceNotTrusted    = "Device not trusted"
)

// NewRBACError creates a new rbac error
func NewRBACError(errType, message string) *RBACError {
	return &RBACError{
		Type:    errType,
		Message: message,
	}
}

// NewRBACErrorWithDetail creates a new rbac error with details
func NewRBACErrorWithDetail(errType, message, detail string) *RBACError {
	return &RBACError{
		Type:    errType,
		Message: message,
		Detail:  detail,
	}
}
