-- name: CreateMerchantWallet :exec
-- @param id uuid
-- @param user_id uuid
-- @param merchant_id uuid
-- @param amount numeric
-- @param locked_amount numeric
-- @param currency text
INSERT INTO merchant.wallet (
    id, user_id, merchant_id, amount, locked_amount, currency, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW(), NOW()
);

-- name: GetMerchantWallet :one
-- @param merchant_id uuid
SELECT * FROM merchant.wallet
WHERE merchant_id = $1;

-- name: GetMerchantWalletForUpdate :one
-- @param merchant_id uuid
SELECT * FROM merchant.wallet
WHERE merchant_id = $1
FOR UPDATE;

-- name: UpdateMerchantWallet :exec
-- @param id uuid
-- @param amount numeric
-- @param locked_amount numeric
UPDATE merchant.wallet
SET amount = $2,
    locked_amount = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateMerchantWalletWithTx :exec
-- @param id uuid
-- @param amount numeric
-- @param locked_amount numeric
UPDATE merchant.wallet
SET amount = $2,
    locked_amount = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateMerchantWalletLastSync :exec
-- @param id uuid
UPDATE merchant.wallet
SET last_sync_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND is_active = true;

-- name: DeactivateMerchantWallet :exec
-- @param id uuid
UPDATE merchant.wallet
SET is_active = false,
    updated_at = NOW()
WHERE id = $1;

-- name: GetMerchantWalletByUserID :one
-- @param user_id uuid
SELECT * FROM merchant.wallet
WHERE user_id = $1;

-- name: GetMerchantWalletByMerchantID :one
-- @param merchant_id uuid
SELECT * FROM merchant.wallet
WHERE merchant_id = $1;
