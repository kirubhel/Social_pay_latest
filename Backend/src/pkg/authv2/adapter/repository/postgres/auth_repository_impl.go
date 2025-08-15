package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/authv2/core/repository"
	"github.com/socialpay/socialpay/src/pkg/authv2/utils"
	"github.com/lib/pq"
)

// AuthRepositoryImpl implements the AuthRepository interface
type AuthRepositoryImpl struct {
	db     *sql.DB
	logger *log.Logger
}

// NewAuthRepository creates a new PostgreSQL auth repository
func NewAuthRepository(db *sql.DB, logger *log.Logger) repository.AuthRepository {
	return &AuthRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// CreateUser creates a new user with phone and password
func (r *AuthRepositoryImpl) CreateUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with phone data directly
	userID := uuid.New()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO auth.users (id, email, sir_name, first_name, last_name, user_type, phone_prefix, phone_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		userID, req.Email, req.Title, req.FirstName, req.LastName, req.UserType, req.PhonePrefix, req.PhoneNumber, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create password identity
	_, err = tx.ExecContext(ctx, `
		INSERT INTO auth.password_identities (id, user_id, password, hint, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), userID, hashedPassword, req.PasswordHint, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create password identity: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return created user
	user := &entity.User{
		ID:          userID,
		SirName:     req.Title,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		UserType:    req.UserType,
		PhonePrefix: req.PhonePrefix,
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return user, nil
}

// GetUserByPhone gets user by phone number
func (r *AuthRepositoryImpl) GetUserByPhone(ctx context.Context, prefix, number string) (*entity.User, error) {
	var user entity.User

	query := `
		SELECT id, sir_name, first_name, last_name, user_type, phone_prefix, phone_number, created_at, updated_at
		FROM auth.users
		WHERE phone_prefix = $1 AND phone_number = $2`

	err := r.db.QueryRowContext(ctx, query, prefix, number).Scan(
		&user.ID, &user.SirName, &user.FirstName, &user.LastName, &user.UserType,
		&user.PhonePrefix, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

// GetUserByID gets user by ID
func (r *AuthRepositoryImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	var user entity.User

	query := `SELECT id, email, sir_name, first_name, last_name, user_type, phone_prefix, phone_number, created_at, updated_at FROM auth.users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.SirName, &user.FirstName, &user.LastName, &user.UserType, &user.PhonePrefix, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// RemoveUser deletes user from db
func (r *AuthRepositoryImpl) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM auth.password_identities WHERE user_id = $1;`, userID)

	if err != nil {
		return fmt.Errorf("failed to delete password_identities: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM auth.users WHERE id = $1;`, userID)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// VerifyPassword verifies user password
func (r *AuthRepositoryImpl) VerifyPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error) {
	var hashedPassword string
	query := `SELECT password FROM auth.password_identities WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to get password: %w", err)
	}

	return utils.CheckPasswordHash(password, hashedPassword), nil
}

// UserExists checks if user exists by phone
func (r *AuthRepositoryImpl) UserExists(ctx context.Context, phonePrefix, phoneNumber string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth.users WHERE phone_prefix = $1 AND phone_number = $2)`
	err := r.db.QueryRowContext(ctx, query, phonePrefix, phoneNumber).Scan(&exists)
	return exists, err
}

// CreateMerchant creates a merchant record
func (r *AuthRepositoryImpl) CreateMerchant(ctx context.Context, userID uuid.UUID, merchantData map[string]interface{}) (*uuid.UUID, error) {
	merchantID := uuid.New()

	// Extract values from merchantData with defaults
	legalName, _ := merchantData["legal_name"].(string)
	tradingName, _ := merchantData["trading_name"].(string)
	businessType, _ := merchantData["business_type"].(string)
	businessRegNumber, _ := merchantData["business_registration_number"].(string)
	taxID, _ := merchantData["tax_identification_number"].(string)
	industryCategory, _ := merchantData["industry_category"].(string)
	isBettingCompany, _ := merchantData["is_betting_company"].(bool)
	lotteryCertNumber, _ := merchantData["lottery_certificate_number"].(string)
	websiteURL, _ := merchantData["website_url"].(string)
	establishedDate, _ := merchantData["established_date"].(time.Time)
	status, _ := merchantData["status"].(string)

	// Get user data for generating defaults
	title, _ := merchantData["title"].(string)
	firstName, _ := merchantData["first_name"].(string)
	lastName, _ := merchantData["last_name"].(string)

	// Set defaults if values are empty
	if legalName == "" {
		legalName = fmt.Sprintf("%s %s %s Business PLC", title, firstName, lastName)
	}
	if tradingName == "" {
		tradingName = fmt.Sprintf("%s %s Enterprises", firstName, lastName)
	}
	if businessRegNumber == "" {
		// Use UUID short form + timestamp + large random number for uniqueness
		uuidShort := strings.Replace(uuid.New().String()[:8], "-", "", -1)
		businessRegNumber = fmt.Sprintf("ET%s%s%06d", time.Now().Format("20060102"), uuidShort, rand.Intn(999999))
	}
	if taxID == "" {
		// Use UUID short form + timestamp + large random number for uniqueness
		uuidShort := strings.Replace(uuid.New().String()[:8], "-", "", -1)
		taxID = fmt.Sprintf("TIN%s%s%06d", time.Now().Format("20060102"), uuidShort, rand.Intn(999999))
	}
	if industryCategory == "" {
		industryCategory = "Retail Trade"
	}
	if businessType == "" {
		businessType = "Sole Proprietorship"
	}
	if lotteryCertNumber == "" {
		lotteryCertNumber = ""
	}
	if websiteURL == "" {
		websiteURL = fmt.Sprintf("www.%s%s.com", strings.ToLower(firstName), strings.ToLower(lastName))
	}
	if establishedDate.IsZero() {
		establishedDate = time.Now().AddDate(-3, 0, 0)
	}
	if status == "" {
		status = "inactive"
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO merchants.merchants (
			id, 
			user_id, 
			legal_name, 
			trading_name, 
			business_registration_number, 
			tax_identification_number, 
			industry_category, 
			business_type, 
			is_betting_company, 
			lottery_certificate_number, 
			website_url, 
			established_date, 
			created_at, 
			updated_at, 
			status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`,
		merchantID,
		userID,
		legalName,
		tradingName,
		businessRegNumber,
		taxID,
		industryCategory,
		businessType,
		isBettingCompany,
		lotteryCertNumber,
		websiteURL,
		establishedDate,
		time.Now(),
		time.Now(),
		status)

	if err != nil {
		return nil, fmt.Errorf("failed to create merchant: %w", err)
	}

	return &merchantID, nil
}

func (r *AuthRepositoryImpl) GetMerchant(ctx context.Context, id uuid.UUID) (*entity.Merchant, error) {
	query := `SELECT id, user_id, legal_name, trading_name, business_registration_number, tax_identification_number, business_type, industry_category, is_betting_company, lottery_certificate_number, website_url, established_date, created_at, updated_at, status FROM merchants.merchants WHERE id = $1 AND deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, query, id)

	var merchant entity.Merchant
	err := row.Scan(
		&merchant.ID,
		&merchant.UserID,
		&merchant.LegalName,
		&merchant.TradingName,
		&merchant.BusinessRegistrationNumber,
		&merchant.TaxIdentificationNumber,
		&merchant.BusinessType,
		&merchant.IndustryCategory,
		&merchant.IsBettingCompany,
		&merchant.LotteryCertificateNumber,
		&merchant.WebsiteURL,
		&merchant.EstablishedDate,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
		&merchant.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan merchant row: %w", err)
	}

	return &merchant, nil
}

// CreateWallet creates a wallet for user
func (r *AuthRepositoryImpl) CreateWallet(ctx context.Context, userID uuid.UUID, walletType string) error {
	walletID := uuid.New()

	// For admin wallets, merchant_id is NULL
	var merchantID *uuid.UUID
	if walletType == "merchant" {
		// Get merchant ID for the user
		var mID uuid.UUID
		err := r.db.QueryRowContext(ctx, `SELECT id FROM merchants.merchants WHERE user_id = $1`, userID).Scan(&mID)
		if err == nil {
			merchantID = &mID
		}
	}

	// Create wallet in merchants.wallets table (not admin.wallet)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO merchant.wallet (id, user_id, merchant_id, amount, locked_amount, currency, wallet_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		walletID, userID, merchantID, 0.0, 0.0, "ETB", walletType, time.Now(), time.Now())

	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

// GetOrCreateSuperAdminGroup gets or creates the super admin group
func (r *AuthRepositoryImpl) GetOrCreateSuperAdminGroup(ctx context.Context) (*entity.Group, error) {
	// Try to get existing super admin group
	var group entity.Group
	query := `SELECT id, title, created_at, updated_at FROM auth.groups WHERE title = 'Super Admin' LIMIT 1`
	err := r.db.QueryRowContext(ctx, query).Scan(&group.ID, &group.Title, &group.CreatedAt, &group.UpdatedAt)

	if err == nil {
		return &group, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query super admin group: %w", err)
	}

	// Create super admin group if it doesn't exist
	groupID := uuid.New()
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO auth.groups (id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`,
		groupID, "Super Admin", time.Now(), time.Now())

	if err != nil {
		return nil, fmt.Errorf("failed to create super admin group: %w", err)
	}

	group.ID = groupID
	group.Title = "Super Admin"
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	// Create and assign ADMIN_ALL permissions to super admin group
	adminAllPermission, err := r.GetOrCreatePermission(ctx, string(entity.RESOURCE_ADMIN_ALL), string(entity.OPERATION_ADMIN_ALL), "allow")
	if err != nil {
		return nil, fmt.Errorf("failed to create ADMIN_ALL permission: %w", err)
	}

	err = r.AssignPermissionToGroup(ctx, groupID, adminAllPermission.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign ADMIN_ALL permission to super admin group: %w", err)
	}

	return &group, nil
}

// AssignUserToGroup assigns a user to a group with optional merchant context
func (r *AuthRepositoryImpl) AssignUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	// Get the group's merchant_id
	var merchantID *uuid.UUID
	err := r.db.QueryRowContext(ctx, `SELECT merchant_id FROM auth.groups WHERE id = $1`, groupID).Scan(&merchantID)
	if err != nil {
		return fmt.Errorf("failed to get group merchant_id: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO auth.user_groups (id, user_id, group_id, merchant_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, group_id) DO NOTHING`,
		uuid.New(), userID, groupID, merchantID, time.Now(), time.Now())

	return err
}

// CheckUserPermission checks if user has a specific permission
func (r *AuthRepositoryImpl) CheckUserPermission(ctx context.Context, userID uuid.UUID, resource, operation string) (bool, error) {
	var hasPermission bool

	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM auth.user_groups ug
			JOIN auth.group_permissions gp ON ug.group_id = gp.group_id
			JOIN auth.permissions p ON gp.permission_id = p.id
			JOIN auth.resources r ON p.resource_id = r.id
			JOIN auth.operations o ON o.id = ANY(p.operations)
			WHERE ug.user_id = $1 
			AND (
				-- Exact resource and operation match
				(r.name = $2 AND o.name = $3)
				OR 
				-- ALL resource and ALL operation (merchant permissions)
				(r.name = 'ALL' AND o.name = 'ALL')
				OR
				-- ADMIN_ALL resource and ADMIN_ALL operation (super admin permissions)
				(r.name = 'ADMIN_ALL' AND o.name = 'ADMIN_ALL')
			)
			AND p.effect = 'allow'
		)`

	err := r.db.QueryRowContext(ctx, query, userID, resource, operation).Scan(&hasPermission)
	return hasPermission, err
}

// LogAuthActivity logs authentication activity
func (r *AuthRepositoryImpl) LogAuthActivity(ctx context.Context, activity *entity.AuthActivity) error {
	detailsJSON, _ := json.Marshal(activity.Details)

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth.auth_activities (id, user_id, activity_type, ip_address, user_agent, device_name, success, details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		activity.ID, activity.UserID, activity.ActivityType, activity.IPAddress,
		activity.UserAgent, activity.DeviceName, activity.Success, detailsJSON, activity.CreatedAt)

	return err
}

// UpdateUser updates user information excluding user_type and device_info
func (r *AuthRepositoryImpl) UpdateUser(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Build dynamic update query for auth.users table
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	// Map of allowed fields to update (excluding user_type and device_info)
	allowedFields := map[string]string{
		"title":        "sir_name",
		"sir_name":     "sir_name",
		"first_name":   "first_name",
		"last_name":    "last_name",
		"email":        "email",
		"phone_prefix": "phone_prefix",
		"phone_number": "phone_number",
	}

	for field, value := range updates {
		if dbField, allowed := allowedFields[field]; allowed && value != nil {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", dbField, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(setParts) > 0 {
		setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
		args = append(args, time.Now())
		argIndex++

		query := fmt.Sprintf("UPDATE auth.users SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
		args = append(args, userID)

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Handle password update separately
	if password, ok := updates["password"]; ok && password != nil {
		if passwordStr, ok := password.(string); ok && passwordStr != "" {
			hashedPassword, err := utils.HashPassword(passwordStr)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}

			_, err = tx.ExecContext(ctx, `
				UPDATE auth.password_identities SET password = $1, updated_at = $2 WHERE user_id = $3`,
				hashedPassword, time.Now(), userID)
			if err != nil {
				return fmt.Errorf("failed to update password: %w", err)
			}
		}
	}

	// Handle password hint update separately
	if passwordHint, ok := updates["password_hint"]; ok && passwordHint != nil {
		_, err = tx.ExecContext(ctx, `
			UPDATE auth.password_identities SET hint = $1, updated_at = $2 WHERE user_id = $3`,
			passwordHint, time.Now(), userID)
		if err != nil {
			return fmt.Errorf("failed to update password hint: %w", err)
		}
	}

	return tx.Commit()
}

func (r *AuthRepositoryImpl) CreatePhone(ctx context.Context, prefix, number string) (*entity.Phone, error) {
	phoneID := uuid.New()
	phone := &entity.Phone{
		ID:        phoneID,
		Prefix:    prefix,
		Number:    number,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth.phones (id, prefix, number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		phoneID, prefix, number, phone.CreatedAt, phone.UpdatedAt)

	return phone, err
}

func (r *AuthRepositoryImpl) GetPhoneByID(ctx context.Context, phoneID uuid.UUID) (*entity.Phone, error) {
	var phone entity.Phone
	query := `SELECT id, prefix, number, created_at, updated_at FROM auth.phones WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, phoneID).Scan(
		&phone.ID, &phone.Prefix, &phone.Number, &phone.CreatedAt, &phone.UpdatedAt)
	return &phone, err
}

func (r *AuthRepositoryImpl) GetPhoneByNumber(ctx context.Context, prefix, number string) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, phone_prefix, phone_number, created_at, updated_at FROM auth.users WHERE phone_prefix = $1 AND phone_number = $2`
	err := r.db.QueryRowContext(ctx, query, prefix, number).Scan(
		&user.ID, &user.PhonePrefix, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func (r *AuthRepositoryImpl) LinkPhoneToUser(ctx context.Context, userID, phoneID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth.phone_identities (id, user_id, phone_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), userID, phoneID, time.Now(), time.Now())
	return err
}

func (r *AuthRepositoryImpl) CreatePasswordIdentity(ctx context.Context, userID uuid.UUID, hashedPassword, hint string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth.password_identities (id, user_id, password, hint, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), userID, hashedPassword, hint, time.Now(), time.Now())
	return err
}

func (r *AuthRepositoryImpl) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE auth.password_identities SET password = $1, updated_at = $2 WHERE user_id = $3`,
		hashedPassword, time.Now(), userID)
	return err
}

