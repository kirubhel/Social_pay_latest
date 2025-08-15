-- Remove accept_tip column from hosted_payments table
ALTER TABLE public.hosted_payments 
DROP COLUMN IF EXISTS accept_tip;
