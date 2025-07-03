package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

// TwoFactorError represents 2FA-specific errors
type TwoFactorError struct {
	Type    string
	Message string
}

func (e TwoFactorError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// GenerateTwoFactorCode generates a 6-digit verification code
func (uc Usecase) GenerateTwoFactorCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// EnableTwoFactor initiates 2FA setup for a user
func (uc Usecase) EnableTwoFactor(userId uuid.UUID, phoneNumber string) error {
	// Check if user exists
	user, err := uc.repo.FindUserById(userId)
	if err != nil {
		return TwoFactorError{
			Type:    "USER_NOT_FOUND",
			Message: "User not found",
		}
	}
	if user == nil {
		return TwoFactorError{
			Type:    "USER_NOT_FOUND",
			Message: "User not found",
		}
	}

	// Check if 2FA is already enabled
	status, err := uc.repo.GetTwoFactorStatus(userId)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to check 2FA status",
		}
	}
	if status.Enabled {
		return TwoFactorError{
			Type:    "ALREADY_ENABLED",
			Message: "Two-factor authentication is already enabled",
		}
	}

	// Generate verification code
	code := uc.GenerateTwoFactorCode()
	expiresAt := time.Now().Add(10 * time.Minute) // 10 minutes expiry

	// Store the code
	twoFactorCode := entity.TwoFactorCode{
		Id:        uuid.New(),
		UserId:    userId,
		Code:      code,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: time.Now(),
	}

	err = uc.repo.StoreTwoFactorCode(twoFactorCode)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to store verification code",
		}
	}

	// Send SMS with verification code
	message := fmt.Sprintf("Your SocialPay 2FA verification code is %s. This code expires in 10 minutes. Do not share this code with anyone.", code)

	// Enhanced logging for SMS sending
	uc.log.Printf("🚀 [2FA] Starting SMS sending process...")
	uc.log.Printf("📱 [2FA] Phone number: %s", phoneNumber)
	uc.log.Printf("📝 [2FA] Message: %s", message)
	uc.log.Printf("🔢 [2FA] Verification code: %s", code)
	uc.log.Printf("⏰ [2FA] Expires at: %s", expiresAt.Format("2006-01-02 15:04:05"))

	go func() {
		uc.log.Printf("📤 [2FA] Sending SMS via provider...")
		if err := uc.sms.SendSMS(phoneNumber, message); err != nil {
			uc.log.Printf("❌ [2FA] Failed to send 2FA SMS to %s: %v", phoneNumber, err)
		} else {
			uc.log.Printf("✅ [2FA] SMS sent successfully to %s", phoneNumber)
			uc.log.Printf("✅ [2FA] Verification code %s delivered to %s", code, phoneNumber)
		}
	}()

	return nil
}

// VerifyTwoFactorCode verifies the 2FA setup code and enables 2FA
func (uc Usecase) VerifyTwoFactorCode(userId uuid.UUID, code string) error {
	// Find the verification code
	twoFactorCode, err := uc.repo.FindTwoFactorCode(userId, code)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to verify code",
		}
	}
	if twoFactorCode == nil {
		return TwoFactorError{
			Type:    "INVALID_CODE",
			Message: "Invalid or expired verification code",
		}
	}

	// Mark code as used
	err = uc.repo.MarkTwoFactorCodeAsUsed(twoFactorCode.Id)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to mark code as used",
		}
	}

	// Enable 2FA for the user
	err = uc.repo.EnableTwoFactor(userId)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to enable two-factor authentication",
		}
	}

	return nil
}

// VerifyTwoFactorLoginCode verifies the 2FA code during login (does not enable 2FA)
func (uc Usecase) VerifyTwoFactorLoginCode(userId uuid.UUID, code string) error {
	// Find the verification code
	twoFactorCode, err := uc.repo.FindTwoFactorCode(userId, code)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to verify code",
		}
	}
	if twoFactorCode == nil {
		return TwoFactorError{
			Type:    "INVALID_CODE",
			Message: "Invalid or expired verification code",
		}
	}

	// Check if code is expired
	if time.Now().After(twoFactorCode.ExpiresAt) {
		return TwoFactorError{
			Type:    "EXPIRED_CODE",
			Message: "Verification code has expired",
		}
	}

	// Check if code is already used
	if twoFactorCode.Used {
		return TwoFactorError{
			Type:    "USED_CODE",
			Message: "Verification code has already been used",
		}
	}

	// Mark code as used
	err = uc.repo.MarkTwoFactorCodeAsUsed(twoFactorCode.Id)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to mark code as used",
		}
	}

	return nil
}

