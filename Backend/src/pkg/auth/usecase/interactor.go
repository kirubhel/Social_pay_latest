package usecase

import (
	"net"

	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"

	"github.com/google/uuid"
)

/*

 */

type AuthNInteractor interface {
	// Responsible for identification // session
	// PreSession
	InitPreSession() (entity.PreSession, error)
	CheckPreSession(token string) error

	// / Auth Steps
	// Device Auth
	AuthDevice(token string, ip net.IPAddr, name string, agent string) error
	CheckDeviceAuth(token string) error

	// Phone Auth
	LoginFindPhone(prefix string, number string) (*entity.Phone, error)
	InitPhoneAuth(token, prefix, number string) (*entity.PhoneAuth, error)
	AuthPhone(token, prefix, number, otp string) error
	CheckPhoneAuth(token string) error

	// 2FA Auth
	CheckPasswordAuth(userId uuid.UUID, token string) error

	// Check PreSession Auth
	InitPasswordAuth(token string, password string, hint string) (*entity.PasswordAuth, error)
	AuthPassword(token string, password string) error

	// Session
	CreateSession(token string) (*entity.Session, string, error)
	CreateSessionWithoutPhoneVerification(token string, userId uuid.UUID) (*entity.Session, string, error)
	CheckSession(token string) (*entity.Session, error)
	//Permission
	CheckPermission(userID uuid.UUID, requiredPermission entity.Permission) (bool, error)
}

type TwoFactorInteractor interface {
	// Two-Factor Authentication
	GetTwoFactorStatus(userId uuid.UUID) (*entity.TwoFactorStatus, error)
	EnableTwoFactor(userId uuid.UUID, phoneNumber string) error
	VerifyTwoFactorCode(userId uuid.UUID, code string) error
	VerifyTwoFactorLoginCode(userId uuid.UUID, code string) error
	DisableTwoFactor(userId uuid.UUID, password string) error
	SendTwoFactorCode(userId uuid.UUID, phoneNumber string) error
	CleanupExpiredCodes() error
}
type AuthZInteractor interface {
	// Responsible for access management
}
type MgmtInteractor interface {
	// Resopnsible for user account management
	/// User
	// Create User
	CreateUser(
		Title string,
		FirstName string,
		LastName string,
		PhonePrefix string,
		PhoneNumber string,
		Password string,
		PasswordHint string,
		UserType string,
	) (*entity.User, error)
	// Find User
	GetUserByPhoneNumber(phoneId uuid.UUID) (entity.User, error)
	GetUserById(id uuid.UUID) (*entity.User, error)
	GetUserWithPhoneById(id uuid.UUID) (*entity.User, error)
	CreatePasswordIdentity(userId uuid.UUID, password string, hint string) (*entity.PasswordIdentity, error)
	// GetPasswordIdentityHint(userId uuid.UUID) (string, error)
}

// Interactor
type Interactor interface {
	AuthNInteractor
	MgmtInteractor
	TwoFactorInteractor
}
