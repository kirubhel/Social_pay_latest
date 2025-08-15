package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apikeyEntity "github.com/socialpay/socialpay/src/pkg/apikey_mgmt/core/entity"
	qrEntity "github.com/socialpay/socialpay/src/pkg/qr/core/entity"
	qrUsecase "github.com/socialpay/socialpay/src/pkg/qr/usecase"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
	middleware "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	"github.com/socialpay/socialpay/src/pkg/socialpayapi/core/entity"
	"github.com/socialpay/socialpay/src/pkg/socialpayapi/usecase"
	txnEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	v2MerchantRepo "github.com/socialpay/socialpay/src/pkg/v2_merchant/core/repository"
	settlementdto "github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
	webhookusecase "github.com/socialpay/socialpay/src/pkg/webhook/usecase"
)

var supportedMediumsDeposit = []txnEntity.TransactionMedium{txnEntity.MPESA, txnEntity.CBE, txnEntity.AWASH,
	txnEntity.TELEBIRR, txnEntity.CYBERSOURCE, txnEntity.ETHSWITCH, txnEntity.KACHA}
var supportedMediumsWithdrawal = []txnEntity.TransactionMedium{txnEntity.CBE, txnEntity.TELEBIRR, txnEntity.CYBERSOURCE, txnEntity.KACHA, txnEntity.MPESA}

// @title           SocialPay Payment API
// @version         1.0
// @description     Payment gateway API for processing payments and withdrawals
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.socialpay.com/support
// @contact.email  support@socialpay.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.socialpay.com
// @BasePath  /payment

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API key authentication

// @tag.name Payments
// @tag.description Payment processing operations including direct payments and withdrawals

// @tag.name Transactions
// @tag.description Transaction management and retrieval operations

// Handler handles HTTP requests for the payment API
type Handler struct {
	paymentUseCase usecase.PaymentUseCase
	middleware     *gin.HandlerFunc
	log            logging.Logger
	ipChecker      ginn.IPCheckerMiddleware
	merchantRepo   v2MerchantRepo.Repository
	qrUseCase      qrUsecase.QRUseCase
	webhookUseCase webhookusecase.WebhookUseCase
}

// NewHandler creates a new payment API handler
func NewHandler(
	uc usecase.PaymentUseCase,
	apiAuth gin.HandlerFunc,
	merchantRepo v2MerchantRepo.Repository,
	ipChecker ginn.IPCheckerMiddleware,
	qrUseCase qrUsecase.QRUseCase,
	webhookUseCase webhookusecase.WebhookUseCase) *Handler {
	return &Handler{
		paymentUseCase: uc,
		middleware:     &apiAuth,
		log:            logging.NewStdLogger("[SOCIALPAY-API]"),
		merchantRepo:   merchantRepo,
		ipChecker:      ipChecker,
		qrUseCase:      qrUseCase,
		webhookUseCase: webhookUseCase,
	}
}

// RegisterRoutes registers the API routes to the gin engine
func (h *Handler) RegisterRoutes(r gin.IRouter) {
	api := r.Group("/payment")
	{
		// Payment processing endpoints require payment processing permission
		api.POST("/direct", *h.middleware, middleware.RequirePaymentProcessingPermission(), h.DirectPay)
		api.POST("/checkout", *h.middleware, middleware.RequirePaymentProcessingPermission(), h.Checkout)
		api.PATCH("/checkout/:id", *h.middleware, middleware.RequirePaymentProcessingPermission(), h.UpdateCheckout)

		api.GET("/transaction/:id", h.GetTransaction)
		// Withdrawal endpoints require withdrawal permission
		api.POST("/withdrawal", *h.middleware, middleware.RequireWithdrawalPermission(), h.RequestWithdrawal)
	}

	// Checkout payment endpoint (no authentication required for hosted checkout)
	checkout := r.Group("/checkout", middleware.CheckoutCors)
	{
		checkout.GET("/:id", h.GetHostedCheckout)
		checkout.POST("/makepayment", h.CheckoutPayment)
	}
}

