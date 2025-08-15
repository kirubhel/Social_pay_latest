package entity

import validation "github.com/go-ozzo/ozzo-validation"

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func (p Pagination) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Page, validation.Min(1)),
		validation.Field(&p.PageSize, validation.Min(1),
			validation.Max(100)),
	)
}
