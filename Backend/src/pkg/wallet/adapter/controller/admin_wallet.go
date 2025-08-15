package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/middleware"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
	walletUsecase "github.com/socialpay/socialpay/src/pkg/wallet/usecase"
)

type AdminWalletController struct {
	logger             logging.Logger
	usecase            walletUsecase.AdminWalletUsecase
	middlewareProvider *middleware.MiddlewareProvider
}

func NewAdminWalletController(
	usecase walletUsecase.AdminWalletUsecase,
	middlewareProvider *middleware.MiddlewareProvider,
) *AdminWalletController {
	return &AdminWalletController{
		logger:             logging.NewStdLogger("[adminWalletController]"),
		usecase:            usecase,
		middlewareProvider: middlewareProvider,
	}
}

func (c *AdminWalletController) RegisterRoutes(router *gin.RouterGroup) {
	adminGroup := router.Group("/admin", ginn.ErrorMiddleWare())

	adminGroup.GET("/wallet",
		c.middlewareProvider.JWTAuth,
		c.middlewareProvider.RBAC.RequirePermissionForMerchant(auth_entity.RESOURCE_WALLET, auth_entity.OPERATION_ADMIN_READ),
		c.GetAdminWallet)

	// Health check endpoint for wallet and transaction balance consistency
	adminGroup.GET("/health/wallet-balance",
		c.middlewareProvider.JWTAuth,
		c.middlewareProvider.RBAC.RequirePermissionForMerchant(auth_entity.RESOURCE_WALLET, auth_entity.OPERATION_ADMIN_READ),
		c.CheckWalletBalanceHealth)
}

// GetAdminWallet godoc
// @Summary      Get admin wallet
// @Description  Retrieves the admin wallet information including total amount
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{} "wallet: entity.MerchantWallet, total_amount: float64"
// @Failure      400 {object} map[string]string "error: error message"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /admin/wallet [get]
func (c *AdminWalletController) GetAdminWallet(ctx *gin.Context) {
	userID, exists := ginn.GetUserIDFromContext(ctx)
	if !exists {
		c.logger.Error("User not authenticated", nil)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	wallet, err := c.usecase.GetAdminWallet(ctx)
	if err != nil {
		c.logger.Error("failed to get admin wallet", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to get admin wallet",
			},
		})
		return
	}

	totalAmount, err := c.usecase.GetTotalAdminWalletAmount(ctx)
	if err != nil {
		c.logger.Error("failed to get total admin wallet amount", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to get total admin wallet amount",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"wallet":       wallet,
			"total_amount": totalAmount,
		},
	})
}

// CheckWalletBalanceHealth godoc
// @Summary      Check wallet balance health
// @Description  Verifies if wallet balances match transaction history (deposits - withdrawals)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{} "health_check: WalletHealthCheck"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /admin/health/wallet-balance [get]
func (c *AdminWalletController) CheckWalletBalanceHealth(ctx *gin.Context) {
	userID, exists := ginn.GetUserIDFromContext(ctx)
	if !exists {
		c.logger.Error("User not authenticated", nil)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	healthCheck, err := c.usecase.CheckWalletBalanceHealth(ctx, userID)
	if err != nil {
		c.logger.Error("failed to check wallet balance health", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to check wallet balance health",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"health_check": healthCheck,
		},
	})
}
