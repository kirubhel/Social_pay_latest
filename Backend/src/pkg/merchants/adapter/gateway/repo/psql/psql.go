package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/socialpay/socialpay/src/pkg/merchants/usecase"
)

type PsqlRepo struct {
	log    *log.Logger
	db     *sql.DB
	schema string
}

func NewPsqlRepo(log *log.Logger, db *sql.DB) (usecase.Repository, error) {
	log.SetPrefix("[AUTH] [GATEWAY] [REPO] [PSQL] ")

	var _schema = "merchants"
	tableScripts := map[string]string{
		"merchants": fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %s.merchants (
				id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id uuid NOT NULL REFERENCES auth.users(id) ON DELETE RESTRICT,
				legal_name VARCHAR(255) NOT NULL,
				trading_name VARCHAR(255),
				business_registration_number VARCHAR(100) NOT NULL UNIQUE,
				tax_identification_number VARCHAR(100) NOT NULL UNIQUE,
				business_type VARCHAR(100) NOT NULL, -- e.g., 'retail', 'ecommerce', 'betting'
				industry_category VARCHAR(100),
				is_betting_company BOOLEAN DEFAULT FALSE,
				lottery_certificate_number VARCHAR(100), -- Only for betting companies
				website_url VARCHAR(255),
				established_date DATE,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				status VARCHAR(50) NOT NULL DEFAULT 'pending_verification' -- pending_verification, active, suspended, terminated
			);`, _schema),
		// Add other tables as needed

		"documents": fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s.documents (
				id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
				merchant_id uuid NOT NULL REFERENCES merchants.merchants(id) ON DELETE CASCADE,
				document_type VARCHAR(100) NOT NULL, -- 'business_license', 'tin_certificate', 'lottery_certificate', 'bank_statement'
				document_number VARCHAR(100),
				file_url VARCHAR(255) NOT NULL,
				verified_by uuid REFERENCES auth.users(id),
				verified_at TIMESTAMPTZ,
				status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, approved, rejected
				rejection_reason TEXT,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);`, _schema),

		"addresses": fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s.addresses (
				id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
				merchant_id uuid NOT NULL REFERENCES merchants.merchants(id) ON DELETE CASCADE,
				personal_name VARCHAR(100) NOT NULL,
				email VARCHAR(100) NOT NULL,
				region VARCHAR(100) NOT NULL,
				city VARCHAR(100) NOT NULL,
				sub_city VARCHAR(100) NOT NULL,
				woreda VARCHAR(100) NOT NULL,
				Phone_number VARCHAR(50) NOT NULL,
				secondary_phone_number VARCHAR(50),
				postal_code VARCHAR(50),
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);`,
			_schema),
	}

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create schema
	if _, err := tx.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", _schema)); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	// Create tables
	for _, script := range tableScripts {
		if _, err := tx.Exec(script); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &PsqlRepo{
		log:    log,
		db:     db,
		schema: _schema,
	}, nil
}
