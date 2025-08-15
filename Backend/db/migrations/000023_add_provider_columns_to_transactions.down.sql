-- Remove provider columns and MerchantPaysFee from transactions table
DROP INDEX IF EXISTS idx_transactions_provider_tx_id;
ALTER TABLE public.transactions 
DROP COLUMN IF EXISTS provider_data,
DROP COLUMN IF EXISTS provider_tx_id,
DROP COLUMN IF EXISTS merchant_pays_fee; 