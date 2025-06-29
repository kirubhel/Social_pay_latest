-- Migration rollback: Remove transaction analytics optimization indexes

DROP INDEX IF EXISTS idx_transactions_analytics_primary;
DROP INDEX IF EXISTS idx_transactions_analytics_status;
DROP INDEX IF EXISTS idx_transactions_analytics_type;
DROP INDEX IF EXISTS idx_transactions_analytics_medium;
DROP INDEX IF EXISTS idx_transactions_analytics_source;
DROP INDEX IF EXISTS idx_transactions_tip_analytics;
DROP INDEX IF EXISTS idx_transactions_merchant_analytics;
DROP INDEX IF EXISTS idx_transactions_amount_range;
DROP INDEX IF EXISTS idx_transactions_qr_analytics;
DROP INDEX IF EXISTS idx_transactions_user_type_amount;
DROP INDEX IF EXISTS idx_transactions_user_status_amount; 