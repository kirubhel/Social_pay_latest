-- Add provider columns and MerchantPaysFee to transactions table for webhook data
ALTER TABLE public.transactions 
ADD COLUMN provider_tx_id VARCHAR(255),
ADD COLUMN provider_data JSONB,
ADD COLUMN merchant_pays_fee BOOLEAN DEFAULT FALSE;

-- Add index for provider transaction ID for better query performance
CREATE INDEX IF NOT EXISTS idx_transactions_provider_tx_id ON public.transactions(provider_tx_id);

-- Add comments for documentation
COMMENT ON COLUMN public.transactions.provider_tx_id IS 'Transaction ID from payment provider (e.g., MPesa transaction ID)';
COMMENT ON COLUMN public.transactions.provider_data IS 'Additional data from payment provider stored as JSON';
COMMENT ON COLUMN public.transactions.merchant_pays_fee IS 'Flag indicating whether merchant pays the transaction fee'; 