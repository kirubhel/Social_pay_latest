CREATE SCHEMA IF NOT EXISTS merchant;

CREATE TABLE IF NOT EXISTS merchant.wallet (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    amount FLOAT NOT NULL DEFAULT 0,
    locked_amount FLOAT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    wallet_type VARCHAR(10) NOT NULL DEFAULT 'merchant' CHECK (wallet_type IN ('merchant', 'admin')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES auth.users(id)
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_wallet_merchant_id ON merchant.wallet(merchant_id);
CREATE INDEX IF NOT EXISTS idx_wallet_user_id ON merchant.wallet(user_id); 