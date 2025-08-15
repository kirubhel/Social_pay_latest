package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

// LinkGatewayRequest represents the request payload
type LinkGatewayRequest struct {
	MerchantID string `json:"merchant_id"`
	GatewayID  string `json:"gateway_id"`
}

// LinkGatewayResponse represents the response structure
type LinkGatewayResponse struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message,omitempty"`
	LinkedAt   time.Time `json:"linked_at,omitempty"`
	MerchantID uuid.UUID `json:"merchant_id,omitempty"`
	GatewayID  uuid.UUID `json:"gateway_id,omitempty"`
	Error      *Error    `json:"error,omitempty"`
}

func (controller *Controller) LinkGatewayToMerchant(w http.ResponseWriter, r *http.Request) {
	var req LinkGatewayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, LinkGatewayResponse{
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
		SendJSONResponse(w, LinkGatewayResponse{
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
		SendJSONResponse(w, LinkGatewayResponse{
			Success: false,
			Error: &Error{
				Type:    "INVALID_INPUT",
				Message: "Invalid gateway ID format",
			},
		}, http.StatusBadRequest)
		return
	}

	// Call interactor to create the link
	linkedAt, err := controller.interactor.LinkGatewayToMerchant(merchantUUID, gatewayUUID)
	if err != nil {
		status := http.StatusInternalServerError
		errType := "INTERNAL_ERROR"

		switch {
		case errors.Is(err, entity.ErrMerchantNotFound):
			status = http.StatusNotFound
			errType = "NOT_FOUND"
		case errors.Is(err, entity.ErrGatewayNotFound):
			status = http.StatusNotFound
			errType = "NOT_FOUND"
		case errors.Is(err, entity.ErrAlreadyLinked):
			status = http.StatusConflict
			errType = "CONFLICT"
		}

		SendJSONResponse(w, LinkGatewayResponse{
			Success: false,
			Error: &Error{
				Type:    errType,
				Message: err.Error(),
			},
		}, status)
		return
	}

	// Success response
	SendJSONResponse(w, LinkGatewayResponse{
		Success:    true,
		Message:    "Gateway successfully linked to merchant",
		LinkedAt:   linkedAt,
		MerchantID: merchantUUID,
		GatewayID:  gatewayUUID,
	}, http.StatusCreated)
}
