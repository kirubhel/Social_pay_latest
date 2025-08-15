package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/erp_v2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/erp_v2/core/repository"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
)

// ERPUseCase defines the interface for ERP business logic
type ERPUseCase interface {
	// Account operations
	CreateAccount(ctx context.Context, account *entity.Account) error
	GetAccount(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, account *entity.Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error

	// Catalog operations
	CreateCatalog(ctx context.Context, req *entity.CreateCatalogRequest, merchantID uuid.UUID, userID uuid.UUID) (*entity.CatalogResponse, error)
	GetCatalog(ctx context.Context, id uuid.UUID) (*entity.CatalogResponse, error)
	GetCatalogs(ctx context.Context, params entity.GetCatalogsParams) (*entity.CatalogsResponse, error)
	UpdateCatalog(ctx context.Context, id uuid.UUID, req *entity.UpdateCatalogRequest, userID uuid.UUID) error
	DeleteCatalog(ctx context.Context, id uuid.UUID) error

	// Customer operations
	CreateCustomer(ctx context.Context, req *entity.CreateCustomerRequest, userID uuid.UUID) (*entity.CustomerResponse, error)
	GetCustomer(ctx context.Context, id uuid.UUID) (*entity.CustomerResponse, error)
	GetCustomers(ctx context.Context, params entity.GetCustomersParams) (*entity.CustomersResponse, error)
	UpdateCustomer(ctx context.Context, id uuid.UUID, req *entity.UpdateCustomerRequest, userID uuid.UUID) error
	DeleteCustomer(ctx context.Context, id uuid.UUID) error

	// Order operations
	CreateOrder(ctx context.Context, req *entity.CreateOrderRequest, userID uuid.UUID) (*entity.OrderResponse, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*entity.OrderResponse, error)
	GetOrders(ctx context.Context, params entity.GetOrdersParams) (*entity.OrdersResponse, error)
	UpdateOrder(ctx context.Context, id uuid.UUID, req *entity.UpdateOrderRequest, userID uuid.UUID) error
	DeleteOrder(ctx context.Context, id uuid.UUID) error

	// Payment Method operations
	CreatePaymentMethod(ctx context.Context, req *entity.CreatePaymentMethodRequest, merchantID uuid.UUID, userID uuid.UUID) (*entity.PaymentMethodResponse, error)
	GetPaymentMethod(ctx context.Context, id uuid.UUID) (*entity.PaymentMethodResponse, error)
	GetPaymentMethods(ctx context.Context, params entity.GetPaymentMethodsParams) (*entity.PaymentMethodsResponse, error)
	UpdatePaymentMethod(ctx context.Context, id uuid.UUID, req *entity.UpdatePaymentMethodRequest, userID uuid.UUID) error
	DeletePaymentMethod(ctx context.Context, id uuid.UUID) error

	// Product operations
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest, merchantID uuid.UUID, userID uuid.UUID) (*entity.ProductResponse, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*entity.ProductResponse, error)
	GetProducts(ctx context.Context, params entity.GetProductsParams) (*entity.ProductsResponse, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req *entity.UpdateProductRequest, userID uuid.UUID) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error

	// Warehouse operations
	CreateWarehouse(ctx context.Context, req *entity.CreateWarehouseRequest, merchantID uuid.UUID, userID uuid.UUID) (*entity.WarehouseResponse, error)
	GetWarehouse(ctx context.Context, id uuid.UUID) (*entity.WarehouseResponse, error)
	GetWarehouses(ctx context.Context, params entity.GetWarehousesParams) (*entity.WarehousesResponse, error)
	UpdateWarehouse(ctx context.Context, id uuid.UUID, req *entity.UpdateWarehouseRequest, userID uuid.UUID) error
	DeleteWarehouse(ctx context.Context, id uuid.UUID) error
}

type erpUseCase struct {
	repo repository.Repository
	log  logging.Logger
}

// NewERPUseCase creates a new ERP use case
func NewERPUseCase(repo repository.Repository) ERPUseCase {
	return &erpUseCase{
		repo: repo,
		log:  logging.NewStdLogger("[ERP_V2]"),
	}
}

// Account operations
func (u *erpUseCase) CreateAccount(ctx context.Context, account *entity.Account) error {
	return u.repo.CreateAccount(ctx, account)
}

func (u *erpUseCase) GetAccount(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	return u.repo.GetAccount(ctx, id)
}

func (u *erpUseCase) UpdateAccount(ctx context.Context, id uuid.UUID, account *entity.Account) error {
	return u.repo.UpdateAccount(ctx, id, account)
}

func (u *erpUseCase) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteAccount(ctx, id)
}