// RegisterQRRoutes registers the QR payment routes on the main router (without /api/v2 prefix)
func (h *Handler) RegisterQRRoutes(r gin.IRouter) {
	// Test endpoint to verify route registration
	r.GET("/qr/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "QR routes are working"})
	})

	qr := r.Group("/qr/payment", middleware.CheckoutCors)
	{
		// QR Payment endpoints (public for customer use)
		qr.GET("/link/:id", h.GetQRLinkForPayment)
		qr.POST("/link/:id", h.ProcessQRPayment)

		qr.POST("/merchant", h.QRMerchantPayment)
	}

	// QR callback endpoint (no CORS restrictions)
	r.POST("/qr/callback", h.QRCallback)
	r.GET("/payment/receipt/:id", h.GetReceipt)
}

// DirectPay godoc
// @Summary      Process a direct payment
// @Description  Process a direct payment transaction with the specified details
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request body entity.DirectPaymentRequest true "Payment request details"
// @Success      200  {object}  entity.PaymentResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     ApiKeyAuth
// @Router       /payment/direct [post]
func (h *Handler) DirectPay(c *gin.Context) {
	h.log.Info("Starting direct payment processing", map[string]interface{}{})
	// Read the raw request body
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Failed to read request body", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		h.log.Info("Raw request body", map[string]interface{}{
			"body": string(rawBody),
		})
	}

	// Restore the request body for later binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	var req entity.DirectPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Failed to bind JSON request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	if !slices.Contains(supportedMediumsDeposit, req.Medium) {
		h.log.Error("Unsupported medium", map[string]interface{}{
			"error": "Unsupported medium",
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("Unsupported medium")))
		return
	}

	h.log.Info("Received direct payment request", map[string]interface{}{
		"medium":       req.Medium,
		"amount":       req.Amount,
		"currency":     req.Currency,
		"callback_url": req.CallbackURL,
	})

	if err := req.Validate(); err != nil {
		h.log.Error("Request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	h.log.Info("Request validation successful, processing payment", map[string]interface{}{})

	apiKeyData, _ := c.Get("apiKey")
	apiKey, _ := apiKeyData.(*apikeyEntity.APIKeyResponse)
	userID := apiKey.UserID
	merchantID := apiKey.MerchantID
	ctx := c.Request.Context()
	apiKeyHeader := c.GetHeader("X-API-Key")
	resp, err := h.paymentUseCase.ProcessDirectPayment(ctx, apiKeyHeader, userID, merchantID, &req)
	if err != nil {
		h.log.Error("Payment processing failed", map[string]interface{}{
			"error":    err.Error(),
			"medium":   req.Medium,
			"amount":   req.Amount,
			"currency": req.Currency,
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	h.log.Info("Payment processed successfully", map[string]interface{}{
		"success":     resp.Success,
		"reference":   resp.Reference,
		"payment_url": resp.PaymentURL,
	})

	c.JSON(http.StatusOK, resp)
}

// Checkout godoc
// @Summary      Create hosted checkout
// @Description  Create a hosted checkout session for payment processing
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request body entity.HostedCheckoutRequest true "Hosted checkout request details"
// @Success      200  {object}  entity.PaymentResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     ApiKeyAuth
// @Router       /payment/checkout [post]
func (h *Handler) Checkout(c *gin.Context) {
	h.log.Info("Starting hosted checkout creation", map[string]interface{}{})

	var req entity.HostedCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Failed to bind JSON request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	if err := req.Validate(); err != nil {
		h.log.Error("Request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	// Validate supported mediums
	for _, medium := range req.SupportedMediums {
		if !slices.Contains(supportedMediumsDeposit, medium) {
			h.log.Error("Unsupported medium in list", map[string]interface{}{
				"medium": medium,
			})
			c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("unsupported medium: %s", medium)))
			return
		}
	}

	h.log.Info("Received hosted checkout request", map[string]interface{}{
		"amount":            req.Amount,
		"currency":          req.Currency,
		"supported_mediums": req.SupportedMediums,
		"reference":         req.Reference,
	})

	apiKeyData, _ := c.Get("apiKey")
	apiKey, _ := apiKeyData.(*apikeyEntity.APIKeyResponse)
	userID := apiKey.UserID
	merchantID := apiKey.MerchantID
	ctx := c.Request.Context()
	apiKeyHeader := c.GetHeader("X-API-Key")

	resp, err := h.paymentUseCase.CreateHostedCheckout(ctx, apiKeyHeader, userID, merchantID, &req)
	if err != nil {
		h.log.Error("Hosted checkout creation failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	h.log.Info("Hosted checkout created successfully", map[string]interface{}{
		"payment_url": resp.PaymentURL,
		"reference":   resp.Reference,
	})

	c.JSON(http.StatusOK, resp)
}

// CheckoutPayment godoc
// @Summary      Process payment from hosted checkout
// @Description  Process payment from hosted checkout page with selected medium and phone number
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request body entity.CheckoutPaymentRequest true "Checkout payment request details"
// @Success      200  {object}  entity.PaymentResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /checkout/makepayment [post]
func (h *Handler) CheckoutPayment(c *gin.Context) {
	h.log.Info("Starting checkout payment processing", map[string]interface{}{})

	var req entity.CheckoutPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Failed to bind JSON request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	if err := req.Validate(); err != nil {
		h.log.Error("Request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	if !slices.Contains(supportedMediumsDeposit, req.Medium) {
		h.log.Error("Unsupported medium", map[string]interface{}{
			"medium": req.Medium,
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("unsupported medium: %s", req.Medium)))
		return
	}

	h.log.Info("Received checkout payment request", map[string]interface{}{
		"hosted_checkout_id": req.HostedCheckoutID,
		"medium":             req.Medium,
		"phone_number":       req.PhoneNumber,
	})

	ctx := c.Request.Context()
	resp, err := h.paymentUseCase.ProcessCheckoutPayment(ctx, &req)
	if err != nil {
		h.log.Error("Checkout payment processing failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	h.log.Info("Checkout payment processed successfully", map[string]interface{}{
		"success":        resp.Success,
		"transaction_id": resp.SocialPayTransactionID,
	})

	c.JSON(http.StatusOK, resp)
}

// GetTransaction godoc
// @Summary      Get transaction details
// @Description  Retrieve details of a specific transaction by ID
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  entity.TransactionResponseDTO
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /payment/receipt/{id} [get]
func (h *Handler) GetReceipt(c *gin.Context) {
	var query entity.TransactionQuery
	query.ID, _ = uuid.Parse(c.Param("id"))

	if err := query.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}
	tx, err := h.paymentUseCase.GetTransactionWithMerchant(c.Request.Context(), query.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	// Map transaction entity to response DTO
	response := entity.TransactionResponseDTO{
		Id:              tx.Id,
		MerchantId:      tx.MerchantId,
		PhoneNumber:     tx.PhoneNumber,
		UserId:          tx.UserId,
		Type:            tx.Type,
		Medium:          tx.Medium,
		Reference:       tx.Reference,
		Comment:         tx.Comment,
		Verified:        tx.Verified,
		Details:         tx.Details,
		CreatedAt:       tx.CreatedAt,
		UpdatedAt:       tx.UpdatedAt,
		ReferenceNumber: tx.ReferenceNumber,
		Test:            tx.Test,
		Status:          tx.Status,
		Description:     tx.Description,
		Token:           tx.Token,
		Amount:          tx.BaseAmount,
		WebhookReceived: tx.WebhookReceived,
		FeeAmount:       tx.FeeAmount,
		AdminNet:        tx.AdminNet,
		VatAmount:       tx.VatAmount,
		MerchantNet:     tx.MerchantNet,
		TotalAmount:     tx.TotalAmount,
		Currency:        tx.Currency,
		CallbackURL:     tx.CallbackURL,
		SuccessURL:      tx.SuccessURL,
		FailedURL:       tx.FailedURL,
		Merchant:        tx.Merchant,
	}

	c.JSON(http.StatusOK, response)
}

// GetTransaction godoc
// @Summary      Get transaction details
// @Description  Retrieve details of a specific transaction by ID
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  entity.TransactionResponseDTO
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /payment/transaction/{id} [get]
func (h *Handler) GetTransaction(c *gin.Context) {
	var query entity.TransactionQuery
	query.ID, _ = uuid.Parse(c.Param("id"))

	if err := query.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}
	h.log.Info("Querying transaction", map[string]interface{}{
		"transaction_id": query.ID,
	})
	tx, err := h.paymentUseCase.GetTransaction(c.Request.Context(), query.ID)
	h.log.Info("Transaction", map[string]interface{}{
		"transaction": tx,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	if tx.Status == txnEntity.PENDING && tx.CreatedAt.Before(time.Now().Add(-1*time.Minute*5)) {
		// Query transaction status from provider
		var id string

		id = tx.ProviderTxId
		if tx.Medium == "AWASH" {
			id = tx.Id.String()
		}
		if tx.Medium == "CBE" {
			id = tx.Id.String()
		}
		queryResp, err := h.paymentUseCase.QueryTransactionStatus(c.Request.Context(), tx.Medium, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, newErrorResponse(err))
			return
		}
		if queryResp != nil {
			if queryResp.Status != tx.Status {
				tx.Status = queryResp.Status
				tx.ProviderTxId = queryResp.ProviderTxId
				tx.ProviderData = queryResp.ProviderData
				h.dispatchWebhookSettlement(c, tx.Id, queryResp)
			}
		}
	}

	// Map transaction entity to response DTO
	response := entity.TransactionResponseDTO{
		Id:              tx.Id,
		MerchantId:      tx.MerchantId,
		PhoneNumber:     tx.PhoneNumber,
		UserId:          tx.UserId,
		Type:            tx.Type,
		Medium:          tx.Medium,
		Reference:       tx.Reference,
		Comment:         tx.Comment,
		Verified:        tx.Verified,
		Details:         tx.Details,
		CreatedAt:       tx.CreatedAt,
		UpdatedAt:       tx.UpdatedAt,
		ReferenceNumber: tx.ReferenceNumber,
		Test:            tx.Test,
		Status:          tx.Status,
		Description:     tx.Description,
		Token:           tx.Token,
		Amount:          tx.BaseAmount,
		WebhookReceived: tx.WebhookReceived,
		FeeAmount:       tx.FeeAmount,
		AdminNet:        tx.AdminNet,
		VatAmount:       tx.VatAmount,
		MerchantNet:     tx.MerchantNet,
		TotalAmount:     tx.TotalAmount,
		Currency:        tx.Currency,
		CallbackURL:     tx.CallbackURL,
		SuccessURL:      tx.SuccessURL,
		FailedURL:       tx.FailedURL,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) dispatchWebhookSettlement(c *gin.Context, transactionID uuid.UUID, transactionStatusQueryResponse *payment.TransactionStatusQueryResponse) {

	providerData, _ := json.Marshal(transactionStatusQueryResponse.ProviderData)

	// Prepare webhook request
	webhookReq := settlementdto.WebhookRequest{
		TransactionID: transactionID.String(),
		Status:        string(transactionStatusQueryResponse.Status),
		Message:       "Transaction status updated by Transaction status query",
		ProviderTxID:  transactionStatusQueryResponse.ProviderTxId,
		ProviderData:  string(providerData),
		Timestamp:     time.Now(),
	}

	h.log.Info("Dispatching webhook", map[string]interface{}{
		"webhook_request": webhookReq,
		"request_id":      c.GetHeader("X-Request-ID"),
	})

	// Dispatch webhook
	if err := h.webhookUseCase.HandleWebhookDispatch(c.Request.Context(), webhookReq); err != nil {
		h.log.Error("Failed to dispatch webhook", map[string]interface{}{
			"error":      err.Error(),
			"request_id": c.GetHeader("X-Request-ID"),
		})
		// Note: We don't return error to client as the payment was processed
		// But we log it for debugging
	}

	h.log.Info("Settlement processing completed", map[string]interface{}{
		"transaction_id": transactionID.String(),
		"status":         transactionStatusQueryResponse.Status,
		"request_id":     c.GetHeader("X-Request-ID"),
	})
}

// RequestWithdrawal godoc
// @Summary      Request a withdrawal
// @Description  Process a withdrawal request for the specified amount
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request body entity.WithdrawalRequest true "Withdrawal request details"
// @Success      200  {object}  entity.PaymentResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     ApiKeyAuth
// @Router       /payment/withdrawal [post]
func (h *Handler) RequestWithdrawal(c *gin.Context) {
	h.log.Info("[Withdrawal] Starting withdrawal request", map[string]interface{}{})

	// Read the raw request body
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Failed to read request body", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		h.log.Info("[Withdrawal] Raw request body", map[string]interface{}{
			"body": string(rawBody),
		})
	}

	// Restore the request body for later binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	var req entity.WithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}
	if !slices.Contains(supportedMediumsWithdrawal, req.Medium) {
		h.log.Error("Unsupported medium", map[string]interface{}{
			"error": "Unsupported medium",
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("Unsupported medium")))
		return
	}

	apiKeyData, _ := c.Get("apiKey")
	apiKey, _ := apiKeyData.(*apikeyEntity.APIKeyResponse)
	userID := apiKey.UserID
	merchantID := apiKey.MerchantID
	apiKeyHeader := c.GetHeader("X-API-Key")

	// The wallet balance check and locking is now handled inside the RequestWithdrawal method
	// in a transaction-safe way using row-level locking
	resp, err := h.paymentUseCase.RequestWithdrawal(c.Request.Context(), apiKeyHeader, userID, merchantID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ErrorResponse represents an error response
// @Description Error response with a message
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message" example:"Invalid request parameters"`
}

// QRMerchantPayment godoc
// @Summary      Process a QR payment for merchant
// @Description  Process a QR payment transaction from checkout page
// @Tags         QR-Payments
// @Accept       json
// @Produce      json
// @Param        request body entity.QRMerchantPaymentRequest true "QR Payment request details"
// @Success      200  {object}  entity.QRPaymentResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /qr/payment/merchant [post]
func (h *Handler) QRMerchantPayment(c *gin.Context) {
	h.log.Info("Starting QR merchant payment processing", map[string]interface{}{})

	// Read the raw request body
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Failed to read request body", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		h.log.Info("Raw QR payment request body", map[string]interface{}{
			"body": string(rawBody),
		})
	}

	// Restore the request body for later binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	var req entity.QRMerchantPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Failed to bind JSON request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	if !slices.Contains(supportedMediumsDeposit, req.Medium) {
		h.log.Error("Unsupported medium", map[string]interface{}{
			"error": "Unsupported medium",
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("Unsupported medium")))
		return
	}

	h.log.Info("Received QR payment request", map[string]interface{}{
		"medium":      req.Medium,
		"amount":      req.Amount,
		"merchant_id": req.MerchantID,
		"currency":    req.Currency,
	})

	if err := req.Validate(); err != nil {
		h.log.Error("Request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	// Get merchant information to get user ID
	ctx := c.Request.Context()
	merchant, err := h.merchantRepo.GetMerchant(ctx, req.MerchantID)
	if err != nil {
		h.log.Error("Failed to get merchant", map[string]interface{}{
			"error":       err.Error(),
			"merchant_id": req.MerchantID,
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("merchant not found")))
		return
	}

	h.log.Info("Found merchant", map[string]interface{}{
		"merchant_id": merchant.ID,
		"user_id":     merchant.UserID,
		"legal_name":  merchant.LegalName,
	})

	// Set defaults
	currency := req.Currency
	if currency == "" {
		currency = "ETB"
	}

	description := req.Description
	if description == "" {
		description = fmt.Sprintf("QR Payment to %s", merchant.LegalName)
	}

	phoneNumber := req.PhoneNumber

	reference := req.Reference
	if reference == "" {
		reference = fmt.Sprintf("QR_%s", uuid.New().String()[:8])
	}

	// Get callback URL from environment
	baseURL := os.Getenv("APP_URL_V2")
	if baseURL == "" {
		baseURL = "http://196.190.251.194:8082" // fallback
	}
	callbackURL := fmt.Sprintf("%s/api/v2/qr/callback", baseURL)

	// Create DirectPaymentRequest
	directPaymentReq := &entity.DirectPaymentRequest{
		Medium:      req.Medium,
		Description: description,
		PhoneNumber: phoneNumber,
		Reference:   reference,
		Amount:      req.Amount,
		Currency:    currency,
		Details:     txnEntity.TransactionDetails{
			// Add any required details
		},
		Redirects: txnEntity.TransactionRedirects{
			Success: fmt.Sprintf("%s/success", baseURL),
			Failed:  fmt.Sprintf("%s/failed", baseURL),
		},
		CallbackURL: callbackURL,
	}

	h.log.Info("QR payment validation successful, processing payment", map[string]interface{}{
		"callback_url": callbackURL,
	})

	// Process payment using merchant ID instead of API key
	resp, err := h.paymentUseCase.ProcessDirectPayment(ctx, req.MerchantID.String(), merchant.UserID, req.MerchantID, directPaymentReq)
	if err != nil {
		h.log.Error("QR payment processing failed", map[string]interface{}{
			"error":       err.Error(),
			"medium":      req.Medium,
			"amount":      req.Amount,
			"currency":    currency,
			"merchant_id": req.MerchantID,
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	h.log.Info("QR payment processed successfully", map[string]interface{}{
		"success":     resp.Success,
		"reference":   resp.Reference,
		"payment_url": resp.PaymentURL,
	})

	// Convert to QR payment response
	qrResp := &entity.QRPaymentResponse{
		Success:                resp.Success,
		Status:                 resp.Status,
		Message:                resp.Message,
		PaymentURL:             resp.PaymentURL,
		Reference:              resp.Reference,
		SocialPayTransactionID: resp.SocialPayTransactionID,
	}

	c.JSON(http.StatusOK, qrResp)
}

// QRCallback godoc
// @Summary      Handle QR payment callback
// @Description  Handle callback from payment processors for QR payments
// @Tags         QR-Payments
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /qr/callback [post]
func (h *Handler) QRCallback(c *gin.Context) {
	h.log.Info("Received QR payment callback", map[string]interface{}{})

	// Read the raw request body
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Failed to read callback request body", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read request body",
		})
		return
	}

	h.log.Info("QR Payment Callback - Raw request body", map[string]interface{}{
		"body":         string(rawBody),
		"content_type": c.GetHeader("Content-Type"),
		"user_agent":   c.GetHeader("User-Agent"),
		"remote_addr":  c.ClientIP(),
	})

	// Print headers
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	h.log.Info("QR Payment Callback - Headers", map[string]interface{}{
		"headers": headers,
	})

	// Try to parse as JSON
	var jsonBody map[string]interface{}
	if err := c.ShouldBindJSON(&jsonBody); err == nil {
		h.log.Info("QR Payment Callback - Parsed JSON", map[string]interface{}{
			"json": jsonBody,
		})
	}

	// Restore the request body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// For now, just acknowledge the callback
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "QR payment callback received and logged",
	})
}

