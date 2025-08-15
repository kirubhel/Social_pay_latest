package entity

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// DateUnit represents the time unit for chart data aggregation
type DateUnit string

const (
	DAY   DateUnit = "day"
	WEEK  DateUnit = "week"
	MONTH DateUnit = "month"
	YEAR  DateUnit = "year"
)

// AnalyticsFilter represents filters for transaction analytics
// @Description Parameters for filtering transaction analytics data
type AnalyticsFilter struct {
	// @Description Start date in ISO 8601 format (YYYY-MM-DD)
	// @Example 2023-01-01
	StartDate time.Time `json:"start_date" binding:"required"`

	// @Description End date in ISO 8601 format (YYYY-MM-DD)
	// @Example 2023-12-31
	EndDate time.Time `json:"end_date" binding:"required"`

	// @Description Transaction status filter
	Status []TransactionStatus `json:"status,omitempty"`

	// @Description Transaction type filter
	Type []TransactionType `json:"type,omitempty"`

	// @Description Payment medium filter
	Medium []TransactionMedium `json:"medium,omitempty"`

	// @Description Transaction source filter (QR, DIRECT, HOSTED_CHECKOUT)
	Source []TransactionSource `json:"source,omitempty"`

	// @Description QR tag filter
	QRTag []string `json:"qr_tag,omitempty"`

	// @Description Amount range filter
	AmountMin *float64 `json:"amount_min,omitempty"`
	AmountMax *float64 `json:"amount_max,omitempty"`

	// @Description Merchant ID filter (for admin analytics)
	MerchantID []string `json:"merchant_id,omitempty"`
}

// ChartFilter represents filters for chart data
// @Description Parameters for generating chart data
type ChartFilter struct {
	AnalyticsFilter

	// @Description Date unit for chart aggregation (day, week, month, year)
	DateUnit DateUnit `json:"date_unit" binding:"required"`

	// @Description Chart type: "amount" for transaction amounts, "count" for transaction counts
	ChartType string `json:"chart_type" binding:"required"`
}

// TransactionAnalytics represents aggregated transaction analytics
// @Description Aggregated transaction analytics data
type TransactionAnalytics struct {
	// Transaction counts and amounts
	TotalTransactions int64   `json:"total_transactions"`
	TotalAmount       float64 `json:"total_amount"`

	// Financial breakdown
	TotalMerchantNet float64 `json:"total_merchant_net"`

	// Transaction type breakdown
	TotalDeposits    TransactionTypeAnalytics `json:"total_deposits"`
	TotalWithdrawals TransactionTypeAnalytics `json:"total_withdrawals"`
	TotalTips        TransactionTypeAnalytics `json:"total_tips"`

	// Period comparison (percentage change from previous period)
	PeriodComparison *PeriodComparison `json:"period_comparison,omitempty"`
}

// AdminTransactionAnalytics represents admin-specific transaction analytics with additional financial data
// @Description Admin transaction analytics with VAT, fees, and admin net information
type AdminTransactionAnalytics struct {
	// Transaction counts and amounts
	TotalTransactions int64   `json:"total_transactions"`
	TotalAmount       float64 `json:"total_amount"`

	// Financial breakdown - Admin specific
	TotalMerchantNet float64 `json:"total_merchant_net"`
	TotalAdminNet    float64 `json:"total_admin_net"`
	TotalVATAmount   float64 `json:"total_vat_amount"`
	TotalFeeAmount   float64 `json:"total_fee_amount"`
	TotalCustomerNet float64 `json:"total_customer_net"`

	// Transaction type breakdown
	TotalDeposits    TransactionTypeAnalytics `json:"total_deposits"`
	TotalWithdrawals TransactionTypeAnalytics `json:"total_withdrawals"`
	TotalTips        TransactionTypeAnalytics `json:"total_tips"`

	// Period comparison (percentage change from previous period)
	PeriodComparison *PeriodComparison `json:"period_comparison,omitempty"`
}

