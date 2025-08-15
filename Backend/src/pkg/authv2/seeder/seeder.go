package seeder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/repository"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	"github.com/lib/pq"
)

// AuthSeeder handles seeding of authentication data
type AuthSeeder struct {
	authService service.AuthService
	repo        repository.AuthRepository
	db          *sql.DB
	logger      *log.Logger
}

// NewAuthSeeder creates a new auth seeder
func NewAuthSeeder(authService service.AuthService, repo repository.AuthRepository, db *sql.DB) *AuthSeeder {
	return &AuthSeeder{
		authService: authService,
		repo:        repo,
		db:          db,
		logger:      log.New(os.Stdout, "", log.LstdFlags),
	}
}

// OperationConfig represents operation configuration from JSON
type OperationConfig struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsAdminOperation bool   `json:"is_admin_operation"`
}

// ResourceConfig represents resource configuration from JSON
type ResourceConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RBACConfig represents the complete RBAC configuration
type RBACConfig struct {
	Operations []OperationConfig `json:"operations"`
	Resources  []ResourceConfig  `json:"resources"`
}

// SeedAll seeds all authentication data including RBAC and super admin
func (s *AuthSeeder) SeedAll(ctx context.Context) error {
	s.logger.Println("Starting complete authentication seeding...")

	// Load RBAC config from JSON
	rbacConfig, err := s.loadRBACConfig()
	if err != nil {
		return fmt.Errorf("failed to load RBAC config: %w", err)
	}

	// Seed RBAC system
	if err := s.seedRBAC(ctx, rbacConfig); err != nil {
		return fmt.Errorf("failed to seed RBAC: %w", err)
	}

	// Create super admin if not exists
	_, user := s.createDefaultSuperAdmin(ctx)
	if user.ID == uuid.Nil {
		return fmt.Errorf("failed to create super admin")
	}

	// Create admin wallet if not exists
	if err := s.createAdminWalletIfNotExists(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to create admin wallet: %w", err)
	}

	s.logger.Println("Authentication seeding completed successfully")
	return nil
}

// loadRBACConfig loads RBAC configuration from JSON file
func (s *AuthSeeder) loadRBACConfig() (*RBACConfig, error) {
	// Look for config file in multiple locations
	configPaths := []string{
		"./src/pkg/authv2/seeder/rbac_config.json",
		"./pkg/authv2/seeder/rbac_config.json",
		"./rbac_config.json",
	}

	var configData []byte
	var err error

	for _, path := range configPaths {
		configData, err = os.ReadFile(path)
		if err == nil {
			s.logger.Printf("Loaded RBAC config from: %s", path)
			break
		}
	}

	if err != nil {
		// If no config file found, use default configuration
		s.logger.Println("No RBAC config file found, using default configuration")
		return s.getDefaultRBACConfig(), nil
	}

	var config RBACConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse RBAC config JSON: %w", err)
	}

	return &config, nil
}

// getDefaultRBACConfig returns default RBAC configuration
func (s *AuthSeeder) getDefaultRBACConfig() *RBACConfig {
	return &RBACConfig{
		Operations: []OperationConfig{
			{Name: "CREATE", Description: "Create operation", IsAdminOperation: false},
			{Name: "READ", Description: "Read operation", IsAdminOperation: false},
			{Name: "UPDATE", Description: "Update operation", IsAdminOperation: false},
			{Name: "DELETE", Description: "Delete operation", IsAdminOperation: false},
			{Name: "ADMIN_READ", Description: "Admin-level read operation", IsAdminOperation: true},
			{Name: "ADMIN_WRITE", Description: "Admin-level write operation", IsAdminOperation: true},
		},
		Resources: []ResourceConfig{
			{Name: "transaction", Description: "Transaction management"},
			{Name: "merchant", Description: "Merchant management"},
			{Name: "user", Description: "User management"},
			{Name: "admin_wallet", Description: "Admin wallet operations"},
			{Name: "ip_whitelist", Description: "IP whitelist management"},
			{Name: "api_key", Description: "API key management"},
			{Name: "webhook", Description: "Webhook management"},
			{Name: "analytics", Description: "Analytics and reports"},
			{Name: "commission", Description: "Commission settings"},
			{Name: "qr", Description: "QR code management"},
			{Name: "checkout", Description: "Checkout management"},
			{Name: "notification", Description: "Notification management"},
		},
	}
}

