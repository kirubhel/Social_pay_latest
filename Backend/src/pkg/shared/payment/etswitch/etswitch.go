package etswitch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

type processor struct {
	userName    string
	credentials string
	isTestMode  bool // test mode
	baseURL     string
	log         logging.Logger
}

// ProcessorConfig  for Etswitch processor
type ProcessorConfig struct {
	UserName   string
	Credential string
	BaseURL    string
	IsTestMode bool
}

// Constructor

func NewEtSwitchProcessor(cfg ProcessorConfig) payment.Processor {

	if cfg.BaseURL == "" {
		cfg.BaseURL = os.Getenv("ETHSWITCH_BASE_URL")
	}
	if cfg.UserName == "" {
		cfg.UserName = os.Getenv("ETHSWITCH_USERNAME")
	}
	if cfg.Credential == "" {
		cfg.Credential = os.Getenv("ETHSWITCH_PASSWORD")
	}

	return &processor{
		userName:    cfg.UserName,
		credentials: cfg.Credential,
		baseURL:     cfg.BaseURL,
		log:         logging.NewStdLogger("ETHSWITH_PROCESSOR_LOG"),
	}

}

// Initiate EthSwitch host-checkout onephase payment
func (p *processor) InitiatePayment(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {

	// Request id
	ReqId := uuid.NewString()

	// Logging the initiation
	p.log.Info("Initiating EthSwitch hosted checkout", map[string]interface{}{
		"transactionID": req.TransactionID,
		"requestID":     ReqId,
		"amount":        req.Amount,
		"currency":      req.Currency,
	})

	// Converting the amount float types to minor deminator
	amount := int(math.Round(req.Amount * 100))
	// Mapint the TransactionId to OrderNumber
	orderNumber := MapTransactionIDToOrderNumber(req.TransactionID)

	// Build query parameters
	params := url.Values{}
	params.Set("userName", p.userName)
	params.Set("password", p.credentials)
	params.Set("amount", strconv.Itoa(amount))
	// params.Set("amount", strconv.FormatFloat(req.Amount, 'f', -1, 64))
	params.Set("currency", CurrencyToISOCode[req.Currency])
	params.Set("orderNumber", orderNumber)
	params.Set("returnUrl", req.SuccessURL)

	fullURL := fmt.Sprintf("%s?%s", p.baseURL, params.Encode())

	// Create request with context
	getReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Use http.Client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,

		// For test case only
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip the ssl for test purpose only
			},
		},
	}

	// Send the request
	resp, err := client.Do(getReq)
	if err != nil {
		return nil, fmt.Errorf("failed to contact EthSwitch: %w", err)
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

	p.log.Info("EthSwitch HostedCheckout Response body", map[string]interface{}{
		"body": string(bodyBytes),
	})

	// Parse response body
	var Response EthSwitchResponse

	if err := json.Unmarshal(bodyBytes, &Response); err != nil {

		p.log.Error("failed to decode EthSwitch response", map[string]interface{}{
			"err": err,
		})

		return nil, err
	}

	// EthSwitch Fallback
	if Response.ErrorCode != 0 {
		p.log.Error("failed to initiata the ethswitch the request", map[string]interface{}{
			"err":       Response.ErrorMessage,
			"operation": "initiate payement",
		})

		return nil, errors.New(Response.ErrorMessage)
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
		Message:       "Host checkout created successfully", // Update
		PaymentURL:    Response.FormUrl,
		// add other if any
	}, nil

}

// SettlePayment settle the payment
func (p *processor) SettlePayment(ctx context.Context, req *payment.CallbackRequest) error {

	p.log.Info("Processing EthSwitch callback", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"status":         req.Status,
	})

	decision, _ := req.Metadata["decision"].(string)

	var status txEntity.HostedPaymentStatus

	switch decision {
	case "SUCCESS":
		status = txEntity.HostedPaymentCompleted
	case "CANCEL":
		status = txEntity.HostedPaymentCanceled
	case "DECLINE":
		status = txEntity.HostedPaymentStatus(txEntity.FAILED)
	default:
		status = txEntity.HostedPaymentStatus(txEntity.CANCELED)
	}

	if status == txEntity.HostedPaymentCompleted {
		return nil
	}

	return fmt.Errorf("payment failed: %s", decision)
}

// GetType return the transaction medium
func (p *processor) GetType() txEntity.TransactionMedium {
	return txEntity.ETHSWITCH
}

func (p *processor) InitiateWithdrawal(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Error("Withdrawal not supported", map[string]interface{}{
		"processor": "EthSwitch",
	})
	return nil, fmt.Errorf("withdrawal not supported for EthSwitch")
}

// MapTransactionIDToOrderNumber maps a UUID to an AN1.32-compatible string (for EthSwitch)
func MapTransactionIDToOrderNumber(txID uuid.UUID) string {
	return strings.ReplaceAll(txID.String(), "-", "") // returns 32-char alphanumeric string
}
