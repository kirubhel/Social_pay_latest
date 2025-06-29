package repo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (repo PsqlRepo) EnableMerchantGateway(
	merchantID uuid.UUID,
	gatewayID uuid.UUID,
	reason string,
) (time.Time, error) {
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

	// Check if the link exists and is active
	var linkExists bool
	var isActive bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM merchants.gateway_merchants 
		WHERE merchant_id = $1 AND gateway_id = $2),
		(SELECT is_active FROM merchants.gateway_merchants 
		WHERE merchant_id = $1 AND gateway_id = $2)
	`, merchantID, gatewayID).Scan(&linkExists, &isActive)
	if err != nil {
		tx.Rollback()
		return time.Time{}, fmt.Errorf("failed to check gateway-merchant link: %v", err)
	}
	if !linkExists {
		tx.Rollback()
		return time.Time{}, entity.ErrNotLinked
	}
	if !isActive {
		tx.Rollback()
		return time.Time{}, entity.ErrAlreadyDisabled
	}

	// Update the link to disabled
	disabledAt := time.Now()
	_, err = tx.Exec(`
		UPDATE merchants.gateway_merchants 
		SET 
			is_active = true,
			enabled_at = $1,
			enbabled_reason = $2
		WHERE 
			merchant_id = $3 AND 
			gateway_id = $4
	`, disabledAt, reason, merchantID, gatewayID)
	if err != nil {
		tx.Rollback()
		return time.Time{}, fmt.Errorf("failed to disable merchant gateway: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return time.Time{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return disabledAt, nil
}
