package sqlc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/erp_v2/core/entity"
	"github.com/socialpay/socialpay/src/pkg/erp_v2/core/repository"
	"github.com/socialpay/socialpay/src/pkg/types"
)

type sqlcRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewSQLCRepository creates a new SQLC-based ERP repository
func NewSQLCRepository(db *sql.DB) repository.Repository {
	return &sqlcRepository{
		db:      db,
		queries: New(db),
	}
}

// Account operations
func (r *sqlcRepository) CreateAccount(ctx context.Context, account *entity.Account) error {
	verificationStatusJSON, err := json.Marshal(account.VerificationStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal verification status: %w", err)
	}

	err = r.queries.CreateAccount(ctx, CreateAccountParams{
		ID:                 account.ID,
		UserID:             account.UserID,
		Title:              account.Title,
		Type:               string(account.Type),
		DefaultAccount:     account.Default,
		VerificationStatus: verificationStatusJSON,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	})
	return err
}

func (r *sqlcRepository) GetAccount(ctx context.Context, id uuid.UUID) (*entity.Account, error) { // Corrected line
	account, err := r.queries.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	var verificationStatus entity.Account_VerificationStatus
	if err := json.Unmarshal(account.VerificationStatus, &verificationStatus); err != nil {
		return nil, fmt.Errorf("failed to unmarshal verification status: %w", err)
	}

	return &entity.Account{
		ID:                 account.ID,
		UserID:             account.UserID,
		Title:              account.Title,
		Type:               entity.AccountType(account.Type),
		Default:            account.DefaultAccount,
		VerificationStatus: verificationStatus,
		CreatedAt:          account.CreatedAt,
		UpdatedAt:          account.UpdatedAt,
	}, nil
}

func (r *sqlcRepository) UpdateAccount(ctx context.Context, id uuid.UUID, account *entity.Account) error {
	verificationStatusJSON, err := json.Marshal(account.VerificationStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal verification status: %w", err)
	}

	err = r.queries.UpdateAccount(ctx, UpdateAccountParams{
		ID:                 id,
		Title:              account.Title,
		Type:               string(account.Type),
		DefaultAccount:     account.Default,
		VerificationStatus: verificationStatusJSON,
		UpdatedAt:          time.Now(),
	})
	return err
}

func (r *sqlcRepository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteAccount(ctx, id)
	return err
}