// Catalog operations
func (u *erpUseCase) CreateCatalog(ctx context.Context, req *entity.CreateCatalogRequest, merchantID uuid.UUID, userID uuid.UUID) (*entity.CatalogResponse, error) {
	catalog := &entity.Catalog{
		ID:          uuid.New(),
		MerchantID:  merchantID,
		Name:        req.Name,
		Description: req.Description,
		Status:      entity.CatalogStatus(req.Status),
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	err := u.repo.CreateCatalog(ctx, catalog)
	if err != nil {
		return nil, err
	}

	return &entity.CatalogResponse{
		ID:          catalog.ID,
		MerchantID:  catalog.MerchantID,
		Name:        catalog.Name,
		Description: catalog.Description,
		Status:      catalog.Status,
		CreatedAt:   catalog.CreatedAt,
		UpdatedAt:   catalog.UpdatedAt,
		CreatedBy:   catalog.CreatedBy,
		UpdatedBy:   catalog.UpdatedBy,
	}, nil
}

func (u *erpUseCase) GetCatalog(ctx context.Context, id uuid.UUID) (*entity.CatalogResponse, error) {
	catalog, err := u.repo.GetCatalog(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.CatalogResponse{
		ID:          catalog.ID,
		MerchantID:  catalog.MerchantID,
		Name:        catalog.Name,
		Description: catalog.Description,
		Status:      catalog.Status,
		CreatedAt:   catalog.CreatedAt,
		UpdatedAt:   catalog.UpdatedAt,
		CreatedBy:   catalog.CreatedBy,
		UpdatedBy:   catalog.UpdatedBy,
	}, nil
}

func (u *erpUseCase) GetCatalogs(ctx context.Context, params entity.GetCatalogsParams) (*entity.CatalogsResponse, error) {
	return u.repo.GetCatalogs(ctx, params)
}

func (u *erpUseCase) UpdateCatalog(ctx context.Context, id uuid.UUID, req *entity.UpdateCatalogRequest, userID uuid.UUID) error {
	return u.repo.UpdateCatalog(ctx, id, req)
}

func (u *erpUseCase) DeleteCatalog(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteCatalog(ctx, id)
}

// Customer operations
func (u *erpUseCase) CreateCustomer(ctx context.Context, req *entity.CreateCustomerRequest, userID uuid.UUID) (*entity.CustomerResponse, error) {
	merchantUUID, err := uuid.Parse(req.MerchantID)
	if err != nil {
		return nil, err
	}

	customer := &entity.Customer{
		ID:            uuid.New(),
		CustomerID:    uuid.New(),
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.PhoneNumber,
		Address:       req.Address,
		LoyaltyPoints: req.LoyaltyPoints,
		DateOfBirth:   req.DateOfBirth,
		Status:        entity.CustomerStatus(req.Status),
		CreatedBy:     userID,
		UpdatedBy:     userID,
		MerchantID:    merchantUUID,
	}

	err = u.repo.CreateCustomer(ctx, customer)
	if err != nil {
		return nil, err
	}

	return &entity.CustomerResponse{
		ID:            customer.ID,
		CustomerID:    customer.CustomerID,
		Name:          customer.Name,
		Email:         customer.Email,
		Phone:         customer.Phone,
		Address:       customer.Address,
		LoyaltyPoints: customer.LoyaltyPoints,
		DateOfBirth:   customer.DateOfBirth,
		Status:        customer.Status,
		CreatedAt:     customer.CreatedAt,
		UpdatedAt:     customer.UpdatedAt,
		CreatedBy:     customer.CreatedBy,
		UpdatedBy:     customer.UpdatedBy,
		MerchantID:    customer.MerchantID,
	}, nil
}

func (u *erpUseCase) GetCustomer(ctx context.Context, id uuid.UUID) (*entity.CustomerResponse, error) {
	customer, err := u.repo.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.CustomerResponse{
		ID:            customer.ID,
		CustomerID:    customer.CustomerID,
		Name:          customer.Name,
		Email:         customer.Email,
		Phone:         customer.Phone,
		Address:       customer.Address,
		LoyaltyPoints: customer.LoyaltyPoints,
		DateOfBirth:   customer.DateOfBirth,
		Status:        customer.Status,
		CreatedAt:     customer.CreatedAt,
		UpdatedAt:     customer.UpdatedAt,
		CreatedBy:     customer.CreatedBy,
		UpdatedBy:     customer.UpdatedBy,
		MerchantID:    customer.MerchantID,
	}, nil
}

func (u *erpUseCase) GetCustomers(ctx context.Context, params entity.GetCustomersParams) (*entity.CustomersResponse, error) {
	return u.repo.GetCustomers(ctx, params)
}

func (u *erpUseCase) UpdateCustomer(ctx context.Context, id uuid.UUID, req *entity.UpdateCustomerRequest, userID uuid.UUID) error {
	return u.repo.UpdateCustomer(ctx, id, req)
}

func (u *erpUseCase) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteCustomer(ctx, id)
}

// Order operations
func (u *erpUseCase) CreateOrder(ctx context.Context, req *entity.CreateOrderRequest, userID uuid.UUID) (*entity.OrderResponse, error) {
	order := &entity.Order{
		ID:              uuid.New(),
		CustomerDetails: req.CustomerDetails,
		OrderDetails:    req.OrderDetails,
		OrderItems:      req.OrderItems,
		Metadata:        req.Metadata,
		Tracking:        req.Tracking,
		MerchantID:      req.MerchantID,
	}

	err := u.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return &entity.OrderResponse{
		ID:              order.ID,
		CustomerDetails: order.CustomerDetails,
		OrderDetails:    order.OrderDetails,
		OrderItems:      order.OrderItems,
		Metadata:        order.Metadata,
		Tracking:        order.Tracking,
		MerchantID:      order.MerchantID,
	}, nil
}

func (u *erpUseCase) GetOrder(ctx context.Context, id uuid.UUID) (*entity.OrderResponse, error) {
	order, err := u.repo.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.OrderResponse{
		ID:              order.ID,
		CustomerDetails: order.CustomerDetails,
		OrderDetails:    order.OrderDetails,
		OrderItems:      order.OrderItems,
		Metadata:        order.Metadata,
		Tracking:        order.Tracking,
		MerchantID:      order.MerchantID,
	}, nil
}

func (u *erpUseCase) GetOrders(ctx context.Context, params entity.GetOrdersParams) (*entity.OrdersResponse, error) {
	return u.repo.GetOrders(ctx, params)
}

func (u *erpUseCase) UpdateOrder(ctx context.Context, id uuid.UUID, req *entity.UpdateOrderRequest, userID uuid.UUID) error {
	return u.repo.UpdateOrder(ctx, id, req)
}

func (u *erpUseCase) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteOrder(ctx, id)
}

