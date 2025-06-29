package psql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/socialpay/socialpay/src/pkg/notifications/usecase"
)

type PsqlRepo struct {
	log *log.Logger
	db  *sql.DB
}

func New(log *log.Logger, db *sql.DB) (usecase.Repo, error) {

	var _schema = "accounts"
	fmt.Println("Using schema:", _schema)

	return PsqlRepo{log: log, db: db}, nil
}
