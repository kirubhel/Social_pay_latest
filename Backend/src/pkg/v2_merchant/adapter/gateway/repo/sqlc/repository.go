package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/repository"
)

type merchantRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewMerchantRepository creates a new merchant repository
func NewMerchantRepository(db *sql.DB) repository.Repository {
	return &merchantRepository{
		db:      db,
		queries: New(db),
	}
}

// GetMerchant retrieves a merchant by its ID
func (r *merchantRepository) GetMerchant(ctx context.Context, id uuid.UUID) (*entity.Merchant, error) {
	merchant, err := r.queries.GetMerchant(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant not found")
		}
		return nil, fmt.Errorf("failed to get merchant: %w", err)
	}

	return r.convertMerchantToEntity(merchant), nil
}

// GetMerchantDetails retrieves complete merchant information with related data
func (r *merchantRepository) GetMerchantDetails(ctx context.Context, id uuid.UUID) (*entity.MerchantDetails, error) {
	// Get merchant
	merchant, err := r.GetMerchant(ctx, id)
	if err != nil {
		return nil, err
	}

	details := &entity.MerchantDetails{
		Merchant: *merchant,
	}

	// Get addresses
	addresses, err := r.GetMerchantAddresses(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant addresses: %w", err)
	}
	details.Addresses = addresses

	// Get contacts
	contacts, err := r.GetMerchantContacts(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant contacts: %w", err)
	}
	details.Contacts = contacts

	// Get documents
	documents, err := r.GetMerchantDocuments(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant documents: %w", err)
	}
	details.Documents = documents

	// Get bank accounts
	bankAccounts, err := r.GetMerchantBankAccounts(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant bank accounts: %w", err)
	}
	details.BankAccounts = bankAccounts

	// Get settings
	settings, err := r.GetMerchantSettings(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get merchant settings: %w", err)
	}
	if settings != nil {
		details.Settings = settings
	}

	return details, nil
}

// GetMerchantByUserID retrieves a merchant by user ID
func (r *merchantRepository) GetMerchantByUserID(ctx context.Context, userID uuid.UUID) (*entity.Merchant, error) {
	merchant, err := r.queries.GetMerchantByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant not found")
		}
		return nil, fmt.Errorf("failed to get merchant by user ID: %w", err)
	}

	return r.convertMerchantToEntity(merchant), nil
}

// GetMerchantAddresses retrieves all addresses for a merchant
func (r *merchantRepository) GetMerchantAddresses(ctx context.Context, merchantID uuid.UUID) ([]entity.MerchantAddress, error) {
	addresses, err := r.queries.GetMerchantAddresses(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant addresses: %w", err)
	}

	result := make([]entity.MerchantAddress, len(addresses))
	for i, addr := range addresses {
		result[i] = r.convertAddressToEntity(addr)
	}

	return result, nil
}

// GetMerchantContacts retrieves all contacts for a merchant
func (r *merchantRepository) GetMerchantContacts(ctx context.Context, merchantID uuid.UUID) ([]entity.MerchantContact, error) {
	contacts, err := r.queries.GetMerchantContacts(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant contacts: %w", err)
	}

	result := make([]entity.MerchantContact, len(contacts))
	for i, contact := range contacts {
		result[i] = r.convertContactToEntity(contact)
	}

	return result, nil
}

// GetMerchantDocuments retrieves all documents for a merchant
func (r *merchantRepository) GetMerchantDocuments(ctx context.Context, merchantID uuid.UUID) ([]entity.MerchantDocument, error) {
	documents, err := r.queries.GetMerchantDocuments(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant documents: %w", err)
	}

	result := make([]entity.MerchantDocument, len(documents))
	for i, doc := range documents {
		result[i] = r.convertDocumentToEntity(doc)
	}

	return result, nil
}

// GetMerchantBankAccounts retrieves all bank accounts for a merchant
func (r *merchantRepository) GetMerchantBankAccounts(ctx context.Context, merchantID uuid.UUID) ([]entity.MerchantBankAccount, error) {
	bankAccounts, err := r.queries.GetMerchantBankAccounts(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant bank accounts: %w", err)
	}

	result := make([]entity.MerchantBankAccount, len(bankAccounts))
	for i, account := range bankAccounts {
		result[i] = r.convertBankAccountToEntity(account)
	}

	return result, nil
}

// GetMerchantSettings retrieves settings for a merchant
func (r *merchantRepository) GetMerchantSettings(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantSettings, error) {
	settings, err := r.queries.GetMerchantSettings(ctx, merchantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get merchant settings: %w", err)
	}

	return r.convertSettingsToEntity(settings), nil
}

// Helper methods to convert SQLC models to entities

func (r *merchantRepository) convertMerchantToEntity(m MerchantsMerchant) *entity.Merchant {
	var tradingName *string
	if m.TradingName.Valid {
		tradingName = &m.TradingName.String
	}

	var industryCategory *string
	if m.IndustryCategory.Valid {
		industryCategory = &m.IndustryCategory.String
	}

	var lotteryCertificateNumber *string
	if m.LotteryCertificateNumber.Valid {
		lotteryCertificateNumber = &m.LotteryCertificateNumber.String
	}

	var websiteURL *string
	if m.WebsiteUrl.Valid {
		websiteURL = &m.WebsiteUrl.String
	}

	var establishedDate *time.Time
	if m.EstablishedDate.Valid {
		establishedDate = &m.EstablishedDate.Time
	}

	return &entity.Merchant{
		ID:                         m.ID,
		UserID:                     m.UserID,
		LegalName:                  m.LegalName,
		TradingName:                tradingName,
		BusinessRegistrationNumber: m.BusinessRegistrationNumber,
		TaxIdentificationNumber:    m.TaxIdentificationNumber,
		BusinessType:               m.BusinessType,
		IndustryCategory:           industryCategory,
		IsBettingCompany:           m.IsBettingCompany.Bool,
		LotteryCertificateNumber:   lotteryCertificateNumber,
		WebsiteURL:                 websiteURL,
		EstablishedDate:            establishedDate,
		CreatedAt:                  m.CreatedAt,
		UpdatedAt:                  m.UpdatedAt,
		Status:                     entity.MerchantStatus(m.Status),
	}
}

