package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/core/repository"
	"github.com/socialpay/socialpay/src/pkg/v2_merchant/utils"
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

// GetMerchants retrieves list of merchants
func (r *merchantRepository) GetMerchants(ctx context.Context, params entity.GetMerchantsParams) (*entity.MerchantsResponse, error) {
	fmt.Println("Req params -> ", params)
	rows, err := r.queries.SearchMerchants(ctx, SearchMerchantsParams{
		Lower:   params.Text,
		Offset:  int32(params.Skip),
		Limit:   int32(params.Take),
		Column4: params.StartDate,
		Column5: params.EndDate,
		Status:  params.Status,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get merchant: %w", err)
	}

	var merchants []entity.MerchantDetails

	for _, row := range rows {
		merchant := r.convertMerchantToEntity(MerchantsMerchant{
			ID:                         row.ID,
			UserID:                     row.UserID,
			LegalName:                  row.LegalName,
			TradingName:                row.TradingName,
			BusinessRegistrationNumber: row.BusinessRegistrationNumber,
			TaxIdentificationNumber:    row.TaxIdentificationNumber,
			BusinessType:               row.BusinessType,
			IndustryCategory:           row.IndustryCategory,
			IsBettingCompany:           row.IsBettingCompany,
			LotteryCertificateNumber:   row.LotteryCertificateNumber,
			WebsiteUrl:                 row.WebsiteUrl,
			EstablishedDate:            row.EstablishedDate,
			CreatedAt:                  row.CreatedAt,
			UpdatedAt:                  row.UpdatedAt,
			Status:                     row.Status,
		})

		// Get contacts
		contacts, err := r.GetMerchantContacts(ctx, merchant.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get merchant contacts: %w", err)
		}

		// Get documents
		documents, err := r.GetMerchantDocuments(ctx, merchant.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get merchant documents: %w", err)
		}

		merchants = append(merchants, entity.MerchantDetails{
			Merchant:  *merchant,
			Contacts:  contacts,
			Documents: documents,
		})
	}

	var totalCount int
	if len(rows) > 0 {
		totalCount = int(rows[0].TotalCount)
	} else {
		totalCount = 0
	}

	return &entity.MerchantsResponse{
		Count:     totalCount,
		Merchants: merchants,
	}, nil
}

// GetAllMerchants retrieves all merchants
func (r *merchantRepository) GetAllMerchants(ctx context.Context) ([]entity.Merchant, error) {
	// Get all merchants
	merchantsMerchants, err := r.queries.GetAllMerchants(ctx)
	if err != nil {
		return nil, err
	}

	var merchants []entity.Merchant
	for _, merchant := range merchantsMerchants {
		merchant := &entity.Merchant{
			ID: merchant.ID,
		}

		merchants = append(merchants, *merchant)
	}

	return merchants, nil
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

	fmt.Println(details)

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

// UpdateMerchant updates merchant
func (r *merchantRepository) UpdateMerchant(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantRequest) error {
	businessInfo := req.BusinessInfo
	personalInfo := req.PersonalInfo
	documents := req.Documents

	err := r.queries.UpdateMerchant(ctx, UpdateMerchantParams{
		ID:                         merchantID,
		LegalName:                  *businessInfo.LegalName,
		TradingName:                utils.ToNullString(businessInfo.TradingName),
		BusinessRegistrationNumber: *businessInfo.BusinessRegistrationNumber,
		BusinessType:               *businessInfo.BusinessType,
		TaxIdentificationNumber:    *businessInfo.TaxIdentificationNumber,
		IndustryCategory:           utils.ToNullString(businessInfo.IndustryCategory),
		IsBettingCompany:           utils.ToNullBool(businessInfo.IsBettingCompany),
		LotteryCertificateNumber:   utils.ToNullString(businessInfo.LotteryCertificateNumber),
		WebsiteUrl:                 utils.ToNullString(businessInfo.WebsiteURL),
		EstablishedDate:            utils.ToNullTime(businessInfo.EstablishedDate),
		Status:                     string(*businessInfo.Status),
	})

	if err != nil {
		return fmt.Errorf("failed to update merchant business info: %w", err)
	}

	id := uuid.New()
	err = r.queries.CreateMerchantContact(ctx, CreateMerchantContactParams{
		ID:          id,
		MerchantID:  merchantID,
		FirstName:   personalInfo.FirstName,
		LastName:    personalInfo.LastName,
		Email:       personalInfo.Email,
		PhoneNumber: personalInfo.PhoneNumber,
	})

	if err != nil {
		return fmt.Errorf("failed to add merchant personal info: %w", err)
	}

	for _, document := range documents {
		id := uuid.New()
		err := r.queries.CreateMerchantDocument(ctx, CreateMerchantDocumentParams{
			ID:           id,
			MerchantID:   merchantID,
			DocumentType: document.DocumentType,
			FileUrl:      document.FileUrl,
			Status:       document.Status,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})

		if err != nil {
			return fmt.Errorf("failed to create merchant document: %w", err)
		}
	}

	return nil
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
		AddressType:    a.AddressType.String,
		StreetAddress1: a.StreetAddress1.String,
		StreetAddress2: streetAddress2,
		City:           a.City.String,
		Region:         a.Region.String,
		PostalCode:     postalCode,
		Country:        a.Country.String,
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

// UpdateMerchantStatus updates merchant status
func (r *merchantRepository) UpdateMerchantStatus(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantStatusRequest) error {
	err := r.queries.UpdateMerchantStatus(ctx, UpdateMerchantStatusParams{
		ID:     merchantID,
		Status: req.Status,
	})

	if err != nil {
		return fmt.Errorf("failed to update merchant status: %s", err.Error())
	}
	return nil
}

// UpdateMerchantContact updates merchant contact info
func (r *merchantRepository) UpdateMerchantContact(ctx context.Context, id uuid.UUID, req *entity.UpdateMerchantContactRequest) error {
	err := r.queries.UpdateMerchantContact(ctx, UpdateMerchantContactParams{
		ID:          id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		IsVerified:  utils.ToNullBool(&req.IsVerified),
	})

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"idx_contacts_email_unique\"" {
			return fmt.Errorf("failed to update merchant contact: email already used")
		}
		return fmt.Errorf("failed to update merchant contact: %s", err.Error())
	}
	return nil
}

// UpdateMerchantDocument updates merchant document info
func (r *merchantRepository) UpdateMerchantDocument(ctx context.Context, id uuid.UUID, req *entity.UpdateMerchantDocumentRequest) error {
	now := time.Now()
	err := r.queries.UpdateMerchantDocument(ctx, UpdateMerchantDocumentParams{
		ID:              id,
		VerifiedBy:      utils.ToNullUUID(req.VerifiedBy),
		FileUrl:         req.FileUrl,
		Status:          req.Status,
		VerifiedAt:      utils.ToNullTime(&now),
		RejectionReason: utils.ToNullString(req.RejectionReason),
	})

	if err != nil {
		return fmt.Errorf("failed to update merchant document: %w", err)
	}

	return nil
}

// UpdateAdminMerchant updates merchant by admin
func (r *merchantRepository) UpdateAdminMerchant(ctx context.Context, merchantID uuid.UUID, req *entity.UpdateMerchantRequest) error {
	businessInfo := req.BusinessInfo
	personalInfo := req.PersonalInfo
	documents := req.Documents

	// Update merchant business information
	err := r.queries.UpdateMerchant(ctx, UpdateMerchantParams{
		ID:                         merchantID,
		LegalName:                  *businessInfo.LegalName,
		TradingName:                utils.ToNullString(businessInfo.TradingName),
		BusinessRegistrationNumber: *businessInfo.BusinessRegistrationNumber,
		BusinessType:               *businessInfo.BusinessType,
		TaxIdentificationNumber:    *businessInfo.TaxIdentificationNumber,
		IndustryCategory:           utils.ToNullString(businessInfo.IndustryCategory),
		IsBettingCompany:           utils.ToNullBool(businessInfo.IsBettingCompany),
		LotteryCertificateNumber:   utils.ToNullString(businessInfo.LotteryCertificateNumber),
		WebsiteUrl:                 utils.ToNullString(businessInfo.WebsiteURL),
		EstablishedDate:            utils.ToNullTime(businessInfo.EstablishedDate),
		Status:                     string(*businessInfo.Status),
	})

	if err != nil {
		return fmt.Errorf("failed to update merchant business info: %w", err)
	}

	// Get existing contacts for this merchant
	existingContacts, err := r.queries.GetMerchantContacts(ctx, merchantID)
	if err != nil {
		return fmt.Errorf("failed to get existing contacts: %w", err)
	}

	// Update or create contact
	if len(existingContacts) > 0 {
		// Update the first contact (assuming it's the primary contact)
		contact := existingContacts[0]
		err = r.queries.UpdateMerchantContact(ctx, UpdateMerchantContactParams{
			ID:          contact.ID,
			FirstName:   personalInfo.FirstName,
			LastName:    personalInfo.LastName,
			Email:       personalInfo.Email,
			PhoneNumber: personalInfo.PhoneNumber,
			IsVerified:  sql.NullBool{Valid: true, Bool: true}, // Admin updates are verified
		})
		if err != nil {
			return fmt.Errorf("failed to update merchant contact: %w", err)
		}
	} else {
		// Create new contact if none exists
		id := uuid.New()
		err = r.queries.CreateMerchantContact(ctx, CreateMerchantContactParams{
			ID:          id,
			MerchantID:  merchantID,
			ContactType: "primary",
			FirstName:   personalInfo.FirstName,
			LastName:    personalInfo.LastName,
			Email:       personalInfo.Email,
			PhoneNumber: personalInfo.PhoneNumber,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		if err != nil {
			return fmt.Errorf("failed to create merchant contact: %w", err)
		}
	}

	// Handle documents - update existing ones by ID
	for _, document := range documents {
		// Check if document exists
		_, err := r.queries.GetMerchantDocument(ctx, document.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Document doesn't exist, create new one
				err := r.queries.CreateMerchantDocument(ctx, CreateMerchantDocumentParams{
					ID:           document.ID,
					MerchantID:   merchantID,
					DocumentType: document.DocumentType,
					FileUrl:      document.FileUrl,
					Status:       document.Status,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				})
				if err != nil {
					return fmt.Errorf("failed to create merchant document: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check document existence: %w", err)
			}
		} else {
			// Document exists, update it
			err := r.queries.UpdateMerchantDocumentWithType(ctx, UpdateMerchantDocumentWithTypeParams{
				ID:              document.ID,
				DocumentType:    document.DocumentType,
				FileUrl:         document.FileUrl,
				Status:          document.Status,
				VerifiedBy:      utils.ToNullUUID(nil),          // Admin updates are verified
				VerifiedAt:      utils.ToNullTime(&time.Time{}), // Will be set by database trigger
				RejectionReason: utils.ToNullString(nil),
			})
			if err != nil {
				return fmt.Errorf("failed to update merchant document: %w", err)
			}
		}
	}

	return nil
}

// DeleteMerchant soft deletes merchant
func (r *merchantRepository) DeleteMerchant(ctx context.Context, merchantID uuid.UUID) error {
	now := time.Now()
	err := r.queries.DeleteMerchant(ctx, DeleteMerchantParams{
		ID:        merchantID,
		DeletedAt: utils.ToNullTime(&now),
	})
	if err != nil {
		return fmt.Errorf("failed to delete merchant: %w", err)
	}

	return nil
}

// DeleteMerchants soft deletes list of merchants
func (r *merchantRepository) DeleteMerchants(ctx context.Context, ids []uuid.UUID) error {
	fmt.Println("merchants to be deleted -> ", ids)
	err := r.queries.DeleteMerchants(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to delete merchants: %w", err)
	}

	return nil
}

// GetMerchantStats retrieves merchant statistics
func (r *merchantRepository) GetMerchantStats(ctx context.Context) (*entity.MerchantStats, error) {
	stats, err := r.queries.GetMerchantStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant stats: %w", err)
	}

	return &entity.MerchantStats{
		TotalMerchants:  stats.TotalMerchants,
		ActiveMerchants: stats.ActiveMerchants,
		PendingKyc:      stats.PendingKyc,
		NewThisMonth:    stats.NewThisMonth,
	}, nil
}
