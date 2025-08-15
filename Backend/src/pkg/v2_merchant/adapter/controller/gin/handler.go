package gin

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	auth_service "github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/repository"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler manages merchant HTTP requests
type Handler struct {
	authService auth_service.AuthService
	log         logging.Logger
	useCase     usecase.MerchantUseCase
	rbac        *ginn.RBACV2
	repository  repository.Repository
}

// NewHandler creates a new merchant handler
func NewHandler(
	authService auth_service.AuthService,
	useCase usecase.MerchantUseCase,
	rbac *ginn.RBACV2,
	repository repository.Repository,
) *Handler {
	return &Handler{
		authService: authService,
		log:         logging.NewStdLogger("[V2_MERCHANT] [HANDLER]"),
		useCase:     useCase,
		rbac:        rbac,
		repository:  repository,
	}
}

// AddMerchantRequest represents the add merchant request payload
type AddMerchantRequest struct {
	Title        string `json:"title" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	PhonePrefix  string `json:"phone_prefix" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	Password     string `json:"password" binding:"required"`
	PasswordHint string `json:"password_hint,omitempty"`
}

// ExportMerchantsRequest contains export request body
type ExportMerchantsRequest struct {
	FileType  entity.SupportedFileType `json:"fileType" binding:"required"`
	Data      []string                 `json:"data" binding:"required"`
	Merchants []uuid.UUID              `json:"merchants" binding:"required"`
}

// UpdateMerchantRequest contains update merchant request body
type UpdateMerchantRequest struct {
	BusinessInfo entity.UpdateMerchantBusinessInformationRequest `json:"business_info"`
	PersonalInfo entity.CreateMerchantPersonalInformationRequest `json:"personal_info"`
	Documents    []entity.CreateMerchantDocumnetRequest          `json:"documents"`
}

// UpdateMerchantContactRequest contains update merchant document request body
type UpdateMerchantContactRequest struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
}

// UpdateMerchantDocumentRequest contains update merchant document request body
type UpdateMerchantDocumentRequest struct {
	FileUrl         string  `json:"file_url"`
	Status          string  `json:"status"`
	RejectionReason *string `json:"rejection_reason"`
}

// UpdateAdminMerchantRequest contains update merchant request body for admin
type UpdateAdminMerchantRequest struct {
	BusinessInfo entity.UpdateMerchantBusinessInformationRequest `json:"business_info"`
	PersonalInfo entity.CreateMerchantPersonalInformationRequest `json:"personal_info"`
	Documents    []entity.UpdateMerchantDocumentWithIDRequest    `json:"documents"`
}

// UpdateMerchantStatusRequest contains update merchant status request body
type UpdateMerchantStatusRequest struct {
	Status string `json:"status"`
}

// DeleteMerchantsRequest contains delete merchants request body
type DeleteMerchantsRequest struct {
	IDs []uuid.UUID `json:"ids"`
}

// ImpersonateMerchantRequest represents impersonate merchant request
type ImpersonateMerchantRequest struct {
	MerchantID uuid.UUID `json:"merchant_id"`
}

// RegisterRoutes registers merchant management routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	// Public routes - no authentication required as requested
	merchants := r.Group("/merchants")
	{
		merchants.GET("/:id", h.GetMerchant)
		merchants.GET("/:id/details", h.GetMerchantDetails)
	}

	jwtConfig := ginn.JWTAuthMiddlewareConfig{
		AuthService: h.authService,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Public:      false,
	}

	// Private routes
	{
		merchants.POST("/admin/add",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_ADMIN_CREATE),
			h.AddMerchant)

		merchants.GET("/admin/all",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_ADMIN_READ),
			h.GetMerchants)

		merchants.GET("/admin/stats",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_ADMIN_READ),
			h.GetMerchantStats)

		merchants.POST("/admin/export",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_EXPORT),
			h.ExportMerchants)

		merchants.PUT("/update",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequireMerchantOwner(),
			h.UpdateMerchant)

		merchants.PUT("/admin/status/update/:id",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_ADMIN_UPDATE),
			h.UpdateMerchantStatus)

		merchants.PUT("/admin/contact/update/:id",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_ADMIN_UPDATE),
			h.UpdateMerchantContact)

		merchants.PUT("/admin/document/update/:id",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_ADMIN_UPDATE),
			h.UpdateMerchantDocument)

		merchants.PUT("/admin/update/:id",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_ADMIN_UPDATE),
			h.UpdateAdminMerchant)

		merchants.DELETE("/admin/delete/:id",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_ADMIN_DELETE),
			h.DeleteMerchant)

		merchants.DELETE("/admin/delete/all",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_MERCHANT, auth_entity.OPERATION_ADMIN_DELETE),
			h.DeleteMerchants)

		merchants.POST("/admin/impersonate",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.rbac.RequireAdminAccess(),
			h.ImpersonateMerchant)
	}

}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   ApiError `json:"error"`
}

