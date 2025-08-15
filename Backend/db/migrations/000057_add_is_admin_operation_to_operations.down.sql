-- Remove email column from auth.users table
ALTER TABLE auth.users DROP COLUMN IF EXISTS email;

-- Remove is_admin_operation column from auth.operations table  
ALTER TABLE auth.operations DROP COLUMN IF EXISTS is_admin_operation;

-- Drop the email index
DROP INDEX IF EXISTS idx_users_email; 