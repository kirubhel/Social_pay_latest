package gin

import (
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/repository"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler manages merchant HTTP requests
type Handler struct {
	log        logging.Logger
	useCase    usecase.MerchantUseCase
	repository repository.Repository
}

// NewHandler creates a new merchant handler
func NewHandler(
	useCase usecase.MerchantUseCase,
	repository repository.Repository,
) *Handler {
	return &Handler{
		log:        logging.NewStdLogger("[V2_MERCHANT] [HANDLER]"),
		useCase:    useCase,
		repository: repository,
	}
}

// RegisterRoutes registers merchant management routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	// Public routes - no authentication required as requested
	merchants := r.Group("/merchants")
	{
		merchants.GET("/:id", h.GetMerchant)
		merchants.GET("/:id/details", h.GetMerchantDetails)
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

// GetMerchant godoc
// @Summary Get merchant by ID
// @Description Get a merchant by its ID
// @Tags v2-merchants
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 200 {object} SuccessResponse
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
// @Router /merchants/{id}/details [get]
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
