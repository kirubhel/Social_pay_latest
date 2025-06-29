package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

type DisableGatewayRequest struct {
	MerchantID string `json:"merchant_id"`
	GatewayID  string `json:"gateway_id"`
	Reason     string `json:"reason,omitempty"`
}

type DisableGatewayResponse struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message,omitempty"`
	DisabledAt time.Time `json:"disabled_at,omitempty"`
	MerchantID uuid.UUID `json:"merchant_id,omitempty"`
	GatewayID  uuid.UUID `json:"gateway_id,omitempty"`
	IsActive   bool      `json:"is_active,omitempty"`
	Error      *Error    `json:"error,omitempty"`
}

func (controller *Controller) DisableMerchantGateway(w http.ResponseWriter, r *http.Request) {
	var req DisableGatewayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, DisableGatewayResponse{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request payload",
			},
		}, http.StatusBadRequest)
		return
	}

	// Convert string IDs to UUID
	merchantUUID, err := uuid.Parse(req.MerchantID)
	if err != nil {
		SendJSONResponse(w, DisableGatewayResponse{
			Success: false,
			Error: &Error{
				Type:    "INVALID_INPUT",
				Message: "Invalid merchant ID format",
			},
		}, http.StatusBadRequest)
		return
	}

	gatewayUUID, err := uuid.Parse(req.GatewayID)
	if err != nil {
		SendJSONResponse(w, DisableGatewayResponse{
			Success: false,
			Error: &Error{
				Type:    "INVALID_INPUT",
				Message: "Invalid gateway ID format",
			},
		}, http.StatusBadRequest)
		return
	}

	// Call interactor to disable
	disabledAt, err := controller.interactor.DisableMerchantGateway(
		merchantUUID,
		gatewayUUID,
		req.Reason,
	)
	if err != nil {
		status := http.StatusInternalServerError
		errType := "INTERNAL_ERROR"

		switch {
		case errors.Is(err, entity.ErrNotLinked):
			status = http.StatusNotFound
			errType = "NOT_FOUND"
		case errors.Is(err, entity.ErrMerchantNotFound):
			status = http.StatusNotFound
			errType = "NOT_FOUND"
		case errors.Is(err, entity.ErrGatewayNotFound):
			status = http.StatusNotFound
			errType = "NOT_FOUND"
		case errors.Is(err, entity.ErrAlreadyDisabled):
			status = http.StatusConflict
			errType = "CONFLICT"
		}

		SendJSONResponse(w, DisableGatewayResponse{
			Success: false,
			Error: &Error{
				Type:    errType,
				Message: err.Error(),
			},
		}, status)
		return
	}

	// Success response
	SendJSONResponse(w, DisableGatewayResponse{
		Success:    true,
		Message:    "Merchant gateway disabled successfully",
		DisabledAt: disabledAt,
		MerchantID: merchantUUID,
		GatewayID:  gatewayUUID,
		IsActive:   false,
	}, http.StatusOK)
}
