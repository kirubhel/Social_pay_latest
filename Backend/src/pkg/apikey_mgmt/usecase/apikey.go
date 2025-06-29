package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/entity"
	"github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/repository"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"

	"github.com/google/uuid"
)

// APIKeyUseCase defines the API key management use case interface
type APIKeyUseCase interface {
	CreateAPIKey(ctx context.Context, userID, createdBy uuid.UUID, merchantID uuid.UUID, request entity.CreateAPIKeyRequest) (*entity.APIKeyResponse, string, error)
	GetAPIKey(ctx context.Context, id uuid.UUID) (*entity.APIKeyResponse, error)
	GetAPIKeysByUserID(ctx context.Context, userID uuid.UUID) ([]entity.APIKeyResponse, error)
	GetAPIKeysByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]entity.APIKeyResponse, error)
	UpdateAPIKey(ctx context.Context, id uuid.UUID, request entity.UpdateAPIKeyRequest) (*entity.APIKeyResponse, error)
	DeleteAPIKey(ctx context.Context, id uuid.UUID) error
	ValidateAPIKey(ctx context.Context, publicKey, secretKey string) (*entity.APIKeyResponse, error)
	RotateAPIKeySecret(ctx context.Context, id uuid.UUID) (*entity.APIKeyRotateResponse, error)
}

type apiKeyUseCase struct {
	log  logging.Logger
	repo repository.Repository
}

// NewAPIKeyUseCase creates a new API key management use case
func NewAPIKeyUseCase(repo repository.Repository) APIKeyUseCase {
	return &apiKeyUseCase{
		log:  logging.NewStdLogger("[APIKEY]"),
		repo: repo,
	}
}

