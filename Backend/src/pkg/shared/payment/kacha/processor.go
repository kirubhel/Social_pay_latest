package kacha

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
	shortCode   string
	isTestMode  bool
	baseURL     string
	callbackURL string
	log         logging.Logger
}

// ProcessorConfig holds the configuration for Kacha processor
type ProcessorConfig struct {
	IsTestMode  bool
	BaseURL     string
	CallbackURL string
}

// KachaPaymentResponse represents the response from Kacha payment API
type KachaPaymentResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		Detail      string  `json:"detail"`
		Fee         float64 `json:"fee"`
		FromName    string  `json:"from_name"`
		ID          string  `json:"id"`
		Message     string  `json:"message"`
		Phone       string  `json:"phone"`
		Process     string  `json:"process"`
		Reason      string  `json:"reason"`
		Reference   string  `json:"reference"`
		Status      string  `json:"status"`
		StatusCode  int     `json:"status_code"`
		TraceNumber string  `json:"trace_number"`
		UpdatedAt   string  `json:"updated_at"`
	} `json:"data"`
}

// KachaWithdrawalResponse represents the response from Kacha withdrawal API
type KachaWithdrawalResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		Detail      string  `json:"detail"`
		Fee         float64 `json:"fee"`
		FromName    string  `json:"from_name"`
		ID          string  `json:"id"`
		Message     string  `json:"message"`
		Phone       string  `json:"phone"`
		Process     string  `json:"process"`
		Reason      string  `json:"reason"`
		Reference   string  `json:"reference"`
		Status      string  `json:"status"`
		StatusCode  int     `json:"status_code"`
		TraceNumber string  `json:"trace_number"`
		UpdatedAt   string  `json:"updated_at"`
		Error       *struct {
			Detail     string `json:"detail"`
			Message    string `json:"message"`
			Status     string `json:"status"`
			StatusCode string `json:"status_code"`
		} `json:"error,omitempty"`
	} `json:"data"`
}

// KachaCallback represents the callback data from Kacha
type KachaCallback struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Detail      string  `json:"detail"`
	Fee         float64 `json:"fee"`
	FromName    string  `json:"from_name"`
	ID          string  `json:"id"`
	Message     string  `json:"message"`
	Phone       string  `json:"phone"`
	Process     string  `json:"process"`
	Reason      string  `json:"reason"`
	Reference   string  `json:"reference"`
	Status      string  `json:"status"`
	StatusCode  int     `json:"status_code"`
	TraceNumber string  `json:"trace_number"`
	UpdatedAt   string  `json:"updated_at"`
}

// NewProcessor creates a new Kacha payment processor
func NewProcessor(config ProcessorConfig) payment.Processor {
	if config.BaseURL == "" {
		config.BaseURL = os.Getenv("KACHA_BASE_URL")
		if config.BaseURL == "" {
			config.BaseURL = "https://kacha-sv.socialpay.co"
		}
	}

	return &processor{
		isTestMode:  config.IsTestMode,
		baseURL:     config.BaseURL,
		callbackURL: config.CallbackURL,
		log:         logging.NewStdLogger("[KACHA] [PROCESSOR]"),
	}
}

