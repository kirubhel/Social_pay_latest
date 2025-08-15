DROP INDEX IF EXISTS idx_api_keys_merchant_id;

ALTER TABLE api_keys
DROP CONSTRAINT IF EXISTS fk_merchant,
DROP COLUMN IF EXISTS merchant_id;
