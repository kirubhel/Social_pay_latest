package repo

import (
	"github.com/socialpay/socialpay/src/pkg/gateways/usecase"

	"database/sql"
	"log"
)

type PsqlRepo struct {
	log    *log.Logger
	db     *sql.DB
	schema string
}

func NewPsqlRepo(log *log.Logger, db *sql.DB) (usecase.Repo, error) {

	log.SetPrefix("[AUTH] [GATEWAY] [REPO] [PSQL] ")

	var _schema = "merchants"

	return PsqlRepo{log, db, _schema}, nil
}

/*

	// Device


	// // Phone


*/
