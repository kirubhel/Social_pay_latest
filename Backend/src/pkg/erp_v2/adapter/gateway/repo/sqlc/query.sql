-- name: CreateCatalog :exec
INSERT INTO erp.catalogs (id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetCatalog :one
SELECT id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by
FROM erp.catalogs
WHERE id = $1;

-- name: GetCatalogs :many
SELECT id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by
FROM erp.catalogs
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateCatalog :exec
UPDATE erp.catalogs
SET name = $2, description = $3, status = $4, updated_at = $5, updated_by = $6
WHERE id = $1;

-- name: DeleteCatalog :exec
DELETE FROM erp.catalogs
WHERE id = $1;

-- name: CreateCustomer :exec
INSERT INTO erp.customers (id, customer_id, merchant_id, name, email, phone, address, loyalty_points, date_of_birth, status, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);

-- name: GetCustomer :one
SELECT id, customer_id, merchant_id, name, email, phone, address, loyalty_points, date_of_birth, status, created_at, updated_at, created_by, updated_by
FROM erp.customers
WHERE id = $1;

-- name: GetCustomers :many
SELECT id, customer_id, merchant_id, name, email, phone, address, loyalty_points, date_of_birth, status, created_at, updated_at, created_by, updated_by
FROM erp.customers
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateCustomer :exec
UPDATE erp.customers
SET name = $2, email = $3, phone = $4, address = $5, loyalty_points = $6, date_of_birth = $7, status = $8, updated_at = $9, updated_by = $10
WHERE id = $1;

-- name: DeleteCustomer :exec
DELETE FROM erp.customers
WHERE id = $1;

-- name: CreateProduct :exec
INSERT INTO erp.products (id, merchant_id, name, description, price, currency, sku, weight, dimensions, image_url, status, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);

-- name: GetProduct :one
SELECT id, merchant_id, name, description, price, currency, sku, weight, dimensions, image_url, status, created_at, updated_at, created_by, updated_by
FROM erp.products
WHERE id = $1;

-- name: GetProducts :many
SELECT id, merchant_id, name, description, price, currency, sku, weight, dimensions, image_url, status, created_at, updated_at, created_by, updated_by
FROM erp.products
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProduct :exec
UPDATE erp.products
SET name = $2, description = $3, price = $4, currency = $5, sku = $6, weight = $7, dimensions = $8, image_url = $9, status = $10, updated_at = $11, updated_by = $12
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM erp.products
WHERE id = $1;

-- name: CreateWarehouse :exec
INSERT INTO erp.warehouses (id, merchant_id, name, location, capacity, is_active, description, status, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: GetWarehouse :one
SELECT id, merchant_id, name, location, capacity, is_active, description, status, created_at, updated_at, created_by, updated_by
FROM erp.warehouses
WHERE id = $1;

-- name: GetWarehouses :many
SELECT id, merchant_id, name, location, capacity, is_active, description, status, created_at, updated_at, created_by, updated_by
FROM erp.warehouses
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateWarehouse :exec
UPDATE erp.warehouses
SET name = $2, location = $3, capacity = $4, is_active = $5, description = $6, status = $7, updated_at = $8, updated_by = $9
WHERE id = $1;

-- name: DeleteWarehouse :exec
DELETE FROM erp.warehouses
WHERE id = $1;

-- name: CreatePaymentMethod :exec
INSERT INTO erp.payment_methods (id, merchant_id, name, type, commission, details, is_active, created_at, updated_at, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetPaymentMethod :one
SELECT id, merchant_id, name, type, commission, details, is_active, created_at, updated_at, created_by, updated_by
FROM erp.payment_methods
WHERE id = $1;

-- name: GetPaymentMethods :many
SELECT id, merchant_id, name, type, commission, details, is_active, created_at, updated_at, created_by, updated_by
FROM erp.payment_methods
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePaymentMethod :exec
UPDATE erp.payment_methods
SET name = $2, type = $3, commission = $4, details = $5, is_active = $6, updated_at = $7, updated_by = $8
WHERE id = $1;

-- name: DeletePaymentMethod :exec
DELETE FROM erp.payment_methods
WHERE id = $1;

-- name: CreateOrder :exec
INSERT INTO erp.orders (id, merchant_id, customer_details, order_details, order_items, metadata, tracking, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetOrder :one
SELECT id, merchant_id, customer_details, order_details, order_items, metadata, tracking, created_at, updated_at
FROM erp.orders
WHERE id = $1;

-- name: GetOrders :many
SELECT id, merchant_id, customer_details, order_details, order_items, metadata, tracking, created_at, updated_at
FROM erp.orders
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrder :exec
UPDATE erp.orders
SET customer_details = $2, order_details = $3, order_items = $4, metadata = $5, tracking = $6, updated_at = $7
WHERE id = $1;

-- name: DeleteOrder :exec
DELETE FROM erp.orders
WHERE id = $1;

-- name: CreateAccount :exec
INSERT INTO erp.accounts (id, user_id, title, type, default_account, verification_status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetAccount :one
SELECT id, user_id, title, type, default_account, verification_status, created_at, updated_at
FROM erp.accounts
WHERE id = $1;

-- name: UpdateAccount :exec
UPDATE erp.accounts
SET title = $2, type = $3, default_account = $4, verification_status = $5, updated_at = $6
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM erp.accounts
WHERE id = $1;

