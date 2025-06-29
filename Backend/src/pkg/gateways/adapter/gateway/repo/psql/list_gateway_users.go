package repo

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

type GatewayMerchantsResult struct {
	Merchants   []entity.GatewayMerchant
	TotalCount  int
	ActiveCount int
}

func (repo PsqlRepo) GetGatewayMerchants(
	gatewayID uuid.UUID,
	includeInactive bool,
	limit int,
	offset int,
) (*entity.GatewayMerchantsResult, error) {
	// First check if the gateway exists
	var exists bool
	err := repo.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM merchants.payment_gateways WHERE id = $1)
	`, gatewayID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check gateway existence: %v", err)
	}
	if !exists {
		return nil, entity.ErrGatewayNotFound
	}

	// Build the base query
	query := `
		SELECT 
			gm.merchant_id,
			m.legal_name,
			m.business_reg_number,
			gm.is_active,
			gm.linked_at,
			gm.disabled_at
		FROM merchants.gateway_merchants gm
		JOIN merchants.merchants m ON gm.merchant_id = m.id
		WHERE gm.gateway_id = $1
	`

	// Add condition for active/inactive
	if !includeInactive {
		query += " AND gm.is_active = true"
	}

	// Add pagination
	query += " ORDER BY gm.linked_at DESC LIMIT $2 OFFSET $3"

	// Execute the query
	rows, err := repo.db.Query(query, gatewayID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query gateway merchants: %v", err)
	}
	defer rows.Close()

	var merchants []entity.GatewayMerchant
	for rows.Next() {
		var gm entity.GatewayMerchant
		var disabledAt sql.NullTime

		err := rows.Scan(
			&gm.MerchantID,
			&gm.Name,
			&gm.BusinessID,
			&gm.IsActive,
			&gm.LinkedAt,
			&disabledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan merchant row: %v", err)
		}

		if disabledAt.Valid {
			gm.DisabledAt = &disabledAt.Time
		}

		merchants = append(merchants, gm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %v", err)
	}

	// Get total count (without pagination)
	var totalCount int
	countQuery := `
		SELECT COUNT(*) 
		FROM merchants.gateway_merchants 
		WHERE gateway_id = $1
	`
	if !includeInactive {
		countQuery += " AND is_active = true"
	}

	err = repo.db.QueryRow(countQuery, gatewayID).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	// Get active count
	var activeCount int
	err = repo.db.QueryRow(`
		SELECT COUNT(*) 
		FROM merchants.gateway_merchants 
		WHERE gateway_id = $1 AND is_active = true
	`, gatewayID).Scan(&activeCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get active count: %v", err)
	}

	return &entity.GatewayMerchantsResult{
		Merchants:   merchants,
		TotalCount:  totalCount,
		ActiveCount: activeCount,
	}, nil
}
