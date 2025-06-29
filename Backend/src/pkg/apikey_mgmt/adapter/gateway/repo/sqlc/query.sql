-- Common columns for reference:
-- id, user_id, created_by, name, description, public_key, secret_key,
-- can_withdrawal, can_process_payment, created_at, updated_at, 
-- expires_at, last_used_at, is_active, merchant_id

-- name: CreateAPIKey :one
INSERT INTO api_keys (
    id, user_id, created_by, name, description, public_key, secret_key,
    can_withdrawal, can_process_payment, created_at, updated_at, expires_at, 
    is_active, merchant_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING *;

-- name: GetAPIKey :one
SELECT * FROM api_keys
WHERE id = $1;

-- name: GetAPIKeyByPublicKey :one
SELECT * FROM api_keys
WHERE public_key = $1;

-- name: GetAPIKeyBySecretKey :one
SELECT * FROM api_keys
WHERE secret_key = $1;

-- name: ListAPIKeys :many
SELECT * FROM api_keys
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListMerchantAPIKeys :many
SELECT * FROM api_keys
WHERE merchant_id = $1
ORDER BY created_at DESC;

-- name: UpdateAPIKey :one
UPDATE api_keys
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    can_withdrawal = COALESCE($4, can_withdrawal),
    can_process_payment = COALESCE($5, can_process_payment),
    is_active = COALESCE($6, is_active),
    expires_at = COALESCE($7, expires_at),
    merchant_id = COALESCE($8, merchant_id),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys
WHERE id = $1;

-- name: RotateAPIKeySecret :one
UPDATE api_keys
SET secret_key = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateLastUsedAt :one
UPDATE api_keys
SET last_used_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ValidateAPIKey :one
SELECT * FROM api_keys
WHERE public_key = $1 AND secret_key = $2 AND is_active = true; 