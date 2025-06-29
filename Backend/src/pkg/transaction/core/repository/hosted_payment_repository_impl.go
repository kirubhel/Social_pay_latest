package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	db "github.com/socialpay/socialpay/src/pkg/transaction/core/repository/generated"
)

type HostedPaymentRepositoryImpl struct {
	Queries *db.Queries
}

func NewHostedPaymentRepository(dbConn *sql.DB) HostedPaymentRepository {
	return &HostedPaymentRepositoryImpl{
		Queries: db.New(dbConn),
	}
}

func (r *HostedPaymentRepositoryImpl) Create(ctx context.Context, hostedPayment *entity.HostedPayment) error {
	// Convert supported mediums to JSON
	supportedMediums, err := json.Marshal(hostedPayment.SupportedMediums)
	if err != nil {
		return err
	}

	// Convert amount to string (as expected by SQLC)
	amountStr := decimal.NewFromFloat(hostedPayment.Amount).String()

	_, err = r.Queries.CreateHostedPayment(ctx, db.CreateHostedPaymentParams{
		ID:               hostedPayment.ID,
		UserID:           hostedPayment.UserID,
		MerchantID:       hostedPayment.MerchantID,
		Amount:           amountStr,
		Currency:         hostedPayment.Currency,
		Description:      sql.NullString{String: hostedPayment.Description, Valid: hostedPayment.Description != ""},
		Reference:        hostedPayment.Reference,
		SupportedMediums: supportedMediums,
		PhoneNumber:      sql.NullString{String: hostedPayment.PhoneNumber, Valid: hostedPayment.PhoneNumber != ""},
		SuccessUrl:       hostedPayment.SuccessURL,
		FailedUrl:        hostedPayment.FailedURL,
		CallbackUrl:      sql.NullString{String: hostedPayment.CallbackURL, Valid: hostedPayment.CallbackURL != ""},
	})

	return err
}

func (r *HostedPaymentRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.HostedPayment, error) {
	dbHostedPayment, err := r.Queries.GetHostedPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	return toEntityHostedPayment(&dbHostedPayment), nil
}

func (r *HostedPaymentRepositoryImpl) GetByReference(ctx context.Context, reference string, merchantID uuid.UUID) (*entity.HostedPayment, error) {
	dbHostedPayment, err := r.Queries.GetHostedPaymentByReference(ctx, db.GetHostedPaymentByReferenceParams{
		Reference:  reference,
		MerchantID: merchantID,
	})
	if err != nil {
		return nil, err
	}

	return toEntityHostedPayment(&dbHostedPayment), nil
}

func (r *HostedPaymentRepositoryImpl) ValidateReferenceId(ctx context.Context, merchantID uuid.UUID, reference string) error {
	// Try to get an existing hosted payment with the same reference and merchant ID
	_, err := r.Queries.GetHostedPaymentByReference(ctx, db.GetHostedPaymentByReferenceParams{
		Reference:  reference,
		MerchantID: merchantID,
	})

	// If no error, it means a hosted payment with this reference already exists
	if err == nil {
		return fmt.Errorf("reference '%s' already exists. Please use a unique reference", reference)
	}

	// If the error is "no rows found", then the reference is unique (which is what we want)
	if err == sql.ErrNoRows {
		return nil
	}

	// For any other error, return it as is
	return fmt.Errorf("failed to validate reference: %w", err)
}

func (r *HostedPaymentRepositoryImpl) UpdateWithTransaction(ctx context.Context, id uuid.UUID, transactionID uuid.UUID, selectedMedium string, selectedPhoneNumber string, status entity.HostedPaymentStatus) error {
	return r.Queries.UpdateHostedPaymentWithTransaction(ctx, db.UpdateHostedPaymentWithTransactionParams{
		ID:                  id,
		TransactionID:       uuid.NullUUID{UUID: transactionID, Valid: true},
		SelectedMedium:      sql.NullString{String: selectedMedium, Valid: true},
		SelectedPhoneNumber: sql.NullString{String: selectedPhoneNumber, Valid: true},
		Status:              db.HostedPaymentStatus(status),
	})
}

