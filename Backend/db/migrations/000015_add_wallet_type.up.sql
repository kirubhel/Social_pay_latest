-- Add wallet_type column to merchant.wallet table (Safe version)
DO $$
BEGIN
    -- Check if wallet_type column doesn't exist before adding it
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_schema = 'merchant' 
        AND table_name = 'wallet' 
        AND column_name = 'wallet_type'
    ) THEN
        ALTER TABLE merchant.wallet ADD COLUMN wallet_type VARCHAR(20) NOT NULL DEFAULT 'merchant';
        ALTER TABLE merchant.wallet ADD CONSTRAINT chk_wallet_type CHECK (wallet_type IN ('merchant', 'admin', 'super_admin'));
        
        -- Create index on wallet_type for faster lookups
        CREATE INDEX IF NOT EXISTS idx_merchant_wallet_type ON merchant.wallet(wallet_type);
    END IF;
END $$; 