package repo

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"

	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrPasswordHashingFailed  = "ErrPasswordHashingFailed"
	MsgPasswordHashingFailed  = "Password hashing failed"
	ErrUserCreationFailed     = "ErrUserCreationFailed"
	MsgUserCreationFailed     = "User creation failed"
	ErrPhoneCreationFailed    = "ErrPhoneCreationFailed"
	MsgPhoneCreationFailed    = "Phone creation failed"
	ErrPhoneIdentityFailed    = "ErrPhoneIdentityFailed"
	MsgPhoneIdentityFailed    = "Phone identity creation failed"
	ErrPasswordIdentityFailed = "ErrPasswordIdentityFailed"
	MsgPasswordIdentityFailed = "Password identity creation failed"
	ErrCommitFailed           = "ErrCommitFailed"
	MsgCommitFailed           = "Transaction commit failed"
)

func (repo PsqlRepo) CreateUser(
    Title string,
    FirstName string,
    LastName string,
    PhonePrefix string,
    PhoneNumber string,
    Password string,
    PasswordHint string,
    UserType string,
) (*entity.User, error) {

    tx, err := repo.db.Begin()
    if err != nil {
        repo.log.Printf("Failed to begin transaction: %v", err)
        return nil, &entity.Error{
            Type:    entity.ErrAccountCreation,
            Message: entity.MsgAccountCreation,
            Detail:  fmt.Sprintf("failed to begin transaction: %v", err),
        }
    }

    defer func() {
        if err != nil {
            repo.log.Printf("Rolling back transaction due to error: %v", err)
            tx.Rollback()
        }
    }()

    // Phone number validation
    if len(PhoneNumber) < 9 || len(PhoneNumber) > 15 {
        return nil, &entity.Error{
            Type:    entity.ErrInvalidPhoneNumberFormat,
            Message: entity.ErrInvalidPhoneNumber,
            Detail:  "phone number must be between 9 digits long like +251911234567",
        }
    }

    // Check for orphan phone and delete if no user linked
    var phoneIDToDelete uuid.UUID
    err = tx.QueryRow(
        `SELECT p.id
        FROM auth.phones p
        LEFT JOIN auth.phone_identities pi ON p.id = pi.phone_id
        WHERE p.prefix = $1 AND p.number = $2 AND pi.id IS NULL`,
        PhonePrefix, PhoneNumber).Scan(&phoneIDToDelete)

    if err != nil && err != sql.ErrNoRows {
        repo.log.Printf("Failed to check orphan phone: %v", err)
        return nil, &entity.Error{
            Type:    entity.ErrAccountCreation,
            Message: entity.MsgAccountCreation,
            Detail:  fmt.Sprintf("orphan phone lookup failed: %v", err),
        }
    }

    if err == nil {
        _, err = tx.Exec(`DELETE FROM auth.phones WHERE id = $1`, phoneIDToDelete)
        if err != nil {
            repo.log.Printf("Failed to delete orphan phone: %v", err)
            return nil, &entity.Error{
                Type:    entity.ErrAccountCreation,
                Message: entity.MsgAccountCreation,
                Detail:  fmt.Sprintf("orphan phone deletion failed: %v", err),
            }
        }
    }

    // Phone existence check
    var phoneExists bool
    err = tx.QueryRow(
        `SELECT EXISTS(SELECT 1 FROM auth.phones WHERE prefix = $1 AND number = $2)`,
        PhonePrefix,
        PhoneNumber,
    ).Scan(&phoneExists)
    if err != nil {
        repo.log.Printf("Failed to check phone existence: %v", err)
        return nil, &entity.Error{
            Type:    entity.ErrAccountCreation,
            Message: entity.MsgAccountCreation,
            Detail:  fmt.Sprintf("phone validation failed: %v", err),
        }
    }

    if phoneExists {
        return nil, &entity.Error{
            Type:    entity.ErrPhoneAlreadyExists,
            Message: entity.MsgPhoneExists,
        }
    }

    // Password hashing
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
    if err != nil {
        repo.log.Printf("Failed to hash password: %v", err)
        return nil, &entity.Error{
            Type:    ErrPasswordHashingFailed,
            Message: MsgPasswordHashingFailed,
            Detail:  fmt.Sprintf("password hashing failed: %v", err),
        }
    }

    // User creation
    var userID uuid.UUID
    err = tx.QueryRow(
        `INSERT INTO auth.users (
            id,
            sir_name, 
            first_name, 
            last_name, 
            user_type
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING id`,
        uuid.New(),
        Title,
        FirstName,
        LastName,
        UserType,
    ).Scan(&userID)
    if err != nil {
        repo.log.Printf("Failed to create user: %v", err)
        return nil, &entity.Error{
            Type:    ErrUserCreationFailed,
            Message: MsgUserCreationFailed,
            Detail:  fmt.Sprintf("user insertion failed: %v", err),
        }
    }

    // Phone creation
    var phoneID uuid.UUID
    err = tx.QueryRow(
        `INSERT INTO auth.phones (
            created_at,
            id,
            prefix,
            number
        ) VALUES ($1, $2, $3, $4)
        RETURNING id`,
        time.Now(),
        uuid.New(),
        PhonePrefix,
        PhoneNumber,
    ).Scan(&phoneID)
    if err != nil {
        repo.log.Printf("Failed to create phone: %v", err)
        return nil, &entity.Error{
            Type:    ErrPhoneCreationFailed,
            Message: MsgPhoneCreationFailed,
            Detail:  fmt.Sprintf("phone insertion failed: %v", err),
        }
    }

    // Phone identity link
    _, err = tx.Exec(
        `INSERT INTO auth.phone_identities (
            created_at,
            id,
            user_id,
            phone_id
        ) VALUES ($1, $2, $3, $4)`,
        time.Now(),
        uuid.New(),
        userID,
        phoneID,
    )
    if err != nil {
        repo.log.Printf("Failed to create phone identity: %v", err)
        return nil, &entity.Error{
            Type:    ErrPhoneIdentityFailed,
            Message: MsgPhoneIdentityFailed,
            Detail:  fmt.Sprintf("phone identity creation failed: %v", err),
        }
    }

    // Password identity
    _, err = tx.Exec(
        `INSERT INTO auth.password_identities (
            created_at,
            updated_at,
            id,
            user_id,
            password,
            hint
        ) VALUES ($1, $2, $3, $4, $5, $6)`,
        time.Now(),
        time.Now(),
        uuid.New(),
        userID,
        string(hashedPassword),
        PasswordHint,
    )
    if err != nil {
        repo.log.Printf("Failed to create password identity: %v", err)
        return nil, &entity.Error{
            Type:    ErrPasswordIdentityFailed,
            Message: MsgPasswordIdentityFailed,
            Detail:  fmt.Sprintf("password identity creation failed: %v", err),
        }
    }

    // Create wallet based on user type
    walletAmount := 0.0
    var merchantID *uuid.UUID

    if UserType == "merchant" {
        // Create merchant
        newMerchantID := uuid.New()
        merchantID = &newMerchantID

        _, err = tx.Exec(
            `INSERT INTO merchants.merchants (
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
                $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW(), $13
            )`,
            newMerchantID,
            userID,
            fmt.Sprintf("%s %s %s Business PLC", Title, FirstName, LastName),
            fmt.Sprintf("%s %s Enterprises", FirstName, LastName),
            fmt.Sprintf("ET%s%d", time.Now().Format("20060102"), rand.Intn(10000)),
            fmt.Sprintf("TIN%s%d", time.Now().Format("20060102"), rand.Intn(10000)),
            "Retail Trade",
            "Sole Proprietorship",
            false,
            "",
            fmt.Sprintf("www.%s%s.com", strings.ToLower(FirstName), strings.ToLower(LastName)),
            time.Now().AddDate(-3, 0, 0),
            "active",
        )
        if err != nil {
            repo.log.Printf("Failed to create merchant: %v", err)
            return nil, &entity.Error{
                Type:    "ErrMerchantCreationFailed",
                Message: "Merchant creation failed",
                Detail:  fmt.Sprintf("merchant insertion failed: %v", err),
            }
        }

        // Create merchant address
        _, err = tx.Exec(
            `INSERT INTO merchants.addresses (
                merchant_id,
                personal_name,
                phone_number,
                region,
                city,
                sub_city,
                woreda,
                postal_code,
                secondary_phone_number,
                email
            ) VALUES (
                $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
            )`,
            newMerchantID,
            fmt.Sprintf("%s %s %s", Title, FirstName, LastName),
            PhonePrefix+PhoneNumber,
            "Addis Ababa",
            "Addis Ababa",
            "Bole",
            "08",
            "1000",
            "",
            fmt.Sprintf("%s.%s@example.com", strings.ToLower(FirstName), strings.ToLower(LastName)),
        )
        if err != nil {
            repo.log.Printf("Failed to create merchant address: %v", err)
            return nil, &entity.Error{
                Type:    "ErrMerchantAddressCreationFailed",
                Message: "Merchant address creation failed",
                Detail:  fmt.Sprintf("merchant address insertion failed: %v", err),
            }
        }

        // For merchants, double the initial amount
        walletAmount *= 2
    }

    // Create wallet for both merchant and non-merchant users
    _, err = tx.Exec(
        `INSERT INTO merchant.wallet (
            id, 
            user_id,
            merchant_id, 
            amount, 
            locked_amount, 
            currency, 
            created_at, 
            updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, NOW(), NOW()
        )`,
        uuid.New(),
        userID,
        merchantID, // This will be nil for non-merchant users
        walletAmount,
        walletAmount,
        "ETB",
    )
    if err != nil {
        repo.log.Printf("Failed to create wallet: %v", err)
        return nil, &entity.Error{
            Type:    "ErrWalletCreationFailed",
            Message: "Wallet creation failed",
            Detail:  fmt.Sprintf("wallet insertion failed: %v", err),
        }
    }

    if err = tx.Commit(); err != nil {
        repo.log.Printf("Failed to commit transaction: %v", err)
        return nil, &entity.Error{
            Type:    ErrCommitFailed,
            Message: MsgCommitFailed,
            Detail:  fmt.Sprintf("transaction commit failed: %v", err),
        }
    }

    return &entity.User{
        Id:        userID,
        SirName:   Title,
        FirstName: FirstName,
        LastName:  LastName,
        UserType:  UserType,
        CreatedAt: time.Now(),
        Phone: entity.Phone{
            Id:     phoneID,
            Prefix: PhonePrefix,
            Number: PhoneNumber,
        },
    }, nil
}

