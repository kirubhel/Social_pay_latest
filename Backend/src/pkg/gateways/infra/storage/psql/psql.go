package psql

import (
	"database/sql"
	"log"

	sharedDB "github.com/socialpay/socialpay/src/pkg/shared/database"
)

// New returns the shared database connection instead of creating a new one
// This prevents database connection pool fragmentation
func New(log *log.Logger) (*sql.DB, error) {
	log.SetPrefix(log.Prefix() + " [GATEWAYS] [INFRA] [STORAGE] [PSQL] ")
	log.Println("Getting shared database connection")

	db, err := sharedDB.GetSharedConnection()
	if err != nil {
		log.Println("Failed to get shared database connection: " + err.Error())
		return nil, err
	}

	log.Println("Successfully obtained shared database connection")
	return db, nil
}
