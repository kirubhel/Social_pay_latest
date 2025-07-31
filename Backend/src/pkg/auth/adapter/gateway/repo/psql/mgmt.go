package repo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"

	"github.com/google/uuid"
)

func (repo PsqlRepo) StoreUser(user entity.User) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.users (id, sir_name, first_name, last_name, created_at)
	VALUES ($1::UUID, $2, $3, $4, $5);
	`, repo.schema), user.Id,
		sql.NullString{Valid: user.SirName != "", String: user.SirName},
		sql.NullString{Valid: user.FirstName != "", String: user.FirstName},
		sql.NullString{Valid: user.LastName != "", String: user.LastName},
		user.CreatedAt,
	)

	return err
}

func (repo PsqlRepo) StorePhoneIdentity(phoneIdentity entity.PhoneIdentity) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.phone_identities (id, "user_id", phone_id, created_at)
	VALUES ($1::UUID, $2::UUID, $3::UUID, $4);
	`, repo.schema), phoneIdentity.Id, phoneIdentity.User.Id, phoneIdentity.Phone.Id, phoneIdentity.CreatedAt)

	log.Println(err)

	return err
}

func (repo PsqlRepo) FindUserUsingPhoneIdentity(phoneId uuid.UUID) (*entity.User, error) {
	var user entity.User

	var sirName sql.NullString
	var lastName sql.NullString
	var userType sql.NullString // Add variable for user_type

	var twoFactorVerifiedAt sql.NullTime

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT users.id, users.sir_name, users.first_name, users.last_name, users.user_type, users.two_factor_enabled, users.two_factor_verified_at
	FROM %s.phone_identities
	INNER JOIN %s.users ON %s.users.id = phone_identities.user_id
	WHERE phone_id = $1;
	`, repo.schema, repo.schema, repo.schema), phoneId).Scan(
		&user.Id, &sirName, &user.FirstName, &lastName, &userType, &user.TwoFactorEnabled, &twoFactorVerifiedAt,
	)

	if sirName.Valid {
		user.SirName = sirName.String
	}

	if lastName.Valid {
		user.LastName = lastName.String
	}

	if userType.Valid {
		user.UserType = userType.String // Set the user_type if valid
	} else {
		user.UserType = "UNKNOWN" // Default value if user_type is not found
	}

	if twoFactorVerifiedAt.Valid {
		user.TwoFactorVerifiedAt = &twoFactorVerifiedAt.Time
	}

	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			{
				return nil, nil
			}
		}
		return nil, err
	}

	return &user, nil
}

func (repo PsqlRepo) FindUserById(id uuid.UUID) (*entity.User, error) {
	var user entity.User

	var sirName sql.NullString
	var lastName sql.NullString
	var userType sql.NullString
	var twoFactorVerifiedAt sql.NullTime

	fmt.Println("################################################### , repo 1", id)

	sqlStmt := `SELECT id, sir_name, first_name, last_name, user_type, two_factor_enabled, two_factor_verified_at, created_at
	FROM auth.users
	WHERE id = $1;
	`

	err := repo.db.QueryRow(sqlStmt, id).Scan(&user.Id, &sirName, &user.FirstName, &lastName, &userType, &user.TwoFactorEnabled, &twoFactorVerifiedAt, &user.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return &entity.User{}, nil
	}
	fmt.Println("################################################### , repo 10")

	if sirName.Valid {
		user.SirName = sirName.String
	}

	if lastName.Valid {
		user.LastName = lastName.String
	}

	if userType.Valid {
		user.UserType = userType.String // Set the user_type if valid
	} else {
		user.UserType = "UNKNOWN" // Default value if user_type is not found
	}

	if twoFactorVerifiedAt.Valid {
		user.TwoFactorVerifiedAt = &twoFactorVerifiedAt.Time
	}

	fmt.Println("################################################### , repo 2")

	return &user, nil
}

func (repo PsqlRepo) StorePasswordIdentity(passwordIdentity entity.PasswordIdentity) error {
	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.password_identities (id, "user_id", password, hint, created_at)
	VALUES ($1::UUID, $2::UUID, $3, $4, $5);
	`, repo.schema), passwordIdentity.Id, passwordIdentity.User.Id, passwordIdentity.Password, sql.NullString{String: passwordIdentity.Hint, Valid: passwordIdentity.Hint != ""}, passwordIdentity.CreatedAt)

	return err
}
func (repo PsqlRepo) UpdatePasswordIdentity(password string, userId uuid.UUID) error {
	fmt.Println("|||||||||||| 0000 ", password)
	fmt.Println("|||||||||||| 0000 ", userId)

	_, err := repo.db.Exec(fmt.Sprintf(`
	UPDATE %s.password_identities
		SET password = $1, updated_at = NOW()
		WHERE user_id = $2;
	`, repo.schema), password, userId)

	return err
}

func (repo PsqlRepo) FindPasswordIdentityByUser(userId uuid.UUID) (*entity.PasswordIdentity, error) {
	var passwordIdentity entity.PasswordIdentity

	var hint sql.NullString

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT id, user_id, password, hint, created_at, updated_at
	FROM %s.password_identities
	WHERE user_id = $1::UUID;
	`, repo.schema), userId).Scan(
		&passwordIdentity.Id,
		&passwordIdentity.User.Id,
		&passwordIdentity.Password,
		&hint,
		&passwordIdentity.CreatedAt,
		&passwordIdentity.UpdatedAt,
	)

	if hint.Valid {
		passwordIdentity.Hint = hint.String
	}

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
	}

	return &passwordIdentity, err
}

// Two-Factor Authentication Methods

func (repo PsqlRepo) GetTwoFactorStatus(userId uuid.UUID) (*entity.TwoFactorStatus, error) {
	var status entity.TwoFactorStatus
	var verifiedAt sql.NullTime

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT two_factor_enabled, two_factor_verified_at
	FROM %s.users
	WHERE id = $1::UUID;
	`, repo.schema), userId).Scan(&status.Enabled, &verifiedAt)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return &entity.TwoFactorStatus{Enabled: false}, nil
		}
		return nil, err
	}

	if verifiedAt.Valid {
		status.VerifiedAt = &verifiedAt.Time
	}

	return &status, nil
}

