package usecase

import (
	"errors"

	"github.com/socialpay/socialpay/src/pkg/org/core/entity"
)

// Implement the interface with correct signature
func (uc Usecase) CheckTIN(tin string, _ Usecase) (*entity.Organization, error) {
	return uc.tinchecker.CheckTIN(tin, uc)
}

// Additional raw response method
func (uc Usecase) CheckTINRaw(tin string) (map[string]interface{}, error) {
	if rawChecker, ok := uc.tinchecker.(interface {
		CheckTINRaw(tin string) (map[string]interface{}, error)
	}); ok {
		return rawChecker.CheckTINRaw(tin)
	}
	return nil, errors.New("raw checking not supported")
}