// ApiError represents API error details
type ApiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// AddMerchant handles add merchant
// @Summary Adds new merchant
// @Description Adds new with phone verification
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param request body AddMerchantRequest true "AddMerchant details"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 409 {object} map[string]interface{} "Phone already exists"
// AddMerchant godoc
// @Summary Add new merchant
// @Description Add a new merchant to the system (admin only)
// @Tags Merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddMerchantRequest true "Merchant creation request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /merchants/admin/add [post]
func (h *Handler) AddMerchant(c *gin.Context) {
	var req AddMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    auth_entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	// Get device info from request
	deviceInfo := &auth_entity.DeviceInfo{
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		DeviceName: "web", // Default for web requests
	}

	// Convert to entity request
	entityReq := &auth_entity.CreateUserRequest{
		Title:        req.Title,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhonePrefix:  req.PhonePrefix,
		PhoneNumber:  req.PhoneNumber,
		Password:     req.Password,
		PasswordHint: req.PasswordHint,
		UserType:     auth_entity.USER_TYPE_MERCHANT,
		DeviceInfo:   deviceInfo,
	}

	// Add merchant
	err := h.useCase.AddMerchant(c.Request.Context(), entityReq)
	if err != nil {
		if authErr, ok := err.(*auth_entity.AuthError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == auth_entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    auth_entity.ErrInternalServer,
					"message": "Add merchant failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchant added successfully",
	})
}

// GetMerchant godoc
// @Summary Get merchant by ID
// @Description Get a merchant by its ID (public endpoint)
// @Tags Merchants
// @Produce json
// @Param id path string true "Merchant ID" format(uuid)
// @Success 200 {object} SuccessResponse{data=entity.MerchantResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /merchants/{id} [get]
func (h *Handler) GetMerchant(c *gin.Context) {
	// Parse merchant ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid merchant ID", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid merchant ID",
			},
		})
		return
	}

	// Get merchant
	merchant, err := h.useCase.GetMerchant(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "merchant not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "NOT_FOUND",
					Message: "Merchant not found",
				},
			})
			return
		}

		h.log.Error("Failed to get merchant", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to get merchant",
			},
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    merchant,
	})
}

// GetMerchants godoc
// @Summary Get merchants
// @Description Get a list of merchants
// @Tags v2-merchants
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /merchants/all [get]
func (h *Handler) GetMerchants(c *gin.Context) {
	// Parse query parameters
	var skip int
	var take int
	var err error

	text := c.Query("text")
	status := c.Query("status")
	querySkip := c.Query("skip")
	queryTake := c.Query("take")
	queryStartDate := c.Query("startDate")
	queryEndDate := c.Query("endDate")

	if querySkip == "" {
		skip = 0
	} else {
		skip, err = strconv.Atoi(querySkip)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INVALID_REQUEST",
					Message: "Invalid skip query value",
				},
			})
			return
		}
	}

	if queryTake == "" {
		take = 10
	} else {
		take, err = strconv.Atoi(queryTake)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INVALID_REQUEST",
					Message: "Invalid take query value",
				},
			})
			return
		}
	}

	var merchants *entity.MerchantsResponse
	if queryStartDate == "" && queryEndDate == "" {
		// Get merchants
		merchants, err = h.useCase.GetMerchants(c.Request.Context(), entity.GetMerchantsParams{
			Text:   text,
			Skip:   skip,
			Take:   take,
			Status: status,
		})
	} else if queryStartDate == "" && queryEndDate != "" {
		// Parse end date query
		endDate, err := time.Parse("2006-01-02", queryEndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INVALID_REQUEST",
					Message: "Invalid end date query value",
				},
			})
			return
		}

		merchants, err = h.useCase.GetMerchants(c.Request.Context(), entity.GetMerchantsParams{
			Text:    text,
			Skip:    skip,
			Take:    take,
			Status:  status,
			EndDate: endDate,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INTERNAL_SERVER_ERROR",
					Message: "Failed to get merchants",
				},
			})
			return
		}
	} else if queryStartDate != "" && queryEndDate == "" {
		// Parse end date query
		startDate, err := time.Parse("2006-01-02", queryStartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INVALID_REQUEST",
					Message: "Invalid start date query value",
				},
			})
			return
		}

		merchants, err = h.useCase.GetMerchants(c.Request.Context(), entity.GetMerchantsParams{
			Text:      text,
			Skip:      skip,
			Take:      take,
			Status:    status,
			StartDate: startDate,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INTERNAL_SERVER_ERROR",
					Message: "Failed to get merchants",
				},
			})
			return
		}
	}

	if err != nil {
		if err.Error() == "merchant not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "NOT_FOUND",
					Message: "Merchant not found",
				},
			})
			return
		}

		h.log.Error("Failed to get merchant", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to get merchants",
			},
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    merchants,
	})
}

