package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (repo PsqlRepo) ListAllGateways() ([]entity.PaymentGateway, error) {
	startTime := time.Now()
	repo.log.Println("Repository: Starting to list all payment gateways")

	query := `
		SELECT 
			id,
			name,
			description,
			type,
			is_active,
			config,
			created_at,
			updated_at
		FROM merchants.payment_gateways
		ORDER BY name ASC
	`

	repo.log.Printf("Repository: Executing query: %s", query)
	rows, err := repo.db.Query(query)
	if err != nil {
		repo.log.Printf("Repository: Database query failed: %v", err)
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.log.Printf("Repository: Error closing rows: %v", err)
		}
	}()

	var gateways []entity.PaymentGateway
	repo.log.Println("Repository: Scanning rows...")

	for rows.Next() {
		var gateway entity.PaymentGateway
		var configJSON []byte
		var updatedAt sql.NullTime

		err := rows.Scan(
			&gateway.ID,
			&gateway.Name,
			&gateway.Description,
			&gateway.Type,
			&gateway.IsActive,
			&configJSON,
			&gateway.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			repo.log.Printf("Repository: Row scan failed: %v", err)
			return nil, fmt.Errorf("row scan failed: %w", err)
		}

		repo.log.Printf("Repository: Unmarshaling config for gateway ID %s", gateway.ID)
		var config entity.GatewayConfig
		if err := json.Unmarshal(configJSON, &config); err != nil {
			repo.log.Printf("Repository: Config unmarshal failed for gateway ID %s: %v", gateway.ID, err)
			return nil, fmt.Errorf("config unmarshal failed: %w", err)
		}
		gateway.Config = config

		if updatedAt.Valid {
			gateway.UpdatedAt = updatedAt.Time
		} else {
			gateway.UpdatedAt = gateway.CreatedAt
		}

		gateways = append(gateways, gateway)
	}

	if err := rows.Err(); err != nil {
		repo.log.Printf("Repository: Rows iteration error: %v", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	duration := time.Since(startTime)
	repo.log.Printf("Repository: Successfully retrieved %d gateways in %v", len(gateways), duration)
	return gateways, nil
}
