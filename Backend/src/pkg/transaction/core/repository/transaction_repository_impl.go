package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/errorxx"
	"github.com/socialpay/socialpay/src/pkg/shared/filter"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	db "github.com/socialpay/socialpay/src/pkg/transaction/core/repository/generated"
	merchantEntity "github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
	"github.com/shopspring/decimal"
	"github.com/sqlc-dev/pqtype"
)

type TransactionRepositoryImpl struct {
	Queries *db.Queries
	q       sql.DB
}

func NewTransactionRepository(dbConn *sql.DB) TransactionRepository {
	return &TransactionRepositoryImpl{
		Queries: db.New(dbConn),
		q:       *dbConn,
	}
}

func (r *TransactionRepositoryImpl) OverrideTransactionStatus(ctx context.Context, txnID uuid.UUID, newStatus entity.TransactionStatus, reason string, adminID string) error {
	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return err
	}
	return r.Queries.OverrideTransactionStatus(ctx, db.OverrideTransactionStatusParams{
		TransactionID: txnID,
		Status:        db.TransactionStatus(newStatus),
		Reason:        reason,
		AdminID:       adminUUID,
	})
}

func (r *TransactionRepositoryImpl) GetTransactions(ctx context.Context, user_id uuid.UUID, limit, offset int32) ([]entity.Transaction, int, error) {
	dbTxns, err := r.Queries.GetTransactions(ctx, db.GetTransactionsParams{
		UserID: user_id,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := r.Queries.CountTransactions(ctx, user_id)

	if err != nil {

		return nil, 0, err
	}

	return toEntityTransactions(dbTxns), int(count), nil
}

func (r *TransactionRepositoryImpl) GetTransactionByParamenters(ctx context.Context,
	parameters *entity.FilterParameters) ([]entity.Transaction, error) {

	params := db.GetFilteredTransactionsParams{
		CreatedAt:   parameters.StartDate,
		CreatedAt_2: parameters.EndDate,
		Status:      db.TransactionStatus(parameters.Status),
		Type:        string(parameters.Type),
	}
	dbTxn, err := r.Queries.GetFilteredTransactions(ctx, params)
	if err != nil {
		return nil, err
	}
	return toEntityTransactions(dbTxn), nil
}

func (r *TransactionRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	dbTxn, err := r.Queries.GetTransaction(ctx, id)
	if err != nil {
		return nil, err
	}
	entityTxn := toEntityTransaction(&dbTxn)
	return &entityTxn, nil
}

func (r *TransactionRepositoryImpl) GetByIDWithMerchant(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	dbTxnWithMerchant, err := r.Queries.GetTransactionWithMerchant(ctx, id)
	if err != nil {
		return nil, err
	}
	entityTxn := toEntityTransactionWithMerchant(&dbTxnWithMerchant)
	return &entityTxn, nil
}

func (r *TransactionRepositoryImpl) GetByReferenceID(ctx context.Context, referenceID string) (*entity.Transaction, error) {
	dbTxn, err := r.Queries.GetByReferenceID(ctx, sql.NullString{
		String: referenceID,
		Valid:  true,
	})
	if err != nil {
		return nil, err
	}
	entityTxn := toEntityTransaction(&dbTxn)
	return &entityTxn, nil
}

func (r *TransactionRepositoryImpl) GetByUserIdAndReferenceID(ctx context.Context, userID uuid.UUID, referenceID string) (*entity.Transaction, error) {
	dbTxn, err := r.Queries.GetByUserIdAndReferenceID(ctx, db.GetByUserIdAndReferenceIDParams{
		UserID: userID,
		Reference: sql.NullString{
			String: referenceID,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, err
	}
	entityTxn := toEntityTransaction(&dbTxn)
	return &entityTxn, nil
}

func (r *TransactionRepositoryImpl) GetByMerchantIdAndReferenceID(ctx context.Context, merchantID uuid.UUID, referenceID string) (*entity.Transaction, error) {
	dbTxn, err := r.Queries.GetByMerchantIdAndReferenceID(ctx, db.GetByMerchantIdAndReferenceIDParams{
		MerchantID: uuid.NullUUID{
			UUID:  merchantID,
			Valid: true,
		},
		Reference: sql.NullString{
			String: referenceID,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, err
	}
	entityTxn := toEntityTransaction(&dbTxn)
	return &entityTxn, nil
}

func (r *TransactionRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TransactionStatus) error {
	return r.Queries.UpdateStatus(ctx, db.UpdateStatusParams{
		ID:     id,
		Status: db.TransactionStatus(status),
	})
}

func toEntityTransaction(dbTxn *db.Transaction) entity.Transaction {
	tx := entity.Transaction{
		Id:        dbTxn.ID,
		UserId:    dbTxn.UserID,
		Type:      entity.TransactionType(dbTxn.Type),
		Medium:    entity.TransactionMedium(dbTxn.Medium),
		CreatedAt: dbTxn.CreatedAt,
		UpdatedAt: dbTxn.UpdatedAt,
		Status:    entity.TransactionStatus(dbTxn.Status),
	}

	// Handle nullable/conversion fields
	if dbTxn.PhoneNumber.Valid {
		tx.PhoneNumber = dbTxn.PhoneNumber.String
	}

	if dbTxn.MerchantID.Valid {
		tx.MerchantId = dbTxn.MerchantID.UUID
	}

	if dbTxn.Verified.Valid {
		tx.Verified = dbTxn.Verified.Bool
	}

	if dbTxn.Ttl.Valid {
		tx.TTL = dbTxn.Ttl.Int64
	}

	if dbTxn.ConfirmTimestamp.Valid {
		tx.Confirm_Timestamp = dbTxn.ConfirmTimestamp.Time
	}

	if dbTxn.Test.Valid {
		tx.Test = dbTxn.Test.Bool
	}

	// Convert decimal to float64
	amount, _ := dbTxn.Amount.Float64()
	tx.Amount = amount

	if dbTxn.HasChallenge.Valid {
		tx.HasChallenge = dbTxn.HasChallenge.Bool
	}

	if dbTxn.WebhookReceived.Valid {
		tx.WebhookReceived = dbTxn.WebhookReceived.Bool
	}

	if dbTxn.FeeAmount.Valid {
		amount, _ := dbTxn.FeeAmount.Decimal.Float64()
		tx.FeeAmount = amount
	}

	if dbTxn.AdminNet.Valid {
		amount, _ := dbTxn.AdminNet.Decimal.Float64()
		tx.AdminNet = amount
	}

	if dbTxn.VatAmount.Valid {
		amount, _ := dbTxn.VatAmount.Decimal.Float64()
		tx.VatAmount = amount
	}

	if dbTxn.MerchantNet.Valid {
		amount, _ := dbTxn.MerchantNet.Decimal.Float64()
		tx.MerchantNet = amount
	}

	if dbTxn.TotalAmount.Valid {
		amount, _ := dbTxn.TotalAmount.Decimal.Float64()
		tx.TotalAmount = amount
	}

	if dbTxn.Currency.Valid {
		tx.Currency = dbTxn.Currency.String
	}

	// Handle other nullable string fields
	if dbTxn.Reference.Valid {
		tx.Reference = dbTxn.Reference.String
	}
	if dbTxn.Comment.Valid {
		tx.Comment = dbTxn.Comment.String
	}
	if dbTxn.ReferenceNumber.Valid {
		tx.ReferenceNumber = dbTxn.ReferenceNumber.String
	}
	if dbTxn.Description.Valid {
		tx.Description = dbTxn.Description.String
	}
	if dbTxn.Token.Valid {
		tx.Token = dbTxn.Token.String
	}
	if dbTxn.CallbackUrl.Valid {
		tx.CallbackURL = dbTxn.CallbackUrl.String
	}
	if dbTxn.SuccessUrl.Valid {
		tx.SuccessURL = dbTxn.SuccessUrl.String
	}
	if dbTxn.FailedUrl.Valid {
		tx.FailedURL = dbTxn.FailedUrl.String
	}

	// Handle QR and tip context fields
	if dbTxn.TransactionSource.Valid {
		tx.TransactionSource = entity.TransactionSource(dbTxn.TransactionSource.TransactionSource)
	}

	if dbTxn.QrLinkID.Valid {
		tx.QRLinkID = &dbTxn.QrLinkID.UUID
	}
	if dbTxn.HostedCheckoutID.Valid {
		tx.HostedCheckoutID = &dbTxn.HostedCheckoutID.UUID
	}
	if dbTxn.QrTag.Valid {
		tx.QRTag = &dbTxn.QrTag.String
	}

	if dbTxn.HasTip.Valid {
		tx.HasTip = dbTxn.HasTip.Bool
	}

	if dbTxn.TipAmount.Valid {
		amount := 0.0
		fmt.Sscanf(dbTxn.TipAmount.String, "%f", &amount)
		tx.TipAmount = &amount
	}
	if dbTxn.TipeePhone.Valid {
		tx.TipeePhone = &dbTxn.TipeePhone.String
	}
	if dbTxn.TipMedium.Valid {
		tx.TipMedium = &dbTxn.TipMedium.String
	}
	if dbTxn.TipTransactionID.Valid {
		tx.TipTransactionID = &dbTxn.TipTransactionID.UUID
	}
	if dbTxn.TipProcessed.Valid {
		tx.TipProcessed = dbTxn.TipProcessed.Bool
	}

	// Handle details JSON
	if dbTxn.Details.Valid {
		json.Unmarshal(dbTxn.Details.RawMessage, &tx.Details)
	}

	return tx
}

func toEntityTransactions(dbTxns []db.Transaction) []entity.Transaction {
	entities := make([]entity.Transaction, len(dbTxns))
	for i, dbTxn := range dbTxns {
		entities[i] = toEntityTransaction(&dbTxn)
	}
	return entities
}

func toEntityTransactionWithMerchant(dbTxnWithMerchant *db.GetTransactionWithMerchantRow) entity.Transaction {
	tx := entity.Transaction{
		Id:        dbTxnWithMerchant.ID,
		UserId:    dbTxnWithMerchant.UserID,
		Type:      entity.TransactionType(dbTxnWithMerchant.Type),
		Medium:    entity.TransactionMedium(dbTxnWithMerchant.Medium),
		CreatedAt: dbTxnWithMerchant.CreatedAt,
		UpdatedAt: dbTxnWithMerchant.UpdatedAt,
		Status:    entity.TransactionStatus(dbTxnWithMerchant.Status),
	}

	// Handle nullable/conversion fields
	if dbTxnWithMerchant.PhoneNumber.Valid {
		tx.PhoneNumber = dbTxnWithMerchant.PhoneNumber.String
	}

	if dbTxnWithMerchant.MerchantID.Valid {
		tx.MerchantId = dbTxnWithMerchant.MerchantID.UUID
	}

	if dbTxnWithMerchant.Verified.Valid {
		tx.Verified = dbTxnWithMerchant.Verified.Bool
	}

	if dbTxnWithMerchant.Ttl.Valid {
		tx.TTL = dbTxnWithMerchant.Ttl.Int64
	}

	if dbTxnWithMerchant.ConfirmTimestamp.Valid {
		tx.Confirm_Timestamp = dbTxnWithMerchant.ConfirmTimestamp.Time
	}

	if dbTxnWithMerchant.Test.Valid {
		tx.Test = dbTxnWithMerchant.Test.Bool
	}

	// Convert decimal to float64
	amount, _ := dbTxnWithMerchant.Amount.Float64()
	tx.Amount = amount

	if dbTxnWithMerchant.HasChallenge.Valid {
		tx.HasChallenge = dbTxnWithMerchant.HasChallenge.Bool
	}

	if dbTxnWithMerchant.WebhookReceived.Valid {
		tx.WebhookReceived = dbTxnWithMerchant.WebhookReceived.Bool
	}

	if dbTxnWithMerchant.FeeAmount.Valid {
		amount, _ := dbTxnWithMerchant.FeeAmount.Decimal.Float64()
		tx.FeeAmount = amount
	}

	if dbTxnWithMerchant.AdminNet.Valid {
		amount, _ := dbTxnWithMerchant.AdminNet.Decimal.Float64()
		tx.AdminNet = amount
	}

	if dbTxnWithMerchant.VatAmount.Valid {
		amount, _ := dbTxnWithMerchant.VatAmount.Decimal.Float64()
		tx.VatAmount = amount
	}

	if dbTxnWithMerchant.MerchantNet.Valid {
		amount, _ := dbTxnWithMerchant.MerchantNet.Decimal.Float64()
		tx.MerchantNet = amount
	}

	if dbTxnWithMerchant.TotalAmount.Valid {
		amount, _ := dbTxnWithMerchant.TotalAmount.Decimal.Float64()
		tx.TotalAmount = amount
	}

	if dbTxnWithMerchant.Currency.Valid {
		tx.Currency = dbTxnWithMerchant.Currency.String
	}

	// Handle other nullable string fields
	if dbTxnWithMerchant.Reference.Valid {
		tx.Reference = dbTxnWithMerchant.Reference.String
	}
	if dbTxnWithMerchant.Comment.Valid {
		tx.Comment = dbTxnWithMerchant.Comment.String
	}
	if dbTxnWithMerchant.ReferenceNumber.Valid {
		tx.ReferenceNumber = dbTxnWithMerchant.ReferenceNumber.String
	}
	if dbTxnWithMerchant.Description.Valid {
		tx.Description = dbTxnWithMerchant.Description.String
	}
	if dbTxnWithMerchant.Token.Valid {
		tx.Token = dbTxnWithMerchant.Token.String
	}
	if dbTxnWithMerchant.CallbackUrl.Valid {
		tx.CallbackURL = dbTxnWithMerchant.CallbackUrl.String
	}
	if dbTxnWithMerchant.SuccessUrl.Valid {
		tx.SuccessURL = dbTxnWithMerchant.SuccessUrl.String
	}
	if dbTxnWithMerchant.FailedUrl.Valid {
		tx.FailedURL = dbTxnWithMerchant.FailedUrl.String
	}

	// Handle QR and tip context fields
	if dbTxnWithMerchant.TransactionSource.Valid {
		tx.TransactionSource = entity.TransactionSource(dbTxnWithMerchant.TransactionSource.TransactionSource)
	}

	if dbTxnWithMerchant.QrLinkID.Valid {
		tx.QRLinkID = &dbTxnWithMerchant.QrLinkID.UUID
	}
	if dbTxnWithMerchant.HostedCheckoutID.Valid {
		tx.HostedCheckoutID = &dbTxnWithMerchant.HostedCheckoutID.UUID
	}
	if dbTxnWithMerchant.QrTag.Valid {
		tx.QRTag = &dbTxnWithMerchant.QrTag.String
	}

	if dbTxnWithMerchant.HasTip.Valid {
		tx.HasTip = dbTxnWithMerchant.HasTip.Bool
	}

	if dbTxnWithMerchant.TipAmount.Valid {
		amount := 0.0
		fmt.Sscanf(dbTxnWithMerchant.TipAmount.String, "%f", &amount)
		tx.TipAmount = &amount
	}
	if dbTxnWithMerchant.TipeePhone.Valid {
		tx.TipeePhone = &dbTxnWithMerchant.TipeePhone.String
	}
	if dbTxnWithMerchant.TipMedium.Valid {
		tx.TipMedium = &dbTxnWithMerchant.TipMedium.String
	}
	if dbTxnWithMerchant.TipTransactionID.Valid {
		tx.TipTransactionID = &dbTxnWithMerchant.TipTransactionID.UUID
	}
	if dbTxnWithMerchant.TipProcessed.Valid {
		tx.TipProcessed = dbTxnWithMerchant.TipProcessed.Bool
	}

	// Handle details JSON
	if dbTxnWithMerchant.Details.Valid {
		json.Unmarshal(dbTxnWithMerchant.Details.RawMessage, &tx.Details)
	}

	// Handle merchant data with proper null checking
	var merchantID uuid.UUID
	if dbTxnWithMerchant.MerchantID.Valid {
		merchantID = dbTxnWithMerchant.MerchantID.UUID
	}

	var legalName string
	if dbTxnWithMerchant.MerchantLegalName.Valid {
		legalName = dbTxnWithMerchant.MerchantLegalName.String
	}

	var tradingName *string
	if dbTxnWithMerchant.MerchantTradingName.Valid {
		tradingName = &dbTxnWithMerchant.MerchantTradingName.String
	}

	var businessRegNum string
	if dbTxnWithMerchant.MerchantBusinessRegistrationNumber.Valid {
		businessRegNum = dbTxnWithMerchant.MerchantBusinessRegistrationNumber.String
	}

	var taxIdNum string
	if dbTxnWithMerchant.MerchantTaxIdentificationNumber.Valid {
		taxIdNum = dbTxnWithMerchant.MerchantTaxIdentificationNumber.String
	}

	var businessType string
	if dbTxnWithMerchant.MerchantBusinessType.Valid {
		businessType = dbTxnWithMerchant.MerchantBusinessType.String
	}

	var industryCategory *string
	if dbTxnWithMerchant.MerchantIndustryCategory.Valid {
		industryCategory = &dbTxnWithMerchant.MerchantIndustryCategory.String
	}

	var isBettingCompany bool
	if dbTxnWithMerchant.MerchantIsBettingCompany.Valid {
		isBettingCompany = dbTxnWithMerchant.MerchantIsBettingCompany.Bool
	}

	var lotteryCertNum *string
	if dbTxnWithMerchant.MerchantLotteryCertificateNumber.Valid {
		lotteryCertNum = &dbTxnWithMerchant.MerchantLotteryCertificateNumber.String
	}

	var websiteURL *string
	if dbTxnWithMerchant.MerchantWebsiteUrl.Valid {
		websiteURL = &dbTxnWithMerchant.MerchantWebsiteUrl.String
	}

	var establishedDate *time.Time
	if dbTxnWithMerchant.MerchantEstablishedDate.Valid {
		establishedDate = &dbTxnWithMerchant.MerchantEstablishedDate.Time
	}

	var createdAt time.Time
	if dbTxnWithMerchant.MerchantCreatedAt.Valid {
		createdAt = dbTxnWithMerchant.MerchantCreatedAt.Time
	}

	var updatedAt time.Time
	if dbTxnWithMerchant.MerchantUpdatedAt.Valid {
		updatedAt = dbTxnWithMerchant.MerchantUpdatedAt.Time
	}

	var status merchantEntity.MerchantStatus
	if dbTxnWithMerchant.MerchantStatus.Valid {
		status = merchantEntity.MerchantStatus(dbTxnWithMerchant.MerchantStatus.String)
	}

	tx.Merchant = &merchantEntity.Merchant{
		ID:                         merchantID,
		UserID:                     dbTxnWithMerchant.UserID,
		LegalName:                  legalName,
		TradingName:                tradingName,
		BusinessRegistrationNumber: businessRegNum,
		TaxIdentificationNumber:    taxIdNum,
		BusinessType:               businessType,
		IndustryCategory:           industryCategory,
		IsBettingCompany:           isBettingCompany,
		LotteryCertificateNumber:   lotteryCertNum,
		WebsiteURL:                 websiteURL,
		EstablishedDate:            establishedDate,
		CreatedAt:                  createdAt,
		UpdatedAt:                  updatedAt,
		Status:                     status,
	}
	return tx
}

func (r *TransactionRepositoryImpl) Create(ctx context.Context, tx *entity.Transaction) error {
	params := r.toCreateParams(tx)
	return r.Queries.CreateTransaction(ctx, params)
}

// Helper function to convert entity Transaction to generated Transaction params
func (r *TransactionRepositoryImpl) toCreateParams(tx *entity.Transaction) db.CreateTransactionParams {
	params := db.CreateTransactionParams{
		ID:     tx.Id,
		UserID: tx.UserId,
		Type:   string(tx.Type),
		Medium: string(tx.Medium),
		Status: db.TransactionStatus(tx.Status),
		Amount: decimal.NewFromFloat(tx.Amount),
	}

	// Handle nullable fields
	if tx.PhoneNumber != "" {
		params.PhoneNumber = sql.NullString{String: tx.PhoneNumber, Valid: true}
	}
	if tx.MerchantId != uuid.Nil {
		params.MerchantID = uuid.NullUUID{UUID: tx.MerchantId, Valid: true}
	}
	if tx.Reference != "" {
		params.Reference = sql.NullString{String: tx.Reference, Valid: true}
	}
	if tx.Comment != "" {
		params.Comment = sql.NullString{String: tx.Comment, Valid: true}
	}
	if tx.ReferenceNumber != "" {
		params.ReferenceNumber = sql.NullString{String: tx.ReferenceNumber, Valid: true}
	}
	if tx.Description != "" {
		params.Description = sql.NullString{String: tx.Description, Valid: true}
	}
	if tx.Token != "" {
		params.Token = sql.NullString{String: tx.Token, Valid: true}
	}
	if tx.CallbackURL != "" {
		params.CallbackUrl = sql.NullString{String: tx.CallbackURL, Valid: true}
	}
	if tx.SuccessURL != "" {
		params.SuccessUrl = sql.NullString{String: tx.SuccessURL, Valid: true}
	}
	if tx.FailedURL != "" {
		params.FailedUrl = sql.NullString{String: tx.FailedURL, Valid: true}
	}
	if tx.Currency != "" {
		params.Currency = sql.NullString{String: tx.Currency, Valid: true}
	}

	// Handle decimal fields properly
	if tx.FeeAmount != 0 {
		params.FeeAmount = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.FeeAmount), Valid: true}
	}
	if tx.AdminNet != 0 {
		params.AdminNet = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.AdminNet), Valid: true}
	}
	if tx.VatAmount != 0 {
		params.VatAmount = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.VatAmount), Valid: true}
	}
	if tx.MerchantNet != 0 {
		params.MerchantNet = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.MerchantNet), Valid: true}
	}
	if tx.TotalAmount != 0 {
		params.TotalAmount = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.TotalAmount), Valid: true}
	}

	// Handle boolean fields
	params.Verified = sql.NullBool{Bool: tx.Verified, Valid: true}
	params.Test = sql.NullBool{Bool: tx.Test, Valid: true}
	params.HasChallenge = sql.NullBool{Bool: tx.HasChallenge, Valid: true}

	// Handle TTL
	if tx.TTL != 0 {
		params.Ttl = sql.NullInt64{Int64: tx.TTL, Valid: true}
	}

	// Handle details JSON
	if tx.Details != nil {
		detailsJSON, _ := json.Marshal(tx.Details)
		params.Details = pqtype.NullRawMessage{RawMessage: detailsJSON, Valid: true}
	}

	// Handle timestamp fields
	if !tx.Confirm_Timestamp.IsZero() {
		params.ConfirmTimestamp = sql.NullTime{Time: tx.Confirm_Timestamp, Valid: true}
	}

	// Handle context-specific fields like QR and tips
	if tx.TransactionSource != "" {
		params.TransactionSource = db.NullTransactionSource{
			TransactionSource: db.TransactionSource(tx.TransactionSource),
			Valid:             true,
		}
	}

	if tx.QRLinkID != nil {
		params.QrLinkID = uuid.NullUUID{UUID: *tx.QRLinkID, Valid: true}
	}
	if tx.HostedCheckoutID != nil {
		params.HostedCheckoutID = uuid.NullUUID{UUID: *tx.HostedCheckoutID, Valid: true}
	}
	if tx.QRTag != nil {
		params.QrTag = sql.NullString{String: *tx.QRTag, Valid: true}
	}

	params.HasTip = sql.NullBool{Bool: tx.HasTip, Valid: true}

	if tx.TipAmount != nil {
		params.TipAmount = sql.NullString{String: fmt.Sprintf("%f", *tx.TipAmount), Valid: true}
	}
	if tx.TipeePhone != nil {
		params.TipeePhone = sql.NullString{String: *tx.TipeePhone, Valid: true}
	}
	if tx.TipMedium != nil {
		params.TipMedium = sql.NullString{String: *tx.TipMedium, Valid: true}
	}

	return params
}