func (r *sqlcRepository) CreateCatalog(ctx context.Context, catalog *entity.Catalog) error {
	sqlDesc := catalog.Description.ToSqlNullString()

	err := r.queries.CreateCatalog(ctx, CreateCatalogParams{
		ID:          catalog.ID,
		MerchantID:  catalog.MerchantID,
		Name:        catalog.Name,
		Description: sqlDesc,
		Status:      string(catalog.Status),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   catalog.CreatedBy,
		UpdatedBy:   catalog.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) GetCatalog(ctx context.Context, id uuid.UUID) (*entity.Catalog, error) {
	dbCatalog, err := r.queries.GetCatalog(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Catalog{
		ID:          dbCatalog.ID,
		MerchantID:  dbCatalog.MerchantID,
		Name:        dbCatalog.Name,
		Description: types.FromSqlNullString(dbCatalog.Description),
		Status:      entity.CatalogStatus(dbCatalog.Status),
		CreatedAt:   dbCatalog.CreatedAt,
		UpdatedAt:   dbCatalog.UpdatedAt,
		CreatedBy:   dbCatalog.CreatedBy,
		UpdatedBy:   dbCatalog.UpdatedBy,
	}, nil
}

func (r *sqlcRepository) GetCatalogs(ctx context.Context, params entity.GetCatalogsParams) (*entity.CatalogsResponse, error) {
	dbCatalogs, err := r.queries.GetCatalogs(ctx, GetCatalogsParams{
		MerchantID: params.MerchantID,
		Limit:      int32(params.Take),
		Offset:     int32(params.Skip),
	})
	if err != nil {
		return nil, err
	}

	var catalogResponses []entity.CatalogResponse
	for _, c := range dbCatalogs {
		catalogResponses = append(catalogResponses, entity.CatalogResponse{
			ID:          c.ID,
			MerchantID:  c.MerchantID,
			Name:        c.Name,
			Description: types.FromSqlNullString(c.Description),
			Status:      entity.CatalogStatus(c.Status),
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
			CreatedBy:   c.CreatedBy,
			UpdatedBy:   c.UpdatedBy,
		})
	}

	return &entity.CatalogsResponse{
		Catalogs: catalogResponses,
		Count:    len(catalogResponses),
	}, nil
}

func (r *sqlcRepository) UpdateCatalog(ctx context.Context, id uuid.UUID, req *entity.UpdateCatalogRequest) error {
	catalog, err := r.queries.GetCatalog(ctx, id)
	if err != nil {
		return err
	}

	// Initialize with existing values
	name := catalog.Name
	description := catalog.Description
	status := catalog.Status

	// Apply updates from request
	if req.Name != nil {
		name = *req.Name
	}

	if req.Description.Valid {
		description = req.Description.ToSqlNullString()
	}

	if req.Status != nil {
		status = string(*req.Status)
	}

	err = r.queries.UpdateCatalog(ctx, UpdateCatalogParams{
		ID:          id,
		Name:        name,
		Description: description,
		Status:      status,
		UpdatedAt:   time.Now(),
		UpdatedBy:   catalog.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) DeleteCatalog(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteCatalog(ctx, id)
	return err
}

// Customer operations
func (r *sqlcRepository) CreateCustomer(ctx context.Context, customer *entity.Customer) error {
	err := r.queries.CreateCustomer(ctx, CreateCustomerParams{
		ID:            customer.ID,
		CustomerID:    customer.CustomerID,
		MerchantID:    customer.MerchantID,
		Name:          customer.Name,
		Email:         customer.Email,
		Phone:         sql.NullString{String: customer.Phone, Valid: customer.Phone != ""},
		Address:       sql.NullString{String: customer.Address, Valid: customer.Address != ""},
		LoyaltyPoints: sql.NullInt32{Int32: int32(customer.LoyaltyPoints), Valid: true},
		DateOfBirth:   sql.NullString{String: customer.DateOfBirth, Valid: customer.DateOfBirth != ""},
		Status:        string(customer.Status),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		CreatedBy:     customer.CreatedBy,
		UpdatedBy:     customer.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) GetCustomer(ctx context.Context, id uuid.UUID) (*entity.Customer, error) {
	customer, err := r.queries.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Customer{
		ID:            customer.ID,
		CustomerID:    customer.CustomerID,
		MerchantID:    customer.MerchantID,
		Name:          customer.Name,
		Email:         customer.Email,
		Phone:         customer.Phone.String,
		Address:       customer.Address.String,
		LoyaltyPoints: int(customer.LoyaltyPoints.Int32),
		DateOfBirth:   customer.DateOfBirth.String,
		Status:        entity.CustomerStatus(customer.Status),
		CreatedAt:     customer.CreatedAt,
		UpdatedAt:     customer.UpdatedAt,
		CreatedBy:     customer.CreatedBy,
		UpdatedBy:     customer.UpdatedBy,
	}, nil
}

func (r *sqlcRepository) GetCustomers(ctx context.Context, params entity.GetCustomersParams) (*entity.CustomersResponse, error) {
	customers, err := r.queries.GetCustomers(ctx, GetCustomersParams{
		MerchantID: params.MerchantID,
		Limit:      int32(params.Take),
		Offset:     int32(params.Skip),
	})
	if err != nil {
		return nil, err
	}

	var customerResponses []entity.CustomerResponse
	for _, c := range customers {
		customerResponses = append(customerResponses, entity.CustomerResponse{
			ID:            c.ID,
			CustomerID:    c.CustomerID,
			MerchantID:    c.MerchantID,
			Name:          c.Name,
			Email:         c.Email,
			Phone:         c.Phone.String,
			Address:       c.Address.String,
			LoyaltyPoints: int(c.LoyaltyPoints.Int32),
			DateOfBirth:   c.DateOfBirth.String,
			Status:        entity.CustomerStatus(c.Status),
			CreatedAt:     c.CreatedAt,
			UpdatedAt:     c.UpdatedAt,
			CreatedBy:     c.CreatedBy,
			UpdatedBy:     c.UpdatedBy,
		})
	}

	return &entity.CustomersResponse{
		Count:     len(customerResponses),
		Customers: customerResponses,
	}, nil
}

func (r *sqlcRepository) UpdateCustomer(ctx context.Context, id uuid.UUID, req *entity.UpdateCustomerRequest) error {
	customer, err := r.queries.GetCustomer(ctx, id)
	if err != nil {
		return err
	}

	name := customer.Name
	if req.Name != nil {
		name = *req.Name
	}
	email := customer.Email
	if req.Email != nil {
		email = *req.Email
	}
	phone := customer.Phone
	if req.Phone != nil {
		phone = sql.NullString{String: *req.Phone, Valid: true}
	}
	address := customer.Address
	if req.Address != nil {
		address = sql.NullString{String: *req.Address, Valid: true}
	}
	loyaltyPoints := customer.LoyaltyPoints
	if req.LoyaltyPoints != nil {
		loyaltyPoints = sql.NullInt32{Int32: int32(*req.LoyaltyPoints), Valid: true}
	}
	dateOfBirth := customer.DateOfBirth
	if req.DateOfBirth != nil {
		dateOfBirth = sql.NullString{String: *req.DateOfBirth, Valid: true}
	}
	status := customer.Status
	if req.Status != nil {
		status = *req.Status
	}

	err = r.queries.UpdateCustomer(ctx, UpdateCustomerParams{
		ID:            id,
		Name:          name,
		Email:         email,
		Phone:         phone,
		Address:       address,
		LoyaltyPoints: loyaltyPoints,
		DateOfBirth:   dateOfBirth,
		Status:        status,
		UpdatedAt:     time.Now(),
		UpdatedBy:     customer.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteCustomer(ctx, id)
	return err
}

// Order operations
func (r *sqlcRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	customerDetailsJSON, err := json.Marshal(order.CustomerDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal customer details: %w", err)
	}
	orderDetailsJSON, err := json.Marshal(order.OrderDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal order details: %w", err)
	}
	orderItemsJSON, err := json.Marshal(order.OrderItems)
	if err != nil {
		return fmt.Errorf("failed to marshal order items: %w", err)
	}
	metadataJSON, err := json.Marshal(order.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	trackingJSON, err := json.Marshal(order.Tracking)
	if err != nil {
		return fmt.Errorf("failed to marshal tracking: %w", err)
	}

	err = r.queries.CreateOrder(ctx, CreateOrderParams{
		ID:              order.ID,
		MerchantID:      order.MerchantID,
		CustomerDetails: customerDetailsJSON,
		OrderDetails:    orderDetailsJSON,
		OrderItems:      orderItemsJSON,
		Metadata:        metadataJSON,
		Tracking:        trackingJSON,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	return err
}

func (r *sqlcRepository) GetOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	order, err := r.queries.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	var customerDetails entity.CustomerDetails
	if err := json.Unmarshal(order.CustomerDetails, &customerDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer details: %w", err)
	}
	var orderDetails entity.OrderDetails
	if err := json.Unmarshal(order.OrderDetails, &orderDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order details: %w", err)
	}
	var orderItems []entity.OrderItem
	if err := json.Unmarshal(order.OrderItems, &orderItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order items: %w", err)
	}
	var metadata entity.Metadata
	if err := json.Unmarshal(order.Metadata, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}
	var tracking entity.Tracking
	if err := json.Unmarshal(order.Tracking, &tracking); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tracking: %w", err)
	}

	return &entity.Order{
		ID:              order.ID,
		MerchantID:      order.MerchantID,
		CustomerDetails: customerDetails,
		OrderDetails:    orderDetails,
		OrderItems:      orderItems,
		Metadata:        metadata,
		Tracking:        tracking,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}, nil
}

func (r *sqlcRepository) GetOrders(ctx context.Context, params entity.GetOrdersParams) (*entity.OrdersResponse, error) {
	orders, err := r.queries.GetOrders(ctx, GetOrdersParams{
		MerchantID: params.MerchantID,
		Limit:      int32(params.Take),
		Offset:     int32(params.Skip),
	})
	if err != nil {
		return nil, err
	}

	var orderResponses []entity.OrderResponse
	for _, o := range orders {
		var customerDetails entity.CustomerDetails
		if err := json.Unmarshal(o.CustomerDetails, &customerDetails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal customer details: %w", err)
		}
		var orderDetails entity.OrderDetails
		if err := json.Unmarshal(o.OrderDetails, &orderDetails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order details: %w", err)
		}
		var orderItems []entity.OrderItem
		if err := json.Unmarshal(o.OrderItems, &orderItems); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order items: %w", err)
		}
		var metadata entity.Metadata
		if err := json.Unmarshal(o.Metadata, &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		var tracking entity.Tracking
		if err := json.Unmarshal(o.Tracking, &tracking); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tracking: %w", err)
		}

		orderResponses = append(orderResponses, entity.OrderResponse{
			ID:              o.ID,
			MerchantID:      o.MerchantID,
			CustomerDetails: customerDetails,
			OrderDetails:    orderDetails,
			OrderItems:      orderItems,
			Metadata:        metadata,
			Tracking:        tracking,
		})
	}

	return &entity.OrdersResponse{
		Count:  len(orderResponses),
		Orders: orderResponses,
	}, nil
}

func (r *sqlcRepository) UpdateOrder(ctx context.Context, id uuid.UUID, req *entity.UpdateOrderRequest) error {
	order, err := r.queries.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	customerDetails := order.CustomerDetails
	if req.CustomerDetails != nil {
		cd, err := json.Marshal(req.CustomerDetails)
		if err != nil {
			return fmt.Errorf("failed to marshal customer details: %w", err)
		}
		customerDetails = cd
	}
	orderDetails := order.OrderDetails
	if req.OrderDetails != nil {
		od, err := json.Marshal(req.OrderDetails)
		if err != nil {
			return fmt.Errorf("failed to marshal order details: %w", err)
		}
		orderDetails = od
	}
	orderItems := order.OrderItems
	if req.OrderItems != nil {
		oi, err := json.Marshal(req.OrderItems)
		if err != nil {
			return fmt.Errorf("failed to marshal order items: %w", err)
		}
		orderItems = oi
	}
	metadata := order.Metadata
	if req.Metadata != nil {
		md, err := json.Marshal(req.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadata = md
	}
	tracking := order.Tracking
	if req.Tracking != nil {
		tr, err := json.Marshal(req.Tracking)
		if err != nil {
			return fmt.Errorf("failed to marshal tracking: %w", err)
		}
		tracking = tr
	}

	err = r.queries.UpdateOrder(ctx, UpdateOrderParams{
		ID:              id,
		CustomerDetails: customerDetails,
		OrderDetails:    orderDetails,
		OrderItems:      orderItems,
		Metadata:        metadata,
		Tracking:        tracking,
		UpdatedAt:       time.Now(),
	})
	return err
}

func (r *sqlcRepository) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteOrder(ctx, id)
	return err
}

// Payment Method operations
func (r *sqlcRepository) CreatePaymentMethod(ctx context.Context, paymentMethod *entity.PaymentMethod) error {
	err := r.queries.CreatePaymentMethod(ctx, CreatePaymentMethodParams{
		ID:         paymentMethod.ID,
		MerchantID: paymentMethod.MerchantID,
		Name:       paymentMethod.Name,
		Type:       paymentMethod.Type,
		Commission: paymentMethod.Commission,
		Details:    sql.NullString{String: paymentMethod.Details, Valid: paymentMethod.Details != ""},
		IsActive:   paymentMethod.IsActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		CreatedBy:  paymentMethod.CreatedBy,
		UpdatedBy:  paymentMethod.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) GetPaymentMethod(ctx context.Context, id uuid.UUID) (*entity.PaymentMethod, error) {
	paymentMethod, err := r.queries.GetPaymentMethod(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.PaymentMethod{
		ID:         paymentMethod.ID,
		MerchantID: paymentMethod.MerchantID,
		Name:       paymentMethod.Name,
		Type:       paymentMethod.Type,
		Commission: paymentMethod.Commission,
		Details:    paymentMethod.Details.String,
		IsActive:   paymentMethod.IsActive,
		CreatedAt:  paymentMethod.CreatedAt,
		UpdatedAt:  paymentMethod.UpdatedAt,
		CreatedBy:  paymentMethod.CreatedBy,
		UpdatedBy:  paymentMethod.UpdatedBy,
	}, nil
}

func (r *sqlcRepository) GetPaymentMethods(ctx context.Context, params entity.GetPaymentMethodsParams) (*entity.PaymentMethodsResponse, error) {
	paymentMethods, err := r.queries.GetPaymentMethods(ctx, GetPaymentMethodsParams{
		MerchantID: params.MerchantID,
		Limit:      int32(params.Take),
		Offset:     int32(params.Skip),
	})
	if err != nil {
		return nil, err
	}

	var paymentMethodResponses []entity.PaymentMethodResponse
	for _, pm := range paymentMethods {
		paymentMethodResponses = append(paymentMethodResponses, entity.PaymentMethodResponse{
			ID:         pm.ID,
			MerchantID: pm.MerchantID,
			Name:       pm.Name,
			Type:       pm.Type,
			Commission: pm.Commission,
			Details:    pm.Details.String,
			IsActive:   pm.IsActive,
			CreatedAt:  pm.CreatedAt,
			UpdatedAt:  pm.UpdatedAt,
			CreatedBy:  pm.CreatedBy,
			UpdatedBy:  pm.UpdatedBy,
		})
	}

	return &entity.PaymentMethodsResponse{
		Count:          len(paymentMethodResponses),
		PaymentMethods: paymentMethodResponses,
	}, nil
}

func (r *sqlcRepository) UpdatePaymentMethod(ctx context.Context, id uuid.UUID, req *entity.UpdatePaymentMethodRequest) error {
	paymentMethod, err := r.queries.GetPaymentMethod(ctx, id)
	if err != nil {
		return err
	}

	name := paymentMethod.Name
	if req.Name != nil {
		name = *req.Name
	}
	type_ := paymentMethod.Type
	if req.Type != nil {
		type_ = *req.Type
	}
	commission := paymentMethod.Commission
	if req.Commission != nil {
		commission = *req.Commission
	}
	details := paymentMethod.Details
	if req.Details != nil {
		details = sql.NullString{String: *req.Details, Valid: true}
	}
	isActive := paymentMethod.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	err = r.queries.UpdatePaymentMethod(ctx, UpdatePaymentMethodParams{
		ID:         id,
		Name:       name,
		Type:       type_,
		Commission: commission,
		Details:    details,
		IsActive:   isActive,
		UpdatedAt:  time.Now(),
		UpdatedBy:  paymentMethod.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) DeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeletePaymentMethod(ctx, id)
	return err
}

// Product operations
func (r *sqlcRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	err := r.queries.CreateProduct(ctx, CreateProductParams{
		ID:          product.ID,
		MerchantID:  product.MerchantID,
		Name:        product.Name,
		Description: sql.NullString{String: product.Description, Valid: product.Description != ""},
		Price:       product.Price,
		Currency:    sql.NullString{String: product.Currency, Valid: product.Currency != ""},
		Sku:         sql.NullString{String: product.SKU, Valid: product.SKU != ""},
		Weight:      sql.NullString{String: fmt.Sprintf("%f", product.Weight), Valid: true},
		Dimensions:  sql.NullString{String: product.Dimensions, Valid: product.Dimensions != ""},
		ImageUrl:    sql.NullString{String: product.ImageURL, Valid: product.ImageURL != ""},
		Status:      string(product.Status),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   product.CreatedBy,
		UpdatedBy:   product.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) GetProduct(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	product, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Product{
		ID:          product.ID,
		MerchantID:  product.MerchantID,
		Name:        product.Name,
		Description: product.Description.String,
		Price:       product.Price,
		Currency:    product.Currency,
		SKU:         product.Sku.String,
		Weight:      product.Weight.Float64,
		Dimensions:  product.Dimensions.String,
		ImageURL:    product.ImageUrl.String,
		Status:      entity.ProductStatus(product.Status),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		CreatedBy:   product.CreatedBy,
		UpdatedBy:   product.UpdatedBy,
	}, nil
}

func (r *sqlcRepository) GetProducts(ctx context.Context, params entity.GetProductsParams) (*entity.ProductsResponse, error) {
	products, err := r.queries.GetProducts(ctx, GetProductsParams{
		MerchantID: params.MerchantID,
		Limit:      int32(params.Take),
		Offset:     int32(params.Skip),
	})
	if err != nil {
		return nil, err
	}

	var productResponses []entity.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, entity.ProductResponse{
			ID:          p.ID,
			MerchantID:  p.MerchantID,
			Name:        p.Name,
			Description: p.Description.String,
			Price:       p.Price,
			Currency:    p.Currency,
			SKU:         p.Sku.String,
			Weight:      p.Weight.Float64,
			Dimensions:  p.Dimensions.String,
			ImageURL:    p.ImageUrl.String,
			Status:      entity.ProductStatus(p.Status),
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			CreatedBy:   p.CreatedBy,
			UpdatedBy:   p.UpdatedBy,
		})
	}

	return &entity.ProductsResponse{
		Count:    len(productResponses),
		Products: productResponses,
	}, nil
}

func (r *sqlcRepository) UpdateProduct(ctx context.Context, id uuid.UUID, req *entity.UpdateProductRequest) error {
	product, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		return err
	}

	name := product.Name
	if req.Name != nil {
		name = *req.Name
	}
	description := product.Description
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	price := product.Price
	if req.Price != nil {
		price = *req.Price // Assuming Price is float64 in entity.Product and req.Price is *float64
	}
	currency := product.Currency
	if req.Currency != nil { // Corrected: Use sql.NullString for Currency
		currency = *req.Currency
	}
	sku := product.Sku
	if req.SKU != nil {
		sku = sql.NullString{String: *req.SKU, Valid: true}
	}
	weight := product.Weight
	if req.Weight != nil {
		weight = sql.NullFloat64{Float64: *req.Weight, Valid: true}
	}
	dimensions := product.Dimensions
	if req.Dimensions != nil {
		dimensions = sql.NullString{String: *req.Dimensions, Valid: true}
	}
	imageURL := product.ImageUrl
	if req.ImageURL != nil {
		imageURL = sql.NullString{String: *req.ImageURL, Valid: true}
	}
	status := product.Status
	if req.Status != nil {
		status = *req.Status
	}

	err = r.queries.UpdateProduct(ctx, UpdateProductParams{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price, // Corrected: Price is float64, no change needed here
		Currency:    currency,
		Sku:         sku,
		Weight:      weight,
		Dimensions:  dimensions,
		ImageUrl:    imageURL,
		Status:      status,
		UpdatedAt:   time.Now(),
		UpdatedBy:   product.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteProduct(ctx, id)
	return err
}

// Warehouse operations
func (r *sqlcRepository) CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	err := r.queries.CreateWarehouse(ctx, CreateWarehouseParams{
		ID:          warehouse.ID,
		MerchantID:  warehouse.MerchantID,
		Name:        warehouse.Name,
		Location:    warehouse.Location,
		Capacity:    int32(warehouse.Capacity),
		IsActive:    warehouse.IsActive,
		Description: sql.NullString{String: warehouse.Description, Valid: warehouse.Description != ""},
		Status:      string(warehouse.Status),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   warehouse.CreatedBy,
		UpdatedBy:   warehouse.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) GetWarehouse(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error) {
	warehouse, err := r.queries.GetWarehouse(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Warehouse{
		ID:          warehouse.ID,
		MerchantID:  warehouse.MerchantID,
		Name:        warehouse.Name,
		Location:    warehouse.Location,
		Capacity:    int(warehouse.Capacity),
		IsActive:    warehouse.IsActive,
		Description: warehouse.Description.String,
		Status:      entity.WarehouseStatus(warehouse.Status),
		CreatedAt:   warehouse.CreatedAt,
		UpdatedAt:   warehouse.UpdatedAt,
		CreatedBy:   warehouse.CreatedBy,
		UpdatedBy:   warehouse.UpdatedBy,
	}, nil
}

func (r *sqlcRepository) GetWarehouses(ctx context.Context, params entity.GetWarehousesParams) (*entity.WarehousesResponse, error) {
	warehouses, err := r.queries.GetWarehouses(ctx, GetWarehousesParams{
		MerchantID: params.MerchantID,
		Limit:      int32(params.Take),
		Offset:     int32(params.Skip),
	})
	if err != nil {
		return nil, err
	}

	var warehouseResponses []entity.WarehouseResponse
	for _, w := range warehouses {
		warehouseResponses = append(warehouseResponses, entity.WarehouseResponse{
			ID:          w.ID,
			MerchantID:  w.MerchantID,
			Name:        w.Name,
			Location:    w.Location,
			Capacity:    int(w.Capacity),
			IsActive:    w.IsActive,
			Description: w.Description.String,
			Status:      entity.WarehouseStatus(w.Status),
			CreatedAt:   w.CreatedAt,
			UpdatedAt:   w.UpdatedAt,
			CreatedBy:   w.CreatedBy,
			UpdatedBy:   w.UpdatedBy,
		})
	}

	return &entity.WarehousesResponse{
		Count:      len(warehouseResponses),
		Warehouses: warehouseResponses,
	}, nil
}

func (r *sqlcRepository) UpdateWarehouse(ctx context.Context, id uuid.UUID, req *entity.UpdateWarehouseRequest) error {
	warehouse, err := r.queries.GetWarehouse(ctx, id)
	if err != nil {
		return err
	}

	name := warehouse.Name
	if req.Name != nil {
		name = *req.Name
	}
	location := warehouse.Location
	if req.Location != nil {
		location = *req.Location
	}
	capacity := warehouse.Capacity
	if req.Capacity != nil {
		capacity = int32(*req.Capacity)
	}
	isActive := warehouse.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	description := warehouse.Description
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	status := warehouse.Status
	if req.Status != nil {
		status = *req.Status
	}

	err = r.queries.UpdateWarehouse(ctx, UpdateWarehouseParams{
		ID:          id,
		Name:        name,
		Location:    location,
		Capacity:    capacity,
		IsActive:    isActive,
		Description: description,
		Status:      status,
		UpdatedAt:   time.Now(),
		UpdatedBy:   warehouse.UpdatedBy,
	})
	return err
}

func (r *sqlcRepository) DeleteWarehouse(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteWarehouse(ctx, id)
	return err
}
