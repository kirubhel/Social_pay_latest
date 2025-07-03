-- Remove 2FA support from users table
ALTER TABLE auth.users DROP COLUMN IF EXISTS two_factor_enabled;
ALTER TABLE auth.users DROP COLUMN IF EXISTS two_factor_verified_at;

-- Drop 2FA verification codes table
DROP TABLE IF EXISTS auth.two_factor_codes;

-- Drop trigger and function
DROP TRIGGER IF EXISTS update_two_factor_codes_updated_at ON auth.two_factor_codes;
DROP FUNCTION IF EXISTS update_updated_at_column(); 