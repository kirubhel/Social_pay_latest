package gin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authv2Entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/shared/errorxx"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	middleware "github.com/socialpay/socialpay/src/pkg/shared/middleware"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	"github.com/socialpay/socialpay/src/pkg/shared/response"
	"github.com/socialpay/socialpay/src/pkg/transaction/adapter/dto"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/exporter"
	"github.com/socialpay/socialpay/src/pkg/transaction/usecase"
	settlementdto "github.com/socialpay/socialpay/src/pkg/webhook/adapter/dto"
)

// WebhookDispatcher interface for webhook dispatch to avoid import cycles
type WebhookDispatcher interface {
	HandleWebhookDispatch(ctx context.Context, req settlementdto.WebhookRequest) error
}

// Handler handles HTTP requests for transactions
type Handler struct {
	useCase            usecase.TransactionUseCase
	webhookDispatcher  WebhookDispatcher
	log                logging.Logger
	middlewareProvider *middleware.MiddlewareProvider
}

func NewTransactionHistoryHandler(
	useCase usecase.TransactionUseCase,
	middlewareProvider *middleware.MiddlewareProvider,
	webhookDispatcher WebhookDispatcher,
) Handler {

	return Handler{
		log:                logging.NewStdLogger("[TRANSACTION] [HANDLER]"),
		useCase:            useCase,
		middlewareProvider: middlewareProvider,
		webhookDispatcher:  webhookDispatcher,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// defining routes
	g := rg.Group("/transactions", ginn.ErrorMiddleWare(), h.middlewareProvider.JWTAuth, h.middlewareProvider.RBAC.RequirePermissionForMerchant(authv2Entity.RESOURCE_TRANSACTION, authv2Entity.OPERATION_READ))
	g.GET("/history", h.GetTransactions)
	g.POST("/history", func(c *gin.Context) { h.GetTransactionsByParameter(c, false) })
	g.POST("/analytics", h.GetTransactionAnalytics)
	g.POST("/chart", h.GetChartData)
	g.POST("/data/export", func(c *gin.Context) { h.GetTransactionData(c, false) })
}

func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/transactions/admin", ginn.ErrorMiddleWare(), h.middlewareProvider.JWTAuth, h.middlewareProvider.RBAC.RequirePermissionForAdmin(authv2Entity.RESOURCE_TRANSACTION, authv2Entity.OPERATION_ADMIN_READ))
	g.POST("/history", func(c *gin.Context) { h.GetTransactionsByParameter(c, true) })
	g.POST("/data/export", func(ctx *gin.Context) { h.GetTransactionData(ctx, true) })
	g.POST("/override/:transactionID", h.middlewareProvider.RBAC.RequirePermissionForAdmin(authv2Entity.RESOURCE_TRANSACTION, authv2Entity.OPERATION_ADMIN_OVERRIDE), h.OverrideTransactionStatus)

	// Admin analytics endpoints - using RequireAdmin for now
	adminAnalytics := rg.Group("/admin/analytics", ginn.ErrorMiddleWare(), h.middlewareProvider.JWTAuth, h.middlewareProvider.RBAC.RequirePermissionForAdmin(authv2Entity.RESOURCE_TRANSACTION, authv2Entity.OPERATION_ADMIN_READ))
	adminAnalytics.GET("/transactions", h.GetAdminTransactionAnalytics)
	adminAnalytics.GET("/chart", h.GetAdminChartData)
	adminAnalytics.GET("/merchant-growth", h.GetMerchantGrowthAnalytics)
}

// get transactions
// @Tags transactions
// @Summary get transactions
// @Description get transacton history
// @Accept json
// @Produce json
// @Param page query int true "page number"
// @Param page_size query int  true "page size"
// @Success 200 {object} response.PaginatedResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router /transactions/history [get]
// @Security BearerAuth
func (h *Handler) GetTransactions(c *gin.Context) {

	// request context
	requestContext := c.Request.Context()

	pagination, err := pagination.NewPagination(c, h.log)

	if err != nil {

		err = errorxx.ErrAppBadInput.Wrap(err, "pagination binding error").
			WithProperty(errorxx.ErrorCode, 400)

		c.Error(err)
		return

	}

	userId, _ := ginn.GetUserIDFromContext(c)

	data, count, err := h.useCase.GetTransactions(requestContext, userId, *pagination)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination.GetInfo(count),
	})

}

