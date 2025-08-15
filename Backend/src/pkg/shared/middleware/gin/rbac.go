package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/service"
)

// AuthV2RBACConfig holds configuration for the AuthV2 RBAC middleware
type AuthV2RBACConfig struct {
	AuthService service.AuthService
}

type RBACV2 struct {
	AuthService service.AuthService
}

func NewRBACV2(authService service.AuthService) *RBACV2 {
	return &RBACV2{
		AuthService: authService,
	}
}

// RequirePermissionForMerchant creates a middleware that checks if the user has the required permission within a merchant context
func (rbac *RBACV2) RequirePermissionForMerchant(resource entity.Resource, operation entity.Operation) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "User not authenticated",
				},
			})
			c.Abort()
			return
		}

		// Get user to check user type
		user, exists := GetUserFromContext(c)
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

		// For admins and super admins, check global permissions (no merchant context needed)
		if user.UserType == entity.USER_TYPE_ADMIN || user.UserType == entity.USER_TYPE_SUPER_ADMIN {
			hasPermission, err := rbac.AuthService.CheckPermission(c.Request.Context(), userID, resource, operation)
			if err != nil {
				fmt.Printf("[RBAC-V2] Failed to check global admin permission: %v\n", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Success: false,
					Error: ApiError{
						Type:    "INTERNAL_SERVER_ERROR",
						Message: "Failed to check permissions",
					},
				})
				c.Abort()
				return
			}

			if !hasPermission {
				fmt.Printf("[RBAC-V2] Global permission denied for admin user %s: resource=%s, operation=%s\n",
					userID, resource, operation)
				c.JSON(http.StatusForbidden, ErrorResponse{
					Success: false,
					Error: ApiError{
						Type:    "FORBIDDEN",
						Message: "Insufficient admin permissions",
					},
				})
				c.Abort()
				return
			}

			fmt.Printf("[RBAC-V2] Global permission granted for admin user %s: resource=%s, operation=%s\n",
				userID, resource, operation)
			c.Next()
			return
		}

		// For merchants and members, require merchant context
		merchantID, exists := GetMerchantIDFromContext(c)
		if !exists {
			// Try to get merchant ID from X-MERCHANT-ID header
			merchantIDStr := c.GetHeader("X-MERCHANT-ID")
			if merchantIDStr == "" {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Success: false,
					Error: ApiError{
						Type:    "MERCHANT_CONTEXT_REQUIRED",
						Message: "Merchant context required. Please provide X-MERCHANT-ID header",
					},
				})
				c.Abort()
				return
			}

			var err error
			merchantID, err = uuid.Parse(merchantIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Success: false,
					Error: ApiError{
						Type:    "INVALID_MERCHANT_ID",
						Message: "Invalid merchant ID format",
					},
				})
				c.Abort()
				return
			}

			// Set merchant ID in context for future use
			c.Set(ContextKeyMerchantID, merchantID)
		}

		// Check merchant-specific permission
		hasPermission, err := rbac.AuthService.CheckPermissionForMerchant(c.Request.Context(), userID, merchantID, resource, operation)
		if err != nil {
			fmt.Printf("[RBAC-V2] Failed to check merchant permission: %v\n", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INTERNAL_SERVER_ERROR",
					Message: "Failed to check merchant permissions",
				},
			})
			c.Abort()
			return
		}

		if !hasPermission {
			fmt.Printf("[RBAC-V2] Merchant permission denied for user %s in merchant %s: resource=%s, operation=%s\n",
				userID, merchantID, resource, operation)
			c.JSON(http.StatusForbidden, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "FORBIDDEN",
					Message: "Insufficient permissions for this merchant",
				},
			})
			c.Abort()
			return
		}

		fmt.Printf("[RBAC-V2] Merchant permission granted for user %s in merchant %s: resource=%s, operation=%s\n",
			userID, merchantID, resource, operation)
		c.Next()
	}
}

// RequirePermissionForAdmin creates a middleware that checks if the user is admin/super_admin and has required admin permission
func (rbac *RBACV2) RequirePermissionForAdmin(resource entity.Resource, operation entity.Operation) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "UNAUTHORIZED",
					Message: "User not authenticated",
				},
			})
			c.Abort()
			return
		}

		// Get user to check user type
		user, exists := GetUserFromContext(c)
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

		// Only allow admin and super_admin user types
		if user.UserType != entity.USER_TYPE_ADMIN && user.UserType != entity.USER_TYPE_SUPER_ADMIN {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INSUFFICIENT_PRIVILEGES",
					Message: "Admin privileges required",
				},
			})
			c.Abort()
			return
		}

		// Check admin permission (no merchant context needed)
		hasPermission, err := rbac.AuthService.CheckPermission(c.Request.Context(), userID, resource, operation)
		if err != nil {
			fmt.Printf("[RBAC-V2] Failed to check admin permission: %v\n", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "INTERNAL_SERVER_ERROR",
					Message: "Failed to check admin permissions",
				},
			})
			c.Abort()
			return
		}

		if !hasPermission {
			fmt.Printf("[RBAC-V2] Admin permission denied for user %s: resource=%s, operation=%s\n",
				userID, resource, operation)
			c.JSON(http.StatusForbidden, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "FORBIDDEN",
					Message: "Insufficient admin permissions",
				},
			})
			c.Abort()
			return
		}

		fmt.Printf("[RBAC-V2] Admin permission granted for user %s: resource=%s, operation=%s\n",
			userID, resource, operation)
		c.Next()
	}
}

// PermissionPair represents a resource-operation permission pair
type PermissionPair struct {
	Resource  entity.Resource
	Operation entity.Operation
}

// NewPermissionPair creates a new permission pair
func NewPermissionPair(resource entity.Resource, operation entity.Operation) PermissionPair {
	return PermissionPair{
		Resource:  resource,
		Operation: operation,
	}
}

// Convenience functions for common permission combinations

// RequireTransactionRead creates middleware for transaction read permission
func (rbac *RBACV2) RequireTransactionRead() gin.HandlerFunc {
	return rbac.RequirePermissionForMerchant(entity.RESOURCE_TRANSACTION, entity.OPERATION_READ)
}

// RequireTransactionCreate creates middleware for transaction create permission
func (rbac *RBACV2) RequireTransactionCreate() gin.HandlerFunc {
	return rbac.RequirePermissionForMerchant(entity.RESOURCE_TRANSACTION, entity.OPERATION_CREATE)
}

// RequireUserManagement creates middleware for user management permissions (admin only)
func (rbac *RBACV2) RequireUserManagement() gin.HandlerFunc {
	return rbac.RequirePermissionForAdmin(entity.RESOURCE_USER, entity.OPERATION_ADMIN_CREATE)
}

// RequireMerchantOwner creates middleware that checks if user is owner of the merchant
func (rbac *RBACV2) RequireMerchantOwner() gin.HandlerFunc {
	return rbac.RequirePermissionForMerchant(entity.RESOURCE_ALL, entity.OPERATION_ALL)
}

// RequireAdminAccess creates middleware for admin-level operations
func (rbac *RBACV2) RequireAdminAccess() gin.HandlerFunc {
	return rbac.RequirePermissionForAdmin(entity.RESOURCE_ADMIN_ALL, entity.OPERATION_ADMIN_ALL)
}
