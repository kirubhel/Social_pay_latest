package filter

import (
	"reflect"
	"testing"

	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
)

func TestFilterBuild(t *testing.T) {
	tests := []struct {
		name      string
		filter    Filter
		wantSQL   string
		wantArgs  []interface{}
		expectErr bool
	}{
		{
			name: "simple AND fields",
			filter: Filter{
				Group: FilterGroup{
					Linker: "AND",
					Fields: []FilterItem{
						Field{Name: "age", Operator: ">", Value: 18},
						Field{Name: "status", Operator: "=", Value: "active"},
					},
				},
			},
			wantSQL:  "WHERE (age > $1) AND (status = $2);",
			wantArgs: []interface{}{18, "active"},
		},
		{
			name: "nested OR group",
			filter: Filter{
				Group: FilterGroup{
					Linker: "AND",
					Fields: []FilterItem{
						Field{Name: "age", Operator: ">", Value: 18},
						FilterGroup{
							Linker: "OR",
							Fields: []FilterItem{
								Field{Name: "name", Operator: "LIKE", Value: "%john%"},
								Field{Name: "email", Operator: "LIKE", Value: "%john%"},
							},
						},
					},
				},
			},
			wantSQL:  "WHERE (age > $1) AND ((name LIKE $2) OR (email LIKE $3));",
			wantArgs: []interface{}{18, "%john%", "%john%"},
		},
		{
			name: "IS NULL operator",
			filter: Filter{
				Group: FilterGroup{
					Linker: "AND",
					Fields: []FilterItem{
						Field{Name: "deleted_at", Operator: "IS NULL", Value: nil},
					},
				},
			},
			wantSQL:  "WHERE (deleted_at IS NULL);",
			wantArgs: nil,
		},
		{
			name: "empty group",
			filter: Filter{
				Group: FilterGroup{
					Linker: "AND",
					Fields: []FilterItem{},
				},
			},
			wantSQL:  "WHERE 1=1",
			wantArgs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotArgs, err := tt.filter.Build()
			if (err != nil) != tt.expectErr {
				t.Fatalf("Build() error = %v, expectErr %v", err, tt.expectErr)
			}
			if gotSQL != tt.wantSQL {
				t.Errorf("Build() gotSQL = %q, wantSQL %q", gotSQL, tt.wantSQL)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("Build() gotArgs = %v, wantArgs %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestFilterBuild_Advanced(t *testing.T) {
	filter := Filter{
		Group: FilterGroup{
			Linker: "AND",
			Fields: []FilterItem{
				// age > 21
				Field{Name: "age", Operator: ">", Value: 21},

				// nested OR group: (name LIKE '%Alice%' OR email LIKE '%example.com')
				FilterGroup{
					Linker: "OR",
					Fields: []FilterItem{
						Field{Name: "name", Operator: "LIKE", Value: "%Alice%"},
						Field{Name: "email", Operator: "LIKE", Value: "%example.com"},
					},
				},

				// nested AND group: (status = 'active' AND deleted_at IS NOT NULL)
				FilterGroup{
					Linker: "AND",
					Fields: []FilterItem{
						Field{Name: "status", Operator: "=", Value: "active"},
						Field{Name: "deleted_at", Operator: "IS NOT NULL", Value: nil},
					},
				},

				// simple field: score >= 70
				Field{Name: "score", Operator: ">=", Value: 70},
			},
		},
	}

	wantSQL := "WHERE (age > $1) AND ((name LIKE $2) OR (email LIKE $3)) AND ((status = $4) AND (deleted_at IS NOT NULL)) AND (score >= $5);"
	wantArgs := []interface{}{21, "%Alice%", "%example.com", "active", 70}

	gotSQL, gotArgs, err := filter.Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	if gotSQL != wantSQL {
		t.Errorf("Build() gotSQL = %q, wantSQL %q", gotSQL, wantSQL)
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("Build() gotArgs = %v, wantArgs %v", gotArgs, wantArgs)
	}
}


func TestFilterBuild_WithSortAndPagination_PagePerPage(t *testing.T) {
	filter := Filter{
		Group: FilterGroup{
			Linker: "AND",
			Fields: []FilterItem{
				Field{Name: "age", Operator: ">", Value: 30},
				Field{Name: "status", Operator: "=", Value: "active"},
			},
		},
		Sort: []Sort{
			{Field: "created_at", Operator: "DESC"},
		
		},
		Pagination: pagination.Pagination{
			Page:    3,
			PerPage: 5,
		},
	}

	wantSQL := "WHERE (age > $1) AND (status = $2) ORDER BY created_at DESC LIMIT 5 OFFSET 10;"
	wantArgs := []interface{}{30, "active"}

	gotSQL, gotArgs, err := filter.Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	if gotSQL != wantSQL {
		t.Errorf("Build() gotSQL = %q, wantSQL %q", gotSQL, wantSQL)
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("Build() gotArgs = %v, wantArgs %v", gotArgs, wantArgs)
	}
}

func TestFilterBuild_WithSearch(t *testing.T) {
	filter := Filter{
		Group: FilterGroup{
			Linker: "AND",
			Fields: []FilterItem{
				Field{Name: "status", Operator: "=", Value: "active"},
			},
		},
		Search: &Search{
			Queries: []SearchQuery{
				{Field: "name", Term: "john"},
				{Field: "email", Term: "john"},
			},
		},
	}

	wantSQL := "WHERE (status = $1) AND ((name ILIKE $2) OR (email ILIKE $3));"
	wantArgs := []interface{}{"active", "%john%", "%john%"}

	gotSQL, gotArgs, err := filter.Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	if gotSQL != wantSQL {
		t.Errorf("Build() gotSQL = %q, wantSQL %q", gotSQL, wantSQL)
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("Build() gotArgs = %v, wantArgs %v", gotArgs, wantArgs)
	}
}


