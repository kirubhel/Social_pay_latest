-- name: CreateMerchantWallet :exec
INSERT INTO merchant.wallet (
    id, user_id, merchant_id, amount, locked_amount, currency, wallet_type, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
);

-- name: GetMerchantWalletByUserID :one
SELECT * FROM merchant.wallet
WHERE user_id = $1 AND wallet_type = 'merchant';

-- name: UpdateMerchantWallet :exec
UPDATE merchant.wallet
SET amount = $2,
    locked_amount = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: GetMerchantWalletByMerchantID :one
SELECT * FROM merchant.wallet
WHERE merchant_id = $1 AND wallet_type = 'merchant';

-- name: GetMerchantWalletByMerchantIDForUpdate :one
SELECT * FROM merchant.wallet
WHERE merchant_id = $1 AND wallet_type = 'merchant'
FOR UPDATE;

-- name: UpdateMerchantWalletAmountByMerchantID :exec
UPDATE merchant.wallet
SET amount = $2,
    updated_at = NOW()
WHERE merchant_id = $1 AND wallet_type = 'merchant';

-- name: GetAdminWallet :one
SELECT * FROM merchant.wallet
WHERE  wallet_type = 'super_admin';

-- name: GetAdminWalletForUpdate :one
SELECT * FROM merchant.wallet
WHERE wallet_type = 'super_admin'
FOR UPDATE;

-- name: GetTotalAdminWalletAmount :one
SELECT 
    CAST(SUM(amount) AS FLOAT) as total_amount,
    CAST(SUM(locked_amount) AS FLOAT) as total_locked_amount,
    currency
FROM merchant.wallet
WHERE wallet_type = 'super_admin'
GROUP BY currency;