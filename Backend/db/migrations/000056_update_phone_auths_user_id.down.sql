-- Remove index
DROP INDEX IF EXISTS idx_phone_auths_user_id;

-- Remove foreign key constraint
ALTER TABLE auth.phone_auths DROP CONSTRAINT IF EXISTS fk_phone_auths_user_id;

-- Remove user_id column
ALTER TABLE auth.phone_auths DROP COLUMN IF EXISTS user_id; 