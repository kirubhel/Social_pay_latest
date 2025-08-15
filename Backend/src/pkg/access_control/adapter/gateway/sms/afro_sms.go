package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AfroSMS struct {
	log   *log.Logger
	Name  string
	Token string
	URL   string
}

func New(log *log.Logger) AfroSMS {
	return AfroSMS{
		log:   log,
		Name:  "Afro SMS",
		Token: "eyJhbGciOiJIUzI1NiJ9.eyJpZGVudGlmaWVyIjoicUdFc1VjQVN2WTU3SDB5Vm5jMlVVWnJ0S2FKRUxFVW8iLCJleHAiOjE4NDI1OTcxNzQsImlhdCI6MTY4NDc0NDM3NCwianRpIjoiY2I3MzFhYzEtNWNjOC00YTRkLTg3NTEtMjMxMzc1ZTIwNWM3In0.gka4m6qu_Wx6sNdDHWzggcmxPWAY_gG4kFj2kUfcJPo",
		URL:   "https://api.afromessage.com/api/send",
	}
}

// SendSMS sends an SMS using the AfroMessage API.
func (sms AfroSMS) SendSMS(phone, message string) error {
	sms.log.Println("[SendSMS] Sending SMS to:", phone)
	sms.log.Println("[SendSMS] Message:", message)

	// Define the request body structure
	type SMSBody struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Message string `json:"message"`
	}

	// Create the request body
	body := SMSBody{
		From:    "e80ad9d8-adf3-463f-80f4-7c4b39f7f164", // Replace with your sender ID
		To:      phone,
		Message: message,
	}

	// Marshal the request body to JSON
	serBody, err := json.Marshal(body)
	if err != nil {
		sms.log.Println("[SendSMS] Failed to marshal request body:", err)
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	sms.log.Println("[SendSMS] Request body:", string(serBody))

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, sms.URL, bytes.NewReader(serBody))
	if err != nil {
		sms.log.Println("[SendSMS] Failed to create request:", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+sms.Token)

	// Send the request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		sms.log.Println("[SendSMS] Failed to send request:", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	// Log the raw response body
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		sms.log.Println("[SendSMS] Failed to read response body:", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}
	sms.log.Println("[SendSMS] Raw API response:", string(rawBody))

	// Define the expected API response structure
	type APIResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	// Attempt to decode the response as JSON
	var apiResponse APIResponse
	err = json.Unmarshal(rawBody, &apiResponse)
	if err != nil {
		sms.log.Println("[SendSMS] Failed to decode JSON response:", err)
		return fmt.Errorf("invalid API response: %s", string(rawBody))
	}

	sms.log.Println("[SendSMS] Decoded API response:", apiResponse)

	// Check if the API request was successful
	if res.StatusCode != http.StatusOK || apiResponse.Status != "success" {
		sms.log.Println("[SendSMS] API request failed with status code:", res.StatusCode)
		return fmt.Errorf("API request failed: %s", apiResponse.Message)
	}

	sms.log.Println("[SendSMS] SMS sent successfully")
	return nil
}
