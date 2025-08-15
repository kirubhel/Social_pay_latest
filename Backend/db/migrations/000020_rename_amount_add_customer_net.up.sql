-- Rename amount column to base_amount
ALTER TABLE public.transactions 
RENAME COLUMN amount TO base_amount;

-- Add customer_net column
ALTER TABLE public.transactions
ADD COLUMN customer_net DECIMAL(20,2);

-- Add comments to document the columns
COMMENT ON COLUMN public.transactions.base_amount IS 'The original amount in the request';
COMMENT ON COLUMN public.transactions.customer_net IS 'Amount deducted or added from customer account'; 