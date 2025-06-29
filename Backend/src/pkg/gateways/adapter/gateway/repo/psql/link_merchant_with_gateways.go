package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (repo PsqlRepo) LinkGatewayToMerchant(merchantID uuid.UUID, gatewayID uuid.UUID) (time.Time, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Check if merchant exists
	var merchantExists bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM merchants.merchants WHERE id = $1)
	`, merchantID).Scan(&merchantExists)
	if err != nil {
		tx.Rollback()
		return time.Time{}, fmt.Errorf("failed to check merchant existence: %v", err)
	}
	if !merchantExists {
		tx.Rollback()
		return time.Time{}, entity.ErrMerchantNotFound
	}

	// Check if gateway exists
	var gatewayExists bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM merchants.payment_gateways WHERE id = $1)
	`, gatewayID).Scan(&gatewayExists)
	if err != nil {
		tx.Rollback()
		return time.Time{}, fmt.Errorf("failed to check gateway existence: %v", err)
	}
	if !gatewayExists {
		tx.Rollback()
		return time.Time{}, entity.ErrGatewayNotFound
	}

	// Check if link already exists and is active
	var linkExists bool
	var isActive bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM merchants.gateway_merchants WHERE merchant_id = $1 AND gateway_id = $2),
		(SELECT is_active FROM merchants.gateway_merchants WHERE merchant_id = $1 AND gateway_id = $2)
	`, merchantID, gatewayID).Scan(&linkExists, &isActive)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return time.Time{}, fmt.Errorf("failed to check existing link: %v", err)
	}

	if linkExists && isActive {
		tx.Rollback()
		return time.Time{}, entity.ErrAlreadyLinked
	}

	linkedAt := time.Now()

	// Upsert the link (insert or update if exists but inactive)
	_, err = tx.Exec(`
		INSERT INTO merchants.gateway_merchants (merchant_id, gateway_id, is_active, linked_at, disabled_at)
		VALUES ($1, $2, true, $3, NULL)
		ON CONFLICT (merchant_id, gateway_id) 
		DO UPDATE SET 
			is_active = true,
			linked_at = $3,
			disabled_at = NULL,
			disabled_reason = NULL
	`, merchantID, gatewayID, linkedAt)
	if err != nil {
		tx.Rollback()
		return time.Time{}, fmt.Errorf("failed to create merchant-gateway link: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return time.Time{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return linkedAt, nil
}
