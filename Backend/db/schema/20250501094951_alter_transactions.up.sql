ALTER TABLE public.transactions
ALTER COLUMN status TYPE transaction_status
USING status::transaction_status;