-- Quick performance indexes for admin analytics
-- Run these first to improve query performance

-- Core admin analytics index (covers date range + main aggregation fields)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_admin_analytics_core 
ON public.transactions (created_at DESC, base_amount, merchant_net, admin_net, vat_amount, fee_amount, customer_net);

-- Type-specific indexes for parallel queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_admin_deposits 
ON public.transactions (created_at DESC, base_amount) 
WHERE type = 'deposit';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_admin_withdrawals 
ON public.transactions (created_at DESC, base_amount) 
WHERE type = 'withdrawal';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_admin_tips 
ON public.transactions (created_at DESC, base_amount) 
WHERE type = 'tip';

-- Status-based index for filtering
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_admin_status_date 
ON public.transactions (status, created_at DESC, base_amount);
