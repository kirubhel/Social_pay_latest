-- Drop admin analytics indexes
DROP INDEX IF EXISTS idx_transactions_admin_analytics_optimized;
DROP INDEX IF EXISTS idx_transactions_admin_deposits;
DROP INDEX IF EXISTS idx_transactions_admin_withdrawals;
DROP INDEX IF EXISTS idx_transactions_admin_tips;