// Pre Session
func (repo PsqlRepo) StorePreSession(preSession entity.PreSession) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.pre_sessions (id, token, created_at)
	VALUES ($1::UUID, $2, $3);
	`, repo.schema), preSession.Id, sql.NullString{Valid: preSession.Token != "", String: preSession.Token}, preSession.CreatedAt)

	return err
}

func (repo PsqlRepo) UpdatePreSession(id string, token string) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	UPDATE %s.pre_sessions
	SET hash = $2
	WHERE id = $1
	RETURNING id, hash;
	`, repo.schema), id, token)

	return err
}

// Device

func (repo PsqlRepo) StoreDevice(device entity.Device) error {
	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.devices (id, ip, name, agent, created_at)
	VALUES ($1::UUID,$2,$3, $4, $5);
	`, repo.schema), device.Id, device.IP.String(), device.Name, device.Agent, device.CreatedAt)

	return err
}

func (repo PsqlRepo) StoreDeviceAuth(deviceAuth entity.DeviceAuth) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.device_auths (id, token, device, status, created_at)
	VALUES($1::UUID,$2, $3, $4, $5);
	`, repo.schema), deviceAuth.Id, deviceAuth.Token, deviceAuth.Device.Id, deviceAuth.Status, deviceAuth.CreatedAt)

	return err
}

