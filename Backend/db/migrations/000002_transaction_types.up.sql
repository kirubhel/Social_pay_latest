-- create transaction type  
CREATE TYPE public.transaction_type AS ENUM (
    'deposit',
    'withdrawal',
    'transfer',
    'payment',
    'refund'
);
