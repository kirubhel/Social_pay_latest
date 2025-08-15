package gin

import (
	"net/http"
	"os"

	auth_service "github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	"github.com/socialpay/socialpay/src/pkg/erp_v2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/erp_v2/usecase"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type IDResponse struct {
	ID uuid.UUID `json:"id"`
}

// ApiError represents an error response structure.
type ApiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ErrorResponse represents a standardized error response.
type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   ApiError `json:"error"`
}

// ERPHandler manages ERP HTTP requests
type ERPHandler struct {
	authService auth_service.AuthService
	log         logging.Logger
	useCase     usecase.ERPUseCase
	rbac        *ginn.RBACV2
}

// SuccessResponse represents a successful API response
// @Schema
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty" swaggertype:"object"`
}

// NewERPHandler creates a new ERP handler
func NewERPHandler(
	authService auth_service.AuthService,
	useCase usecase.ERPUseCase,
	rbac *ginn.RBACV2,
) *ERPHandler {
	return &ERPHandler{
		authService: authService,
		log:         logging.NewStdLogger("[ERP_V2] [HANDLER]"),
		useCase:     useCase,
		rbac:        rbac,
	}
}

// RegisterRoutes registers ERP management routes under the provided router group
func (h *ERPHandler) RegisterRoutes(r *gin.RouterGroup) {
	jwtConfig := ginn.JWTAuthMiddlewareConfig{
		AuthService: h.authService,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Public:      false,
	}

	// ERP routes group - will be mounted under /api/v2/erp
	erpRoutes := r.Group("/erp")
	erpRoutes.Use(ginn.JWTAuthMiddleware(jwtConfig), h.rbac.RequireMerchantOwner())

	// Catalog routes
	catalogs := erpRoutes.Group("/catalogs")
	{
		catalogs.POST("", h.CreateCatalog)
		catalogs.GET("/:id", h.GetCatalog)
		catalogs.GET("", h.GetCatalogs)
		catalogs.PUT("/:id", h.UpdateCatalog)
		catalogs.DELETE("/:id", h.DeleteCatalog)
	}

	// Customer routes
	customers := erpRoutes.Group("/customers")
	{
		customers.POST("", h.CreateCustomer)
		customers.GET("/:id", h.GetCustomer)
		customers.GET("", h.GetCustomers)
		customers.PUT("/:id", h.UpdateCustomer)
		customers.DELETE("/:id", h.DeleteCustomer)
	}

	// Product routes
	products := erpRoutes.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.GET("/:id", h.GetProduct)
		products.GET("", h.GetProducts)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}

	// Warehouse routes
	warehouses := erpRoutes.Group("/warehouses")
	{
		warehouses.POST("", h.CreateWarehouse)
		warehouses.GET("/:id", h.GetWarehouse)
		warehouses.GET("", h.GetWarehouses)
		warehouses.PUT("/:id", h.UpdateWarehouse)
		warehouses.DELETE("/:id", h.DeleteWarehouse)
	}

	// Payment Method routes
	paymentMethods := erpRoutes.Group("/payment-methods")
	{
		paymentMethods.POST("", h.CreatePaymentMethod)
		paymentMethods.GET("/:id", h.GetPaymentMethod)
		paymentMethods.GET("", h.GetPaymentMethods)
		paymentMethods.PUT("/:id", h.UpdatePaymentMethod)
		paymentMethods.DELETE("/:id", h.DeletePaymentMethod)
	}

	// Order routes
	orders := erpRoutes.Group("/orders")
	{
		orders.POST("", h.CreateOrder)
		orders.GET("/:id", h.GetOrder)
		orders.GET("", h.GetOrders)
		orders.PUT("/:id", h.UpdateOrder)
		orders.DELETE("/:id", h.DeleteOrder)
	}
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Create a new customer in the ERP system
// @Tags ERP
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.CreateCustomerRequest true "Customer creation request"
// @Success 201 {object} entity.CustomerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /erp/customers [post]
func (h *ERPHandler) CreateCustomer(c *gin.Context) {
	var req entity.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	userID := user.ID

	customer, err := h.useCase.CreateCustomer(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// GetCustomer godoc
// @Summary Get a customer by ID
// @Description Get a customer from the ERP system by their ID
// @Tags ERP
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID" format(uuid)
// @Success 200 {object} entity.CustomerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /erp/customers/{id} [get]
func (h *ERPHandler) GetCustomer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	customer, err := h.useCase.GetCustomer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// GetCustomers godoc
// @Summary Get a list of customers
// @Description Get a list of customers from the ERP system with pagination
// @Tags ERP
// @Produce json
// @Security BearerAuth
// @Param merchantID query string true "Merchant ID" format(uuid)
// @Param skip query int false "Number of items to skip" default(0)
// @Param take query int false "Number of items to take" default(10)
// @Success 200 {object} entity.CustomersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/customers [get]
func (h *ERPHandler) GetCustomers(c *gin.Context) {
	// Implement query parameter parsing for pagination, filtering, etc.
	var params entity.GetCustomersParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customers, err := h.useCase.GetCustomers(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// UpdateCustomer godoc
// @Summary Update an existing customer
// @Description Update an existing customer in the ERP system by their ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param request body entity.UpdateCustomerRequest true "Customer update request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/customers/{id} [put]
func (h *ERPHandler) UpdateCustomer(c *gin.Context) {
	// Parse customer ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid customer ID format",
			},
		})
		return
	}

	// Parse request body
	var req entity.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	// Call use case
	if err := h.useCase.UpdateCustomer(c.Request.Context(), id, &req, user.ID); err != nil {
		// Handle specific error types
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "NOT_FOUND",
					Message: "Customer not found",
				},
			})
			return
		}

		h.log.Error("Failed to update customer", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to update customer",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

// DeleteCustomer godoc
// @Summary Delete a customer by ID
// @Description Delete a customer from the ERP system by their ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/customers/{id} [delete]
func (h *ERPHandler) DeleteCustomer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.useCase.DeleteCustomer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

// CreateCatalog godoc
// @Summary Create a new catalog
// @Description Create a new catalog in the ERP system
// @Tags Catalogs
// @Accept json
// @Produce json
// @Param request body entity.CreateCatalogRequest true "Catalog creation request"
// @Success 201 {object} SuccessResponse{data=entity.Catalog}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/catalogs [post]
func (h *ERPHandler) CreateCatalog(c *gin.Context) {
	// Parse request body
	var req entity.CreateCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// Get user from context
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	// Create catalog
	_, err := h.useCase.CreateCatalog(c.Request.Context(), &req, *user.MerchantID, user.ID)
	if err != nil {
		h.log.Error("Failed to create catalog", map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to create catalog",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Catalog created successfully"})

}

// GetCatalog godoc
// @Summary Get a catalog by ID
// @Description Get a catalog from the ERP system by their ID
// @Tags Catalogs
// @Accept json
// @Produce json
// @Param id path string true "Catalog ID"
// @Success 200 {object} SuccessResponse{data=entity.Catalog}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/catalogs/{id} [get]
func (h *ERPHandler) GetCatalog(c *gin.Context) {
	// Parse catalog ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid catalog ID format",
			},
		})
		return
	}

	// Get user from context for permission checks
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}
	userID := user.ID

	// Get catalog
	_, err = h.useCase.GetCatalog(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "catalog not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "NOT_FOUND",
					Message: "Catalog not found",
				},
			})
			return
		}

		h.log.Error("Failed to get catalog", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to get catalog",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Catalog listed successfully"})

}

