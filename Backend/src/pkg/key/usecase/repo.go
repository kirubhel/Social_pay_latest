// github.com/socialpay/socialpay/src/pkg/key/core/usecase/repository.go
package usecase

import "github.com/socialpay/socialpay/src/pkg/key/core/entity"

type KeyRepository interface {
	Save(apiKey *entity.APIKey) error
	FindByToken(token string) (*entity.APIKey, error)
	UpdateStatus(token string, enabled bool) error
	// Add any other required methods
}
