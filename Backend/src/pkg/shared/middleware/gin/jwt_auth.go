package gin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	"github.com/socialpay/socialpay/src/pkg/authv2/utils"
)

// ErrorResponse represents the standard error response structure
type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   ApiError `json:"error"`
}

// ApiError represents detailed error information
type ApiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// JWTAuthMiddlewareConfig holds configuration for the JWT auth middleware
type JWTAuthMiddlewareConfig struct {
	// AuthService for JWT validation using authv2
	AuthService service.AuthService
	// JWT Secret for token validation
	JWTSecret string
	// If true, routes will be accessible without authentication
	Public bool
}

// Constants for context keys
const (
	ContextKeySession = "session"
	ContextKeyUserID  = "userID"
	ContextKeyUser    = "user"
)

// JWTAuthMiddleware creates a new JWT authentication middleware for Gin using authv2
func JWTAuthMiddleware(config JWTAuthMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("[JWT-Auth-V2] Processing request for path: %s\n", c.Request.URL.Path)
		fmt.Printf("[JWT-Auth-V2] Method: %s\n", c.Request.Method)
		fmt.Printf("[JWT-Auth-V2] Is public route: %v\n", config.Public)

		// If route is public, skip authentication
		if config.Public {
			fmt.Println("[JWT-Auth-V2] Public route, skipping authentication")
			c.Next()
			return
		}

		// Get JWT token from Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("[JWT-Auth-V2] Authorization header: %s\n", maskToken(authHeader))

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("[JWT-Auth-V2] Missing or invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "Authentication token missing in header",
				},
			})
			c.Abort()
			return
		}

		// Extract token from header
		token := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Printf("[JWT-Auth-V2] Extracted token: %s\n", maskToken(token))

		// Validate token using authv2 JWT utilities
		fmt.Println("[JWT-Auth-V2] Validating token with authv2 JWT utilities")
		claims, err := utils.ValidateJWT(token, config.JWTSecret)
		if err != nil {
			fmt.Printf("[JWT-Auth-V2] Token validation failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "Invalid authentication token",
				},
			})
			c.Abort()
			return
		}

		// Parse user ID from claims
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			fmt.Printf("[JWT-Auth-V2] Invalid user ID in token: %v\n", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "Invalid token format",
				},
			})
			c.Abort()
			return
		}

		// Get user profile from authv2 service
		user, err := config.AuthService.GetUserProfile(c.Request.Context(), userID)
		if err != nil {
			fmt.Printf("[JWT-Auth-V2] Failed to get user profile: %v\n", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "User not found",
				},
			})
			c.Abort()
			return
		}

		fmt.Printf("[JWT-Auth-V2] Token validated successfully for user ID: %s\n", user.ID)

		// Create session-like structure for backwards compatibility
		session := &entity.Session{
			UserID: user.ID,
			Token:  token,
		}

		// Set context values
		c.Set(ContextKeySession, session)
		c.Set(ContextKeyUserID, user.ID)
		c.Set(ContextKeyUser, user)

		fmt.Printf("[JWT-Auth-V2] Session, UserID, and User set in context. User: %+v\n", user)
		fmt.Printf("[JWT-Auth-V2] Context keys available: %v\n", c.Keys)

		// Continue to the next middleware/handler
		c.Next()
	}
}

// GetSessionFromContext extracts the session from the Gin context
func GetSessionFromContext(c *gin.Context) (*entity.Session, bool) {
	fmt.Printf("[JWT-Auth-V2] Getting session from context. Available keys: %v\n", c.Keys)
	session, exists := c.Get(ContextKeySession)
	if !exists {
		fmt.Println("[JWT-Auth-V2] Session not found in context")
		return nil, false
	}

	sessionObj, ok := session.(*entity.Session)
	if !ok {
		fmt.Printf("[JWT-Auth-V2] Invalid session type in context. Got: %T\n", session)
		return nil, false
	}

	fmt.Printf("[JWT-Auth-V2] Session retrieved successfully: %+v\n", sessionObj)
	return sessionObj, true
}