// GetCatalogs godoc
// @Summary Get a list of catalogs
// @Description Get a list of catalogs from the ERP system with pagination
// @Tags Catalogs
// @Accept json
// @Produce json
// @Param merchantID query string false "Merchant ID"
// @Param skip query int false "Number of items to skip" default(0)
// @Param take query int false "Number of items to take" default(10)
// @Success 200 {object} SuccessResponse{data=entity.CatalogsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/catalogs [get]
func (h *ERPHandler) GetCatalogs(c *gin.Context) {
	// Parse query parameters
	var params entity.GetCatalogsParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// Get user and merchant from context
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	// Set default pagination values if not provided
	if params.Skip == 0 {
		params.Skip = 0
	}
	if params.Take == 0 {
		params.Take = 10
	}

	// Get merchant ID from query or context
	merchantID := uuid.Nil
	if params.MerchantID != uuid.Nil {
		merchantID = params.MerchantID
	} else if merchantIDVal, exists := c.Get("merchant_id"); exists {
		merchantID = merchantIDVal.(uuid.UUID)
	}

	// Get catalogs
	catalogs, err := h.useCase.GetCatalogs(c.Request.Context(), entity.GetCatalogsParams{
		MerchantID: merchantID,
		Skip:       params.Skip,
		Take:       params.Take,
	})

	if err != nil {
		h.log.Error("Failed to get catalogs", map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to get catalogs",
			},
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Catalogs retrieved successfully",
		Data:    catalogs,
	})
}

