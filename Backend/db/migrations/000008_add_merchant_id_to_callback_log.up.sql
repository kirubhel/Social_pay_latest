-- Add merchant_id column to webhook.callback_logs if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_schema = 'webhook' 
        AND table_name = 'callback_logs' 
        AND column_name = 'merchant_id'
    ) THEN
        ALTER TABLE webhook.callback_logs
        ADD COLUMN merchant_id UUID,
        ADD CONSTRAINT fk_callback_logs_merchant 
            FOREIGN KEY (merchant_id) 
            REFERENCES merchants.merchants(id) 
            ON DELETE SET NULL;

        -- Create index for better query performance
        CREATE INDEX IF NOT EXISTS idx_callback_logs_merchant_id 
        ON webhook.callback_logs(merchant_id);
    END IF;
END $$; 