// Payment Method operations
func (u *erpUseCase) CreatePaymentMethod(ctx context.Context, req *entity.CreatePaymentMethodRequest, merchantID uuid.UUID, userID uuid.UUID) (*entity.PaymentMethodResponse, error) {
	paymentMethod := &entity.PaymentMethod{
		ID:         uuid.New(),
		MerchantID: merchantID,
		Name:       req.Name,
		Type:       req.Type,
		Commission: req.Commission,
		Details:    req.Details,
		IsActive:   req.IsActive,
		CreatedBy:  userID,
		UpdatedBy:  userID,
	}

	err := u.repo.CreatePaymentMethod(ctx, paymentMethod)
	if err != nil {
		return nil, err
	}

	return &entity.PaymentMethodResponse{
		ID:         paymentMethod.ID,
		MerchantID: paymentMethod.MerchantID,
		Name:       paymentMethod.Name,
		Type:       paymentMethod.Type,
		Commission: paymentMethod.Commission,
		Details:    paymentMethod.Details,
		IsActive:   paymentMethod.IsActive,
		CreatedAt:  paymentMethod.CreatedAt,
		UpdatedAt:  paymentMethod.UpdatedAt,
		CreatedBy:  paymentMethod.CreatedBy,
		UpdatedBy:  paymentMethod.UpdatedBy,
	}, nil
}

func (u *erpUseCase) GetPaymentMethod(ctx context.Context, id uuid.UUID) (*entity.PaymentMethodResponse, error) {
	paymentMethod, err := u.repo.GetPaymentMethod(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.PaymentMethodResponse{
		ID:         paymentMethod.ID,
		MerchantID: paymentMethod.MerchantID,
		Name:       paymentMethod.Name,
		Type:       paymentMethod.Type,
		Commission: paymentMethod.Commission,
		Details:    paymentMethod.Details,
		IsActive:   paymentMethod.IsActive,
		CreatedAt:  paymentMethod.CreatedAt,
		UpdatedAt:  paymentMethod.UpdatedAt,
		CreatedBy:  paymentMethod.CreatedBy,
		UpdatedBy:  paymentMethod.UpdatedBy,
	}, nil
}

func (u *erpUseCase) GetPaymentMethods(ctx context.Context, params entity.GetPaymentMethodsParams) (*entity.PaymentMethodsResponse, error) {
	return u.repo.GetPaymentMethods(ctx, params)
}

func (u *erpUseCase) UpdatePaymentMethod(ctx context.Context, id uuid.UUID, req *entity.UpdatePaymentMethodRequest, userID uuid.UUID) error {
	return u.repo.UpdatePaymentMethod(ctx, id, req)
}

func (u *erpUseCase) DeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeletePaymentMethod(ctx, id)
}

