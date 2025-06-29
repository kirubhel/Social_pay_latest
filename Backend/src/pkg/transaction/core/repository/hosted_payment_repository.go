package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

// HostedPaymentRepository defines the interface for hosted payment operations
type HostedPaymentRepository interface {
	// Create creates a new hosted payment
	Create(ctx context.Context, hostedPayment *entity.HostedPayment) error

	// GetByID retrieves a hosted payment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.HostedPayment, error)

	// GetByReference retrieves a hosted payment by reference and merchant ID
	GetByReference(ctx context.Context, reference string, merchantID uuid.UUID) (*entity.HostedPayment, error)

	// ValidateReferenceId validates that the reference is unique for the merchant
	ValidateReferenceId(ctx context.Context, merchantID uuid.UUID, reference string) error

	// Update updates a hosted payment
	Update(ctx context.Context, hostedPayment *entity.HostedPayment) error

	// UpdateWithTransaction updates hosted payment with transaction details
	UpdateWithTransaction(ctx context.Context, id uuid.UUID, transactionID uuid.UUID, selectedMedium string, selectedPhoneNumber string, status entity.HostedPaymentStatus) error

	// UpdateStatus updates the status of a hosted payment
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.HostedPaymentStatus) error

	// GetExpiredPayments retrieves expired hosted payments
	GetExpiredPayments(ctx context.Context) ([]entity.HostedPayment, error)
}
