package gin

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CheckoutCors is a CORS middleware that only allows requests from checkout pages
func CheckoutCors(c *gin.Context) {
	origin := c.GetHeader("Origin")

	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
	} else {
		allowedOrigins = []string{
			"http://localhost:3000",
			"https://localhost:3000",
			"https://checkout.Socialpay.co",
			"http://checkout.Socialpay.co", // for development
		}
	}
	log.Println("allowedOrgins:", allowedOrigins)
	// Check if the origin is allowed
	isAllowed := false
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			isAllowed = true
			break
		}
	}

	if isAllowed {
		c.Header("Access-Control-Allow-Origin", origin)
	}

	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Max-Age", "86400")

	// Handle preflight requests
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	// If origin is not allowed, return 403
	if !isAllowed && origin != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "FORBIDDEN",
				"message": "Origin not allowed",
				"origin":  origin,
			},
		})
		c.Abort()
		return
	}

	c.Next()
}