func (repo PsqlRepo) UpdateDeviceAuthStatus(deviceAuthId uuid.UUID, status bool) error {

	_, err := repo.db.Query(fmt.Sprintf(`
	UPDATE %s.device_auths
	SET status = $2
	WHERE id = $1::UUID;
	`, repo.schema), deviceAuthId, status)

	return err
}

func (repo PsqlRepo) FindDeviceAuth(token string) (entity.DeviceAuth, error) {
	var deviceAuth entity.DeviceAuth

	var ip sql.NullString

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT device_auths.id, device_auths.token, devices.id, devices.ip, devices.name, devices.agent ,device_auths.status
	FROM %s.device_auths
	INNER JOIN %s.devices ON %s.devices.id = device_auths.device
	WHERE device_auths.token = $1
	`, repo.schema, repo.schema, repo.schema), token).Scan(
		&deviceAuth.Id, &deviceAuth.Token,
		&deviceAuth.Device.Id, &ip, &deviceAuth.Device.Name, &deviceAuth.Device.Agent,
		&deviceAuth.Status,
	)

	return deviceAuth, err
}

// Phone

func (repo PsqlRepo) StorePhone(phone entity.Phone) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.phones (id, prefix, number, created_at)
	VALUES ($1::UUID,$2,$3,$4);
	`, repo.schema), phone.Id, phone.Prefix, phone.Number, phone.CreatedAt)

	return err
}

