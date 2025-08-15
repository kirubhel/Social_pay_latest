package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/repository"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	"github.com/socialpay/socialpay/src/pkg/authv2/utils"
	"github.com/socialpay/socialpay/src/pkg/notifications"
	"github.com/socialpay/socialpay/src/pkg/notifications/usecase"
)

// AuthServiceImpl implements the AuthService interface
type AuthServiceImpl struct {
	repo                repository.AuthRepository
	notificationService usecase.NotificationService
	jwtSecret           string
	logger              *log.Logger
}

// NewAuthService creates a new authentication service
func NewAuthService(repo repository.AuthRepository, jwtSecret string, logger *log.Logger) service.AuthService {
	return &AuthServiceImpl{
		repo:                repo,
		notificationService: notifications.NewNotificationService(logger),
		jwtSecret:           jwtSecret,
		logger:              logger,
	}
}

// Register creates a new user account
func (s *AuthServiceImpl) Register(ctx context.Context, req *entity.CreateUserRequest) (*entity.AuthResponse, error) {
	// Validate request
	if err := utils.ValidateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// Sanitize inputs
	req.FirstName = utils.SanitizeInput(req.FirstName)
	req.LastName = utils.SanitizeInput(req.LastName)
	req.Title = utils.SanitizeInput(req.Title)
	req.PhonePrefix = utils.SanitizeInput(req.PhonePrefix)
	req.PhoneNumber = utils.SanitizeInput(req.PhoneNumber)
	req.PasswordHint = utils.SanitizeInput(req.PasswordHint)

	// Check if user already exists
	exists, err := s.repo.UserExists(ctx, req.PhonePrefix, req.PhoneNumber)
	if err != nil {
		s.logger.Printf("Error checking user existence: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}
	if exists {
		return nil, entity.NewAuthError(entity.ErrPhoneAlreadyExists, entity.MsgPhoneAlreadyExists)
	}

	// Create user
	user, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		s.logger.Printf("Error creating user: %v", err)
		return nil, entity.NewAuthError(entity.ErrAccountCreationFailed, entity.MsgAccountCreationFailed)
	}

	// If merchant, create merchant record and wallet
	if req.UserType == entity.USER_TYPE_MERCHANT {
		merchantData := map[string]interface{}{
			"legal_name":    req.FirstName + " " + req.LastName,
			"trading_name":  req.FirstName + " " + req.LastName,
			"business_type": "individual",
			"title":         req.Title,
			"first_name":    req.FirstName,
			"last_name":     req.LastName,
		}

		merchantID, err := s.repo.CreateMerchant(ctx, user.ID, merchantData)
		if err != nil {
			s.logger.Printf("Error creating merchant: %v", err)
			return nil, entity.NewAuthError(entity.ErrMerchantCreationFailed, entity.MsgMerchantCreationFailed)
		}
		user.MerchantID = merchantID

		// Create wallet
		err = s.repo.CreateWallet(ctx, user.ID, "merchant")
		if err != nil {
			s.logger.Printf("Error creating wallet: %v", err)
			return nil, entity.NewAuthError(entity.ErrWalletCreationFailed, entity.MsgWalletCreationFailed)
		}

		// Create merchant owner group and assign ALL permissions
		merchantOwnerGroup, err := s.repo.GetOrCreateMerchantOwnerGroup(ctx, *merchantID)
		if err != nil {
			s.logger.Printf("Error creating merchant owner group: %v", err)
			return nil, entity.NewAuthError(entity.ErrInternalServer, "Failed to assign merchant permissions")
		}

		// Assign user to merchant owner group (gives ALL permissions)
		err = s.repo.AssignUserToGroup(ctx, user.ID, merchantOwnerGroup.ID)
		if err != nil {
			s.logger.Printf("Error assigning merchant to owner group: %v", err)
			return nil, entity.NewAuthError(entity.ErrInternalServer, "Failed to assign merchant permissions")
		}

		s.logger.Printf("Merchant user %s assigned to owner group %s with ALL permissions", user.ID, merchantOwnerGroup.ID)
	}

	// Create session and generate tokens
	device, err := s.repo.CreateDevice(ctx, "unknown", "registration", "registration")
	if err != nil {
		s.logger.Printf("Error creating device: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Generate tokens
	sessionID := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	merchantIDStr := ""
	if user.MerchantID != nil {
		merchantIDStr = user.MerchantID.String()
	}

	token, err := utils.GenerateJWT(user.ID.String(), string(req.UserType), merchantIDStr, sessionID.String(), s.jwtSecret, 24)
	if err != nil {
		s.logger.Printf("Error generating JWT: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		s.logger.Printf("Error generating refresh token: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Create session
	_, err = s.repo.CreateSession(ctx, user.ID, device.ID, token, refreshToken, expiresAt)
	if err != nil {
		s.logger.Printf("Error creating session: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Log activity
	s.LogActivity(ctx, user.ID, entity.ACTIVITY_ACCOUNT_CREATED, "unknown", "registration", "registration", true, map[string]interface{}{
		"user_type": req.UserType,
	})

	// Get user's merchants for response
	merchants, err := s.repo.GetMerchantsByUser(ctx, user.ID)
	if err != nil {
		s.logger.Printf("Error getting user merchants for registration response: %v", err)
		// Don't fail the registration, just set empty merchants
		merchants = []entity.Merchant{}
	}

	return &entity.AuthResponse{
		User:         user,
		Merchants:    merchants,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Login authenticates a user and sends OTP
func (s *AuthServiceImpl) Login(ctx context.Context, req *entity.LoginRequest) (string, error) {
	// Validate request
	if err := utils.ValidateLoginRequest(req); err != nil {
		return "", err
	}

	// Get user by phone
	user, err := s.repo.GetUserByPhone(ctx, req.PhonePrefix, req.PhoneNumber)
	if err != nil {
		s.logger.Printf("Error getting user by phone: %v", err)
		s.LogActivity(ctx, uuid.Nil, entity.ACTIVITY_LOGIN_FAILED, "unknown", "unknown", "unknown", false, map[string]interface{}{
			"reason": "user_not_found",
			"phone":  "+" + req.PhonePrefix + req.PhoneNumber,
		})
		return "", entity.NewAuthError(entity.ErrInvalidCredentials, entity.MsgInvalidCredentials)
	}

	// Verify password
	valid, err := s.repo.VerifyPassword(ctx, user.ID, req.Password)
	if err != nil || !valid {
		s.logger.Printf("Password verification failed for user %s", user.ID)
		s.LogActivity(ctx, user.ID, entity.ACTIVITY_LOGIN_FAILED, "unknown", "unknown", "unknown", false, map[string]interface{}{
			"reason": "invalid_password",
		})
		return "", entity.NewAuthError(entity.ErrInvalidCredentials, entity.MsgInvalidCredentials)
	}

	// Send OTP
	otpToken, err := s.SendOTP(ctx, req.PhonePrefix, req.PhoneNumber)
	if err != nil {
		s.logger.Printf("Error sending OTP: %v", err)
		return "", err
	}

	s.LogActivity(ctx, user.ID, entity.ACTIVITY_OTP_SENT, "unknown", "unknown", "unknown", true, map[string]interface{}{
		"otp_token": otpToken,
	})

	return otpToken, nil
}

// SendOTP sends an OTP code to the specified phone number
func (s *AuthServiceImpl) SendOTP(ctx context.Context, phonePrefix, phoneNumber string) (string, error) {
	// Get phone record
	user, err := s.repo.GetPhoneByNumber(ctx, phonePrefix, phoneNumber)
	if err != nil {
		return "", entity.NewAuthError(entity.ErrPhoneNotFound, entity.MsgPhoneNotFound)
	}

	// Generate OTP code and token
	code, err := utils.GenerateOTPCode()
	if err != nil {
		s.logger.Printf("Error generating OTP code: %v", err)
		return "", entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	token, err := utils.GenerateOTPToken()
	if err != nil {
		s.logger.Printf("Error generating OTP token: %v", err)
		return "", entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Store OTP (expires in 5 minutes)
	expiresAt := time.Now().Add(5 * time.Minute).Unix()
	_, err = s.repo.CreateOTP(ctx, user.ID, code, token, expiresAt)
	if err != nil {
		s.logger.Printf("Error creating OTP: %v", err)
		return "", entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Send SMS using notification service
	phoneNumberFull := "+" + phonePrefix + phoneNumber
	message := fmt.Sprintf("Your SocialPay verification code is: %s. This code will expire in 5 minutes. Do not share this code with anyone.", code)
	err = s.notificationService.SendSMS(ctx, phoneNumberFull, message)
	if err != nil {
		s.logger.Printf("Error sending SMS: %v", err)
		return "", entity.NewAuthError(entity.ErrInternalServer, "Failed to send OTP via SMS")
	}

	return token, nil
}

func (s AuthServiceImpl) verifyOTP(ctx context.Context, req *entity.VerifyOTPRequest) (*entity.User, error) {
	// Validate request
	if err := utils.ValidateOTPRequest(req); err != nil {
		return nil, err
	}

	// Get OTP record
	otp, err := s.repo.GetOTPByToken(ctx, req.Token)
	if err != nil {
		return nil, entity.NewAuthError(entity.ErrOTPInvalid, entity.MsgOTPInvalid)
	}

	// Check if OTP is expired
	if time.Now().Unix() > otp.ExpiresAt.Unix() {
		return nil, entity.NewAuthError(entity.ErrOTPExpired, entity.MsgOTPExpired)
	}

	// Check if OTP is already used
	if otp.Status {
		return nil, entity.NewAuthError(entity.ErrOTPInvalid, entity.MsgOTPInvalid)
	}

	// Verify code
	if otp.Code != req.Code {
		return nil, entity.NewAuthError(entity.ErrOTPInvalid, entity.MsgOTPInvalid)
	}

	// Mark OTP as used
	err = s.repo.MarkOTPAsUsed(ctx, otp.ID)
	if err != nil {
		s.logger.Printf("Error marking OTP as used: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Get user directly from OTP record
	user, err := s.repo.GetUserByID(ctx, otp.UserID)
	if err != nil {
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	return user, nil
}

// Sends OTP
func (s *AuthServiceImpl) RequestOTP(ctx context.Context, phonePrefix string, phoneNumber string) (string, error) {
	// Get user by phone
	user, err := s.repo.GetUserByPhone(ctx, phonePrefix, phoneNumber)
	if err != nil {
		s.logger.Printf("Error getting user by phone: %v", err)
		s.LogActivity(ctx, uuid.Nil, entity.ACTIVITY_LOGIN_FAILED, "unknown", "unknown", "unknown", false, map[string]interface{}{
			"reason": "user_not_found",
			"phone":  "+" + phonePrefix + phoneNumber,
		})
		return "", entity.NewAuthError(entity.ErrInvalidCredentials, entity.MsgInvalidCredentials)
	}

	// Send OTP
	otpToken, err := s.SendOTP(ctx, phonePrefix, phoneNumber)
	if err != nil {
		s.logger.Printf("Error sending OTP: %v", err)
		return "", err
	}

	s.LogActivity(ctx, user.ID, entity.ACTIVITY_OTP_SENT, "unknown", "unknown", "unknown", true, map[string]interface{}{
		"otp_token": otpToken,
	})

	return otpToken, nil
}

// VerifyOTP verifies an OTP code and completes authentication
func (s *AuthServiceImpl) VerifyOTP(ctx context.Context, req *entity.VerifyOTPRequest) (*entity.AuthResponse, error) {
	user, err := s.verifyOTP(ctx, req)
	if err != nil {
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Create device and session
	device, err := s.repo.CreateDevice(ctx, "unknown", "login", "login")
	if err != nil {
		s.logger.Printf("Error creating device: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Generate tokens
	sessionID := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	merchantIDStr := ""
	if user.MerchantID != nil {
		merchantIDStr = user.MerchantID.String()
	}

	token, err := utils.GenerateJWT(user.ID.String(), string(user.UserType), merchantIDStr, sessionID.String(), s.jwtSecret, 24)
	if err != nil {
		s.logger.Printf("Error generating JWT: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		s.logger.Printf("Error generating refresh token: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Create session
	_, err = s.repo.CreateSession(ctx, user.ID, device.ID, token, refreshToken, expiresAt)
	if err != nil {
		s.logger.Printf("Error creating session: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Log successful login
	s.LogActivity(ctx, user.ID, entity.ACTIVITY_LOGIN_SUCCESS, "unknown", "unknown", "unknown", true, map[string]interface{}{
		"otp_token": req.Token,
	})

	var merchants []entity.Merchant

	fmt.Println("user type -> ", user.UserType)
	if user.UserType == entity.USER_TYPE_MERCHANT {
		// Get user's merchants for response
		merchants, err = s.repo.GetMerchantsByUser(ctx, user.ID)
		if err != nil {
			s.logger.Printf("Error getting user merchants for registration response: %v", err)
			// Don't fail the registration, just set empty merchants
			merchants = []entity.Merchant{}
		}
	} else if user.UserType == entity.USER_TYPE_MEMBER {
		groups, err := s.repo.GetAllUserGroups(ctx, user.ID)
		if err != nil {
			s.logger.Printf("Errof getting user groups: %v", err)
			merchants = []entity.Merchant{}
		}

		fmt.Println("user id %w and groups %w", user.ID, groups)
		for merchantID := range groups {
			merchantUUID, err := uuid.Parse(merchantID)
			if err != nil {
				s.logger.Printf("Error parsing merchant ID: %w", err)
				continue
			}

			merchant, err := s.repo.GetMerchant(ctx, merchantUUID)
			if err != nil {
				s.logger.Printf("Error getting merchant by id: %v", err)
				continue
			}

			merchants = append(merchants, *merchant)
		}
	}

	return &entity.AuthResponse{
		User:         user,
		Merchants:    merchants,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// ValidateOTP validates an OTP without completing authentication
func (s *AuthServiceImpl) ValidateOTP(ctx context.Context, token, code string) (bool, error) {
	otp, err := s.repo.GetOTPByToken(ctx, token)
	if err != nil {
		return false, nil
	}

	// Check expiry and code
	if time.Now().Unix() > otp.ExpiresAt.Unix() || otp.Status || otp.Code != code {
		return false, nil
	}

	return true, nil
}

// VerifyResetPasswordOTP verifies reset password otp code
func (s *AuthServiceImpl) VerifyResetPasswordOTP(ctx context.Context, req *entity.VerifyOTPRequest) error {
	_, err := s.verifyOTP(ctx, req)
	if err != nil {
		return entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	return nil
}

func (s *AuthServiceImpl) UpdatePassword(ctx context.Context, req *entity.UpdatePasswordRequest) error {
	// Get OTP record
	otp, err := s.repo.GetOTPByToken(ctx, req.OTPToken)
	if err != nil {
		return entity.NewAuthError(entity.ErrOTPInvalid, entity.MsgOTPInvalid)
	}

	// Check OTP status
	if !otp.Status {
		return entity.NewAuthError(entity.ErrInvalidRequest, entity.MsgInvalidRequest)
	}

	// Check if OTP is verified within 3 mins
	if time.Since(otp.UpdatedAt.UTC()) > 3*time.Minute {
		return fmt.Errorf("failed to update password: OTP expired")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	err = s.repo.UpdatePassword(ctx, otp.UserID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// CheckPermission checks if a user has a specific permission
func (s *AuthServiceImpl) CheckPermission(ctx context.Context, userID uuid.UUID, resource entity.Resource, operation entity.Operation) (bool, error) {
	hasPermission, err := s.repo.CheckUserPermission(ctx, userID, string(resource), string(operation))
	if err != nil {
		s.logger.Printf("Error checking permission for user %s: %v", userID, err)
		return false, err
	}

	if !hasPermission {
		s.LogActivity(ctx, userID, entity.ACTIVITY_PERMISSION_DENIED, "unknown", "unknown", "unknown", false, map[string]interface{}{
			"resource":  resource,
			"operation": operation,
		})
	}

	return hasPermission, nil
}

// GetUserPermissions gets all permissions for a user
func (s *AuthServiceImpl) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.repo.GetUserPermissions(ctx, userID)
}

// GetUserProfile retrieves user profile information
func (s *AuthServiceImpl) GetUserProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// GetUserByPhone retrieves user by phone number
func (s *AuthServiceImpl) GetUserByPhone(ctx context.Context, phonePrefix, phoneNumber string) (*entity.User, error) {
	return s.repo.GetUserByPhone(ctx, phonePrefix, phoneNumber)
}

// GetMerchantsByUser retrieves all merchants for a user
func (s *AuthServiceImpl) GetMerchantsByUser(ctx context.Context, userID uuid.UUID) ([]entity.Merchant, error) {
	return s.repo.GetMerchantsByUser(ctx, userID)
}

// UpdateUser updates user information
func (s *AuthServiceImpl) UpdateUser(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	return s.repo.UpdateUser(ctx, userID, updates)
}

// UpdateUserProfile updates user profile information
func (s *AuthServiceImpl) UpdateUserProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	// This is a simple implementation - in production you'd want more validation
	err := s.repo.UpdateUser(ctx, userID, updates)
	if err != nil {
		s.logger.Printf("Error updating user profile: %v", err)
		return entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}
	return nil
}

// ChangePassword changes user password
func (s *AuthServiceImpl) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Verify old password
	valid, err := s.repo.VerifyPassword(ctx, userID, oldPassword)
	if err != nil || !valid {
		return entity.NewAuthError(entity.ErrInvalidCredentials, "Current password is incorrect")
	}

	// Validate new password
	if err := utils.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Update password
	err = s.repo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return entity.NewAuthError(entity.ErrInternalServer, "Failed to update password")
	}

	// Log activity
	s.LogActivity(ctx, userID, entity.ACTIVITY_PASSWORD_CHANGED, "unknown", "unknown", "unknown", true, nil)

	return nil
}

// RefreshToken refreshes an authentication token
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*entity.AuthResponse, error) {
	// TODO: Implement refresh token logic
	return nil, entity.NewAuthError(entity.ErrInternalServer, "Refresh token not implemented yet")
}

// Logout logs out a user
func (s *AuthServiceImpl) Logout(ctx context.Context, token string) error {
	// Validate token and get session
	claims, err := utils.ValidateJWT(token, s.jwtSecret)
	if err != nil {
		return entity.NewAuthError(entity.ErrInvalidToken, entity.MsgInvalidToken)
	}

	userID, _ := uuid.Parse(claims.UserID)
	sessionID, _ := uuid.Parse(claims.SessionID)

	// Revoke session
	err = s.repo.RevokeSession(ctx, sessionID)
	if err != nil {
		s.logger.Printf("Error revoking session: %v", err)
		return entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Log activity
	s.LogActivity(ctx, userID, entity.ACTIVITY_LOGOUT, "unknown", "unknown", "unknown", true, nil)

	return nil
}

// CreateAdminUser creates a new admin user
func (s *AuthServiceImpl) CreateSuperAdminUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.User, error) {

	// Validate request
	if err := utils.ValidateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// Create user
	user, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		s.logger.Printf("Error creating admin user: %v", err)
		return nil, entity.NewAuthError(entity.ErrAccountCreationFailed, entity.MsgAccountCreationFailed)
	}

	// Get or create super admin group
	group, err := s.repo.GetOrCreateSuperAdminGroup(ctx)
	if err != nil {
		s.logger.Printf("Error getting super admin group: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, entity.MsgInternalServer)
	}

	// Assign user to super admin group
	err = s.repo.AssignUserToGroup(ctx, user.ID, group.ID)
	if err != nil {
		s.logger.Printf("Error assigning user to group: %v", err)
		return nil, entity.NewAuthError(entity.ErrInternalServer, "Failed to assign admin permissions")
	}

	return user, nil
}

// AssignUserToGroup assigns a user to a group
func (s *AuthServiceImpl) AssignUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	return s.repo.AssignUserToGroup(ctx, userID, groupID)
}

// LogActivity logs an authentication activity
func (s *AuthServiceImpl) LogActivity(ctx context.Context, userID uuid.UUID, activityType entity.AuthActivityType, ip, userAgent, deviceName string, success bool, details map[string]interface{}) error {
	activity := &entity.AuthActivity{
		ID:           uuid.New(),
		UserID:       userID,
		ActivityType: activityType,
		IPAddress:    ip,
		UserAgent:    userAgent,
		DeviceName:   deviceName,
		Success:      success,
		Details:      details,
		CreatedAt:    time.Now(),
	}

	return s.repo.LogAuthActivity(ctx, activity)
}

// GetUserActivities gets user activities
func (s *AuthServiceImpl) GetUserActivities(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.AuthActivity, error) {
	return s.repo.GetUserActivities(ctx, userID, limit, offset)
}

// SeedSuperAdmin creates the default super admin user
func (s *AuthServiceImpl) SeedSuperAdmin(ctx context.Context) error {
	// Check if super admin already exists
	exists, err := s.repo.UserExists(ctx, "251", "961186323")
	if err != nil {
		return err
	}
	if exists {
		s.logger.Println("Super admin user already exists, skipping seed")
		return nil
	}

	// Create super admin user
	req := &entity.CreateUserRequest{
		Title:        "Mr",
		FirstName:    "SocialPay",
		LastName:     "SuperAdmin",
		PhonePrefix:  "251",
		PhoneNumber:  "961186323",
		Password:     "SocialPay$123SuperAdmiN",
		PasswordHint: "superadmin",
		UserType:     entity.USER_TYPE_ADMIN,
	}

	user, err := s.CreateSuperAdminUser(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create super admin: %v", err)
	}

	s.logger.Printf("Super admin user created successfully with ID: %s", user.ID)
	return nil
}

// GetUserPermissionsByMerchant gets all permissions for a user within a specific merchant context
func (s *AuthServiceImpl) GetUserPermissionsByMerchant(ctx context.Context, userID, merchantID uuid.UUID) ([]string, error) {
	return s.repo.GetUserPermissionsByMerchant(ctx, userID, merchantID)
}

// GetAllUserPermissions gets all permissions for a user grouped by merchant ID
func (s *AuthServiceImpl) GetAllUserPermissions(ctx context.Context, userID uuid.UUID) (map[string][]string, error) {
	return s.repo.GetAllUserPermissions(ctx, userID)
}

// GetUserGroupsByGroupedByMerchant gets all groups for a user grouped by merchant ID
func (s *AuthServiceImpl) GetUserGroupsByGroupedByMerchant(ctx context.Context, userID uuid.UUID) (map[string][]entity.Group, error) {
	return s.repo.GetAllUserGroups(ctx, userID)
}

// CheckPermissionForMerchant checks if a user has a specific permission within a merchant context
func (s *AuthServiceImpl) CheckPermissionForMerchant(ctx context.Context, userID, merchantID uuid.UUID, resource entity.Resource, operation entity.Operation) (bool, error) {
	hasPermission, err := s.repo.CheckUserPermissionForMerchant(ctx, userID, merchantID, string(resource), string(operation))
	if err != nil {
		s.logger.Printf("Error checking permission for user %s in merchant %s: %v", userID, merchantID, err)
		return false, err
	}

	if !hasPermission {
		s.LogActivity(ctx, userID, entity.ACTIVITY_PERMISSION_DENIED, "unknown", "unknown", "unknown", false, map[string]interface{}{
			"resource":    resource,
			"operation":   operation,
			"merchant_id": merchantID,
		})
	}

	return hasPermission, nil
}

// CreateDevice create new device
func (s *AuthServiceImpl) CreateDevice(ctx context.Context, args entity.CreateDeviceArgs) (*entity.Device, error) {
	return s.repo.CreateDevice(ctx, args.IP, args.Name, args.Agent)
}

// CreateSesssion create new session
func (s *AuthServiceImpl) CreateSession(ctx context.Context, args entity.CreateSessionArgs) (*entity.Session, error) {
	return s.repo.CreateSession(ctx, args.UserID, args.DeviceID, args.Token, args.RefreshToken, args.ExpiresAt)
}
