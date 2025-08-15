// github.com/socialpay/socialpay/src/pkg/key/repo/psql_repo.go
package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/socialpay/socialpay/src/pkg/merchants/core/entity"

	"github.com/google/uuid"
)

func (r *PsqlRepo) Save(merchant *entity.Merchant) error {
	query := fmt.Sprintf(`
        INSERT INTO merchants.%s (
			id, 
			user_id, 
			legal_name, 
			trading_name, 
			business_registration_number, 
			tax_identification_number, 
			industry_category, 
			business_type, 
			is_betting_company, 
			lottery_certificate_number, 
			website_url, 
			established_date, 
			created_at, 
			updated_at, 
			status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW(), $13
		)`, r.schema)

	_, err := r.db.Exec(query,
		uuid.New(),
		merchant.UserID,
		merchant.LegalName,
		merchant.TradingName,
		merchant.BusinessRegNumber,
		merchant.TaxIdentifier,
		merchant.IndustryType,
		merchant.BusinessType,
		merchant.IsBettingClient,
		merchant.LoyaltyCertificate,
		merchant.WebsiteURL,
		merchant.EstablishedDate,
		merchant.Status,
	)
	return err
}

func (r *PsqlRepo) FindByUserID(userID uuid.UUID) (*entity.Merchant, error) {
	query := fmt.Sprintf(`
		SELECT id, user_id, legal_name, trading_name, business_registration_number, 
			   tax_identification_number, industry_category, business_type, 
			   is_betting_company, lottery_certificate_number, website_url, 
			   established_date, created_at, updated_at, status
		FROM merchants.%s 
		WHERE user_id = $1`, r.schema)

	var merchant entity.Merchant
	err := r.db.QueryRow(query, userID).Scan(
		&merchant.MerchantID,
		&merchant.UserID,
		&merchant.LegalName,
		&merchant.TradingName,
		&merchant.BusinessRegNumber,
		&merchant.TaxIdentifier,
		&merchant.IndustryType,
		&merchant.BusinessType,
		&merchant.IsBettingClient,
		&merchant.LoyaltyCertificate,
		&merchant.WebsiteURL,
		&merchant.EstablishedDate,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
		&merchant.Status,
	)

	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (r *PsqlRepo) GetMerchantDetails(userID uuid.UUID) (*entity.MerchantDetails, error) {
	// Get merchant basic info
	merchantQuery := `
        SELECT id, user_id, legal_name, trading_name, business_registration_number, 
               tax_identification_number, industry_category, business_type, 
               is_betting_company, lottery_certificate_number, website_url, 
               established_date, created_at, updated_at, status
        FROM merchants.merchants 
        WHERE user_id = $1`

	var merchant entity.MerchantDetails
	err := r.db.QueryRow(merchantQuery, userID).Scan(
		&merchant.MerchantID,
		&merchant.UserID,
		&merchant.LegalName,
		&merchant.TradingName,
		&merchant.BusinessRegNumber,
		&merchant.TaxIdentifier,
		&merchant.IndustryType,
		&merchant.BusinessType,
		&merchant.IsBettingClient,
		&merchant.LoyaltyCertificate,
		&merchant.WebsiteURL,
		&merchant.EstablishedDate,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
		&merchant.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant not found")
		}
		return nil, fmt.Errorf("failed to get merchant: %w", err)
	}

	// Get address
	addressQuery := `
        SELECT personal_name, phone_number, region, city, sub_city, 
               woreda, postal_code, secondary_phone_number, email
        FROM merchants.addresses 
        WHERE merchant_id = $1`

	var address entity.MerchantAdditionalInfo
	err = r.db.QueryRow(addressQuery, merchant.MerchantID).Scan(
		&address.PersonalName,
		&address.PhoneNumber,
		&address.Region,
		&address.City,
		&address.SubCity,
		&address.Woreda,
		&address.PostalCode,
		&address.SecondaryPhoneNumber,
		&address.Email,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	merchant.Address = &address

	// Get documents
	docsQuery := `
        SELECT document_type, document_number, file_url, status
        FROM merchants.documents 
        WHERE merchant_id = $1`

	rows, err := r.db.Query(docsQuery, merchant.MerchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	defer rows.Close()

	var documents []entity.MerchantDocument
	for rows.Next() {
		var doc entity.MerchantDocument
		if err := rows.Scan(
			&doc.DocumentType,
			&doc.DocumentNumber,
			&doc.FileURL,
			&doc.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}
	merchant.Documents = documents

	return &merchant, nil
}

func (r *PsqlRepo) GetMerchants() ([]entity.Merchant, error) {
	query := `
        SELECT 
            id, user_id, legal_name, trading_name, 
            business_registration_number, tax_identification_number, 
            industry_category, business_type, is_betting_company, 
            lottery_certificate_number, website_url, established_date, 
            created_at, updated_at, status
        FROM merchants.merchants`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var merchants []entity.Merchant
	for rows.Next() {
		var m entity.Merchant
		err := rows.Scan(
			&m.MerchantID,
			&m.UserID,
			&m.LegalName,
			&m.TradingName,
			&m.BusinessRegNumber,
			&m.TaxIdentifier,
			&m.IndustryType,
			&m.BusinessType,
			&m.IsBettingClient,
			&m.LoyaltyCertificate,
			&m.WebsiteURL,
			&m.EstablishedDate,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("row scanning error: %w", err)
		}
		merchants = append(merchants, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return merchants, nil
}

func (r *PsqlRepo) GetMerchantsByUserID(userID uuid.UUID) ([]entity.Merchant, error) {
	query := `
        SELECT 
            id, user_id, legal_name, trading_name, 
            business_registration_number, tax_identification_number, 
            industry_category, business_type, is_betting_company, 
            lottery_certificate_number, website_url, established_date, 
            created_at, updated_at, status
        FROM merchants.merchants
		WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var merchants []entity.Merchant
	for rows.Next() {
		var m entity.Merchant
		err := rows.Scan(
			&m.MerchantID,
			&m.UserID,
			&m.LegalName,
			&m.TradingName,
			&m.BusinessRegNumber,
			&m.TaxIdentifier,
			&m.IndustryType,
			&m.BusinessType,
			&m.IsBettingClient,
			&m.LoyaltyCertificate,
			&m.WebsiteURL,
			&m.EstablishedDate,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("row scanning error: %w", err)
		}
		merchants = append(merchants, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return merchants, nil
}

func (r *PsqlRepo) UpdateMerchantStatus(merchantID uuid.UUID, status entity.MerchantStatus) error {
	query := `
        UPDATE merchants.merchants 
        SET status = $1, updated_at = NOW() 
        WHERE id = $2`

	result, err := r.db.Exec(query, status, merchantID)
	if err != nil {
		return fmt.Errorf("database update error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no merchant found with ID %s", merchantID)
	}

	return nil
}

func (r *PsqlRepo) DeleteMerchant(ctx context.Context, merchantID uuid.UUID) error {

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete from addresses table
	if _, err = tx.ExecContext(ctx,
		`DELETE FROM merchants.addresses WHERE merchant_id = $1`, merchantID); err != nil {
		return fmt.Errorf("failed to delete addresses: %w", err)
	}

	// Delete from documents table
	if _, err = tx.ExecContext(ctx,
		`DELETE FROM merchants.documents WHERE merchant_id = $1`, merchantID); err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}

	if _, err = tx.ExecContext(ctx,
		`DELETE FROM merchants.additional_info WHERE merchant_id = $1`, merchantID); err != nil {
		return fmt.Errorf("failed to delete additional info: %w", err)
	}

	// Delete from main merchants table
	result, err := tx.ExecContext(ctx,
		`DELETE FROM merchants.merchants WHERE id = $1`, merchantID)
	if err != nil {
		return fmt.Errorf("failed to delete merchant: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no merchant found with ID %s", merchantID)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PsqlRepo) UpdateFullMerchant(ctx context.Context, merchantID uuid.UUID,
	merchant *entity.Merchant, address *entity.MerchantAdditionalInfo,
	documents []entity.MerchantDocument) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update main merchant table (excluding status)
	merchantQuery := `
        UPDATE merchants.merchants SET
            legal_name = $1,
            trading_name = $2,
            business_registration_number = $3,
            tax_identification_number = $4,
            industry_category = $5,
            business_type = $6,
            is_betting_company = $7,
            lottery_certificate_number = $8,
            website_url = $9,
            established_date = $10,
            updated_at = NOW()
        WHERE id = $11`

	_, err = tx.ExecContext(ctx, merchantQuery,
		merchant.LegalName,
		merchant.TradingName,
		merchant.BusinessRegNumber,
		merchant.TaxIdentifier,
		merchant.IndustryType,
		merchant.BusinessType,
		merchant.IsBettingClient,
		merchant.LoyaltyCertificate,
		merchant.WebsiteURL,
		merchant.EstablishedDate,
		merchantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update merchant: %w", err)
	}

	// Update or insert address
	addressQuery := `
        INSERT INTO merchants.addresses (
            merchant_id, personal_name, phone_number, region, city, 
            sub_city, woreda, postal_code, secondary_phone_number, email
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (merchant_id) 
        DO UPDATE SET
            personal_name = EXCLUDED.personal_name,
            phone_number = EXCLUDED.phone_number,
            region = EXCLUDED.region,
            city = EXCLUDED.city,
            sub_city = EXCLUDED.sub_city,
            woreda = EXCLUDED.woreda,
            postal_code = EXCLUDED.postal_code,
            secondary_phone_number = EXCLUDED.secondary_phone_number,
            email = EXCLUDED.email`

	_, err = tx.ExecContext(ctx, addressQuery,
		merchantID,
		address.PersonalName,
		address.PhoneNumber,
		address.Region,
		address.City,
		address.SubCity,
		address.Woreda,
		address.PostalCode,
		address.SecondaryPhoneNumber,
		address.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}

	// Delete existing documents and insert new ones
	if _, err = tx.ExecContext(ctx,
		`DELETE FROM merchants.documents WHERE merchant_id = $1`, merchantID); err != nil {
		return fmt.Errorf("failed to delete existing documents: %w", err)
	}

	for _, doc := range documents {
		docQuery := `
            INSERT INTO merchants.documents (
                merchant_id, document_type, document_number, file_url, status
            ) VALUES ($1, $2, $3, $4, $5)`

		if _, err = tx.ExecContext(ctx, docQuery,
			merchantID,
			doc.DocumentType,
			doc.DocumentNumber,
			doc.FileURL,
			doc.Status,
		); err != nil {
			return fmt.Errorf("failed to insert document: %w", err)
		}
	}

	return tx.Commit()
}

func (r *PsqlRepo) GetMerchantByID(merchantID uuid.UUID) (*entity.Merchant, error) {
	query := `
        SELECT id, status 
        FROM merchants.merchants 
        WHERE id = $1`

	var merchant entity.Merchant
	err := r.db.QueryRow(query, merchantID).Scan(&merchant.MerchantID, &merchant.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant not found")
		}
		return nil, fmt.Errorf("database query error: %w", err)
	}

	return &merchant, nil
}

func (r *PsqlRepo) CreateMerchantAdditionalInfo(ctx context.Context, info entity.MerchantAdditionalInfo) error {
	query := `INSERT INTO merchants.addresses (
		merchant_id,
		personal_name,
		phone_number,
		region,
		city,
		sub_city,
		woreda,
		postal_code,
		secondary_phone_number,
		email
	) VALUES (
		$1, -- merchant_id
		$2, -- personal_name
		$3, -- phone_number
		$4, -- region
		$5, -- city
		$6, -- sub_city
		$7, -- woreda
		$8, -- postal_code
		$9, -- secondary_phone_number
		$10 -- email
	);`

	if _, err := r.db.ExecContext(ctx, query,
		info.MerchantID,
		info.PersonalName,
		info.PhoneNumber,
		info.Region,
		info.City,
		info.SubCity,
		info.Woreda,
		info.PostalCode,
		info.SecondaryPhoneNumber,
		info.Email); err != nil {
		return err
	}
	return nil
}

func (r *PsqlRepo) SaveDocument(ctx context.Context, doc entity.MerchantDocument) error {
	query := `INSERT INTO merchants.documents (
		merchant_id,
		document_type,
		document_number,
		file_url,
		status
	) VALUES (
		$1, -- merchant_id
		$2, -- document_type
		$3, -- document_number
		$4, -- file_url
		$5  -- status
	);`

	_, err := r.db.Exec(query,
		doc.MerchantID,
		doc.DocumentType,
		doc.DocumentNumber,
		doc.FileURL,
		doc.Status,
	)
	return err
}
