package awash

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

type processor struct {
	merchantID         string
	merchantTillNumber string
	// terminalID         string
	credentialKey string
	isTestMode    bool
	baseURL       string
	callbackURL   string
	log           logging.Logger
}

// processor AWASH CONFIG
type ProcessorConfig struct {
	MerchantID string
	// MerchantKey   string
	CredentialKey      string
	IsTestMode         bool
	BaseURL            string
	CallbackURL        string
	MerchantTillNumber string
}

// initiator
func NewProcessor(cfg ProcessorConfig) payment.Processor {

	// BaseUlr is not defined
	if cfg.BaseURL == "" {
		cfg.BaseURL = os.Getenv("AWASH_TEST_BASE_URL")
	}
	// Callback Url
	if cfg.CallbackURL == "" {
		cfg.CallbackURL = os.Getenv("AWASH_TEST_CALLBACK_URL")

	}

	// Awash processor
	return &processor{
		merchantID:         cfg.MerchantID,
		merchantTillNumber: cfg.MerchantTillNumber,
		baseURL:            cfg.BaseURL,
		isTestMode:         cfg.IsTestMode,
		credentialKey:      cfg.CredentialKey,
		callbackURL:        cfg.CallbackURL,
		log:                logging.NewStdLogger("AWASH_LOG::"),
	}
}

// GenerateSignature generates a signature for the request
func (p *processor) GenerateSignature(credintial, requestID string) string {
	// Concatenate the merchant key and request ID
	input := credintial + requestID

	// Generate SHA-256 hash
	hash := sha256.Sum256([]byte(input))

	// Return the hexadecimal representation of the hash
	return hex.EncodeToString(hash[:])
}

// Payment Initiator
func (p *processor) InitiatePayment(c context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {

	// Request id
	ReqId := uuid.NewString()

	// Logging the initiation
	p.log.Info("Initiating AWASH payment", map[string]interface{}{
		"merchantID":    p.merchantID,
		"transactionID": req.TransactionID,
		"requestID":     ReqId,
		"amount":        req.Amount,
		"currency":      req.Currency,
	})

	// Prepare the request body
	awashReq := map[string]interface{}{
		"authorization": map[string]interface{}{
			"merchantCode":       p.merchantID,
			"merchantTillNumber": p.merchantTillNumber,
			"requestId":          ReqId,
			"requestSignature": p.GenerateSignature(p.credentialKey,
				ReqId),
		},
		"paymentRequest": map[string]interface{}{
			"payerPhone":        req.PhoneNumber,
			"reason":            req.Description,
			"amount":            fmt.Sprintf("%.2f", req.Amount),
			"externalReference": req.TransactionID,
			"callbackUrl":       p.callbackURL,
		},
	}

	p.log.Info("Prepared AWASH BIRR request", awashReq)

	body, err := json.Marshal(awashReq)

	if err != nil {
		p.log.Error("failed to marshal the AWASH Request body", map[string]interface{}{
			"err": err.Error(),
		})
		return nil, err
	}

	// Creating HTTP request
	httpReq, err := http.NewRequestWithContext(c,
		http.MethodPost,
		fmt.Sprintf("%v/MerchantRS/DebitRequest", p.baseURL),
		strings.NewReader(string(body)))

	if err != nil {
		p.log.Error("Failed to create http request", map[string]interface{}{
			"err": err.Error(),
		})

		return nil, err
	}
	// Setting headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Client with timeout
	// Note: The timeout is set to 60 seconds, adjust as necessary
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.log.Error("Failed to send request to AWASH BIRR", map[string]interface{}{
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

	p.log.Info("Awash Response body", map[string]interface{}{
		"body": string(bodyBytes),
	})

	// parse response body
	var awashResponse AwashResponse

	if err := json.Unmarshal(bodyBytes, &awashResponse); err != nil {

		p.log.Error("failed to decode AWASH", map[string]interface{}{
			"err": err,
		})

		return nil, err
	}
	// the fallback if awash
	if awashResponse.ReturnCode > 0 {
		err := errors.New(awashResponse.ReturnMessage)
		return nil, err
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
		Message:       awashResponse.ReturnMessage,
		// add other if any
	}, nil
}

func (p *processor) GetType() txEntity.TransactionMedium {
	return txEntity.AWASH
}

// Todo remove
func (p *processor) SettlePayment(ctx context.Context, req *payment.CallbackRequest) error {

	return nil
}

// To be implemented
func (p *processor) InitiateWithdrawal(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {

	return nil, nil
}