func (r *AuthRepositoryImpl) CreateOTP(ctx context.Context, userID uuid.UUID, code, token string, expiresAt int64) (*entity.OTPRequest, error) {
	otpID := uuid.New()
	otp := &entity.OTPRequest{
		ID:        otpID,
		UserID:    userID,
		Code:      code,
		Token:     token,
		Method:    "sms",
		Status:    false,
		ExpiresAt: time.Unix(expiresAt, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth.phone_auths (id, user_id, code, token, method, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		otpID, userID, code, token, "sms", false, otp.CreatedAt, otp.UpdatedAt)

	return otp, err
}

func (r *AuthRepositoryImpl) GetOTPByToken(ctx context.Context, token string) (*entity.OTPRequest, error) {
	var otp entity.OTPRequest
	query := `SELECT id, user_id, code, token, method, status, created_at, updated_at FROM auth.phone_auths WHERE token = $1 ORDER BY created_at DESC LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&otp.ID, &otp.UserID, &otp.Code, &otp.Token, &otp.Method, &otp.Status, &otp.CreatedAt, &otp.UpdatedAt)

	// Set expires at to 5 minutes after creation
	otp.ExpiresAt = otp.CreatedAt.Add(5 * time.Minute)

	return &otp, err
}

func (r *AuthRepositoryImpl) MarkOTPAsUsed(ctx context.Context, otpID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE auth.phone_auths SET status = true, updated_at = $1 WHERE id = $2`,
		time.Now().UTC(), otpID)
	return err
}

func (r *AuthRepositoryImpl) CreateSession(ctx context.Context, userID, deviceID uuid.UUID, token, refreshToken string, expiresAt int64) (*entity.Session, error) {
	sessionID := uuid.New()
	session := &entity.Session{
		ID:           sessionID,
		UserID:       userID,
		DeviceID:     deviceID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Unix(expiresAt, 0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth.sessions (id, user_id, device_id, token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		sessionID, userID, deviceID, token, session.CreatedAt, session.UpdatedAt)

	return session, err
}

func (r *AuthRepositoryImpl) GetSessionByToken(ctx context.Context, token string) (*entity.Session, error) {
	// TODO: Implement if needed
	return nil, nil
}

func (r *AuthRepositoryImpl) UpdateSessionToken(ctx context.Context, sessionID uuid.UUID, token, refreshToken string, expiresAt int64) error {
	// TODO: Implement if needed
	return nil
}

func (r *AuthRepositoryImpl) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM auth.sessions WHERE id = $1`, sessionID)
	return err
}

func (r *AuthRepositoryImpl) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM auth.sessions WHERE user_id = $1`, userID)
	return err
}

func (r *AuthRepositoryImpl) CreateDevice(ctx context.Context, ip, name, agent string) (*entity.Device, error) {
	deviceID := uuid.New()
	device := &entity.Device{
		ID:        deviceID,
		IP:        ip,
		Name:      name,
		Agent:     agent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth.devices (id, ip, name, agent, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		deviceID, ip, name, agent, device.CreatedAt, device.UpdatedAt)

	return device, err
}

func (r *AuthRepositoryImpl) GetDeviceByFingerprint(ctx context.Context, ip, name, agent string) (*entity.Device, error) {
	var device entity.Device
	query := `SELECT id, ip, name, agent, created_at, updated_at FROM auth.devices WHERE ip = $1 AND name = $2 AND agent = $3 LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, ip, name, agent).Scan(
		&device.ID, &device.IP, &device.Name, &device.Agent, &device.CreatedAt, &device.UpdatedAt)
	return &device, err
}

func (r *AuthRepositoryImpl) GetUserGroups(ctx context.Context, userID uuid.UUID) ([]entity.Group, error) {
	var groups []entity.Group

	query := `SELECT id, title, permissions, created_at, updated_at, merchant_id, description FROM auth."groups" WHERE user_id=$1;`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group entity.Group
		var uuidPerms []uuid.UUID

		err := rows.Scan(
			&group.ID,
			&group.Title,
			pq.Array(&uuidPerms),
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.MerchantID,
			&group.Description,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to parse group: %w", err)
		}

		groups = append(groups, group)

	}

	return []entity.Group{}, nil
}

func (r *AuthRepositoryImpl) RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM auth.user_groups WHERE user_id = $1 AND group_id = $2`, userID, groupID)
	return err
}

// UpdateUserGroup updates a user's group assignment by changing from old group to new group
func (r *AuthRepositoryImpl) UpdateUserGroup(ctx context.Context, userID, oldGroupID, newGroupID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the new group's merchant_id
	var newMerchantID *uuid.UUID
	err = tx.QueryRowContext(ctx, `SELECT merchant_id FROM auth.groups WHERE id = $1`, newGroupID).Scan(&newMerchantID)
	if err != nil {
		return fmt.Errorf("failed to get new group merchant_id: %w", err)
	}

	// Update the user_groups record
	_, err = tx.ExecContext(ctx, `
		UPDATE auth.user_groups 
		SET group_id = $1, merchant_id = $2, updated_at = $3 
		WHERE user_id = $4 AND group_id = $5`,
		newGroupID, newMerchantID, time.Now(), userID, oldGroupID)

	if err != nil {
		return fmt.Errorf("failed to update user group: %w", err)
	}

	return tx.Commit()
}

// CreateGroup creates a new group with optional merchant_id
func (r *AuthRepositoryImpl) CreateGroup(ctx context.Context, title string, description *string, merchantID *uuid.UUID) (*entity.Group, error) {
	groupID := uuid.New()
	group := &entity.Group{
		ID:          groupID,
		Title:       title,
		Description: description,
		MerchantID:  merchantID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth.groups (id, title, description, merchant_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		groupID, title, description, merchantID, group.CreatedAt, group.UpdatedAt)

	return group, err
}

func (r *AuthRepositoryImpl) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*entity.Group, error) {
	var group entity.Group
	query := `SELECT id, title, created_at, updated_at FROM auth.groups WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(
		&group.ID, &group.Title, &group.CreatedAt, &group.UpdatedAt)
	return &group, err
}

func (r *AuthRepositoryImpl) GetGroupsByMerchant(ctx context.Context, merchantID *uuid.UUID) ([]entity.Group, error) {
	// TODO: Implement if needed
	return []entity.Group{}, nil
}

func (r *AuthRepositoryImpl) AssignPermissionToGroup(ctx context.Context, groupID, permissionID uuid.UUID) error {
	// Check if the permission is already assigned to avoid duplicates
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM auth.group_permissions WHERE group_id = $1 AND permission_id = $2)`
	err := r.db.QueryRowContext(ctx, checkQuery, groupID, permissionID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check existing permission assignment: %w", err)
	}

	if exists {
		return nil // Permission already assigned, nothing to do
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO auth.group_permissions (id, group_id, permission_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), groupID, permissionID, time.Now(), time.Now())
	return err
}

func (r *AuthRepositoryImpl) GetPermissionByResourceOperation(ctx context.Context, resource, operation string) (*entity.Permission, error) {
	// TODO: Implement if needed
	return nil, nil
}

func (r *AuthRepositoryImpl) GetUserActivities(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.AuthActivity, error) {
	// TODO: Implement if needed
	return []entity.AuthActivity{}, nil
}

// GetUserPermissions gets all permissions for a user
func (r *AuthRepositoryImpl) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `
		SELECT DISTINCT r.name || ':' || o.name as permission
		FROM auth.user_groups ug
		JOIN auth.group_permissions gp ON ug.group_id = gp.group_id
		JOIN auth.permissions p ON gp.permission_id = p.id
		JOIN auth.resources r ON p.resource_id = r.id
		JOIN auth.operations o ON o.id = ANY(p.operations)
		WHERE ug.user_id = $1 
		AND p.effect = 'allow'
		ORDER BY permission`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user permissions: %w", err)
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetUserPermissionsByMerchant gets all permissions for a user within a specific merchant context
func (r *AuthRepositoryImpl) GetUserPermissionsByMerchant(ctx context.Context, userID, merchantID uuid.UUID) ([]string, error) {
	query := `
		SELECT DISTINCT r.name || ':' || o.name as permission
		FROM auth.user_groups ug
		JOIN auth.group_permissions gp ON ug.group_id = gp.group_id
		JOIN auth.permissions p ON gp.permission_id = p.id
		JOIN auth.resources r ON p.resource_id = r.id
		JOIN auth.operations o ON o.id = ANY(p.operations)
		WHERE ug.user_id = $1 
		AND ug.merchant_id = $2
		AND p.effect = 'allow'
		ORDER BY permission`

	rows, err := r.db.QueryContext(ctx, query, userID, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user permissions by merchant: %w", err)
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetAllUserPermissions gets all permissions for a user grouped by merchant ID
func (r *AuthRepositoryImpl) GetAllUserPermissions(ctx context.Context, userID uuid.UUID) (map[string][]string, error) {
	query := `
		SELECT 
			COALESCE(ug.merchant_id::text, 'global') as merchant_id,
			r.name || ':' || o.name as permission
		FROM auth.user_groups ug
		JOIN auth.group_permissions gp ON ug.group_id = gp.group_id
		JOIN auth.permissions p ON gp.permission_id = p.id
		JOIN auth.resources r ON p.resource_id = r.id
		JOIN auth.operations o ON o.id = ANY(p.operations)
		WHERE ug.user_id = $1 
		AND p.effect = 'allow'
		ORDER BY merchant_id, permission`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query all user permissions: %w", err)
	}
	defer rows.Close()

	permissionsByMerchant := make(map[string][]string)
	for rows.Next() {
		var merchantID, permission string
		err := rows.Scan(&merchantID, &permission)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}

		if permissionsByMerchant[merchantID] == nil {
			permissionsByMerchant[merchantID] = []string{}
		}
		permissionsByMerchant[merchantID] = append(permissionsByMerchant[merchantID], permission)
	}

	return permissionsByMerchant, nil
}

// CheckUserPermissionForMerchant checks if user has a specific permission within a merchant context
func (r *AuthRepositoryImpl) CheckUserPermissionForMerchant(ctx context.Context, userID, merchantID uuid.UUID, resource, operation string) (bool, error) {
	var hasPermission bool

	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM auth.user_groups ug
			JOIN auth.group_permissions gp ON ug.group_id = gp.group_id
			JOIN auth.permissions p ON gp.permission_id = p.id
			JOIN auth.resources r ON p.resource_id = r.id
			JOIN auth.operations o ON o.id = ANY(p.operations)
			WHERE ug.user_id = $1 
			AND ug.merchant_id = $2
			AND (
				-- Exact resource and operation match
				(r.name = $3 AND o.name = $4)
				OR 
				-- ALL resource and ALL operation (merchant permissions)
				(r.name = 'ALL' AND o.name = 'ALL')
				OR
				-- ADMIN_ALL resource and ADMIN_ALL operation (super admin permissions)
				(r.name = 'ADMIN_ALL' AND o.name = 'ADMIN_ALL')
			)
			AND p.effect = 'allow'
		)`

	err := r.db.QueryRowContext(ctx, query, userID, merchantID, resource, operation).Scan(&hasPermission)
	return hasPermission, err
}

// GetUserGroupsByMerchant gets all groups for a user within a specific merchant context
func (r *AuthRepositoryImpl) GetUserGroupsByMerchant(ctx context.Context, userID, merchantID uuid.UUID) ([]entity.Group, error) {
	query := `
		SELECT g.id, g.title, g.description, g.merchant_id, g.created_at, g.updated_at
		FROM auth.groups g
		JOIN auth.user_groups ug ON g.id = ug.group_id
		WHERE ug.user_id = $1 AND ug.merchant_id = $2`

	rows, err := r.db.QueryContext(ctx, query, userID, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user groups by merchant: %w", err)
	}
	defer rows.Close()

	var groups []entity.Group
	for rows.Next() {
		var group entity.Group
		err := rows.Scan(&group.ID, &group.Title, &group.Description,
			&group.MerchantID, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// GetAllUserGroups gets all groups for a user grouped by merchant ID
func (r *AuthRepositoryImpl) GetAllUserGroups(ctx context.Context, userID uuid.UUID) (map[string][]entity.Group, error) {
	query := `
		SELECT 
		    ug.merchant_id,
		    g.id, g.title, g.description, g.merchant_id, g.created_at, g.updated_at
		FROM auth.user_groups ug
		JOIN auth.groups g ON ug.group_id = g.id
		WHERE ug.user_id = $1
		ORDER BY ug.merchant_id, g.title`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query all user groups: %w", err)
	}
	defer rows.Close()

	groupsByMerchant := make(map[string][]entity.Group)
	for rows.Next() {
		var merchantID uuid.UUID
		var group entity.Group
		err := rows.Scan(&merchantID, &group.ID, &group.Title, &group.Description,
			&group.MerchantID, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}

		if groupsByMerchant[merchantID.String()] == nil {
			groupsByMerchant[merchantID.String()] = []entity.Group{}
		}
		groupsByMerchant[merchantID.String()] = append(groupsByMerchant[merchantID.String()], group)
	}

	return groupsByMerchant, nil
}

// CreatePermission creates a new permission
func (r *AuthRepositoryImpl) CreatePermission(ctx context.Context, resourceName, operation, effect string) (*entity.Permission, error) {
	// Get or create resource
	var resourceID uuid.UUID
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM auth.resources WHERE name = $1 LIMIT 1`,
		resourceName).Scan(&resourceID)

	if err == sql.ErrNoRows {
		// Create new resource if it doesn't exist
		resourceID = uuid.New()
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO auth.resources (id, name, created_at, updated_at)
			VALUES ($1, $2, $3, $4)`,
			resourceID, resourceName, time.Now(), time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to create resource: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to query resource: %w", err)
	}

	// Get or create operation
	var operationID uuid.UUID
	err = r.db.QueryRowContext(ctx, `
		SELECT id FROM auth.operations WHERE name = $1 LIMIT 1`,
		operation).Scan(&operationID)

	if err == sql.ErrNoRows {
		// Create new operation if it doesn't exist
		operationID = uuid.New()
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO auth.operations (id, name, created_at, updated_at)
			VALUES ($1, $2, $3, $4)`,
			operationID, operation, time.Now(), time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to create operation: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to query operation: %w", err)
	}

	permissionID := uuid.New()
	permission := &entity.Permission{
		ID:           permissionID,
		ResourceName: resourceName,
		Operations:   []uuid.UUID{operationID},
		Effect:       effect,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO auth.permissions (id, resource_id, operations, effect, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		permissionID, resourceID, pq.Array([]uuid.UUID{operationID}), effect, permission.CreatedAt, permission.UpdatedAt)

	return permission, err
}

