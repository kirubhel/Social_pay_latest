package repo

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (repo PsqlRepo) ListMerchantGateways(
	merchantID uuid.UUID,
	includeDisabled bool,
	gatewayType string,
) ([]entity.GatewayMerchant, error) {
	// First check if merchant exists
	var merchantExists bool
	err := repo.db.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM merchants WHERE id = $1)
    `, merchantID).Scan(&merchantExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check merchant existence: %v", err)
	}
	if !merchantExists {
		return nil, entity.ErrMerchantNotFound
	}

	// Build base query
	query := `
        SELECT 
            pg.id,
            pg.name,
            pg.type,
            gm.is_active,
            gm.linked_at,
            gm.disabled_at,
            gm.disabled_reason,
            m.legal_name,
            m.business_reg_number
        FROM merchants.payment_gateways pg
        JOIN merchants.gateway_merchants gm ON pg.id = gm.gateway_id
        JOIN merchants.merchants m ON gm.merchant_id = m.id
        WHERE gm.merchant_id = $1
    `

	// Add conditions based on parameters
	args := []interface{}{merchantID}
	argPos := 2

	if !includeDisabled {
		query += fmt.Sprintf(" AND gm.is_active = true")
	}

	if gatewayType != "" {
		query += fmt.Sprintf(" AND pg.type = $%d", argPos)
		args = append(args, gatewayType)
		argPos++
	}

	query += " ORDER BY gm.linked_at DESC"

	// Execute query
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query merchant gateways: %v", err)
	}
	defer rows.Close()

	var gatewayMerchants []entity.GatewayMerchant
	for rows.Next() {
		var gm entity.PaymentGateway
		var disabledAt sql.NullTime
		var disabledReason sql.NullString

		err := rows.Scan(
			&gm.ID,
			&gm.Name,
			&gm.Type,
			&gm.IsActive,
			&gm.LinkedAt,
			&disabledAt,
			&disabledReason,
			&gm.Merchant.LegalName,
			&gm.Merchant.BusinessRegNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan gateway merchant row: %v", err)
		}

		if disabledAt.Valid {
			gm.DisabledAt = &disabledAt.Time
		}

		if disabledReason.Valid {
			gm.DisabledReason = disabledReason.String
		}
		gatewayMerchants = append(gatewayMerchants, entity.GatewayMerchant{
			Name:       gm.Name,
			IsActive:   gm.IsActive,
			LinkedAt:   gm.LinkedAt,
			DisabledAt: gm.DisabledAt,
			DisabledReason: func(s string) *string {
				if s == "" {
					return nil
				}
				return &s
			}(gm.DisabledReason),
			// Removed Merchant field as it does not exist in GatewayMerchant struct
		})

		// Append the processed gateway merchant to the list
		continue
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %v", err)
	}

	return gatewayMerchants, nil
}
