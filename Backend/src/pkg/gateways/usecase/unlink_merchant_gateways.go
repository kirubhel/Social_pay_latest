package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (uc Usecase) UnlinkGatewayFromMerchant(merchantID uuid.UUID, gatewayID uuid.UUID) (time.Time, error) {
	const (
		ErrUnlinkFailed     = "UNLINK_GATEWAY_FAILED"
		ErrNotLinked        = "NOT_LINKED"
		ErrMerchantNotFound = "MERCHANT_NOT_FOUND"
		ErrGatewayNotFound  = "GATEWAY_NOT_FOUND"
	)

	// Validate UUIDs are not empty
	if merchantID == uuid.Nil || gatewayID == uuid.Nil {
		return time.Time{}, Error{
			Type:    ErrUnlinkFailed,
			Message: "both merchant ID and gateway ID must be provided",
		}
	}

	// Remove the link in repository
	unlinkedAt, err := uc.repo.UnlinkGatewayFromMerchant(merchantID, gatewayID)
	if err != nil {
		// Map repository errors to domain errors
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
		default:
			return time.Time{}, Error{
				Type:    ErrUnlinkFailed,
				Message: "failed to unlink merchant from gateway",
			}
		}
	}

	// Log successful unlinking
	uc.log.Println(
		"Unlinked merchant from gateway",
		"merchant_id", merchantID,
		"gateway_id", gatewayID,
		"unlinked_at", unlinkedAt,
	)

	return unlinkedAt, nil
}
