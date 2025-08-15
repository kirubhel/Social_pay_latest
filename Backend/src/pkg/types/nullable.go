// pkg/types/nullable.go
package types

import (
	"database/sql"
	"encoding/json"
)

// NullString is a swagger-friendly nullable string
// @Schema
type NullString struct {
	String string `json:"string" example:"sample description"`
	Valid  bool   `json:"valid" example:"true"`
}

// ToNullString creates a NullString from regular string
func ToNullString(s string) NullString {
	return NullString{
		String: s,
		Valid:  s != "",
	}
}

// FromSqlNullString converts sql.NullString to our type
func FromSqlNullString(ns sql.NullString) NullString {
	return NullString{
		String: ns.String,
		Valid:  ns.Valid,
	}
}

// ToSqlNullString converts back to sql.NullString
func (ns NullString) ToSqlNullString() sql.NullString {
	return sql.NullString{
		String: ns.String,
		Valid:  ns.Valid,
	}
}

// MarshalJSON implements json.Marshaler
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements json.Unmarshaler
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.String)
}
