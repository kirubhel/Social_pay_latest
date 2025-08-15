-- Update vat_amount column precision in transactions table to 4 decimal places
ALTER TABLE public.transactions
ALTER COLUMN vat_amount TYPE numeric(20,4);

-- Add a comment to document the change
COMMENT ON COLUMN public.transactions.vat_amount IS 'VAT amount with 4 decimal places precision';
