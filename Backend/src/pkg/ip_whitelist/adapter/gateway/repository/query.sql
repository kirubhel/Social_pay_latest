-- name: CreateWhitelistedIP :exec
INSERT INTO ip_whitelist.whitelisted_ips (
    id, merchant_id, ip_address, description, is_active, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW()
);

-- name: GetWhitelistedIP :one
SELECT EXISTS (
    SELECT 1 FROM ip_whitelist.whitelisted_ips
    WHERE merchant_id = $1 AND ip_address = $2
) as ip_whitelisted;

-- name: GetWhitelistedIPByID :one
SELECT * FROM ip_whitelist.whitelisted_ips
WHERE id = $1;

-- name: GetWhitelistedIPsByMerchantID :many
SELECT * FROM ip_whitelist.whitelisted_ips
WHERE merchant_id = $1
ORDER BY created_at DESC;

-- name: UpdateWhitelistedIP :exec
UPDATE ip_whitelist.whitelisted_ips
SET ip_address = $2,
    is_active = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteWhitelistedIP :exec
DELETE FROM ip_whitelist.whitelisted_ips
WHERE id = $1;

-- name: CheckIPWhitelisted :one
SELECT EXISTS (
    SELECT 1 FROM ip_whitelist.whitelisted_ips
    WHERE merchant_id = $1
    AND is_active = true
    AND ip_address >>= $2::inet
) as is_whitelisted; 
