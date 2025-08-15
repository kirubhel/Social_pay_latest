-- Revert vat_amount column precision in transactions table back to 2 decimal places
ALTER TABLE public.transactions
ALTER COLUMN vat_amount TYPE numeric(20,2);

-- Update the comment to reflect the change
COMMENT ON COLUMN public.transactions.vat_amount IS 'VAT amount with 2 decimal places precision';