// GetMerchantDetails godoc
// @Summary Get merchant details by ID
// @Description Get complete merchant information with related data by ID
// @Tags v2-merchants
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /merchants/:id/details [get]
func (h *Handler) GetMerchantDetails(c *gin.Context) {
	// Parse merchant ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid merchant ID", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid merchant ID",
			},
		})
		return
	}

	// Get merchant details
	details, err := h.useCase.GetMerchantDetails(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "merchant not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "NOT_FOUND",
					Message: "Merchant not found",
				},
			})
			return
		}

		h.log.Error("Failed to get merchant details", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to get merchant details",
			},
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    details,
	})
}

// ExportMerchants handles export merchants
// @Summary Export merchants
// @Description Export list of merchants
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param request body ExportMerchantsRequest true "ExportMerchants details"
// @Success 201 {object} map[string]interface{} "Merchants exported successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 409 {object} map[string]interface{} "Phone already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/admin/export [post]
func (h *Handler) ExportMerchants(c *gin.Context) {
	var req ExportMerchantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    auth_entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	exportMerchants := &entity.ExportMerchantsRequest{
		FileType:  req.FileType,
		Data:      req.Data,
		Merchants: req.Merchants,
	}

	filePath, err := h.useCase.ExportMerchants(c, exportMerchants)
	if err != nil {
		h.log.Error("Failed to export merchants", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to export merchants",
			},
		})
		return
	}

	var fileName string

	splittedFilePath := strings.Split(*filePath, "/")
	fileName = splittedFilePath[1]
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", fileName))
	c.FileAttachment(*filePath, fileName)
}

// UpdateMerchant handles update merchant
// @Summary Update merchant
// @Description Update merchant informations
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param request body UpdateMerchantRequest true "UpdateMerchant details"
// @Success 201 {object} map[string]interface{} "Merchants updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/update [put]
func (h *Handler) UpdateMerchant(c *gin.Context) {
	merchantId, exists := ginn.GetMerchantIDFromContext(c)

	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	var req entity.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    auth_entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	err := h.useCase.UpdateMerchant(c, merchantId, &req)

	if err != nil {
		if err.Error() == "merchant not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "NOT_FOUND",
					Message: "Merchant not found",
				},
			})
			return
		}

		h.log.Error("Failed to get merchant", map[string]interface{}{
			"error": err.Error(),
			"id":    merchantId,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to update merchant",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchant updated successfully",
	})
}

// UpdateMerchantStatus handles update merchant status
// @Summary Update merchant status
// @Description Update merchant status
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param request body UpdateMerchantStatusRequest true "UpdateMerchantStatus details"
// @Success 201 {object} map[string]interface{} "Merchant status updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/admin/status/update/:id [put]
func (h *Handler) UpdateMerchantStatus(c *gin.Context) {
	// Parse merchant ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid ID", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid ID",
			},
		})
		return
	}

	var req entity.UpdateMerchantStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    auth_entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	err = h.useCase.UpdateMerchantStatus(c, id, &req)

	if err != nil {
		if err.Error() == "merchant not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "NOT_FOUND",
					Message: "Merchant not found",
				},
			})
			return
		}

		h.log.Error("Failed to get merchant", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to update merchant status",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchant status updated successfully",
	})
}

