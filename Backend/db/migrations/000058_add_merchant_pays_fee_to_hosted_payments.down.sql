-- Remove merchant_pays_fee column from hosted_payments table
ALTER TABLE public.hosted_payments 
DROP COLUMN IF EXISTS merchant_pays_fee; 