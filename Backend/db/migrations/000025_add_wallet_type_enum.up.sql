-- Migration to add wallet_type column to merchant.wallet table
-- This migration ensures the column exists before any data migration operations

-- First check if the column already exists, if not add it
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_schema = 'merchant' 
        AND table_name = 'wallet' 
        AND column_name = 'wallet_type'
    ) THEN
        -- Add wallet_type column with default value
        ALTER TABLE merchant.wallet 
        ADD COLUMN wallet_type VARCHAR(10) NOT NULL DEFAULT 'merchant';
        
        -- Add check constraint to enforce enum-like behavior
        ALTER TABLE merchant.wallet 
        ADD CONSTRAINT chk_wallet_type CHECK (wallet_type IN ('merchant', 'admin'));
        
        -- Create index for better performance
        CREATE INDEX IF NOT EXISTS idx_merchant_wallet_type ON merchant.wallet(wallet_type);
    END IF;
END $$; 