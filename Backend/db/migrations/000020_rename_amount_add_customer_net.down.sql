-- Remove customer_net column
ALTER TABLE public.transactions
DROP COLUMN customer_net;

-- Rename base_amount back to amount
ALTER TABLE public.transactions 
RENAME COLUMN base_amount TO amount; 