# AuthV2 Seeder

The AuthV2 seeder provides functionality to seed the authentication system with RBAC (Role-Based Access Control) data and create the default super admin user.

## Features

- JSON-configurable operations and resources for easy editing
- Creates super admin user with specified credentials
- Creates admin wallet for super admin user
- Sets up complete RBAC system with permissions and groups
- Never removes existing super admin (safe operation)
- Idempotent operations (can be run multiple times safely)

## Usage

### Basic Usage

```go
import (
    "context"
    "database/sql"
    "log"
    
    "github.com/socialpay/socialpay/src/pkg/authv2/seeder"
    "github.com/socialpay/socialpay/src/pkg/authv2/core/service"
    "github.com/socialpay/socialpay/src/pkg/authv2/core/repository"
)

// Initialize seeder
authSeeder := seeder.NewAuthSeeder(authService, authRepo, db, logger)

// Seed everything
ctx := context.Background()
err := authSeeder.SeedAll(ctx)
if err != nil {
    log.Fatalf("Failed to seed auth system: %v", err)
}
```

### RBAC Configuration

The seeder supports JSON configuration for operations and resources. It looks for the config file in these locations (in order):

1. `src/pkg/authv2/seeder/rbac_config.json`
2. `pkg/authv2/seeder/rbac_config.json`
3. `rbac_config.json`

If no config file is found, it uses default configuration.

#### JSON Config Format

```json
{
  "operations": [
    {
      "name": "CREATE",
      "description": "Create operation"
    },
    {
      "name": "READ", 
      "description": "Read operation"
    }
  ],
  "resources": [
    {
      "name": "transaction",
      "description": "Transaction management"
    },
    {
      "name": "merchant",
      "description": "Merchant management"
    }
  ]
}
```

### Default Super Admin

The seeder creates a default super admin with these credentials:

- **Title**: mr
- **First Name**: SocialPay
- **Last Name**: SuperAdmin
- **Phone**: +251961186323
- **Password**: SocialPay$123SuperAdmiN
- **Password Hint**: superadmin
- **User Type**: admin

**Important**: The seeder will never remove an existing super admin. If the super admin already exists, it will skip creation and log a message.

### Admin Wallet

The seeder also creates an admin wallet for the super admin user with these properties:

- **Wallet Type**: admin_wallet
- **Currency**: ETB
- **Initial Amount**: 0.00
- **Locked Amount**: 0.00
- **Merchant ID**: NULL (admin wallets don't belong to merchants)

The wallet creation is also idempotent - if an admin wallet already exists for the user, it will skip creation.

## What Gets Created

### Operations
- CREATE, READ, UPDATE, DELETE
- ADMIN_READ, ADMIN_WRITE
- APPROVE, REJECT, PROCESS, EXPORT

### Resources
- transaction, merchant, user
- admin_wallet, ip_whitelist, api_key
- webhook, analytics, commission
- qr, checkout, notification
- team, role, system_settings, audit_log

### Groups
- **super_admin**: Group with all permissions to all resources

### Permissions
- All possible combinations of operations and resources
- Assigned to super_admin group

### Wallets
- **Admin wallet**: Created for super admin user with wallet_type 'super_admin'

## Safety Features

- **Idempotent**: Can be run multiple times without issues
- **Never removes super admin**: Existing super admin accounts are preserved
- **Conflict handling**: Uses `ON CONFLICT DO NOTHING` for database inserts
- **Validation**: Uses auth service for user creation (includes validation)

## Error Handling

The seeder provides detailed error messages and logs all operations. If any step fails, it returns a descriptive error indicating what went wrong.

## Integration with Auth Service

The seeder integrates with the AuthService interface and uses proper business logic for user creation, ensuring:

- Password hashing
- Input validation
- Proper user type assignment
- Activity logging (if implemented in auth service) 