// get transactions data
// @Tags transactions
// @Summary get transaction by specific filter parameter
// @Description get transacton history with given specific parameter. Use date format YYYY-MM-DD (e.g. 2023-01-01)
// @Accept json
// @Produce json
// @Param filterParameter body entity.FilterParameters true "request body with date format YYYY-MM-DD"
// @Param page query int true "page number"
// @Param page_size query int  true "page size"
// @Success 200 {object} response.PaginatedResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router /transactions/history [post]
// @Security BearerAuth
func (h *Handler) GetTransactionsByParameter(c *gin.Context, queryForAllUsers bool) {

	var filterParameters entity.FilterParameters
	requestContext := c.Request.Context()

	pagination, err := pagination.NewPagination(c, h.log)

	if err != nil {

		err = errorxx.ErrAppBadInput.Wrap(err, "pagination binding error").
			WithProperty(errorxx.ErrorCode, 400)

		c.Error(err)
		return

	}

	if err := c.ShouldBind(&filterParameters); err != nil {

		err = errorxx.ErrAppBadInput.Wrap(err, "binding filter parameters").
			WithProperty(errorxx.ErrorCode, 400)

		h.log.Info("error while binding filter parameter",
			map[string]interface{}{
				"type":    "binding",
				"error":   err.Error(),
				"context": requestContext,
			})

		// setting err on context
		c.Error(err)
		return
	}

	// get user_id from the context
	userId, _ := ginn.GetUserIDFromContext(c)

	// usecase layer
	data, count, err := h.useCase.GetTransactionByParamenters(requestContext, userId,
		&filterParameters, *pagination, queryForAllUsers)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination.GetInfo(count),
	})
}

// get transactions data
// @Tags transactions
// @Summary exporting transaction data
// @Description export the transaction data with pdf format. Use date format YYYY-MM-DD (e.g. 2023-01-01)
// @Accept json
// @Produce application/pdf
// @Param filterParameter body entity.FilterParameters true "request body with date format YYYY-MM-DD"
// @Param page query int true "page number"
// @Param page_size query int true "page size"
// @Success 200 {string} binary "PDF file"
// @Router /transactions/data/export [post]
// @Security BearerAuth
func (h *Handler) GetTransactionData(c *gin.Context, queryForAllUsers bool) {
	var filterParameters entity.FilterParameters

	requestContext := c.Request.Context()

	pagination, err := pagination.NewPagination(c, h.log)

	if err != nil {

		err = errorxx.ErrAppBadInput.Wrap(err, "pagination binding error").
			WithProperty(errorxx.ErrorCode, 400)

		c.Error(err)
		return

	}

	if err := c.ShouldBind(&filterParameters); err != nil {

		err = errorxx.ErrAppBadInput.Wrap(err, "filter parameter binding error").
			WithProperty(errorxx.ErrorCode, 400)

		h.log.Error("error while binding filter parameter",
			map[string]interface{}{
				"type":    "binding",
				"error":   err.Error(),
				"context": requestContext,
			})
		// setting err on context
		c.Error(err)
		return
	}
	// get user_id from the context
	userId, _ := ginn.GetUserIDFromContext(c)

	data, _, err := h.useCase.GetTransactionByParamenters(requestContext, userId,
		&filterParameters,
		*pagination, queryForAllUsers)

	if err != nil {
		c.Error(err)

		return
	}

	pdf := exporter.CreatePDFReport("Transaction Report", data)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=transactions.pdf")
	_ = pdf.Output(c.Writer)

}

