package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (uc Usecase) LinkGatewayToMerchant(merchantID uuid.UUID, gatewayID uuid.UUID) (time.Time, error) {
	const (
		ErrLinkFailed       = "LINK_GATEWAY_FAILED"
		ErrAlreadyLinked    = "ALREADY_LINKED"
		ErrMerchantNotFound = "MERCHANT_NOT_FOUND"
		ErrGatewayNotFound  = "GATEWAY_NOT_FOUND"
	)

	// Validate UUIDs are not empty
	if merchantID == uuid.Nil || gatewayID == uuid.Nil {
		return time.Time{}, Error{
			Type:    ErrLinkFailed,
			Message: "both merchant ID and gateway ID must be provided",
		}
	}

	// Create the link in repository
	linkedAt, err := uc.repo.LinkGatewayToMerchant(merchantID, gatewayID)
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
		case errors.Is(err, entity.ErrAlreadyLinked):
			return time.Time{}, Error{
				Type:    ErrAlreadyLinked,
				Message: "this merchant and gateway are already linked",
			}
		default:
			return time.Time{}, Error{
				Type:    ErrLinkFailed,
				Message: "failed to create merchant-gateway link",
			}
		}
	}

	// Log successful linking
	uc.log.Println(
		"Linked merchant to gateway",
		"merchant_id", merchantID,
		"gateway_id", gatewayID,
		"linked_at", linkedAt,
	)

	return linkedAt, nil
}