// Helper function to create error responses
func newErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: err.Error(),
	}
}

// GetHostedCheckout godoc
// @Summary      Get hosted checkout details
// @Description  Retrieve hosted checkout details by ID for the checkout page
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Hosted Checkout ID"
// @Success      200  {object}  entity.HostedCheckoutWithMerchantResponseDTO
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /checkout/{id} [get]
func (h *Handler) GetHostedCheckout(c *gin.Context) {
	h.log.Info("Getting hosted checkout details", map[string]interface{}{})

	// Parse hosted checkout ID from URL parameter
	hostedCheckoutIDStr := c.Param("id")
	hostedCheckoutID, err := uuid.Parse(hostedCheckoutIDStr)
	if err != nil {
		h.log.Error("Invalid hosted checkout ID", map[string]interface{}{
			"id":    hostedCheckoutIDStr,
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("invalid hosted checkout ID")))
		return
	}

	ctx := c.Request.Context()
	hostedCheckout, err := h.paymentUseCase.GetHostedCheckoutWithMerchant(ctx, hostedCheckoutID)
	if err != nil {
		h.log.Error("Failed to get hosted checkout", map[string]interface{}{
			"id":    hostedCheckoutID,
			"error": err.Error(),
		})
		c.JSON(http.StatusNotFound, newErrorResponse(fmt.Errorf("hosted checkout not found")))
		return
	}

	h.log.Info("Successfully retrieved hosted checkout", map[string]interface{}{
		"id":     hostedCheckout.ID,
		"status": hostedCheckout.Status,
	})

	c.JSON(http.StatusOK, hostedCheckout)
}

