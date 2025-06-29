package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (uc Usecase) ListMerchantGateways(
	merchantID uuid.UUID,
	includeDisabled bool,
	gatewayType string,
) ([]entity.GatewayMerchant, error) {
	const (
		ErrListFailed       = "LIST_MERCHANT_GATEWAYS_FAILED"
		ErrMerchantNotFound = "MERCHANT_NOT_FOUND"
	)

	// Validate merchant ID
	if merchantID == uuid.Nil {
		return nil, Error{
			Type:    ErrListFailed,
			Message: "merchant ID must be provided",
		}
	}

	// Validate gateway type if provided
	if gatewayType != "" {
		validTypes := map[string]bool{
			"bank":           true,
			"mobile_money":   true,
			"card_processor": true,
		}
		if !validTypes[gatewayType] {
			return nil, Error{
				Type:    ErrListFailed,
				Message: "invalid gateway type specified",
			}
		}
	}

	// Get merchant gateways from repository
	gateways, err := uc.repo.ListMerchantGateways(merchantID, includeDisabled, gatewayType)
	if err != nil {
		if errors.Is(err, entity.ErrMerchantNotFound) {
			return nil, Error{
				Type:    ErrMerchantNotFound,
				Message: "specified merchant does not exist",
			}
		}
		return nil, Error{
			Type:    ErrListFailed,
			Message: "failed to retrieve merchant gateways",
		}
	}

	// Log the operation
	uc.log.Println(
		"Retrieved merchant gateways",
		"merchant_id", merchantID,
		"count", len(gateways),
		"timestamp", time.Now(),
	)

	return gateways, nil
}
