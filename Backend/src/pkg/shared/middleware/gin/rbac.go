package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	acusecase "github.com/socialpay/socialpay/src/pkg/access_control/usecase"
)

// RBACMiddleware returns a Gin middleware that checks if the user has the required permission
func RBACMiddleware(accessControl acusecase.Interactor, resourceName, operationName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if !exists || userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		allowed := accessControl.CheckUserPermission(userID, resourceName, operationName)
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// WrapWithRBAC is a helper to wrap a handler with RBAC middleware for DRY route registration
func WrapWithRBAC(accessControl acusecase.Interactor, resourceName, operationName string, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		RBACMiddleware(accessControl, resourceName, operationName)(c)
		if c.IsAborted() {
			return
		}
		handler(c)
	}
}