// GetQRLinkForPayment godoc
// @Summary      Get QR link for payment
// @Description  Get QR link details for payment processing (public endpoint)
// @Tags         QR-Payments
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "QR Link ID"
// @Success      200  {object}  qrEntity.QRLinkResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /qr/payment/link/{id} [get]
func (h *Handler) GetQRLinkForPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("invalid QR link ID")))
		return
	}

	response, err := h.qrUseCase.GetQRLink(c.Request.Context(), id)
	if err != nil {
		h.log.Error("Failed to get QR link for payment", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusNotFound, newErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ProcessQRPayment godoc
// @Summary      Process QR payment
// @Description  Process a payment using QR link
// @Tags         QR-Payments
// @Accept       json
// @Produce      json
// @Param        id      path      string                    true  "QR Link ID"
// @Param        request body      qrEntity.QRPaymentRequest   true  "QR payment request"
// @Success      200  {object}  qrEntity.QRPaymentResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /qr/payment/link/{id} [post]
func (h *Handler) ProcessQRPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("invalid QR link ID")))
		return
	}

	var req qrEntity.QRPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Failed to bind JSON request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	response, err := h.qrUseCase.ProcessQRPayment(c.Request.Context(), id, &req)
	if err != nil {
		h.log.Error("Failed to process QR payment", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCheckout godoc
// @Summary      Update hosted checkout
// @Description  Update hosted checkout details (only allowed when status is PENDING and not expired)
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id      path      string                    true  "Hosted Checkout ID"
// @Param        request body      entity.UpdateHostedCheckoutRequest true "Update request details"
// @Success      200  {object}  entity.PaymentResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     ApiKeyAuth
// @Router       /payment/checkout/{id} [patch]
func (h *Handler) UpdateCheckout(c *gin.Context) {
	h.log.Info("Starting hosted checkout update", map[string]interface{}{})

	// Parse hosted checkout ID from URL parameter
	hostedCheckoutIDStr := c.Param("id")
	hostedCheckoutID, err := uuid.Parse(hostedCheckoutIDStr)
	if err != nil {
		h.log.Error("Invalid hosted checkout ID", map[string]interface{}{
			"id":    hostedCheckoutIDStr,
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("invalid hosted checkout ID")))
		return
	}

	var req entity.UpdateHostedCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Failed to bind JSON request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	if err := req.Validate(); err != nil {
		h.log.Error("Request validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	// Validate supported mediums if provided
	if len(req.SupportedMediums) > 0 {
		for _, medium := range req.SupportedMediums {
			if !slices.Contains(supportedMediumsDeposit, medium) {
				h.log.Error("Unsupported medium in update list", map[string]interface{}{
					"medium": medium,
				})
				c.JSON(http.StatusBadRequest, newErrorResponse(fmt.Errorf("unsupported medium: %s", medium)))
				return
			}
		}
	}

	h.log.Info("Received hosted checkout update request", map[string]interface{}{
		"hosted_checkout_id": hostedCheckoutID,
		"amount":             req.Amount,
		"currency":           req.Currency,
		"supported_mediums":  req.SupportedMediums,
	})

	apiKeyData, _ := c.Get("apiKey")
	apiKey, _ := apiKeyData.(*apikeyEntity.APIKeyResponse)
	userID := apiKey.UserID
	merchantID := apiKey.MerchantID
	ctx := c.Request.Context()
	apiKeyHeader := c.GetHeader("X-API-Key")

	resp, err := h.paymentUseCase.UpdateHostedCheckout(ctx, apiKeyHeader, userID, merchantID, hostedCheckoutID, &req)
	if err != nil {
		h.log.Error("Hosted checkout update failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	h.log.Info("Hosted checkout updated successfully", map[string]interface{}{
		"payment_url": resp.PaymentURL,
		"reference":   resp.Reference,
	})

	c.JSON(http.StatusOK, resp)
}
