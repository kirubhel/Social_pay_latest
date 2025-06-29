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

func ProcessCybersource(id string, amount float64, host string) (string, error) {
	var url string
	var err error

	log.Printf("[CYBERSOURCE] Processing payment for ID: %s, Amount: %v", id, amount)

	UNSIGNED_FIELD_NAMES := []string{}
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

	reqParams := map[string]string{
		"access_key":           "66ad734a971a3f79b84f183c4e52b790",
		"amount":               fmt.Sprintf("%v", amount),
		"currency":             "ETB",
		"locale":               "en-US",
		"payment_method":       "card",
		"profile_id":           "2674CE2D-EA15-4D9F-85D0-713FC5F6329F",
		"reference_number":     uuid.New().String(),
		"signed_date_time":     time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"signed_field_names":   strings.Join(SIGNED_FIELD_NAMES, ","),
		"transaction_type":     "sale",
		"transaction_uuid":     id,
		"unsigned_field_names": strings.Join(UNSIGNED_FIELD_NAMES, ","),
		"bill_to_forename":          "NOREAL",
		"bill_to_surname":          "NAME",
		"bill_to_address_line1":    "1295 Charleston rd",
		"bill_to_address_city":     "Mountain View",
		"bill_to_address_state":    "CA",
		"bill_to_address_postal_code": "94043",
		"bill_to_address_country":  "US",
		"bill_to_email":           "null@cybersource.com",
		"bill_to_phone":           "6509656000",
	}

	walletId := id
	log.Printf("[CYBERSOURCE] Generating payment form for wallet ID: %s", walletId)

	filePath, err := utils.GetPublicFilePath(fmt.Sprintf("%s.html", walletId))
	if err != nil {
		log.Printf("[CYBERSOURCE] Error getting public file path: %v", err)
		return "", fmt.Errorf("failed to get public file path: %w", err)
	}
	log.Printf("[CYBERSOURCE] Will write payment form to: %s", filePath)

	formHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>Cybersource Payment</title>
			</head>
			<body>
				<form id='payment_form' method='post' action="https://testsecureacceptance.cybersource.com/pay">
					<input type='hidden' name='access_key' value='%s'/>
					<input type='hidden' name='amount' value='%s'/>
					<input type='hidden' name='currency' value='%s'/>
					<input type='hidden' name='locale' value='%s'/>
					<input type='hidden' name='payment_method' value='%s'/>
					<input type='hidden' name='profile_id' value='%s'/>
					<input type='hidden' name='reference_number' value='%s'/>
					<input type='hidden' name='signed_date_time' value='%s'/>
					<input type='hidden' name='signed_field_names' value='%s'/>
					<input type='hidden' name='transaction_type' value='%s'/>
					<input type='hidden' name='transaction_uuid' value='%s'/>
					<input type='hidden' name='unsigned_field_names' value='%s'/>
					<input type='hidden' name='signature' value='%s'/>
					<!-- Static billing fields -->
					<input type='hidden' name='bill_to_forename' value='NOREAL'/>
					<input type='hidden' name='bill_to_surname' value='NAME'/>
					<input type='hidden' name='bill_to_address_line1' value='1295 Charleston rd'/>
					<input type='hidden' name='bill_to_address_city' value='Mountain View'/>
					<input type='hidden' name='bill_to_address_state' value='CA'/>
					<input type='hidden' name='bill_to_address_postal_code' value='94043'/>
					<input type='hidden' name='bill_to_address_country' value='US'/>
					<input type='hidden' name='bill_to_email' value='null@cybersource.com'/>
					<input type='hidden' name='bill_to_phone' value='6509656000'/>
				</form>
				<script>
					document.getElementById('payment_form').submit();
				</script>
			</body>
		</html>
		`,
		reqParams["access_key"],
		reqParams["amount"],
		reqParams["currency"],
		reqParams["locale"],
		reqParams["payment_method"],
		reqParams["profile_id"],
		reqParams["reference_number"],
		reqParams["signed_date_time"],
		reqParams["signed_field_names"],
		reqParams["transaction_type"],
		reqParams["transaction_uuid"],
		reqParams["unsigned_field_names"],
		sign(reqParams, "328be89eb0ca4c53845594974c09a17cd4aa0c561bc147f0b12f6fd612cb85a21ddd85b92e0b4ff3b39153c2ad904c9966375e00572f450ba68b19a72c2055a11ea890496b9a4eaab8748fc93d7bef65a37c01d94d6d43c2ab8e7e90314ce098c6978a1d2ceb4e17b0a995d46f90099676bb45f923e64258b7d6d856e00487cc"),
	)

	log.Printf("[CYBERSOURCE] Writing payment form to file...")
	if err := os.WriteFile(filePath, []byte(formHTML), 0644); err != nil {
		log.Printf("[CYBERSOURCE] Error writing payment form: %v", err)
		return "", fmt.Errorf("failed to write payment form: %w", err)
	}
	log.Printf("[CYBERSOURCE] Successfully wrote payment form to: %s", filePath)

	url = fmt.Sprintf("%s/static/%s.html", host, walletId)
	log.Printf("[CYBERSOURCE] Generated payment URL: %s", url)

	return url, nil
}