func (repo PsqlRepo) EnableTwoFactor(userId uuid.UUID) error {
	_, err := repo.db.Exec(fmt.Sprintf(`
	UPDATE %s.users
	SET two_factor_enabled = TRUE, two_factor_verified_at = NOW()
	WHERE id = $1::UUID;
	`, repo.schema), userId)
	return err
}

func (repo PsqlRepo) DisableTwoFactor(userId uuid.UUID) error {
	_, err := repo.db.Exec(fmt.Sprintf(`
	UPDATE %s.users
	SET two_factor_enabled = FALSE, two_factor_verified_at = NULL
	WHERE id = $1::UUID;
	`, repo.schema), userId)
	return err
}

func (repo PsqlRepo) StoreTwoFactorCode(code entity.TwoFactorCode) error {
	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.two_factor_codes (id, user_id, code, expires_at, used, created_at)
	VALUES ($1::UUID, $2::UUID, $3, $4, $5, $6);
	`, repo.schema), code.Id, code.UserId, code.Code, code.ExpiresAt, code.Used, code.CreatedAt)
	return err
}

func (repo PsqlRepo) FindTwoFactorCode(userId uuid.UUID, code string) (*entity.TwoFactorCode, error) {
	var twoFactorCode entity.TwoFactorCode

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT id, user_id, code, expires_at, used, created_at, updated_at
	FROM %s.two_factor_codes
	WHERE user_id = $1::UUID AND code = $2 AND used = FALSE AND expires_at > NOW()
	ORDER BY created_at DESC
	LIMIT 1;
	`, repo.schema), userId, code).Scan(
		&twoFactorCode.Id,
		&twoFactorCode.UserId,
		&twoFactorCode.Code,
		&twoFactorCode.ExpiresAt,
		&twoFactorCode.Used,
		&twoFactorCode.CreatedAt,
		&twoFactorCode.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &twoFactorCode, nil
}

func (repo PsqlRepo) MarkTwoFactorCodeAsUsed(codeId uuid.UUID) error {
	_, err := repo.db.Exec(fmt.Sprintf(`
	UPDATE %s.two_factor_codes
	SET used = TRUE, updated_at = NOW()
	WHERE id = $1::UUID;
	`, repo.schema), codeId)
	return err
}

func (repo PsqlRepo) CleanupExpiredTwoFactorCodes() error {
	_, err := repo.db.Exec(fmt.Sprintf(`
	DELETE FROM %s.two_factor_codes
	WHERE expires_at < NOW() OR used = TRUE;
	`, repo.schema))
	return err
}

func (repo PsqlRepo) FindUserWithPhoneById(userId uuid.UUID) (*entity.User, error) {
	var user entity.User
	var sirName sql.NullString
	var lastName sql.NullString
	var userType sql.NullString
	var twoFactorVerifiedAt sql.NullTime
	var phoneId sql.NullString
	var phonePrefix sql.NullString
	var phoneNumber sql.NullString

	// First, let's check if the user exists and get their basic info
	userSqlStmt := fmt.Sprintf(`
	SELECT id, sir_name, first_name, last_name, user_type, 
		two_factor_enabled, two_factor_verified_at, created_at
	FROM %s.users
	WHERE id = $1::UUID;
	`, repo.schema)

	err := repo.db.QueryRow(userSqlStmt, userId).Scan(
		&user.Id, &sirName, &user.FirstName, &lastName, &userType,
		&user.TwoFactorEnabled, &twoFactorVerifiedAt, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Set user fields
	if sirName.Valid {
		user.SirName = sirName.String
	}

	if lastName.Valid {
		user.LastName = lastName.String
	}

	if userType.Valid {
		user.UserType = userType.String
	} else {
		user.UserType = "UNKNOWN"
	}

	if twoFactorVerifiedAt.Valid {
		user.TwoFactorVerifiedAt = &twoFactorVerifiedAt.Time
	}

	// Now get the phone information using LEFT JOIN to handle missing phone data gracefully
	phoneSqlStmt := fmt.Sprintf(`
	SELECT p.id, p.prefix, p.number
	FROM %s.phone_identities pi
	LEFT JOIN %s.phones p ON pi.phone_id = p.id
	WHERE pi.user_id = $1::UUID;
	`, repo.schema, repo.schema)

	err = repo.db.QueryRow(phoneSqlStmt, userId).Scan(
		&phoneId, &phonePrefix, &phoneNumber,
	)
	if err != nil {
		// Log the error for debugging
		repo.log.Printf("Failed to get phone for user %s: %v", userId, err)
		// Return user without phone data instead of failing completely
		return &user, nil
	}

	// Set phone information
	if phoneId.Valid {
		user.PhoneID, _ = uuid.Parse(phoneId.String)
	}

	if phonePrefix.Valid && phoneNumber.Valid {
		user.Phone = entity.Phone{
			Id:     user.PhoneID,
			Prefix: phonePrefix.String,
			Number: phoneNumber.String,
		}
	}

	return &user, nil
}
