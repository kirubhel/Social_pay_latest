package filter

import (
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
)

type Filter struct {
	Pagination pagination.Pagination 		`json:"pagination,omitempty"`
	Sort       []Sort             			`json:"sort"`
	Group      FilterGroup           		`json:"group"`
	Search 		*Search   					`json:"search,omitempty"`
}

type FilterGroup struct {
	Linker string       `json:"linker"` // AND or OR
	Fields []FilterItem `json:"fields"` // Fields or nested groups
}

type Field struct {
	Name     string      `json:"name"`
	Operator string      `json:"operator"` // =, >, <, LIKE, NOT LIKE, IS NULL, IS NOT NULL
	Value    interface{} `json:"value"`
}
type Search struct {
	Queries []SearchQuery
}

type SearchQuery struct {
	Term string `json:"term"` // The search keywords 
	Field string `json:"fields"` // The DB fields to apply the search on 
}


type Sort struct {
	// Field name to sort with
	Field string `json:"field"`
	// ASC DESC
	Operator string 	`json:"operator"`	
}
