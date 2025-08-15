package gin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	"github.com/socialpay/socialpay/src/pkg/shared/payment/etswitch"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	"github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
	usecase "github.com/socialpay/socialpay/src/pkg/webhook/usecase"
)

// Package gin provides HTTP handlers for payment settlement
// @title Payment Settlement API
// @version 1.0
// @description API for handling payment settlement callbacks from various payment providers
// @BasePath /

type SettlementHandler struct {
	log        logging.Logger
	processors map[txEntity.TransactionMedium]payment.Processor
	usecase    usecase.WebhookUseCase
}

func NewSettlementHandler(processors map[txEntity.TransactionMedium]payment.Processor, usecase usecase.WebhookUseCase) *SettlementHandler {
	return &SettlementHandler{
		log:        logging.NewStdLogger("[SETTLEMENT] [HANDLER]"),
		processors: processors,
		usecase:    usecase,
	}
}

func (h *SettlementHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/settle/std", h.HandleSettlement)
	router.POST("/settle/mpesa", h.HandleMPESASettlement)
	router.POST("/settle/telebirr", h.HandleTelebirrSettlement)
	router.POST("/settle/cbe", h.HandleCBESettlement)
	router.POST("/settle/cybersource", h.HandleCybersourceSettlement)
	router.POST("/settle/awash", h.HandleAwashSettlement)
	router.POST("/settle/epg/payment", h.HandleEthSwitchSettlement)
}

