package gin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"
	"github.com/socialpay/socialpay/src/pkg/auth/usecase"
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
	// UseCase for JWT validation
	AuthUseCase usecase.AuthNInteractor
	// If true, routes will be accessible without authentication
	Public bool
}

// Constants for context keys
const (
	ContextKeySession = "session"
	ContextKeyUserID  = "userID"
)

// JWTAuthMiddleware creates a new JWT authentication middleware for Gin
func JWTAuthMiddleware(config JWTAuthMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("[JWT-Auth] Processing request for path: %s\n", c.Request.URL.Path)
		fmt.Printf("[JWT-Auth] Method: %s\n", c.Request.Method)
		fmt.Printf("[JWT-Auth] Is public route: %v\n", config.Public)

		// If route is public, skip authentication
		if config.Public {
			fmt.Println("[JWT-Auth] Public route, skipping authentication")
			c.Next()
			return
		}

		// Get JWT token from Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("[JWT-Auth] Authorization header: %s\n", maskToken(authHeader))

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("[JWT-Auth] Missing or invalid Authorization header format")
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
		fmt.Printf("[JWT-Auth] Extracted token: %s\n", maskToken(token))

		// Validate token
		fmt.Println("[JWT-Auth] Validating token with auth service")
		session, err := config.AuthUseCase.CheckSession(token)
		if err != nil {
			fmt.Printf("[JWT-Auth] Token validation failed: %v\n", err)
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

		fmt.Printf("[JWT-Auth] Token validated successfully for user ID: %s\n", session.User.Id)

		// Set session in context
		c.Set(ContextKeySession, session)
		// Also set user ID for convenience
		c.Set(ContextKeyUserID, session.User.Id)

		fmt.Printf("[JWT-Auth] Session and UserID set in context. Session: %+v\n", session)
		fmt.Printf("[JWT-Auth] Context keys available: %v\n", c.Keys)

		// Continue to the next middleware/handler
		c.Next()
	}
}

// GetSessionFromContext extracts the session from the Gin context
func GetSessionFromContext(c *gin.Context) (*entity.Session, bool) {
	fmt.Printf("[JWT-Auth] Getting session from context. Available keys: %v\n", c.Keys)
	session, exists := c.Get(ContextKeySession)
	if !exists {
		fmt.Println("[JWT-Auth] Session not found in context")
		return nil, false
	}

	sessionObj, ok := session.(*entity.Session)
	if !ok {
		fmt.Printf("[JWT-Auth] Invalid session type in context. Got: %T\n", session)
		return nil, false
	}

	fmt.Printf("[JWT-Auth] Session retrieved successfully: %+v\n", sessionObj)
	return sessionObj, true
}

// GetUserIDFromContext extracts the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	fmt.Printf("[JWT-Auth] Getting userID from context. Available keys: %v\n", c.Keys)
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		fmt.Println("[JWT-Auth] UserID not found in context, trying session fallback")
		// Try to get it from session as fallback
		session, ok := GetSessionFromContext(c)
		if !ok {
			fmt.Println("[JWT-Auth] Session fallback failed")
			return uuid.Nil, false
		}
		fmt.Printf("[JWT-Auth] Got userID from session: %s\n", session.User.Id)
		return session.User.Id, true
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		fmt.Printf("[JWT-Auth] Invalid userID type in context. Got: %T\n", userID)
		return uuid.Nil, false
	}

	fmt.Printf("[JWT-Auth] UserID retrieved successfully: %s\n", id)
	return id, true
}

// RequireJWTAuthentication is a helper middleware that ensures a user is authenticated with JWT
func RequireJWTAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("[JWT-Auth] Checking JWT authentication for path: %s\n", c.Request.URL.Path)
		_, exists := GetSessionFromContext(c)
		if !exists {
			fmt.Println("[JWT-Auth] Authentication check failed - no session found")
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
		fmt.Println("[JWT-Auth] Authentication check passed")
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
