package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)


type DeleteGatewayRequest struct {
	GatewayID uuid.UUID `json:"gateway_id"`
}

type DeleteGatewayResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Error     *Error    `json:"error,omitempty"`
}

func (controller *Controller) DeleteGateway(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req DeleteGatewayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, DeleteGatewayResponse{
			Success:   false,
			Message:   "Invalid request payload",
			Timestamp: time.Now(),
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Failed to parse request body",
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate gateway ID
	if req.GatewayID == uuid.Nil {
		SendJSONResponse(w, DeleteGatewayResponse{
			Success:   false,
			Message:   "Validation failed",
			Timestamp: time.Now(),
			Error: &Error{
				Type:    "INVALID_INPUT",
				Message: "Gateway ID is required",
			},
		}, http.StatusBadRequest)
		return
	}

	// Call interactor to delete gateway
	err := controller.interactor.DeleteGateway(req.GatewayID)
	if err != nil {
		controller.log.Printf("Failed to delete gateway %s: %v", req.GatewayID, err)

		status := http.StatusInternalServerError
		errType := "INTERNAL_ERROR"

		if errors.Is(err, entity.ErrGatewayNotFound) {
			status = http.StatusNotFound
			errType = "NOT_FOUND"
		} else if errors.Is(err, entity.ErrGatewayInUse) {
			status = http.StatusConflict
			errType = "CONFLICT"
		}

		SendJSONResponse(w, DeleteGatewayResponse{
			Success:   false,
			Message:   "Failed to delete gateway",
			Timestamp: time.Now(),
			Error: &Error{
				Type:    errType,
				Message: err.Error(),
			},
		}, status)
		return
	}

	// Success response
	SendJSONResponse(w, DeleteGatewayResponse{
		Success:   true,
		Message:   fmt.Sprintf("Gateway %s deleted successfully", req.GatewayID),
		Timestamp: time.Now(),
	}, http.StatusOK)
}