func (repo PsqlRepo) FindPhone(prefix string, number string) (*entity.Phone, error) {

	var phone entity.Phone
	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT id, prefix, number 
	FROM %s.phones
	WHERE prefix = $1 AND number = $2;
	`, repo.schema), prefix, number).Scan(
		&phone.Id, &phone.Prefix, &phone.Number,
	)

	repo.log.Println("Find phone")

	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			{
				return nil, nil
			}
		}
		return nil, err
	}

	return &phone, err
}

func (repo PsqlRepo) LoginFindPhone(prefix string, number string) (*entity.Phone, error) {
	var phone entity.Phone

	err := repo.db.QueryRow(fmt.Sprintf(`
		SELECT id, prefix, number 
		FROM %s.phones
		WHERE prefix = $1 AND number = $2;
	`, repo.schema), prefix, number).Scan(&phone.Id, &phone.Prefix, &phone.Number)

	repo.log.Println("executed FindPhone query")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &phone, nil
}

func (repo PsqlRepo) StorePhoneAuth(phoneAuth entity.PhoneAuth) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.phone_auths (id, token, phone_id, code, method, status, created_at)
	VALUES ($1::UUID, $2, $3::UUID, $4, $5, $6, $7);
	`, repo.schema), phoneAuth.Id, phoneAuth.Token, phoneAuth.Phone.Id, phoneAuth.Code, phoneAuth.Method, phoneAuth.Status, time.Now())

	return err
}

func (repo PsqlRepo) FindPhoneAuth(token string) (entity.PhoneAuth, error) {

	var phoneAuth entity.PhoneAuth

	fmt.Println("||||||||||||||||||| ", token)

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT phone_auths.id, phone_auths.token, 
		phones.id, phones.prefix, phones.number,
		phone_auths.code, phone_auths.method, phone_auths.status
	FROM %s.phone_auths
	INNER JOIN %s.phones ON %s.phones.id = phone_auths.phone_id
	WHERE token = $1
	`, repo.schema, repo.schema, repo.schema), token).Scan(
		&phoneAuth.Id, &phoneAuth.Token,
		&phoneAuth.Phone.Id, &phoneAuth.Phone.Prefix, &phoneAuth.Phone.Number,
		&phoneAuth.Code, &phoneAuth.Method, &phoneAuth.Status,
	)

	return phoneAuth, err
}

func (repo PsqlRepo) FindPhoneAuthWithoutPhone(token string) (entity.PhoneAuth, error) {

	var phoneAuth entity.PhoneAuth

	fmt.Println("||||||||||||||||||| ", token)

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT phone_auths.id, phone_auths.token, 
		phone_auths.code, phone_auths.method, phone_auths.status
	FROM %s.phone_auths
	WHERE token = $1
	order by created_at DESC
	limit 1
	`, repo.schema), token).Scan(
		&phoneAuth.Id, &phoneAuth.Token,
		&phoneAuth.Code, &phoneAuth.Method, &phoneAuth.Status,
	)

	return phoneAuth, err
}

