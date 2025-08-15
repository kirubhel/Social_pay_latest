package gin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService service.AuthService
	logger      logging.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logging.NewStdLogger("auth_handler"),
	}
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Title        string `json:"title" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	PhonePrefix  string `json:"phone_prefix" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	Password     string `json:"password" binding:"required"`
	PasswordHint string `json:"password_hint,omitempty"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	PhonePrefix string `json:"phone_prefix" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// VerifyOTPRequest represents the OTP verification request payload
type VerifyOTPRequest struct {
	OTPToken string `json:"otp_token" binding:"required"`
	OTPCode  string `json:"otp_code" binding:"required"`
}

// RefreshTokenRequest represents the token refresh request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdatePasswordRequest represents args for updating password
type UpdatePasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required"`
	OTPToken    string `json:"otp_token" validate:"required"`
}

// AdminLoginRequest represents the admin login request payload
type AdminLoginRequest struct {
	PhonePrefix string `json:"phone_prefix" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// ResetPasswordRequest represents args for reset password
type RequestOTPRequest struct {
	PhonePrefix string `json:"phone_prefix" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

// ResetPasswordRequest represents args for update password
type ResetPasswordRequest struct {
	UpdatePasswordRequest
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account with phone verification
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 409 {object} map[string]interface{} "Phone already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	// Get device info from request
	deviceInfo := &entity.DeviceInfo{
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		DeviceName: "web", // Default for web requests
	}

	// Convert to entity request
	entityReq := &entity.CreateUserRequest{
		Title:        req.Title,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhonePrefix:  req.PhonePrefix,
		PhoneNumber:  req.PhoneNumber,
		Password:     req.Password,
		PasswordHint: req.PasswordHint,
		UserType:     entity.USER_TYPE_MERCHANT,
		DeviceInfo:   deviceInfo,
	}

	// Register user
	authResponse, err := h.authService.Register(c.Request.Context(), entityReq)
	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Registration failed",
				},
			})
		}
		return
	}

	// Get all user permissions grouped by merchant ID
	allPermissions, err := h.authService.GetAllUserPermissions(c.Request.Context(), authResponse.User.ID)
	if err != nil {
		// Log error but don't fail the request
		allPermissions = make(map[string][]string)
	}

	// Build merchant-aware resource-operation map for better frontend UX
	merchantPermissions := make(map[string]map[string][]string)
	globalPermissions := []string{}
	globalResourceOperations := make(map[string][]string)

	for merchantID, permissions := range allPermissions {
		if merchantID == "global" {
			globalPermissions = permissions
			// Build global resource-operation map
			for _, perm := range permissions {
				parts := strings.Split(perm, ":")
				if len(parts) == 2 {
					resource := parts[0]
					operation := parts[1]
					if _, exists := globalResourceOperations[resource]; !exists {
						globalResourceOperations[resource] = []string{}
					}
					globalResourceOperations[resource] = append(globalResourceOperations[resource], operation)
				}
			}
		} else {
			// Build merchant-specific resource-operation map
			resourceOperations := make(map[string][]string)
			for _, perm := range permissions {
				parts := strings.Split(perm, ":")
				if len(parts) == 2 {
					resource := parts[0]
					operation := parts[1]
					if _, exists := resourceOperations[resource]; !exists {
						resourceOperations[resource] = []string{}
					}
					resourceOperations[resource] = append(resourceOperations[resource], operation)
				}
			}
			merchantPermissions[merchantID] = resourceOperations
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data": gin.H{
			"user":                         authResponse.User,
			"merchants":                    authResponse.Merchants,
			"access_token":                 authResponse.Token,
			"refresh_token":                authResponse.RefreshToken,
			"expires_at":                   authResponse.ExpiresAt,
			"permissions":                  globalPermissions,        // Deprecated: backward compatibility
			"resource_operations":          globalResourceOperations, // Deprecated: backward compatibility
			"merchant_permissions":         allPermissions,           // New: permissions grouped by merchant_id
			"merchant_resource_operations": merchantPermissions,      // New: resource-operations grouped by merchant_id
		},
	})
}

// Login handles user login and sends OTP
// @Summary Login user
// @Description Authenticate user credentials and send OTP
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "OTP sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Check if user is admin/super_admin and reject
	user, err := h.authService.GetUserByPhone(c.Request.Context(), req.PhonePrefix, req.PhoneNumber)

	h.logger.Info("user", map[string]interface{}{
		"user": user,
	})

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "UNAUTHORIZED",
			"message": "User not found",
		})
		return
	}

	if err == nil && (user.UserType == entity.USER_TYPE_ADMIN || user.UserType == entity.USER_TYPE_SUPER_ADMIN) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "ADMIN_LOGIN_REQUIRED",
			"message": "Wrong user type",
		})
		return
	}

	// Convert to entity request
	entityReq := &entity.LoginRequest{
		PhonePrefix: req.PhonePrefix,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		DeviceInfo: &entity.DeviceInfo{
			IPAddress:  c.ClientIP(),
			DeviceName: c.GetHeader("X-Device-Name"),
			UserAgent:  c.GetHeader("User-Agent"),
		},
	}

	// Authenticate and send OTP
	otpToken, err := h.authService.Login(c.Request.Context(), entityReq)
	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"error":   authErr.Type,
				"message": authErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "OTP sent successfully",
		"otp_token": otpToken,
	})
}

