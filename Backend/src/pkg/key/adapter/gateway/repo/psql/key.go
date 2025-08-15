// github.com/socialpay/socialpay/src/pkg/key/repo/psql_repo.go
package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/key/core/entity"

	"github.com/google/uuid"
)

func (r *PsqlRepo) Save(apiKey *entity.APIKey) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.api_keys 
        (id, merchant_id, private_key, public_key, api_key, service, expiry_date, store, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, r.schema)
	apiKey.APIKey = uuid.New().String()

	_, err := r.db.Exec(query,
		uuid.New(),
		apiKey.MerchantID,
		apiKey.PrivateKey,
		apiKey.PublicKey,
		apiKey.APIKey,
		apiKey.Service,
		apiKey.ExpiryDate,
		apiKey.Store,
		apiKey.IsActive,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *PsqlRepo) FindByToken(token string) (*entity.APIKey, error) {
	query := fmt.Sprintf(`
        SELECT id, merchant_id, private_key, public_key, api_key, service, 
               expiry_date, store, is_active, created_at, updated_at
        FROM %s.api_keys 
        WHERE api_key = $1`, r.schema)

	var key entity.APIKey
	err := r.db.QueryRow(query, token).Scan(
		&key.ID,
		&key.MerchantID,
		&key.PrivateKey,
		&key.PublicKey,
		&key.APIKey,
		&key.Service,
		&key.ExpiryDate,
		&key.Store,
		&key.IsActive,
		&key.CreatedAt,
		&key.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No API key found with the given token")
			return nil, nil
		}
		fmt.Println("Error scanning API key:", err)
		return nil, fmt.Errorf("failed to find API key: %w", err)
	}
	return &key, nil
}

func (r *PsqlRepo) UpdateStatus(token string, enabled bool) error {
	query := fmt.Sprintf(`
        UPDATE %s.api_keys 
        SET is_enabled = $1, updated_at = NOW() 
        WHERE token = $2`, r.schema)

	_, err := r.db.Exec(query, enabled, token)
	if err != nil {
		return fmt.Errorf("failed to update API key status: %w", err)
	}
	return nil
}