func (repo PsqlRepo) UpdatePhoneAuthStatus(phoneAuthId uuid.UUID, status bool) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	UPDATE %s.phone_auths
	SET status = $2
	WHERE id = $1::UUID
	`, repo.schema), phoneAuthId, status)
	return err
}

func (repo PsqlRepo) StoreSession(session entity.Session) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.sessions (id, token, user_id, device_id, created_at)
	VALUES ($1::UUID, $2, $3::UUID, $4::UUID, $5);
	`, repo.schema), session.Id, session.Token, session.User.Id, session.Device.Id, session.CreatedAt)

	return err
}

func (repo PsqlRepo) FindSessionById(id uuid.UUID) (*entity.Session, error) {
	var session entity.Session

	var sirName sql.NullString
	var lastName sql.NullString
	var userType sql.NullString

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT sessions.id, sessions.token, users.id, users.sir_name, users.first_name, users.last_name, users.user_type
	FROM %s.sessions
	INNER JOIN %s.users ON %s.users.id = sessions.user_id
	WHERE sessions.id = $1::UUID
	`, repo.schema, repo.schema, repo.schema), id).Scan(
		&session.Id, &session.Token,
		&session.User.Id, &sirName, &session.User.FirstName, &lastName, &userType,
	)

	if sirName.Valid {
		session.User.SirName = sirName.String
	}

	if lastName.Valid {
		session.User.LastName = lastName.String
	}

	if userType.Valid {
		session.User.UserType = userType.String
	} else {
		session.User.UserType = "UNKNOWN"
	}

	return &session, err
}

func (repo PsqlRepo) StorePasswordAuth(passwordAuth entity.PasswordAuth) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.password_auths (id, token, password_id, status, created_at)
	VALUES ($1::UUID, $2, $3::UUID, $4, $5);
	`, repo.schema), passwordAuth.Id, passwordAuth.Token, passwordAuth.Password.Id, passwordAuth.Status, passwordAuth.CreatedAt)

	return err
}

