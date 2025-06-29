package psql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/account/core/entity"

	"github.com/google/uuid"
)

func (repo PsqlRepo) StoreAirtimeTransaction(merchantID uuid.UUID, amount int, msisdn string, transactionRef string) (*entity.AirtimeTransaction, error) {
	query := `
		INSERT INTO accounts.airtime_transactions (
			merchant_id, 
			payment_method, 
			amount, 
			currency, 
			phone_number,
			status, 
			reference_code, 
			webhook_notified, 
			created_at, 
			updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, merchant_id, amount, currency, phone_number, status, reference_code, created_at, updated_at;
	`

	// Default values
	currency := "ETB"
	status := "Pending"
	paymentMethod := "Telebirr"
	webhookNotified := false
	createdAt := time.Now()
	updatedAt := createdAt

	var transaction entity.AirtimeTransaction
	err := repo.db.QueryRow(
		query,
		merchantID,
		paymentMethod,
		amount,
		currency,
		msisdn,
		status,
		transactionRef,
		webhookNotified,
		createdAt,
		updatedAt,
	).Scan(
		&transaction.ID,
		&transaction.MerchantID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.MSISDN,
		&transaction.Status,
		&transaction.TransactionRef,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to store airtime transaction: %w", err)
	}

	return &transaction, nil
}

func (repo PsqlRepo) UpdateAirtimeSuccessTransaction(transactionRef string, telebirrRef string, responseData map[string]interface{}) error {
	query := `
		UPDATE accounts.airtime_transactions
		SET 
			status = $1,
			yimulu_transaction_ref = $2,
			transaction_details = $3,
			updated_at = $4
		WHERE reference_code = $5
		RETURNING id, status, yimulu_transaction_ref, transaction_details, updated_at;
	`

	status := "Completed"
	updatedAt := time.Now()
	transactionDetails, err := json.Marshal(responseData)
	if err != nil {
		return fmt.Errorf("failed to marshal responseData: %w", err)
	}

	var transaction entity.AirtimeTransaction
	// Perform the query and scan
	err = repo.db.QueryRow(
		query,
		status,
		telebirrRef,
		transactionDetails,
		updatedAt,
		transactionRef,
	).Scan(
		&transaction.ID,
		&transaction.Status,
		&transaction.YimuluTransactionRef,
		&transaction.TransactionDetails,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction status to completed: %w", err)
	}

	if transaction.TransactionDetails != nil {
		var detailsMap map[string]interface{}
		if err := json.Unmarshal(transaction.TransactionDetails, &detailsMap); err != nil {
			return fmt.Errorf("failed to unmarshal transaction_details: %w", err)
		}
		if transaction.TransactionDetails, err = json.Marshal(detailsMap); err != nil {
			return fmt.Errorf("failed to marshal transaction_details: %w", err)
		}
	}

	return nil
}

func (repo PsqlRepo) GetAirtimeTransactions() ([]entity.AirtimeTransaction, error) {
	query := `
        SELECT 
            id,
            merchant_id,
            payment_method,
            amount,
            currency,
            phone_number,
            status,
            reference_code,
            yimulu_transaction_ref,
            transaction_details,
            webhook_notified,
            created_at,
            updated_at
        FROM accounts.airtime_transactions
        ORDER BY created_at DESC
    `

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query airtime transactions: %w", err)
	}
	defer rows.Close()

	var transactions []entity.AirtimeTransaction
	for rows.Next() {
		var transaction entity.AirtimeTransaction
		var yimuluRef sql.NullString
		var transactionDetails sql.NullString

		err := rows.Scan(
			&transaction.ID,
			&transaction.MerchantID,
			&transaction.PaymentMethod,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.MSISDN,
			&transaction.Status,
			&transaction.TransactionRef,
			&yimuluRef,
			&transactionDetails,
			&transaction.WebhookNotified,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan airtime transaction: %w", err)
		}

		// Convert NullString to regular string
		transaction.YimuluTransactionRef = yimuluRef.String
		transaction.TransactionDetails = []byte(transactionDetails.String)

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning airtime transactions: %w", err)
	}

	return transactions, nil
}