// VerifyOTP handles OTP verification and completes login
// @Summary Verify OTP
// @Description Verify OTP code and complete login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "OTP verification details"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid or expired OTP"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Convert to entity request
	entityReq := &entity.VerifyOTPRequest{
		Token: req.OTPToken,
		Code:  req.OTPCode,
		DeviceInfo: &entity.DeviceInfo{
			IPAddress:  c.ClientIP(),
			DeviceName: c.GetHeader("X-Device-Name"),
			UserAgent:  c.GetHeader("User-Agent"),
		},
	}

	// Verify OTP and complete login
	authResponse, err := h.authService.VerifyOTP(c.Request.Context(), entityReq)
	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   authErr.Type,
				"message": authErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	// Get all user permissions grouped by merchant ID
	allPermissions, err := h.authService.GetAllUserPermissions(c.Request.Context(), authResponse.User.ID)
	if err != nil {
		// Log error but don't fail the request
		allPermissions = make(map[string][]string)
	}

	h.logger.Info("allPermissions", map[string]interface{}{
		"allPermissions": allPermissions,
	})

	// Build merchant-aware resource-operation map for better frontend UX
	merchantPermissions := make(map[string]map[string][]string)
	globalResourceOperations := make(map[string][]string)

	for merchantID, permissions := range allPermissions {
		if merchantID == "global" {
			// Build global resource-operation map
			for _, perm := range permissions {
				parts := strings.Split(perm, ":")
				if len(parts) == 2 {
					resource := parts[0]
					operation := parts[1]
					if _, exists := globalResourceOperations[resource]; !exists {
						globalResourceOperations[resource] = []string{}
					}
					globalResourceOperations[resource] = append(globalResourceOperations[resource], operation)
				}
			}
		} else {
			// Build merchant-specific resource-operation map
			resourceOperations := make(map[string][]string)
			for _, perm := range permissions {
				parts := strings.Split(perm, ":")
				if len(parts) == 2 {
					resource := parts[0]
					operation := parts[1]
					if _, exists := resourceOperations[resource]; !exists {
						resourceOperations[resource] = []string{}
					}
					resourceOperations[resource] = append(resourceOperations[resource], operation)
				}
			}
			merchantPermissions[merchantID] = resourceOperations
		}
	}
	h.logger.Info("merchantPermissions", map[string]interface{}{
		"user":                authResponse.User,
		"merchantPermissions": merchantPermissions,
		"allPermissions":      allPermissions,
	})

	// Get user groups
	groups, err := h.authService.GetUserGroupsByGroupedByMerchant(c.Request.Context(), authResponse.User.ID)
	if err != nil {
		// Log error but don't fail the request
		groups = make(map[string][]entity.Group)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": gin.H{
			"user":                         authResponse.User,
			"merchants":                    authResponse.Merchants,
			"access_token":                 authResponse.Token,
			"refresh_token":                authResponse.RefreshToken,
			"expires_at":                   authResponse.ExpiresAt,
			"groups":                       groups,
			"merchant_permissions":         allPermissions,      // New: permissions grouped by merchant_id
			"merchant_resource_operations": merchantPermissions, // New: resource-operations grouped by merchant_id
			"global_resource_operations":   globalResourceOperations,
		},
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid refresh token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Refresh token
	authResponse, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   authErr.Type,
				"message": authErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Token refreshed successfully",
		"access_token":  authResponse.Token,
		"refresh_token": authResponse.RefreshToken,
		"expires_at":    authResponse.ExpiresAt,
	})
}

// Logout handles user logout
// @Summary Logout user
// @Description Logout user and invalidate session
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logged out successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "UNAUTHORIZED",
			"message": "Authorization header required",
		})
		return
	}

	// Extract token (remove "Bearer " prefix)
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// Logout
	err := h.authService.Logout(c.Request.Context(), token)
	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   authErr.Type,
				"message": authErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User profile"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user from context (set by middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "UNAUTHORIZED",
			"message": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*entity.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Invalid user context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param updates body map[string]interface{} true "Profile updates"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "UNAUTHORIZED",
			"message": "User ID not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Invalid user ID context",
		})
		return
	}

	// Parse updates
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Update profile
	err := h.authService.UpdateUser(c.Request.Context(), userID, updates)
	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   authErr.Type,
				"message": authErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}