func (r *merchantRepository) convertAddressToEntity(a MerchantsAddress) entity.MerchantAddress {
	var streetAddress2 *string
	if a.StreetAddress2.Valid {
		streetAddress2 = &a.StreetAddress2.String
	}

	var postalCode *string
	if a.PostalCode.Valid {
		postalCode = &a.PostalCode.String
	}

	return entity.MerchantAddress{
		ID:             a.ID,
		MerchantID:     a.MerchantID,
		AddressType:    a.AddressType,
		StreetAddress1: a.StreetAddress1,
		StreetAddress2: streetAddress2,
		City:           a.City,
		Region:         a.Region,
		PostalCode:     postalCode,
		Country:        a.Country,
		IsPrimary:      a.IsPrimary.Bool,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

func (r *merchantRepository) convertContactToEntity(c MerchantsContact) entity.MerchantContact {
	var position *string
	if c.Position.Valid {
		position = &c.Position.String
	}

	return entity.MerchantContact{
		ID:          c.ID,
		MerchantID:  c.MerchantID,
		ContactType: c.ContactType,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
		Position:    position,
		IsVerified:  c.IsVerified.Bool,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func (r *merchantRepository) convertDocumentToEntity(d MerchantsDocument) entity.MerchantDocument {
	var documentNumber *string
	if d.DocumentNumber.Valid {
		documentNumber = &d.DocumentNumber.String
	}

	var fileHash *string
	if d.FileHash.Valid {
		fileHash = &d.FileHash.String
	}

	var verifiedBy *uuid.UUID
	if d.VerifiedBy.Valid {
		verifiedBy = &d.VerifiedBy.UUID
	}

	var verifiedAt *time.Time
	if d.VerifiedAt.Valid {
		verifiedAt = &d.VerifiedAt.Time
	}

	var rejectionReason *string
	if d.RejectionReason.Valid {
		rejectionReason = &d.RejectionReason.String
	}

	return entity.MerchantDocument{
		ID:              d.ID,
		MerchantID:      d.MerchantID,
		DocumentType:    d.DocumentType,
		DocumentNumber:  documentNumber,
		FileURL:         d.FileUrl,
		FileHash:        fileHash,
		VerifiedBy:      verifiedBy,
		VerifiedAt:      verifiedAt,
		Status:          d.Status,
		RejectionReason: rejectionReason,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

func (r *merchantRepository) convertBankAccountToEntity(b MerchantsBankAccount) entity.MerchantBankAccount {
	var branchCode *string
	if b.BranchCode.Valid {
		branchCode = &b.BranchCode.String
	}

	var verificationDocumentID *uuid.UUID
	if b.VerificationDocumentID.Valid {
		verificationDocumentID = &b.VerificationDocumentID.UUID
	}

	return entity.MerchantBankAccount{
		ID:                     b.ID,
		MerchantID:             b.MerchantID,
		AccountHolderName:      b.AccountHolderName,
		BankName:               b.BankName,
		BankCode:               b.BankCode,
		BranchCode:             branchCode,
		AccountNumber:          b.AccountNumber,
		AccountType:            b.AccountType,
		Currency:               b.Currency,
		IsPrimary:              b.IsPrimary.Bool,
		IsVerified:             b.IsVerified.Bool,
		VerificationDocumentID: verificationDocumentID,
		CreatedAt:              b.CreatedAt,
		UpdatedAt:              b.UpdatedAt,
	}
}

func (r *merchantRepository) convertSettingsToEntity(s MerchantsSetting) *entity.MerchantSettings {
	var riskSettings *string
	if s.RiskSettings.Valid {
		riskSettingsStr := string(s.RiskSettings.RawMessage)
		riskSettings = &riskSettingsStr
	}

	var checkoutTheme *string
	if s.CheckoutTheme.Valid {
		checkoutTheme = &s.CheckoutTheme.String
	}

	var webhookURL *string
	if s.WebhookUrl.Valid {
		webhookURL = &s.WebhookUrl.String
	}

	var webhookSecret *string
	if s.WebhookSecret.Valid {
		webhookSecret = &s.WebhookSecret.String
	}

	settlementFrequency := "daily" // default value
	if s.SettlementFrequency.Valid {
		settlementFrequency = s.SettlementFrequency.String
	}

	return &entity.MerchantSettings{
		MerchantID:          s.MerchantID,
		DefaultCurrency:     s.DefaultCurrency,
		DefaultLanguage:     s.DefaultLanguage,
		CheckoutTheme:       checkoutTheme,
		EnableWebhooks:      s.EnableWebhooks.Bool,
		WebhookURL:          webhookURL,
		WebhookSecret:       webhookSecret,
		AutoSettlement:      s.AutoSettlement.Bool,
		SettlementFrequency: settlementFrequency,
		RiskSettings:        riskSettings,
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
	}
}
