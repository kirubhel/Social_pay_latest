package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/usecase"
	authUseCase "github.com/socialpay/socialpay/src/pkg/auth/usecase"
	ginMiddleware "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
)

// MiddlewareProvider provides middleware functions for the application
type MiddlewareProvider struct {
	// The authentication middleware for protected routes
	JWTAuth gin.HandlerFunc
	// The authentication middleware for public routes
	Public gin.HandlerFunc
	// The authentication middleware by api key
	APIKey gin.HandlerFunc
	// CORS middleware for JWT authenticated routes
	CORS gin.HandlerFunc
	// Merchant ID middleware for merchant validation
	MerchantID gin.HandlerFunc
}

// NewMiddlewareProvider creates a new middleware provider
func NewMiddlewareProvider(
	auth authUseCase.AuthNInteractor,
	api usecase.APIKeyUseCase,
) *MiddlewareProvider {
	// JWT auth middleware
	jwtAuth := ginMiddleware.JWTAuthMiddleware(ginMiddleware.JWTAuthMiddlewareConfig{
		AuthUseCase: auth,
		Public:      false,
	})

	// Public middleware (no auth required)
	public := ginMiddleware.JWTAuthMiddleware(ginMiddleware.JWTAuthMiddlewareConfig{
		AuthUseCase: auth,
		Public:      true,
	})

	// API key middleware
	apiKey := ginMiddleware.APIKeyAuth(api)

	// CORS middleware
	cors := ginMiddleware.CORSMiddleware(ginMiddleware.CORSConfig{})

	// Merchant ID middleware
	merchantID := ginMiddleware.MerchantIDMiddleware()

	return &MiddlewareProvider{
		JWTAuth:    jwtAuth,
		Public:     public,
		APIKey:     apiKey,
		CORS:       cors,
		MerchantID: merchantID,
	}
}
