package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/qr/core/entity"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
)

// QRRepository defines the interface for QR link operations
type QRRepository interface {
	// Create creates a new QR link
	Create(ctx context.Context, qrLink *entity.QRLink) error

	// GetByID retrieves a QR link by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.QRLink, error)

	// GetByMerchant retrieves QR links for a merchant with pagination
	GetByMerchant(ctx context.Context, merchantID uuid.UUID, pagination *pagination.Pagination) ([]entity.QRLink, int64, error)

	// GetByUser retrieves QR links for a user with pagination
	GetByUser(ctx context.Context, userID uuid.UUID, pagination *pagination.Pagination) ([]entity.QRLink, int64, error)

	// Update updates an existing QR link
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates *entity.UpdateQRLinkRequest) (*entity.QRLink, error)

	// Delete soft deletes a QR link
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
