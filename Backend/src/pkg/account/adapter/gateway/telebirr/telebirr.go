package telebirr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// obtains an access token from the Safaricom API.
func getAccessToken() (string, error) {
	username := os.Getenv("SAFARICOM_USERNAME")
	password := os.Getenv("SAFARICOM_PASSWORD")
	grantType := "client_credentials"

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	url := "https://apisandbox.safaricom.et/v1/token/generate?grant_type=" + grantType

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
