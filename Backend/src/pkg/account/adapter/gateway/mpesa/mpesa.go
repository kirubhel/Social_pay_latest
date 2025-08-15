package mpesa

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

type TransactionType string

const (
	C2B TransactionType = "CustomerToBusiness"
	B2C TransactionType = "BusinessToCustomer"
	B2B TransactionType = "BusinessToBusiness"
)

type USSDPushRequest struct {
	MerchantRequestID string          `json:"MerchantRequestID"`
	BusinessShortCode string          `json:"BusinessShortCode"`
	Password          string          `json:"Password"`
	Timestamp         string          `json:"Timestamp"`
	TransactionType   string          `json:"TransactionType"`
	Amount            float64         `json:"Amount"`
	PartyA            string          `json:"PartyA"`
	PartyB            string          `json:"PartyB"`
	PhoneNumber       string          `json:"PhoneNumber"`
	CallBackURL       string          `json:"CallBackURL"`
	AccountReference  string          `json:"AccountReference"`
	TransactionDesc   string          `json:"TransactionDesc"`
	ReferenceData     []ReferenceData `json:"ReferenceData"`
	MerchantName      string          `json:"MerchantName"`
}

type B2CPaymentRequest struct {
	OriginatorConversationID string  `json:"OriginatorConversationID"`
	InitiatorName            string  `json:"InitiatorName"`
	SecurityCredential       string  `json:"SecurityCredential"`
	CommandID                string  `json:"CommandID"`
	PartyA                   string  `json:"PartyA"`
	PartyB                   string  `json:"PartyB"`
	Amount                   float64 `json:"Amount"`
	Remarks                  string  `json:"Remarks"`
	Occasion                 string  `json:"Occassion"`
	QueueTimeOutURL          string  `json:"QueueTimeOutURL"`
	ResultURL                string  `json:"ResultURL"`
}

type ReferenceData struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func generatePassword(businessShortCode, passkey, timestamp string) string {
	rawPassword := businessShortCode + passkey + timestamp
	log.Printf("Raw Password (before hashing): %s", rawPassword)
	hashed := sha256.Sum256([]byte(rawPassword))
	hashedHex := hex.EncodeToString(hashed[:])
	log.Printf("Hashed Password: %s", hashedHex)
	password := base64.StdEncoding.EncodeToString([]byte(hashedHex))
	log.Printf("Generated Password (after base64 encoding): %s", password)

	return password
}

// processes the STK push request and returns the full response
func HandleSTKPushRequest(req USSDPushRequest) (map[string]interface{}, error) {
	timestamp := time.Now().Format("20060102150405")
	passkey := "141e624db2261261ac66fb74edecb5657aded57ad5358476605c7f6d7e199145"
	password := generatePassword(req.BusinessShortCode, passkey, timestamp)

	formattedAmount := formatAmount(req.Amount)
	reqWithCredentials := map[string]interface{}{
		"MerchantRequestID": req.MerchantRequestID,
		"BusinessShortCode": req.BusinessShortCode,
		"Password":          password,
		"Timestamp":         timestamp,
		"TransactionType":   req.TransactionType,
		"Amount":            formattedAmount,
		"PartyA":            req.PartyA,
		"PartyB":            req.PartyB,
		"PhoneNumber":       req.PhoneNumber,
		"CallBackURL":       req.CallBackURL,
		"AccountReference":  req.AccountReference,
		"TransactionDesc":   req.TransactionDesc,
		"ReferenceData":     req.ReferenceData,
	}
	// Marshal the request payload to JSON
	jsonData, err := json.MarshalIndent(reqWithCredentials, "", "  ")
	if err != nil {
		log.Printf("[ERROR] 🛑 Failed to marshal STK push request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("[INFO] 📤 Sending STK Push Request:\n%s", jsonData)

	// Create the HTTP request
	request, err := http.NewRequest("POST", "https://api.safaricom.et/mpesa/stkpush/v3/processrequest", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[ERROR] 🛑 Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Obtain the access token
	accessToken, err := getAccessToken()
	if err != nil {
		log.Printf("[ERROR] 🛑 Error getting access token: %v", err)
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Set the request headers
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("[ERROR] 🛑 Error sending STK push request: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] 🛑 Error reading response body %v", err)
		return nil, fmt.Errorf("failed to read response body %w", err)
	}

	log.Printf("[INFO] 📩 Response Body:\n%s", body)

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("[ERROR] 🛑 Failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] 🛑 STK push request failed: Status Code: %d, Response: %s", resp.StatusCode, body)
		return response, fmt.Errorf("STK push request failed: Status Code: %d", resp.StatusCode)
	}

	log.Printf("[SUCCESS] ✅ STK Push Request successful. Status Code: %d", resp.StatusCode)
	return response, nil
}

func HandleB2CPaymentRequest(req B2CPaymentRequest) error {
	req.Amount = formatAmount(req.Amount)
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("Problem marshalling B2C payment request: %v", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Sending B2C Payment Request: %s", jsonData)
	request, err := http.NewRequest("POST", "https://api.safaricom.et/mpesa/b2c/v2/paymentrequest", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Problem creating new request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	accessToken, err := getAccessToken()
	if err != nil {
		log.Printf("Error getting access token: %v", err)
		return fmt.Errorf("failed to get access token: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending B2C payment request: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Response body: %s", body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send B2C payment request: %s, Status Code: %d", body, resp.StatusCode)
		return fmt.Errorf("B2C payment request failed: Status Code: %d, response: %s", resp.StatusCode, body)
	}

	log.Printf("B2C payment request successful: Status Code: %d", resp.StatusCode)
	return nil
}

func getAccessToken() (string, error) {
	username := os.Getenv("SAFARICOM_USERNAME")
	password := os.Getenv("SAFARICOM_PASSWORD")
	grantType := "client_credentials"
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	url := "https://api.safaricom.et/v1/token/generate?grant_type=" + grantType

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Response body: %s", body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s, response: %s", resp.Status, body)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   string `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func formatAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}
