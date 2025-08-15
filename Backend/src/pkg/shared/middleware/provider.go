package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/usecase"
	authUseCase "github.com/socialpay/socialpay/src/pkg/auth/usecase"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	ipWhitelistUsecase "github.com/socialpay/socialpay/src/pkg/ip_whitelist/usecase"
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
	// IPChecker middleware for merchant routes
	IPChecker *ginMiddleware.IPCheckerMiddleware
	// Merchant ID middleware for merchant validation
	MerchantID gin.HandlerFunc
	// RBAC middleware for permission checking
	RBAC *ginMiddleware.RBACV2
}

// NewMiddlewareProvider creates a new middleware provider
func NewMiddlewareProvider(
	auth authUseCase.AuthNInteractor,
	api usecase.APIKeyUseCase,
	authService service.AuthService,
	ipWhitelistUsecase ipWhitelistUsecase.IPWhitelistUseCase,
) *MiddlewareProvider {
	// JWT auth middleware
	jwtAuth := ginMiddleware.JWTAuthMiddleware(ginMiddleware.JWTAuthMiddlewareConfig{
		Public:      false,
		AuthService: authService,
	})

	// Public middleware (no auth required)
	public := ginMiddleware.JWTAuthMiddleware(ginMiddleware.JWTAuthMiddlewareConfig{
		Public:      true,
		AuthService: authService,
	})

	// API key middleware
	apiKey := ginMiddleware.APIKeyAuth(api)

	// CORS middleware
	cors := ginMiddleware.CORSMiddleware(ginMiddleware.CORSConfig{})

	// Merchant ID middleware
	merchantID := ginMiddleware.MerchantIDMiddleware()

	rbac := ginMiddleware.NewRBACV2(authService)

	ipChecker := ginMiddleware.NewIPCheckerMiddleware(ipWhitelistUsecase)

	return &MiddlewareProvider{
		JWTAuth:    jwtAuth,
		Public:     public,
		APIKey:     apiKey,
		CORS:       cors,
		IPChecker:  ipChecker,
		MerchantID: merchantID,
		RBAC:       rbac,
	}
}
