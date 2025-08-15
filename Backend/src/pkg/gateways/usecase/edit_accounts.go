package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (uc Usecase) UpdateGateway(
	gatewayID uuid.UUID,
	name *string,
	description *string,
	isActive *bool,
	config *entity.GatewayConfig,
) (*entity.PaymentGateway, error) {
	const (
		ErrGatewayUpdateFailed = "GATEWAY_UPDATE_FAILED"
		ErrGatewayNotFound     = "GATEWAY_NOT_FOUND"
		ErrInvalidConfig       = "INVALID_CONFIG"
		ErrNoUpdatesProvided   = "NO_UPDATES_PROVIDED"
	)

	// Validate gateway ID
	if gatewayID == uuid.Nil {
		return nil, Error{
			Type:    ErrGatewayUpdateFailed,
			Message: "invalid gateway ID",
		}
	}

	// Check if at least one field is being updated
	if name == nil && description == nil && isActive == nil && config == nil {
		return nil, Error{
			Type:    ErrNoUpdatesProvided,
			Message: "at least one field must be provided for update",
		}
	}
	// Save updates
	updatedGateway, err := uc.repo.UpdateGateway(gatewayID, name, description, isActive, config)
	if err != nil {
		return nil, Error{
			Type:    ErrGatewayUpdateFailed,
			Message: fmt.Sprintf("failed to update gateway: %v", err),
		}
	}

	return updatedGateway, nil
}
