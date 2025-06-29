package gin

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds configuration for CORS middleware
type CORSConfig struct {
	AllowedOrigins []string
}

// CORSMiddleware creates a middleware that handles CORS for different authentication scenarios
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	// Get allowed origins from environment variable
	allowedOrigins := config.AllowedOrigins
	if len(allowedOrigins) == 0 {
		// Default to environment variable
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		if allowedOriginsEnv != "" {
			allowedOrigins = strings.Split(allowedOriginsEnv, ",")
			fmt.Printf("[CORS] Loaded origins from env: %v\n", allowedOrigins)
		} else {
			// Fallback to default
			allowedOrigins = []string{"http://localhost:3001", "http://localhost:3003"}
			fmt.Printf("[CORS] Using default origins: %v\n", allowedOrigins)
		}
	} else {
		fmt.Printf("[CORS] Using configured origins: %v\n", allowedOrigins)
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		fmt.Printf("[CORS] Request from origin: %s\n", origin)
		fmt.Printf("[CORS] Request method: %s\n", c.Request.Method)
		fmt.Printf("[CORS] Request path: %s\n", c.Request.URL.Path)

		// Check if the origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		fmt.Printf("[CORS] Origin allowed: %v\n", allowed)

		// Always set CORS headers for OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			fmt.Println("[CORS] Handling OPTIONS request")
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key, , x-merchant-id")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.AbortWithStatus(204)
			fmt.Println("[CORS] OPTIONS request handled, response headers set")
			return
		}

		if allowed {
			fmt.Println("[CORS] Setting CORS headers for allowed origin")
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key, X-MERCHANT-ID")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		c.Next()
	}
}
