package usecase

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"

	"github.com/google/uuid"
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

func (uc Usecase) CreatePasswordIdentity(userId uuid.UUID, password string, hint string) (*entity.PasswordIdentity, error) {
	var identity entity.PasswordIdentity

	hasher := sha256.New()
	_, err := hasher.Write([]byte(password))
	if err != nil {
		return nil, Error{
			Type:    "ERRCRATINGPASS",
			Message: err.Error(),
		}
	}

	identity = entity.PasswordIdentity{
		Id:        uuid.New(),
		User:      entity.User{Id: userId},
		Password:  base64.URLEncoding.EncodeToString(hasher.Sum(nil)),
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
