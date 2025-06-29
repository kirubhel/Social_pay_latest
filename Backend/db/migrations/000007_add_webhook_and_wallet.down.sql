-- Drop indexes for webhook.callback_logs
DROP INDEX IF EXISTS webhook.idx_callback_logs_user_id;
DROP INDEX IF EXISTS webhook.idx_callback_logs_txn_id;
DROP INDEX IF EXISTS webhook.idx_callback_logs_status;
DROP INDEX IF EXISTS webhook.idx_callback_logs_created_at;

-- Drop webhook.callback_logs table
DROP TABLE IF EXISTS webhook.callback_logs;

-- Drop indexes for merchant.wallet
DROP INDEX IF EXISTS merchant.idx_wallet_user_id;
DROP INDEX IF EXISTS merchant.idx_wallet_merchant_id;
DROP INDEX IF EXISTS merchant.idx_wallet_currency;
DROP INDEX IF EXISTS merchant.idx_wallet_created_at;

-- Drop merchant.wallet table
DROP TABLE IF EXISTS merchant.wallet;

-- Drop schemas
DROP SCHEMA IF EXISTS webhook CASCADE;
DROP SCHEMA IF EXISTS merchant CASCADE;