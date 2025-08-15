# Authentication v2 System

A comprehensive authentication and authorization system for SocialPay that provides better developer experience, simplified architecture, and robust security features.

## Features

- ✅ **Simplified Authentication Flow**: Streamlined login/register with OTP verification
- ✅ **Role-Based Access Control (RBAC)**: Fine-grained permissions with enums for better DX
- ✅ **Multi-User Support**: Supports both merchants and super admins with team members
- ✅ **Activity Logging**: Comprehensive authentication activity tracking
- ✅ **Better DX**: Type-safe enums and intuitive middleware API
- ✅ **Phone Validation**: Robust phone number validation and OTP verification
- ✅ **JWT Authentication**: Secure token-based authentication
- ✅ **Auto-Seeding**: Automatic super admin account creation

## Quick Start

### 1. Basic Setup

```go
import (
    "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
    "github.com/socialpay/socialpay/src/pkg/authv2/service"
    ginMiddleware "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
)

// Create auth service (dependency injection)
authService := service.NewAuthService(repo, smsService, jwtSecret, logger)

// Create middleware
authV2 := ginMiddleware.NewAuthV2Middleware(authService, jwtSecret)
```

### 2. Authentication Endpoints

```go
// Register new user
auth.POST("/register", handleRegister(authService))

// Login (sends OTP)
auth.POST("/login", handleLogin(authService))

// Verify OTP (completes authentication)
auth.POST("/verify-otp", handleVerifyOTP(authService))

// Logout
auth.POST("/logout", authV2.Authenticate(), handleLogout(authService))
```

### 3. Protected Routes with RBAC

```go
// Require authentication
api.Use(authV2.Authenticate())

// Require specific user type
adminRoutes.Use(authV2.RequireAdmin())
merchantRoutes.Use(authV2.RequireMerchant())

// Require specific permissions (THIS IS THE MAGIC!)
transactionRoutes.GET("/", 
    authV2.RequirePermission(entity.RESOURCE_TRANSACTION, entity.OPERATION_READ),
    handleListTransactions())

transactionRoutes.POST("/", 
    authV2.RequirePermission(entity.RESOURCE_TRANSACTION, entity.OPERATION_CREATE),
    handleCreateTransaction())

// Admin-level operations
transactionRoutes.GET("/analytics", 
    authV2.RequirePermission(entity.RESOURCE_TRANSACTION, entity.OPERATION_ADMIN_READ),
    handleTransactionAnalytics())
```

## Core Concepts

### User Types

```go
entity.USER_TYPE_ADMIN    // Super admin and team members
entity.USER_TYPE_MERCHANT // Merchant users
```

### Resources (What can be accessed)

```go
entity.RESOURCE_TRANSACTION    // Transaction management
entity.RESOURCE_MERCHANT       // Merchant management
entity.RESOURCE_USER           // User management
entity.RESOURCE_ADMIN_WALLET   // Admin wallet operations
entity.RESOURCE_API_KEY        // API key management
entity.RESOURCE_WEBHOOK        // Webhook management
entity.RESOURCE_ANALYTICS      // Analytics and reports
entity.RESOURCE_COMMISSION     // Commission settings
entity.RESOURCE_QR             // QR code management
entity.RESOURCE_CHECKOUT       // Checkout management
entity.RESOURCE_NOTIFICATION   // Notification management
```

### Operations (What can be done)

```go
// Basic operations
entity.OPERATION_CREATE   // Create new records
entity.OPERATION_READ     // Read/view records
entity.OPERATION_UPDATE   // Update existing records
entity.OPERATION_DELETE   // Delete records

// Admin-level operations (for super admins)
entity.OPERATION_ADMIN_CREATE   // Admin-level creation
entity.OPERATION_ADMIN_READ     // Admin-level reading
entity.OPERATION_ADMIN_UPDATE   // Admin-level updates
entity.OPERATION_ADMIN_DELETE   // Admin-level deletion
```

## Authentication Flow

### 1. Registration

```json
POST /api/auth/register
{
  "title": "Mr",
  "first_name": "John",
  "last_name": "Doe",
  "phone_prefix": "251",
  "phone_number": "911234567",
  "password": "SecurePass123!",
  "password_hint": "my secure password",
  "user_type": "merchant"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": { ... },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "...",
    "expires_at": 1640995200
  }
}
```

### 2. Login (Two-Step Process)

**Step 1: Login with credentials**
```json
POST /api/auth/login
{
  "phone_prefix": "251",
  "phone_number": "911234567",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "OTP sent successfully",
  "data": {
    "otp_token": "abcd1234..."
  }
}
```

**Step 2: Verify OTP**
```json
POST /api/auth/verify-otp
{
  "token": "abcd1234...",
  "code": "123456"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": { ... },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "...",
    "expires_at": 1640995200
  }
}
```

## Middleware Usage

### 1. Basic Authentication

```go
// Require any authenticated user
api.Use(authV2.Authenticate())
```

### 2. User Type Restrictions

```go
// Admin only
adminRoutes.Use(authV2.RequireAdmin())

// Merchant only
merchantRoutes.Use(authV2.RequireMerchant())

// Multiple user types
api.Use(authV2.RequireUserType(entity.USER_TYPE_ADMIN, entity.USER_TYPE_MERCHANT))
```

### 3. Permission-Based Access Control

```go
// Check specific permissions
router.GET("/transactions", 
    authV2.RequirePermission(entity.RESOURCE_TRANSACTION, entity.OPERATION_READ),
    handleTransactions)

// Chain multiple middleware
router.POST("/admin/users",
    authV2.RequireAdmin(),  // Must be admin
    authV2.RequirePermission(entity.RESOURCE_USER, entity.OPERATION_ADMIN_CREATE),  // And have permission
    handleCreateUser)
```

