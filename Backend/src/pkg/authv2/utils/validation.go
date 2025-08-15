package utils

import (
	"regexp"
	"strings"

	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
)

// ValidatePhoneNumber validates phone number format
func ValidatePhoneNumber(prefix, number string) *entity.AuthError {
	// Validate prefix (should be 1-3 digits)
	if len(prefix) == 0 || len(prefix) > 3 {
		return entity.NewAuthError(entity.ErrInvalidPhoneFormat, "Phone prefix must be 1-3 digits")
	}

	// Check if prefix contains only digits
	prefixRegex := regexp.MustCompile(`^\d+$`)
	if !prefixRegex.MatchString(prefix) {
		return entity.NewAuthError(entity.ErrInvalidPhoneFormat, "Phone prefix must contain only digits")
	}

	// Validate number (should be 7-15 digits)
	if len(number) != 9 {
		return entity.NewAuthError(entity.ErrInvalidPhoneFormat, "Phone number must be 9 digits long")
	}

	// Check if number contains only digits
	numberRegex := regexp.MustCompile(`^\d+$`)
	if !numberRegex.MatchString(number) {
		return entity.NewAuthError(entity.ErrInvalidPhoneFormat, "Phone number must contain only digits")
	}

	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) *entity.AuthError {
	if len(password) < 8 {
		return entity.NewAuthError(entity.ErrInvalidPasswordFormat, "Password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return entity.NewAuthError(entity.ErrInvalidPasswordFormat, "Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
	}

	return nil
}

// ValidateUserType validates user type
func ValidateUserType(userType entity.UserType) *entity.AuthError {
	switch userType {
	case entity.USER_TYPE_SUPER_ADMIN, entity.USER_TYPE_ADMIN, entity.USER_TYPE_MERCHANT:
		return nil
	default:
		return entity.NewAuthError(entity.ErrUserTypeInvalid, "User type must be 'super_admin', 'admin' or 'merchant'")
	}
}

// ValidateCreateUserRequest validates create user request
func ValidateCreateUserRequest(req *entity.CreateUserRequest) *entity.AuthError {
	// Required fields
	if strings.TrimSpace(req.FirstName) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "First name is required")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "Last name is required")
	}
	if strings.TrimSpace(req.PhonePrefix) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "Phone prefix is required")
	}
	if strings.TrimSpace(req.PhoneNumber) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "Phone number is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "Password is required")
	}

	// Validate phone number
	if err := ValidatePhoneNumber(req.PhonePrefix, req.PhoneNumber); err != nil {
		return err
	}

	// Validate password
	if err := ValidatePassword(req.Password); err != nil {
		return err
	}

	// Validate user type
	if err := ValidateUserType(req.UserType); err != nil {
		return err
	}

	return nil
}

// ValidateLoginRequest validates login request
func ValidateLoginRequest(req *entity.LoginRequest) *entity.AuthError {
	if strings.TrimSpace(req.PhonePrefix) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "Phone prefix is required")
	}
	if strings.TrimSpace(req.PhoneNumber) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "Phone number is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "Password is required")
	}

	// Validate phone number format
	if err := ValidatePhoneNumber(req.PhonePrefix, req.PhoneNumber); err != nil {
		return err
	}

	return nil
}

// ValidateOTPRequest validates OTP verification request
func ValidateOTPRequest(req *entity.VerifyOTPRequest) *entity.AuthError {
	if strings.TrimSpace(req.Token) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "OTP token is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return entity.NewAuthError(entity.ErrMissingRequiredData, "OTP code is required")
	}

	// Validate OTP code format (should be 6 digits)
	if len(req.Code) != 6 {
		return entity.NewAuthError(entity.ErrOTPInvalid, "OTP code must be 6 digits")
	}

	codeRegex := regexp.MustCompile(`^\d{6}$`)
	if !codeRegex.MatchString(req.Code) {
		return entity.NewAuthError(entity.ErrOTPInvalid, "OTP code must contain only digits")
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