// @Summary Handle M-PESA payment settlement
// @Description Process M-PESA payment callback and update transaction status
// @Tags Settlement
// @Accept json
// @Produce json
// @Param callback body MPESACallback true "M-PESA callback data"
// @Success 200 {object} SettlementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /settle/mpesa [post]
func (h *SettlementHandler) HandleMPESASettlement(c *gin.Context) {
	processor, exists := h.processors[txEntity.MPESA]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MPESA processor not configured"})
		return
	}

	// Parse M-PESA specific callback
	var mpesaCallback struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"`
		Message       string `json:"message"`
	}
	if err := c.BindJSON(&mpesaCallback); err != nil {
		h.log.Error("Failed to decode M-PESA callback", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback data"})
		return
	}

	transactionID, err := uuid.Parse(mpesaCallback.TransactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}
	transactionStatus := txEntity.FAILED
	if mpesaCallback.Status == "0" {
		transactionStatus = txEntity.SUCCESS
	}

	// Create generic callback request
	callbackReq := &payment.CallbackRequest{
		TransactionID: transactionID,
		Status:        transactionStatus,
		Metadata:      map[string]interface{}{"message": mpesaCallback.Message},
	}

	if err := processor.SettlePayment(c.Request.Context(), callbackReq); err != nil {
		h.log.Error("Settlement failed", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Settlement failed"})
		return
	}
	// Dispatch webhook
	h.usecase.HandleWebhookDispatch(c.Request.Context(), dto.WebhookRequest{
		TransactionID: callbackReq.TransactionID.String(),
		Status:        string(callbackReq.Status),
		Message:       callbackReq.Metadata["message"].(string),
		ProviderTxID:  callbackReq.ProcessorRef,
		ProviderData:  callbackReq.Metadata["message"].(string),
		Timestamp:     time.Now(),
	})
	c.JSON(http.StatusOK, gin.H{"message": "Settlement processed successfully"})
}

// @Summary Handle Telebirr payment settlement
// @Description Process Telebirr payment callback and update transaction status
// @Tags Settlement
// @Accept json
// @Produce json
// @Param callback body TelebirrCallback true "Telebirr callback data"
// @Success 200 {object} SettlementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /settle/telebirr [post]
func (h *SettlementHandler) HandleTelebirrSettlement(c *gin.Context) {
	processor, exists := h.processors[txEntity.TELEBIRR]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telebirr processor not configured"})
		return
	}

	// Parse Telebirr specific callback
	var telebirrCallback struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"`
		Message       string `json:"message"`
	}
	if err := c.BindJSON(&telebirrCallback); err != nil {
		h.log.Error("Failed to decode Telebirr callback", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback data"})
		return
	}

	transactionID, err := uuid.Parse(telebirrCallback.TransactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}

	transactionStatus := txEntity.FAILED
	if telebirrCallback.Status == "0" {
		transactionStatus = txEntity.SUCCESS
	}

	// Create generic callback request
	callbackReq := &payment.CallbackRequest{
		TransactionID: transactionID,
		Status:        transactionStatus,
		Metadata:      map[string]interface{}{"message": telebirrCallback.Message},
	}

	if err := processor.SettlePayment(c.Request.Context(), callbackReq); err != nil {
		h.log.Error("Settlement failed", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Settlement failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settlement processed successfully"})
}

// @Summary Handle CBE payment settlement
// @Description Process CBE payment callback and update transaction status
// @Tags Settlement
// @Accept json
// @Produce json
// @Param callback body payment.STDCallbackRequest true "CBE callback data"
// @Success 200 {object} SettlementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /settle/cbe [post]
func (h *SettlementHandler) HandleCBESettlement(c *gin.Context) {

	// Parse CBE specific callback
	var cbeCallback payment.STDCallbackRequest
	if err := c.BindJSON(&cbeCallback); err != nil {
		h.log.Error("Failed to decode CBE callback", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback data"})
		return
	}

	transactionID, err := uuid.Parse(cbeCallback.ReferenceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}
	transactionStatus := txEntity.FAILED
	if cbeCallback.Status == "SUCCESS" {
		transactionStatus = txEntity.SUCCESS
	}

	// Dispatch webhook
	h.usecase.HandleWebhookDispatch(c.Request.Context(), dto.WebhookRequest{
		TransactionID: transactionID.String(),
		Status:        string(transactionStatus),
		Message:       cbeCallback.Message,
		ProviderTxID:  cbeCallback.ProviderTxId,
		ProviderData:  cbeCallback.ProviderData,
		Timestamp:     time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Settlement processed successfully"})
}

// @Summary Handle Awash settlement
// @Description Process Awash settlement callback
// @Tags Settlement
// @Accept json
// @Produce json
// @Param callback body AwashCallback true "Awash settlement callback data"
// @Success 200 {object} SettlementResponse
// @Success 501 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /settle/awash [post]
func (h *SettlementHandler) HandleAwashSettlement(c *gin.Context) {

	// TODO Check if the callback come from Awash
	// Using the awash origin header
	// Log raw request body
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Failed to read raw request body", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request"})
		return
	}
	// Restore the request body for later binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// Log the raw request details
	h.log.Info("Received Awash settlement callback", map[string]interface{}{
		"raw_body":    string(rawBody),
		"method":      c.Request.Method,
		"url":         c.Request.URL.String(),
		"headers":     c.Request.Header,
		"remote_addr": c.Request.RemoteAddr,
	})

	var AwashCallback AwashCallback
	if err := c.BindJSON(&AwashCallback); err != nil {
		h.log.Error("Failed to decode Awash callback", map[string]interface{}{
			"error":    err.Error(),
			"raw_body": string(rawBody),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback data"})
		return
	}
	h.log.Info("Successfully parsed Awash callback data", map[string]interface{}{
		"transaction_id":     AwashCallback.TransactionID,
		"external_reference": AwashCallback.ExternalReference,
		"status":             AwashCallback.Status,
		"payer_phone":        AwashCallback.PayerPhone,
		"return_code":        AwashCallback.ReturnCode,
		"return_message":     AwashCallback.ReturnMessage,
	})

	// Validate transaction ID format
	transactionID, err := uuid.Parse(AwashCallback.ExternalReference)
	if err != nil {
		h.log.Error("Invalid transaction ID format", map[string]interface{}{
			"external_reference": AwashCallback.ExternalReference,
			"error":              err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}
	h.log.Info("Parsed transaction ID", map[string]interface{}{
		"transaction_id": transactionID.String(),
	})

	var transactionStatus txEntity.TransactionStatus
	if AwashCallback.ReturnCode == 0 {
		transactionStatus = txEntity.SUCCESS
	}

	// Determine transaction status based on Awash callback
	if AwashCallback.ReturnCode != 0 {
		h.log.Error("Awash settlement failed", map[string]interface{}{
			"external_reference": AwashCallback.ExternalReference,
			"return_code":        AwashCallback.ReturnCode,
			"return_message":     AwashCallback.ReturnMessage,
		})

		transactionStatus = txEntity.FAILED
	}

	// Webhook
	h.usecase.HandleWebhookDispatch(c.Request.Context(), dto.WebhookRequest{

		TransactionID: transactionID.String(),
		Status:        string(transactionStatus),
		Message:       AwashCallback.ReturnMessage,
		ProviderTxID:  AwashCallback.TransactionID,
		ProviderData:  string(rawBody),
		Timestamp:     time.Now(),
	})
	// Log successful processing
	h.log.Info("Awash settlement processing completed", map[string]interface{}{
		"transaction_id": transactionID.String(),
		"status":         transactionStatus,
	})

	// Write
	c.JSON(http.StatusOK, gin.H{"message": "Awash Settlement processed successfully"})

}

// @Summary Handle payment settlement
// @Description Process payment callback and update transaction status
// @Tags Settlement
// @Accept json
// @Produce json
// @Param callback body payment.STDCallbackRequest true "Settlement callback data"
// @Success 200 {object} SettlementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /settle/std [post]
func (h *SettlementHandler) HandleSettlement(c *gin.Context) {
	// Log raw request body
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Failed to read raw request body", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request"})
		return
	}
	// Restore the request body for later binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// Log the raw request details
	h.log.Info("Received settlement callback", map[string]interface{}{
		"raw_body":    string(rawBody),
		"method":      c.Request.Method,
		"url":         c.Request.URL.String(),
		"headers":     c.Request.Header,
		"remote_addr": c.Request.RemoteAddr,
	})

	// Parse CBE specific callback
	var stdCallback payment.STDCallbackRequest
	if err := c.BindJSON(&stdCallback); err != nil {
		h.log.Error("Failed to decode callback", map[string]interface{}{
			"error":      err.Error(),
			"raw_body":   string(rawBody),
			"request_id": c.GetHeader("X-Request-ID"),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback data"})
		return
	}

	// Log parsed callback data
	h.log.Info("Successfully parsed callback data", map[string]interface{}{
		"reference_id": stdCallback.ReferenceId,
		"status":       stdCallback.Status,
		"provider_id":  stdCallback.ProviderTxId,
		"message":      stdCallback.Message,
		"request_id":   c.GetHeader("X-Request-ID"),
	})

	transactionID, err := uuid.Parse(stdCallback.ReferenceId)
	if err != nil {
		h.log.Error("Invalid transaction ID format", map[string]interface{}{
			"reference_id": stdCallback.ReferenceId,
			"error":        err.Error(),
			"request_id":   c.GetHeader("X-Request-ID"),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}

	h.log.Info("Parsed transaction ID", map[string]interface{}{
		"transaction_id": transactionID.String(),
		"request_id":     c.GetHeader("X-Request-ID"),
	})

	transactionStatus := txEntity.FAILED
	if stdCallback.Status == "SUCCESS" {
		transactionStatus = txEntity.SUCCESS
	}

	h.log.Info("Determined transaction status", map[string]interface{}{
		"transaction_id": transactionID.String(),
		"raw_status":     stdCallback.Status,
		"final_status":   transactionStatus,
		"request_id":     c.GetHeader("X-Request-ID"),
	})

	// Prepare webhook request
	webhookReq := dto.WebhookRequest{
		TransactionID: transactionID.String(),
		Status:        string(transactionStatus),
		Message:       stdCallback.Message,
		ProviderTxID:  stdCallback.ProviderTxId,
		ProviderData:  stdCallback.ProviderData,
		Timestamp:     time.Now(),
	}

	h.log.Info("Dispatching webhook", map[string]interface{}{
		"webhook_request": webhookReq,
		"request_id":      c.GetHeader("X-Request-ID"),
	})

	// Dispatch webhook
	if err := h.usecase.HandleWebhookDispatch(c.Request.Context(), webhookReq); err != nil {
		h.log.Error("Failed to dispatch webhook", map[string]interface{}{
			"error":      err.Error(),
			"request_id": c.GetHeader("X-Request-ID"),
		})
		// Note: We don't return error to client as the payment was processed
		// But we log it for debugging
	}

	h.log.Info("Settlement processing completed", map[string]interface{}{
		"transaction_id": transactionID.String(),
		"status":         transactionStatus,
		"request_id":     c.GetHeader("X-Request-ID"),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Standard Settlement processed successfully"})
}

// @Summary Handle Cybersource payment settlement
// @Description Process Cybersource payment callback and update transaction status
// @Tags Settlement
// @Accept x-www-form-urlencoded
// @Produce json
// @Param req_transaction_uuid formData string true "Transaction UUID"
// @Param reason_code formData string true "Reason code from Cybersource"
// @Param decision formData string true "Decision from Cybersource (ACCEPT/DECLINE/CANCEL)"
// @Param signed_field_names formData string true "List of signed fields"
// @Param signature formData string true "Request signature"
// @Success 200 {object} SettlementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /settle/cybersource [post]
func (h *SettlementHandler) HandleCybersourceSettlement(c *gin.Context) {
	processor, exists := h.processors[txEntity.CYBERSOURCE]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cybersource processor not configured"})
		return
	}

	// Parse Cybersource form data
	if err := c.Request.ParseForm(); err != nil {
		h.log.Error("Failed to parse form data", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	transactionID := c.Request.Form.Get("req_transaction_uuid")
	reasonCode := c.Request.Form.Get("reason_code")
	decision := c.Request.Form.Get("decision")
	rawBody := c.Request.PostForm.Encode()
	parsedTransactionID, err := uuid.Parse(transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}

	// Map Cybersource status
	transactionStatus := txEntity.FAILED

	switch decision {
	case "ACCEPT":
		if reasonCode == "100" {
			transactionStatus = txEntity.SUCCESS
		} else {
			transactionStatus = txEntity.FAILED
		}
	case "CANCEL":
		transactionStatus = txEntity.CANCELED
	case "DECLINE":
		transactionStatus = txEntity.FAILED
	default:
		transactionStatus = txEntity.FAILED
	}

	// Create generic callback request
	callbackReq := &payment.CallbackRequest{
		TransactionID: parsedTransactionID,
		Status:        transactionStatus,
		Metadata: map[string]interface{}{
			"reason_code": reasonCode,
			"decision":    decision,
		},
	}

	if err := processor.SettlePayment(c.Request.Context(), callbackReq); err != nil {
		h.log.Error("Settlement failed", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Settlement failed"})
		return
	}

	// Dispatch webhook
	h.usecase.HandleWebhookDispatch(c.Request.Context(), dto.WebhookRequest{
		TransactionID: callbackReq.TransactionID.String(),
		Status:        string(callbackReq.Status),
		Message:       decision,
		ProviderTxID:  transactionID,
		ProviderData:  rawBody,
		Timestamp:     time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Settlement processed successfully"})
}

// @Summary Handle the ethswitch settlemnet
// @Description process the ethswtich callback
// @Tags Settlement
// @Accept json
// @Produce json
// @Param callback body EthSwitchCallback true "ethswitch callback body"
// @Success 200 {object} SettlementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /settle/epg/payment [post]
func (h *SettlementHandler) HandleEthSwitchSettlement(c *gin.Context) {

	processor, exists := h.processors[txEntity.ETHSWITCH]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EthSwitch processor not configured"})
		return
	}

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Failed to read raw request body from ethswithc callback", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request"})
		return
	}
	// Restore the request body for later binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// Log the raw request details
	h.log.Info("Received EthSwitch settlement callback", map[string]interface{}{
		"raw_body":    string(rawBody),
		"method":      c.Request.Method,
		"url":         c.Request.URL.String(),
		"headers":     c.Request.Header,
		"remote_addr": c.Request.RemoteAddr,
	})

	var res EthSwitchCallback
	if err := c.BindJSON(&res); err != nil {
		h.log.Error("Failed to decode EthSwitch callback", map[string]interface{}{
			"error":    err.Error(),
			"raw_body": string(rawBody),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback data"})
		return
	}
	h.log.Info("Successfully parsed EthSwitch callback data", map[string]interface{}{
		"transaction_id": res.OrderNumber,
		"status":         res.Operation,
	})

	// Validate transaction ID format
	transactionID, err := ParseOrderNumberToUUID(res.OrderNumber)
	if err != nil {
		h.log.Error("Invalid transaction ID format", map[string]interface{}{
			"id":        res.OrderNumber,
			"operation": "parsing id to uuid",
			"error":     err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID format"})
		return
	}
	h.log.Info("Parsed transaction ID", map[string]interface{}{
		"transaction_id": transactionID.String(),
	})

	// Safe mapping with fallback
	status, ok := etswitch.ETHStatusToConstant[res.Operation]
	if !ok {
		h.log.Warn("Received unknown EthSwitch operation, defaulting to FAILED", map[string]interface{}{
			"operation": res.Operation,
		})
		status = "FAILED"
	}

	callbackReq := &payment.CallbackRequest{
		TransactionID: transactionID,
		Status:        txEntity.TransactionStatus(status),
		Metadata: map[string]interface{}{
			"reason_code": res.Status,
			"decision":    status,
		},
	}

	if err := processor.SettlePayment(c.Request.Context(), callbackReq); err != nil {
		h.log.Error("Settlement failed", map[string]interface{}{"error": err.Error()})

		// Updating the transaction status for failed scenario
		h.usecase.ProcessTransactionStatus(c.Request.Context(), transactionID, status)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Settlement failed"})
		return
	}

	// Dispatch webhook
	h.usecase.HandleWebhookDispatch(c.Request.Context(), dto.WebhookRequest{
		TransactionID:    callbackReq.TransactionID.String(),
		Status:           string(callbackReq.Status),
		Message:          "EthSwitch settlement", // Update
		ProviderTxID:     res.MdOrder,
		ProviderData:     string(rawBody),
		Timestamp:        time.Now(),
		IsHostedCheckout: true,
	})

	// response
	c.JSON(http.StatusOK, gin.H{"message": "EthSwitch Settlement processed successfully"})
}

func ParseOrderNumberToUUID(an132 string) (uuid.UUID, error) {
	// Insert hyphens at standard UUID positions
	if len(an132) != 32 {
		return uuid.Nil, fmt.Errorf("invalid AN1.32 length")
	}
	formatted := fmt.Sprintf("%s-%s-%s-%s-%s",
		an132[0:8],
		an132[8:12],
		an132[12:16],
		an132[16:20],
		an132[20:32],
	)
	return uuid.Parse(formatted)
}

// EthSwitchCallback represents the callback data from ethswitch
type EthSwitchCallback struct {
	PaymentWay  string `jsonL:"PaymentWay"`
	OrderNumber string `json:"orderNumber"`
	MdOrder     string `json:"mdOrder"`
	Operation   string `json:"operation"`
	Status      string `json:"status"`
}

// AwashCallback represents the callback data from Awash
type AwashCallback struct {
	TransactionID     string  `json:"transactionId,omitempty"`
	Amount            float64 `json:"amount"`
	DateRequested     string  `json:"dateRequested"`
	DateApproved      string  `json:"dateApproved,omitempty"`
	ExternalReference string  `json:"externalReference"`
	PayerPhone        string  `json:"payerPhone"`
	ReturnCode        int     `json:"returnCode"`
	ReturnMessage     string  `json:"returnMessage"`
	Status            string  `json:"status"`
}

// MPESACallback represents the callback data from M-PESA
type MPESACallback struct {
	TransactionID string `json:"transaction_id" example:"123e4567-e89b-12d3-a456-426614174000" binding:"required"`
	Status        string `json:"status" example:"0" binding:"required" enums:"0,1"`
	Message       string `json:"message" example:"Payment successful"`
}

// TelebirrCallback represents the callback data from Telebirr
type TelebirrCallback struct {
	TransactionID string `json:"transaction_id" example:"123e4567-e89b-12d3-a456-426614174000" binding:"required"`
	Status        string `json:"status" example:"0" binding:"required" enums:"0,1"`
	Message       string `json:"message" example:"Payment successful"`
}

// SettlementResponse represents the response for settlement endpoints
type SettlementResponse struct {
	Message string `json:"message" example:"Settlement processed successfully"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid transaction ID format"`
}