// GetUserIDFromContext extracts the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	fmt.Printf("[JWT-Auth-V2] Getting userID from context. Available keys: %v\n", c.Keys)
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		fmt.Println("[JWT-Auth-V2] UserID not found in context, trying session fallback")
		// Try to get it from session as fallback
		session, ok := GetSessionFromContext(c)
		if !ok {
			fmt.Println("[JWT-Auth-V2] Session fallback failed")
			return uuid.Nil, false
		}
		fmt.Printf("[JWT-Auth-V2] Got userID from session: %s\n", session.UserID)
		return session.UserID, true
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		fmt.Printf("[JWT-Auth-V2] Invalid userID type in context. Got: %T\n", userID)
		return uuid.Nil, false
	}

	fmt.Printf("[JWT-Auth-V2] UserID retrieved successfully: %s\n", id)
	return id, true
}

// GetUserFromContext extracts the user from the Gin context
func GetUserFromContext(c *gin.Context) (*entity.User, bool) {
	fmt.Printf("[JWT-Auth-V2] Getting user from context. Available keys: %v\n", c.Keys)
	user, exists := c.Get(ContextKeyUser)
	if !exists {
		fmt.Println("[JWT-Auth-V2] User not found in context")
		return nil, false
	}

	userObj, ok := user.(*entity.User)
	if !ok {
		fmt.Printf("[JWT-Auth-V2] Invalid user type in context. Got: %T\n", user)
		return nil, false
	}

	fmt.Printf("[JWT-Auth-V2] User retrieved successfully: %+v\n", userObj)
	return userObj, true
}

// RequireJWTAuthentication is a helper middleware that ensures a user is authenticated with JWT
func RequireJWTAuthentication(authService service.AuthService, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("[JWT-Auth-V2] Checking JWT authentication for path: %s\n", c.Request.URL.Path)

		// Get and validate Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("[JWT-Auth-V2] Authorization header: %s\n", maskToken(authHeader))

		if authHeader == "" {
			fmt.Println("[JWT-Auth-V2] No Authorization header found")
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("[JWT-Auth-V2] Invalid Authorization header format: %s\n", authHeader)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "Invalid Authorization header format. Expected: Bearer <token>",
				},
			})
			c.Abort()
			return
		}

		token := parts[1]
		fmt.Printf("[JWT-Auth-V2] Processing token: %s\n", maskToken(token))

		// Validate JWT token
		claims, err := utils.ValidateJWT(token, jwtSecret)
		if err != nil {
			fmt.Printf("[JWT-Auth-V2] Failed to validate JWT: %v\n", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		// Parse user ID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			fmt.Printf("[JWT-Auth-V2] Invalid user ID in claims: %v\n", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "Invalid token format",
				},
			})
			c.Abort()
			return
		}

		// Get user profile
		user, err := authService.GetUserProfile(c.Request.Context(), userID)
		if err != nil {
			fmt.Printf("[JWT-Auth-V2] Failed to get user profile: %v\n", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "User not found",
				},
			})
			c.Abort()
			return
		}

		// Create session and set context
		session := &entity.Session{
			UserID: user.ID,
			Token:  token,
		}

		c.Set(ContextKeySession, session)
		c.Set(ContextKeyUserID, user.ID)
		c.Set(ContextKeyUser, user)

		fmt.Printf("[JWT-Auth-V2] Authentication successful for user: %s\n", user.ID)
		c.Next()
	}
}

// maskToken masks the token for safe logging
func maskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 10 {
		return "***"
	}
	return token[:5] + "..." + token[len(token)-5:]
}

// RequireUserType creates middleware that checks for specific user types
func RequireUserType(userTypes ...entity.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		user, exists := GetUserFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "User not authenticated",
				},
			})
			c.Abort()
			return
		}

		// Check if user type is allowed
		allowed := false
		for _, allowedType := range userTypes {
			if user.UserType == allowedType {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "FORBIDDEN",
					Message: "Insufficient user type privileges",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware that ensures user is an admin
func RequireAdmin() gin.HandlerFunc {
	return RequireUserType(entity.USER_TYPE_ADMIN)
}

// RequireMerchant middleware that ensures user is a merchant
func RequireMerchant() gin.HandlerFunc {
	return RequireUserType(entity.USER_TYPE_MERCHANT)
}
