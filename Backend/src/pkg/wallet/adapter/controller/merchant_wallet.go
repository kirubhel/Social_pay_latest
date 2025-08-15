package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	ginMiddleware "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
	walletUseCase "github.com/socialpay/socialpay/src/pkg/wallet/usecase"
)

// ErrorResponse represents an error response
// @Description Error response
type ErrorResponse struct {
	// @Description Error message
	Error string `json:"error" example:"error message"`
}

// WalletResponse represents the wallet response
// @Description Wallet information response
type WalletResponse struct {
	// @Description Unique identifier for the wallet
	ID string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description ID of the merchant who owns this wallet
	MerchantID string `json:"merchant_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Current balance in the wallet
	Balance float64 `json:"balance" example:"1000.50"`
	// @Description Currency code
	Currency string `json:"currency" example:"ETB"`
	// @Description When the wallet was created
	CreatedAt string `json:"created_at" example:"2024-05-22T09:00:00Z"`
	// @Description When the wallet was last updated
	UpdatedAt string `json:"updated_at" example:"2024-05-22T09:00:00Z"`
	// @Description When the wallet was last synchronized
	LastSyncAt string `json:"last_sync_at" example:"2024-05-22T09:00:00Z"`
	// @Description Whether the wallet is active
	IsActive bool `json:"is_active" example:"true"`
	// @Description Optional description of the wallet
	Description string `json:"description,omitempty" example:"Main merchant wallet"`
}

// WalletController handles wallet-related HTTP requests
type WalletController struct {
	logger        logging.Logger
	walletUseCase walletUseCase.MerchantWalletUsecase
	middleware    *gin.HandlerFunc
	rbac          *ginMiddleware.RBACV2
}

// NewWalletController creates a new instance of WalletController
func NewWalletController(walletUseCase walletUseCase.MerchantWalletUsecase, middleware gin.HandlerFunc, rbac *ginMiddleware.RBACV2) *WalletController {
	fmt.Print("Initializing WalletController...\n")
	return &WalletController{
		logger:        logging.NewStdLogger("[walletController]"),
		walletUseCase: walletUseCase,
		middleware:    &middleware,
		rbac:          rbac,
	}
}

// RegisterRoutes registers the wallet routes
func (c *WalletController) RegisterRoutes(router *gin.RouterGroup) {
	wallet := router.Group("/wallet", ginMiddleware.ErrorMiddleWare(), *c.middleware, ginMiddleware.MerchantIDMiddleware())
	{
		wallet.GET("", c.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_WALLET, auth_entity.OPERATION_READ), c.GetMerchantWallet)
	}
}

// GetMerchantWallet godoc
// @Summary Get merchant wallet
// @Description Get the wallet information for the authenticated merchant
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security MerchantID
// @Success 200 {object} WalletResponse "Wallet information"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Wallet Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /wallet [get]
func (c *WalletController) GetMerchantWallet(ctx *gin.Context) {
	merchantID, exists := ginMiddleware.GetMerchantIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "merchant ID not found in context"})
		return
	}

	wallet, err := c.walletUseCase.GetMerchantWallet(ctx, merchantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if wallet == nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "wallet not found"})
		return
	}

	ctx.JSON(http.StatusOK, wallet)
}
