-- name: CreateCallbackLog :exec
INSERT INTO webhook.callback_logs (
    id, user_id, txn_id, merchant_id, status, request_body, response_body, retry_count, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
);

-- name: GetCallbackLogByID :one
SELECT * FROM webhook.callback_logs
WHERE id = $1;

-- name: UpdateCallbackLog :exec
UPDATE webhook.callback_logs
SET status = $2,
    response_body = $3,
    retry_count = $4,
    updated_at = NOW()
WHERE id = $1;

-- name: GetCallbackLogsByTransactionID :many
SELECT * FROM webhook.callback_logs
WHERE txn_id = $1
ORDER BY created_at DESC;

-- name: GetCallbackLogsByStatus :many
SELECT * FROM webhook.callback_logs
WHERE status = $1
ORDER BY created_at DESC;

-- name: GetCallbackLogsByMerchantID :many
SELECT * FROM webhook.callback_logs
WHERE merchant_id = $1
ORDER BY created_at DESC;

-- name: GetAllCallbackLogs :many
SELECT * FROM webhook.callback_logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateMerchantWallet :exec
INSERT INTO merchant.wallet (
    id, user_id, merchant_id, amount, locked_amount, currency, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW(), NOW()
);

-- name: GetMerchantWalletByUserID :one
SELECT * FROM merchant.wallet
WHERE user_id = $1;

-- name: UpdateMerchantWallet :exec
UPDATE merchant.wallet
SET amount = $2,
    locked_amount = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: GetMerchantWalletByMerchantID :one
SELECT * FROM merchant.wallet
WHERE merchant_id = $1;

-- name: GetMerchantWalletByMerchantIDForUpdate :one
SELECT * FROM merchant.wallet
WHERE merchant_id = $1
FOR UPDATE;

-- name: UpdateMerchantWalletAmountByMerchantID :exec
UPDATE merchant.wallet
SET amount = $2,
    updated_at = NOW()
WHERE merchant_id = $1;