// OverrideTransactionStatus overrides the status of a transaction
// @tags transactions
// @Summary Override transaction status
// @Description Override the status of a transaction with admin approval
// @Accept multipart/form-data
// @Produce json
// @Param transactionID path string true "Transaction ID" format(uuid)
// @Param status formData string true "New status (INITIATED, PENDING, SUCCESS, FAILED, REFUNDED, EXPIRED, CANCELED)"
// @Param reason formData string true "Reason for override (minimum 10 characters)"
// @Param adminID formData string true "Admin ID" format(uuid)
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transactions/override/{transactionID} [post]
// @security BearerAuth
func (h *Handler) OverrideTransactionStatus(c *gin.Context) {
	adminID, exists := ginn.GetUserIDFromContext(c)
	if !exists {
		err := errorxx.ErrAuthUnauthorized.Wrap(nil, "user ID not found in context").
			WithProperty(errorxx.ErrorCode, 401)
		c.Error(err)
		return
	}
	// Create override request DTO
	overrideReq := dto.OverrideTransactionStatusRequest{
		TransactionID: c.Param("transactionID"),
		Status:        entity.TransactionStatus(c.PostForm("status")),
		Reason:        c.PostForm("reason"),
		AdminID:       adminID.String(),
	}

	// Validate the request
	if err := overrideReq.Validate(); err != nil {
		h.log.Error("Validation failed for override request", map[string]interface{}{
			"error":          err.Error(),
			"transaction_id": overrideReq.TransactionID,
			"admin_id":       overrideReq.AdminID,
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedTxnID, err := uuid.Parse(overrideReq.TransactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	// Dispatch webhook instead of direct database update
	if err := h.dispatchWebhookForOverride(c.Request.Context(), parsedTxnID, &overrideReq); err != nil {
		h.log.Error("Failed to dispatch webhook for transaction override", map[string]interface{}{
			"error":          err.Error(),
			"transaction_id": overrideReq.TransactionID,
			"admin_id":       overrideReq.AdminID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process override: %v", err)})
		return
	}

	h.log.Info("Transaction status override dispatched successfully", map[string]interface{}{
		"transaction_id": overrideReq.TransactionID,
		"new_status":     overrideReq.Status,
		"admin_id":       overrideReq.AdminID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status override dispatched successfully"})
}

// dispatchWebhookForOverride creates and dispatches a webhook for manual transaction status override
func (h *Handler) dispatchWebhookForOverride(ctx context.Context, transactionID uuid.UUID, overrideReq *dto.OverrideTransactionStatusRequest) error {
	// Create a TransactionStatusQueryResponse for manual override
	transactionStatusQueryResponse := &payment.TransactionStatusQueryResponse{
		Status:       overrideReq.Status,
		ProviderTxId: overrideReq.GetProviderTxID(),
		ProviderData: map[string]interface{}{
			"override_type": "manual",
			"admin_id":      overrideReq.AdminID,
			"reason":        overrideReq.Reason,
		},
	}

	// Marshal provider data
	providerData, err := json.Marshal(transactionStatusQueryResponse.ProviderData)
	if err != nil {
		return fmt.Errorf("failed to marshal provider data: %w", err)
	}

	// Prepare webhook request with manual override information
	webhookReq := settlementdto.WebhookRequest{
		TransactionID: transactionID.String(),
		Status:        string(transactionStatusQueryResponse.Status),
		Message:       overrideReq.GetMessage(),
		ProviderTxID:  transactionStatusQueryResponse.ProviderTxId,
		ProviderData:  string(providerData),
		Timestamp:     time.Now(),
	}

	h.log.Info("Dispatching webhook for manual override", map[string]interface{}{
		"webhook_request": webhookReq,
		"transaction_id":  transactionID.String(),
		"admin_id":        overrideReq.AdminID,
		"reason":          overrideReq.Reason,
	})

	// Dispatch webhook
	if err := h.webhookDispatcher.HandleWebhookDispatch(ctx, webhookReq); err != nil {
		h.log.Error("Failed to dispatch webhook for manual override", map[string]interface{}{
			"error":          err.Error(),
			"transaction_id": transactionID.String(),
			"admin_id":       overrideReq.AdminID,
		})
		return fmt.Errorf("failed to dispatch webhook: %w", err)
	}

	h.log.Info("Successfully dispatched webhook for manual override", map[string]interface{}{
		"transaction_id": transactionID.String(),
		"status":         transactionStatusQueryResponse.Status,
		"admin_id":       overrideReq.AdminID,
	})

	return nil
}

// GetTransactionAnalytics handles transaction analytics requests
// @Tags transactions
// @Summary Get transaction analytics
// @Description Get aggregated transaction analytics with filters and period comparison
// @Accept json
// @Produce json
// @Param analyticsFilter body entity.AnalyticsFilter true "Analytics filter parameters"
// @Success 200 {object} response.SuccessResponse{data=entity.TransactionAnalytics}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transactions/analytics [post]
// @Security BearerAuth
func (h *Handler) GetTransactionAnalytics(c *gin.Context) {
	var filter entity.AnalyticsFilter
	requestContext := c.Request.Context()

	// Bind the analytics filter from request body
	if err := c.ShouldBind(&filter); err != nil {
		err = errorxx.ErrAppBadInput.Wrap(err, "binding analytics filter").
			WithProperty(errorxx.ErrorCode, 400)

		h.log.Error("error while binding analytics filter", map[string]interface{}{
			"type":    "binding",
			"error":   err.Error(),
			"context": requestContext,
		})

		c.Error(err)
		return
	}

	// Get user ID from context and add it to request context
	merchantID, exists := ginn.GetMerchantIDFromContext(c)
	if !exists {
		err := errorxx.ErrAuthUnauthorized.Wrap(nil, "user ID not found in context").
			WithProperty(errorxx.ErrorCode, 401)

		h.log.Error("user ID not found in context", map[string]interface{}{
			"context": requestContext,
		})

		c.Error(err)
		return
	}

	// Add user ID to context for usecase
	ctx := context.WithValue(requestContext, "merchant_id", merchantID)

	// Call usecase layer
	analytics, err := h.useCase.GetTransactionAnalytics(ctx, &filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Data:    analytics,
	})
}

// GetChartData handles chart data requests
// @Tags transactions
// @Summary Get transaction chart data
// @Description Get chart data for transaction analytics with date aggregation
// @Accept json
// @Produce json
// @Param chartFilter body entity.ChartFilter true "Chart filter parameters"
// @Success 200 {object} response.SuccessResponse{data=entity.ChartData}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transactions/chart [post]
// @Security BearerAuth
func (h *Handler) GetChartData(c *gin.Context) {
	var filter entity.ChartFilter
	requestContext := c.Request.Context()

	// Bind the chart filter from request body
	if err := c.ShouldBind(&filter); err != nil {
		err = errorxx.ErrAppBadInput.Wrap(err, "binding chart filter").
			WithProperty(errorxx.ErrorCode, 400)

		h.log.Error("error while binding chart filter", map[string]interface{}{
			"type":    "binding",
			"error":   err.Error(),
			"context": requestContext,
		})

		c.Error(err)
		return
	}

	// Get user ID from context and add it to request context
	merchantID, exists := ginn.GetMerchantIDFromContext(c)
	if !exists {
		err := errorxx.ErrAuthUnauthorized.Wrap(nil, "user ID not found in context").
			WithProperty(errorxx.ErrorCode, 401)

		h.log.Error("user ID not found in context", map[string]interface{}{
			"context": requestContext,
		})

		c.Error(err)
		return
	}

	// Add user ID to context for usecase
	ctx := context.WithValue(requestContext, "merchant_id", merchantID)

	// Call usecase layer
	chartData, err := h.useCase.GetChartData(ctx, &filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Data:    chartData,
	})
}

// Admin Analytics Endpoints

// GetAdminTransactionAnalytics godoc
// @Summary Get admin transaction analytics
// @Description Get comprehensive transaction analytics for admin with VAT, fees, and admin net information
// @Tags Admin Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param merchant_id query string false "Filter by merchant ID"
// @Param transaction_type query string false "Filter by transaction type"
// @Param status query string false "Filter by status"
// @Success 200 {object} entity.AdminTransactionAnalytics
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/analytics/transactions [get]
func (h *Handler) GetAdminTransactionAnalytics(c *gin.Context) {
	// Parse query parameters
	filter, err := h.parseAnalyticsFilter(c)
	if err != nil {
		err = errorxx.ErrAppBadInput.Wrap(err, "invalid analytics filter parameters").
			WithProperty(errorxx.ErrorCode, 400)
		c.Error(err)
		return
	}

	// Get analytics data
	analytics, err := h.useCase.GetAdminTransactionAnalytics(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Data:    analytics,
	})
}

// GetAdminChartData godoc
// @Summary Get admin chart data
// @Description Get chart data for admin dashboard with comprehensive financial information
// @Tags Admin Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param merchant_id query string false "Filter by merchant ID"
// @Param chart_type query string false "Chart type (daily, weekly, monthly)"
// @Param date_unit query string false "Date unit (day, week, month)"
// @Success 200 {object} entity.ChartData
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/analytics/chart [get]
func (h *Handler) GetAdminChartData(c *gin.Context) {
	// Parse query parameters
	filter, err := h.parseChartFilter(c)
	if err != nil {
		err = errorxx.ErrAppBadInput.Wrap(err, "invalid chart filter parameters").
			WithProperty(errorxx.ErrorCode, 400)
		c.Error(err)
		return
	}

	// Get chart data
	chartData, err := h.useCase.GetAdminChartData(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Data:    chartData,
	})
}

// GetMerchantGrowthAnalytics godoc
// @Summary Get merchant growth analytics
// @Description Get merchant growth statistics and trends
// @Tags Admin Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param date_unit query string false "Date unit (day, week, month)" default(month)
// @Success 200 {object} entity.MerchantGrowthAnalytics
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/analytics/merchant-growth [get]
func (h *Handler) GetMerchantGrowthAnalytics(c *gin.Context) {
	// Parse date parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	dateUnitStr := c.DefaultQuery("date_unit", "month")

	// Parse dates
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			err = errorxx.ErrAppBadInput.Wrap(err, "invalid start_date format").
				WithProperty(errorxx.ErrorCode, 400)
			c.Error(err)
			return
		}
	} else {
		// Default to 12 months ago
		startDate = time.Now().AddDate(-1, 0, 0)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			err = errorxx.ErrAppBadInput.Wrap(err, "invalid end_date format").
				WithProperty(errorxx.ErrorCode, 400)
			c.Error(err)
			return
		}
	} else {
		// Default to now
		endDate = time.Now()
	}

	// Parse date unit
	var dateUnit entity.DateUnit
	switch dateUnitStr {
	case "day":
		dateUnit = entity.DAY
	case "week":
		dateUnit = entity.WEEK
	case "month":
		dateUnit = entity.MONTH
	default:
		err = errorxx.ErrAppBadInput.Wrap(nil, "invalid date_unit").
			WithProperty(errorxx.ErrorCode, 400)
		c.Error(err)
		return
	}

	// Get merchant growth analytics
	analytics, err := h.useCase.GetMerchantGrowthAnalytics(c.Request.Context(), startDate, endDate, dateUnit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Data:    analytics,
	})
}

