package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"

	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	ginMiddleware "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
	"github.com/socialpay/socialpay/src/pkg/shared/response"
	"github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
	webhookEntity "github.com/socialpay/socialpay/src/pkg/webhook/core/entity"
	"github.com/socialpay/socialpay/src/pkg/webhook/core/validation"
	usecase "github.com/socialpay/socialpay/src/pkg/webhook/usecase"
)

type WebhookController struct {
	logger  logging.Logger
	usecase usecase.WebhookUseCase
	jwtAuth gin.HandlerFunc
	rbac    *ginMiddleware.RBACV2
}

func NewWebhookController(usecase usecase.WebhookUseCase, jwtAuth gin.HandlerFunc, rbac *ginMiddleware.RBACV2) *WebhookController {
	return &WebhookController{
		logger:  logging.NewStdLogger("[webhookController]"),
		usecase: usecase,
		jwtAuth: jwtAuth,
		rbac:    rbac,
	}
}

func (c *WebhookController) RegisterRoutes(router *gin.RouterGroup) {
	webhookGroup := router.Group("/webhooks", ginMiddleware.ErrorMiddleWare(), c.jwtAuth)
	webhookGroup.POST("/callback", c.HandleWebhook)
	webhookGroup.GET("/callback/:id", c.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_WEBHOOK, auth_entity.OPERATION_READ), c.GetCallbackLogByID)
	webhookGroup.GET("/callback/merchant", ginMiddleware.MerchantIDMiddleware(), c.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_WEBHOOK, auth_entity.OPERATION_READ), c.GetCallbackLogsByMerchantID)
	webhookGroup.GET("/callback", c.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_WEBHOOK, auth_entity.OPERATION_READ), c.GetAllCallbackLogs)
}

// HandleWebhook godoc
// @Summary      Receive webhook callback
// @Description  Receives a webhook callback and produces an event to Kafka
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        request body dto.WebhookRequest true "Webhook callback payload"
// @Success      200 {object} map[string]string "status: success"
// @Failure      400 {object} map[string]string "error: error message"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /webhooks/callback [post]
func (c *WebhookController) HandleWebhook(ctx *gin.Context) {
	var req dto.WebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validation.IsValidStatus(req.Status) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	err := c.usecase.HandleWebhookDispatch(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetCallbackLogByID godoc
// @Summary      Get callback log by ID
// @Description  Retrieves a specific callback log by its ID
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id path string true "Callback Log ID"
// @Success      200 {object} webhookEntity.CallbackLog
// @Failure      400 {object} map[string]string "error: error message"
// @Failure      404 {object} map[string]string "error: error message"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /webhooks/callback/{id} [get]
func (c *WebhookController) GetCallbackLogByID(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid callback log ID"})
		return
	}
	var log *webhookEntity.CallbackLog
	log, err = c.usecase.GetCallbackLogByID(ctx, parsedID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if log == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "callback log not found"})
		return
	}

	ctx.JSON(http.StatusOK, log)
}

// GetCallbackLogsByMerchantID godoc
// @Summary      Get callback logs for authenticated merchant
// @Description  Retrieves all callback logs for the merchant associated with the API key. The merchant ID is automatically determined from the API key context.
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        page query int true "Page number (min: 1)"
// @Param        page_size query int true "Number of items per page (min: 1, max: 100)"
// @Security     ApiKeyAuth
// @Success      200 {object} response.PaginatedResponse
// @Failure      400 {object} map[string]string "error: invalid merchant ID"
// @Failure      401 {object} map[string]string "error: unauthorized"
// @Failure      500 {object} map[string]string "error: internal server error"
// @Router       /webhooks/callback/merchant [get]
func (c *WebhookController) GetCallbackLogsByMerchantID(ctx *gin.Context) {
	merchantID, exists := ginMiddleware.GetMerchantIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "merchant ID not found in context"})
		return
	}

	// Use shared pagination like transaction module
	pag, err := pagination.NewPagination(ctx, c.logger)
	if err != nil {
		c.logger.Error("pagination binding error", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination parameters"})
		return
	}

	// Convert to transaction pagination for usecase
	txPagination := &txEntity.Pagination{
		Page:     pag.Page,
		PageSize: pag.PerPage,
	}

	logs, err := c.usecase.GetCallbackLogsByMerchantID(ctx.Request.Context(), merchantID, txPagination)
	if err != nil {
		c.logger.Error("failed to get callback logs by merchant ID", map[string]interface{}{
			"error":      err.Error(),
			"merchantID": merchantID,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.PaginatedResponse{
		Success:    true,
		Data:       logs,
		Pagination: pag.GetInfo(len(logs)),
	})
}

// GetAllCallbackLogs godoc
// @Summary      Get all callback logs
// @Description  Retrieves all callback logs ordered by creation date
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        page query int true "Page number (min: 1)"
// @Param        page_size query int true "Number of items per page (min: 1, max: 100)"
// @Success      200 {object} response.PaginatedResponse
// @Failure      400 {object} map[string]string "error: error message"
// @Failure      500 {object} map[string]string "error: error message"
// @Router       /webhooks/callback [get]
func (c *WebhookController) GetAllCallbackLogs(ctx *gin.Context) {
	// Use shared pagination like transaction module
	pag, err := pagination.NewPagination(ctx, c.logger)
	if err != nil {
		c.logger.Error("pagination binding error", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination parameters"})
		return
	}

	// Convert to transaction pagination for usecase
	txPagination := &txEntity.Pagination{
		Page:     pag.Page,
		PageSize: pag.PerPage,
	}

	logs, err := c.usecase.GetAllCallbackLogs(ctx.Request.Context(), txPagination)
	if err != nil {
		c.logger.Error("failed to get all callback logs", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.PaginatedResponse{
		Success:    true,
		Data:       logs,
		Pagination: pag.GetInfo(len(logs)),
	})
}
