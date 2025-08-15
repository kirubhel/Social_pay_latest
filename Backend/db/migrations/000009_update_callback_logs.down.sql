-- Drop the foreign key constraint
ALTER TABLE webhook.callback_logs DROP CONSTRAINT IF EXISTS fk_callback_logs_merchant;