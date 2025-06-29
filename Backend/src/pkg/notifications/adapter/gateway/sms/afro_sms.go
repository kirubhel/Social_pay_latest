package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

// AfroSMS implements the SMSProvider interface using KMI Cloud API
type AfroSMS struct {
	log       *log.Logger
	Name      string
	AccessKey string
	SecretKey string
	From      string
	URL       string
}

// New creates a new AfroSMS provider instance
func New(log *log.Logger) *AfroSMS {
	accessKey := os.Getenv("KMI_ACCESS_KEY")
	secretKey := os.Getenv("KMI_SECRET_KEY")
	from := os.Getenv("KMI_SMS_FROM")
	url := os.Getenv("KMI_SMS_URL")

	// Validate required environment variables
	if accessKey == "" || secretKey == "" || from == "" || url == "" {
		log.Printf("[AfroSMS] Warning: Missing required environment variables for SMS service")
		log.Printf("[AfroSMS] Required: KMI_ACCESS_KEY, KMI_SECRET_KEY, KMI_SMS_FROM, KMI_SMS_URL")
	}

	return &AfroSMS{
		log:       log,
		Name:      "Social Pay",
		AccessKey: accessKey,
		SecretKey: secretKey,
		From:      from,
		URL:       url,
	}
}

// SendSMS sends an SMS using the KMI Cloud API
func (sms *AfroSMS) SendSMS(phoneNumber, message string) error {
	sms.log.Printf("[AfroSMS] Send SMS via KMI Cloud")

	// Format phone number to use proper country code
	formattedPhone := sms.formatPhoneNumber(phoneNumber)
	if formattedPhone == "" {
		errMsg := "invalid phone number format"
		sms.log.Printf("[AfroSMS] %s", errMsg)
		return errors.New(errMsg)
	}

	// Define the request body structure
	type SMSBody struct {
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
		From      string `json:"from,omitempty"`
		To        string `json:"to"`
		Message   string `json:"message"`
	}

	sms.log.Printf("[AfroSMS] Recipient: %s", formattedPhone)
	sms.log.Printf("[AfroSMS] Message: %s", message)

	sms.log.Printf("[AfroSMS] AccessKey: %s", sms.AccessKey)
	sms.log.Printf("[AfroSMS] SecretKey: %s", sms.SecretKey)
	sms.log.Printf("[AfroSMS] From: %s", sms.From)
	sms.log.Printf("[AfroSMS] To: %s", formattedPhone)
	sms.log.Printf("[AfroSMS] Message: %s", message)

	// Create the request body
	body := SMSBody{
		AccessKey: sms.AccessKey,
		SecretKey: sms.SecretKey,
		From:      sms.From,
		To:        formattedPhone,
		Message:   message,
	}

	// Marshal the request body to JSON
	serBody, err := json.Marshal(body)
	if err != nil {
		sms.log.Printf("[AfroSMS] Error marshaling request: %v", err)
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, sms.URL, bytes.NewReader(serBody))
	if err != nil {
		sms.log.Printf("[AfroSMS] Error creating request: %v", err)
		return err
	}

	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Send the request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		sms.log.Printf("[AfroSMS] Error sending request: %v", err)
		return err
	}
	defer res.Body.Close()

	// Define the expected API response structure
	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Code    int    `json:"code"`
		Result  struct {
			To    string `json:"to"`
			SmsID string `json:"smsId"`
		} `json:"result"`
	}

	// Decode the response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		sms.log.Printf("[AfroSMS] Error decoding response: %v", err)
		return err
	}

	// Check if the API request was successful
	if !response.Success {
		sms.log.Printf("[AfroSMS] API Error: %s", response.Message)
		return errors.New(response.Message)
	}

	sms.log.Printf("[AfroSMS] SMS sent successfully. ID: %s", response.Result.SmsID)
	return nil
}

func (sms AfroSMS) formatPhoneNumber(phone string) string {
	cleaned := strings.ReplaceAll(phone, "+", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	if strings.HasPrefix(cleaned, "251") && len(cleaned) > 3 {
		return "00251" + cleaned[3:]
	}

	return cleaned
}
