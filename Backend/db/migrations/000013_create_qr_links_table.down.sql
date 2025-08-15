-- Drop indexes first
DROP INDEX IF EXISTS idx_qr_links_created_at;
DROP INDEX IF EXISTS idx_qr_links_is_active;
DROP INDEX IF EXISTS idx_qr_links_tag;
DROP INDEX IF EXISTS idx_qr_links_type;
DROP INDEX IF EXISTS idx_qr_links_merchant_id;
DROP INDEX IF EXISTS idx_qr_links_user_id;

-- Drop foreign key constraints
ALTER TABLE public.qr_links DROP CONSTRAINT IF EXISTS fk_qr_links_merchant;
ALTER TABLE public.qr_links DROP CONSTRAINT IF EXISTS fk_qr_links_user;

-- Drop the qr_links table
DROP TABLE IF EXISTS public.qr_links;

-- Drop the enum types
DROP TYPE IF EXISTS qr_link_tag;
DROP TYPE IF EXISTS qr_link_type; 