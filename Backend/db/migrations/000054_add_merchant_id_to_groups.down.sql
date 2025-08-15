-- Remove indexes
DROP INDEX IF EXISTS idx_groups_merchant_id;
DROP INDEX IF EXISTS idx_user_groups_merchant_id;

-- Remove merchant_id columns
ALTER TABLE auth.user_groups DROP COLUMN IF EXISTS merchant_id;
ALTER TABLE auth.groups DROP COLUMN IF EXISTS merchant_id; 