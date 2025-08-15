package usecase

import "github.com/socialpay/socialpay/src/pkg/org/core/entity"

type TINChecker interface {
	CheckTIN(tin string, usecase Usecase) (*entity.Organization, error)
	CheckTINRaw(tin string) (map[string]interface{}, error)
}
