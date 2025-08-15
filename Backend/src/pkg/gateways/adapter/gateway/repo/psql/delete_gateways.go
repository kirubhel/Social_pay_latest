package repo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (repo PsqlRepo) DeleteGateway(gatewayID uuid.UUID) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// First check if the gateway exists
	var exists bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM merchants.payment_gateways WHERE id = $1)
	`, gatewayID).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check gateway existence: %v", err)
	}
	if !exists {
		tx.Rollback()
		return entity.ErrGatewayNotFound
	}

	// Check if the gateway is in use (has linked merchants)
	var inUse bool
	err = tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM merchants.gateway_merchants 
			WHERE gateway_id = $1 AND is_active = true
		)
	`, gatewayID).Scan(&inUse)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check gateway usage: %v", err)
	}
	if inUse {
		tx.Rollback()
		return entity.ErrGatewayInUse
	}

	// Delete the gateway
	_, err = tx.Exec(`
		DELETE FROM merchants.payment_gateways 
		WHERE id = $1
	`, gatewayID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete gateway: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
