# RBAC Middleware Documentation

## Overview
The RBAC (Role-Based Access Control) middleware provides a robust way to implement permission-based access control in your Gin-based applications. It ensures that only authenticated users with the appropriate permissions can access specific resources and operations.

## Features
- Authentication verification
- Permission-based access control
- Resource and operation-level granularity
- Easy-to-use middleware and helper functions
- Clean error handling with appropriate HTTP status codes

## Prerequisites
- Gin web framework
- Access Control Interactor implementation
- User authentication system (must provide user ID in context)

## Installation
The RBAC middleware is part of the `gin` package in the shared middleware directory:

```go
import (
    "github.com/gin-gonic/gin"
    acusecase "github.com/socialpay/socialpay/src/pkg/access_control/usecase"
)
```

## Usage

### Basic Middleware Usage
```go
router.GET("/api/resource", 
    gin.RBACMiddleware(accessControl, "resource_name", "operation_name"),
    handlerFunction
)
```

### Using the Helper Function
```go
router.GET("/api/resource", 
    gin.WrapWithRBAC(accessControl, "resource_name", "operation_name", handlerFunction)
)
```

## Setup Guide

### 1. Create Operations
First, define the operations that can be performed on resources:

```go
// Create basic CRUD operations
createOp, err := accessControl.CreateOperations("create", "Create operation")
readOp, err := accessControl.CreateOperations("read", "Read operation")
updateOp, err := accessControl.CreateOperations("update", "Update operation")
deleteOp, err := accessControl.CreateOperations("delete", "Delete operation")
```

### 2. Create Resources
Define resources and associate them with operations:

```go
// Create a resource with its allowed operations
resource, err := accessControl.CreateResource(
    "users",                                    // resource name
    "User management resource",                 // description
    []uuid.UUID{createOp.ID, readOp.ID, updateOp.ID, deleteOp.ID}, // operations
)
```

### 3. Create Permissions
Create permissions that combine resources and operations:

```go
// Create a permission for reading users
permission, err := accessControl.CreatePermission(
    "users",                    // resource name
    resource.ID,               // resource ID
    []uuid.UUID{readOp.ID},    // operations
    "allow",                   // effect
)
```

### 4. Assign Permissions
Grant permissions to users or groups:

```go
// Grant permission to a user
err := accessControl.GrantPermissionToUser(userID, permission.ID)

// Or grant to a group
err := accessControl.GrantPermissionToGroup(groupID, permission.ID)
```

### 5. Protect Routes
Apply the RBAC middleware to your routes:

```go
// Example route protection
router.GET("/api/users", 
    gin.WrapWithRBAC(accessControl, "users", "read", getUserHandler)
)

router.POST("/api/users", 
    gin.WrapWithRBAC(accessControl, "users", "create", createUserHandler)
)

router.PUT("/api/users/:id", 
    gin.WrapWithRBAC(accessControl, "users", "update", updateUserHandler)
)

router.DELETE("/api/users/:id", 
    gin.WrapWithRBAC(accessControl, "users", "delete", deleteUserHandler)
)
```

## Error Handling
The middleware handles two types of errors:

1. **Unauthorized (401)**
   - Triggered when user is not authenticated
   - Response: `{"error": "User not authenticated"}`

2. **Forbidden (403)**
   - Triggered when user lacks required permissions
   - Response: `{"error": "Forbidden: insufficient permissions"}`

## Best Practices

1. **Resource Naming**
   - Use consistent, lowercase names for resources
   - Use plural form for resource names (e.g., "users" instead of "user")
   - Keep resource names descriptive and clear

2. **Operation Naming**
   - Use standard CRUD operations when possible
   - Use lowercase for operation names
   - Be specific about custom operations

3. **Permission Management**
   - Group related permissions together
   - Use groups for managing permissions at scale
   - Regularly audit and review permissions

4. **Route Organization**
   - Group routes by resource
   - Apply consistent permission patterns
   - Document permission requirements

## Example Implementation

```go
func setupRBACRoutes(router *gin.Engine, accessControl acusecase.Interactor) {
    // User management routes
    userGroup := router.Group("/api/users")
    {
        userGroup.GET("", 
            gin.WrapWithRBAC(accessControl, "users", "read", getUserHandler))
        userGroup.POST("", 
            gin.WrapWithRBAC(accessControl, "users", "create", createUserHandler))
        userGroup.PUT("/:id", 
            gin.WrapWithRBAC(accessControl, "users", "update", updateUserHandler))
        userGroup.DELETE("/:id", 
            gin.WrapWithRBAC(accessControl, "users", "delete", deleteUserHandler))
    }

    // Transaction management routes
    transactionGroup := router.Group("/api/transactions")
    {
        transactionGroup.GET("", 
            gin.WrapWithRBAC(accessControl, "transactions", "read", getTransactionHandler))
        transactionGroup.POST("", 
            gin.WrapWithRBAC(accessControl, "transactions", "create", createTransactionHandler))
    }
}
```

## Troubleshooting

1. **Permission Not Working**
   - Verify user authentication
   - Check if permission is correctly assigned
   - Confirm resource and operation names match exactly
   - Check permission effect is set to "allow"

2. **Unauthorized Errors**
   - Ensure user ID is properly set in context
   - Verify authentication middleware is running before RBAC
   - Check if user ID is valid UUID

3. **Forbidden Errors**
   - Verify user has the correct permission
   - Check if permission is assigned to user or their group
   - Confirm resource and operation names are correct

## Security Considerations

1. **Always verify authentication before RBAC**
2. **Use HTTPS for all API endpoints**
3. **Regularly audit permissions and access patterns**
4. **Implement proper logging for security events**
5. **Follow principle of least privilege**

## Contributing
When contributing to the RBAC middleware:

1. Follow existing code style
2. Add tests for new features
3. Update documentation
4. Consider backward compatibility
5. Add examples for new features 