// seedRBAC seeds the RBAC system with operations, resources, permissions, and roles
func (s *AuthSeeder) seedRBAC(ctx context.Context, config *RBACConfig) error {
	s.logger.Println("Seeding RBAC system...")

	// Seed operations
	for _, op := range config.Operations {
		if err := s.createOperationIfNotExists(ctx, op.Name, op.Description, op.IsAdminOperation); err != nil {
			return fmt.Errorf("failed to create operation %s: %w", op.Name, err)
		}
	}

	// Seed resources
	for _, res := range config.Resources {
		if err := s.createResourceIfNotExists(ctx, res.Name, res.Description); err != nil {
			return fmt.Errorf("failed to create resource %s: %w", res.Name, err)
		}
	}

	// Create super admin group with all permissions
	if err := s.createSuperAdminGroup(ctx); err != nil {
		return fmt.Errorf("failed to create super admin group: %w", err)
	}

	return nil
}

// createDefaultSuperAdmin creates the default super admin user if not exists
func (s *AuthSeeder) createDefaultSuperAdmin(ctx context.Context) (error, entity.User) {
	phonePrefix := "251"
	phoneNumber := "911237975"

	// Check if super admin already exists
	exists, user := s.superAdminExists(ctx, phonePrefix, phoneNumber)

	if exists {
		s.logger.Println("Super admin already exists, skipping creation")
		return nil, user
	}

	// Create super admin
	req := &entity.CreateUserRequest{
		Title:        "mr",
		FirstName:    "SocialPay",
		LastName:     "SuperAdmin",
		Email:        "superadmin@socialpay.com",
		PhonePrefix:  phonePrefix,
		PhoneNumber:  phoneNumber,
		Password:     "SocialPay$123SuperAdmiN",
		PasswordHint: "superadmin",
		UserType:     "super_admin",
	}

	// Create user
	userRef, err := s.authService.CreateSuperAdminUser(ctx, req)
	if err != nil {
		return err, *userRef
	}

	s.logger.Printf("Super admin created: %v", userRef)

	return nil, *userRef
}

