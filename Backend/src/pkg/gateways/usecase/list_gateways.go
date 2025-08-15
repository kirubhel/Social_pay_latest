package usecase

import (
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/gateways/core/entity"
)

func (uc Usecase) ListAllGateways() ([]entity.PaymentGateway, error) {
	const ErrGatewayListFailed = "GATEWAY_LIST_FAILED"
	startTime := time.Now()
	uc.log.Println("Usecase: Starting to list all gateways")

	gateways, err := uc.repo.ListAllGateways()
	if err != nil {
		uc.log.Printf("Usecase: Failed to retrieve gateways: %v", err)
		return nil, Error{
			Type:    ErrGatewayListFailed,
			Message: fmt.Sprintf("Failed to retrieve gateways list: %v", err),
			Details: err.Error(),
		}
	}

	duration := time.Since(startTime)
	uc.log.Printf(
		"Usecase: Successfully retrieved %d gateways in %v",
		len(gateways),
		duration,
	)

	return gateways, nil
}