func (r *TransactionRepositoryImpl) toCreateWithContextParams(tx *entity.Transaction) db.CreateTransactionWithContextParams {
	params := db.CreateTransactionWithContextParams{
		ID:     tx.Id,
		UserID: tx.UserId,
		Type:   string(tx.Type),
		Medium: string(tx.Medium),
		Status: db.TransactionStatus(tx.Status),
		Amount: decimal.NewFromFloat(tx.Amount),
	}

	// Handle all the same nullable fields as toCreateParams
	if tx.PhoneNumber != "" {
		params.PhoneNumber = sql.NullString{String: tx.PhoneNumber, Valid: true}
	}
	if tx.MerchantId != uuid.Nil {
		params.MerchantID = uuid.NullUUID{UUID: tx.MerchantId, Valid: true}
	}
	if tx.Reference != "" {
		params.Reference = sql.NullString{String: tx.Reference, Valid: true}
	}
	if tx.Comment != "" {
		params.Comment = sql.NullString{String: tx.Comment, Valid: true}
	}
	if tx.ReferenceNumber != "" {
		params.ReferenceNumber = sql.NullString{String: tx.ReferenceNumber, Valid: true}
	}
	if tx.Description != "" {
		params.Description = sql.NullString{String: tx.Description, Valid: true}
	}
	if tx.Token != "" {
		params.Token = sql.NullString{String: tx.Token, Valid: true}
	}
	if tx.CallbackURL != "" {
		params.CallbackUrl = sql.NullString{String: tx.CallbackURL, Valid: true}
	}
	if tx.SuccessURL != "" {
		params.SuccessUrl = sql.NullString{String: tx.SuccessURL, Valid: true}
	}
	if tx.FailedURL != "" {
		params.FailedUrl = sql.NullString{String: tx.FailedURL, Valid: true}
	}
	if tx.Currency != "" {
		params.Currency = sql.NullString{String: tx.Currency, Valid: true}
	}

	// Handle decimal fields
	if tx.FeeAmount != 0 {
		params.FeeAmount = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.FeeAmount), Valid: true}
	}
	if tx.AdminNet != 0 {
		params.AdminNet = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.AdminNet), Valid: true}
	}
	if tx.VatAmount != 0 {
		params.VatAmount = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.VatAmount), Valid: true}
	}
	if tx.MerchantNet != 0 {
		params.MerchantNet = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.MerchantNet), Valid: true}
	}
	if tx.TotalAmount != 0 {
		params.TotalAmount = decimal.NullDecimal{Decimal: decimal.NewFromFloat(tx.TotalAmount), Valid: true}
	}

	// Handle boolean fields
	params.Verified = sql.NullBool{Bool: tx.Verified, Valid: true}
	params.Test = sql.NullBool{Bool: tx.Test, Valid: true}
	params.HasChallenge = sql.NullBool{Bool: tx.HasChallenge, Valid: true}

	// Handle TTL
	if tx.TTL != 0 {
		params.Ttl = sql.NullInt64{Int64: tx.TTL, Valid: true}
	}

	// Handle details JSON
	if tx.Details != nil {
		detailsJSON, _ := json.Marshal(tx.Details)
		params.Details = pqtype.NullRawMessage{RawMessage: detailsJSON, Valid: true}
	}

	// Handle timestamp fields
	if !tx.Confirm_Timestamp.IsZero() {
		params.ConfirmTimestamp = sql.NullTime{Time: tx.Confirm_Timestamp, Valid: true}
	}

	// Handle context-specific fields like QR and tips
	if tx.TransactionSource != "" {
		params.TransactionSource = db.NullTransactionSource{
			TransactionSource: db.TransactionSource(tx.TransactionSource),
			Valid:             true,
		}
	}

	if tx.QRLinkID != nil {
		params.QrLinkID = uuid.NullUUID{UUID: *tx.QRLinkID, Valid: true}
	}
	if tx.HostedCheckoutID != nil {
		params.HostedCheckoutID = uuid.NullUUID{UUID: *tx.HostedCheckoutID, Valid: true}
	}
	if tx.QRTag != nil {
		params.QrTag = sql.NullString{String: *tx.QRTag, Valid: true}
	}

	params.HasTip = sql.NullBool{Bool: tx.HasTip, Valid: true}

	if tx.TipAmount != nil {
		params.TipAmount = sql.NullString{String: fmt.Sprintf("%f", *tx.TipAmount), Valid: true}
	}
	if tx.TipeePhone != nil {
		params.TipeePhone = sql.NullString{String: *tx.TipeePhone, Valid: true}
	}
	if tx.TipMedium != nil {
		params.TipMedium = sql.NullString{String: *tx.TipMedium, Valid: true}
	}

	return params
}

