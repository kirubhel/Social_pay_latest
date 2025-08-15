package psql

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/socialpay/socialpay/src/pkg/checkout/service"
)

type CheckoutPSQLRepo struct {
	log *log.Logger
	db  *sql.DB
}

func New(log *log.Logger, db *sql.DB) (service.CheckoutRepo, error) {
	log.SetPrefix("[AUTH] [GATEWAY] [REPO] [PSQL] ")

	// Get the absolute path to the init.sql file
	_, filename, _, _ := runtime.Caller(0) // Gets the path of the current file
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..", "..", "..", "..")
	initPath := filepath.Join(projectRoot, "checkout", "adapter", "gateway", "repo", "psql", "init.sql")

	// Verify the path
	log.Printf("Looking for init.sql at: %s", initPath)
	if _, err := os.Stat(initPath); os.IsNotExist(err) {
		return nil, err
	}

	// Read the SQL script
	script, err := os.ReadFile(initPath)
	if err != nil {
		return nil, err
	}

	// Execute in a transaction
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Defer rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Execute the script
	if _, err = tx.Exec(string(script)); err != nil {
		log.Printf("Failed to execute init script: %v", err)
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &CheckoutPSQLRepo{log, db}, nil
}