// UpdateMerchantContact handles update merchant contact
// @Summary Update merchant contact
// @Description Update merchant contact informations
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param request body UpdateMerchantContactRequest true "UpdateMerchantContact details"
// @Success 201 {object} map[string]interface{} "Merchant contact updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/admin/contact/update/:id [put]
func (h *Handler) UpdateMerchantContact(c *gin.Context) {
	// Parse merchant ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid ID", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid ID",
			},
		})
		return
	}

	var req entity.UpdateMerchantContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    auth_entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	err = h.useCase.UpdateMerchantContact(c, id, &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to update merchant contact",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchant contact updated successfully",
	})
}

// UpdateMerchantDocument handles update merchant document
// @Summary Update merchant document
// @Description Update merchant document informations
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param request body UpdateMerchantDocumentRequest true "UpdateMerchantContact details"
// @Success 201 {object} map[string]interface{} "Merchant document updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/admin/document/update/:id [put]
func (h *Handler) UpdateMerchantDocument(c *gin.Context) {
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User context not found",
			},
		})
		c.Abort()
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid ID", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid ID",
			},
		})
		return
	}

	var req entity.UpdateMerchantDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    auth_entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	err = h.useCase.UpdateMerchantDocument(c, id, &entity.UpdateMerchantDocumentRequest{
		FileUrl:         req.FileUrl,
		VerifiedBy:      &user.ID,
		RejectionReason: req.RejectionReason,
		Status:          req.Status,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to update merchant document",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchant document updated successfully",
	})
}

// UpdateAdminMerchant handles update merchant by admin
// @Summary Update merchant by admin
// @Description Update merchant information, contact information and documents by admin
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param id path string true "Merchant ID"
// @Param request body UpdateAdminMerchantRequest true "UpdateAdminMerchant details"
// @Success 200 {object} map[string]interface{} "Merchant updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Merchant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/admin/update/:id [put]
func (h *Handler) UpdateAdminMerchant(c *gin.Context) {
	// Parse merchant ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid merchant ID", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid merchant ID",
			},
		})
		return
	}

	var req UpdateAdminMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    auth_entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	// Convert to entity request
	entityReq := &entity.UpdateMerchantRequest{
		BusinessInfo: req.BusinessInfo,
		PersonalInfo: req.PersonalInfo,
		Documents:    req.Documents,
	}

	err = h.useCase.UpdateAdminMerchant(c, id, entityReq)

	if err != nil {
		if err.Error() == "merchant not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "NOT_FOUND",
					Message: "Merchant not found",
				},
			})
			return
		}

		h.log.Error("Failed to update merchant", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to update merchant",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchant updated successfully",
	})
}

// DeleteMerchant handles delete merchant by admin
// @Summary Delete merchant by admin
// @Description Admin deletes merchant
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 200 {object} map[string]interface{} "Merchant deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Merchant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/admin/delete/:id [delete]
func (h *Handler) DeleteMerchant(c *gin.Context) {
	// Parse merchant ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid merchant ID", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid merchant ID",
			},
		})
		return
	}

	err = h.useCase.DeleteMerchant(c, id)
	if err != nil {
		h.log.Error("Failed to delete merchant", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to delete merchant",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchant deleted successfully",
	})
}

// DeleteMerchants handles delete merchants by admin
// @Summary Delete merchants by admin
// @Description Admin deletes list of merchants
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param request body DeleteMerchantsRequest true "DeleteMerchants details"
// @Success 200 {object} map[string]interface{} "Merchantd deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Merchant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/admin/delete/all [delete]
func (h *Handler) DeleteMerchants(c *gin.Context) {
	var req entity.DeleteMerchantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    auth_entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	err := h.useCase.DeleteMerchants(c, &req)
	if err != nil {
		h.log.Error("Failed to delete merchants", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to delete merchants",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchants deleted successfully",
	})
}

