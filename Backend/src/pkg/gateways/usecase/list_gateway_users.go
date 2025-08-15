package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

type GatewayMerchantsResult struct {
	Merchants   []entity.GatewayMerchant
	TotalCount  int
	ActiveCount int
}

func (uc Usecase) GetGatewayMerchants(
	gatewayID uuid.UUID,
	includeInactive bool,
	limit int,
	offset int,
) (*entity.PaymentGateway, error) {
	const (
		ErrListFailed      = "LIST_GATEWAY_MERCHANTS_FAILED"
		ErrGatewayNotFound = "GATEWAY_NOT_FOUND"
	)

	// Validate gateway ID
	if gatewayID == uuid.Nil {
		return nil, Error{
			Type:    ErrListFailed,
			Message: "gateway ID must be provided",
		}
	}

	// Set default values for pagination
	if limit == 0 {
		limit = 50 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	// Get merchants from repository
	merchantsSlice, err := uc.repo.GetGatewayMerchants(
		gatewayID,
		includeInactive,
		limit,
		offset,
	)
	if err != nil {
		if errors.Is(err, entity.ErrGatewayNotFound) {
			return nil, Error{
				Type:    ErrGatewayNotFound,
				Message: "specified gateway does not exist",
			}
		}
		return nil, Error{
			Type:    ErrListFailed,
			Message: "failed to retrieve gateway merchants: " + err.Error(),
		}
	}
	merchants := merchantsSlice.Merchants // Ensure merchants is a slice
	if err != nil {
		if errors.Is(err, entity.ErrGatewayNotFound) {
			return nil, Error{
				Type:    ErrGatewayNotFound,
				Message: "specified gateway does not exist",
			}
		}
		return nil, Error{
			Type:    ErrListFailed,
			Message: "failed to retrieve gateway merchants: " + err.Error(),
		}
	}
	// Calculate active count
	activeCount := 0
	for _, m := range merchants {
		if m.IsActive {
			activeCount++
		}
	}

	// Log the operation
	uc.log.Println(
		"Retrieved gateway merchants",
		"gateway_id", gatewayID,
		"total_count", len(merchants),
		"timestamp", time.Now(),
	)

	return &entity.PaymentGateway{
		Merchants: merchants,
	}, nil
}
