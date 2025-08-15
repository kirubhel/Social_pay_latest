package gin

import (
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/entity"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/repository"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/usecase"
	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	jwtMiddleware "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler manages API key HTTP requests
type Handler struct {
	log              logging.Logger
	useCase          usecase.APIKeyUseCase
	repository       repository.Repository
	middleware       *gin.HandlerFunc
	publicMiddleware *gin.HandlerFunc
	rbac             *jwtMiddleware.RBACV2
}

// NewHandler creates a new API key handler
func NewHandler(
	useCase usecase.APIKeyUseCase,
	repository repository.Repository,
	jwtAuth gin.HandlerFunc,
	publicMiddleware gin.HandlerFunc,
	rbac *jwtMiddleware.RBACV2,
) *Handler {
	return &Handler{
		log:              logging.NewStdLogger("[APIKEY] [HANDLER]"),
		useCase:          useCase,
		repository:       repository,
		middleware:       &jwtAuth,
		publicMiddleware: &publicMiddleware,
		rbac:             rbac,
	}
}

// RegisterRoutes registers API key management routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	// Protected routes with RBAC
	keys := r.Group("/keys", *h.middleware, jwtMiddleware.MerchantIDMiddleware())
	{
		keys.GET("",
			h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_API_KEY, auth_entity.OPERATION_READ),
			h.ListAPIKeys)
		keys.POST("",
			h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_API_KEY, auth_entity.OPERATION_CREATE),
			h.CreateAPIKey)
		keys.GET("/:id",
			h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_API_KEY, auth_entity.OPERATION_READ),
			h.GetAPIKey)
		keys.PATCH("/:id",
			h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_API_KEY, auth_entity.OPERATION_UPDATE),
			h.UpdateAPIKey)
		keys.DELETE("/:id",
			h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_API_KEY, auth_entity.OPERATION_DELETE),
			h.DeleteAPIKey)
		keys.POST("/:id/rotate",
			h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_API_KEY, auth_entity.OPERATION_UPDATE),
			h.RotateAPIKeySecret)
	}

	// Public routes (no RBAC needed)
	public := r.Group("/public/keys", *h.publicMiddleware)
	{
		public.POST("/validate", h.ValidateAPIKey)
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
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

// CreateAPIKey godoc
// @Summary Create API key
// @Description Create a new API key for a merchant with specified permissions
// @Tags API Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.CreateAPIKeyRequest true "API key creation request"
// @Success 201 {object} entity.APIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys [post]
// @Router /keys [post]
// @Security BearerAuth
func (h *Handler) CreateAPIKey(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.UserID

	// Get merchant ID from context
	merchantID, exists := jwtMiddleware.GetMerchantIDFromContext(c)
	h.log.Info("merchantID", map[string]interface{}{
		"merchantID": merchantID,
		"exists":     exists,
	})
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

	// Parse request
	var req entity.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// Create API key with userID as both the owner and creator
	apiKey, secretKey, err := h.useCase.CreateAPIKey(c.Request.Context(), userID, userID, merchantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Return response with secret key
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data: map[string]interface{}{
			"api_key":    apiKey,
			"secret_key": secretKey,
		},
	})
}

// GetAPIKey godoc
// @Summary Get API key
// @Description Get an API key by ID for the authenticated merchant
// @Tags API Keys
// @Produce json
// @Security BearerAuth
// @Param id path string true "API key ID" format(uuid)
// @Success 200 {object} entity.APIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys/{id} [get]
func (h *Handler) GetAPIKey(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.UserID

	// Parse API key ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid API key ID",
			},
		})
		return
	}

	// Get API key
	apiKey, err := h.useCase.GetAPIKey(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Check if the API key belongs to the user
	if apiKey.UserID != userID {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "NOT_FOUND",
				Message: "API key not found",
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data: map[string]interface{}{
			"api_key": apiKey,
		},
	})
}

