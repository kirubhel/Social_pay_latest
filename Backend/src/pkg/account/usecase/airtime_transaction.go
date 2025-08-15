package usecase

import (
	"fmt"

	"github.com/socialpay/socialpay/src/pkg/account/core/entity"

	"github.com/google/uuid"
)

// custom error struct to handle structured error responses
type JSONError struct {
	ErrorType string `json:"errorType"`
	Message   string `json:"message"`
}

// is responsible for storing an airtime transaction
func (uc Usecase) StoreAirtimeTransaction(merchantID uuid.UUID, amount int, msisdn string, transactionRef string) (*entity.AirtimeTransaction, error) {
	const ErrFailedToStoreTransaction = "FAILED_TO_STORE_TRANSACTION"
	const ErrInvalidInput = "INVALID_INPUT"
	if err := validateAirtimeTransactionInput(amount, msisdn, transactionRef); err != nil {
		return nil, err
	}

	uc.log.Println("CREATING AIRTIME TRANSACTION")
	transaction, err := uc.repo.StoreAirtimeTransaction(
		merchantID,
		amount,
		msisdn,
		transactionRef,
	)
	if err != nil {
		uc.log.Println("ERROR STORING AIRTIME TRANSACTION")
		return nil, fmt.Errorf("%s: %v", ErrFailedToStoreTransaction, err)
	}

	uc.log.Println("AIRTIME TRANSACTION CREATED SUCCESSFULLY")
	return transaction, nil
}

func validateAirtimeTransactionInput(amount int, msisdn, transactionRef string) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount must be greater than zero")
	}
	if msisdn == "" {
		return fmt.Errorf("invalid msisdn cannot be empty")
	}
	if transactionRef == "" {
		return fmt.Errorf("invalid transaction reference cannot be empty")
	}
	return nil
}

func (uc Usecase) UpdateAirtimeSuccessTransaction(transactionRef string, telebirrRef string, responseData map[string]interface{}) error {
	const ErrFailedToUpdateTransaction = "FAILED_TO_UPDATE_TRANSACTION"
	const ErrInvalidInput = "INVALID_INPUT"

	if transactionRef == "" || telebirrRef == "" {
		return fmt.Errorf("%s: transactionRef and telebirrRef cannot be empty", ErrInvalidInput)
	}

	uc.log.Println("UPDATING AIRTIME TRANSACTION")
	err := uc.repo.UpdateAirtimeSuccessTransaction(
		transactionRef,
		telebirrRef,
		responseData,
	)
	if err != nil {
		uc.log.Println("ERROR UPDATING AIRTIME TRANSACTION")
		return fmt.Errorf("%s: %v", ErrFailedToUpdateTransaction, err)
	}

	uc.log.Println("AIRTIME TRANSACTION UPDATED SUCCESSFULLY")
	return nil
}

func (uc Usecase) GetAirtimeTransactions() ([]entity.AirtimeTransaction, error) {
	const ErrFailedToGetTransaction = "FAILED_TO_UPDATE_TRANSACTION"

	uc.log.Println("UPDATING AIRTIME TRANSACTION")
	transactions, err := uc.repo.GetAirtimeTransactions()
	if err != nil {
		uc.log.Println("ERROR UPDATING AIRTIME TRANSACTION")
		return nil, fmt.Errorf("%s: %v", ErrFailedToGetTransaction, err)
	}

	uc.log.Println("AIRTIME TRANSACTION UPDATED SUCCESSFULLY")
	return transactions, nil
}