// DisableTwoFactor disables 2FA for a user after password verification
func (uc Usecase) DisableTwoFactor(userId uuid.UUID, password string) error {
	// Verify user's password first
	passwordIdentity, err := uc.repo.FindPasswordIdentityByUser(userId)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to verify password",
		}
	}
	if passwordIdentity == nil {
		return TwoFactorError{
			Type:    "NO_PASSWORD",
			Message: "No password found for user",
		}
	}

	// Compare password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(passwordIdentity.Password), []byte(password))
	if err != nil {
		return TwoFactorError{
			Type:    "INCORRECT_PASSWORD",
			Message: "Password is incorrect",
		}
	}

	// Disable 2FA
	err = uc.repo.DisableTwoFactor(userId)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to disable two-factor authentication",
		}
	}

	return nil
}

// GetTwoFactorStatus returns the 2FA status for a user
func (uc Usecase) GetTwoFactorStatus(userId uuid.UUID) (*entity.TwoFactorStatus, error) {
	status, err := uc.repo.GetTwoFactorStatus(userId)
	if err != nil {
		return nil, TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to get 2FA status",
		}
	}
	return status, nil
}

// SendTwoFactorCode sends a new verification code for 2FA
func (uc Usecase) SendTwoFactorCode(userId uuid.UUID, phoneNumber string) error {
	// Check if 2FA is enabled
	status, err := uc.repo.GetTwoFactorStatus(userId)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to check 2FA status",
		}
	}
	if !status.Enabled {
		return TwoFactorError{
			Type:    "NOT_ENABLED",
			Message: "Two-factor authentication is not enabled",
		}
	}

	// Generate new verification code
	code := uc.GenerateTwoFactorCode()
	expiresAt := time.Now().Add(10 * time.Minute) // 10 minutes expiry

	// Store the code
	twoFactorCode := entity.TwoFactorCode{
		Id:        uuid.New(),
		UserId:    userId,
		Code:      code,
		ExpiresAt: expiresAt,
		Used:      false,
		CreatedAt: time.Now(),
	}

	err = uc.repo.StoreTwoFactorCode(twoFactorCode)
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to store verification code",
		}
	}

	// Send SMS with verification code
	message := fmt.Sprintf("Your SocialPay 2FA verification code is %s. This code expires in 10 minutes. Do not share this code with anyone.", code)

	// Enhanced logging for SMS sending
	uc.log.Printf("🚀 [2FA] Starting SMS sending process...")
	uc.log.Printf("📱 [2FA] Phone number: %s", phoneNumber)
	uc.log.Printf("📝 [2FA] Message: %s", message)
	uc.log.Printf("🔢 [2FA] Verification code: %s", code)
	uc.log.Printf("⏰ [2FA] Expires at: %s", expiresAt.Format("2006-01-02 15:04:05"))

	go func() {
		uc.log.Printf("📤 [2FA] Sending SMS via provider...")
		if err := uc.sms.SendSMS(phoneNumber, message); err != nil {
			uc.log.Printf("❌ [2FA] Failed to send 2FA SMS to %s: %v", phoneNumber, err)
		} else {
			uc.log.Printf("✅ [2FA] SMS sent successfully to %s", phoneNumber)
			uc.log.Printf("✅ [2FA] Verification code %s delivered to %s", code, phoneNumber)
		}
	}()

	return nil
}

// CleanupExpiredCodes removes expired 2FA codes from the database
func (uc Usecase) CleanupExpiredCodes() error {
	err := uc.repo.CleanupExpiredTwoFactorCodes()
	if err != nil {
		return TwoFactorError{
			Type:    "DATABASE_ERROR",
			Message: "Failed to cleanup expired codes",
		}
	}
	return nil
}
