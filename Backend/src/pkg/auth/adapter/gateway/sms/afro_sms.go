package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type AfroSMS struct {
	log       *log.Logger
	Name      string
	AccessKey string
	SecretKey string
	From      string
	URL       string
}

func New(log *log.Logger) AfroSMS {
	return AfroSMS{
		log:       log,
		Name:      "Social Pay",
		AccessKey: "cc4025404f1b43d5a066d003b6b816ce",
		SecretKey: "4a5ea5ca21984ad7a40a995a66c561d3",
		From:      "Social Pay",
		URL:       "http://api.kmicloud.com/sms/send/v1/otp",
	}
}

func (sms AfroSMS) SendSMS(phone, message string) error {
	sms.log.Println("Send SMS via KMI Cloud")

	// Format phone number to use proper country code
	formattedPhone := sms.formatPhoneNumber(phone)
	if formattedPhone == "" {
		errMsg := "invalid phone number format"
		sms.log.Println(errMsg)
		return errors.New(errMsg)
	}

	type SMSBody struct {
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
		From      string `json:"from,omitempty"`
		To        string `json:"to"`
		Message   string `json:"message"`
	}

	sms.log.Println("Recipient:", formattedPhone)
	sms.log.Println("Message:", message)

	sms.log.Printf("[AfroSMS] AccessKey: %s", sms.AccessKey)
	sms.log.Printf("[AfroSMS] SecretKey: %s", sms.SecretKey)
	sms.log.Printf("[AfroSMS] From: %s", sms.From)
	sms.log.Printf("[AfroSMS] To: %s", formattedPhone)
	sms.log.Printf("[AfroSMS] Message: %s", message)

	body := SMSBody{
		AccessKey: sms.AccessKey,
		SecretKey: sms.SecretKey,
		From:      sms.From,
		To:        formattedPhone,
		Message:   message,
	}

	serBody, err := json.Marshal(body)
	if err != nil {
		sms.log.Println("Error marshaling request:", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, sms.URL, bytes.NewReader(serBody))
	if err != nil {
		sms.log.Println("Error creating request:", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		sms.log.Println("Error sending request:", err)
		return err
	}
	defer res.Body.Close()

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Code    int    `json:"code"`
		Result  struct {
			To    string `json:"to"`
			SmsID string `json:"smsId"`
		} `json:"result"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		sms.log.Println("Error decoding response:", err)
		return err
	}

	if !response.Success {
		sms.log.Println("API Error:", response.Message)
		return errors.New(response.Message)
	}

	sms.log.Println("SMS sent successfully. ID:", response.Result.SmsID)
	return nil
}

// converts phone numbers to KMI Cloud's required format
func (sms AfroSMS) formatPhoneNumber(phone string) string {
	cleaned := strings.ReplaceAll(phone, "+", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	if strings.HasPrefix(cleaned, "251") && len(cleaned) > 3 {
		return "00251" + cleaned[3:]
	}

	return cleaned
}
