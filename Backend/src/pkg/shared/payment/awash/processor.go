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
	"github.com/socialpay/socialpay/src/pkg/shared/payment/controller/gin"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	transactionRepo "github.com/socialpay/socialpay/src/pkg/transaction/core/repository"
)

type processor struct {
	merchantID         string
	merchantTillNumber string
	// terminalID         string
	credentialKey string
	isTestMode    bool
	baseURL       string
	callbackURL   string
	txn           transactionRepo.TransactionRepository
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
	TxnRepository      transactionRepo.TransactionRepository
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
		txn:                cfg.TxnRepository,
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
	// the fallback if awash fails to intiate the request
	if awashResponse.ReturnCode > 0 {
		err := errors.New(awashResponse.ReturnMessage)
		return nil, err
	}

	// Map response status
	status := txEntity.PENDING
	if resp.StatusCode != http.StatusOK {
		status = txEntity.FAILED
	}

	if status == txEntity.PENDING {
		go p.scheduleFailFallback(context.Background(), *req, 5*time.Minute)
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

func (p *processor) QueryTransactionStatus(ctx context.Context, transactionID string) (*payment.TransactionStatusQueryResponse, error) {

	ReqId := uuid.NewString()
	p.log.Info("Querying Awash transaction status", map[string]interface{}{
		"transaction_id": transactionID,
	})
	// Building Url
	url := fmt.Sprintf(`%v/%v`, p.baseURL, "/MerchantRS/DebitStatus")
	// Preparing body
	body := map[string]interface{}{
		"authorization": map[string]interface{}{
			"merchantCode":       p.merchantID,
			"merchantTillNumber": p.merchantTillNumber,
			"requestId":          ReqId,
			"requestSignature": p.GenerateSignature(p.credentialKey,
				ReqId),
		},
		"paymentRequest": map[string]interface{}{
			"externalReference": transactionID,
		},
	}

	bytes, err := json.Marshal(body)

	if err != nil {
		// Logging
		p.log.Error("failed to marshal the body for the request", map[string]interface{}{
			"operation": "QueryTrasactionStatus::Marshal",
			"error":     err.Error(),
		})

		return nil, err
	}
	// Creating the http request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(bytes)))

	if err != nil {
		// Logging error
		p.log.Error("failed to create debit checking request ", map[string]interface{}{
			"operation": "QueryTransaction:createrequest",
			"error":     err.Error(),
		})

		return nil, err
	}
	// Making request
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		// Logging error
		p.log.Error("failed to do request with transaction id ", map[string]interface{}{
			"operation": "Client request",
			"error":     err.Error(),
		})

		return nil, err

	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	p.log.Info("get this response for querytransactionstatus", map[string]interface{}{
		"operations": "QueryTransactionStatus",
		"body":       string(respBody),
	})

	var awashResp gin.AwashCallback
	if err := json.Unmarshal(respBody, &awashResp); err != nil {

		p.log.Error("failed to unmarshall the response body ", map[string]interface{}{
			"operation": "QueryTransactionStatus::",
			"error":     err.Error(),
			"body":      string(respBody),
		})

	}

	var status txEntity.TransactionStatus
	switch awashResp.Status {
	case "APPROVED":
		status = txEntity.SUCCESS
	case "EXPIRED":
		status = txEntity.EXPIRED
	default:
		status = txEntity.FAILED
	}

	// Preparing the providerData
	providerData := make(map[string]interface{})
	providerData["TransactionId"] = awashResp.TransactionID
	providerData["Status"] = awashResp.Status
	providerData["ErrorCode"] = awashResp.ReturnCode

	return &payment.TransactionStatusQueryResponse{
		Status:       status,
		ProviderTxId: awashResp.TransactionID,
		ProviderData: providerData,
	}, nil
}

// scheduleFailFallback triggers a self-callback after a delay if the transaction is still pending.
func (p *processor) scheduleFailFallback(ctx context.Context, req payment.PaymentRequest, delay time.Duration) {
	select {
	case <-time.After(delay):
		// Proceed after delay
	case <-ctx.Done():
		p.log.Warn("scheduleFailFallback aborted due to context cancellation", map[string]interface{}{
			"txn_id": req.TransactionID,
		})
		return
	}

	txn, err := p.txn.GetByID(ctx, req.TransactionID)
	if err != nil {
		p.log.Error("Failed to fetch transaction during fallback", map[string]interface{}{
			"operation": "scheduleFailFallback::AWASH",
			"txn_id":    req.TransactionID,
			"err":       err.Error(),
		})
		return
	}

	if txn.Status != txEntity.PENDING {
		p.log.Info("Transaction already completed; skipping fallback", map[string]interface{}{
			"txn_id": req.TransactionID,
			"status": txn.Status,
		})
		return
	}

	callback := gin.AwashCallback{
		ReturnCode:        1,
		PayerPhone:        req.PhoneNumber,
		Amount:            req.Amount,
		ReturnMessage:     "Failed to process your transaction",
		Status:            string(txEntity.FAILED),
		ExternalReference: req.TransactionID.String(),
		TransactionID:     "",
		DateApproved:      time.Now().Format("2006-01-02 15:04:05"),
		DateRequested:     "",
	}

	bodyBytes, err := json.Marshal(callback)
	if err != nil {
		p.log.Error("Failed to marshal fallback callback body", map[string]interface{}{
			"operation": "scheduleFailFallback::Marshal",
			"err":       err.Error(),
		})
		return
	}

	p.log.Info("Triggering fallback callback", map[string]interface{}{
		"url":  p.callbackURL,
		"body": string(bodyBytes),
	})

	res, err := http.Post(p.callbackURL, "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		p.log.Error("Failed to send fallback callback", map[string]interface{}{
			"operation": "scheduleFailFallback::HTTP_POST",
			"err":       err.Error(),
		})
		return
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)

	p.log.Info("Fallback callback response", map[string]interface{}{
		"status_code": res.StatusCode,
		"body":        string(respBody),
	})
}
