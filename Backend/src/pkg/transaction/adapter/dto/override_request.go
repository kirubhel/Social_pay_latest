package dto

import (
	"errors"
	"strings"

	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

// OverrideTransactionStatusRequest represents the request to manually override a transaction status
type OverrideTransactionStatusRequest struct {
	TransactionID string                   `json:"transactionId" form:"transactionId" binding:"required" validate:"uuid"`
	Status        entity.TransactionStatus `json:"status" form:"status" binding:"required"`
	Reason        string                   `json:"reason" form:"reason" binding:"required,min=10"`
	AdminID       string                   `json:"adminId" form:"adminId" binding:"required"`
}

// Validate performs validation on the override request
func (r *OverrideTransactionStatusRequest) Validate() error {
	// Validate required fields
	if strings.TrimSpace(r.TransactionID) == "" {
		return errors.New("transaction ID is required")
	}

	if strings.TrimSpace(string(r.Status)) == "" {
		return errors.New("status is required")
	}

	if strings.TrimSpace(r.Reason) == "" {
		return errors.New("reason is required")
	}

	if len(strings.TrimSpace(r.Reason)) < 10 {
		return errors.New("reason must be at least 10 characters long")
	}

	if strings.TrimSpace(r.AdminID) == "" {
		return errors.New("admin ID is required")
	}

	// Validate status is one of the allowed values
	validStatuses := []entity.TransactionStatus{
		entity.SUCCESS,
		entity.FAILED,
	}

	isValidStatus := false
	for _, validStatus := range validStatuses {
		if r.Status == validStatus {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		return errors.New("invalid status. Allowed statuses: SUCCESS, FAILED, CANCELED, REFUNDED")
	}

	return nil
}

// GetMessage returns a formatted message for the webhook indicating manual override
func (r *OverrideTransactionStatusRequest) GetMessage() string {
	return "Transaction status manually overridden by admin. Reason: " + r.Reason + " | Admin ID: " + r.AdminID
}

// GetProviderTxID returns a formatted provider transaction ID for manual override
func (r *OverrideTransactionStatusRequest) GetProviderTxID() string {
	return "MANUAL_OVERRIDE_" + r.AdminID
}

// GetProviderData returns formatted provider data for manual override
func (r *OverrideTransactionStatusRequest) GetProviderData() string {
	return `{"override_type": "manual", "admin_id": "` + r.AdminID + `", "reason": "` + r.Reason + `"}`
}
