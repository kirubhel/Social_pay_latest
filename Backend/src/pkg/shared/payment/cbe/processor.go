package cbe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

type processor struct {
	merchantID    string
	merchantKey   string
	terminalID    string
	credentialKey string
	isTestMode    bool
	baseURL       string
	callbackURL   string
	log           logging.Logger
}

// ProcessorConfig holds the configuration for CBE processor
type ProcessorConfig struct {
	MerchantID    string
	MerchantKey   string
	TerminalID    string
	CredentialKey string
	IsTestMode    bool
	BaseURL       string
	CallbackURL   string
}

type CBEQueryTransactionStatusResponse struct {
	Status            string      `json:"status"`
	Message           string      `json:"message"`
	TransactionId     string      `json:"transactionId"`
	TransactionStatus string      `json:"transactionStatus,omitempty"`
	ReceiptNumber     string      `json:"receiptNumber,omitempty"`
	CompletedTime     string      `json:"completedTime,omitempty"`
	IsReversed        string      `json:"isReversed,omitempty"`
	Data              interface{} `json:"data,omitempty"`
}

// NewProcessor creates a new CBE payment processor
func NewProcessor(config ProcessorConfig) payment.Processor {
	if config.BaseURL == "" {
		config.BaseURL = os.Getenv("CBE_BASE_URL")
	}
	fmt.Println("config.MerchantID", config.MerchantID)
	return &processor{
		merchantID:    config.MerchantID,
		merchantKey:   config.MerchantKey,
		terminalID:    config.TerminalID,
		credentialKey: config.CredentialKey,
		isTestMode:    config.IsTestMode,
		baseURL:       config.BaseURL,
		callbackURL:   config.CallbackURL,
		log:           logging.NewStdLogger("[CBE] [PROCESSOR]"),
	}
}

