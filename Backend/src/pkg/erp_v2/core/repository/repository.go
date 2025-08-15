package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/erp_v2/core/entity"
)

// Repository defines the interface for ERP data access
type Repository interface {
	// Account operations
	CreateAccount(ctx context.Context, account *entity.Account) error
	GetAccount(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, account *entity.Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error

	// Catalog operations
	CreateCatalog(ctx context.Context, catalog *entity.Catalog) error
	GetCatalog(ctx context.Context, id uuid.UUID) (*entity.Catalog, error)
	GetCatalogs(ctx context.Context, params entity.GetCatalogsParams) (*entity.CatalogsResponse, error)
	UpdateCatalog(ctx context.Context, id uuid.UUID, req *entity.UpdateCatalogRequest) error
	DeleteCatalog(ctx context.Context, id uuid.UUID) error

	// Customer operations
	CreateCustomer(ctx context.Context, customer *entity.Customer) error
	GetCustomer(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
	GetCustomers(ctx context.Context, params entity.GetCustomersParams) (*entity.CustomersResponse, error)
	UpdateCustomer(ctx context.Context, id uuid.UUID, req *entity.UpdateCustomerRequest) error
	DeleteCustomer(ctx context.Context, id uuid.UUID) error

	// Order operations
	CreateOrder(ctx context.Context, order *entity.Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	GetOrders(ctx context.Context, params entity.GetOrdersParams) (*entity.OrdersResponse, error)
	UpdateOrder(ctx context.Context, id uuid.UUID, req *entity.UpdateOrderRequest) error
	DeleteOrder(ctx context.Context, id uuid.UUID) error

	// Payment Method operations
	CreatePaymentMethod(ctx context.Context, paymentMethod *entity.PaymentMethod) error
	GetPaymentMethod(ctx context.Context, id uuid.UUID) (*entity.PaymentMethod, error)
	GetPaymentMethods(ctx context.Context, params entity.GetPaymentMethodsParams) (*entity.PaymentMethodsResponse, error)
	UpdatePaymentMethod(ctx context.Context, id uuid.UUID, req *entity.UpdatePaymentMethodRequest) error
	DeletePaymentMethod(ctx context.Context, id uuid.UUID) error

	// Product operations
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProduct(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	GetProducts(ctx context.Context, params entity.GetProductsParams) (*entity.ProductsResponse, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req *entity.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error

	// Warehouse operations
	CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error
	GetWarehouse(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error)
	GetWarehouses(ctx context.Context, params entity.GetWarehousesParams) (*entity.WarehousesResponse, error)
	UpdateWarehouse(ctx context.Context, id uuid.UUID, req *entity.UpdateWarehouseRequest) error
	DeleteWarehouse(ctx context.Context, id uuid.UUID) error
}

