package cybersource

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

type processor struct {
	accessKey  string
	profileID  string
	secretKey  string
	isTestMode bool
	baseURL    string
	log        logging.Logger
}

// ProcessorConfig holds the configuration for Cybersource processor
type ProcessorConfig struct {
	AccessKey  string
	ProfileID  string
	SecretKey  string
	IsTestMode bool
}

// NewProcessor creates a new Cybersource payment processor
func NewProcessor(config ProcessorConfig) payment.Processor {
	if config.AccessKey == "" {
		config.AccessKey = os.Getenv("CYBERSOURCE_ACCESS_KEY")
	}
	if config.ProfileID == "" {
		config.ProfileID = os.Getenv("CYBERSOURCE_PROFILE_ID")
	}
	if config.SecretKey == "" {
		config.SecretKey = os.Getenv("CYBERSOURCE_SECRET_KEY")
	}

	baseURL := os.Getenv("CYBERSOURCE_BASE_URL")
	if baseURL == "" {
		if config.IsTestMode {
			baseURL = "https://testsecureacceptance.cybersource.com"
		} else {
			baseURL = "https://secureacceptance.cybersource.com"
		}
	}

	return &processor{
		accessKey:  config.AccessKey,
		profileID:  config.ProfileID,
		secretKey:  config.SecretKey,
		isTestMode: config.IsTestMode,
		baseURL:    baseURL,
		log:        logging.NewStdLogger("[CYBERSOURCE] [PROCESSOR]"),
	}
}

func (p *processor) InitiatePayment(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Info("Initiating Cybersource payment", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"amount":         req.Amount,
		"currency":       req.Currency,
	})

	SIGNED_FIELD_NAMES := []string{
		"access_key",
		"amount",
		"bill_to_address_city",
		"bill_to_address_country",
		"bill_to_address_line1",
		"bill_to_address_postal_code",
		"bill_to_address_state",
		"bill_to_email",
		"bill_to_forename",
		"bill_to_surname",
		"bill_to_phone",
		"currency",
		"locale",
		"payment_method",
		"profile_id",
		"reference_number",
		"signed_date_time",
		"signed_field_names",
		"transaction_type",
		"transaction_uuid",
		"unsigned_field_names",
	}

	deviceFingerprint := uuid.New()

	// Prepare request parameters
	reqParams := map[string]string{
		"access_key":                  p.accessKey,
		"amount":                      fmt.Sprintf("%.2f", req.Amount),
		"bill_to_forename":            "NOREAL",
		"bill_to_surname":             "NAME",
		"bill_to_email":               "null@cybersource.com",
		"bill_to_phone":               "6509656000",
		"bill_to_address_line1":       "1295 Charleston rd",
		"bill_to_address_city":        "Mountain View",
		"bill_to_address_state":       "CA",
		"bill_to_address_country":     "US",
		"bill_to_address_postal_code": "94043",
		"currency":                    req.Currency,
		"locale":                      "en-US",
		"payment_method":              "card",
		"profile_id":                  p.profileID,
		"reference_number":            deviceFingerprint.String(),
		"signed_date_time":            time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"signed_field_names":          strings.Join(SIGNED_FIELD_NAMES, ","),
		"transaction_type":            "authorization",
		"transaction_uuid":            req.TransactionID.String(),
		"unsigned_field_names":        "",
	}

	// Generate signature
	signature := sign(reqParams, p.secretKey)
	// Generate signature

	// DEBUG: Log the signature and signed fields
	keys := make([]string, 0, len(reqParams))
	for k := range reqParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var encodedFields []string
	for _, k := range keys {
		encodedFields = append(encodedFields, k+"="+reqParams[k])
	}

	p.log.Info("Generated Signature", map[string]interface{}{
		"signature":            signature,
		"signed_fields":        strings.Join(encodedFields, ","),
		"secret_key_truncated": p.secretKey[:4] + "..." + p.secretKey[len(p.secretKey)-4:],
	})

	// Create HTML form
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Cybersource Payment</title>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/crypto-js/4.0.0/crypto-js.min.js"></script>
	</head>
	<body>
		<form id='payment_form' method='post' action="%s/embedded/pay">
			<input type='hidden' id='access_key' name='access_key' value='%s' />
			<input type='hidden' id='amount' name='amount' value='%s' />
			<input type='hidden' id='bill_to_forename' name='bill_to_forename' value='%s' />
			<input type='hidden' id='bill_to_surname' name='bill_to_surname' value='%s' />
			<input type='hidden' id='bill_to_email' name='bill_to_email' value='%s' />
			<input type='hidden' id='bill_to_phone' name='bill_to_phone' value='%s' />
			<input type='hidden' id='bill_to_address_line1' name='bill_to_address_line1' value='%s' />
			<input type='hidden' id='bill_to_address_city' name='bill_to_address_city' value='%s' />
			<input type='hidden' id='bill_to_address_state' name='bill_to_address_state' value='%s' />
			<input type='hidden' id='bill_to_address_country' name='bill_to_address_country' value='%s' />
			<input type='hidden' id='bill_to_address_postal_code' name='bill_to_address_postal_code' value='%s' />
			<input type='hidden' id='currency' name='currency' value='%s' />
			<input type='hidden' id='locale' name='locale' value='%s' />
			<input type='hidden' id='payment_method' name='payment_method' value='%s' />
			<input type='hidden' id='profile_id' name='profile_id' value='%s' />
			<input type='hidden' id='reference_number' name='reference_number' value='%s' />
			<input type='hidden' id='signature' name='signature' value='%s' />
			<input type='hidden' id='signed_date_time' name='signed_date_time' value='%s' />
			<input type='hidden' id='signed_field_names' name='signed_field_names' value='%s' />
			<input type='hidden' id='transaction_type' name='transaction_type' value='%s' />
			<input type='hidden' id='transaction_uuid' name='transaction_uuid' value='%s' />
			<input type='hidden' id='unsigned_field_names' name='unsigned_field_names' value='%s' />
		</form>
		<script type="text/javascript">
			window.onload = function() {
				document.getElementById("payment_form").submit();
			};
		</script>
	</body>
