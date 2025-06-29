package usecase

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/socialpay/socialpay/src/pkg/key/core/entity"
)

type Interactor interface {
	CreateAPIKey(merchantID, service string, expiryDate time.Time, store string, commissionFrom bool) (*entity.APIKey, error)
	ValidateAPIToken(token string) (*entity.APIKey, error)
	ToggleKey(token string, enable bool) error
}

type MerchantRepository interface {
	GetMerchantID(cookie http.Cookie) (string, error)
}

type MySQLKeyRepository struct {
	DB *sql.DB
}

func NewMySQLKeyRepository(db *sql.DB) *MySQLKeyRepository {
	return &MySQLKeyRepository{DB: db}
}

// Implement all KeyRepository methods for MySQLKeyRepository...