// Product operations
func (u *erpUseCase) CreateProduct(ctx context.Context, req *entity.CreateProductRequest, merchantID uuid.UUID, userID uuid.UUID) (*entity.ProductResponse, error) {
	product := &entity.Product{
		ID:          uuid.New(),
		MerchantID:  merchantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		SKU:         req.SKU,
		Weight:      req.Weight,
		Dimensions:  req.Dimensions,
		ImageURL:    req.ImageURL,
		Status:      entity.ProductStatus(req.Status),
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	err := u.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &entity.ProductResponse{
		ID:          product.ID,
		MerchantID:  product.MerchantID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Currency:    product.Currency,
		SKU:         product.SKU,
		Weight:      product.Weight,
		Dimensions:  product.Dimensions,
		ImageURL:    product.ImageURL,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		CreatedBy:   product.CreatedBy,
		UpdatedBy:   product.UpdatedBy,
	}, nil
}

func (u *erpUseCase) GetProduct(ctx context.Context, id uuid.UUID) (*entity.ProductResponse, error) {
	product, err := u.repo.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.ProductResponse{
		ID:          product.ID,
		MerchantID:  product.MerchantID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Currency:    product.Currency,
		SKU:         product.SKU,
		Weight:      product.Weight,
		Dimensions:  product.Dimensions,
		ImageURL:    product.ImageURL,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		CreatedBy:   product.CreatedBy,
		UpdatedBy:   product.UpdatedBy,
	}, nil
}

func (u *erpUseCase) GetProducts(ctx context.Context, params entity.GetProductsParams) (*entity.ProductsResponse, error) {
	return u.repo.GetProducts(ctx, params)
}

func (u *erpUseCase) UpdateProduct(ctx context.Context, id uuid.UUID, req *entity.UpdateProductRequest, userID uuid.UUID) error {
	return u.repo.UpdateProduct(ctx, id, req)
}

func (u *erpUseCase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteProduct(ctx, id)
}

// Warehouse operations
func (u *erpUseCase) CreateWarehouse(ctx context.Context, req *entity.CreateWarehouseRequest, merchantID uuid.UUID, userID uuid.UUID) (*entity.WarehouseResponse, error) {
	warehouse := &entity.Warehouse{
		ID:          uuid.New(),
		MerchantID:  merchantID,
		Name:        req.Name,
		Location:    req.Location,
		Capacity:    req.Capacity,
		CreatedBy:   userID,
		UpdatedBy:   userID,
		IsActive:    req.IsActive,
		Description: req.Description,
		Status:      entity.WarehouseStatus(req.Status),
	}

	err := u.repo.CreateWarehouse(ctx, warehouse)
	if err != nil {
		return nil, err
	}

	return &entity.WarehouseResponse{
		ID:          warehouse.ID,
		MerchantID:  warehouse.MerchantID,
		Name:        warehouse.Name,
		Location:    warehouse.Location,
		Capacity:    warehouse.Capacity,
		CreatedAt:   warehouse.CreatedAt,
		UpdatedAt:   warehouse.UpdatedAt,
		CreatedBy:   warehouse.CreatedBy,
		UpdatedBy:   warehouse.UpdatedBy,
		IsActive:    warehouse.IsActive,
		Description: warehouse.Description,
		Status:      warehouse.Status,
	}, nil
}

func (u *erpUseCase) GetWarehouse(ctx context.Context, id uuid.UUID) (*entity.WarehouseResponse, error) {
	warehouse, err := u.repo.GetWarehouse(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.WarehouseResponse{
		ID:          warehouse.ID,
		MerchantID:  warehouse.MerchantID,
		Name:        warehouse.Name,
		Location:    warehouse.Location,
		Capacity:    warehouse.Capacity,
		CreatedAt:   warehouse.CreatedAt,
		UpdatedAt:   warehouse.UpdatedAt,
		CreatedBy:   warehouse.CreatedBy,
		UpdatedBy:   warehouse.UpdatedBy,
		IsActive:    warehouse.IsActive,
		Description: warehouse.Description,
		Status:      warehouse.Status,
	}, nil
}

func (u *erpUseCase) GetWarehouses(ctx context.Context, params entity.GetWarehousesParams) (*entity.WarehousesResponse, error) {
	return u.repo.GetWarehouses(ctx, params)
}

func (u *erpUseCase) UpdateWarehouse(ctx context.Context, id uuid.UUID, req *entity.UpdateWarehouseRequest, userID uuid.UUID) error {
	return u.repo.UpdateWarehouse(ctx, id, req)
}

func (u *erpUseCase) DeleteWarehouse(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteWarehouse(ctx, id)
}
