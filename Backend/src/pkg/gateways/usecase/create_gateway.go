package usecase

import (
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (uc Usecase) CreateGateway(
	name string,
	description string,
	gatewayType string,
	config entity.GatewayConfig,
	isActive bool,
) (*entity.ListPaymentGateway, error) {
	const ErrGatewayCreationFailed = "GATEWAY_CREATION_FAILED"

	// Validate gateway type
	validTypes := map[string]bool{
		"bank":           true,
		"mobile_money":   true,
		"card_processor": true,
	}
	if !validTypes[gatewayType] {
		return nil, Error{
			Type:    ErrGatewayCreationFailed,
			Message: "Invalid gateway type",
		}
	}

	// Create the gateway entity
	gateway := entity.ListPaymentGateway{
		Name:        name,
		Description: description,
		Type:        gatewayType,
		Config:      config,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Store the gateway
	err := uc.repo.CreateGateway(gateway)
	if err != nil {
		return nil, Error{
			Type:    ErrGatewayCreationFailed,
			Message: fmt.Sprintf("Failed to store gateway: %v", err),
		}
	}

	return &gateway, nil
}