func (p *processor) InitiatePayment(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Info("Initiating CBE payment", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"amount":         req.Amount,
		"currency":       req.Currency,
	})

	// Prepare CBE specific request following the existing implementation
	cbeReq := map[string]interface{}{
		"amount":        req.Amount,
		"description":   "Deposit Request from " + req.PhoneNumber + " ID: " + req.TransactionID.String(),
		"referenceId":   req.TransactionID.String(),
		"callbackUrl":   os.Getenv("APP_URL_V2") + "/api/v2/settle/std",
		"phoneNumber":   req.PhoneNumber,
		"merchantId":    p.merchantID,
		"merchantKey":   p.merchantKey,
		"terminalId":    p.terminalID,
		"credentialKey": p.credentialKey,
		"userId":        "HABTAMUTA",
	}

	p.log.Info("Prepared CBE request", cbeReq)

	// Convert request to JSON
	jsonData, err := json.Marshal(cbeReq)
	if err != nil {
		p.log.Error("Failed to marshal request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/payments/create", p.baseURL),
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		p.log.Error("Failed to create HTTP request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	p.log.Info("Sending request to CBE", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	// Send request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.log.Error("Failed to send request to CBE", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

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
	var cbeResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			TransactionID string `json:"transactionId"`
			PaymentURL    string `json:"paymentUrl"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &cbeResp); err != nil {
		p.log.Error("Failed to decode CBE response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Map response status
	status := txEntity.PENDING
	if resp.StatusCode != http.StatusOK {
		status = txEntity.FAILED
	}

	return &payment.PaymentResponse{
		Success:       resp.StatusCode == http.StatusOK,
		TransactionID: req.TransactionID,
		Status:        status,
		ProcessorRef:  cbeResp.Data.TransactionID,
		PaymentURL:    cbeResp.Data.PaymentURL,
		Message:       cbeResp.Message,
	}, nil
}

// TODO; remove
func (p *processor) SettlePayment(ctx context.Context, req *payment.CallbackRequest) error {
	// Parse CBE-specific data from metadata
	var cbeCallback payment.STDCallbackRequest
	metadataBytes, err := json.Marshal(req.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := json.Unmarshal(metadataBytes, &cbeCallback); err != nil {
		return fmt.Errorf("failed to parse CBE callback data: %w", err)
	}

	p.log.Info("Processing CBE callback", map[string]interface{}{
		"reference_id": cbeCallback.ReferenceId,
		"status":       cbeCallback.Status,
		"message":      cbeCallback.Message,
	})

	// Map CBE status to our payment status
	var status txEntity.TransactionStatus
	switch cbeCallback.Status {
	case "SUCCESS":
		status = txEntity.SUCCESS
	case "FAILURE":
		status = txEntity.FAILED
	default:
		status = txEntity.FAILED
	}

	// Create metadata
	metadata := map[string]interface{}{
		"provider_tx_id": cbeCallback.ProviderTxId,
		"provider_data":  cbeCallback.ProviderData,
		"timestamp":      cbeCallback.Timestamp,
		"message":        cbeCallback.Message,
	}

	// Return appropriate error based on status
	if status == txEntity.SUCCESS {
		p.log.Info("Payment successful", metadata)
		return nil
	}

	p.log.Error("Payment failed", metadata)
	return fmt.Errorf("payment failed: %s", cbeCallback.Message)
}

func (p *processor) GetType() txEntity.TransactionMedium {
	return txEntity.CBE
}

func (p *processor) InitiateWithdrawal(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Info("Initiating CBE withdrawal", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"amount":         req.Amount,
		"currency":       req.Currency,
		"phone_number":   req.PhoneNumber,
	})

	// Prepare CBE specific request for withdrawal
	cbeReq := map[string]interface{}{
		"amount":        req.Amount,
		"description":   "Withdrawal Request from " + req.PhoneNumber + " ID: " + req.TransactionID.String(),
		"referenceId":   req.TransactionID.String(),
		"callbackUrl":   os.Getenv("APP_URL_V2") + "/api/v2/settle/std",
		"recipientId":   req.PhoneNumber,
		"merchantId":    p.merchantID,
		"merchantKey":   p.merchantKey,
		"terminalId":    p.terminalID,
		"credentialKey": p.credentialKey,
		"userId":        "HABTAMUTA",
	}

	p.log.Info("Prepared CBE withdrawal request", cbeReq)

	// Convert request to JSON
	jsonData, err := json.Marshal(cbeReq)
	if err != nil {
		p.log.Error("Failed to marshal withdrawal request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to marshal withdrawal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/withdrawals/create", p.baseURL),
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		p.log.Error("Failed to create HTTP withdrawal request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create withdrawal request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	p.log.Info("Sending withdrawal request to CBE", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	// Send request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.log.Error("Failed to send withdrawal request to CBE", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to send withdrawal request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		p.log.Error("Failed to read withdrawal response body", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to read withdrawal response body: %w", err)
	}

	p.log.Info("Withdrawal response body", map[string]interface{}{
		"body": string(bodyBytes),
	})

	// Parse response
	var cbeResp struct {
		Status      string `json:"status"`
		TraceId     string `json:"traceId"`
		ReferenceId string `json:"referenceId"`
		Message     string `json:"message"`
		Data        struct {
			OriginatorConversationID string `json:"OriginatorConversationID"`
			ConversationID           string `json:"ConversationID"`
			ResponseCode             string `json:"ResponseCode"`
			ResponseDesc             string `json:"ResponseDesc"`
			ServiceStatus            string `json:"ServiceStatus"`
			Timestamp                string `json:"Timestamp"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &cbeResp); err != nil {
		p.log.Error("Failed to decode CBE withdrawal response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to decode withdrawal response: %w", err)
	}

	// Map response status
	status := txEntity.PENDING
	if resp.StatusCode != http.StatusOK || cbeResp.Status != "SUCCESS" {
		status = txEntity.FAILED
		p.log.Error("Withdrawal request failed", map[string]interface{}{
			"status":       cbeResp.Status,
			"message":      cbeResp.Message,
			"responseCode": cbeResp.Data.ResponseCode,
			"responseDesc": cbeResp.Data.ResponseDesc,
		})
	} else {
		p.log.Info("Withdrawal request successful", map[string]interface{}{
			"traceId":        cbeResp.TraceId,
			"referenceId":    cbeResp.ReferenceId,
			"conversationID": cbeResp.Data.ConversationID,
			"responseCode":   cbeResp.Data.ResponseCode,
			"responseDesc":   cbeResp.Data.ResponseDesc,
		})
	}

	return &payment.PaymentResponse{
		Success:       resp.StatusCode == http.StatusOK && cbeResp.Status == "SUCCESS",
		TransactionID: req.TransactionID,
		Status:        status,
		ProcessorRef:  cbeResp.Data.ConversationID,
		Message:       cbeResp.Message,
		Metadata: map[string]interface{}{
			"originatorConversationID": cbeResp.Data.OriginatorConversationID,
			"conversationID":           cbeResp.Data.ConversationID,
			"responseCode":             cbeResp.Data.ResponseCode,
			"responseDesc":             cbeResp.Data.ResponseDesc,
			"serviceStatus":            cbeResp.Data.ServiceStatus,
			"timestamp":                cbeResp.Data.Timestamp,
		},
	}, nil
}

