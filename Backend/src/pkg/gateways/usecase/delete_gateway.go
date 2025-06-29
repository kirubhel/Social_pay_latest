package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (uc Usecase) DeleteGateway(gatewayID uuid.UUID) error {
	const (
		ErrGatewayDeletionFailed = "GATEWAY_DELETION_FAILED"
		ErrGatewayNotFound       = "GATEWAY_NOT_FOUND"
		ErrGatewayInUse          = "GATEWAY_IN_USE"
	)

	// Validate gateway ID
	if gatewayID == uuid.Nil {
		return Error{
			Type:    ErrGatewayDeletionFailed,
			Message: "invalid gateway ID",
		}
	}

	// Delete the gateway
	var err error
	err = uc.repo.DeleteGateway(gatewayID)
	if err != nil {
		return Error{
			Type:    ErrGatewayDeletionFailed,
			Message: fmt.Sprintf("failed to delete gateway: %v", err),
		}
	}

	// Log the deletion
	uc.log.Printf("Gateway %s deleted successfully at %v", gatewayID, time.Now())

	return nil
}
