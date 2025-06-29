-- Drop indexes first
DROP INDEX IF EXISTS webhook.idx_callback_logs_user_id;
DROP INDEX IF EXISTS webhook.idx_callback_logs_txn_id;
DROP INDEX IF EXISTS webhook.idx_callback_logs_merchant_id;
DROP INDEX IF EXISTS webhook.idx_callback_logs_status;
DROP INDEX IF EXISTS webhook.idx_callback_logs_created_at;

-- Drop the table
DROP TABLE IF EXISTS webhook.callback_logs;

-- Drop index first
DROP INDEX IF EXISTS webhook.idx_callback_logs_merchant_id;

-- Drop foreign key constraint
ALTER TABLE webhook.callback_logs
DROP CONSTRAINT IF EXISTS fk_callback_logs_merchant;

-- Drop merchant_id column
ALTER TABLE webhook.callback_logs
DROP COLUMN IF EXISTS merchant_id; 