BEGIN;

-- Drop triggers first
DROP TRIGGER IF EXISTS update_api_keys_updated_at ON api_keys;

-- Drop indexes
DROP INDEX IF EXISTS idx_api_keys_user_id;
DROP INDEX IF EXISTS idx_api_keys_public_key;
DROP INDEX IF EXISTS idx_api_keys_created_at;
DROP INDEX IF EXISTS idx_api_keys_is_active;

-- Drop the table (this will automatically drop its constraints)
DROP TABLE IF EXISTS api_keys;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

COMMIT; 