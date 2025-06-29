package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/socialpay/socialpay/src/pkg/key/usecase"
)

type PsqlRepo struct {
	log    *log.Logger
	db     *sql.DB
	schema string
}

func NewPsqlRepo(log *log.Logger, db *sql.DB) (usecase.KeyRepository, error) {
	log.SetPrefix("[AUTH] [GATEWAY] [REPO] [PSQL] ")

	var _schema = "merchants"
	tableScripts := map[string]string{
		"api_keys": fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %s.api_keys (
                id UUID PRIMARY KEY,
                merchant_id VARCHAR(255) NOT NULL,
                private_key TEXT NOT NULL,
                public_key TEXT NOT NULL,
                api_key VARCHAR(255) UNIQUE NOT NULL,
                service VARCHAR(255),
                expiry_date TIMESTAMP,
                store VARCHAR(255),
                is_active BOOLEAN DEFAULT true,
                created_at TIMESTAMP NOT NULL,
                updated_at TIMESTAMP NOT NULL DEFAULT NOW()
            );`, _schema),
		// Add other tables as needed
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
