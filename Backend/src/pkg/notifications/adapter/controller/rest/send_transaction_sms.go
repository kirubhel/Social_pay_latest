package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (controller Controller) SendTransactionSMS(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		PhoneNumber  string  `json:"phone_number"`
		CustomerName string  `json:"customer_name"`
		MerchantName string  `json:"merchant_name"`
		Amount       float64 `json:"amount"`
		Currency     string  `json:"currency"`
		ReferenceID  string  `json:"reference_id"`
	}

	var req Request
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request payload",
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate and normalize phone number
	normalizedPhone, err := normalizeEthiopianPhoneNumber(req.PhoneNumber)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_PHONE_NUMBER",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Generate current timestamp in Addis Ababa timezone
	loc, _ := time.LoadLocation("Africa/Addis_Ababa")
	now := time.Now().In(loc)
	dateStr := now.Format("02 Jan 2006")
	timeStr := now.Format("3:04 PM")

	// Clean SMS message
	message := fmt.Sprintf(
		`Dear %s,
Payment of %.2f %s to %s successful.
Reference: %s
Date: %s at %s

Winners choose and use Social Pay!`,
		req.CustomerName,
		req.Amount,
		req.Currency,
		req.MerchantName,
		req.ReferenceID,
		dateStr,
		timeStr,
	)

	// Send SMS
	if err := controller.sms.SendSMS(normalizedPhone, message); err != nil {
		controller.log.Printf("Failed to send transaction SMS: %v", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "SMS_SEND_FAILED",
				Message: "Failed to send transaction notification",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Success response
	SendJSONResponse(w, Response{
		Success: true,
		Data: map[string]string{
			"message":          "Transaction SMS sent successfully",
			"normalized_phone": normalizedPhone,
			"timestamp":        now.Format(time.RFC3339),
		},
	}, http.StatusOK)
}

func formatDate(timestamp string) string {
	t, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		return timestamp
	}
	return t.Format("02 Jan 2006")
}

func formatTime(timestamp string) string {
	t, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		return timestamp
	}
	return t.Format("3:04 PM")
}

// accepts various formats and returns 251xxxxxxxxx
func normalizeEthiopianPhoneNumber(phone string) (string, error) {
	// Remove all non-digit characters
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	// Handle different formats
	switch {
	case len(cleaned) == 9 && cleaned[0] == '9': // 9xxxxxxxx
		return "251" + cleaned, nil
	case len(cleaned) == 10 && cleaned[0] == '0': // 09xxxxxxxx
		return "251" + cleaned[1:], nil
	case len(cleaned) == 12 && strings.HasPrefix(cleaned, "251"): // 2519xxxxxxxx
		return cleaned, nil
	case len(cleaned) == 13 && strings.HasPrefix(cleaned, "+251"): // +2519xxxxxxxx
		return cleaned[1:], nil
	default:
		return "", fmt.Errorf("invalid Ethiopian phone number format. Accepted formats: 9xxxxxxxx, 09xxxxxxxx, 2519xxxxxxxx, +2519xxxxxxxx")
	}
}
