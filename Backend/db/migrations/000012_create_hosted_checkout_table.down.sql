-- Drop indexes first
DROP INDEX IF EXISTS idx_hosted_payments_transaction_id;
DROP INDEX IF EXISTS idx_hosted_payments_reference;
DROP INDEX IF EXISTS idx_hosted_payments_expires_at;
DROP INDEX IF EXISTS idx_hosted_payments_created_at;
DROP INDEX IF EXISTS idx_hosted_payments_status;
DROP INDEX IF EXISTS idx_hosted_payments_merchant_id;
DROP INDEX IF EXISTS idx_hosted_payments_user_id;

-- Drop foreign key constraints
ALTER TABLE public.hosted_payments DROP CONSTRAINT IF EXISTS fk_hosted_payments_transaction;
ALTER TABLE public.hosted_payments DROP CONSTRAINT IF EXISTS fk_hosted_payments_merchant;
ALTER TABLE public.hosted_payments DROP CONSTRAINT IF EXISTS fk_hosted_payments_user;

-- Drop the hosted_payments table
DROP TABLE IF EXISTS public.hosted_payments;

-- Drop the hosted_payment_status enum type
DROP TYPE IF EXISTS hosted_payment_status;
