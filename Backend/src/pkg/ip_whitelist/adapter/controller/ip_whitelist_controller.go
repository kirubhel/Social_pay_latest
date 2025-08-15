package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	auth_service "github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	"github.com/socialpay/socialpay/src/pkg/ip_whitelist/core/entity"
	ip_whitelist_usecase "github.com/socialpay/socialpay/src/pkg/ip_whitelist/usecase"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
)

type IPWhitelistController struct {
	authService auth_service.AuthService
	logger      logging.Logger
	rbac        *ginn.RBACV2
	usecase     ip_whitelist_usecase.IPWhitelistUseCase
}

// CreateIPWhitelistRequest contains create IP Whitelist request body
type CreateIPWhitelistRequest struct {
	IPAddress string `json:"ip_address"`
}

// UpdateIPWhitelistRequest contains create IP Whitelist request body
type UpdateIPWhitelistRequest struct {
	IPAddress string `json:"ip_address"`
	IsActive  bool   `json:"is_active"`
}

func NewIPWhitelistController(
	authService auth_service.AuthService,
	rbac *ginn.RBACV2,
	usecase ip_whitelist_usecase.IPWhitelistUseCase,
) *IPWhitelistController {
	return &IPWhitelistController{
		authService: authService,
		logger:      logging.NewStdLogger("[ipWhitelistController]"),
		rbac:        rbac,
		usecase:     usecase,
	}
}

func (c *IPWhitelistController) RegisterRoutes(router *gin.RouterGroup) {
	merchantGroup := router.Group("/merchant")

	jwtConfig := ginn.JWTAuthMiddlewareConfig{
		AuthService: c.authService,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Public:      false,
	}

	// IP whitelist routes
	merchantGroup.GET("/ip-whitelist",
		ginn.JWTAuthMiddleware(jwtConfig),
		c.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_IP_WHITELIST, auth_entity.OPERATION_READ),
		c.GetIPWhitelist)

	merchantGroup.POST("/ip-whitelist",
		ginn.JWTAuthMiddleware(jwtConfig),
		c.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_IP_WHITELIST, auth_entity.OPERATION_CREATE),
		c.CreateIPWhitelist)

	merchantGroup.PUT("/ip-whitelist/:id",
		ginn.JWTAuthMiddleware(jwtConfig),
		c.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_IP_WHITELIST, auth_entity.OPERATION_UPDATE),
		c.UpdateIPWhitelist)

	merchantGroup.DELETE("/ip-whitelist/:id",
		ginn.JWTAuthMiddleware(jwtConfig),
		c.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_IP_WHITELIST, auth_entity.OPERATION_DELETE),
		c.DeleteIPWhitelist)
}

// GetIPWhitelist godoc
// @Summary      Get merchant IP whitelist
// @Description  Retrieves the IP whitelist for the authenticated merchant
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "whitelist: entity.IPWhitelist"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /merchant/ip-whitelist [get]
func (c *IPWhitelistController) GetIPWhitelist(ctx *gin.Context) {
	merchantID, exists := ginn.GetMerchantIDFromContext(ctx)
	fmt.Println(merchantID, exists)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "UNAUTHORIZED",
				"message": "Merchant not authenticated",
			},
		})
		return
	}

	whitelist, err := c.usecase.GetIPWhitelist(ctx.Request.Context(), merchantID)
	if err != nil {
		c.logger.Error("failed to get IP whitelist", map[string]interface{}{
			"error":      err.Error(),
			"merchantID": merchantID,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to get IP whitelist",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    whitelist,
	})
}

// CreateIPWhitelist godoc
// @Summary      Create merchant IP whitelist
// @Description  Create the IP whitelist for the authenticated merchant
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateIPWhitelistRequest true "IP Whitelist Request"
// @Success      200 {object} map[string]interface{} "message: IP whitelist created successfully"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /merchant/ip-whitelist [post]
func (c *IPWhitelistController) CreateIPWhitelist(ctx *gin.Context) {
	merchantID, exists := ginn.GetMerchantIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "UNAUTHORIZED",
				"message": "Merchant not authenticated",
			},
		})
		return
	}

	var request entity.CreateIPWhitelistRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
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

	err := c.usecase.CreateIPWhitelist(ctx, merchantID, request)
	if err != nil {
		c.logger.Error("failed to whitelist IP address", map[string]interface{}{
			"error":      err.Error(),
			"merchantID": merchantID,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed whitelist an IP address: " + err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "IP Whitelisted successfully",
	})
}

// UpdateIPWhitelist godoc
// @Summary      Update merchant IP whitelist
// @Description  Updates the IP whitelist for the authenticated merchant
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateIPWhitelistRequest true "IP Whitelist Request"
// @Success      200 {object} map[string]interface{} "whitelist: entity.IPWhitelist"
// @Failure      400 {object} map[string]string "error: invalid request"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /merchant/ip-whitelist/:id [put]
func (c *IPWhitelistController) UpdateIPWhitelist(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		c.logger.Error("Invalid IP Whitelist ID", map[string]interface{}{
			"error": err.Error(),
			"id":    ctx.Param("id"),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid IP Whitelist ID",
			},
		})
		return
	}

	merchantID, exists := ginn.GetMerchantIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "UNAUTHORIZED",
				"message": "Merchant not authenticated",
			},
		})
		return
	}

	var request UpdateIPWhitelistRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
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

	err = c.usecase.UpdateIPWhitelist(ctx.Request.Context(), merchantID, entity.UpdateIPWhitelistRequest{
		ID:        id,
		IPAddress: request.IPAddress,
		IsActive:  request.IsActive,
	})
	if err != nil {
		c.logger.Error("failed to update IP whitelist", map[string]interface{}{
			"error":      err.Error(),
			"merchantID": merchantID,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to update IP whitelist",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "IP Whitelist successfully updated",
	})
}

// DeleteIPWhitelist godoc
// @Summary      Delete merchant IP whitelist
// @Description  Deletes the IP whitelist for the authenticated merchant
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "message: IP whitelist deleted successfully"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /merchant/ip-whitelist/:id [delete]
func (c *IPWhitelistController) DeleteIPWhitelist(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		c.logger.Error("Invalid IP Whitelist ID", map[string]interface{}{
			"error": err.Error(),
			"id":    ctx.Param("id"),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid IP Whitelist ID",
			},
		})
		return
	}

	merchantID, exists := ginn.GetMerchantIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "UNAUTHORIZED",
				"message": "Merchant not authenticated",
			},
		})
		return
	}

	if err := c.usecase.DeleteIPWhitelist(ctx.Request.Context(), id, merchantID); err != nil {
		c.logger.Error("failed to delete IP whitelist", map[string]interface{}{
			"error":      err.Error(),
			"merchantID": merchantID,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_SERVER_ERROR",
				"message": "Failed to delete IP whitelist",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "IP whitelist deleted successfully",
	})
}
