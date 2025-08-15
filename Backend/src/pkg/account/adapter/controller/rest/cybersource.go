package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func (controller Controller) CybersourceTransactionStatus(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	envFilePath := ".env"
	if err := godotenv.Load(envFilePath); err != nil {
		controller.log.Println("Warning: .env file not found, proceeding without it.")
	}

	forwardURL := os.Getenv("SOCIALPAY_API_URL")
	if forwardURL == "" {
		controller.log.Println("SOCIALPAY_API_URL is not set in .env file")
		http.Error(w, "SOCIALPAY_API_URL is not configured", http.StatusInternalServerError)
		return
	}

	controller.log.Println("||||||| Received Cybersource Transaction Status Callback ||||||||")

	// Parse form data
	if err := r.ParseForm(); err != nil {
		controller.log.Printf("Failed to parse form data: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Log raw callback data
	controller.log.Println("Raw Callback Data:")
	for key, values := range r.Form {
		controller.log.Printf("%s: %v\n", key, values)
	}

	// Extract key parameters
	referenceId := r.Form.Get("req_transaction_uuid") // Our internal transaction UUID
	if referenceId == "" {
		// Fallback to reference number if UUID not available
		referenceId = r.Form.Get("req_reference_number")
		controller.log.Printf("Warning: Using req_reference_number as fallback for referenceId: %s", referenceId)
	}

	providerTxId := r.Form.Get("transaction_id") // Cybersource's internal transaction ID
	if providerTxId == "" {
		// Fallback to authorization reference if transaction_id not available
		providerTxId = r.Form.Get("auth_trans_ref_no")
		controller.log.Printf("Warning: Using auth_trans_ref_no as fallback for providerTxId: %s", providerTxId)
	}

	decision := r.Form.Get("decision")
	reasonCode := r.Form.Get("reason_code")

	// Marshal the original form data for providerData
	formDataBytes, err := json.Marshal(r.Form)
	if err != nil {
		controller.log.Printf("Failed to marshal form data: %v", err)
		http.Error(w, "Failed to process callback data", http.StatusInternalServerError)
		return
	}
	providerData := string(formDataBytes)

	// Determine transaction status
	status := "FAILURE"
	if decision == "ACCEPT" && reasonCode == "100" {
		status = "SUCCESS"
	}

	// Prepare payload for forwarding
	newPayload := map[string]interface{}{
		"referenceId":  referenceId,     // Our internal transaction reference
		"status":       status,          // SUCCESS/FAILURE
		"message":      fmt.Sprintf("Cybersource: %s (Code: %s)", decision, reasonCode),
		"providerTxId": providerTxId,    // Cybersource's transaction ID
		"providerData": providerData,    // Original callback data
		"timestamp":    time.Now().Format(time.RFC3339),
		"type":         "CYBERSOURCE",
	}

	// Marshal payload for forwarding
	payloadBytes, err := json.Marshal(newPayload)
	if err != nil {
		controller.log.Printf("Error marshalling payload: %v", err)
		http.Error(w, "Failed to prepare forwarding payload", http.StatusInternalServerError)
		return
	}

	// Debug log the payload
	controller.log.Printf("Forwarding payload: %s", string(payloadBytes))

	// Create and send forward request
	req, err := http.NewRequest("POST", forwardURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		controller.log.Printf("Failed to create forward request: %v", err)
		http.Error(w, "Failed to create forward request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	forwardResponse, err := client.Do(req)
	if err != nil {
		controller.log.Printf("Forwarding failed: %v", err)
		http.Error(w, "Failed to forward callback", http.StatusInternalServerError)
		return
	}
	defer forwardResponse.Body.Close()

	// Respond to Cybersource
	response := map[string]interface{}{
		"success": true,
		"message": "Callback processed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Printf("Failed to encode response: %v", err)
	}
}