// ListAPIKeys godoc
// @Summary List API keys
// @Description List all API keys for the authenticated user
// @Tags apikeys
// @Produce json
// ListAPIKeys godoc
// @Summary List API keys
// @Description Get all API keys for the authenticated merchant
// @Tags API Keys
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse{data=[]entity.APIKeyResponse}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys [get]
func (h *Handler) ListAPIKeys(c *gin.Context) {
	// Get session from context (middleware ensures this exists)

	merchantID, exists := jwtMiddleware.GetMerchantIDFromContext(c)
	h.log.Info("merchantID", map[string]interface{}{
		"merchantID": merchantID,
		"exists":     exists,
	})

	// Get API keys
	apiKeys, err := h.useCase.GetAPIKeysByMerchantID(c.Request.Context(), merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data: map[string]interface{}{
			"api_keys": apiKeys,
		},
	})
}

// UpdateAPIKey godoc
// @Summary Update API key
// @Description Update an existing API key for the authenticated merchant
// @Tags API Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "API key ID" format(uuid)
// @Param request body entity.UpdateAPIKeyRequest true "API key update request"
// @Success 200 {object} entity.APIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys/{id} [patch]
func (h *Handler) UpdateAPIKey(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.UserID

	// Parse API key ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid API key ID",
			},
		})
		return
	}

	// First, get the API key to verify ownership
	apiKey, err := h.useCase.GetAPIKey(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Check if API key exists and belongs to the user
	if apiKey == nil || apiKey.UserID != userID {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "NOT_FOUND",
				Message: "API key not found",
			},
		})
		return
	}

	// Parse request
	var req entity.UpdateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// Update API key
	updatedAPIKey, err := h.useCase.UpdateAPIKey(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data: map[string]interface{}{
			"api_key": updatedAPIKey,
		},
	})
}

// DeleteAPIKey godoc
// @Summary Delete API key
// @Description Delete an existing API key
// @Tags apikeys
// DeleteAPIKey godoc
// @Summary Delete API key
// @Description Delete an API key for the authenticated merchant
// @Tags API Keys
// @Produce json
// @Security BearerAuth
// @Param id path string true "API key ID" format(uuid)
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys/{id} [delete]
func (h *Handler) DeleteAPIKey(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.UserID

	// Parse API key ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid API key ID",
			},
		})
		return
	}

	// First, get the API key to verify ownership
	apiKey, err := h.useCase.GetAPIKey(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Check if API key exists and belongs to the user
	if apiKey == nil || apiKey.UserID != userID {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "NOT_FOUND",
				Message: "API key not found",
			},
		})
		return
	}

	// Delete API key
	err = h.useCase.DeleteAPIKey(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    map[string]interface{}{},
	})
}

// RotateAPIKeySecret godoc
// @Summary Rotate API key secret
// @Description Generate a new secret for an API key for the authenticated merchant
// @Tags API Keys
// @Produce json
// @Security BearerAuth
// @Param id path string true "API key ID" format(uuid)
// @Success 200 {object} entity.APIKeyRotateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys/{id}/rotate [post]
func (h *Handler) RotateAPIKeySecret(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.UserID

	// Parse API key ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid API key ID",
			},
		})
		return
	}

	// First, get the API key to verify ownership
	apiKey, err := h.useCase.GetAPIKey(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Check if API key exists and belongs to the user
	if apiKey == nil || apiKey.UserID != userID {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "NOT_FOUND",
				Message: "API key not found",
			},
		})
		return
	}

	// Rotate API key secret
	response, err := h.useCase.RotateAPIKeySecret(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Return response with new secret key
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data: map[string]interface{}{
			"api_key":    response.APIKey,
			"secret_key": response.SecretKey,
		},
	})
}

// ValidateAPIKey godoc
// @Summary Validate API key
// @Description Validate an API key with public and secret keys (public endpoint)
// @Tags API Keys
// @Accept json
// @Produce json
// @Param request body entity.APIKeyValidateRequest true "API key validate request"
// @Success 200 {object} SuccessResponse{data=entity.APIKeyResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /public/keys/validate [post]
func (h *Handler) ValidateAPIKey(c *gin.Context) {
	// Parse request
	var req entity.APIKeyValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// Validate API key
	apiKey, err := h.useCase.ValidateAPIKey(c.Request.Context(), req.PublicKey, req.SecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "Invalid API key",
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data: map[string]interface{}{
			"api_key": apiKey,
		},
	})
}
