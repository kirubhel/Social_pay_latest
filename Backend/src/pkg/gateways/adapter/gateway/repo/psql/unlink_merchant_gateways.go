package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (repo PsqlRepo) UnlinkGatewayFromMerchant(merchantID uuid.UUID, gatewayID uuid.UUID) (time.Time, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if merchant exists
	var merchantExists bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM merchants.merchants WHERE id = $1)
	`, merchantID).Scan(&merchantExists)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to check merchant existence: %v", err)
	}
	if !merchantExists {
		return time.Time{}, entity.ErrMerchantNotFound
	}

	// Check if gateway exists
	var gatewayExists bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM merchants.payment_gateways WHERE id = $1)
	`, gatewayID).Scan(&gatewayExists)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to check gateway existence: %v", err)
	}
	if !gatewayExists {
		return time.Time{}, entity.ErrGatewayNotFound
	}

	// Check if link exists and is active
	var linkExists bool
	var isActive bool
	err = tx.QueryRow(`
		SELECT 
			EXISTS(SELECT 1 FROM merchants.gateway_merchants WHERE merchant_id = $1 AND gateway_id = $2),
			(SELECT is_active FROM merchants.gateway_merchants WHERE merchant_id = $1 AND gateway_id = $2)
	`, merchantID, gatewayID).Scan(&linkExists, &isActive)
	if err != nil && err != sql.ErrNoRows {
		return time.Time{}, fmt.Errorf("failed to check existing link: %v", err)
	}

	if !linkExists {
		return time.Time{}, entity.ErrNotLinked
	}

	unlinkedAt := time.Now()

	// Update the link to mark as inactive
	_, err = tx.Exec(`
		UPDATE merchants.gateway_merchants 
		SET 
			is_active = false,
			disabled_at = $1,
			disabled_reason = 'Manually unlinked'
		WHERE merchant_id = $2 AND gateway_id = $3
	`, unlinkedAt, merchantID, gatewayID)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to unlink gateway from merchant: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return time.Time{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return unlinkedAt, nil
}
