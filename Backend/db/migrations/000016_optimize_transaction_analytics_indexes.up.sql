-- Migration: Optimize transaction analytics indexes for billion-scale performance
-- This migration adds composite indexes specifically designed for analytics queries
-- Compatible with PostgreSQL 9.6+ (removed INCLUDE clauses for compatibility)
-- Note: CONCURRENTLY removed to allow running in migration transaction

-- Primary analytics index: user_id + created_at (covers most analytics queries)
-- This is the most important index for user-specific analytics
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_primary 
ON public.transactions(user_id, created_at DESC);

-- Analytics by status: user_id + status + created_at
-- Optimizes queries filtering by transaction status
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_status 
ON public.transactions(user_id, status, created_at DESC);

-- Analytics by type: user_id + type + created_at  
-- Optimizes queries filtering by transaction type (DEPOSIT/WITHDRAWAL)
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_type 
ON public.transactions(user_id, type, created_at DESC);

-- Analytics by medium: user_id + medium + created_at
-- Optimizes queries filtering by payment medium
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_medium 
ON public.transactions(user_id, medium, created_at DESC);

-- Analytics by source: user_id + transaction_source + created_at
-- Optimizes queries filtering by transaction source (QR, DIRECT, etc.)
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_source 
ON public.transactions(user_id, transaction_source, created_at DESC);

-- Tip analytics: user_id + has_tip + created_at
-- Optimizes tip-related analytics with partial index
CREATE INDEX IF NOT EXISTS idx_transactions_tip_analytics 
ON public.transactions(user_id, has_tip, created_at DESC)
WHERE has_tip = true;

-- Merchant analytics: merchant_id + created_at
-- Optimizes merchant-specific analytics queries
CREATE INDEX IF NOT EXISTS idx_transactions_merchant_analytics 
ON public.transactions(merchant_id, created_at DESC)
WHERE merchant_id IS NOT NULL;

-- Amount range queries: user_id + amount + created_at
-- Optimizes queries with amount range filters
CREATE INDEX IF NOT EXISTS idx_transactions_amount_range 
ON public.transactions(user_id, amount, created_at DESC);

-- QR tag analytics: user_id + qr_tag + created_at
-- Optimizes QR-specific analytics with partial index
CREATE INDEX IF NOT EXISTS idx_transactions_qr_analytics 
ON public.transactions(user_id, qr_tag, created_at DESC)
WHERE qr_tag IS NOT NULL;

-- Additional indexes for better analytics performance
-- User + type + amount (for type-specific amount analytics)
CREATE INDEX IF NOT EXISTS idx_transactions_user_type_amount 
ON public.transactions(user_id, type, amount);

-- User + status + amount (for status-specific amount analytics)
CREATE INDEX IF NOT EXISTS idx_transactions_user_status_amount 
ON public.transactions(user_id, status, amount);

-- Add table statistics refresh
ANALYZE public.transactions;