// superAdminExists checks if super admin with given phone already exists
func (s *AuthSeeder) superAdminExists(ctx context.Context, phonePrefix, phoneNumber string) (bool, entity.User) {
	var exists bool
	var user entity.User
	err := s.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM auth.users WHERE phone_prefix = $1 AND phone_number = $2 AND user_type = 'super_admin')`,
		phonePrefix, phoneNumber).Scan(&exists)
	if err != nil {
		return false, user
	}

	if exists {
		err = s.db.QueryRowContext(ctx,
			`SELECT id, first_name, last_name, phone_prefix, phone_number, user_type, created_at, updated_at 
			FROM auth.users WHERE phone_prefix = $1 AND phone_number = $2 AND user_type = 'super_admin'`,
			phonePrefix, phoneNumber).Scan(
			&user.ID, &user.FirstName, &user.LastName,
			&user.PhonePrefix, &user.PhoneNumber, &user.UserType, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			s.logger.Printf("Error getting super admin user: %v", err)
			return false, user
		}
	}

	return exists, user
}

// Helper methods for seeding

func (s *AuthSeeder) createOperationIfNotExists(ctx context.Context, name, description string, isAdminOperation bool) error {
	var exists bool
	err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM auth.operations WHERE name = $1)`, name).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		// Create new operation
		_, err = s.db.ExecContext(ctx, `
			INSERT INTO auth.operations (id, name, description, is_admin_operation, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			uuid.New(), name, description, isAdminOperation, time.Now(), time.Now())
		if err != nil {
			return err
		}
		s.logger.Printf("Created operation: %s (admin: %v)", name, isAdminOperation)
	} else {
		// Update existing operation if description or is_admin_operation changed
		_, err = s.db.ExecContext(ctx, `
			UPDATE auth.operations 
			SET description = $2, is_admin_operation = $3, updated_at = $4
			WHERE name = $1`,
			name, description, isAdminOperation, time.Now())
		if err != nil {
			return err
		}
		s.logger.Printf("Updated operation: %s (admin: %v)", name, isAdminOperation)
	}

	return nil
}

func (s *AuthSeeder) createResourceIfNotExists(ctx context.Context, name, description string) error {
	var exists bool
	err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM auth.resources WHERE name = $1)`, name).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		// Get all operation IDs
		rows, err := s.db.QueryContext(ctx, `SELECT id FROM auth.operations`)
		if err != nil {
			return err
		}
		defer rows.Close()

		var operationIDs []uuid.UUID
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				return err
			}
			operationIDs = append(operationIDs, id)
		}

		// Use pq.Array for PostgreSQL array format
		_, err = s.db.ExecContext(ctx, `
			INSERT INTO auth.resources (id, name, description, operations, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			uuid.New(), name, description, pq.Array(operationIDs), time.Now(), time.Now())
		if err != nil {
			return err
		}
		s.logger.Printf("Created resource: %s", name)
	}

	return nil
}

func (s *AuthSeeder) createSuperAdminGroup(ctx context.Context) error {
	var exists bool
	err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM auth.groups WHERE title = 'super_admin')`).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		groupID := uuid.New()
		_, err = s.db.ExecContext(ctx, `
			INSERT INTO auth.groups (id, title, created_at, updated_at)
			VALUES ($1, $2, $3, $4)`,
			groupID, "super_admin", time.Now(), time.Now())
		if err != nil {
			return err
		}

		// Assign all permissions to super admin group
		rows, err := s.db.QueryContext(ctx, `
			SELECT r.id, array_agg(o.id)
			FROM auth.resources r
			CROSS JOIN auth.operations o
			GROUP BY r.id`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var resourceID uuid.UUID
			var operationIDs []uuid.UUID
			err := rows.Scan(&resourceID, pq.Array(&operationIDs))
			if err != nil {
				return err
			}

			// Create permission
			permissionID := uuid.New()
			_, err = s.db.ExecContext(ctx, `
				INSERT INTO auth.permissions (id, resource_id, operations, effect, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)`,
				permissionID, resourceID, pq.Array(operationIDs), "allow", time.Now(), time.Now())
			if err != nil {
				return err
			}

			// Assign permission to group
			_, err = s.db.ExecContext(ctx, `
				INSERT INTO auth.group_permissions (id, group_id, permission_id, created_at)
				VALUES ($1, $2, $3, $4)`,
				uuid.New(), groupID, permissionID, time.Now())
			if err != nil {
				return err
			}
		}

		s.logger.Println("Created super admin group with all permissions")
	}

	return nil
}

// createAdminWalletIfNotExists creates an admin wallet for the super admin user if not exists
func (s *AuthSeeder) createAdminWalletIfNotExists(ctx context.Context, userID uuid.UUID) error {
	s.logger.Println("Checking admin wallet...")

	// Check how many super_admin wallets exist
	var walletCount int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM merchant.wallet WHERE wallet_type = 'super_admin'`).Scan(&walletCount)
	if err != nil {
		return fmt.Errorf("failed to check admin wallet count: %w", err)
	}

	if walletCount > 1 {
		s.logger.Printf("WARNING: Found %d admin wallets, expected exactly 1. This may indicate data inconsistency.", walletCount)
		return nil // Don't fail, just warn
	}

	if walletCount == 1 {
		s.logger.Println("Admin wallet already exists, skipping creation")
		return nil
	}

	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Println("Super admin user not found, skipping admin wallet creation")
			return nil // Don't fail if super admin doesn't exist yet
		}
		return fmt.Errorf("failed to get super admin user ID: %w", err)
	}

	s.logger.Printf("Creating admin wallet for super admin user ID: %s", userID)

	// Create admin wallet with default values
	walletID := uuid.New()
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO merchant.wallet (id, user_id, merchant_id, amount, locked_amount, currency, wallet_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		walletID,
		userID,
		nil,           // merchant_id is NULL for admin wallets
		0.0,           // initial amount
		0.0,           // initial locked_amount
		"ETB",         // default currency
		"super_admin", // wallet_type
		time.Now(),
		time.Now())

	if err != nil {
		return fmt.Errorf("failed to create admin wallet: %w", err)
	}

	s.logger.Printf("Admin wallet created successfully with ID: %s for super admin user: %s", walletID, userID)
	return nil
}
