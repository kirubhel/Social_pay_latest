-- name: GetDefaultCommission :one
SELECT value::text
FROM admin.settings
WHERE key = 'default_commission';

-- name: GetMerchantCommission :one
SELECT commission_active, commission_percent, commission_cent
FROM merchants.merchants
WHERE id = $1;

-- name: UpdateMerchantCommission :exec
UPDATE merchants.merchants
SET commission_active = $1,
    commission_percent = $2,
    commission_cent = $3,
    updated_at = NOW()
WHERE id = $4;

-- name: UpdateDefaultCommission :exec
UPDATE admin.settings
SET value = $1::jsonb,
    updated_at = NOW()
WHERE key = 'default_commission'; 