// GetActivities handles getting user authentication activities
// @Summary Get user activities
// @Description Get authenticated user's authentication activities
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit number of activities" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} map[string]interface{} "User activities"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/activities [get]
func (h *AuthHandler) GetActivities(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "UNAUTHORIZED",
			"message": "User ID not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Invalid user ID context",
		})
		return
	}

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get activities
	activities, err := h.authService.GetUserActivities(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to get user activities",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(activities),
		},
	})
}

// AdminLogin handles admin login and sends OTP (restricted to super_admin and admin only)
// @Summary Admin Login
// @Description Authenticate admin credentials and send OTP (super_admin and admin only)
// @Tags admin-auth
// @Accept json
// @Produce json
// @Param request body AdminLoginRequest true "Admin login credentials"
// @Success 200 {object} map[string]interface{} "OTP sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials or insufficient privileges"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/auth/login [post]
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Convert to entity request
	entityReq := &entity.LoginRequest{
		PhonePrefix: req.PhonePrefix,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		DeviceInfo: &entity.DeviceInfo{
			IPAddress:  c.ClientIP(),
			DeviceName: c.GetHeader("X-Device-Name"),
			UserAgent:  c.GetHeader("User-Agent"),
		},
	}

	// First, get user to check user type before sending OTP
	user, err := h.authService.GetUserByPhone(c.Request.Context(), req.PhonePrefix, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "INVALID_CREDENTIALS",
			"message": "Invalid credentials",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "INVALID_CREDENTIALS",
			"message": "Invalid credentials",
		})
		return
	}
	h.logger.Info("user logged in", map[string]interface{}{
		"user": user,
	})

	// Check if user is admin or super_admin
	if user.UserType != entity.USER_TYPE_ADMIN && user.UserType != entity.USER_TYPE_SUPER_ADMIN {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "INSUFFICIENT_PRIVILEGES",
			"message": "Access denied. Admin privileges required.",
		})
		return
	}

	// Authenticate and send OTP
	otpToken, err := h.authService.Login(c.Request.Context(), entityReq)
	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"error":   authErr.Type,
				"message": authErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "OTP sent successfully",
		"otp_token": otpToken,
		"user_type": string(user.UserType),
	})
}

// RequestOTP handles request otp
// @Summary Request OTP
// @Description Users can request otp code for password management
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RequestOTPRequest true "Request OTP request body"
// @Success 200 {object} map[string]interface{} "OTP sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials or insufficient privileges"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/request-otp [post]
func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	otpToken, err := h.authService.RequestOTP(c, req.PhonePrefix, req.PhoneNumber)
	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"error":   authErr.Type,
				"message": authErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "OTP sent successfully",
		"otp_token": otpToken,
	})
}

// VerifyResetPasswordOTP handles OTP verification and completes reset password
// @Summary Verify OTP
// @Description Verify OTP code and complete reset password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "OTP verification details"
// @Success 200 {object} map[string]interface{} "OTP verified successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid or expired OTP"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/password/verify-otp [post]
func (h *AuthHandler) VerifyResetPasswordOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Convert to entity request
	entityReq := &entity.VerifyOTPRequest{
		Token: req.OTPToken,
		Code:  req.OTPCode,
		DeviceInfo: &entity.DeviceInfo{
			IPAddress:  c.ClientIP(),
			DeviceName: c.GetHeader("X-Device-Name"),
			UserAgent:  c.GetHeader("User-Agent"),
		},
	}

	err := h.authService.VerifyResetPasswordOTP(c, entityReq)

	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"error":   authErr.Type,
				"message": authErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
	})
}

// ResetPassword handles reset password
// @Summary Reset Password
// @Description Users can reset their password.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset password details"
// @Success 200 {object} map[string]interface{} "Password resetted successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid or expired OTP"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	err := h.authService.UpdatePassword(c, &entity.UpdatePasswordRequest{
		NewPassword: req.NewPassword,
		OTPToken:    req.OTPToken,
	})

	if err != nil {
		if authErr, ok := err.(*entity.AuthError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"error":   authErr.Type,
				"message": authErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password resetted successfully",
	})
}

// RegisterRoutes registers all authentication routes
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/verify-otp", h.VerifyOTP)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
		auth.GET("/profile", h.GetProfile)
		auth.PUT("/profile", h.UpdateProfile)
		auth.GET("/activities", h.GetActivities)
		auth.POST("/request-otp", h.RequestOTP)
		auth.POST("/password/verify-otp", h.VerifyResetPasswordOTP)
		auth.POST("/reset-password", h.ResetPassword)
	}
}

// RegisterAdminRoutes registers admin authentication routes
func (h *AuthHandler) RegisterAdminRoutes(router *gin.RouterGroup) {
	adminAuth := router.Group("/admin/auth")
	{
		adminAuth.POST("/login", h.AdminLogin)
		// Admin can also use regular verify-otp and refresh endpoints
		adminAuth.POST("/verify-otp", h.VerifyOTP)
		adminAuth.POST("/refresh", h.RefreshToken)
		adminAuth.POST("/logout", h.Logout)
	}
}
