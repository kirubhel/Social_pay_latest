package gin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/entity"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/usecase"
)

const (
	// HeaderAPIKey is the header key for API key authentication
	HeaderAPIKey = "X-API-Key"
)

// APIKeyAuth middleware for API key authentication
func APIKeyAuth(useCase usecase.APIKeyUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("[API-Key-Auth] Processing request for path: %s\n", c.Request.URL.Path)
		fmt.Printf("[API-Key-Auth] Method: %s\n", c.Request.Method)

		// Get API key from header
		apiKeyHeader := c.GetHeader(HeaderAPIKey)
		if apiKeyHeader == "" {
			fmt.Println("[API-Key-Auth] Missing API key header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "UNAUTHORIZED",
					"message": "API key is required",
				},
			})
			c.Abort()
			return
		}

		fmt.Printf("[API-Key-Auth] API key header present: %s\n", maskAPIKey(apiKeyHeader))

		// Split the API key into public and secret parts
		parts := strings.Split(apiKeyHeader, ":")
		if len(parts) != 2 {
			fmt.Println("[API-Key-Auth] Invalid API key format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "UNAUTHORIZED",
					"message": "Invalid API key format",
				},
			})
			c.Abort()
			return
		}

		publicKey, secretKey := parts[0], parts[1]
		fmt.Printf("[API-Key-Auth] Public key: %s\n", maskAPIKey(publicKey))

		// Validate API key
		fmt.Println("[API-Key-Auth] Validating API key with auth service")
		apiKeyData, err := useCase.ValidateAPIKey(c.Request.Context(), publicKey, secretKey)
		if err != nil {
			fmt.Printf("[API-Key-Auth] API key validation failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "UNAUTHORIZED",
					"message": "Invalid API key",
				},
			})
			c.Abort()
			return
		}

		fmt.Printf("[API-Key-Auth] API key validated successfully for user ID: %s\n", apiKeyData.UserID)
		fmt.Printf("[API-Key-Auth] API key permissions - CanWithdrawal: %v, CanProcessPayment: %v\n",
			apiKeyData.CanWithdrawal, apiKeyData.CanProcessPayment)

		// Store API key data in context
		c.Set("apiKey", apiKeyData)
		c.Set("userID", apiKeyData.UserID.String())

		fmt.Printf("[API-Key-Auth] Context keys set - Available keys: %v\n", c.Keys)

		c.Next()
	}
}

// RequireWithdrawalPermission middleware to check withdrawal permission
func RequireWithdrawalPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("[API-Key-Auth] Checking withdrawal permission for path: %s\n", c.Request.URL.Path)

		apiKeyData, exists := c.Get("apiKey")
		if !exists {
			fmt.Println("[API-Key-Auth] API key data not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "UNAUTHORIZED",
					"message": "API key data not found",
				},
			})
			c.Abort()
			return
		}

		apiKey, ok := apiKeyData.(*entity.APIKeyResponse)
		if !ok {
			fmt.Printf("[API-Key-Auth] Invalid API key type in context. Got: %T\n", apiKeyData)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "INTERNAL_SERVER_ERROR",
					"message": "Invalid API key data type",
				},
			})
			c.Abort()
			return
		}

		fmt.Printf("[API-Key-Auth] Checking withdrawal permission - CanWithdrawal: %v\n", apiKey.CanWithdrawal)
		if !apiKey.CanWithdrawal {
			fmt.Printf("[API-Key-Auth] Withdrawal permission denied for API key: %s\n", maskAPIKey(apiKey.PublicKey))
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "FORBIDDEN",
					"message": "API key does not have withdrawal permission",
				},
			})
			c.Abort()
			return
		}

		fmt.Println("[API-Key-Auth] Withdrawal permission granted")
		c.Next()
	}
}

// RequirePaymentProcessingPermission middleware to check payment processing permission
func RequirePaymentProcessingPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("[API-Key-Auth] Checking payment processing permission for path: %s\n", c.Request.URL.Path)

		apiKeyData, exists := c.Get("apiKey")
		if !exists {
			fmt.Println("[API-Key-Auth] API key data not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "UNAUTHORIZED",
					"message": "API key data not found",
				},
			})
			c.Abort()
			return
		}

		apiKey, ok := apiKeyData.(*entity.APIKeyResponse)
		if !ok {
			fmt.Printf("[API-Key-Auth] Invalid API key type in context. Got: %T\n", apiKeyData)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "INTERNAL_SERVER_ERROR",
					"message": "Invalid API key data type",
				},
			})
			c.Abort()
			return
		}

		fmt.Printf("[API-Key-Auth] Checking payment processing permission - CanProcessPayment: %v\n", apiKey.CanProcessPayment)
		if !apiKey.CanProcessPayment {
			fmt.Printf("[API-Key-Auth] Payment processing permission denied for API key: %s\n", maskAPIKey(apiKey.PublicKey))
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "FORBIDDEN",
					"message": "API key does not have payment processing permission",
				},
			})
			c.Abort()
			return
		}

		fmt.Println("[API-Key-Auth] Payment processing permission granted")
		c.Next()
	}
}

// maskAPIKey masks the API key for safe logging
func maskAPIKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
