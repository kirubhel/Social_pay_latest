package entity

import (
	"github.com/google/uuid"
)

type CommissionSettings struct {
	Percent float64 `json:"percent"`
	Cent    float64 `json:"cent"`
}

type MerchantCommission struct {
	MerchantID        uuid.UUID `json:"merchant_id"`
	CommissionActive  bool      `json:"commission_active"`
	CommissionPercent *float64  `json:"commission_percent,omitempty"`
	CommissionCent    *float64  `json:"commission_cent,omitempty"`
}

// GetFloat64OrDefault returns the value of the float64 pointer if it's not nil,
// otherwise returns the default value
func GetFloat64OrDefault(f *float64, defaultValue float64) float64 {
	if f == nil {
		return defaultValue
	}
	return *f
}
