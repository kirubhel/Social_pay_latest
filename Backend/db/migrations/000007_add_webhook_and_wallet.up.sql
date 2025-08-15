-- Create schemas
CREATE SCHEMA IF NOT EXISTS webhook;
CREATE SCHEMA IF NOT EXISTS merchant;

-- Create webhook.callback_logs table
CREATE TABLE IF NOT EXISTS webhook.callback_logs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    txn_id UUID NOT NULL,
    status INTEGER NOT NULL,
    request_body TEXT NOT NULL,
    response_body TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (txn_id) REFERENCES public.transactions(id)
);

-- Add indexes to optimize queries on webhook.callback_logs
CREATE INDEX IF NOT EXISTS idx_callback_logs_user_id ON webhook.callback_logs (user_id);
CREATE INDEX IF NOT EXISTS idx_callback_logs_txn_id ON webhook.callback_logs (txn_id);
CREATE INDEX IF NOT EXISTS idx_callback_logs_status ON webhook.callback_logs (status);
CREATE INDEX IF NOT EXISTS idx_callback_logs_created_at ON webhook.callback_logs (created_at);

-- Create merchant.wallet table
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

-- Add indexes to optimize queries on merchant.wallet
CREATE INDEX IF NOT EXISTS idx_wallet_user_id ON merchant.wallet (user_id);
CREATE INDEX IF NOT EXISTS idx_wallet_merchant_id ON merchant.wallet (merchant_id);
CREATE INDEX IF NOT EXISTS idx_wallet_currency ON merchant.wallet (currency);
CREATE INDEX IF NOT EXISTS idx_wallet_created_at ON merchant.wallet (created_at);