package gin

import (
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/entity"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/repository"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/usecase"
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
}

// NewHandler creates a new API key handler
func NewHandler(
	useCase usecase.APIKeyUseCase,
	repository repository.Repository,
	jwtAuth gin.HandlerFunc,
	publicMiddleware gin.HandlerFunc,
) *Handler {
	return &Handler{
		log:              logging.NewStdLogger("[APIKEY] [HANDLER]"),
		useCase:          useCase,
		repository:       repository,
		middleware:       &jwtAuth,
		publicMiddleware: &publicMiddleware,
	}
}

// RegisterRoutes registers API key management routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	// Protected routes
	keys := r.Group("/keys", *h.middleware, jwtMiddleware.MerchantIDMiddleware())
	{
		keys.GET("", h.ListAPIKeys)
		keys.POST("", h.CreateAPIKey)
		keys.GET("/:id", h.GetAPIKey)
		keys.PATCH("/:id", h.UpdateAPIKey)
		keys.DELETE("/:id", h.DeleteAPIKey)
		keys.POST("/:id/rotate", h.RotateAPIKeySecret)
	}

	// Public routes
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
// @Description Create a new API key for a merchant
// @Tags apikeys
// @Accept json
// @Produce json
// @Param request body entity.CreateAPIKeyRequest true "API key creation request"
// @Success 201 {object} entity.APIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys [post]
// @Security BearerAuth
func (h *Handler) CreateAPIKey(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.User.Id

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
// @Description Get an API key by ID
// @Tags apikeys
// @Produce json
// @Param id path string true "API key ID"
// @Success 200 {object} entity.APIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys/{id} [get]
// @Security BearerAuth
func (h *Handler) GetAPIKey(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.User.Id

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
// @Success 200 {array} entity.APIKeyResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys [get]
// @Security BearerAuth
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
// @Description Update an existing API key
// @Tags apikeys
// @Accept json
// @Produce json
// @Param id path string true "API key ID"
// @Param request body entity.UpdateAPIKeyRequest true "API key update request"
// @Success 200 {object} entity.APIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys/{id} [patch]
// @Security BearerAuth
func (h *Handler) UpdateAPIKey(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.User.Id

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
// @Produce json
// @Param id path string true "API key ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys/{id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteAPIKey(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.User.Id

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
// @Description Generate a new secret for an API key
// @Tags apikeys
// @Produce json
// @Param id path string true "API key ID"
// @Success 200 {object} entity.APIKeyRotateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /keys/{id}/rotate [post]
// @Security BearerAuth
func (h *Handler) RotateAPIKeySecret(c *gin.Context) {
	// Get session from context (middleware ensures this exists)
	session, _ := jwtMiddleware.GetSessionFromContext(c)
	userID := session.User.Id

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
// @Description Validate an API key with public and secret keys
// @Tags apikeys
// @Accept json
// @Produce json
// @Param request body entity.APIKeyValidateRequest true "API key validate request"
// @Success 200 {object} entity.APIKeyResponse
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
