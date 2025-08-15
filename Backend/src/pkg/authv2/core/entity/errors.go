package entity

import "fmt"

// AuthError represents an authentication error
type AuthError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"-"`
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Error types
const (
	ErrInvalidRequest         = "INVALID_REQUEST"
	ErrInvalidCredentials     = "INVALID_CREDENTIALS"
	ErrPhoneAlreadyExists     = "PHONE_ALREADY_EXISTS"
	ErrPhoneNotFound          = "PHONE_NOT_FOUND"
	ErrUserNotFound           = "USER_NOT_FOUND"
	ErrInvalidPhoneFormat     = "INVALID_PHONE_FORMAT"
	ErrInvalidPasswordFormat  = "INVALID_PASSWORD_FORMAT"
	ErrPasswordMismatch       = "PASSWORD_MISMATCH"
	ErrMissingRequiredData    = "MISSING_REQUIRED_DATA"
	ErrAccountCreationFailed  = "ACCOUNT_CREATION_FAILED"
	ErrInternalServer         = "INTERNAL_ERROR"
	ErrSessionExpired         = "SESSION_EXPIRED"
	ErrInvalidToken           = "INVALID_TOKEN"
	ErrPermissionDenied       = "PERMISSION_DENIED"
	ErrOTPExpired             = "OTP_EXPIRED"
	ErrOTPInvalid             = "OTP_INVALID"
	ErrOTPVerificationFailed  = "OTP_VERIFICATION_FAILED"
	ErrTooManyAttempts        = "TOO_MANY_ATTEMPTS"
	ErrDeviceNotTrusted       = "DEVICE_NOT_TRUSTED"
	ErrUserTypeInvalid        = "USER_TYPE_INVALID"
	ErrMerchantCreationFailed = "MERCHANT_CREATION_FAILED"
	ErrWalletCreationFailed   = "WALLET_CREATION_FAILED"
)

// Error messages
const (
	MsgInvalidRequest         = "Invalid request format"
	MsgInvalidCredentials     = "Invalid credentials provided"
	MsgPhoneAlreadyExists     = "This phone number is already registered"
	MsgPhoneNotFound          = "Phone number not found"
	MsgUserNotFound           = "User not found"
	MsgInvalidPhoneFormat     = "Phone number format is invalid"
	MsgInvalidPasswordFormat  = "Password must be at least 8 characters long"
	MsgPasswordMismatch       = "Password and confirmation do not match"
	MsgMissingRequiredData    = "Required information was not provided"
	MsgAccountCreationFailed  = "Failed to create account. Please try again."
	MsgInternalServer         = "An unexpected error occurred"
	MsgSessionExpired         = "Your session has expired. Please login again."
	MsgInvalidToken           = "Invalid authentication token"
	MsgPermissionDenied       = "You don't have permission to perform this action"
	MsgOTPExpired             = "OTP has expired. Please request a new one."
	MsgOTPInvalid             = "Invalid OTP code"
	MsgOTPVerificationFailed  = "OTP verification failed"
	MsgTooManyAttempts        = "Too many failed attempts. Please try again later."
	MsgDeviceNotTrusted       = "Device not trusted"
	MsgUserTypeInvalid        = "Invalid user type"
	MsgMerchantCreationFailed = "Failed to create merchant account"
	MsgWalletCreationFailed   = "Failed to create wallet"
)

// NewAuthError creates a new authentication error
func NewAuthError(errType, message string) *AuthError {
	return &AuthError{
		Type:    errType,
		Message: message,
	}
}

// NewAuthErrorWithDetail creates a new authentication error with details
func NewAuthErrorWithDetail(errType, message, detail string) *AuthError {
	return &AuthError{
		Type:    errType,
		Message: message,
		Detail:  detail,
	}
}
