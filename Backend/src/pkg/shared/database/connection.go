package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	sharedDB *sql.DB
)

// GetSharedConnection returns a shared database connection instance
// This should be used by all modules instead of creating separate connections
func GetSharedConnection() (*sql.DB, error) {
	if sharedDB != nil {
		return sharedDB, nil
	}

	return initializeSharedConnection()
}

// initializeSharedConnection creates and configures the shared database connection
func initializeSharedConnection() (*sql.DB, error) {
	envFilePath := ".env"
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Println("Error loading .env file:", err)
		// Don't return error, environment variables might be set elsewhere
	}

	var (
		host, _     = os.LookupEnv("DB_HOST")
		user, _     = os.LookupEnv("DB_USER")
		password, _ = os.LookupEnv("DB_PASS")
		dbName, _   = os.LookupEnv("DB_NAME")
		sslMode, _  = os.LookupEnv("SSL_MODE")
		port, _     = os.LookupEnv("DB_PORT")
	)

	// Create connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslMode)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool for production use
	// These settings are appropriate for handling concurrent operations including cron jobs
	db.SetMaxIdleConns(20)                  // Keep 20 idle connections ready
	db.SetMaxOpenConns(100)                 // Allow up to 100 total connections
	db.SetConnMaxLifetime(15 * time.Minute) // Close connections after 15 minutes
	db.SetConnMaxIdleTime(5 * time.Minute)  // Close idle connections after 5 minutes

	// Test the connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Shared database connection established successfully")
	log.Printf("Connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v",
		100, 20, 15*time.Minute, 5*time.Minute)

	sharedDB = db
	return sharedDB, nil
}

// CloseSharedConnection closes the shared database connection
// This should be called during application shutdown
func CloseSharedConnection() error {
	if sharedDB != nil {
		log.Println("Closing shared database connection")
		err := sharedDB.Close()
		sharedDB = nil
		return err
	}
	return nil
}

// GetConnectionStats returns current connection pool statistics
func GetConnectionStats() sql.DBStats {
	if sharedDB == nil {
		return sql.DBStats{}
	}
	return sharedDB.Stats()
}
