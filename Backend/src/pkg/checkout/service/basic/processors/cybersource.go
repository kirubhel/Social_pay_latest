package processors

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/utils"
)

// sign creates the HMAC SHA-256 signature from signed fields
func sign(fields map[string]string, secretKey string) string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var encodedFields []string
	for _, k := range keys {
		encodedFields = append(encodedFields, k+"="+fields[k])
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(strings.Join(encodedFields, ",")))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature
}

// generateHiddenFields converts map[string]string to HTML <input type="hidden"> tags
func generateHiddenFields(fields map[string]string) string {
	var htmlFields strings.Builder
	for k, v := range fields {
		htmlFields.WriteString(fmt.Sprintf("<input type='hidden' name='%s' value='%s'/>\n", k, v))
	}
	return htmlFields.String()
}

// ProcessCybersource generates the payment form HTML and saves it to a public file
func ProcessCybersource(id string, amount float64, host string) (string, error) {
	log.Printf("[CYBERSOURCE] Processing socialpay payment for ID: %s, Amount: %v", id, amount)

	// Define signed and unsigned fields
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

	// Define parameters
	reqParams := map[string]string{
		"access_key":                  "cef3ef87ed743dec8563b723ac403da7",
		"profile_id":                  "F39EE97D-B53B-46DF-924E-8EDF122480F5",
		"transaction_uuid":            id,
		"amount":                      fmt.Sprintf("%.2f", amount),
		"currency":                    "USD",
		"locale":                      "en-US",
		"payment_method":              "card",
		"reference_number":            uuid.New().String(),
		"signed_date_time":            time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"signed_field_names":          strings.Join(SIGNED_FIELD_NAMES, ","),
		"unsigned_field_names":        "",
		"transaction_type":            "sale",
		"bill_to_forename":            "NOREAL",
		"bill_to_surname":             "NAME",
		"bill_to_address_line1":       "1295 Charleston rd",
		"bill_to_address_city":        "Mountain View",
		"bill_to_address_state":       "CA",
		"bill_to_address_postal_code": "94043",
		"bill_to_address_country":     "US",
		"bill_to_email":               "null@cybersource.com",
	}

	// Prepare signature
	signedFields := make(map[string]string)
	for _, k := range SIGNED_FIELD_NAMES {
		signedFields[k] = reqParams[k]
	}
	signature := sign(signedFields, "fafe570b811d421aa19ef9530b0c5dda5dd02f61dd1a4b0184a0f26addd4109c657a7d12df634da0a42211b94a821bc4156a4ea2f0b44ec49486302e1ab24d5850ce4980bc46496ea1e90fe19686609543c49814b84940c0ae2970eeb301bb2442eb1a6c726f4ac39a1beba1c49a02b83760625c4f744420876ed98e9b945214")

	// Prepare form HTML
	formHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Cybersource Payment</title>
</head>
<body>
	<form id="payment_form" method="post" action="https://secureacceptance.cybersource.com/pay">
		%s
		<input type="hidden" name="signature" value="%s"/>
	</form>
	<script>
		document.getElementById("payment_form").submit();
	</script>
</body>
</html>
`, generateHiddenFields(reqParams), signature)

	// Write HTML to file
	walletId := id
	filePath, err := utils.GetPublicFilePath(fmt.Sprintf("%s.html", walletId))
	if err != nil {
		log.Printf("[CYBERSOURCE] Error getting public file path: %v", err)
		return "", fmt.Errorf("failed to get public file path: %w", err)
	}
	if err := os.WriteFile(filePath, []byte(formHTML), 0644); err != nil {
		log.Printf("[CYBERSOURCE] Error writing payment form: %v", err)
		return "", fmt.Errorf("failed to write payment form: %w", err)
	}

	// Return public URL to HTML form
	url := fmt.Sprintf("%s/static/%s.html", host, walletId)
	log.Printf("[CYBERSOURCE] Generated payment URL: %s", url)
	return url, nil
}
