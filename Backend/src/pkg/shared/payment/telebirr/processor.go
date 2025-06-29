package telebirr

import (
	"bytes"
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
	securityCredential string
	password           string
	isTestMode         bool
	shortCode          string
	identityID         string
	baseURL            string
	callbackURL        string
	log                logging.Logger
}

// ProcessorConfig holds the configuration for Telebirr processor
type ProcessorConfig struct {
	SecurityCredential string
	Password           string
	IsTestMode         bool
	ShortCode          string
	IdentityID         string
	BaseURL            string
	CallbackURL        string
}

// NewProcessor creates a new Telebirr payment processor
func NewProcessor(config ProcessorConfig) payment.Processor {
	if config.BaseURL == "" {
		if config.IsTestMode {
			config.BaseURL = "https://api-test.telebirr.com"
		} else {
			config.BaseURL = "https://api.telebirr.com"
		}
	}

	return &processor{
		securityCredential: config.SecurityCredential,
		password:           config.Password,
		isTestMode:         config.IsTestMode,
		shortCode:          config.ShortCode,
		identityID:         config.IdentityID,
		baseURL:            config.BaseURL,
		callbackURL:        config.CallbackURL,
		log:                logging.NewStdLogger("[TELEBIRR] [PROCESSOR]"),
	}
}

