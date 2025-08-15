-- Add accept_tip column to hosted_payments table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND table_name = 'hosted_payments' 
        AND column_name = 'accept_tip'
    ) THEN
        ALTER TABLE public.hosted_payments 
        ADD COLUMN accept_tip BOOLEAN NOT NULL DEFAULT false;
    END IF;
END $$;

-- Add comment for documentation (will update even if column exists)
COMMENT ON COLUMN public.hosted_payments.accept_tip IS 'Indicates whether the hosted payment accepts tips from customers';
