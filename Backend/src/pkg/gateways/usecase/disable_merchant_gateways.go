package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (uc Usecase) DisableMerchantGateway(
	merchantID uuid.UUID,
	gatewayID uuid.UUID,
	reason string,
) (time.Time, error) {
	const (
		ErrDisableFailed    = "DISABLE_GATEWAY_FAILED"
		ErrAlreadyDisabled  = "ALREADY_DISABLED"
		ErrNotLinked        = "NOT_LINKED"
		ErrMerchantNotFound = "MERCHANT_NOT_FOUND"
		ErrGatewayNotFound  = "GATEWAY_NOT_FOUND"
	)

	// Validate UUIDs
	if merchantID == uuid.Nil || gatewayID == uuid.Nil {
		return time.Time{}, Error{
			Type:    ErrDisableFailed,
			Message: "both merchant ID and gateway ID must be provided",
		}
	}

	// Disable the gateway for merchant
	disabledAt, err := uc.repo.DisableMerchantGateway(merchantID, gatewayID, reason)
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
		case errors.Is(err, entity.ErrAlreadyDisabled):
			return time.Time{}, Error{
				Type:    ErrAlreadyDisabled,
				Message: "this merchant gateway is already disabled",
			}
		default:
			return time.Time{}, Error{
				Type:    ErrDisableFailed,
				Message: "failed to disable merchant gateway",
			}
		}
	}

	// Log the operation
	uc.log.Println(
		"Disabled merchant gateway",
		"merchant_id", merchantID,
		"gateway_id", gatewayID,
		"reason", reason,
		"disabled_at", disabledAt,
	)

	return disabledAt, nil
}
