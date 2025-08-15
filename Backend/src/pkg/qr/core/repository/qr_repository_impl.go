package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/socialpay/socialpay/src/pkg/qr/core/entity"
	db "github.com/socialpay/socialpay/src/pkg/qr/core/repository/generated"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

type QRRepositoryImpl struct {
	Queries *db.Queries
}

func NewQRRepository(dbConn *sql.DB) QRRepository {
	return &QRRepositoryImpl{
		Queries: db.New(dbConn),
	}
}

func (r *QRRepositoryImpl) Create(ctx context.Context, qrLink *entity.QRLink) error {
	supportedMethodsJSON, err := json.Marshal(qrLink.SupportedMethods)
	if err != nil {
		return fmt.Errorf("failed to marshal supported methods: %w", err)
	}

	var amount sql.NullString
	if qrLink.Amount != nil {
		amount = sql.NullString{String: fmt.Sprintf("%.2f", *qrLink.Amount), Valid: true}
	}

	var title, description, imageURL sql.NullString
	if qrLink.Title != nil {
		title = sql.NullString{String: *qrLink.Title, Valid: true}
	}
	if qrLink.Description != nil {
		description = sql.NullString{String: *qrLink.Description, Valid: true}
	}
	if qrLink.ImageURL != nil {
		imageURL = sql.NullString{String: *qrLink.ImageURL, Valid: true}
	}

	params := db.CreateQRLinkParams{
		ID:               qrLink.ID,
		UserID:           qrLink.UserID,
		MerchantID:       qrLink.MerchantID,
		Type:             db.QrLinkType(qrLink.Type),
		Amount:           amount,
		SupportedMethods: supportedMethodsJSON,
		Tag:              db.QrLinkTag(qrLink.Tag),
		Title:            title,
		Description:      description,
		ImageUrl:         imageURL,
		IsTipEnabled:     sql.NullBool{Bool: qrLink.IsTipEnabled, Valid: true},
	}

	dbQRLink, err := r.Queries.CreateQRLink(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create QR link: %w", err)
	}

	// Update the entity with database values
	qrLink.CreatedAt = dbQRLink.CreatedAt
	qrLink.UpdatedAt = dbQRLink.UpdatedAt

	return nil
}

func (r *QRRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.QRLink, error) {
	dbQRLink, err := r.Queries.GetQRLink(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("QR link not found")
		}
		return nil, fmt.Errorf("failed to get QR link: %w", err)
	}

	return r.mapDBQRLinkToEntity(dbQRLink)
}

