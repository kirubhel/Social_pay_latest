-- Rollback: Remove phone_prefix and phone_number columns from auth.users table

-- Drop check constraints
ALTER TABLE auth.users DROP CONSTRAINT IF EXISTS chk_users_phone_prefix_format;
ALTER TABLE auth.users DROP CONSTRAINT IF EXISTS chk_users_phone_number_format;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_phone_prefix_number;
DROP INDEX IF EXISTS idx_users_phone_prefix;
DROP INDEX IF EXISTS idx_users_phone_number;

-- Drop unique constraint
ALTER TABLE auth.users DROP CONSTRAINT IF EXISTS unique_users_phone_prefix_number;

-- Drop columns
ALTER TABLE auth.users DROP COLUMN IF EXISTS phone_prefix;
ALTER TABLE auth.users DROP COLUMN IF EXISTS phone_number; 