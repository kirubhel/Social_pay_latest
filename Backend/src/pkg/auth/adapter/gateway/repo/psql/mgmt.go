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

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT users.id, users.sir_name, users.first_name, users.last_name, users.user_type
	FROM %s.phone_identities
	INNER JOIN %s.users ON %s.users.id = phone_identities.user_id
	WHERE phone_id = $1;
	`, repo.schema, repo.schema, repo.schema), phoneId).Scan(
		&user.Id, &sirName, &user.FirstName, &lastName, &userType,
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
	var userType sql.NullString // Add variable for user_type

	fmt.Println("################################################### , repo 1", id)

	sqlStmt := `SELECT id, sir_name, first_name, last_name, user_type, created_at
	FROM auth.users
	WHERE id = $1;
	`

	err := repo.db.QueryRow(sqlStmt, id).Scan(&user.Id, &sirName, &user.FirstName, &lastName, &userType, &user.CreatedAt)
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
		SET finger_password = CASE WHEN $1 <> '' THEN $1 ELSE finger_password END
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
