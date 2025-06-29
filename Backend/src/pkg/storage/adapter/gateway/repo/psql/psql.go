package psql

import (
	"database/sql"
	"log"

	"github.com/socialpay/socialpay/src/pkg/storage/usecase"
)

type PsqlRepo struct {
	log *log.Logger
	db  *sql.DB
}

func New(log *log.Logger, db *sql.DB) (usecase.Repo, error) {
	return PsqlRepo{log: log, db: db}, nil
}