// GetOrCreatePermission gets or creates a permission
func (r *AuthRepositoryImpl) GetOrCreatePermission(ctx context.Context, resourceName, operation, effect string) (*entity.Permission, error) {
	// Get resource ID
	var resourceID uuid.UUID
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM auth.resources WHERE name = $1 LIMIT 1`,
		resourceName).Scan(&resourceID)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query resource: %w", err)
	}

	// Get operation ID
	var operationID uuid.UUID
	err = r.db.QueryRowContext(ctx, `
		SELECT id FROM auth.operations WHERE name = $1 LIMIT 1`,
		operation).Scan(&operationID)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query operation: %w", err)
	}

	// Try to get existing permission
	var permission entity.Permission
	query := `SELECT p.id, p.operations, p.effect, p.created_at, p.updated_at 
			  FROM auth.permissions p
			  WHERE p.resource_id = $1 AND $2 = ANY(p.operations) AND p.effect = $3 LIMIT 1`

	err = r.db.QueryRowContext(ctx, query, resourceID, operationID, effect).Scan(
		&permission.ID, pq.Array(&permission.Operations), &permission.Effect,
		&permission.CreatedAt, &permission.UpdatedAt)

	if err == nil {
		permission.ResourceName = resourceName
		return &permission, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query permission: %w", err)
	}

	// Create new permission if it doesn't exist
	return r.CreatePermission(ctx, resourceName, operation, effect)
}

// GetMerchantsByUser gets all merchants for a user
func (r *AuthRepositoryImpl) GetMerchantsByUser(ctx context.Context, userID uuid.UUID) ([]entity.Merchant, error) {
	query := `
		SELECT id, user_id, legal_name, trading_name, business_registration_number, 
			   tax_identification_number, industry_category, business_type, 
			   is_betting_company, lottery_certificate_number, website_url, 
			   established_date, status, created_at, updated_at
		FROM merchants.merchants 
		WHERE user_id = $1 AND deleted_at IS NULL`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query merchants: %w", err)
	}
	defer rows.Close()

	var merchants []entity.Merchant
	for rows.Next() {
		var merchant entity.Merchant
		err := rows.Scan(
			&merchant.ID, &merchant.UserID, &merchant.LegalName, &merchant.TradingName,
			&merchant.BusinessRegistrationNumber, &merchant.TaxIdentificationNumber,
			&merchant.IndustryCategory, &merchant.BusinessType, &merchant.IsBettingCompany,
			&merchant.LotteryCertificateNumber, &merchant.WebsiteURL, &merchant.EstablishedDate,
			&merchant.Status, &merchant.CreatedAt, &merchant.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan merchant: %w", err)
		}
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

// GetOrCreateMerchantOwnerGroup gets or creates the merchant owner group for a specific merchant
func (r *AuthRepositoryImpl) GetOrCreateMerchantOwnerGroup(ctx context.Context, merchantID uuid.UUID) (*entity.Group, error) {
	// Try to get existing merchant owner group
	var group entity.Group
	query := `SELECT id, title, description, merchant_id, created_at, updated_at 
			  FROM auth.groups 
			  WHERE title = $1 AND merchant_id = $2 LIMIT 1`

	groupTitle := fmt.Sprintf("Merchant Owner - %s", merchantID.String()[:8])
	err := r.db.QueryRowContext(ctx, query, groupTitle, merchantID).Scan(
		&group.ID, &group.Title, &group.Description, &group.MerchantID, &group.CreatedAt, &group.UpdatedAt)

	if err == nil {
		return &group, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query merchant owner group: %w", err)
	}

	// Create merchant owner group if it doesn't exist
	groupID := uuid.New()
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO auth.groups (id, title, description, merchant_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		groupID, groupTitle, "Owner group for merchant", &merchantID, time.Now(), time.Now())

	if err != nil {
		return nil, fmt.Errorf("failed to create merchant owner group: %w", err)
	}

	description := "Owner group for merchant"

	group.ID = groupID
	group.Title = groupTitle
	group.Description = &description
	group.MerchantID = &merchantID
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	// Create and assign ALL permissions to merchant owner group
	allPermission, err := r.GetOrCreatePermission(ctx, string(entity.RESOURCE_ALL), string(entity.OPERATION_ALL), "allow")
	if err != nil {
		return nil, fmt.Errorf("failed to create ALL permission: %w", err)
	}

	err = r.AssignPermissionToGroup(ctx, groupID, allPermission.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign ALL permission to merchant owner group: %w", err)
	}

	return &group, nil
}
