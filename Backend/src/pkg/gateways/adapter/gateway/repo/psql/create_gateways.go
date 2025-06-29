package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
	_ "github.com/lib/pq"
)

func (repo PsqlRepo) CreateGateway(gateway entity.ListPaymentGateway) error {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if gateway name already exists
	var nameExists bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM merchants.payment_gateways WHERE name = $1)
	`, gateway.Name).Scan(&nameExists)
	if err != nil {
		return fmt.Errorf("failed to check gateway name existence: %v", err)
	}
	if nameExists {
		return entity.ErrGatewayAlreadyExists
	}

	// Marshal config to JSON
	configJSON, err := json.Marshal(gateway.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal gateway config: %v", err)
	}

	// Insert the new gateway
	query := `
		INSERT INTO merchants.payment_gateways (
			name, 
			description, 
			type,
			is_active,
			config,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, created_at, updated_at
	`

	var (
		id        uuid.UUID
		createdAt time.Time
		updatedAt time.Time
	)

	err = tx.QueryRow(
		query,
		gateway.Name,
		gateway.Description,
		gateway.Type,
		gateway.IsActive,
		configJSON, // Marshaled JSON config
		time.Now(), // created_at
		time.Now(), // updated_at
	).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Failed to insert payment gateway: %v", err)
		return fmt.Errorf("failed to insert payment gateway: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (repo PsqlRepo) GetGatewayByID(
	gatewayID uuid.UUID,
	includeInactive bool,
	limit int,
	offset int,
) ([]entity.GatewayMerchant, int, error) {
	// Base query
	query := `
		SELECT 
			gm.merchant_id,
			m.trading_name as name,
			m.business_registration_number as business_id,
			gm.is_active,
			gm.linked_at,
			gm.disabled_at,
			gm.disabled_reason,
			COUNT(*) OVER() as total_count
		FROM merchants.merchant_gateways gm
		JOIN merchants.merchants m ON gm.merchant_id = m.id
		WHERE gm.gateway_id = $1
	`

	// Add inactive filter if needed
	if !includeInactive {
		query += " AND gm.is_active = true"
	}

	// Add pagination
	query += " LIMIT $2 OFFSET $3"

	rows, err := repo.db.Query(query, gatewayID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query gateway merchants: %v", err)
	}
	defer rows.Close()

	var merchants []entity.GatewayMerchant
	var totalCount int

	for rows.Next() {
		var gm entity.GatewayMerchant
		var disabledReason sql.NullString
		err := rows.Scan(
			&gm.MerchantID,
			&gm.Name,
			&gm.BusinessID,
			&gm.IsActive,
			&gm.LinkedAt,
			&gm.DisabledAt,
			&disabledReason,
			&totalCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan gateway merchant: %v", err)
		}

		if disabledReason.Valid {
			gm.DisabledReason = &disabledReason.String
		}

		merchants = append(merchants, gm)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error after scanning rows: %v", err)
	}

	return merchants, totalCount, nil
}