</html>`,
		p.baseURL,
		reqParams["access_key"],
		reqParams["amount"],
		reqParams["bill_to_forename"],
		reqParams["bill_to_surname"],
		reqParams["bill_to_email"],
		reqParams["bill_to_phone"],
		reqParams["bill_to_address_line1"],
		reqParams["bill_to_address_city"],
		reqParams["bill_to_address_state"],
		reqParams["bill_to_address_country"],
		reqParams["bill_to_address_postal_code"],
		reqParams["currency"],
		reqParams["locale"],
		reqParams["payment_method"],
		reqParams["profile_id"],
		reqParams["reference_number"],
		signature,
		reqParams["signed_date_time"],
		reqParams["signed_field_names"],
		reqParams["transaction_type"],
		reqParams["transaction_uuid"],
		reqParams["unsigned_field_names"],
	)

	// Save the HTML file
	fileName := fmt.Sprintf("%s.html", req.TransactionID)
	filePath := fmt.Sprintf("./public/%s", fileName)
	if err := os.WriteFile(filePath, []byte(htmlContent), 0666); err != nil {
		p.log.Error("Failed to create payment form", map[string]interface{}{
			"error": err.Error(),
			"path":  filePath,
		})
		return nil, fmt.Errorf("failed to create payment form: %v", err)
	}

	// Return the URL to the static file
	paymentURL := fmt.Sprintf("%s/api/v2/static/%s", "https://api.socialpay.co", fileName)

	p.log.Info("Payment form created", map[string]interface{}{
		"url":            paymentURL,
		"transaction_id": req.TransactionID,
	})

	return &payment.PaymentResponse{
		Success:       true,
		TransactionID: req.TransactionID,
		Status:        txEntity.PENDING,
		PaymentURL:    paymentURL,
		Message:       "Redirect user to payment URL",
	}, nil
}

func (p *processor) SettlePayment(ctx context.Context, req *payment.CallbackRequest) error {
	p.log.Info("Processing Cybersource callback", map[string]interface{}{
		"transaction_id": req.TransactionID,
		"status":         req.Status,
	})

	// Extract Cybersource-specific data from metadata
	reasonCode, _ := req.Metadata["reason_code"].(string)
	decision, _ := req.Metadata["decision"].(string)

	p.log.Info("Cybersource callback details", map[string]interface{}{
		"reason_code": reasonCode,
		"decision":    decision,
	})

	// Map Cybersource status
	var status txEntity.TransactionStatus
	switch decision {
	case "ACCEPT":
		if reasonCode == "100" {
			status = txEntity.SUCCESS
		} else {
			status = txEntity.FAILED
		}
	case "CANCEL":
		status = txEntity.FAILED
	case "DECLINE":
		status = txEntity.FAILED
	default:
		status = txEntity.FAILED
	}

	if status == txEntity.SUCCESS {
		return nil
	}

	return fmt.Errorf("payment failed: %s", decision)
}

func (p *processor) GetType() txEntity.TransactionMedium {
	return txEntity.CYBERSOURCE
}

func (p *processor) InitiateWithdrawal(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	p.log.Error("Withdrawal not supported", map[string]interface{}{
		"processor": "Cybersource",
	})
	return nil, fmt.Errorf("withdrawal not supported for Cybersource")
}

func sign(fields map[string]string, secretKey string) string {
	// Field order MUST match SIGNED_FIELD_NAMES exactly
	fieldOrder := []string{
		"access_key",
		"amount",
		"bill_to_address_city",
		"bill_to_address_country",
		"bill_to_address_line1",
		"bill_to_address_postal_code",
		"bill_to_address_state",
		"bill_to_email",
		"bill_to_forename",
		"bill_to_surname",
		"bill_to_phone",
		"currency",
		"locale",
		"payment_method",
		"profile_id",
		"reference_number",
		"signed_date_time",
		"signed_field_names",
		"transaction_type",
		"transaction_uuid",
		"unsigned_field_names",
	}

	var encodedFields []string
	for _, k := range fieldOrder {
		if val, exists := fields[k]; exists {
			encodedFields = append(encodedFields, k+"="+val)
		}
	}

	// Debug log the exact string being signed
	signData := strings.Join(encodedFields, ",")
	fmt.Printf("[DEBUG] Signing data: %s\n", signData)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signData))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Debug log the generated signature
	fmt.Printf("[DEBUG] Generated signature: %s\n", signature)

	return signature
}

func (p *processor) QueryTransactionStatus(ctx context.Context, transactionID string) (*payment.TransactionStatusQueryResponse, error) {
	p.log.Info("Querying Cybersource transaction status", map[string]interface{}{
		"transaction_id": transactionID,
	})
	return nil, nil
}
