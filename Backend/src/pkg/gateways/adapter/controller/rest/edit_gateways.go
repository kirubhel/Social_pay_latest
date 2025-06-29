package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

type UpdateGatewayRequest struct {
	GatewayID   uuid.UUID             `json:"gateway_id"`
	Name        *string               `json:"name,omitempty"`
	Description *string               `json:"description,omitempty"`
	IsActive    *bool                 `json:"is_active,omitempty"`
	Config      *entity.GatewayConfig `json:"config,omitempty"`
}

type UpdateGatewayResponse struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message,omitempty"`
	Data      *entity.PaymentGateway `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Error     *Error                 `json:"error,omitempty"`
}

func (controller *Controller) UpdateGateway(w http.ResponseWriter, r *http.Request) {
	var req UpdateGatewayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, UpdateGatewayResponse{
			Success:   false,
			Message:   "Invalid request payload",
			Timestamp: time.Now(),
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate Gateway ID
	if req.GatewayID == uuid.Nil {
		SendJSONResponse(w, UpdateGatewayResponse{
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

	// Check if at least one field is being updated
	if req.Name == nil && req.Description == nil && req.IsActive == nil && req.Config == nil {
		SendJSONResponse(w, UpdateGatewayResponse{
			Success:   false,
			Message:   "No fields provided for update",
			Timestamp: time.Now(),
			Error: &Error{
				Type:    "INVALID_INPUT",
				Message: "At least one field must be provided for update",
			},
		}, http.StatusBadRequest)
		return
	}

	// Call interactor to update gateway
	updatedGateway, err := controller.interactor.UpdateGateway(
		req.GatewayID,
		req.Name,
		req.Description,
		req.IsActive,
		req.Config,
	)
	if err != nil {
		controller.log.Printf("Failed to update gateway %s: %v", req.GatewayID, err)

		status := http.StatusInternalServerError
		errType := "INTERNAL_ERROR"

		if errors.Is(err, entity.ErrGatewayNotFound) {
			status = http.StatusNotFound
			errType = "NOT_FOUND"
		} else if errors.Is(err, entity.ErrInvalidConfig) {
			status = http.StatusBadRequest
			errType = "INVALID_CONFIG"
		}

		SendJSONResponse(w, UpdateGatewayResponse{
			Success:   false,
			Message:   "Failed to update gateway",
			Timestamp: time.Now(),
			Error: &Error{
				Type:    errType,
				Message: err.Error(),
			},
		}, status)
		return
	}

	// Success response
	SendJSONResponse(w, UpdateGatewayResponse{
		Success:   true,
		Message:   "Gateway updated successfully",
		Data:      updatedGateway,
		Timestamp: time.Now(),
	}, http.StatusOK)
}