// UpdateCatalog godoc
// @Summary Update an existing catalog
// @Description Update an existing catalog in the ERP system by its ID
// @Tags Catalogs
// @Accept json
// @Produce json
// @Param id path string true "Catalog ID" format(uuid)
// @Param request body entity.UpdateCatalogRequest true "Catalog update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/catalogs/{id} [put]
func (h *ERPHandler) UpdateCatalog(c *gin.Context) {
	// Parse catalog ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Invalid catalog ID format",
			},
		})
		return
	}

	// Parse request body
	var req entity.UpdateCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// Get user from context
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	// Update catalog
	err = h.useCase.UpdateCatalog(c.Request.Context(), id, &req, user.ID)
	if err != nil {
		if err.Error() == "catalog not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error: ApiError{
					Type:    "NOT_FOUND",
					Message: "Catalog not found",
				},
			})
			return
		}

		h.log.Error("Failed to update catalog", map[string]interface{}{
			"error":   err.Error(),
			"id":      id,
			"user_id": user.ID,
		})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to update catalog",
			},
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Catalog updated successfully",
	})
}

// DeleteCatalog godoc
// @Summary Delete a catalog by ID
// @Description Delete a catalog from the ERP system by its ID
// @Tags Catalogs
// @Accept json
// @Produce json
// @Param id path string true "Catalog ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/catalogs/{id} [delete]
func (h *ERPHandler) DeleteCatalog(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.useCase.DeleteCatalog(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Catalog deleted successfully"})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product in the ERP system
// @Tags Products
// @Accept json
// @Produce json
// @Param request body entity.CreateProductRequest true "Product creation request"
// @Success 201 {object} entity.Product
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/products [post]
func (h *ERPHandler) CreateProduct(c *gin.Context) {
	var req entity.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	product, err := h.useCase.CreateProduct(c.Request.Context(), &req, *user.MerchantID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get a product from the ERP system by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/products/{id} [get]
func (h *ERPHandler) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	product, err := h.useCase.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts godoc
// @Summary Get a list of products
// @Description Get a list of products from the ERP system with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param merchantID query string true "Merchant ID"
// @Param skip query int false "Number of items to skip"
// @Param take query int false "Number of items to take"
// @Success 200 {object} entity.ProductsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/products [get]
func (h *ERPHandler) GetProducts(c *gin.Context) {
	var params entity.GetProductsParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}
	merchantID := user.MerchantID
	// Dereference merchantID before assigning it to params.MerchantID
	params.MerchantID = *merchantID

	products, err := h.useCase.GetProducts(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update an existing product in the ERP system by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body entity.UpdateProductRequest true "Product update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/products/{id} [put]
func (h *ERPHandler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req entity.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	if err := h.useCase.UpdateProduct(c.Request.Context(), id, &req, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete a product from the ERP system by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/products/{id} [delete]
func (h *ERPHandler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.useCase.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// CreateWarehouse godoc
// @Summary Create a new warehouse
// @Description Create a new warehouse in the ERP system
// @Tags Warehouses
// @Accept json
// @Produce json
// @Param request body entity.CreateWarehouseRequest true "Warehouse creation request"
// @Success 201 {object} entity.Warehouse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/warehouses [post]
func (h *ERPHandler) CreateWarehouse(c *gin.Context) {
	var req entity.CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}
	warehouse, err := h.useCase.CreateWarehouse(c.Request.Context(), &req, *user.MerchantID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, warehouse)
}

// GetWarehouse godoc
// @Summary Get a warehouse by ID
// @Description Get a warehouse from the ERP system by its ID
// @Tags Warehouses
// @Accept json
// @Produce json
// @Param id path string true "Warehouse ID"
// @Success 200 {object} entity.Warehouse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/warehouses/{id} [get]
func (h *ERPHandler) GetWarehouse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	warehouse, err := h.useCase.GetWarehouse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

// GetWarehouses godoc
// @Summary Get a list of warehouses
// @Description Get a list of warehouses from the ERP system with pagination
// @Tags Warehouses
// @Accept json
// @Produce json
// @Param merchantID query string true "Merchant ID"
// @Param skip query int false "Number of items to skip"
// @Param take query int false "Number of items to take"
// @Success 200 {object} entity.WarehousesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/warehouses [get]
func (h *ERPHandler) GetWarehouses(c *gin.Context) {
	var params entity.GetWarehousesParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}
	merchantID := user.MerchantID
	params.MerchantID = *merchantID

	warehouses, err := h.useCase.GetWarehouses(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, warehouses)
}

// UpdateWarehouse godoc
// @Summary Update an existing warehouse
// @Description Update an existing warehouse in the ERP system by its ID
// @Tags Warehouses
// @Accept json
// @Produce json
// @Param id path string true "Warehouse ID"
// @Param request body entity.UpdateWarehouseRequest true "Warehouse update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/warehouses/{id} [put]
func (h *ERPHandler) UpdateWarehouse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req entity.UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	if err := h.useCase.UpdateWarehouse(c.Request.Context(), id, &req, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse updated successfully"})
}

// DeleteWarehouse godoc
// @Summary Delete a warehouse by ID
// @Description Delete a warehouse from the ERP system by its ID
// @Tags Warehouses
// @Accept json
// @Produce json
// @Param id path string true "Warehouse ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/warehouses/{id} [delete]
func (h *ERPHandler) DeleteWarehouse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.useCase.DeleteWarehouse(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse deleted successfully"})
}

// CreatePaymentMethod godoc
// @Summary Create a new payment method
// @Description Create a new payment method in the ERP system
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Param request body entity.CreatePaymentMethodRequest true "Payment method creation request"
// @Success 201 {object} entity.PaymentMethod
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/payment-methods [post]
func (h *ERPHandler) CreatePaymentMethod(c *gin.Context) {
	var req entity.CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	paymentMethod, err := h.useCase.CreatePaymentMethod(c.Request.Context(), &req, *user.MerchantID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, paymentMethod)
}

// GetPaymentMethod godoc
// @Summary Get a payment method by ID
// @Description Get a payment method from the ERP system by its ID
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Param id path string true "Payment Method ID"
// @Success 200 {object} entity.PaymentMethod
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/payment-methods/{id} [get]
func (h *ERPHandler) GetPaymentMethod(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	paymentMethod, err := h.useCase.GetPaymentMethod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}

	c.JSON(http.StatusOK, paymentMethod)
}

// GetPaymentMethods godoc
// @Summary Get a list of payment methods
// @Description Get a list of payment methods from the ERP system with pagination
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Param merchantID query string true "Merchant ID"
// @Param skip query int false "Number of items to skip"
// @Param take query int false "Number of items to take"
// @Success 200 {object} entity.PaymentMethodsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/payment-methods [get]
func (h *ERPHandler) GetPaymentMethods(c *gin.Context) {
	var params entity.GetPaymentMethodsParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}
	merchantID := user.MerchantID
	params.MerchantID = *merchantID

	paymentMethods, err := h.useCase.GetPaymentMethods(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, paymentMethods)
}

// UpdatePaymentMethod godoc
// @Summary Update an existing payment method
// @Description Update an existing payment method in the ERP system by its ID
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Param id path string true "Payment Method ID"
// @Param request body entity.UpdatePaymentMethodRequest true "Payment method update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/payment-methods/{id} [put]
func (h *ERPHandler) UpdatePaymentMethod(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req entity.UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	if err := h.useCase.UpdatePaymentMethod(c.Request.Context(), id, &req, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method updated successfully"})
}

// DeletePaymentMethod godoc
// @Summary Delete a payment method by ID
// @Description Delete a payment method from the ERP system by its ID
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Param id path string true "Payment Method ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/payment-methods/{id} [delete]
func (h *ERPHandler) DeletePaymentMethod(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.useCase.DeletePaymentMethod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted successfully"})
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order in the ERP system
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body entity.CreateOrderRequest true "Order creation request"
// @Success 201 {object} entity.Order
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/orders [post]
func (h *ERPHandler) CreateOrder(c *gin.Context) {
	var req entity.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (using team basiles existing middleware pattern)
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	order, err := h.useCase.CreateOrder(c.Request.Context(), &req, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get an order from the ERP system by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} entity.Order
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/orders/{id} [get]
func (h *ERPHandler) GetOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	order, err := h.useCase.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrders godoc
// @Summary Get a list of orders
// @Description Get a list of orders from the ERP system with pagination
// @Tags Orders
// @Accept json
// @Produce json
// @Param merchantID query string true "Merchant ID"
// @Param skip query int false "Number of items to skip"
// @Param take query int false "Number of items to take"
// @Success 200 {object} entity.OrdersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/orders [get]
func (h *ERPHandler) GetOrders(c *gin.Context) {
	var params entity.GetOrdersParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}
	params.MerchantID = *user.MerchantID

	orders, err := h.useCase.GetOrders(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrder godoc
// @Summary Update an existing order
// @Description Update an existing order in the ERP system by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body entity.UpdateOrderRequest true "Order update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/orders/{id} [put]
func (h *ERPHandler) UpdateOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req entity.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, exists := ginn.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}
	userID := user.ID
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.useCase.UpdateOrder(c.Request.Context(), id, &req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

// DeleteOrder godoc
// @Summary Delete an order by ID
// @Description Delete an order from the ERP system by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /erp/orders/{id} [delete]
func (h *ERPHandler) DeleteOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.useCase.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
