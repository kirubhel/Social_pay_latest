package migration

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func InitMigiration(path string, db_url string) *migrate.Migrate {

	m_file, err := migrate.New(fmt.Sprintf("file://%s", path), db_url)

	if err != nil {

		log.Fatal(err)
	}

	return m_file

}

func UpMigiration(m *migrate.Migrate) {

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// TODO update to structural loging
		log.Printf("ERR :: migration error HAPPEN ")

	}
}
