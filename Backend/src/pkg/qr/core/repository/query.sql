-- name: CreateQRLink :one
INSERT INTO public.qr_links (
    id, user_id, merchant_id, type, amount, supported_methods, 
    tag, title, description, image_url, is_tip_enabled
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetQRLink :one
SELECT * FROM public.qr_links 
WHERE id = $1 AND is_active = true;

-- name: GetQRLinksByMerchant :many
SELECT * FROM public.qr_links 
WHERE merchant_id = $1 AND is_active = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetQRLinksByUser :many
SELECT * FROM public.qr_links 
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC  
LIMIT $2 OFFSET $3;

-- name: UpdateQRLink :one
UPDATE public.qr_links 
SET 
    amount = COALESCE($2, amount),
    supported_methods = COALESCE($3, supported_methods),
    tag = COALESCE($4, tag),
    title = COALESCE($5, title),
    description = COALESCE($6, description),
    image_url = COALESCE($7, image_url),
    is_tip_enabled = COALESCE($8, is_tip_enabled),
    is_active = COALESCE($9, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $10
RETURNING *;

-- name: DeleteQRLink :exec
UPDATE public.qr_links 
SET is_active = false, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2;

-- name: CountQRLinksByMerchant :one
SELECT COUNT(*) FROM public.qr_links 
WHERE merchant_id = $1 AND is_active = true;

-- name: CountQRLinksByUser :one
SELECT COUNT(*) FROM public.qr_links 
WHERE user_id = $1 AND is_active = true; 