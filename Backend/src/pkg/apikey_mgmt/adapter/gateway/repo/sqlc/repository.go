package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/entity"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/repository"

	"github.com/google/uuid"
)

// Repository implements the repository.Repository interface for API key management
type Repository struct {
	q *Queries
}

// NewRepository creates a new API key repository
func NewRepository(db *sql.DB) repository.Repository {
	return &Repository{
		q: New(db),
	}
}

// CreateAPIKey creates a new API key
func (r *Repository) CreateAPIKey(ctx context.Context, apiKey *entity.APIKey) (*entity.APIKey, error) {
	var expiresAt sql.NullTime
	if apiKey.ExpiresAt != nil {
		expiresAt = sql.NullTime{
			Time:  *apiKey.ExpiresAt,
			Valid: true,
		}
	}
	fmt.Println("apiKey.MerchantID", apiKey.MerchantID)
	params := CreateAPIKeyParams{
		ID:                apiKey.ID,
		UserID:            apiKey.UserID,
		CreatedBy:         apiKey.CreatedBy,
		MerchantID:        uuid.NullUUID{UUID: apiKey.MerchantID, Valid: true},
		Name:              apiKey.Name,
		Description:       sql.NullString{String: apiKey.Description, Valid: true},
		PublicKey:         apiKey.PublicKey,
		SecretKey:         apiKey.SecretKey,
		CanWithdrawal:     apiKey.CanWithdrawal,
		CanProcessPayment: apiKey.CanProcessPayment,
		CreatedAt:         apiKey.CreatedAt,
		UpdatedAt:         apiKey.UpdatedAt,
		ExpiresAt:         expiresAt,
		IsActive:          apiKey.IsActive,
	}

	result, err := r.q.CreateAPIKey(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return convertToEntity(result), nil
}

// GetAPIKey gets an API key by ID
func (r *Repository) GetAPIKey(ctx context.Context, id uuid.UUID) (*entity.APIKey, error) {
	result, err := r.q.GetAPIKey(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return convertToEntity(result), nil
}

// GetAPIKeyByPublicKey gets an API key by public key
func (r *Repository) GetAPIKeyByPublicKey(ctx context.Context, publicKey string) (*entity.APIKey, error) {
	result, err := r.q.GetAPIKeyByPublicKey(ctx, publicKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get API key by public key: %w", err)
	}

	return convertToEntity(result), nil
}

// ListAPIKeys gets all API keys for a user
func (r *Repository) ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]*entity.APIKey, error) {
	results, err := r.q.ListAPIKeys(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	apiKeys := make([]*entity.APIKey, len(results))
	for i, result := range results {
		apiKeys[i] = convertToEntity(result)
	}

	return apiKeys, nil
}

// ListMerchantAPIKeys gets all API keys for a merchant
func (r *Repository) ListMerchantAPIKeys(ctx context.Context, merchantID uuid.UUID) ([]*entity.APIKey, error) {
	results, err := r.q.ListMerchantAPIKeys(ctx, uuid.NullUUID{UUID: merchantID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	apiKeys := make([]*entity.APIKey, len(results))
	for i, result := range results {
		apiKeys[i] = convertToEntity(result)
	}

	return apiKeys, nil
}

// UpdateAPIKey updates an API key
func (r *Repository) UpdateAPIKey(ctx context.Context, id uuid.UUID, request entity.UpdateAPIKeyRequest) (*entity.APIKey, error) {
	params := UpdateAPIKeyParams{
		ID: id,
	}

	if request.Name != nil {
		params.Name = *request.Name
	}

	if request.Description != nil {
		params.Description = sql.NullString{String: *request.Description, Valid: true}
	}

	if request.CanWithdrawal != nil {
		params.CanWithdrawal = *request.CanWithdrawal
	}

	if request.CanProcessPayment != nil {
		params.CanProcessPayment = *request.CanProcessPayment
	}

	if request.IsActive != nil {
		params.IsActive = *request.IsActive
	}

	expiresAt := request.GetExpiresAt()
	if expiresAt != nil {
		params.ExpiresAt = sql.NullTime{Time: *expiresAt, Valid: true}
	}

	result, err := r.q.UpdateAPIKey(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update API key: %w", err)
	}

	return convertToEntity(result), nil
}

// DeleteAPIKey deletes an API key
func (r *Repository) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteAPIKey(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}

// ValidateAPIKey validates an API key's public and secret key combination
func (r *Repository) ValidateAPIKey(ctx context.Context, request entity.APIKeyValidateRequest) (*entity.APIKey, error) {
	// Get API key by public key and secret key
	result, err := r.q.ValidateAPIKey(ctx, ValidateAPIKeyParams{
		PublicKey: request.PublicKey,
		SecretKey: request.SecretKey,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid API key")
		}
		return nil, fmt.Errorf("failed to validate API key: %w", err)
	}

	apiKey := convertToEntity(result)

	// Check if API key has expired
	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return nil, errors.New("API key has expired")
	}

	// Update last used at
	result, err = r.q.UpdateLastUsedAt(ctx, apiKey.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update last used at: %w", err)
	}
	return convertToEntity(result), nil
}

// RotateAPIKeySecret rotates an API key's secret
func (r *Repository) RotateAPIKeySecret(ctx context.Context, id uuid.UUID, newSecret string) (*entity.APIKey, error) {
	params := RotateAPIKeySecretParams{
		ID:        id,
		SecretKey: newSecret,
	}

	result, err := r.q.RotateAPIKeySecret(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to rotate API key secret: %w", err)
	}

	return convertToEntity(result), nil
}

// convertToEntity converts an SQLC model to an entity
func convertToEntity(model ApiKey) *entity.APIKey {
	var expiresAt *time.Time
	if model.ExpiresAt.Valid {
		expiresAt = &model.ExpiresAt.Time
	}

	var lastUsedAt *time.Time
	if model.LastUsedAt.Valid {
		lastUsedAt = &model.LastUsedAt.Time
	}

	return &entity.APIKey{
		ID:                model.ID,
		UserID:            model.UserID,
		MerchantID:        model.MerchantID.UUID,
		CreatedBy:         model.CreatedBy,
		Name:              model.Name,
		Description:       model.Description.String,
		PublicKey:         model.PublicKey,
		SecretKey:         model.SecretKey,
		CanWithdrawal:     model.CanWithdrawal,
		CanProcessPayment: model.CanProcessPayment,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		ExpiresAt:         expiresAt,
		LastUsedAt:        lastUsedAt,
		IsActive:          model.IsActive,
	}
}
