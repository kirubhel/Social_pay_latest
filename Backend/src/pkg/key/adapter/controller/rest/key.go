package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	merchant "github.com/socialpay/socialpay/src/pkg/merchants/core/entity"
)

// Response
type AuthResponse struct {
	APIKey *struct {
		Token      string `json:"token"`
		Active     string `json:"active"`
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
	} `json:"apikey,omitempty"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (controller Controller) CreateKey(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [GetSignIn] ")
	// Request
	type Request struct {
		MerchantID       string    `json:"merchant_id" example:"123456"`
		Service          string    `json:"service" example:"hosted checkout"`
		ExpiresAt        time.Time `json:"expiry_date" example:"2024-12-04"`
		CommissionFromMe bool      `json:"commission_from_me" example:"true"`
		Store            string    `json:"store" example:"store 1"`
		LicenseNumber    string    `json:"license_number" example:"AA/1313/2014"`
	}
	var token string

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide a valid header token",
			},
		}, http.StatusUnauthorized)
		return
	}

	token = strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}
	fmt.Println("session")
	fmt.Println(session.User.Id)
	//request the endpoint and get merchant id
	host := os.Getenv("HOST")
	if host == "" {
		host = "http://196.190.251.194:8082" // Default value if the environment variable is not set
	}
	reqMerchant, err := http.NewRequest("GET", host+"/merchant-by-user", nil)
	if err != nil {
		controller.log.Println("Error creating request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to create request to merchant service",
			},
		}, http.StatusInternalServerError)
		return
	}
	// Pass the token
	reqMerchant.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(reqMerchant)
	if err != nil {
		controller.log.Println("Error making request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to retrieve merchant information",
			},
		}, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		controller.log.Println("Error: received non-200 response code")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to retrieve merchant information",
			},
		}, http.StatusInternalServerError)
		return
	}
	var merchantResp struct {
		Success bool              `json:"success"`
		Data    merchant.Merchant `json:"data"`
	}
	fmt.Println("response", response.Body)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&merchantResp)
	if err != nil {
		controller.log.Println("Error decoding response:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to decode merchant information",
			},
		}, http.StatusInternalServerError)
		return
	}
	fmt.Println("response", response.Body)

	merchant := merchantResp.Data

	var req Request
	// Parse request
	defer r.Body.Close()
	decoder = json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, AuthResponse{
			Error: &struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	//merchantID, service, expiryDate, store string, commissionFrom bool
	req.MerchantID = merchant.MerchantID
	apiKey, err := controller.interactor.CreateAPIKey(req.MerchantID, req.Service, req.ExpiresAt, req.Store, req.CommissionFromMe)
	if err != nil {
		SendJSONResponse(w, AuthResponse{
			Error: &struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{
				Type:    err.Error(),
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	SendJSONResponse(w, AuthResponse{
		APIKey: &struct {
			Token      string `json:"token"`
			Active     string `json:"active"`
			PrivateKey string `json:"private_key"`
			PublicKey  string `json:"public_key"`
		}{

			PrivateKey: apiKey.PrivateKey,
			PublicKey:  apiKey.PublicKey,
		},
	}, http.StatusOK)
}
func (controller Controller) GetApiKeyByToken(w http.ResponseWriter, r *http.Request) {
	// 1. Get token from query
	token := r.URL.Query().Get("token")

	// 2. Check if token is missing
	if token == "" {
		SendJSONResponse(w, AuthResponse{
			Error: &struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{
				Type:    "INVALID_REQUEST",
				Message: "Token is required",
			},
		}, http.StatusBadRequest)
		return
	}
	fmt.Println("Token from request:", token)

	// 3. Validate API token
	apiKey, err := controller.interactor.ValidateAPIToken(token)
	if err != nil {
		fmt.Println("Error validating API token:", err)
		SendJSONResponse(w, AuthResponse{
			Error: &struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{
				Type:    "NOT_FOUND",
				Message: "Invalid API Token",
			},
		}, http.StatusNotFound)
		return
	}

	// 4. VERY IMPORTANT: Check apiKey is not nil
	if apiKey == nil {
		fmt.Println("API Key is nil!")
		SendJSONResponse(w, AuthResponse{
			Error: &struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{
				Type:    "NOT_FOUND",
				Message: "API Key not found",
			},
		}, http.StatusNotFound)
		return
	}

	// 5. Now it's safe to use apiKey

	SendJSONResponse(w, AuthResponse{
		APIKey: &struct {
			Token      string `json:"token"`
			Active     string `json:"active"`
			PrivateKey string `json:"private_key"`
			PublicKey  string `json:"public_key"`
		}{
			Token:      apiKey.APIKey,
			PublicKey:  apiKey.PublicKey,
			PrivateKey: apiKey.PrivateKey,
		},
	}, http.StatusOK)
}
