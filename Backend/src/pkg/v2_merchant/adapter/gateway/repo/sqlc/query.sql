-- Merchant queries for v2_merchant module

-- name: GetMerchant :one
SELECT * FROM merchants.merchants
WHERE id = $1;

-- name: SearchMerchants :many
WITH search_results AS (
    SELECT
        *,
        COUNT(*) OVER() AS total_count
    FROM
        merchants.merchants
    WHERE
        (
            LOWER(coalesce(legal_name, '')) LIKE '%' || LOWER($1) || '%'
            OR LOWER(coalesce(trading_name, '')) LIKE '%' || LOWER($1) || '%'
            OR LOWER(coalesce(business_registration_number, '')) LIKE '%' || LOWER($1) || '%'
            OR LOWER(coalesce(tax_identification_number, '')) LIKE '%' || LOWER($1) || '%'
            OR LOWER(coalesce(business_type, '')) LIKE '%' || LOWER($1) || '%'
            OR LOWER(coalesce(industry_category, '')) LIKE '%' || LOWER($1) || '%'
            OR $1 = ''
        )
        AND (
            ($4 != '0001-01-01 00:00:00+00'::timestamptz AND created_at >= $4::timestamptz)
            OR $4 = '0001-01-01 00:00:00+00'
        )
        AND (
            ($5 != '0001-01-01 00:00:00+00'::timestamptz AND created_at <= $5::timestamptz)
            OR $5 = '0001-01-01 00:00:00+00'
        )
        AND (
            -- Filter by status if provided
            (status = $6 OR $6 = '')
        ) 
        AND deleted_at IS NULL
    ORDER BY
        created_at DESC
    LIMIT $2
    OFFSET $3
)
SELECT
    id,
    user_id,
    legal_name,
    trading_name,
    business_registration_number,
    tax_identification_number,
    business_type,
    industry_category,
    is_betting_company,
    lottery_certificate_number,
    website_url,
    established_date,
    created_at,
    updated_at,
    status,
    total_count
FROM
    search_results; 

-- name: GetAllMerchants :many
SELECT * FROM merchants.merchants;

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

-- name: GetMerchantDocument :one
SELECT * FROM merchants.documents
WHERE id = $1;

-- name: GetMerchantBankAccounts :many
SELECT * FROM merchants.bank_accounts
WHERE merchant_id = $1
ORDER BY is_primary DESC, created_at ASC;

-- name: GetMerchantSettings :one
SELECT * FROM merchants.settings
WHERE merchant_id = $1;

-- name: UpdateMerchant :exec
UPDATE merchants.merchants
SET 
    legal_name=$2, 
    trading_name=$3, 
    business_registration_number=$4, 
    business_type=$5, 
    industry_category=$6, 
    is_betting_company=$7, 
    lottery_certificate_number=$8, 
    tax_identification_number=$9,
    website_url=$10, 
    established_date=$11, 
    status=$12
WHERE id = $1;

-- name: UpdateMerchantStatus :exec
UPDATE merchants.merchants SET status = $2 WHERE id = $1;

-- name: CreateMerchantContact :exec
INSERT INTO merchants.contacts(
    id, 
    merchant_id, 
    contact_type, 
    first_name, 
    last_name, 
    email, 
    phone_number,
    created_at, 
    updated_at
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: UpdateMerchantContact :exec
UPDATE merchants.contacts
SET 
    first_name=$2, 
    last_name=$3, 
    phone_number=$4, 
    email=$5, 
    is_verified=$6
WHERE id=$1;

-- name: CreateMerchantDocument :exec
INSERT INTO merchants.documents(
    id, 
    merchant_id, 
    document_type, 
    file_url, 
    status, 
    created_at, 
    updated_at
)
VALUES($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateMerchantDocument :exec
UPDATE merchants.documents
SET
    file_url=$2,
    verified_by=$3, 
    verified_at=$4, 
    status=$5, 
    rejection_reason=$6
WHERE id=$1;

-- name: UpdateMerchantDocumentWithType :exec
UPDATE merchants.documents
SET
    document_type=$2,
    file_url=$3,
    status=$4,
    verified_by=$5, 
    verified_at=$6, 
    rejection_reason=$7
WHERE id=$1;

-- name: DeleteMerchant :exec
UPDATE merchants.merchants SET deleted_at = $2 WHERE id = $1;

-- name: DeleteMerchants :exec

-- name: BatchSoftDeleteMerchants :exec
UPDATE merchants.merchants
SET deleted_at = NOW()
WHERE id = ANY($1::uuid[]);

-- name: GetMerchantStats :one
SELECT 
    COUNT(*) as total_merchants,
    COUNT(CASE WHEN status = 'active' THEN 1 END) as active_merchants,
    COUNT(CASE WHEN status = 'pending_verification' THEN 1 END) as pending_kyc,
    COUNT(CASE WHEN created_at >= date_trunc('month', CURRENT_DATE) THEN 1 END) as new_this_month
FROM merchants.merchants
WHERE deleted_at IS NULL;

