-- Add merchant_pays_fee column to hosted_payments table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND table_name = 'hosted_payments' 
        AND column_name = 'merchant_pays_fee'
    ) THEN
        ALTER TABLE public.hosted_payments 
        ADD COLUMN merchant_pays_fee BOOLEAN NOT NULL DEFAULT false;
    END IF;
END $$;

-- Add comment for documentation (will update even if column exists)
COMMENT ON COLUMN public.hosted_payments.merchant_pays_fee IS 'Indicates whether the merchant pays the transaction fee instead of the customer'; 