func (p *processor) InitiatePayment(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Info("Initiating Kacha payment", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"amount":         req.Amount,
		"currency":       req.Currency,
		"phone_number":   req.PhoneNumber,
	})

	// Prepare Kacha specific request
	kachaReq := map[string]interface{}{
		"callback_url": req.CallbackURL,
		"phone":        req.PhoneNumber,
		"amount":       req.Amount,
		"trace_number": req.TransactionID.String(),
		"reason":       "payment",
	}

	p.log.Info("Prepared Kacha request", kachaReq)

	// Convert request to JSON
	jsonData, err := json.Marshal(kachaReq)
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
		fmt.Sprintf("%s/api/kacha/push-ussd-payment", p.baseURL),
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

	p.log.Info("Sending request to Kacha", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	// Send request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.log.Error("Failed to send request to Kacha", map[string]interface{}{
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
		"body":        string(bodyBytes),
		"status_code": resp.StatusCode,
	})

	// Parse response
	var kachaResp KachaPaymentResponse
	if err := json.Unmarshal(bodyBytes, &kachaResp); err != nil {
		p.log.Error("Failed to decode Kacha response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Map response status
	status := txEntity.PENDING
	if !kachaResp.Success || resp.StatusCode != http.StatusOK {
		status = txEntity.FAILED
	} else if kachaResp.Data.Status == "SUCCESS" {
		status = txEntity.SUCCESS
	} else if kachaResp.Data.Status == "FAILED" {
		status = txEntity.FAILED
	}

	return &payment.PaymentResponse{
		Success:       kachaResp.Success,
		TransactionID: req.TransactionID,
		Status:        status,
		ProcessorRef:  kachaResp.Data.Reference,
		Message:       kachaResp.Data.Message,
		Metadata: map[string]interface{}{
			"kacha_id":     kachaResp.Data.ID,
			"trace_number": kachaResp.Data.TraceNumber,
			"from_name":    kachaResp.Data.FromName,
			"description":  kachaResp.Data.Description,
			"detail":       kachaResp.Data.Detail,
			"process":      kachaResp.Data.Process,
			"fee":          kachaResp.Data.Fee,
			"status_code":  kachaResp.Data.StatusCode,
			"updated_at":   kachaResp.Data.UpdatedAt,
		},
	}, nil
}

func (p *processor) SettlePayment(ctx context.Context, req *payment.CallbackRequest) error {
	// Parse Kacha-specific data from metadata
	var kachaCallback KachaCallback
	metadataBytes, err := json.Marshal(req.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := json.Unmarshal(metadataBytes, &kachaCallback); err != nil {
		return fmt.Errorf("failed to parse Kacha callback data: %w", err)
	}

	p.log.Info("Processing Kacha callback", map[string]interface{}{
		"reference":    kachaCallback.Reference,
		"status":       kachaCallback.Status,
		"message":      kachaCallback.Message,
		"trace_number": kachaCallback.TraceNumber,
	})

	// Map Kacha status to our payment status
	var status txEntity.TransactionStatus
	switch strings.ToUpper(kachaCallback.Status) {
	case "SUCCESS":
		status = txEntity.SUCCESS
	case "FAILED", "FAILURE":
		status = txEntity.FAILED
	case "PENDING":
		status = txEntity.PENDING
	default:
		status = txEntity.FAILED
	}

	// Create metadata
	metadata := map[string]interface{}{
		"kacha_id":     kachaCallback.ID,
		"trace_number": kachaCallback.TraceNumber,
		"from_name":    kachaCallback.FromName,
		"description":  kachaCallback.Description,
		"detail":       kachaCallback.Detail,
		"process":      kachaCallback.Process,
		"fee":          kachaCallback.Fee,
		"status_code":  kachaCallback.StatusCode,
		"updated_at":   kachaCallback.UpdatedAt,
		"message":      kachaCallback.Message,
	}

	// Return appropriate error based on status
	if status == txEntity.SUCCESS {
		p.log.Info("Payment successful", metadata)
		return nil
	}

	p.log.Error("Payment failed", metadata)
	return fmt.Errorf("payment failed: %s", kachaCallback.Message)
}

func (p *processor) GetType() txEntity.TransactionMedium {
	return txEntity.KACHA
}

func (p *processor) InitiateWithdrawal(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Info("Initiating Kacha withdrawal", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"amount":         req.Amount,
		"currency":       req.Currency,
		"phone_number":   req.PhoneNumber,
	})

	// Prepare Kacha specific request for withdrawal
	kachaReq := map[string]interface{}{
		"to":         req.PhoneNumber,
		"amount":     req.Amount,
		"reason":     "salary_payments",
		"short_code": p.shortCode,
	}

	p.log.Info("Prepared Kacha withdrawal request", kachaReq)

	// Convert request to JSON
	jsonData, err := json.Marshal(kachaReq)
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
		fmt.Sprintf("%s/api/b2c/kacha/initiate-transfer", p.baseURL),
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

	p.log.Info("Sending withdrawal request to Kacha", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	// Send request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.log.Error("Failed to send withdrawal request to Kacha", map[string]interface{}{
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
		"body":        string(bodyBytes),
		"status_code": resp.StatusCode,
	})

	// Parse response
	var kachaResp KachaWithdrawalResponse
	if err := json.Unmarshal(bodyBytes, &kachaResp); err != nil {
		p.log.Error("Failed to decode Kacha withdrawal response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to decode withdrawal response: %w", err)
	}

	// Map response status
	status := txEntity.PENDING
	message := "Withdrawal initiated"

	if !kachaResp.Success || resp.StatusCode != http.StatusOK {
		status = txEntity.FAILED
		if kachaResp.Data.Error != nil {
			message = kachaResp.Data.Error.Detail
			p.log.Error("Withdrawal request failed", map[string]interface{}{
				"error_detail":  kachaResp.Data.Error.Detail,
				"error_message": kachaResp.Data.Error.Message,
				"error_status":  kachaResp.Data.Error.Status,
				"error_code":    kachaResp.Data.Error.StatusCode,
			})
		}
	} else {
		// Success case
		if kachaResp.Data.Status == "SUCCESS" {
			status = txEntity.SUCCESS
		} else if kachaResp.Data.Status == "FAILED" {
			status = txEntity.FAILED
		}
		message = kachaResp.Data.Message

		p.log.Info("Withdrawal request successful", map[string]interface{}{
			"kacha_id":     kachaResp.Data.ID,
			"reference":    kachaResp.Data.Reference,
			"trace_number": kachaResp.Data.TraceNumber,
			"status":       kachaResp.Data.Status,
			"message":      kachaResp.Data.Message,
		})
	}

	metadata := map[string]interface{}{}
	if kachaResp.Data.Error != nil {
		metadata["error"] = map[string]interface{}{
			"detail":      kachaResp.Data.Error.Detail,
			"message":     kachaResp.Data.Error.Message,
			"status":      kachaResp.Data.Error.Status,
			"status_code": kachaResp.Data.Error.StatusCode,
		}
	} else {
		metadata = map[string]interface{}{
			"kacha_id":     kachaResp.Data.ID,
			"trace_number": kachaResp.Data.TraceNumber,
			"from_name":    kachaResp.Data.FromName,
			"description":  kachaResp.Data.Description,
			"detail":       kachaResp.Data.Detail,
			"process":      kachaResp.Data.Process,
			"fee":          kachaResp.Data.Fee,
			"status_code":  kachaResp.Data.StatusCode,
			"updated_at":   kachaResp.Data.UpdatedAt,
		}
	}

	return &payment.PaymentResponse{
		Success:       kachaResp.Success,
		TransactionID: req.TransactionID,
		Status:        status,
		ProcessorRef:  kachaResp.Data.Reference,
		Message:       message,
		Metadata:      metadata,
	}, nil
}

func (p *processor) QueryTransactionStatus(ctx context.Context, transactionID string) (*payment.TransactionStatusQueryResponse, error) {
	p.log.Info("Querying Kacha transaction status", map[string]interface{}{
		"transaction_id": transactionID,
	})

	// Note: Kacha doesn't provide a specific transaction status query endpoint in the examples
	// This is a placeholder implementation - you may need to adjust this based on actual Kacha API
	p.log.Warn("Kacha transaction status query not implemented - no API endpoint provided", map[string]interface{}{
		"transaction_id": transactionID,
	})

	return &payment.TransactionStatusQueryResponse{
		Status:       txEntity.PENDING,
		ProviderTxId: "",
		ProviderData: map[string]interface{}{
			"message": "Transaction status query not supported by Kacha API",
		},
	}, nil
}