// Helper methods for parsing filters

func (h *Handler) parseAnalyticsFilter(c *gin.Context) (*entity.AnalyticsFilter, error) {
	filter := &entity.AnalyticsFilter{}

	// Parse dates from query parameters
	if startDate := c.Query("start_date"); startDate != "" {
		parsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, errorxx.ErrAppBadInput.Wrap(err, "invalid start_date format")
		}
		filter.StartDate = parsed
	}

	if endDate := c.Query("end_date"); endDate != "" {
		parsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, errorxx.ErrAppBadInput.Wrap(err, "invalid end_date format")
		}
		filter.EndDate = parsed
	}

	// Parse other filters
	if merchantID := c.Query("merchant_id"); merchantID != "" {
		filter.MerchantID = []string{merchantID}
	}

	// Parse status filters (can be multiple)
	statusValues := c.QueryArray("status")
	if len(statusValues) > 0 {
		var statuses []entity.TransactionStatus
		for _, status := range statusValues {
			statuses = append(statuses, entity.TransactionStatus(status))
		}
		filter.Status = statuses
	}

	// Parse type filters (can be multiple)
	typeValues := c.QueryArray("type")
	if len(typeValues) > 0 {
		var types []entity.TransactionType
		for _, typ := range typeValues {
			types = append(types, entity.TransactionType(typ))
		}
		filter.Type = types
	}

	// Parse medium filters (can be multiple)
	mediumValues := c.QueryArray("medium")
	if len(mediumValues) > 0 {
		var mediums []entity.TransactionMedium
		for _, medium := range mediumValues {
			mediums = append(mediums, entity.TransactionMedium(medium))
		}
		filter.Medium = mediums
	}

	// Parse source filters (can be multiple)
	sourceValues := c.QueryArray("source")
	if len(sourceValues) > 0 {
		var sources []entity.TransactionSource
		for _, source := range sourceValues {
			sources = append(sources, entity.TransactionSource(source))
		}
		filter.Source = sources
	}

	// Parse QR tag filters (can be multiple)
	qrTagValues := c.QueryArray("qr_tag")
	if len(qrTagValues) > 0 {
		filter.QRTag = qrTagValues
	}

	// Parse amount range filters
	if amountMinStr := c.Query("amount_min"); amountMinStr != "" {
		if amountMin, err := strconv.ParseFloat(amountMinStr, 64); err == nil {
			filter.AmountMin = &amountMin
		}
	}

	if amountMaxStr := c.Query("amount_max"); amountMaxStr != "" {
		if amountMax, err := strconv.ParseFloat(amountMaxStr, 64); err == nil {
			filter.AmountMax = &amountMax
		}
	}

	return filter, nil
}

