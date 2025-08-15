package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
)

// AuthRepository defines the interface for authentication data operations
type AuthRepository interface {
	// User operations
	CreateUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetUserByPhone(ctx context.Context, prefix, number string) (*entity.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error
	RemoveUser(ctx context.Context, userID uuid.UUID) error

	// Phone operations
	CreatePhone(ctx context.Context, prefix, number string) (*entity.Phone, error)
	GetPhoneByID(ctx context.Context, phoneID uuid.UUID) (*entity.Phone, error)
	GetPhoneByNumber(ctx context.Context, prefix, number string) (*entity.User, error)
	LinkPhoneToUser(ctx context.Context, userID, phoneID uuid.UUID) error

	// Password operations
	CreatePasswordIdentity(ctx context.Context, userID uuid.UUID, hashedPassword, hint string) error
	VerifyPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error

	// OTP operations
	CreateOTP(ctx context.Context, userID uuid.UUID, code, token string, expiresAt int64) (*entity.OTPRequest, error)
	GetOTPByToken(ctx context.Context, token string) (*entity.OTPRequest, error)
	MarkOTPAsUsed(ctx context.Context, otpID uuid.UUID) error

	// Session operations
	CreateSession(ctx context.Context, userID, deviceID uuid.UUID, token, refreshToken string, expiresAt int64) (*entity.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*entity.Session, error)
	UpdateSessionToken(ctx context.Context, sessionID uuid.UUID, token, refreshToken string, expiresAt int64) error
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error

	// Device operations
	CreateDevice(ctx context.Context, ip, name, agent string) (*entity.Device, error)
	GetDeviceByFingerprint(ctx context.Context, ip, name, agent string) (*entity.Device, error)

	// RBAC operations (merchant-aware)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) // Deprecated: Use GetUserPermissionsByMerchant
	GetUserPermissionsByMerchant(ctx context.Context, userID, merchantID uuid.UUID) ([]string, error)
	GetAllUserPermissions(ctx context.Context, userID uuid.UUID) (map[string][]string, error)            // Returns permissions grouped by merchant_id
	CheckUserPermission(ctx context.Context, userID uuid.UUID, resource, operation string) (bool, error) // Deprecated: Use CheckUserPermissionForMerchant
	CheckUserPermissionForMerchant(ctx context.Context, userID, merchantID uuid.UUID, resource, operation string) (bool, error)
	GetUserGroups(ctx context.Context, userID uuid.UUID) ([]entity.Group, error)
	GetUserGroupsByMerchant(ctx context.Context, userID, merchantID uuid.UUID) ([]entity.Group, error)
	GetAllUserGroups(ctx context.Context, userID uuid.UUID) (map[string][]entity.Group, error) // Returns groups grouped by merchant_id
	AssignUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error
	UpdateUserGroup(ctx context.Context, userID, oldGroupID, newGroupID uuid.UUID) error

	// Group/Role operations
	CreateGroup(ctx context.Context, title string, description *string, merchantID *uuid.UUID) (*entity.Group, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*entity.Group, error)
	GetGroupsByMerchant(ctx context.Context, merchantID *uuid.UUID) ([]entity.Group, error)
	AssignPermissionToGroup(ctx context.Context, groupID, permissionID uuid.UUID) error

	// Permission operations
	GetPermissionByResourceOperation(ctx context.Context, resource, operation string) (*entity.Permission, error)
	CreatePermission(ctx context.Context, resourceName, operation, effect string) (*entity.Permission, error)
	GetOrCreatePermission(ctx context.Context, resourceName, operation, effect string) (*entity.Permission, error)

	// Activity logging
	LogAuthActivity(ctx context.Context, activity *entity.AuthActivity) error
	GetUserActivities(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.AuthActivity, error)

	// Merchant operations
	CreateMerchant(ctx context.Context, userID uuid.UUID, merchantData map[string]interface{}) (*uuid.UUID, error)
	GetMerchant(ctx context.Context, id uuid.UUID) (*entity.Merchant, error)
	GetMerchantsByUser(ctx context.Context, userID uuid.UUID) ([]entity.Merchant, error)
	CreateWallet(ctx context.Context, userID uuid.UUID, walletType string) error

	// Admin operations
	UserExists(ctx context.Context, phonePrefix, phoneNumber string) (bool, error)
	GetOrCreateSuperAdminGroup(ctx context.Context) (*entity.Group, error)
	GetOrCreateMerchantOwnerGroup(ctx context.Context, merchantID uuid.UUID) (*entity.Group, error)
}
