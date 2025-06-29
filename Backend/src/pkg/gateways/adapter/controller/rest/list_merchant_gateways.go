package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

// MerchantGatewaysResponse represents the response structure
type MerchantGatewaysResponse struct {
	Success     bool                  `json:"success"`
	MerchantID  uuid.UUID             `json:"merchant_id"`
	Count       int                   `json:"count"`
	ActiveCount int                   `json:"active_count"`
	Gateways    []MerchantGatewayInfo `json:"gateways"`
	Timestamp   time.Time             `json:"timestamp"`
	Error       *Error                `json:"error,omitempty"`
}

// MerchantGatewaysRequest defines the expected request parameters
type MerchantGatewaysRequest struct {
	MerchantID      string `json:"merchant_id" validate:"required,uuid4"`
	IncludeDisabled bool   `json:"include_disabled"`
	Type            string `json:"type" validate:"omitempty,oneof=bank mobile_money card_processor"`
}

type MerchantGatewayInfo struct {
	MerchantID     uuid.UUID  `json:"merchant_id"`
	GatewayID      uuid.UUID  `json:"gateway_id"`
	Name           string     `json:"name"`
	Type           string     `json:"type"`
	IsActive       bool       `json:"is_active"`
	LinkedAt       time.Time  `json:"linked_at"`
	DisabledAt     *time.Time `json:"disabled_at,omitempty"`
	DisabledReason string     `json:"disabled_reason,omitempty"`
}

func (controller *Controller) ListMerchantGateways(w http.ResponseWriter, r *http.Request) {
	// Parse request parameters
	var req MerchantGatewaysRequest
	if err := controller.ParseRequest(r, &req); err != nil {
		SendJSONResponse(w, MerchantGatewaysResponse{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Convert to UUID (validation already done by parseRequest)
	merchantID := uuid.MustParse(req.MerchantID)

	// Call interactor
	gateways, err := controller.interactor.ListMerchantGateways(
		merchantID,
		req.IncludeDisabled,
		req.Type,
	)
	if err != nil {
		status := http.StatusInternalServerError
		errType := "INTERNAL_ERROR"

		if errors.Is(err, entity.ErrMerchantNotFound) {
			status = http.StatusNotFound
			errType = "NOT_FOUND"
		}

		SendJSONResponse(w, MerchantGatewaysResponse{
			Success: false,
			Error: &Error{
				Type:    errType,
				Message: err.Error(),
			},
		}, status)
		return
	}
	// Count active gateways
	activeCount := 0
	for _, gw := range gateways {
		if gw.IsActive {
			activeCount++
		}
	}

	// Prepare response
	response := MerchantGatewaysResponse{
		Success:     true,
		MerchantID:  merchantID,
		Count:       len(gateways),
		ActiveCount: activeCount,
		Gateways:    make([]MerchantGatewayInfo, 0, len(gateways)),
		Timestamp:   time.Now(),
	}

	for _, gw := range gateways {
		response.Gateways = append(response.Gateways, MerchantGatewayInfo{
			MerchantID: gw.MerchantID,
			Name:       gw.Name,
			Type:       gw.Type,
			IsActive:   gw.IsActive,
			LinkedAt:   gw.LinkedAt,
			DisabledAt: gw.DisabledAt,
			DisabledReason: func() string {
				if gw.DisabledReason != nil {
					return *gw.DisabledReason
				}
				return ""
			}(),
		})
	}

	SendJSONResponse(w, response, http.StatusOK)
}

func (controller *Controller) ParseRequest(r *http.Request, req interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		return errors.New("failed to parse request")
	}
	return nil
}
