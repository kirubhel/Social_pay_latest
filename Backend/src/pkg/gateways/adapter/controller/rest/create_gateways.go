package rest

import (
	"encoding/json"
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

type GatewayResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type CreateGatewayRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"type"`
	Config      entity.GatewayConfig `json:"config"`
	IsActive    bool                 `json:"is_active"`
}

func (controller *Controller) CreateGateway(w http.ResponseWriter, r *http.Request) {
	var req CreateGatewayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, GatewayResponse{
			Success: false,
			Error: &Error{
				Message: "Invalid request payload",
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Description == "" || req.Type == "" {
		SendJSONResponse(w, GatewayResponse{
			Success: false,
			Error: &Error{
				Message: "Name, description and type are required fields",
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate gateway type
	validTypes := map[string]bool{
		"bank":           true,
		"mobile_money":   true,
		"card_processor": true,
	}
	if !validTypes[req.Type] {
		SendJSONResponse(w, GatewayResponse{
			Success: false,
			Error: &Error{
				Message: "Invalid gateway type. Must be one of: bank, mobile_money, card_processor",
			},
		}, http.StatusBadRequest)
		return
	}

	// Create gateway using interactor
	gateway, err := controller.interactor.CreateGateway(
		req.Name,
		req.Description,
		req.Type,
		req.Config,
		req.IsActive,
	)
	if err != nil {
		controller.log.Printf("Failed to create gateway: %v", err)
		SendJSONResponse(w, GatewayResponse{
			Success: false,
			Error: &Error{
				Message: "Failed to create payment gateway",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Success response
	SendJSONResponse(w, GatewayResponse{
		Success: true,
		Data:    gateway,
	}, http.StatusCreated)
}