func (r *QRRepositoryImpl) GetByMerchant(ctx context.Context, merchantID uuid.UUID, pag *pagination.Pagination) ([]entity.QRLink, int64, error) {
	// Get QR links
	dbQRLinks, err := r.Queries.GetQRLinksByMerchant(ctx, db.GetQRLinksByMerchantParams{
		MerchantID: merchantID,
		Limit:      int32(pag.GetLimit()),
		Offset:     int32(pag.GetOffset()),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get QR links by merchant: %w", err)
	}

	// Get total count
	total, err := r.Queries.CountQRLinksByMerchant(ctx, merchantID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count QR links by merchant: %w", err)
	}

	qrLinks := make([]entity.QRLink, len(dbQRLinks))
	for i, dbQRLink := range dbQRLinks {
		qrLink, err := r.mapDBQRLinkToEntity(dbQRLink)
		if err != nil {
			return nil, 0, err
		}
		qrLinks[i] = *qrLink
	}

	return qrLinks, total, nil
}

func (r *QRRepositoryImpl) GetByUser(ctx context.Context, userID uuid.UUID, pag *pagination.Pagination) ([]entity.QRLink, int64, error) {
	// Get QR links
	dbQRLinks, err := r.Queries.GetQRLinksByUser(ctx, db.GetQRLinksByUserParams{
		UserID: userID,
		Limit:  int32(pag.GetLimit()),
		Offset: int32(pag.GetOffset()),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get QR links by user: %w", err)
	}

	// Get total count
	total, err := r.Queries.CountQRLinksByUser(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count QR links by user: %w", err)
	}

	qrLinks := make([]entity.QRLink, len(dbQRLinks))
	for i, dbQRLink := range dbQRLinks {
		qrLink, err := r.mapDBQRLinkToEntity(dbQRLink)
		if err != nil {
			return nil, 0, err
		}
		qrLinks[i] = *qrLink
	}

	return qrLinks, total, nil
}

func (r *QRRepositoryImpl) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates *entity.UpdateQRLinkRequest) (*entity.QRLink, error) {
	// First get the current QR link to use existing values for fields that aren't being updated
	currentQRLink, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get current QR link: %w", err)
	}

	var supportedMethodsJSON json.RawMessage
	if updates.SupportedMethods != nil {
		bytes, err := json.Marshal(updates.SupportedMethods)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal supported methods: %w", err)
		}
		supportedMethodsJSON = bytes
	} else {
		// Use existing supported methods
		bytes, err := json.Marshal(currentQRLink.SupportedMethods)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal existing supported methods: %w", err)
		}
		supportedMethodsJSON = bytes
	}

	var amount sql.NullString
	if updates.Amount != nil {
		amount = sql.NullString{String: fmt.Sprintf("%.2f", *updates.Amount), Valid: true}
	} else if currentQRLink.Amount != nil {
		amount = sql.NullString{String: fmt.Sprintf("%.2f", *currentQRLink.Amount), Valid: true}
	}

	tag := db.QrLinkTag(currentQRLink.Tag)
	if updates.Tag != nil {
		tag = db.QrLinkTag(*updates.Tag)
	}

	var title, description, imageURL sql.NullString
	if updates.Title != nil {
		title = sql.NullString{String: *updates.Title, Valid: true}
	} else if currentQRLink.Title != nil {
		title = sql.NullString{String: *currentQRLink.Title, Valid: true}
	}

	if updates.Description != nil {
		description = sql.NullString{String: *updates.Description, Valid: true}
	} else if currentQRLink.Description != nil {
		description = sql.NullString{String: *currentQRLink.Description, Valid: true}
	}

	if updates.ImageURL != nil {
		imageURL = sql.NullString{String: *updates.ImageURL, Valid: true}
	} else if currentQRLink.ImageURL != nil {
		imageURL = sql.NullString{String: *currentQRLink.ImageURL, Valid: true}
	}

	var isTipEnabled, isActive sql.NullBool
	if updates.IsTipEnabled != nil {
		isTipEnabled = sql.NullBool{Bool: *updates.IsTipEnabled, Valid: true}
	} else {
		isTipEnabled = sql.NullBool{Bool: currentQRLink.IsTipEnabled, Valid: true}
	}

	if updates.IsActive != nil {
		isActive = sql.NullBool{Bool: *updates.IsActive, Valid: true}
	} else {
		isActive = sql.NullBool{Bool: currentQRLink.IsActive, Valid: true}
	}

	params := db.UpdateQRLinkParams{
		ID:               id,
		Amount:           amount,
		SupportedMethods: supportedMethodsJSON,
		Tag:              tag,
		Title:            title,
		Description:      description,
		ImageUrl:         imageURL,
		IsTipEnabled:     isTipEnabled,
		IsActive:         isActive,
		UserID:           userID,
	}

	dbQRLink, err := r.Queries.UpdateQRLink(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("QR link not found or not authorized")
		}
		return nil, fmt.Errorf("failed to update QR link: %w", err)
	}

	return r.mapDBQRLinkToEntity(dbQRLink)
}

func (r *QRRepositoryImpl) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	err := r.Queries.DeleteQRLink(ctx, db.DeleteQRLinkParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete QR link: %w", err)
	}
	return nil
}

func (r *QRRepositoryImpl) mapDBQRLinkToEntity(dbQRLink db.QrLink) (*entity.QRLink, error) {
	var supportedMethods []txEntity.TransactionMedium
	if err := json.Unmarshal(dbQRLink.SupportedMethods, &supportedMethods); err != nil {
		return nil, fmt.Errorf("failed to unmarshal supported methods: %w", err)
	}

	var amount *float64
	if dbQRLink.Amount.Valid {
		if amountFloat, err := strconv.ParseFloat(dbQRLink.Amount.String, 64); err == nil {
			amount = &amountFloat
		}
	}

	var title, description, imageURL *string
	if dbQRLink.Title.Valid {
		title = &dbQRLink.Title.String
	}
	if dbQRLink.Description.Valid {
		description = &dbQRLink.Description.String
	}
	if dbQRLink.ImageUrl.Valid {
		imageURL = &dbQRLink.ImageUrl.String
	}

	return &entity.QRLink{
		ID:               dbQRLink.ID,
		UserID:           dbQRLink.UserID,
		MerchantID:       dbQRLink.MerchantID,
		Type:             entity.QRLinkType(dbQRLink.Type),
		Amount:           amount,
		SupportedMethods: supportedMethods,
		Tag:              entity.QRLinkTag(dbQRLink.Tag),
		Title:            title,
		Description:      description,
		ImageURL:         imageURL,
		IsTipEnabled:     dbQRLink.IsTipEnabled.Bool,
		IsActive:         dbQRLink.IsActive.Bool,
		CreatedAt:        dbQRLink.CreatedAt,
		UpdatedAt:        dbQRLink.UpdatedAt,
	}, nil
}
