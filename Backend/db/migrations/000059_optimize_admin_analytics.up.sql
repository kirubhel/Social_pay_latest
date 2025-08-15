-- Specialized index for admin analytics (covers all fields needed)
CREATE INDEX IF NOT EXISTS idx_transactions_admin_analytics_optimized 
ON public.transactions (
    created_at DESC,
    type,
    status,
    base_amount,
    merchant_net,
    admin_net
); 

-- Partial index for deposits
CREATE INDEX IF NOT EXISTS idx_transactions_admin_deposits
ON public.transactions (created_at DESC, base_amount)
WHERE type = 'deposit';

-- Partial index for withdrawals  
CREATE INDEX IF NOT EXISTS idx_transactions_admin_withdrawals
ON public.transactions (created_at DESC, base_amount)
WHERE type = 'withdrawal';

-- Partial index for tips
CREATE INDEX IF NOT EXISTS idx_transactions_admin_tips
ON public.transactions (created_at DESC, base_amount)
WHERE type = 'tip';