func (r *TransactionRepositoryImpl) Update(ctx context.Context, tx *entity.Transaction) error {
	return r.Queries.UpdateTransaction(ctx, db.UpdateTransactionParams{
		ID:     tx.Id,
		Status: db.TransactionStatus(tx.Status),
	})
}

func (r *TransactionRepositoryImpl) GetTransactionsByParameters(ctx context.Context, filterParam filter.Filter, userID uuid.UUID) ([]entity.Transaction, error) {
	// status := string(params.Status)

	// txs, err := r.Queries.GetFilteredTransactions(ctx, db.GetFilteredTransactionsParams{
	// 	UserID:      userID,
	// 	CreatedAt:   params.StartDate,
	// 	CreatedAt_2: params.EndDate,
	// 	Column4:     status,
	// 	Limit:       limit,
	// 	Offset:      offset,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	clause, args, err := filterParam.Build()

	if err != nil {

		return nil, err
	}

	txs, err := r.GetTransactionWithParameter(clause, args)

	if err != nil {

		return nil, err
	}

	return toEntityTransactions(txs), nil
}

func (r *TransactionRepositoryImpl) GetTransactionByParametersCount(ctx context.Context,
	filterParam filter.Filter, userID uuid.UUID) (int, error) {

	clause, args, err := filterParam.Build()

	if err != nil {

		return 0, err
	}

	total, err := r.CountTransactionWithParameter(ctx, clause, args)

	if err != nil {

		return 0, err

	}

	return total, nil

}

