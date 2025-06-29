-- Merchant queries for v2_merchant module

-- name: GetMerchant :one
SELECT * FROM merchants.merchants
WHERE id = $1;

-- name: GetMerchantByUserID :one
SELECT * FROM merchants.merchants
WHERE user_id = $1;

-- name: GetMerchantAddresses :many
SELECT * FROM merchants.addresses
WHERE merchant_id = $1
ORDER BY is_primary DESC, created_at ASC;

-- name: GetMerchantContacts :many
SELECT * FROM merchants.contacts
WHERE merchant_id = $1
ORDER BY contact_type, created_at ASC;

-- name: GetMerchantDocuments :many
SELECT * FROM merchants.documents
WHERE merchant_id = $1
ORDER BY document_type, created_at DESC;

-- name: GetMerchantBankAccounts :many
SELECT * FROM merchants.bank_accounts
WHERE merchant_id = $1
ORDER BY is_primary DESC, created_at ASC;

-- name: GetMerchantSettings :one
SELECT * FROM merchants.settings
WHERE merchant_id = $1; 