func (p *processor) InitiatePayment(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Info("Initiating Telebirr payment", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"amount":         req.Amount,
		"currency":       req.Currency,
	})

	// Prepare Telebirr specific request
	telebirrReq := map[string]interface{}{
		"command_id":         "InitTrans_BuyGoodsForCustomer",
		"conversation_id":    req.TransactionID.String(),
		"thirdPartyID":       "Social-Pay",
		"password":           p.password,
		"resultURL":          os.Getenv("APP_URL_V2") + "/api/v2/settle/std",
		"timestamp":          time.Now().Format("20060102150405"),
		"identifier_type":    12,
		"identifier":         p.identityID,
		"securityCredential": p.securityCredential,
		"shortCode":          p.shortCode,
		"primary_party":      req.PhoneNumber,
		"receiver_party":     p.shortCode,
		"amount":             req.Amount,
		"currency":           req.Currency,
	}

	p.log.Info("Prepared Telebirr request", map[string]interface{}{
		"command_id":    telebirrReq["CommandID"],
		"originator_id": telebirrReq["OriginatorConversationID"],
	})

	// Convert request to JSON
	jsonData, err := json.Marshal(telebirrReq)
	if err != nil {
		p.log.Error("Failed to marshal request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	p.log.Info("Telebirr request", map[string]interface{}{
		"request": string(jsonData),
	})

	// Create HTTP request
	p.log.Info("[TELEBIRR] Creating request to Telebirr", map[string]interface{}{
		"url": fmt.Sprintf("%s/api/account/telebirr/ussd-push", p.baseURL),
	})
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/account/telebirr/ussd-push", p.baseURL),
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
	httpReq.Header.Set("X-API-KEY", apikey)
	p.log.Info("[TELEBIRR] Added headers - Content-Type: application/json, X-API-KEY: [REDACTED]", map[string]interface{}{
		"x_api_key": maskAPIKey(apikey),
	})

	p.log.Info("Sending request to Telebirr", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	// Send request
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.log.Error("Failed to send request to Telebirr", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	p.log.Info("Received response from Telebirr", map[string]interface{}{
		"status": resp.Status,
	})

	// Read response body into string
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		p.log.Error("Failed to read response body", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	// Create new reader with the bytes for subsequent operations
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	p.log.Info("Response body", map[string]interface{}{
		"body": string(bodyBytes),
	})

	// Parse response
	var telebirrResp struct {
		Success bool `json:"success"`
		Data    struct {
			ConversationID string `json:"conversation_id"`
			Message        string `json:"message"`
			ResponseCode   int    `json:"response_code"`
			ResponseDesc   string `json:"response_desc"`
			ServiceStatus  int    `json:"service_status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&telebirrResp); err != nil {
		p.log.Error("Failed to decode Telebirr response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("something went wrong while initiating payment")
	}

	p.log.Info("Received response from Telebirr", map[string]interface{}{
		"success":         telebirrResp.Success,
		"message":         telebirrResp.Data.Message,
		"conversation_id": telebirrResp.Data.ConversationID,
		"response_code":   telebirrResp.Data.ResponseCode,
	})

	// Map response status
	status := txEntity.PENDING
	if telebirrResp.Data.ResponseCode != 0 {
		status = txEntity.FAILED
	}

	return &payment.PaymentResponse{
		Success:       telebirrResp.Success,
		TransactionID: req.TransactionID,
		Status:        status,
		ProcessorRef:  telebirrResp.Data.ConversationID,
		Message:       telebirrResp.Data.Message,
	}, nil
}

func (p *processor) SettlePayment(ctx context.Context, req *payment.CallbackRequest) error {
	p.log.Info("Processing Telebirr callback", map[string]interface{}{
		"status": req.Status,
	})

	// The callback request should contain the status from Telebirr
	// Map Telebirr status to our payment status
	switch req.Status {
	case txEntity.SUCCESS:
		p.log.Info("Payment successful", map[string]interface{}{
			"metadata": req.Metadata,
		})
		return nil
	case txEntity.FAILED:
		p.log.Error("Payment failed", map[string]interface{}{
			"reason": req.Metadata["reason"],
		})
		return fmt.Errorf("payment failed: %s", req.Metadata["reason"])
	case txEntity.EXPIRED:
		p.log.Info("Payment expired", nil)
		return fmt.Errorf("payment expired")
	default:
		p.log.Error("Unknown payment status", map[string]interface{}{
			"status": req.Status,
		})
		return fmt.Errorf("unknown payment status: %s", req.Status)
	}
}

func (p *processor) GetType() txEntity.TransactionMedium {
	return txEntity.TELEBIRR
}

func (p *processor) InitiateWithdrawal(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Info("Initiating Telebirr withdrawal", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"amount":         req.Amount,
		"currency":       req.Currency,
		"phone_number":   req.PhoneNumber,
	})

	// Prepare Telebirr specific request for B2C (withdrawal)
	telebirrReq := map[string]interface{}{
		"OriginatorConversationID": req.TransactionID.String(),
		"InitiatorName":            "51437702",
		"SecurityCredential":       "Hh3kw29A1beSRK2cgALN52bX7Io8GCqEvQ9CFkqn/Qg=",
		"CommandID":                "InitTrans_2003",
		"PartyA":                   "514377",
		"PartyB":                   req.PhoneNumber,
		"Amount":                   req.Amount,
		"Currency":                 req.Currency,
		"Remarks":                  "Withdrawal Request from " + req.PhoneNumber + " ID: " + req.TransactionID.String(),
		"Occasion":                 "Pay for Individual",
		"QueueTimeOutURL":          "/api/accounts/telebirr/b2c/transactionstatus",
		"ResultURL":                "http://api.socialpay.et:6080/api/accounts/telebirr/b2c/transactionstatus", //TODO: get from env
		"ThirdPartyID":             "Social-Pay",
		"Password":                 "jBq7JfxTs0C5ji0VPKakmRSgBbeh4NO0juJ1LXnPIOw=",
	}

	p.log.Info("Prepared Telebirr withdrawal request", map[string]interface{}{
		"telebirr_req": telebirrReq,
	})
	p.log.Info("Prepared Telebirr withdrawal request", map[string]interface{}{
		"originator_id": telebirrReq["OriginatorConversationID"],
		"amount":        telebirrReq["Amount"],
		"currency":      telebirrReq["Currency"],
		"phone_number":  telebirrReq["PartyB"],
	})

	// Convert request to JSON
	jsonData, err := json.Marshal(telebirrReq)
	if err != nil {
		p.log.Error("Failed to marshal withdrawal request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to marshal withdrawal request: %w", err)
	}

	p.log.Info("Telebirr withdrawal request", map[string]interface{}{
		"request": string(jsonData),
	})

	// Create HTTP request
	p.log.Info("[TELEBIRR] Creating withdrawal request to Telebirr", map[string]interface{}{
		"url": fmt.Sprintf("%s/api/accounts/telebirr/payment/b2c", p.baseURL),
	})
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/accounts/telebirr/payment/b2c", p.baseURL),
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
	httpReq.Header.Set("X-API-KEY", apikey)
	p.log.Info("[TELEBIRR] Added headers - Content-Type: application/json, X-API-KEY: [REDACTED]", map[string]interface{}{
		"x_api_key": maskAPIKey(apikey),
	})

	p.log.Info("Sending withdrawal request to Telebirr", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	// Send request
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.log.Error("Failed to send withdrawal request to Telebirr", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to send withdrawal request: %w", err)
	}
	defer resp.Body.Close()

	p.log.Info("Received response from Telebirr for withdrawal", map[string]interface{}{
		"status": resp.Status,
	})

	// Read response body into string
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		p.log.Error("Failed to read withdrawal response body", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to read withdrawal response body: %w", err)
	}
	// Create new reader with the bytes for subsequent operations
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	p.log.Info("Withdrawal response body", map[string]interface{}{
		"body": string(bodyBytes),
	})

	// Parse response - using the same structure as payment since response is expected to be similar
	var telebirrResp struct {
		Success bool   `json:"success"`
		Data    string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&telebirrResp); err != nil {
		p.log.Error("Failed to decode Telebirr withdrawal response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("something went wrong while initiating withdrawal")
	}

	p.log.Info("Processed withdrawal response from Telebirr", map[string]interface{}{
		"success": telebirrResp.Success,
		"message": telebirrResp.Data,
	})

	// Map response status
	status := txEntity.PENDING
	if resp.StatusCode != 200 {
		status = txEntity.FAILED
	}

	return &payment.PaymentResponse{
		Success:       telebirrResp.Success,
		TransactionID: req.TransactionID,
		Status:        status,
		ProcessorRef:  req.TransactionID.String(),
		Message:       telebirrResp.Data,
		Metadata: map[string]interface{}{
			"reason": telebirrResp.Data,
		},
	}, nil
}

func maskAPIKey(key string) string {
	if len(key) <= 4 {
		return key
	}
	return "****" + key[len(key)-4:]
}
