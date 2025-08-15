package usecase

import (
	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"

	"github.com/google/uuid"
)

type AuthNRepo interface {
	FindSessionById(uuid.UUID) (*entity.Session, error)
	// PreSession
	StorePreSession(entity.PreSession) error

	// Device
	StoreDevice(entity.Device) error
	StoreDeviceAuth(deviceAuth entity.DeviceAuth) error
	UpdateDeviceAuthStatus(deviceAuthId uuid.UUID, status bool) error
	FindDeviceAuth(token string) (entity.DeviceAuth, error)

	// // Phone
	StorePhone(phone entity.Phone) error
	FindPhone(prefix string, number string) (*entity.Phone, error)
	LoginFindPhone(prefix string, number string) (*entity.Phone, error)

	StorePhoneAuth(phoneAuth entity.PhoneAuth) error
	FindPhoneAuth(token string) (entity.PhoneAuth, error)
	FindPhoneAuthWithoutPhone(token string) (entity.PhoneAuth, error)
	UpdatePhoneAuthStatus(phoneAuthId uuid.UUID, status bool) error

	// Password
	StorePasswordAuth(passwordAuth entity.PasswordAuth) error
	FindPasswordAuth(token string) (*entity.PasswordAuth, error)

	// // Session
	StoreSession(session entity.Session) error
	UpdatePasswordIdentity(password string, userId uuid.UUID) error
	// Permission
	CheckPermission(userID uuid.UUID, requiredPermission entity.Permission) (bool, error)
	FindUserPermissions(userID uuid.UUID, requiredPermission entity.Permission) ([]entity.Permission, error)
}

// Repo
type Repo interface {
	AuthNRepo

	// User
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
	StoreUser(user entity.User) error
	FindUserById(id uuid.UUID) (*entity.User, error)
	StorePhoneIdentity(phoneIdentity entity.PhoneIdentity) error
	FindUserUsingPhoneIdentity(phoneId uuid.UUID) (*entity.User, error)
	StorePasswordIdentity(passwordIdentity entity.PasswordIdentity) error
	FindPasswordIdentityByUser(userId uuid.UUID) (*entity.PasswordIdentity, error)
}
