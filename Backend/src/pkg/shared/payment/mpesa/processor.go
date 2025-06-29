package mpesa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

type processor struct {
	username    string
	password    string
	isTestMode  bool
	baseURL     string
	callbackURL string
	log         logging.Logger
}

// ProcessorConfig holds the configuration for M-PESA processor
type ProcessorConfig struct {
	Username    string
	Password    string
	IsTestMode  bool
	BaseURL     string
	CallbackURL string
}

// NewProcessor creates a new M-PESA payment processor
func NewProcessor(config ProcessorConfig) payment.Processor {
	if config.BaseURL == "" {
		config.BaseURL = os.Getenv("MPESA_BASE_URL")
	}

	return &processor{
		username:    config.Username,
		password:    config.Password,
		isTestMode:  config.IsTestMode,
		baseURL:     config.BaseURL,
		callbackURL: config.CallbackURL,
		log:         logging.NewStdLogger("[MPESA] [PROCESSOR]"),
	}
}

func (p *processor) InitiatePayment(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Info("Initiating M-PESA payment", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"amount":         req.Amount,
		"currency":       req.Currency,
	})

	// Prepare M-PESA specific request
	mpesaReq := map[string]interface{}{
		"MerchantRequestID": req.TransactionID.String(),
		"AccountReference":  req.TransactionID.String(),
		"Amount":            req.Amount,
		"BusinessShortCode": "1883",
		"CallBackURL":       os.Getenv("MPESA_CALLBACK_BASE_URL"),
		"Password":          p.password,
		"PartyA":            req.PhoneNumber,
		"PartyB":            "1883",
		"PhoneNumber":       req.PhoneNumber,
		"TransactionDesc":   req.Description,
		"Timestamp":         time.Now().Format("20060102150405"),
		"TransactionType":   "CustomerPayBillOnline",
	}

	p.log.Info("Prepared M-PESA request", map[string]interface{}{
		"reference": mpesaReq["AccountReference"],
		"phone":     mpesaReq["PhoneNumber"],
	})

	// Convert request to JSON
	jsonData, err := json.Marshal(mpesaReq)
	if err != nil {
		p.log.Error("Failed to marshal request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	p.log.Info("Creating request to M-PESA", map[string]interface{}{
		"url":  fmt.Sprintf("%s/account/mpesa/ussd-push", p.baseURL),
		"body": string(jsonData),
	})
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/account/mpesa/ussd-push", p.baseURL),
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		log.Printf("[MPESA] Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", apikey)
	log.Printf("[MPESA] Added headers - Content-Type: application/json, X-API-KEY: %s", maskAPIKey(apikey))
	httpReq.Header.Set("Authorization", "Basic "+p.username)

	p.log.Info("Sending request to M-PESA", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	// Send request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.log.Error("Failed to send request to M-PESA", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	p.log.Info("M-PESA response status", map[string]interface{}{
		"status": resp.Status,
	})
	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		p.log.Error("Failed to read response body", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	p.log.Info("Response body", map[string]interface{}{
		"body": string(bodyBytes),
	})

	// Parse response
	var mpesaResp struct {
		Success bool `json:"success"`
		Data    struct {
			CheckoutRequestID   string `json:"CheckoutRequestID"`
			CustomerMessage     string `json:"CustomerMessage"`
			MerchantRequestID   string `json:"MerchantRequestID"`
			ResponseCode        string `json:"ResponseCode"`
			ResponseDescription string `json:"ResponseDescription"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &mpesaResp); err != nil {
		p.log.Error("Failed to decode M-PESA response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Map response status
	status := txEntity.PENDING
	if !mpesaResp.Success {
		status = txEntity.FAILED
	}

	return &payment.PaymentResponse{
		Success:       mpesaResp.Success,
		TransactionID: req.TransactionID,
		Status:        status,
		ProcessorRef:  req.Reference,
		Message:       mpesaResp.Data.CustomerMessage,
	}, nil
}

func (p *processor) SettlePayment(ctx context.Context, req *payment.CallbackRequest) error {
	p.log.Info("Processing M-PESA callback", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"status":         req.Status,
	})

	// Extract M-PESA-specific data from metadata
	message, _ := req.Metadata["message"].(string)

	p.log.Info("M-PESA callback details", map[string]interface{}{
		"message": message,
	})

	if req.Status == txEntity.SUCCESS {
		return nil
	}

	return fmt.Errorf("payment failed: %s", message)
}

func (p *processor) GetType() txEntity.TransactionMedium {
	return txEntity.MPESA
}

func (p *processor) InitiateWithdrawal(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Error("Withdrawal not supported", map[string]interface{}{
		"processor": "M-PESA",
	})
	return nil, fmt.Errorf("withdrawal not supported for M-PESA")
}

func maskAPIKey(key string) string {
	if len(key) <= 4 {
		return key
	}
	return "****" + key[len(key)-4:]
}
