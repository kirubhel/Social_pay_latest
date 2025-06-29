CREATE SCHEMA IF NOT EXISTS webhook;

CREATE SCHEMA IF NOT EXISTS merchant;
CREATE TABLE IF NOT EXISTS webhook.callback_logs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    txn_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    status INTEGER NOT NULL,
    request_body TEXT NOT NULL,
    response_body TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (txn_id) REFERENCES public.transactions(id)
); 

CREATE TABLE IF NOT EXISTS merchant.wallet (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    amount FLOAT NOT NULL DEFAULT 0,
    locked_amount FLOAT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES auth.users(id)
);
