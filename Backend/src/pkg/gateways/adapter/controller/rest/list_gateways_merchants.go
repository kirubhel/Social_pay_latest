package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

// GatewayMerchantsRequest defines the request parameters
type GatewayMerchantsRequest struct {
	GatewayID       string `json:"gateway_id" validate:"required,uuid4"`
	IncludeInactive bool   `json:"include_inactive"`
	Limit           int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset          int    `json:"offset" validate:"omitempty,min=0"`
}

// GatewayMerchantsResponse defines the response structure
type GatewayMerchantsResponse struct {
	Success     bool              `json:"success"`
	GatewayID   uuid.UUID         `json:"gateway_id"`
	TotalCount  int               `json:"total_count"`
	ActiveCount int               `json:"active_count"`
	Merchants   []GatewayMerchant `json:"merchants"`
	Timestamp   time.Time         `json:"timestamp"`
	Error       *Error            `json:"error,omitempty"`
}

type GatewayMerchant struct {
	MerchantID uuid.UUID  `json:"merchant_id"`
	Name       string     `json:"name"`
	BusinessID string     `json:"business_id"`
	IsActive   bool       `json:"is_active"`
	LinkedAt   time.Time  `json:"linked_at"`
	DisabledAt *time.Time `json:"disabled_at,omitempty"`
	Merchants entity.Merchant
}

func (controller *Controller) ListGatewayMerchants(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request
	var req GatewayMerchantsRequest
	if err := controller.ParseRequest(r, &req); err != nil {
		SendJSONResponse(w, GatewayMerchantsResponse{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	gatewayID := uuid.MustParse(req.GatewayID)

	// Call interactor
	result, err := controller.interactor.GetGatewayMerchants(
		gatewayID,
		req.IncludeInactive,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		status := http.StatusInternalServerError
		errType := "INTERNAL_ERROR"

		if errors.Is(err, entity.ErrGatewayNotFound) {
			status = http.StatusNotFound
			errType = "NOT_FOUND"
		}

		SendJSONResponse(w, GatewayMerchantsResponse{
			Success: false,
			Error: &Error{
				Type:    errType,
				Message: err.Error(),
			},
		}, status)
		return
	}

	// Prepare response
	response := GatewayMerchantsResponse{
		Success:     true,
		GatewayID:   gatewayID,
		Merchants:   make([]GatewayMerchant, 0, len(result.Merchants)),
		Timestamp:   time.Now(),
	}

	for _, m := range result.Merchants {
		response.Merchants = append(response.Merchants, GatewayMerchant{
			MerchantID: m.MerchantID,
			Name:       m.Name,
			BusinessID: m.BusinessID,
			IsActive:   m.IsActive,
			LinkedAt:   m.LinkedAt,
			DisabledAt: m.DisabledAt,
		})
	}

	SendJSONResponse(w, response, http.StatusOK)
}
