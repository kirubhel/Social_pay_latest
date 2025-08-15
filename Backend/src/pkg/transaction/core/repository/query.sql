-- Common columns for reference:
-- id, phone_number, user_id, merchant_id, type, medium, reference, comment, 
-- reference_number, description, verified, status, test, has_challenge, ttl, 
-- created_at, updated_at, confirm_timestamp, base_amount, fee_amount, admin_net,
-- vat_amount, merchant_net, total_amount, currency, details, token, 
-- callback_url, success_url, failed_url

-- name: CreateTransaction :exec
INSERT INTO public.transactions (
    id, phone_number, user_id, merchant_id, type, medium, reference, comment, verified,
    ttl, details, confirm_timestamp, reference_number, test, status,
    description, token, base_amount, has_challenge, fee_amount, admin_net,
    vat_amount, merchant_net, total_amount, customer_net, currency, callback_url,
    success_url, failed_url, transaction_source, qr_link_id, hosted_checkout_id, qr_tag,
    has_tip, tip_amount, tipee_phone, tip_medium, merchant_pays_fee
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
    $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
    $29, $30, $31, $32, $33, $34, $35, $36, $37, $38
);

-- name: CreateTransactionWithContext :exec
INSERT INTO public.transactions (
    id, phone_number, user_id, merchant_id, type, medium, reference, comment, verified,
    ttl, details, confirm_timestamp, reference_number, test, status,
    description, token, base_amount, has_challenge, fee_amount, admin_net,
    vat_amount, merchant_net, total_amount, currency, callback_url,
    success_url, failed_url, transaction_source, qr_link_id, hosted_checkout_id, qr_tag,
    has_tip, tip_amount, tipee_phone, tip_medium, merchant_pays_fee
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
    $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
    $29, $30, $31, $32, $33, $34, $35, $36, $37
);

-- name: UpdateTransaction :exec
UPDATE public.transactions 
SET 
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetTransaction :one
SELECT * FROM public.transactions 
WHERE id = $1;

-- name: GetTransactionWithMerchant :one
SELECT 
    t.*,
    m.id as merchant_id,
    m.legal_name as merchant_legal_name,
    m.trading_name as merchant_trading_name,
    m.business_registration_number as merchant_business_registration_number,
    m.tax_identification_number as merchant_tax_identification_number,
    m.business_type as merchant_business_type,
    m.industry_category as merchant_industry_category,
    m.is_betting_company as merchant_is_betting_company,
    m.lottery_certificate_number as merchant_lottery_certificate_number,
    m.website_url as merchant_website_url,
    m.established_date as merchant_established_date,
    m.status as merchant_status,
    m.created_at as merchant_created_at,
    m.updated_at as merchant_updated_at
FROM public.transactions t
LEFT JOIN merchants.merchants m ON t.merchant_id = m.id
WHERE t.id = $1;

-- Hosted Payments Queries

-- name: CreateHostedPayment :one
INSERT INTO public.hosted_payments (
    id, user_id, merchant_id, amount, currency, description, reference, 
    supported_mediums, phone_number, success_url, failed_url, callback_url, merchant_pays_fee, accept_tip
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: GetHostedPayment :one
SELECT * FROM public.hosted_payments 
WHERE id = $1;

-- name: GetHostedPaymentByReference :one
SELECT * FROM public.hosted_payments 
WHERE reference = $1 AND merchant_id = $2;

-- name: UpdateHostedPaymentWithTransaction :exec
UPDATE public.hosted_payments 
SET 
    transaction_id = $2,
    selected_medium = $3,
    selected_phone_number = $4,
    status = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateHostedPaymentStatus :exec
UPDATE public.hosted_payments 
SET 
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetExpiredHostedPayments :many
SELECT * FROM public.hosted_payments 
WHERE status = 'PENDING' AND expires_at < CURRENT_TIMESTAMP;

-- name: GetTransactions :many
SELECT * FROM public.transactions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetMerchantTransactions :many
SELECT * FROM public.transactions
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetFilteredTransactions :many
SELECT * FROM public.transactions
WHERE user_id = $1
    AND created_at BETWEEN $2 AND $3
    AND (status = $4)
    AND (type = $5)
ORDER BY created_at DESC
LIMIT $6 OFFSET $7;

-- name: GetFilteredMerchantTransactions :many
SELECT * FROM public.transactions
WHERE merchant_id = $1
    AND created_at BETWEEN $2 AND $3
    AND (status = $4)
    AND (type = $5)
ORDER BY created_at DESC
LIMIT $6 OFFSET $7;


-- name: GetTransactionsByStatus :many
SELECT * FROM public.transactions
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTransactionsByType :many
SELECT * FROM public.transactions
WHERE type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetByReferenceID :one
SELECT * FROM public.transactions
WHERE reference = $1
LIMIT 1;

-- name: GetByUserIdAndReferenceID :one
SELECT * FROM public.transactions
WHERE user_id = $1 AND reference = $2
LIMIT 1;

-- name: GetByMerchantIdAndReferenceID :one
SELECT * FROM public.transactions
WHERE merchant_id = $1 AND reference = $2
LIMIT 1;

-- name: UpdateStatus :exec
UPDATE public.transactions
SET status = $2, webhook_received = TRUE, updated_at = NOW()
WHERE id = $1;

-- name: OverrideTransactionStatus :exec
WITH updated_transaction AS (
    UPDATE public.transactions
    SET status = $2, updated_at = NOW()
    WHERE id = $1
    RETURNING id
)
INSERT INTO public.transaction_status_overrides (
    id,
    transaction_id,
    reason,
    admin_id
)
VALUES (
    gen_random_uuid(),
    $1,
    $3,
    $4
);

-- name: UpdateTipProcessing :exec
UPDATE public.transactions 
SET 
    tip_transaction_id = $2,
    tip_processed = true,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetTransactionsWithPendingTips :many
SELECT * FROM public.transactions 
WHERE has_tip = true AND tip_processed = false AND status = 'SUCCESS'
ORDER BY created_at ASC;

-- name: GetTransactionsByQRLink :many
SELECT * FROM public.transactions 
WHERE qr_link_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: CountTransactions :one
SELECT COUNT(*) FROM public.transactions
WHERE   user_id=$1;

-- name: UpdateHostedPayment :exec
UPDATE public.hosted_payments 
SET 
    amount = $2,
    currency = $3,
    description = $4,
    supported_mediums = $5,
    phone_number = $6,
    success_url = $7,
    failed_url = $8,
    callback_url = $9,
    expires_at = $10,
    merchant_pays_fee = $11,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
