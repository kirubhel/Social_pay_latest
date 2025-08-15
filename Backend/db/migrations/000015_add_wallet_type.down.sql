-- Drop index
DROP INDEX IF EXISTS idx_merchant_wallet_type;

-- Remove wallet_type column
ALTER TABLE merchant.wallet
DROP COLUMN wallet_type; 