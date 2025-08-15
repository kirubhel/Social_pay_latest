package controller

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/socialpay/socialpay/src/pkg/bank/core"
	"github.com/socialpay/socialpay/src/pkg/merchants/errors"
)

func (controller Controller) LoadEnv() {

	envFilePath := ".env"
	err := godotenv.Overload(envFilePath)
	if err != nil {
		log.Println("Error loading .env file:", err)
	}
}

func GetSignature(password, requestId string) string {
	input := password + requestId
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (controller Controller) extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		err := errors.Error{
			Type:    "UNAUTHORIZED",
			Message: "Please provide a valid header token",
			Code:    http.StatusUnauthorized,
		}
		return "", err
	}
	return parts[1], nil
}

func (controller Controller) SendRequest(URL string, jsonData []byte, w http.ResponseWriter) {

	forwardReq, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Failed to create outbound request:", err)
		http.Error(w, "Failed to create outbound request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	forwardReq.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // only for test
			},
		},
	}
	log.Println("Sending request to external API...")
	resp, err := client.Do(forwardReq)
	if err != nil {
		log.Println("Failed to reach external API:", err)
		http.Error(w, "Failed to reach external API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	log.Println("External API responded with status:", resp.StatusCode)
	w.WriteHeader(resp.StatusCode)

	var apiResponse map[string]interface{}
	log.Println("Decoding external API response...")
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Println("Failed to decode external API response:", err)
		http.Error(w, "Failed to decode external API response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Forwarding response to client: %+v\n", apiResponse)

	if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
		log.Println("Failed to encode response to client:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (controller Controller) GetRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	// creating request id
	reqId := uuid.NewString()

	// getting env
	MerchantCode := os.Getenv("AWASH_TEST_MERCHANT_CODE")
	TinNumber := os.Getenv("AWASH_TEST_TIN_NUMBER")
	Password := os.Getenv("AWASH_TEST_PASSWORD")

	log.Println("Decoding incoming JSON body...")
	var req core.DebitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return nil, err
	}
	log.Printf("Decoded request: %+v\n", req)

	log.Println("Marshalling request to JSON for forwarding...")

	// creating request authorization
	req.Authorization.RequestID = reqId
	req.Authorization.MerchantCode = MerchantCode
	req.Authorization.MerchantTillNumber = TinNumber

	// signing the request
	requestSign := GetSignature(Password, reqId)
	req.Authorization.RequestSignature = requestSign

	//loging request body
	log.Println("request body", req)

	// marshalling the request
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Failed to marshal request: "+err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return jsonData, nil
}

func (controller Controller) TestDebitHandler(w http.ResponseWriter, r *http.Request) {

	controller.LoadEnv()
	token, err := controller.extractToken(r)
	if err != nil {
		log.Println("UnAuthorized::Attempt")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		log.Println("UnAuthorized::Attempt")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Println("LogIn::user session ", session.User.Id)
	jsonData, err := controller.GetRequestBody(w, r)

	if err != nil {
		return
	}

	baseURL := os.Getenv("AWASH_TEST_BASE_URL")
	log.Println("Read BaseUlr from env:", baseURL)
	if baseURL == "" {
		log.Println("API_BASE_URL not set in environment")
		http.Error(w, "API_BASE_URL not set", http.StatusInternalServerError)
		return
	}

	targetURL := baseURL + "/MerchantRS/DebitRequest"
	log.Println("targetUrl:", targetURL)
	controller.SendRequest(targetURL, jsonData, w)
}

func (controller Controller) TestDebitStatus(w http.ResponseWriter, r *http.Request) {

	controller.LoadEnv()

	jsonData, err := controller.GetRequestBody(w, r)

	if err != nil {
		return
	}

	baseURL := os.Getenv("AWASH_TEST_BASE_URL")
	log.Println("Read BaseUlr from env:", baseURL)
	if baseURL == "" {
		log.Println("API_BASE_URL not set in environment")
		http.Error(w, "API_BASE_URL not set", http.StatusInternalServerError)
		return
	}

	targetURL := baseURL + "/MerchantRS/DebitStatus"
	controller.SendRequest(targetURL, jsonData, w)

}

func (controller Controller) SuccessCallBack(w http.ResponseWriter, r *http.Request) {

	var payload core.AwashPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// logging the request response
	log.Println("AWASH::CallBack Payload:", payload)

	// TODO save to db later
	res := map[string]interface{}{
		"status": "received",
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		// log err
		log.Println("err::callback response encode:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