func (h *Handler) parseChartFilter(c *gin.Context) (*entity.ChartFilter, error) {
	filter := &entity.ChartFilter{}

	// Parse dates
	if startDate := c.Query("start_date"); startDate != "" {
		parsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, errorxx.ErrAppBadInput.Wrap(err, "invalid start_date format")
		}
		filter.StartDate = parsed
	}

	if endDate := c.Query("end_date"); endDate != "" {
		parsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, errorxx.ErrAppBadInput.Wrap(err, "invalid end_date format")
		}
		filter.EndDate = parsed
	}

	// Parse other filters
	if merchantID := c.Query("merchant_id"); merchantID != "" {
		filter.MerchantID = []string{merchantID}
	}

	// Parse status filters (can be multiple)
	statusValues := c.QueryArray("status")
	if len(statusValues) > 0 {
		var statuses []entity.TransactionStatus
		for _, status := range statusValues {
			statuses = append(statuses, entity.TransactionStatus(status))
		}
		filter.Status = statuses
	}

	// Parse type filters (can be multiple)
	typeValues := c.QueryArray("type")
	if len(typeValues) > 0 {
		var types []entity.TransactionType
		for _, typ := range typeValues {
			types = append(types, entity.TransactionType(typ))
		}
		filter.Type = types
	}

	// Parse medium filters (can be multiple)
	mediumValues := c.QueryArray("medium")
	if len(mediumValues) > 0 {
		var mediums []entity.TransactionMedium
		for _, medium := range mediumValues {
			mediums = append(mediums, entity.TransactionMedium(medium))
		}
		filter.Medium = mediums
	}

	// Parse source filters (can be multiple)
	sourceValues := c.QueryArray("source")
	if len(sourceValues) > 0 {
		var sources []entity.TransactionSource
		for _, source := range sourceValues {
			sources = append(sources, entity.TransactionSource(source))
		}
		filter.Source = sources
	}

	// Parse QR tag filters (can be multiple)
	qrTagValues := c.QueryArray("qr_tag")
	if len(qrTagValues) > 0 {
		filter.QRTag = qrTagValues
	}

	// Parse amount range filters
	if amountMinStr := c.Query("amount_min"); amountMinStr != "" {
		if amountMin, err := strconv.ParseFloat(amountMinStr, 64); err == nil {
			filter.AmountMin = &amountMin
		}
	}

	if amountMaxStr := c.Query("amount_max"); amountMaxStr != "" {
		if amountMax, err := strconv.ParseFloat(amountMaxStr, 64); err == nil {
			filter.AmountMax = &amountMax
		}
	}

	if chartType := c.Query("chart_type"); chartType != "" {
		filter.ChartType = chartType
	}

	if dateUnit := c.Query("date_unit"); dateUnit != "" {
		switch dateUnit {
		case "day":
			filter.DateUnit = entity.DAY
		case "week":
			filter.DateUnit = entity.WEEK
		case "month":
			filter.DateUnit = entity.MONTH
		default:
			return nil, errorxx.ErrAppBadInput.New("invalid date_unit")
		}
	}

	return filter, nil
}
