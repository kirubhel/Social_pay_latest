package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/commission/core/entity"
	commission_usecase "github.com/socialpay/socialpay/src/pkg/commission/usecase"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/middleware"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
)

type CommissionController struct {
	logger             logging.Logger
	usecase            commission_usecase.CommissionUseCase
	middlewareProvider *middleware.MiddlewareProvider
}

// CommissionResponse represents the standard API response structure
type CommissionResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}

// CommissionError represents the standard API error structure
type CommissionError struct {
	Type    string `json:"type" example:"INTERNAL_SERVER_ERROR"`
	Message string `json:"message" example:"Failed to get default commission settings"`
}

// CommissionErrorResponse represents the standard API error response structure
type CommissionErrorResponse struct {
	Success bool            `json:"success" example:"false"`
	Error   CommissionError `json:"error"`
}

func NewCommissionController(
	usecase commission_usecase.CommissionUseCase,
	middlewareProvider *middleware.MiddlewareProvider,
) *CommissionController {
	return &CommissionController{
		logger:             logging.NewStdLogger("[commissionController]"),
		usecase:            usecase,
		middlewareProvider: middlewareProvider,
	}
}

func (c *CommissionController) RegisterRoutes(router *gin.RouterGroup) {
	c.logger.Info("Registering commission routes", map[string]interface{}{
		"base_path": router.BasePath(),
	})

	adminGroup := router.Group("/admin", ginn.ErrorMiddleWare())
	c.logger.Info("Created admin group", map[string]interface{}{
		"admin_group_path": adminGroup.BasePath(),
	})

	// Default commission routes
	adminGroup.GET("/commissions/default",
		c.middlewareProvider.JWTAuth,
		c.middlewareProvider.RBAC.RequirePermissionForAdmin(auth_entity.RESOURCE_COMMISSION, auth_entity.OPERATION_ADMIN_READ),
		c.GetDefaultCommission)

	c.logger.Info("Registered GET /commissions/default route", map[string]interface{}{
		"full_path": adminGroup.BasePath() + "/commissions/default",
	})

	adminGroup.PUT("/commissions/default",
		c.middlewareProvider.JWTAuth,
		c.middlewareProvider.RBAC.RequirePermissionForAdmin(auth_entity.RESOURCE_COMMISSION, auth_entity.OPERATION_ADMIN_UPDATE),
		c.UpdateDefaultCommission)

	// Merchant commission routes
	adminGroup.GET("/commissions/merchant/:merchantID",
		c.middlewareProvider.JWTAuth,
		c.middlewareProvider.RBAC.RequirePermissionForAdmin(auth_entity.RESOURCE_COMMISSION, auth_entity.OPERATION_ADMIN_READ),
		c.GetMerchantCommission)

	adminGroup.PUT("/commissions/merchant/:merchantID",
		c.middlewareProvider.JWTAuth,
		c.middlewareProvider.RBAC.RequirePermissionForAdmin(auth_entity.RESOURCE_COMMISSION, auth_entity.OPERATION_ADMIN_UPDATE),
		c.UpdateMerchantCommission)
}

// GetDefaultCommission godoc
// @Summary      Get default commission settings
// @Description  Retrieves the default commission settings for all merchants
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "settings: entity.CommissionSettings"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      403 {object} map[string]string "error: forbidden"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /admin/commissions/default [get]
func (c *CommissionController) GetDefaultCommission(ctx *gin.Context) {
	c.logger.Info("Received request for default commission", map[string]interface{}{
		"path":   ctx.Request.URL.Path,
		"method": ctx.Request.Method,
	})

	settings, err := c.usecase.GetDefaultCommission(ctx.Request.Context())
	if err != nil {
		c.logger.Error("failed to get default commission settings", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to get default commission settings",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// UpdateDefaultCommission godoc
// @Summary      Update default commission settings
// @Description  Updates the default commission settings for all merchants
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        settings body entity.CommissionSettings true "Commission Settings"
// @Success      200 {object} map[string]interface{} "settings: entity.CommissionSettings"
// @Failure      400 {object} map[string]string "error: invalid request"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      403 {object} map[string]string "error: forbidden"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /admin/commissions/default [put]
func (c *CommissionController) UpdateDefaultCommission(ctx *gin.Context) {
	var settings entity.CommissionSettings
	if err := ctx.ShouldBindJSON(&settings); err != nil {
		c.logger.Error("invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
		return
	}

	if err := c.usecase.UpdateDefaultCommission(ctx.Request.Context(), &settings); err != nil {
		c.logger.Error("failed to update default commission settings", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to update default commission settings",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// GetMerchantCommission godoc
// @Summary      Get merchant commission settings
// @Description  Retrieves the commission settings for a specific merchant
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        merchantID path string true "Merchant ID" format(uuid)
// @Success      200 {object} map[string]interface{} "commission: entity.MerchantCommission"
// @Failure      400 {object} map[string]string "error: invalid merchant ID"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      403 {object} map[string]string "error: forbidden"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /admin/commissions/merchant/{merchantID} [get]
func (c *CommissionController) GetMerchantCommission(ctx *gin.Context) {
	merchantID, err := uuid.Parse(ctx.Param("merchantID"))
	if err != nil {
		c.logger.Error("invalid merchant ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid merchant ID",
			},
		})
		return
	}

	commission, err := c.usecase.GetMerchantCommission(ctx.Request.Context(), merchantID)
	if err != nil {
		c.logger.Error("failed to get merchant commission settings", map[string]interface{}{
			"error":      err.Error(),
			"merchantID": merchantID,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to get merchant commission settings",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    commission,
	})
}

// UpdateMerchantCommission godoc
// @Summary      Update merchant commission settings
// @Description  Updates the commission settings for a specific merchant
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        merchantID path string true "Merchant ID" format(uuid)
// @Param        commission body entity.MerchantCommission true "Merchant Commission Settings"
// @Success      200 {object} map[string]interface{} "commission: entity.MerchantCommission"
// @Failure      400 {object} map[string]string "error: invalid request"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      403 {object} map[string]string "error: forbidden"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /admin/commissions/merchant/{merchantID} [put]
func (c *CommissionController) UpdateMerchantCommission(ctx *gin.Context) {
	merchantID, err := uuid.Parse(ctx.Param("merchantID"))
	if err != nil {
		c.logger.Error("invalid merchant ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid merchant ID",
			},
		})
		return
	}

	var commission entity.MerchantCommission
	if err := ctx.ShouldBindJSON(&commission); err != nil {
		c.logger.Error("invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
		return
	}

	commission.MerchantID = merchantID

	if err := c.usecase.UpdateMerchantCommission(ctx.Request.Context(), merchantID, &commission); err != nil {
		c.logger.Error("failed to update merchant commission settings", map[string]interface{}{
			"error":      err.Error(),
			"merchantID": merchantID,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to update merchant commission settings",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    commission,
	})
}
