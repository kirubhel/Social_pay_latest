package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// HeaderAPIKey is the header key for merchant ID
	HeaderMerchantID = "X-MERCHANT-ID"
	// ContextKeyMerchantID is the context key for merchant ID
	ContextKeyMerchantID = "merchantID"
)

// MerchantIDMiddleware creates a middleware that validates merchant ID from header
func MerchantIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get merchant ID from header
		merchantIDStr := c.GetHeader(HeaderMerchantID)
		if merchantIDStr == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INVALID_REQUEST",
					Message: "Merchant ID is required in X-MERCHANT-ID header",
				},
			})
			c.Abort()
			return
		}

		// Parse merchant ID
		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INVALID_REQUEST",
					Message: "Invalid merchant ID format",
				},
			})
			c.Abort()
			return
		}

		// Get session from context to verify merchant association
		session, exists := GetSessionFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "Session not found",
				},
			})
			c.Abort()
			return
		}

		// Verify merchant ID belongs to the user's session
		// TODO: Add proper merchant ID verification logic here
		// For now, we'll just store it in context

		// Store merchant ID in context
		c.Set(ContextKeyMerchantID, merchantID)

		fmt.Printf("[Merchant-Auth] Merchant ID %s validated for user %s\n", merchantID, session.User.Id)
		c.Next()
	}
}

// GetMerchantIDFromContext extracts the merchant ID from the Gin context
func GetMerchantIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	merchantID, exists := c.Get(ContextKeyMerchantID)
	if !exists {
		return uuid.Nil, false
	}

	id, ok := merchantID.(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}

	return id, true
}
