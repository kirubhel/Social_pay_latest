package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/qr/core/entity"
	"github.com/socialpay/socialpay/src/pkg/qr/usecase"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	ginMiddleware "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
)

type Handler struct {
	qrUseCase     usecase.QRUseCase
	log           logging.Logger
	jwtMiddleware gin.HandlerFunc
	rbac          *ginMiddleware.RBACV2
}

// SetupQRManagementRoutes sets up the QR management routes (protected by JWT and RBAC)
func (h *Handler) RegisterRouter(router *gin.RouterGroup) {
	qrMgmt := router.Group("/qr_mgmt", h.jwtMiddleware, ginMiddleware.MerchantIDMiddleware())
	// QR Link management endpoints with RBAC protection
	qrMgmt.POST("/links",
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_QR, auth_entity.OPERATION_CREATE),
		h.CreateQRLink)
	qrMgmt.GET("/links",
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_QR, auth_entity.OPERATION_READ),
		h.GetQRLinks)
	qrMgmt.GET("/links/:id",
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_QR, auth_entity.OPERATION_READ),
		h.GetQRLink)
	qrMgmt.PUT("/links/:id",
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_QR, auth_entity.OPERATION_UPDATE),
		h.UpdateQRLink)
	qrMgmt.DELETE("/links/:id",
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_QR, auth_entity.OPERATION_DELETE),
		h.DeleteQRLink)
}

func NewHandler(qrUseCase usecase.QRUseCase, jwtMiddleware gin.HandlerFunc, rbac *ginMiddleware.RBACV2) *Handler {
	return &Handler{
		qrUseCase:     qrUseCase,
		log:           logging.NewStdLogger("qr_handler"),
		jwtMiddleware: jwtMiddleware,
		rbac:          rbac,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func newErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:   "request_failed",
		Message: err.Error(),
	}
}

// CreateQRLink godoc
// @Summary      Create QR payment link
// @Description  Create a new QR payment link for merchant
// @Tags         QR-Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body entity.CreateQRLinkRequest true "QR link creation request"
// @Success      201  {object}  entity.QRLinkResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /qr_mgmt/links [post]
func (h *Handler) CreateQRLink(c *gin.Context) {
	var req entity.CreateQRLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Failed to bind JSON request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	// Get user ID from JWT context
	userID, exists := ginMiddleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, newErrorResponse(fmt.Errorf("user not authenticated")))
		return
	}

	merchantID, exists := ginMiddleware.GetMerchantIDFromContext(c) // This would be different in production
	if !exists {
		c.JSON(http.StatusUnauthorized, newErrorResponse(fmt.Errorf("merchant not authenticated")))
		return
	}

	response, err := h.qrUseCase.CreateQRLink(c.Request.Context(), userID, merchantID, &req)
	if err != nil {
		h.log.Error("Failed to create QR link", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetQRLinks godoc
// @Summary      Get QR payment links
// @Description  Get paginated list of QR payment links for authenticated user
// @Tags         QR-Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page      query     int  false  "Page number (default: 1)"
// @Param        page_size query     int  false  "Items per page (default: 10)"
// @Success      200  {object}  entity.QRLinksListResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /qr_mgmt/links [get]
func (h *Handler) GetQRLinks(c *gin.Context) {
	// Get pagination parameters using existing pagination package
	pag, err := pagination.NewPagination(c, h.log)
	if err != nil {
		h.log.Error("Failed to parse pagination parameters", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	// Get user ID from JWT context
	merchantID, exists := ginMiddleware.GetMerchantIDFromContext(c) // This would be different in production
	if !exists {
		c.JSON(http.StatusUnauthorized, newErrorResponse(fmt.Errorf("merchant not authenticated")))
		return
	}

	response, err := h.qrUseCase.GetQRLinksByMerchant(c.Request.Context(), merchantID, pag)
	if err != nil {
		h.log.Error("Failed to get QR links", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetQRLink godoc
// @Summary      Get QR payment link
// @Description  Get QR payment link details by ID
// @Tags         QR-Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "QR Link ID"
// @Success      200  {object}  entity.QRLinkResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /qr_mgmt/links/{id} [get]
func (h *Handler) GetQRLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("invalid QR link ID")))
		return
	}

	response, err := h.qrUseCase.GetQRLink(c.Request.Context(), id)
	if err != nil {
		h.log.Error("Failed to get QR link", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusNotFound, newErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateQRLink godoc
// @Summary      Update QR payment link
// @Description  Update an existing QR payment link
// @Tags         QR-Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                        true  "QR Link ID"
// @Param        request body      entity.UpdateQRLinkRequest    true  "QR link update request"
// @Success      200  {object}  entity.QRLinkResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /qr_mgmt/links/{id} [put]
func (h *Handler) UpdateQRLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("invalid QR link ID")))
		return
	}

	var req entity.UpdateQRLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Failed to bind JSON request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	// Get user ID from JWT context
	userID, exists := ginMiddleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, newErrorResponse(fmt.Errorf("user not authenticated")))
		return
	}

	response, err := h.qrUseCase.UpdateQRLink(c.Request.Context(), id, userID, &req)
	if err != nil {
		h.log.Error("Failed to update QR link", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteQRLink godoc
// @Summary      Delete QR payment link
// @Description  Delete (deactivate) a QR payment link
// @Tags         QR-Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "QR Link ID"
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /qr_mgmt/links/{id} [delete]
func (h *Handler) DeleteQRLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("invalid QR link ID")))
		return
	}

	// Get user ID from JWT context
	userID, exists := ginMiddleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, newErrorResponse(fmt.Errorf("user not authenticated")))
		return
	}

	err = h.qrUseCase.DeleteQRLink(c.Request.Context(), id, userID)
	if err != nil {
		h.log.Error("Failed to delete QR link", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}