// GetMerchantStats handles get merchant statistics
// @Summary Get merchant statistics
// @Description Get merchant statistics for admin dashboard
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Merchant statistics"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/admin/stats [get]
func (h *Handler) GetMerchantStats(c *gin.Context) {
	stats, err := h.useCase.GetMerchantStats(c)
	if err != nil {
		h.log.Error("Failed to get merchant stats", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to get merchant statistics",
			},
		})
		return
	}

	// Calculate percentage for active merchants
	var activePercentage float64
	if stats.TotalMerchants > 0 {
		activePercentage = float64(stats.ActiveMerchants) / float64(stats.TotalMerchants) * 100
	}

	// Format the response to match the UI requirements
	response := gin.H{
		"success": true,
		"data": gin.H{
			"total_merchants": gin.H{
				"title":      "Total Merchants",
				"value":      stats.TotalMerchants,
				"sub_text":   fmt.Sprintf("+%d new this month", stats.NewThisMonth),
				"icon":       "grid",
				"icon_color": "gray",
			},
			"active_merchants": gin.H{
				"title":      "Active Merchants",
				"value":      stats.ActiveMerchants,
				"sub_text":   fmt.Sprintf("%.1f%% of total", activePercentage),
				"icon":       "grid",
				"icon_color": "green",
			},
			"pending_kyc": gin.H{
				"title":      "Pending KYC",
				"value":      stats.PendingKyc,
				"sub_text":   "Awaiting verification",
				"icon":       "grid",
				"icon_color": "orange",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// ImpersonateMerchant impersonates merchant
// @Summary Admins impersonates specific merchant
// @Description Admins can impersonate specific merchant
// @Tags v2-merchants
// @Accept json
// @Produce json
// @Param request body ImpersonateMerchantRequest true "ImpersonateMerchantRequest details"
// @Success 201 {object} map[string]interface{} "Merchant impersonated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /merchants/admin/impersonate [post]
func (h *Handler) ImpersonateMerchant(c *gin.Context) {
	var req entity.ImpersonateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    auth_entity.ErrInvalidRequest,
				"message": err.Error(),
			},
		})
		return
	}

	response, err := h.useCase.ImpersonateMerchant(c, req.MerchantID)
	if err != nil {
		h.log.Error("Failed to impersonate merchant", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to impersonate merchant",
			},
		})
		return
	}

	// Get all user permissions grouped by merchant ID
	allPermissions, err := h.authService.GetAllUserPermissions(c.Request.Context(), response.User.ID)
	if err != nil {
		// Log error but don't fail the request
		allPermissions = make(map[string][]string)
	}

	h.log.Info("allPermissions", map[string]interface{}{
		"allPermissions": allPermissions,
	})

	// Build merchant-aware resource-operation map for better frontend UX
	merchantPermissions := make(map[string]map[string][]string)
	globalResourceOperations := make(map[string][]string)

	for merchantID, permissions := range allPermissions {
		if merchantID == "global" {
			// Build global resource-operation map
			for _, perm := range permissions {
				parts := strings.Split(perm, ":")
				if len(parts) == 2 {
					resource := parts[0]
					operation := parts[1]
					if _, exists := globalResourceOperations[resource]; !exists {
						globalResourceOperations[resource] = []string{}
					}
					globalResourceOperations[resource] = append(globalResourceOperations[resource], operation)
				}
			}
		} else {
			// Build merchant-specific resource-operation map
			resourceOperations := make(map[string][]string)
			for _, perm := range permissions {
				parts := strings.Split(perm, ":")
				if len(parts) == 2 {
					resource := parts[0]
					operation := parts[1]
					if _, exists := resourceOperations[resource]; !exists {
						resourceOperations[resource] = []string{}
					}
					resourceOperations[resource] = append(resourceOperations[resource], operation)
				}
			}
			merchantPermissions[merchantID] = resourceOperations
		}
	}

	h.log.Info("merchantPermissions", map[string]interface{}{
		"user":                response.User,
		"merchantPermissions": merchantPermissions,
		"allPermissions":      allPermissions,
	})

	// Get user groups
	groups, err := h.authService.GetUserGroupsByGroupedByMerchant(c.Request.Context(), response.User.ID)
	if err != nil {
		// Log error but don't fail the request
		groups = make(map[string][]auth_entity.Group)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": gin.H{
			"user":                         response.User,
			"merchants":                    response.Merchants,
			"access_token":                 response.Token,
			"refresh_token":                response.RefreshToken,
			"expires_at":                   response.ExpiresAt,
			"groups":                       groups,
			"merchant_permissions":         allPermissions,      // New: permissions grouped by merchant_id
			"merchant_resource_operations": merchantPermissions, // New: resource-operations grouped by merchant_id
			"global_resource_operations":   globalResourceOperations,
		},
	})
}
