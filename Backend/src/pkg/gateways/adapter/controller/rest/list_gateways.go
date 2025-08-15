package rest

import (
	"net/http"
	"time"

	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

type ListSocialPayGatewaysResponse struct {
	Success   bool                        `json:"success"`
	Count     int                         `json:"count"`
	Data      []entity.ListPaymentGateway `json:"data"`
	Timestamp time.Time                   `json:"timestamp"`
	Error     *Error                      `json:"error,omitempty"`
}

// mapToListPaymentGateways converts []entity.PaymentGateway to []entity.ListPaymentGateway
func mapToListPaymentGateways(gateways []entity.PaymentGateway) []entity.ListPaymentGateway {
	var listGateways []entity.ListPaymentGateway
	for _, gateway := range gateways {
		listGateways = append(listGateways, entity.ListPaymentGateway{
			ID:          gateway.ID,
			Name:        gateway.Name,
			Type:        gateway.Type,
			IsActive:    gateway.IsActive,
			CreatedAt:   gateway.CreatedAt,
			UpdatedAt:   gateway.UpdatedAt,
			Description: gateway.Description,
			Config: entity.GatewayConfig{
				APIKey:         gateway.Config.APIKey,
				SecretKey:      gateway.Config.SecretKey,
				BaseURL:        gateway.Config.BaseURL,
				WebhookURL:     gateway.Config.WebhookURL,
				TransactionFee: gateway.Config.TransactionFee,
				IsTest:         gateway.Config.IsTest,
			},
		})
	}
	return listGateways
}

func (controller *Controller) ListGateways(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	controller.log.Println("Controller: Starting to list gateways")

	gateways, err := controller.interactor.ListAllGateways()
	if err != nil {
		controller.log.Printf("Controller: Failed to list gateways: %v", err)

		SendJSONResponse(w, ListSocialPayGatewaysResponse{
			Success:   false,
			Timestamp: time.Now(),
			Error: &Error{
				Type:    "INTERNAL_ERROR",
				Message: "Failed to retrieve gateways",
				Details: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}
	// Sanitize sensitive data before response
	controller.log.Println("Controller: Sanitizing sensitive data")
	for i := range gateways {
		gateways[i].Config.APIKey = ""
		gateways[i].Config.SecretKey = ""
	}

	duration := time.Since(startTime)
	controller.log.Printf("Controller: Successfully processed request in %v", duration)

	SendJSONResponse(w, ListSocialPayGatewaysResponse{
		Success:   true,
		Count:     len(gateways),
		Data:      mapToListPaymentGateways(gateways),
		Timestamp: time.Now(),
	}, http.StatusOK)
}
