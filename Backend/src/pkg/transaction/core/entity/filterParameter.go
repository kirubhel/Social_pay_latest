package entity

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/filter"
)

type DataFormat string

const (
	FORMATPDF DataFormat = "pdf"
	FORMATCSV DataFormat = "csv"
)

// FilterParameters defines the criteria for filtering transactions
// @Description Parameters for filtering transaction data
type FilterParameters struct {
	// @Description Start date in ISO 8601 format (YYYY-MM-DD)
	// @Example 2023-01-01
	StartDate time.Time `json:"start_date"`
	// @Description End date in ISO 8601 format (YYYY-MM-DD)
	// @Example 2023-12-31
	EndDate time.Time `json:"end_date"`

	Status TransactionStatus `json:"status"`
	Type   TransactionType   `json:"type"`
	Medium TransactionMedium `json:"medium"`
	Format DataFormat        `json:"format"`

	// Searching
	SocialPayTransactionID  string `json:"socialpay_tx_id,omitempty"`  // optional
	ReferenceId           string `json:"reference_id,omitempty"`   // optional
	ProviderTransactionId string `json:"provider_tx_id,omitempty"` // optional

	MerchantID uuid.UUID `json:"merchant_id,omitempty"` // optional

	// Sort
	//@Description ASC or DESC
	Sort string `json:"sort,omitempty"` // optional

}

// Validate validates the filter parameters
func (f *FilterParameters) Validate() error {
	if err := validation.ValidateStruct(f,
		validation.Field(&f.StartDate, validation.Required, validation.By(func(value interface{}) error {
			startDate, ok := value.(time.Time)
			if !ok {
				return errors.New("invalid start date format")
			}
			if startDate.IsZero() {
				return errors.New("start date cannot be empty")
			}
			return nil
		})),
		validation.Field(&f.EndDate, validation.Required, validation.By(func(value interface{}) error {
			endDate, ok := value.(time.Time)
			if !ok {
				return errors.New("invalid end date format")
			}
			if endDate.IsZero() {
				return errors.New("end date cannot be empty")
			}
			return nil
		})),
	); err != nil {
		return err
	}

	// Additional validation for date range
	if f.EndDate.Before(f.StartDate) {
		return errors.New("end date must be after start date")
	}

	// Validate status if provided
	if f.Status != "" {
		if err := validation.Validate(f.Status, validation.In(
			INITIATED,
			PENDING,
			FAILED,
			SUCCESS,
			REFUNDED,
			EXPIRED,
			CANCELED,
		)); err != nil {
			return errors.New("invalid status value")
		}
	}

	// Validate type if provided
	if f.Type != "" {
		if err := validation.Validate(f.Type, validation.In(
			DEPOSIT,
			SETTLEMENT,
			SALE,
			WITHDRAWAL,
			REFUND,
		)); err != nil {
			return errors.New("invalid type value")
		}
	}

	// Validate medium if provided
	if f.Medium != "" {
		if err := validation.Validate(f.Medium, validation.In(
			MPESA,
			TELEBIRR,
			CBE,
			CYBERSOURCE,
			KACHA,
		)); err != nil {
			return errors.New("invalid medium value")
		}
	}

	// Validate format if provided
	if f.Format != "" {
		if err := validation.Validate(f.Format, validation.In(
			FORMATPDF,
			FORMATCSV,
		)); err != nil {
			return errors.New("invalid format value")
		}
	}

	// validating TransactionId
	if f.SocialPayTransactionID != "" {

		if _, err := uuid.Parse(f.SocialPayTransactionID); err != nil {

			return errors.New("invalid socialpay_txn_id format")
		}

	}

	if f.Sort != "" {

		if err := validation.Validate(f.Sort, validation.In("ASC", "DESC")); err != nil {

			return err
		}
	}

	// validating ReferenceId
	// if f.ReferenceId != ""{

	// 	if _,err :=uuid.Parse(f.ReferenceId);err != nil{

	// 		return errors.New("invalid reference id format")
	// 	}
	// }

	return nil
}

func (fp FilterParameters) ToFilter() filter.Filter {
	var fields []filter.FilterItem

	if !fp.StartDate.IsZero() && !fp.EndDate.IsZero() {
		fields = append(fields, filter.Field{
			Name:     "created_at",
			Operator: "between",
			Value:    []interface{}{fp.StartDate, fp.EndDate},
		})
	}

	if fp.Status != "" {
		fields = append(fields, filter.Field{
			Name:     "status",
			Operator: "=",
			Value:    fp.Status,
		})
	}

	if fp.Type != "" {
		fields = append(fields, filter.Field{
			Name:     "type",
			Operator: "=",
			Value:    fp.Type,
		})
	}

	if fp.Medium != "" {
		fields = append(fields, filter.Field{
			Name:     "medium",
			Operator: "=",
			Value:    fp.Medium,
		})
	}

	if fp.SocialPayTransactionID != "" {

		fields = append(fields, &filter.Field{
			Name:     "id",
			Operator: "=",
			Value:    fp.SocialPayTransactionID,
		})
	}

	if fp.MerchantID != uuid.Nil {
		fields = append(fields, filter.Field{
			Name:     "merchant_id",
			Operator: "=",
			Value:    fp.MerchantID,
		})
	}

	// make default sort desc

	if fp.Sort == "" {
		fp.Sort = "DESC"
	}

	// search query
	var searchQueries []filter.SearchQuery
	if fp.ReferenceId != "" {
		searchQueries = append(searchQueries, filter.SearchQuery{
			Field: "reference", Term: fp.ReferenceId,
		})
	}
	if fp.ProviderTransactionId != "" {
		searchQueries = append(searchQueries, filter.SearchQuery{
			Field: "provider_tx_id", Term: fp.ProviderTransactionId,
		})
	}

	return filter.Filter{
		Group: filter.FilterGroup{
			Linker: "AND",
			Fields: fields,
		},

		Search: &filter.Search{
			Queries: searchQueries,
		},
		Sort: []filter.Sort{
			{Field: "created_at", Operator: fp.Sort},
		},
	}
}