// MerchantGrowthAnalytics represents merchant growth statistics
// @Description Merchant growth analytics and statistics
type MerchantGrowthAnalytics struct {
	// Current counts
	TotalMerchants    int64 `json:"total_merchants"`
	ActiveMerchants   int64 `json:"active_merchants"`
	InactiveMerchants int64 `json:"inactive_merchants"`
	PendingMerchants  int64 `json:"pending_merchants"`

	// Growth rates (percentage change from previous period)
	TotalMerchantGrowth  float64 `json:"total_merchant_growth"`
	ActiveMerchantGrowth float64 `json:"active_merchant_growth"`

	// Period breakdown
	NewMerchantsThisPeriod         int64 `json:"new_merchants_this_period"`
	ActivatedMerchantsThisPeriod   int64 `json:"activated_merchants_this_period"`
	DeactivatedMerchantsThisPeriod int64 `json:"deactivated_merchants_this_period"`

	// Time series data for charts
	GrowthData []MerchantGrowthDataPoint `json:"growth_data,omitempty"`
}

// MerchantGrowthDataPoint represents a single data point in merchant growth chart
type MerchantGrowthDataPoint struct {
	Date            time.Time `json:"date"`
	TotalMerchants  int64     `json:"total_merchants"`
	NewMerchants    int64     `json:"new_merchants"`
	ActiveMerchants int64     `json:"active_merchants"`
	Label           string    `json:"label"`
}

// TransactionTypeAnalytics represents analytics for a specific transaction type
type TransactionTypeAnalytics struct {
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
}

// PeriodComparison represents percentage changes from the previous period
type PeriodComparison struct {
	TransactionCountChange float64 `json:"transaction_count_change"`
	AmountChange           float64 `json:"amount_change"`
	MerchantNetChange      float64 `json:"merchant_net_change"`
}

// ChartDataPoint represents a single data point in a chart
type ChartDataPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
	Label string    `json:"label"`
}

// ChartData represents chart data response
// @Description Chart data for transaction analytics
type ChartData struct {
	ChartType string           `json:"chart_type"`
	DateUnit  DateUnit         `json:"date_unit"`
	Data      []ChartDataPoint `json:"data"`
	Summary   ChartSummary     `json:"summary"`
}

// ChartSummary provides summary statistics for the chart
type ChartSummary struct {
	TotalValue   float64 `json:"total_value"`
	AverageValue float64 `json:"average_value"`
	MaxValue     float64 `json:"max_value"`
	MinValue     float64 `json:"min_value"`
	DataPoints   int     `json:"data_points"`
}

// Validate validates the analytics filter parameters
func (f *AnalyticsFilter) Validate() error {
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

	// Validate amount range
	if f.AmountMin != nil && f.AmountMax != nil && *f.AmountMin > *f.AmountMax {
		return errors.New("amount_min cannot be greater than amount_max")
	}

	return nil
}

// Validate validates the chart filter parameters
func (f *ChartFilter) Validate() error {
	// First validate the base analytics filter
	if err := f.AnalyticsFilter.Validate(); err != nil {
		return err
	}

	// Validate date unit
	if err := validation.Validate(f.DateUnit, validation.In(DAY, WEEK, MONTH, YEAR)); err != nil {
		return errors.New("invalid date_unit value. Must be one of: day, week, month, year")
	}

	// Validate chart type
	if err := validation.Validate(f.ChartType, validation.In("amount", "count")); err != nil {
		return errors.New("invalid chart_type value. Must be either 'amount' or 'count'")
	}

	return nil
}

// GetPreviousPeriod calculates the previous period for comparison
func (f *AnalyticsFilter) GetPreviousPeriod() (time.Time, time.Time) {
	duration := f.EndDate.Sub(f.StartDate)
	previousEnd := f.StartDate.Add(-time.Second)
	previousStart := previousEnd.Add(-duration)
	return previousStart, previousEnd
}
