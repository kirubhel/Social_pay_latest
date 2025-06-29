package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (repo PsqlRepo) UpdateGateway(
	gatewayID uuid.UUID,
	name *string,
	description *string,
	isActive *bool,
	config *entity.GatewayConfig,
) (*entity.PaymentGateway, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	// First check if the gateway exists
	var exists bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM payment_gateways WHERE id = $1)
	`, gatewayID).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to check gateway existence: %v", err)
	}
	if !exists {
		tx.Rollback()
		return nil, entity.ErrGatewayNotFound
	}

	// Build the update query dynamically based on provided fields
	query := `
		UPDATE payment_gateways 
		SET 
			updated_at = $1,
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			is_active = COALESCE($4, is_active),
			config = COALESCE($5, config)
		WHERE id = $6
		RETURNING 
			id, name, description, type, is_active, config, 
			created_at, updated_at, linked_at, 
			disabled_reason, disabled_at
	`

	// Prepare config JSON if provided
	var configJSON []byte
	if config != nil {
		configJSON, err = json.Marshal(config)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to marshal config: %v", err)
		}
	}

	// Execute the update
	var gateway entity.PaymentGateway
	var rawConfig []byte
	var linkedAt sql.NullTime
	var disabledAt sql.NullTime
	var disabledReason sql.NullString

	err = tx.QueryRow(
		query,
		time.Now(),  // updated_at
		name,        // name
		description, // description
		isActive,    // is_active
		configJSON,  // config
		gatewayID,   // id
	).Scan(
		&gateway.ID,
		&gateway.Name,
		&gateway.Description,
		&gateway.Type,
		&gateway.IsActive,
		&rawConfig,
		&gateway.CreatedAt,
		&gateway.UpdatedAt,
		&linkedAt,
		&disabledReason,
		&disabledAt,
	)

	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update gateway: %v", err)
	}

	// Unmarshal the config
	if len(rawConfig) > 0 {
		if err := json.Unmarshal(rawConfig, &gateway.Config); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to unmarshal config: %v", err)
		}
	}

	// Handle nullable fields
	if linkedAt.Valid {
		gateway.LinkedAt = linkedAt.Time
	}
	if disabledAt.Valid {
		gateway.DisabledAt = &disabledAt.Time
	}
	if disabledReason.Valid {
		gateway.DisabledReason = disabledReason.String
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &gateway, nil
}