### 4. Context Helpers

```go
func handleSomeEndpoint(c *gin.Context) {
    // Get user information from context
    userID, exists := ginMiddleware.GetUserIDFromContextV2(c)
    userType, exists := ginMiddleware.GetUserTypeFromContextV2(c)
    merchantID, exists := ginMiddleware.GetMerchantIDFromContextV2(c)
    claims, exists := ginMiddleware.GetClaimsFromContextV2(c)
    
    // Use the information...
}
```

## Super Admin & Team Management

### 1. Seed Super Admin

```bash
# Run the seeder
go run src/pkg/authv2/cmd/seed/main.go
```

**Default Super Admin:**
- **Phone:** +251961186323
- **Password:** SocialPay$123SuperAdmiN
- **User Type:** admin

### 2. Create Team Members

```go
// Admin can create team members with specific roles
teamMember, err := authService.CreateTeamMember(
    ctx, 
    adminUserID, 
    &entity.CreateUserRequest{...}, 
    groupID)
```

## Activity Logging

All authentication activities are automatically logged:

```go
// Activities are logged automatically, but you can also log custom activities
authService.LogActivity(ctx, userID, entity.ACTIVITY_LOGIN_SUCCESS, ip, userAgent, deviceName, true, details)
```

**Activity Types:**
- `ACTIVITY_LOGIN_SUCCESS`
- `ACTIVITY_LOGIN_FAILED`
- `ACTIVITY_LOGOUT`
- `ACTIVITY_OTP_SENT`
- `ACTIVITY_OTP_VERIFIED`
- `ACTIVITY_OTP_FAILED`
- `ACTIVITY_PASSWORD_CHANGED`
- `ACTIVITY_ACCOUNT_CREATED`
- `ACTIVITY_PERMISSION_DENIED`

## Database Schema

The system works with the existing database schema:

```sql
-- Uses existing tables:
auth.users
auth.phones
auth.phone_identities
auth.password_identities
auth.sessions
auth.devices
auth.groups
auth.user_groups
auth.permissions
auth.group_permissions
auth.resources
auth.operations

-- New table for activity logging:
auth.auth_activities
```

## Configuration

### Environment Variables

```bash
DATABASE_URL=postgres://user:pass@localhost/socialpay
JWT_SECRET=your-super-secret-jwt-key
SMS_API_KEY=your-sms-provider-api-key
```

### Validation Rules

- **Phone Numbers:** 7-15 digits, prefix 1-4 digits
- **Passwords:** Minimum 8 characters, must contain uppercase, lowercase, digit, and special character
- **OTP:** 6-digit numeric code, expires in 5 minutes
- **JWT Tokens:** 24-hour expiry by default

## Security Features

- ✅ **Bcrypt Password Hashing**
- ✅ **JWT Token Validation**
- ✅ **OTP Expiration**
- ✅ **Session Management**
- ✅ **Input Sanitization**
- ✅ **SQL Injection Prevention**
- ✅ **Permission-Based Access Control**
- ✅ **Activity Logging for Audit Trail**

## Migration from Old Auth

1. **Gradual Migration**: Old and new auth can coexist
2. **Same Database Schema**: Uses existing tables where possible
3. **New Endpoints**: Create new v2 endpoints alongside existing ones
4. **Middleware**: Use new middleware for new routes

## Best Practices

1. **Always validate inputs** using the provided validation functions
2. **Use enums** for resources and operations for type safety
3. **Chain middleware** in logical order (auth → user type → permissions)
4. **Log activities** for important operations
5. **Use context helpers** to get user information
6. **Handle errors** gracefully with proper HTTP status codes

## Example: Complete Transaction API

```go
func SetupTransactionRoutes(router *gin.Engine, authV2 *ginMiddleware.AuthV2Middleware) {
    txRoutes := router.Group("/api/v2/transactions")
    txRoutes.Use(authV2.Authenticate())  // All routes require authentication
    
    // Anyone can read their own transactions
    txRoutes.GET("/", 
        authV2.RequirePermission(entity.RESOURCE_TRANSACTION, entity.OPERATION_READ),
        handleListUserTransactions)
    
    // Merchants can create transactions
    txRoutes.POST("/", 
        authV2.RequireMerchant(),
        authV2.RequirePermission(entity.RESOURCE_TRANSACTION, entity.OPERATION_CREATE),
        handleCreateTransaction)
    
    // Admin-only routes
    adminTxRoutes := txRoutes.Group("/admin")
    adminTxRoutes.Use(authV2.RequireAdmin())
    {
        // View all transactions (admin permission required)
        adminTxRoutes.GET("/all", 
            authV2.RequirePermission(entity.RESOURCE_TRANSACTION, entity.OPERATION_ADMIN_READ),
            handleListAllTransactions)
        
        // Refund transactions (admin permission required)
        adminTxRoutes.POST("/:id/refund", 
            authV2.RequirePermission(entity.RESOURCE_TRANSACTION, entity.OPERATION_ADMIN_UPDATE),
            handleRefundTransaction)
        
        // Analytics (admin permission required)
        adminTxRoutes.GET("/analytics", 
            authV2.RequirePermission(entity.RESOURCE_ANALYTICS, entity.OPERATION_ADMIN_READ),
            handleTransactionAnalytics)
    }
}
```

This provides a clean, secure, and developer-friendly authentication system that scales from simple merchant applications to complex admin operations with fine-grained permission control. 