func (r *HostedPaymentRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.HostedPaymentStatus) error {
	return r.Queries.UpdateHostedPaymentStatus(ctx, db.UpdateHostedPaymentStatusParams{
		ID:     id,
		Status: db.HostedPaymentStatus(status),
	})
}

func (r *HostedPaymentRepositoryImpl) GetExpiredPayments(ctx context.Context) ([]entity.HostedPayment, error) {
	dbHostedPayments, err := r.Queries.GetExpiredHostedPayments(ctx)
	if err != nil {
		return nil, err
	}

	var hostedPayments []entity.HostedPayment
	for _, dbHostedPayment := range dbHostedPayments {
		hostedPayments = append(hostedPayments, *toEntityHostedPayment(&dbHostedPayment))
	}

	return hostedPayments, nil
}

func (r *HostedPaymentRepositoryImpl) Update(ctx context.Context, hostedPayment *entity.HostedPayment) error {
	// Convert supported mediums to JSON
	supportedMediums, err := json.Marshal(hostedPayment.SupportedMediums)
	if err != nil {
		return err
	}

	// Convert amount to string (as expected by SQLC)
	amountStr := decimal.NewFromFloat(hostedPayment.Amount).String()

	err = r.Queries.UpdateHostedPayment(ctx, db.UpdateHostedPaymentParams{
		ID:               hostedPayment.ID,
		Amount:           amountStr,
		Currency:         hostedPayment.Currency,
		Description:      sql.NullString{String: hostedPayment.Description, Valid: hostedPayment.Description != ""},
		SupportedMediums: supportedMediums,
		PhoneNumber:      sql.NullString{String: hostedPayment.PhoneNumber, Valid: hostedPayment.PhoneNumber != ""},
		SuccessUrl:       hostedPayment.SuccessURL,
		FailedUrl:        hostedPayment.FailedURL,
		CallbackUrl:      sql.NullString{String: hostedPayment.CallbackURL, Valid: hostedPayment.CallbackURL != ""},
		ExpiresAt:        hostedPayment.ExpiresAt,
	})

	return err
}

// Helper function to convert db.HostedPayment to entity.HostedPayment
func toEntityHostedPayment(dbHostedPayment *db.HostedPayment) *entity.HostedPayment {
	// Convert amount from string to float64
	amountDecimal, _ := decimal.NewFromString(dbHostedPayment.Amount)
	amount, _ := amountDecimal.Float64()

	// Parse supported mediums from JSON
	var supportedMediums []entity.TransactionMedium
	if len(dbHostedPayment.SupportedMediums) > 0 {
		json.Unmarshal(dbHostedPayment.SupportedMediums, &supportedMediums)
	}

	hostedPayment := &entity.HostedPayment{
		ID:               dbHostedPayment.ID,
		UserID:           dbHostedPayment.UserID,
		MerchantID:       dbHostedPayment.MerchantID,
		Amount:           amount,
		Currency:         dbHostedPayment.Currency,
		Description:      dbHostedPayment.Description.String,
		Reference:        dbHostedPayment.Reference,
		SupportedMediums: supportedMediums,
		PhoneNumber:      dbHostedPayment.PhoneNumber.String,
		SuccessURL:       dbHostedPayment.SuccessUrl,
		FailedURL:        dbHostedPayment.FailedUrl,
		CallbackURL:      dbHostedPayment.CallbackUrl.String,
		Status:           entity.HostedPaymentStatus(dbHostedPayment.Status),
		CreatedAt:        dbHostedPayment.CreatedAt,
		UpdatedAt:        dbHostedPayment.UpdatedAt,
		ExpiresAt:        dbHostedPayment.ExpiresAt,
	}

	// Handle optional transaction ID
	if dbHostedPayment.TransactionID.Valid {
		hostedPayment.TransactionID = &dbHostedPayment.TransactionID.UUID
	}

	// Handle selected payment details
	hostedPayment.SelectedMedium = dbHostedPayment.SelectedMedium.String
	hostedPayment.SelectedPhoneNumber = dbHostedPayment.SelectedPhoneNumber.String

	return hostedPayment
}