func (r *TransactionRepositoryImpl) GetTransactionsByStatus(ctx context.Context, status entity.TransactionStatus, limit, offset int32) ([]entity.Transaction, error) {
	txs, err := r.Queries.GetTransactionsByStatus(ctx, db.GetTransactionsByStatusParams{
		Status: db.TransactionStatus(status),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toEntityTransactions(txs), nil
}

func (r *TransactionRepositoryImpl) GetTransactionsByType(ctx context.Context, txType entity.TransactionType, limit, offset int32) ([]entity.Transaction, error) {
	dbTxns, err := r.Queries.GetTransactionsByType(ctx, db.GetTransactionsByTypeParams{
		Type:   string(txType),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return toEntityTransactions(dbTxns), nil
}

func (r *TransactionRepositoryImpl) GetMerchantTransactions(ctx context.Context, merchantID uuid.UUID, limit, offset int32) ([]entity.Transaction, error) {
	dbTxns, err := r.Queries.GetMerchantTransactions(ctx, db.GetMerchantTransactionsParams{
		MerchantID: uuid.NullUUID{UUID: merchantID, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {

		err = errorxx.ErrDBRead.Wrap(err, "Getting transaction error").
			WithProperty(errorxx.ErrorCode, 500)

		return nil, err
	}

	return toEntityTransactions(dbTxns), nil
}

func (r *TransactionRepositoryImpl) GetFilteredMerchantTransactions(ctx context.Context, params *entity.FilterParameters, merchantID uuid.UUID, limit, offset int32) ([]entity.Transaction, error) {
	// Validate the parameters
	if err := params.Validate(); err != nil {
		return nil, err
	}

	// Convert parameters to strings, use empty string for null values
	status := string(params.Status)
	if params.Status == "" {
		status = ""
	}

	txType := string(params.Type)
	if params.Type == "" {
		txType = ""
	}

	txs, err := r.Queries.GetFilteredMerchantTransactions(ctx, db.GetFilteredMerchantTransactionsParams{
		MerchantID:  uuid.NullUUID{UUID: merchantID, Valid: true},
		CreatedAt:   params.StartDate,
		CreatedAt_2: params.EndDate,
		Status:      db.TransactionStatus(status),
		Type:        txType,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, err
	}
	return toEntityTransactions(txs), nil
}

func (r *TransactionRepositoryImpl) CreateWithContext(ctx context.Context, tx *entity.Transaction) error {
	params := r.toCreateWithContextParams(tx)
	return r.Queries.CreateTransactionWithContext(ctx, params)
}

func (r *TransactionRepositoryImpl) UpdateTipProcessing(ctx context.Context, transactionID, tipTransactionID uuid.UUID) error {
	return r.Queries.UpdateTipProcessing(ctx, db.UpdateTipProcessingParams{
		ID:               transactionID,
		TipTransactionID: uuid.NullUUID{UUID: tipTransactionID, Valid: true},
	})
}

func (r *TransactionRepositoryImpl) GetTransactionsWithPendingTips(ctx context.Context) ([]entity.Transaction, error) {
	dbTxns, err := r.Queries.GetTransactionsWithPendingTips(ctx)
	if err != nil {
		return nil, err
	}

	return toEntityTransactions(dbTxns), nil
}

func (r *TransactionRepositoryImpl) GetTransactionsByQRLink(ctx context.Context, qrLinkID uuid.UUID, limit, offset int32) ([]entity.Transaction, error) {
	dbTxns, err := r.Queries.GetTransactionsByQRLink(ctx, db.GetTransactionsByQRLinkParams{
		QrLinkID: uuid.NullUUID{UUID: qrLinkID, Valid: true},
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}

	return toEntityTransactions(dbTxns), nil
}

// GetTransactionAnalytics aggregates transaction data based on filters
func (r *TransactionRepositoryImpl) GetTransactionAnalytics(ctx context.Context, filter *entity.AnalyticsFilter, merchantID uuid.UUID) (*entity.TransactionAnalytics, error) {
	// Add query timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Build the base WHERE clause
	whereClause, args := r.buildAnalyticsWhereClause(filter, merchantID)

	fmt.Println("whereClause", whereClause)

	// Optimized analytics query using parallel aggregation
	// Split into separate queries for better performance with proper indexes

	// Main aggregation query
	mainQuery := `
		SELECT 
			COUNT(*) as total_transactions,
			COALESCE(SUM(amount), 0) as total_amount,
			COALESCE(SUM(merchant_net), 0) as total_merchant_net
		FROM transactions 
		WHERE ` + whereClause

	// Transaction type breakdown query (separate for better index usage)
	typeQuery := `
		SELECT 
			type,
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as amount
		FROM transactions 
		WHERE ` + whereClause + `
		GROUP BY type`

	// Tip aggregation query (separate for better performance)
	tipQuery := `
		SELECT 
			COUNT(*) as tip_count,
			COALESCE(SUM(COALESCE(tip_amount::numeric, 0)), 0) as tip_amount
		FROM transactions 
		WHERE ` + whereClause + ` AND has_tip = true`

	var analytics entity.TransactionAnalytics

	// Execute main query
	err := r.q.QueryRowContext(ctx, mainQuery, args...).Scan(
		&analytics.TotalTransactions,
		&analytics.TotalAmount,
		&analytics.TotalMerchantNet,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute main analytics query: %w", err)
	}

	// Execute type breakdown query
	typeRows, err := r.q.QueryContext(ctx, typeQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute type breakdown query: %w", err)
	}
	defer typeRows.Close()

	// Initialize type analytics
	analytics.TotalDeposits = entity.TransactionTypeAnalytics{Count: 0, Amount: 0}
	analytics.TotalWithdrawals = entity.TransactionTypeAnalytics{Count: 0, Amount: 0}

	for typeRows.Next() {
		var txType string
		var count int64
		var amount float64

		if err := typeRows.Scan(&txType, &count, &amount); err != nil {
			return nil, fmt.Errorf("failed to scan type breakdown: %w", err)
		}

		switch txType {
		case "DEPOSIT":
			analytics.TotalDeposits = entity.TransactionTypeAnalytics{Count: count, Amount: amount}
		case "WITHDRAWAL":
			analytics.TotalWithdrawals = entity.TransactionTypeAnalytics{Count: count, Amount: amount}
		}
	}

	// Execute tip query
	var tipCount int64
	var tipAmount float64
	err = r.q.QueryRowContext(ctx, tipQuery, args...).Scan(&tipCount, &tipAmount)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to execute tip query: %w", err)
	}
	analytics.TotalTips = entity.TransactionTypeAnalytics{Count: tipCount, Amount: tipAmount}

	// Calculate period comparison asynchronously if needed (optional optimization)
	// Skip period comparison for now to improve performance - can be added as separate endpoint
	// if comparison, err := r.calculatePeriodComparison(ctx, filter, userID, &analytics); err == nil {
	// 	analytics.PeriodComparison = comparison
	// }

	return &analytics, nil
}

// GetChartData generates chart data based on filters
func (r *TransactionRepositoryImpl) GetChartData(ctx context.Context, filter *entity.ChartFilter, merchantID uuid.UUID) (*entity.ChartData, error) {
	// Add query timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Build the base WHERE clause
	whereClause, args := r.buildAnalyticsWhereClause(&filter.AnalyticsFilter, merchantID)

	// Determine the date truncation based on date unit
	dateTrunc := r.getDateTruncation(filter.DateUnit)

	// Optimized chart query - single query structure for both count and amount
	query := fmt.Sprintf(`
		SELECT 
			DATE_TRUNC('%s', created_at) as period,
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as total_amount
		FROM transactions 
		WHERE %s
		GROUP BY DATE_TRUNC('%s', created_at)
		ORDER BY period ASC
		LIMIT 1000
	`, dateTrunc, whereClause, dateTrunc) // Add LIMIT to prevent excessive data

	rows, err := r.q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute chart query: %w", err)
	}
	defer rows.Close()

	var dataPoints []entity.ChartDataPoint
	var totalValue, maxValue, minValue float64
	var count int
	minValue = math.MaxFloat64

	for rows.Next() {
		var period time.Time
		var transactionCount int64
		var totalAmount float64

		err := rows.Scan(&period, &transactionCount, &totalAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chart data: %w", err)
		}

		var value float64
		if filter.ChartType == "count" {
			value = float64(transactionCount)
		} else {
			value = totalAmount
		}

		dataPoint := entity.ChartDataPoint{
			Date:  period,
			Value: value,
			Label: r.formatDateLabel(period, filter.DateUnit),
		}

		dataPoints = append(dataPoints, dataPoint)
		totalValue += value
		if value > maxValue {
			maxValue = value
		}
		if value < minValue {
			minValue = value
		}
		count++
	}

	if count == 0 {
		minValue = 0
	}

	averageValue := float64(0)
	if count > 0 {
		averageValue = totalValue / float64(count)
	}

	chartData := &entity.ChartData{
		ChartType: filter.ChartType,
		DateUnit:  filter.DateUnit,
		Data:      dataPoints,
		Summary: entity.ChartSummary{
			TotalValue:   totalValue,
			AverageValue: averageValue,
			MaxValue:     maxValue,
			MinValue:     minValue,
			DataPoints:   count,
		},
	}

	return chartData, nil
}

// Helper function to build WHERE clause for analytics
func (r *TransactionRepositoryImpl) buildAnalyticsWhereClause(filter *entity.AnalyticsFilter, merchantID uuid.UUID) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// User filter - most selective condition first
	conditions = append(conditions, fmt.Sprintf("merchant_id = $%d", argIndex))
	args = append(args, merchantID)
	argIndex++

	// Date range filter - second most selective
	conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
	args = append(args, filter.StartDate)
	argIndex++

	conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
	args = append(args, filter.EndDate)
	argIndex++

	// Status filter - use ANY for better performance with arrays
	if len(filter.Status) > 0 {
		conditions = append(conditions, fmt.Sprintf("status = ANY($%d)", argIndex))
		// Convert to string array for PostgreSQL
		statusStrings := make([]string, len(filter.Status))
		for i, status := range filter.Status {
			statusStrings[i] = string(status)
		}
		args = append(args, statusStrings)
		argIndex++
	}

	// Type filter - use ANY for better performance
	if len(filter.Type) > 0 {
		conditions = append(conditions, fmt.Sprintf("type = ANY($%d)", argIndex))
		typeStrings := make([]string, len(filter.Type))
		for i, txType := range filter.Type {
			typeStrings[i] = string(txType)
		}
		args = append(args, typeStrings)
		argIndex++
	}

	// Medium filter - use ANY for better performance
	if len(filter.Medium) > 0 {
		conditions = append(conditions, fmt.Sprintf("medium = ANY($%d)", argIndex))
		mediumStrings := make([]string, len(filter.Medium))
		for i, medium := range filter.Medium {
			mediumStrings[i] = string(medium)
		}
		args = append(args, mediumStrings)
		argIndex++
	}

	// Source filter - use ANY for better performance
	if len(filter.Source) > 0 {
		conditions = append(conditions, fmt.Sprintf("transaction_source = ANY($%d)", argIndex))
		sourceStrings := make([]string, len(filter.Source))
		for i, source := range filter.Source {
			sourceStrings[i] = string(source)
		}
		args = append(args, sourceStrings)
		argIndex++
	}

	// QR Tag filter - use ANY for better performance
	if len(filter.QRTag) > 0 {
		conditions = append(conditions, fmt.Sprintf("qr_tag = ANY($%d)", argIndex))
		args = append(args, filter.QRTag)
		argIndex++
	}

	// Amount range filter - use BETWEEN for better performance
	if filter.AmountMin != nil && filter.AmountMax != nil {
		conditions = append(conditions, fmt.Sprintf("amount BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, *filter.AmountMin, *filter.AmountMax)
		argIndex += 2
	} else if filter.AmountMin != nil {
		conditions = append(conditions, fmt.Sprintf("amount >= $%d", argIndex))
		args = append(args, *filter.AmountMin)
		argIndex++
	} else if filter.AmountMax != nil {
		conditions = append(conditions, fmt.Sprintf("amount <= $%d", argIndex))
		args = append(args, *filter.AmountMax)
		argIndex++
	}

	// Merchant ID filter - use ANY for better performance
	if len(filter.MerchantID) > 0 {
		conditions = append(conditions, fmt.Sprintf("merchant_id = ANY($%d)", argIndex))
		args = append(args, filter.MerchantID)
		argIndex++
	}

	return strings.Join(conditions, " AND "), args
}

// Helper function to get date truncation string
func (r *TransactionRepositoryImpl) getDateTruncation(unit entity.DateUnit) string {
	switch unit {
	case entity.DAY:
		return "day"
	case entity.WEEK:
		return "week"
	case entity.MONTH:
		return "month"
	case entity.YEAR:
		return "year"
	default:
		return "day"
	}
}

// Helper function to format date labels
func (r *TransactionRepositoryImpl) formatDateLabel(date time.Time, unit entity.DateUnit) string {
	switch unit {
	case entity.DAY:
		return date.Format("2006-01-02")
	case entity.WEEK:
		return fmt.Sprintf("Week of %s", date.Format("2006-01-02"))
	case entity.MONTH:
		return date.Format("2006-01")
	case entity.YEAR:
		return date.Format("2006")
	default:
		return date.Format("2006-01-02")
	}
}

// Helper function to calculate period comparison
func (r *TransactionRepositoryImpl) calculatePeriodComparison(ctx context.Context, filter *entity.AnalyticsFilter, merchantID uuid.UUID, current *entity.TransactionAnalytics) (*entity.PeriodComparison, error) {
	// Get previous period dates
	prevStart, prevEnd := filter.GetPreviousPeriod()

	// Create filter for previous period
	prevFilter := &entity.AnalyticsFilter{
		StartDate:  prevStart,
		EndDate:    prevEnd,
		Status:     filter.Status,
		Type:       filter.Type,
		Medium:     filter.Medium,
		Source:     filter.Source,
		QRTag:      filter.QRTag,
		AmountMin:  filter.AmountMin,
		AmountMax:  filter.AmountMax,
		MerchantID: filter.MerchantID,
	}

	// Get previous period analytics
	previous, err := r.GetTransactionAnalytics(ctx, prevFilter, merchantID)
	if err != nil {
		return nil, err
	}

	// Calculate percentage changes
	comparison := &entity.PeriodComparison{}

	if previous.TotalTransactions > 0 {
		comparison.TransactionCountChange = ((float64(current.TotalTransactions) - float64(previous.TotalTransactions)) / float64(previous.TotalTransactions)) * 100
	}

	if previous.TotalAmount > 0 {
		comparison.AmountChange = ((current.TotalAmount - previous.TotalAmount) / previous.TotalAmount) * 100
	}

	if previous.TotalMerchantNet > 0 {
		comparison.MerchantNetChange = ((current.TotalMerchantNet - previous.TotalMerchantNet) / previous.TotalMerchantNet) * 100
	}

	// Note: Success rate calculation removed since we no longer have status breakdown
	// Users can filter by status to get specific success/failure analytics

	return comparison, nil
}