func (repo PsqlRepo) FindPasswordAuth(token string) (*entity.PasswordAuth, error) {

	var passwordAuth entity.PasswordAuth
	var hint sql.NullString

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT password_auths.id, password_auths.token, 
		password_identities.id, password_identities.password, password_identities.hint,
		password_auths.status, password_auths.created_at, password_auths.updated_at
	FROM %s.password_auths
	INNER JOIN %s.password_identities ON %s.password_identities.id = password_auths.password_id
	WHERE token = $1
	`, repo.schema, repo.schema, repo.schema), token).Scan(
		&passwordAuth.Id, &passwordAuth.Token,
		&passwordAuth.Password.Id, &passwordAuth.Password.Password, &hint,
		&passwordAuth.Status, &passwordAuth.CreatedAt, &passwordAuth.UpdatedAt,
	)

	if hint.Valid {
		passwordAuth.Password.Hint = hint.String
	}

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
	}

	return &passwordAuth, err
}

func (repo PsqlRepo) CheckPermission(userID uuid.UUID, requiredPermission entity.Permission) (bool, error) {
	// First check direct user permissions
	query := `
        WITH operation_id AS (
            SELECT id, name
            FROM auth.operations
            WHERE name = $2
        ),
        resource_id AS (
            SELECT id, name
            FROM auth.resources
            WHERE name = $3
        )
        SELECT p.resource, p.resource_id, o.name AS operation, p.effect,
               r.name as resource_name, oi.name as required_operation
        FROM auth.permissions p
        JOIN auth.user_permissions up ON up.permission_id = p.id
        JOIN resource_id r ON r.id = p.resource
        LEFT JOIN auth.operations o ON o.id = ANY(p.operations)
        CROSS JOIN operation_id oi
        WHERE up.user_id = $1
        AND EXISTS (
            SELECT 1
            FROM operation_id oi
            WHERE oi.id = ANY(p.operations)
        )
        AND p.effect = $4;
    `

	rowsUserPermissions, err := repo.db.Query(query, userID, requiredPermission.Operation, requiredPermission.Resource, requiredPermission.Effect)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch user permissions: %v", err)
		return false, fmt.Errorf("failed to fetch user permissions: %v", err)
	}
	defer rowsUserPermissions.Close()

	var permissionFound bool
	for rowsUserPermissions.Next() {
		var permission entity.Permission
		var resourceName, requiredOp string
		if err := rowsUserPermissions.Scan(&permission.Resource, &permission.ResourceIdentifier,
			&permission.Operation, &permission.Effect, &resourceName, &requiredOp); err != nil {
			log.Printf("[ERROR] Failed to scan permission: %v", err)
			return false, fmt.Errorf("failed to scan permission: %v", err)
		}


		if permission.Effect == requiredPermission.Effect {
			log.Printf("[DEBUG] User %s has direct permission to perform the operation", userID)
			permissionFound = true
			break
		}
	}

	if permissionFound {
		return true, nil
	}

	// Check group permissions
	queryGroups := `
        SELECT g.id, g.title
        FROM auth.user_groups ug
        JOIN auth.groups g ON ug.group_id = g.id
        WHERE ug.user_id = $1
    `
	rowsGroups, err := repo.db.Query(queryGroups, userID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch user groups: %v", err)
		return false, fmt.Errorf("failed to fetch user groups: %v", err)
	}
	defer rowsGroups.Close()

	var groups []uuid.UUID
	for rowsGroups.Next() {
		var groupID uuid.UUID
		var groupTitle string
		if err := rowsGroups.Scan(&groupID, &groupTitle); err != nil {
			log.Printf("[ERROR] Failed to scan group: %v", err)
			return false, fmt.Errorf("failed to scan group: %v", err)
		}
		groups = append(groups, groupID)
	}

	if len(groups) == 0 {
		return false, fmt.Errorf("user does not belong to any group")
	}

	queryGroupPermissions := `
        WITH operation_id AS (
            SELECT id, name
            FROM auth.operations
            WHERE name = $3
        ),
        resource_id AS (
            SELECT id, name
            FROM auth.resources
            WHERE name = $2
        )
        SELECT p.resource, p.resource_id, o.name AS operation, p.effect,
               r.name as resource_name, oi.name as required_operation,
               g.title as group_name
        FROM auth.group_permissions gp
        JOIN auth.permissions p ON gp.permission_id = p.id
        JOIN auth.user_groups ug ON ug.group_id = gp.group_id
        JOIN auth.groups g ON g.id = gp.group_id
        JOIN resource_id r ON r.id = p.resource
        LEFT JOIN auth.operations o ON o.id = ANY(p.operations)
        CROSS JOIN operation_id oi
        WHERE ug.user_id = $1
        AND EXISTS (
            SELECT 1
            FROM operation_id oi
            WHERE oi.id = ANY(p.operations)
        )
        AND p.effect = $4;
    `

	rowsGroupPermissions, err := repo.db.Query(queryGroupPermissions, userID, requiredPermission.Resource, requiredPermission.Operation, requiredPermission.Effect)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch group permissions: %v", err)
		return false, fmt.Errorf("failed to fetch group permissions: %v", err)
	}
	defer rowsGroupPermissions.Close()

	for rowsGroupPermissions.Next() {
		var permission entity.Permission
		var resourceName, requiredOp, groupName string
		if err := rowsGroupPermissions.Scan(&permission.Resource, &permission.ResourceIdentifier,
			&permission.Operation, &permission.Effect, &resourceName, &requiredOp, &groupName); err != nil {
			log.Printf("[ERROR] Failed to scan group permission: %v", err)
			return false, fmt.Errorf("failed to scan group permission: %v", err)
		}

		if permission.Effect == requiredPermission.Effect {
			return true, nil
		}
	}

	return false, fmt.Errorf("user does not have the required permission")
}

func (repo PsqlRepo) FindUserPermissions(userID uuid.UUID, requiredPermission entity.Permission) ([]entity.Permission, error) {
	var permissions []entity.Permission
	query := `
		SELECT resource, resource_identifier, operation, effect
		FROM auth.user_permissions
		WHERE user_id = $1
		AND resource = $2
		AND operation = $3
		AND resource_identifier = $4
	`
	rows, err := repo.db.Query(query, userID, requiredPermission.Resource, requiredPermission.Operation, requiredPermission.ResourceIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user permissions: %v", err)
	}
	defer rows.Close()

	// Collect all the permissions for the user
	for rows.Next() {
		var permission entity.Permission
		if err := rows.Scan(&permission.Resource, &permission.ResourceIdentifier, &permission.Operation, &permission.Effect); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %v", err)
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}
