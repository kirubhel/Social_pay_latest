-- Migration to populate customer_net for existing transactions
-- This script calculates customer_net based on the existing transaction data

DO $$
DECLARE
    affected_rows INTEGER;
BEGIN
    -- Update transactions where customer_net is NULL
    -- For deposit transactions: customer_net = total_amount (or fallback to base_amount)
    -- For withdrawal transactions: customer_net = base_amount
    UPDATE public.transactions 
    SET customer_net = CASE 
        WHEN type = 'withdrawal' THEN base_amount
        ELSE COALESCE(total_amount, base_amount)
    END
    WHERE customer_net IS NULL;
    
    GET DIAGNOSTICS affected_rows = ROW_COUNT;
    RAISE NOTICE 'Updated customer_net for % existing transactions', affected_rows;
END $$; 