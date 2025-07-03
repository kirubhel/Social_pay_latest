package usecase

import (
	"time"

	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc Usecase) GetUserByPhoneNumber(phoneId uuid.UUID) (entity.User, error) {
	var ErrUserNotFound string = "USER_NOT_FOUND"

	user, err := uc.repo.FindUserUsingPhoneIdentity(phoneId)
	if err != nil {
		return entity.User{}, Error{
			Type:    ErrUserNotFound,
			Message: err.Error(),
		}
	}
	return *user, nil
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (uc Usecase) CreatePasswordIdentity(userId uuid.UUID, password string, hint string) (*entity.PasswordIdentity, error) {
	var identity entity.PasswordIdentity

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, Error{
			Type:    "ERRCRATINGPASS",
			Message: err.Error(),
		}
	}

	identity = entity.PasswordIdentity{
		Id:        uuid.New(),
		User:      entity.User{Id: userId},
		Password:  hashedPassword,
		Hint:      hint,
		CreatedAt: time.Now(),
	}

	// Store password
	err = uc.repo.StorePasswordIdentity(identity)
	if err != nil {
		return nil, Error{
			Type:    "FAILED_TO_CREATE_2FA",
			Message: err.Error(),
		}
	}

	return &identity, nil
}

func (uc Usecase) GetUserWithPhoneById(id uuid.UUID) (*entity.User, error) {
	user, err := uc.repo.FindUserWithPhoneById(id)
	if err != nil {
		return nil, Error{
			Type:    "USER_NOT_FOUND",
			Message: err.Error(),
		}
	}
	return user, nil
}