// CreateAPIKey creates a new API key
func (u *apiKeyUseCase) CreateAPIKey(ctx context.Context, userID, createdBy uuid.UUID, merchantID uuid.UUID, request entity.CreateAPIKeyRequest) (*entity.APIKeyResponse, string, error) {
	// Generate public and secret keys
	publicKeyBase, err := generateRandomString(16) // 16 chars
	if err != nil {
		u.log.Error("Failed to generate public key", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, "", fmt.Errorf("failed to generate public key: %w", err)
	}
	publicKey := "SocialPUB_" + publicKeyBase

	secretKeyBase, err := generateRandomString(24) // 24 chars
	if err != nil {
		u.log.Error("Failed to generate secret key", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, "", fmt.Errorf("failed to generate secret key: %w", err)
	}
	secretKey := "SocialSEC_" + secretKeyBase

	// Set default values for permissions if not provided
	canWithdrawal := false
	if request.CanWithdrawal != nil {
		canWithdrawal = *request.CanWithdrawal
	}

	canProcessPayment := true
	if request.CanProcessPayment != nil {
		canProcessPayment = *request.CanProcessPayment
	}

	// Create API key
	apiKey := &entity.APIKey{
		ID:                uuid.New(),
		UserID:            userID,
		MerchantID:        merchantID,
		CreatedBy:         createdBy,
		Name:              request.Name,
		Description:       request.Description,
		PublicKey:         publicKey,
		SecretKey:         secretKey,
		CanWithdrawal:     canWithdrawal,
		CanProcessPayment: canProcessPayment,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		ExpiresAt:         request.GetExpiresAt(),
		IsActive:          true,
	}

	// Save to repository
	apiKey, err = u.repo.CreateAPIKey(ctx, apiKey)
	if err != nil {
		u.log.Error("Failed to create API key", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		return nil, "", fmt.Errorf("failed to create API key: %w", err)
	}

	// Create response
	response := &entity.APIKeyResponse{
		ID:                apiKey.ID,
		UserID:            apiKey.UserID,
		Name:              apiKey.Name,
		Description:       apiKey.Description,
		PublicKey:         apiKey.PublicKey,
		SecretKey:         secretKey,
		CanWithdrawal:     apiKey.CanWithdrawal,
		CanProcessPayment: apiKey.CanProcessPayment,
		CreatedAt:         apiKey.CreatedAt,
		UpdatedAt:         apiKey.UpdatedAt,
		ExpiresAt:         apiKey.ExpiresAt,
		LastUsedAt:        apiKey.LastUsedAt,
		IsActive:          apiKey.IsActive,
	}

	return response, secretKey, nil
}

// GetAPIKey gets an API key by ID
func (u *apiKeyUseCase) GetAPIKey(ctx context.Context, id uuid.UUID) (*entity.APIKeyResponse, error) {
	apiKey, err := u.repo.GetAPIKey(ctx, id)
	if err != nil {
		u.log.Error("Failed to get API key", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if apiKey == nil {
		return nil, errors.New("API key not found")
	}

	response := &entity.APIKeyResponse{
		ID:                apiKey.ID,
		UserID:            apiKey.UserID,
		Name:              apiKey.Name,
		Description:       apiKey.Description,
		SecretKey:         apiKey.SecretKey,
		CanWithdrawal:     apiKey.CanWithdrawal,
		CanProcessPayment: apiKey.CanProcessPayment,
		PublicKey:         apiKey.PublicKey,
		CreatedAt:         apiKey.CreatedAt,
		UpdatedAt:         apiKey.UpdatedAt,
		ExpiresAt:         apiKey.ExpiresAt,
		LastUsedAt:        apiKey.LastUsedAt,
		IsActive:          apiKey.IsActive,
	}

	return response, nil
}

// GetAPIKeysByUserID gets all API keys for a user
func (u *apiKeyUseCase) GetAPIKeysByUserID(ctx context.Context, userID uuid.UUID) ([]entity.APIKeyResponse, error) {
	apiKeys, err := u.repo.ListAPIKeys(ctx, userID)
	if err != nil {
		u.log.Error("Failed to get API keys for user", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		return nil, fmt.Errorf("failed to get API keys for user: %w", err)
	}

	response := make([]entity.APIKeyResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		response[i] = entity.APIKeyResponse{
			ID:          apiKey.ID,
			UserID:      apiKey.UserID,
			Name:        apiKey.Name,
			Description: apiKey.Description,
			PublicKey:   apiKey.PublicKey,
			SecretKey:   apiKey.SecretKey,
			MerchantID:  apiKey.MerchantID,
			CreatedAt:   apiKey.CreatedAt,
			UpdatedAt:   apiKey.UpdatedAt,
			ExpiresAt:   apiKey.ExpiresAt,
			LastUsedAt:  apiKey.LastUsedAt,
			IsActive:    apiKey.IsActive,
		}
	}

	return response, nil
}

// GetAPIKeysByMerchantID gets all API keys for a merchant
func (u *apiKeyUseCase) GetAPIKeysByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]entity.APIKeyResponse, error) {
	apiKeys, err := u.repo.ListMerchantAPIKeys(ctx, merchantID)
	if err != nil {
		u.log.Error("Failed to get API keys for merchant", map[string]interface{}{
			"error":      err.Error(),
			"merchantID": merchantID,
		})
		return nil, fmt.Errorf("failed to get API keys for merchant: %w", err)
	}

	response := make([]entity.APIKeyResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		response[i] = entity.APIKeyResponse{
			ID:                apiKey.ID,
			UserID:            apiKey.UserID,
			Name:              apiKey.Name,
			Description:       apiKey.Description,
			CanWithdrawal:     apiKey.CanWithdrawal,
			CanProcessPayment: apiKey.CanProcessPayment,
			PublicKey:         apiKey.PublicKey,
			SecretKey:         apiKey.SecretKey,
			MerchantID:        apiKey.MerchantID,
			CreatedAt:         apiKey.CreatedAt,
			UpdatedAt:         apiKey.UpdatedAt,
			ExpiresAt:         apiKey.ExpiresAt,
			LastUsedAt:        apiKey.LastUsedAt,
			IsActive:          apiKey.IsActive,
		}
	}

	return response, nil
}

// UpdateAPIKey updates an API key
func (u *apiKeyUseCase) UpdateAPIKey(ctx context.Context, id uuid.UUID, request entity.UpdateAPIKeyRequest) (*entity.APIKeyResponse, error) {
	apiKey, err := u.repo.UpdateAPIKey(ctx, id, request)
	if err != nil {
		u.log.Error("Failed to update API key", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to update API key: %w", err)
	}

	if apiKey == nil {
		return nil, errors.New("API key not found")
	}

	response := &entity.APIKeyResponse{
		ID:                apiKey.ID,
		UserID:            apiKey.UserID,
		Name:              apiKey.Name,
		Description:       apiKey.Description,
		PublicKey:         apiKey.PublicKey,
		SecretKey:         apiKey.SecretKey,
		CanWithdrawal:     apiKey.CanWithdrawal,
		CanProcessPayment: apiKey.CanProcessPayment,
		CreatedAt:         apiKey.CreatedAt,
		UpdatedAt:         apiKey.UpdatedAt,
		ExpiresAt:         apiKey.ExpiresAt,
		LastUsedAt:        apiKey.LastUsedAt,
		IsActive:          apiKey.IsActive,
	}

	return response, nil
}

// RotateAPIKeySecret rotates the secret key for an API key
func (u *apiKeyUseCase) RotateAPIKeySecret(ctx context.Context, id uuid.UUID) (*entity.APIKeyRotateResponse, error) {
	// Generate new secret key
	secretKey, err := generateRandomString(64)
	if err != nil {
		u.log.Error("Failed to generate new secret key", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to generate new secret key: %w", err)
	}

	// Update API key with new secret
	apiKey, err := u.repo.RotateAPIKeySecret(ctx, id, secretKey)
	if err != nil {
		u.log.Error("Failed to rotate API key secret", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to rotate API key secret: %w", err)
	}

	// Create response
	response := &entity.APIKeyRotateResponse{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}

	return response, nil
}

// DeleteAPIKey deletes an API key
func (u *apiKeyUseCase) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	err := u.repo.DeleteAPIKey(ctx, id)
	if err != nil {
		u.log.Error("Failed to delete API key", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}

// ValidateAPIKey validates an API key
func (u *apiKeyUseCase) ValidateAPIKey(ctx context.Context, publicKey, secretKey string) (*entity.APIKeyResponse, error) {
	request := entity.APIKeyValidateRequest{
		PublicKey: publicKey,
		SecretKey: secretKey,
	}

	apiKey, err := u.repo.ValidateAPIKey(ctx, request)
	fmt.Printf("[ValidateAPIKey] REPO API Key: %+v\n", apiKey)
	if err != nil {
		u.log.Error("Failed to validate API key", map[string]interface{}{
			"error":     err.Error(),
			"publicKey": publicKey,
		})
		return nil, fmt.Errorf("failed to validate API key: %w", err)
	}

	if apiKey == nil {
		return nil, errors.New("invalid API key")
	}

	response := &entity.APIKeyResponse{
		ID:                apiKey.ID,
		UserID:            apiKey.UserID,
		MerchantID:        apiKey.MerchantID,
		Name:              apiKey.Name,
		Description:       apiKey.Description,
		PublicKey:         apiKey.PublicKey,
		SecretKey:         apiKey.SecretKey,
		CanWithdrawal:     apiKey.CanWithdrawal,
		CanProcessPayment: apiKey.CanProcessPayment,
		CreatedAt:         apiKey.CreatedAt,
		UpdatedAt:         apiKey.UpdatedAt,
		ExpiresAt:         apiKey.ExpiresAt,
		LastUsedAt:        apiKey.LastUsedAt,
		IsActive:          apiKey.IsActive,
	}

	return response, nil
}

// generateRandomString generates a random string of the specified length
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}
