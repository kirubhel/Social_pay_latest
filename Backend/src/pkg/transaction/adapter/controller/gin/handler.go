package gin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/socialpay/socialpay/src/pkg/shared/errorxx"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/pagination"
	"github.com/socialpay/socialpay/src/pkg/shared/response"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	"github.com/socialpay/socialpay/src/pkg/transaction/core/exporter"
	"github.com/socialpay/socialpay/src/pkg/transaction/usecase"
)

// Handler handles HTTP requests for transactions
type Handler struct {
	useCase    usecase.TransactionUseCase
	log        logging.Logger
	middleware *gin.HandlerFunc
}

func NewTransactionHistoryHandler(
	useCase usecase.TransactionUseCase,
	middleware gin.HandlerFunc,
) Handler {

	return Handler{
		log:        logging.NewStdLogger("[TRANSACTION] [HANDLER]"),
		useCase:    useCase,
		middleware: &middleware,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// defining routes
	g := rg.Group("/transactions", ginn.ErrorMiddleWare(), *h.middleware, ginn.MerchantIDMiddleware())
	g.GET("/history", h.GetTransactions)
	g.POST("/history", h.GetTransactionsByParameter)
	g.POST("/analytics", h.GetTransactionAnalytics)
	g.POST("/chart", h.GetChartData)
	g.POST("/data/export", h.GetTransactionData)
	g.POST("/override/:transactionID", h.OverrideTransactionStatus)
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
func (h *Handler) GetTransactionsByParameter(c *gin.Context) {

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
		&filterParameters, *pagination)

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
func (h *Handler) GetTransactionData(c *gin.Context) {
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
		*pagination)

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
	txnID := c.Param("transactionID")
	newStatus := c.PostForm("status")
	reason := c.PostForm("reason")
	adminID := c.PostForm("adminID")

	if newStatus == "" || reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status and reason are required"})
		return
	}

	if len(reason) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason must be at least 10 characters long"})
		return
	}

	parsedTxnID, err := uuid.Parse(txnID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	// Convert string status to entity.TransactionStatus
	status := entity.TransactionStatus(newStatus)

	if err := h.useCase.OverrideTransactionStatus(c.Request.Context(), parsedTxnID, status, reason, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status overridden successfully"})
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
