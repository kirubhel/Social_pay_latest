package usecase

import (
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/key/core/entity"

	//import procedeure "github.com/socialpay/socialpay/src/pkg/key/adapter/controller/procedure"
	"github.com/socialpay/socialpay/src/pkg/key/adapter/controller/procedure"

	"github.com/google/uuid"
)

func (u *Usecase) CreateAPIKey(merchantID, service string, expiryDate time.Time, store string, commissionFrom bool) (*entity.APIKey, error) {
	// Generate a new key pair
	keyPair, err := procedure.GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	apiKey := entity.APIKey{
		MerchantID: merchantID,
		Service:    service,
		ExpiryDate: expiryDate,
		Store:      store,
		IsActive:   true,
		APIKey:     uuid.New().String(),
		PrivateKey: keyPair.PrivateKey,
		PublicKey:  keyPair.PublicKey,
	}

	// Save the key pair to the repository
	err = u.repo.Save(&apiKey)
	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}
func (u *Usecase) ToggleKey(keyID string, active bool) error {
	// TODO: Add real logic here
	return nil
}

func (u *Usecase) ValidateAPIToken(token string) (*entity.APIKey, error) {
	//todo
	apiKey, err := u.repo.FindByToken(token)
	if err != nil {
		fmt.Println("Error finding API key:", err)
		return nil, err
	}
	return apiKey, nil

}
