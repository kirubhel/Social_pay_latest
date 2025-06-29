-- Drop the existing foreign key constraint if it exists
ALTER TABLE webhook.callback_logs DROP CONSTRAINT IF EXISTS fk_callback_logs_merchant;

-- Add the foreign key constraint to reference the merchants.merchants table
ALTER TABLE webhook.callback_logs
ADD CONSTRAINT fk_callback_logs_merchant
FOREIGN KEY (merchant_id) REFERENCES merchants.merchants(id);