func (p *processor) QueryTransactionStatus(ctx context.Context, transactionID string) (*payment.TransactionStatusQueryResponse, error) {
	p.log.Info("Querying CBE transaction status", map[string]interface{}{
		"transaction_id": transactionID,
	})

	// Prepare CBE specific request for transaction status query
	jsonData := map[string]interface{}{
		"merchantId":          p.merchantID,
		"password":            p.merchantKey,
		"initiatorTerminalId": p.terminalID,
		"credentialKey":       p.credentialKey,
		"userId":              "HABTAMUTA",
		"transactionId":       transactionID,
	}

	p.log.Info("Prepared CBE transaction status request", jsonData)

	// Convert request to JSON
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		p.log.Error("Failed to marshal transaction status request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/transaction/status", p.baseURL),
		strings.NewReader(string(jsonBytes)),
	)
	if err != nil {
		p.log.Error("Failed to create HTTP request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	p.log.Info("Sending transaction status request to CBE", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	// Send request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.log.Error("Failed to send transaction status request to CBE", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		p.log.Error("Failed to read transaction status response body", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	p.log.Info("Transaction status response body", map[string]interface{}{
		"body": string(bodyBytes),
	})

	// Parse response
	var httpReqParsed CBEQueryTransactionStatusResponse
	if err := json.Unmarshal(bodyBytes, &httpReqParsed); err != nil {
		p.log.Error("Failed to decode CBE transaction status response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Map CBE status to our transaction status
	status := txEntity.PENDING
	switch httpReqParsed.Status {
	case "Completed":
		status = txEntity.SUCCESS
	case "Failed":
	case "FAILED":
		status = txEntity.FAILED
	case "Pending":
		status = txEntity.PENDING
	default:
		// For any unknown status, consider it as pending
		status = txEntity.PENDING
		p.log.Warn("Unknown CBE transaction status", map[string]interface{}{
			"status":         httpReqParsed.Status,
			"transaction_id": transactionID,
		})
	}

	p.log.Info("Transaction status query completed", map[string]interface{}{
		"transaction_id":     transactionID,
		"cbe_status":         httpReqParsed.Status,
		"mapped_status":      status,
		"receipt_number":     httpReqParsed.ReceiptNumber,
		"transaction_status": httpReqParsed.TransactionStatus,
		"completed_time":     httpReqParsed.CompletedTime,
		"is_reversed":        httpReqParsed.IsReversed,
	})

	// Prepare provider data
	providerData := make(map[string]interface{})
	if httpReqParsed.Data != nil {
		if data, ok := httpReqParsed.Data.(map[string]interface{}); ok {
			providerData = data
		} else {
			// If data is not a map, wrap it in a map
			providerData["data"] = httpReqParsed.Data
		}
	}

	// Add additional fields from the response to provider data
	providerData["message"] = httpReqParsed.Message
	providerData["transaction_id"] = httpReqParsed.TransactionId
	providerData["transaction_status"] = httpReqParsed.TransactionStatus
	providerData["completed_time"] = httpReqParsed.CompletedTime
	providerData["is_reversed"] = httpReqParsed.IsReversed

	return &payment.TransactionStatusQueryResponse{
		Status:       status,
		ProviderTxId: httpReqParsed.ReceiptNumber,
		ProviderData: providerData,
	}, nil
}
