package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	// Authentication
	Register(ctx context.Context, req *entity.CreateUserRequest) (*entity.AuthResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (string, error) // Returns OTP token
	VerifyOTP(ctx context.Context, req *entity.VerifyOTPRequest) (*entity.AuthResponse, error)
	VerifyResetPasswordOTP(ctx context.Context, req *entity.VerifyOTPRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*entity.AuthResponse, error)
	Logout(ctx context.Context, token string) error
	RequestOTP(ctx context.Context, phonePrefix string, phoneNumber string) (string, error)
	UpdatePassword(ctx context.Context, req *entity.UpdatePasswordRequest) error

	// User management
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetUserByPhone(ctx context.Context, phonePrefix, phoneNumber string) (*entity.User, error)
	GetMerchantsByUser(ctx context.Context, userID uuid.UUID) ([]entity.Merchant, error)
	GetUserPermissionsByMerchant(ctx context.Context, userID, merchantID uuid.UUID) ([]string, error)
	GetAllUserPermissions(ctx context.Context, userID uuid.UUID) (map[string][]string, error)
	GetUserGroupsByGroupedByMerchant(ctx context.Context, userID uuid.UUID) (map[string][]entity.Group, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error

	// 2FA
	SendOTP(ctx context.Context, phonePrefix, phoneNumber string) (string, error) // Returns OTP token
	ValidateOTP(ctx context.Context, token, code string) (bool, error)

	// Authorization (merchant-aware)
	CheckPermission(ctx context.Context, userID uuid.UUID, resource entity.Resource, operation entity.Operation) (bool, error) // Deprecated: Use CheckPermissionForMerchant
	CheckPermissionForMerchant(ctx context.Context, userID, merchantID uuid.UUID, resource entity.Resource, operation entity.Operation) (bool, error)

	// Admin operations
	CreateSuperAdminUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.User, error)
	AssignUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error

	// Activity logging
	LogActivity(ctx context.Context, userID uuid.UUID, activityType entity.AuthActivityType, ip, userAgent, deviceName string, success bool, details map[string]interface{}) error
	GetUserActivities(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.AuthActivity, error)

	// Device
	CreateDevice(ctx context.Context, args entity.CreateDeviceArgs) (*entity.Device, error)

	// Session
	CreateSession(ctx context.Context, args entity.CreateSessionArgs) (*entity.Session, error)
}

// PermissionChecker provides a simple interface for permission checking
type PermissionChecker interface {
	RequirePermission(resource entity.Resource, operation entity.Operation) func(userID uuid.UUID) bool
	RequirePermissionForMerchant(resource entity.Resource, operation entity.Operation) func(userID, merchantID uuid.UUID) bool
}
