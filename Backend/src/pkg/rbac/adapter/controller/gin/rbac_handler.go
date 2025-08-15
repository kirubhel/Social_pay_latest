package gin

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	auth_service "github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	"github.com/socialpay/socialpay/src/pkg/rbac/core/entity"
	"github.com/socialpay/socialpay/src/pkg/rbac/core/service"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
)

// RBACHandler handles rbac HTTP requests
type RBACHandler struct {
	rbacService service.RBACService
	rbac        *ginn.RBACV2
	authService auth_service.AuthService
}

// NewRBACHandler creates a new team member handler
func NewRBACHandler(authService auth_service.AuthService, rbacService service.RBACService, rbac *ginn.RBACV2) *RBACHandler {
	return &RBACHandler{
		authService: authService,
		rbacService: rbacService,
		rbac:        rbac,
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

// GetResources godoc
// @Summary Get resources
// @Description Get list of resources
// @Tags v2-rbac
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rbac/resources [get]
func (h *RBACHandler) GetResources(c *gin.Context) {
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

	var isAdmin bool

	if user.UserType == auth_entity.USER_TYPE_ADMIN || user.UserType == auth_entity.USER_TYPE_SUPER_ADMIN {
		isAdmin = true
	} else {
		isAdmin = false
	}

	// Get resources
	resources, err := h.rbacService.GetResources(c, isAdmin)
	log.Println(err)
	if err != nil {
		if authErr, ok := err.(*entity.RBACError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
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
					"type":    entity.ErrInternalServer,
					"message": "Resources query failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Resources fetched successfully",
		"data":    resources,
	})
}

// RegisterRoutes registers all rbac routes
func (h *RBACHandler) RegisterRoutes(router *gin.RouterGroup) {
	team := router.Group("/rbac")

	// Create JWT config
	jwtConfig := ginn.JWTAuthMiddlewareConfig{
		AuthService: h.authService,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Public:      false,
	}

	team.GET("/resources",
		ginn.RequireJWTAuthentication(jwtConfig.AuthService, jwtConfig.JWTSecret),
		h.GetResources)
}
