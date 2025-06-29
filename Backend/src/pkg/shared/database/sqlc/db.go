package sqlc

import (
	"database/sql"
)

// Queries defines all database queries operations
type Queries struct {
	db *sql.DB
}

// New creates a new Queries instance
func New(db *sql.DB) *Queries {
	return &Queries{
		db: db,
	}
}

// DB is a wrapper around sql.DB
type DB struct {
	*sql.DB
	*Queries
}

// NewDB creates a new DB instance
func NewDB(db *sql.DB) *DB {
	return &DB{
		DB:      db,
		Queries: New(db),
	}
}
