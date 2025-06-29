package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (uc Usecase) EnableMerchantGateway(
	merchantID uuid.UUID,
	gatewayID uuid.UUID,
	reason string,
) (time.Time, error) {
	const (
		ErrEnableFailed     = "ENABLE_GATEWAY_FAILED"
		ErrAlreadyEnabled   = "ALREADY_ENABLED"
		ErrNotLinked        = "NOT_LINKED"
		ErrMerchantNotFound = "MERCHANT_NOT_FOUND"
		ErrGatewayNotFound  = "GATEWAY_NOT_FOUND"
	)

	// Validate UUIDs
	if merchantID == uuid.Nil || gatewayID == uuid.Nil {
		return time.Time{}, Error{
			Type:    ErrEnableFailed,
			Message: "both merchant ID and gateway ID must be provided",
		}
	}

	// Enable the gateway for merchant
	enabledAt, err := uc.repo.EnableMerchantGateway(merchantID, gatewayID, reason)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrMerchantNotFound):
			return time.Time{}, Error{
				Type:    ErrMerchantNotFound,
				Message: "specified merchant does not exist",
			}
		case errors.Is(err, entity.ErrGatewayNotFound):
			return time.Time{}, Error{
				Type:    ErrGatewayNotFound,
				Message: "specified gateway does not exist",
			}
		case errors.Is(err, entity.ErrNotLinked):
			return time.Time{}, Error{
				Type:    ErrNotLinked,
				Message: "these merchant and gateway are not currently linked",
			}
		case errors.Is(err, entity.ErrAlreadyEnabled):
			return time.Time{}, Error{
				Type:    ErrAlreadyEnabled,
				Message: "this merchant gateway is already enabled",
			}
		default:
			return time.Time{}, Error{
				Type:    ErrEnableFailed,
				Message: "failed to enable merchant gateway",
			}
		}
	}

	// Log the operation
	uc.log.Println(
		"Enabled merchant gateway",
		"merchant_id", merchantID,
		"gateway_id", gatewayID,
		"reason", reason,
		"enabled_at", enabledAt,
	)

	return enabledAt, nil
}
