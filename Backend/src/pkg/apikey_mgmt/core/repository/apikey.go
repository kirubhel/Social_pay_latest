package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/entity"
)

// Repository defines the interface for API key storage operations
type Repository interface {
	// CreateAPIKey creates a new API key
	CreateAPIKey(ctx context.Context, apiKey *entity.APIKey) (*entity.APIKey, error)

	// GetAPIKey retrieves an API key by its ID
	GetAPIKey(ctx context.Context, id uuid.UUID) (*entity.APIKey, error)

	// GetAPIKeyByPublicKey retrieves an API key by its public key
	GetAPIKeyByPublicKey(ctx context.Context, publicKey string) (*entity.APIKey, error)

	// UpdateAPIKey updates an existing API key
	UpdateAPIKey(ctx context.Context, id uuid.UUID, request entity.UpdateAPIKeyRequest) (*entity.APIKey, error)

	// DeleteAPIKey deletes an API key
	DeleteAPIKey(ctx context.Context, id uuid.UUID) error

	// ListAPIKeys retrieves all API keys for a user
	ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]*entity.APIKey, error)

	// ListMerchantAPIKeys retrieves all API keys for a merchant
	ListMerchantAPIKeys(ctx context.Context, merchantID uuid.UUID) ([]*entity.APIKey, error)

	// ValidateAPIKey validates an API key's public and secret key combination
	ValidateAPIKey(ctx context.Context, request entity.APIKeyValidateRequest) (*entity.APIKey, error)

	// RotateAPIKeySecret updates the secret key for an API key
	RotateAPIKeySecret(ctx context.Context, id uuid.UUID, newSecret string) (*entity.APIKey, error)
}
