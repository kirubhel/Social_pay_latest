package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

type EnableGatewayRequest struct {
	MerchantID string `json:"merchant_id"`
	GatewayID  string `json:"gateway_id"`
	Reason     string `json:"reason,omitempty"`
}

type EnableGatewayResponse struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message,omitempty"`
	EnabledAt  time.Time `json:"Enabled_at,omitempty"`
	MerchantID uuid.UUID `json:"merchant_id,omitempty"`
	GatewayID  uuid.UUID `json:"gateway_id,omitempty"`
	IsActive   bool      `json:"is_active,omitempty"`
	Error      *Error    `json:"error,omitempty"`
}

func (controller *Controller) EnableMerchantGateway(w http.ResponseWriter, r *http.Request) {
	var req EnableGatewayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, EnableGatewayResponse{
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
		SendJSONResponse(w, EnableGatewayResponse{
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
		SendJSONResponse(w, EnableGatewayResponse{
			Success: false,
			Error: &Error{
				Type:    "INVALID_INPUT",
				Message: "Invalid gateway ID format",
			},
		}, http.StatusBadRequest)
		return
	}

	// Call interactor to Enable
	EnabledAt, err := controller.interactor.EnableMerchantGateway(
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
		case errors.Is(err, entity.ErrAlreadyEnabled):
			status = http.StatusConflict
			errType = "CONFLICT"
		}

		SendJSONResponse(w, EnableGatewayResponse{
			Success: false,
			Error: &Error{
				Type:    errType,
				Message: err.Error(),
			},
		}, status)
		return
	}

	// Success response
	SendJSONResponse(w, EnableGatewayResponse{
		Success:    true,
		Message:    "Merchant gateway Enabled successfully",
		EnabledAt:  EnabledAt,
		MerchantID: merchantUUID,
		GatewayID:  gatewayUUID,
		IsActive:   false,
	}, http.StatusOK)
}
