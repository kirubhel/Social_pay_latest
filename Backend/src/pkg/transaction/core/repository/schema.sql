-- Create transaction status enum type
CREATE TYPE transaction_status AS ENUM (
    'INITIATED',
    'PENDING',
    'SUCCESS', 
    'FAILED',
    'REFUNDED',
    'EXPIRED',
    'CANCELED'
);

-- Create transaction source enum type
CREATE TYPE transaction_source AS ENUM (
    'DIRECT',
    'HOSTED_CHECKOUT',
    'QR_PAYMENT',
    'WITHDRAWAL'
);

-- Create hosted payment status enum type
CREATE TYPE hosted_payment_status AS ENUM (
    'PENDING',
    'COMPLETED',
    'EXPIRED',
    'CANCELED'
);


-- Schema definition for transactions table
-- This file is used by SQLC to generate Go code
-- When modifying the table structure:
-- 1. Create a new migration in db/migrations/
-- 2. Update this schema file to match
-- 3. Run sqlc generate

CREATE TABLE IF NOT EXISTS public.transactions (
    -- Primary identifier
    id UUID PRIMARY KEY,
    
    -- User and account information
    phone_number VARCHAR(50),
    user_id UUID NOT NULL,
    merchant_id UUID REFERENCES merchants.merchants(id) ON DELETE SET NULL,
    
    -- Transaction type and medium
    type VARCHAR(50) NOT NULL,
    medium VARCHAR(50) NOT NULL,
    
    -- Reference and description fields
    reference VARCHAR(100),
    comment TEXT,
    reference_number VARCHAR(100),
    description TEXT,
    
    -- Status and verification
    verified BOOLEAN DEFAULT false,
    status transaction_status NOT NULL DEFAULT 'INITIATED',
    test BOOLEAN DEFAULT false,
    has_challenge BOOLEAN DEFAULT false,
    webhook_received BOOLEAN DEFAULT false,
    
    -- Time-related fields
    ttl BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    confirm_timestamp TIMESTAMP WITH TIME ZONE,
    
    -- Financial information
    base_amount DECIMAL(20,2) NOT NULL,
    fee_amount DECIMAL(20,2),
    admin_net DECIMAL(20,2),
    vat_amount DECIMAL(20,2),
    merchant_net DECIMAL(20,2),
    customer_net DECIMAL(20,2),
    total_amount DECIMAL(20,2),
    currency VARCHAR(3) DEFAULT 'ETB',
    
    -- Additional data
    details JSONB,
    token VARCHAR(255),
    
    -- Provider information
    provider_tx_id VARCHAR(255),
    provider_data JSONB,
    merchant_pays_fee BOOLEAN DEFAULT FALSE,
    
    -- URLs for callbacks and redirects
    callback_url TEXT,
    success_url TEXT,
    failed_url TEXT,
    
    -- QR Payment Context
    transaction_source transaction_source DEFAULT 'DIRECT',
    qr_link_id UUID,
    hosted_checkout_id UUID,
    qr_tag VARCHAR(30),
    
    -- Tip Information
    has_tip BOOLEAN DEFAULT FALSE,
    tip_amount DECIMAL(20,2),
    tipee_phone VARCHAR(20),
    tip_medium VARCHAR(20),
    tip_transaction_id UUID,
    tip_processed BOOLEAN DEFAULT FALSE,
    
    -- Foreign key constraints
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES auth.users(id)
);

-- Create hosted payments table
CREATE TABLE IF NOT EXISTS public.hosted_payments (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- User and merchant information
    user_id UUID NOT NULL,
    merchant_id UUID NOT NULL REFERENCES merchants.merchants(id) ON DELETE SET NULL,
    
    -- Payment details
    amount DECIMAL(20,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ETB',
    description TEXT,
    reference VARCHAR(100) NOT NULL,
    
    -- Supported payment mediums (JSON array)
    supported_mediums JSONB NOT NULL,
    
    -- Optional phone number from merchant
    phone_number VARCHAR(50),
    
    -- URLs for redirects and callbacks
    success_url TEXT NOT NULL,
    failed_url TEXT NOT NULL,
    callback_url TEXT,
    
    -- Status and transaction linking
    status hosted_payment_status NOT NULL DEFAULT 'PENDING',
    transaction_id UUID REFERENCES public.transactions(id) ON DELETE SET NULL,
    
    -- Selected payment details (filled when user makes payment)
    selected_medium VARCHAR(50),
    selected_phone_number VARCHAR(50),
    
    -- Fee configuration
    merchant_pays_fee BOOLEAN NOT NULL DEFAULT false,
    accept_tip BOOLEAN NOT NULL DEFAULT false,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours'),
    
    -- Foreign key constraints
    CONSTRAINT fk_hosted_payments_user FOREIGN KEY (user_id) REFERENCES auth.users(id),
    CONSTRAINT fk_hosted_payments_transaction FOREIGN KEY (transaction_id) REFERENCES public.transactions(id)
);

-- Create QR links table
-- QR link type enum
CREATE TYPE qr_link_type AS ENUM (
    'DYNAMIC',
    'STATIC'
);

-- QR link tag enum  
CREATE TYPE qr_link_tag AS ENUM (
    'RESTAURANT',
    'DONATION',
    'SHOP'
);

-- Create QR links table
CREATE TABLE IF NOT EXISTS public.qr_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    type qr_link_type NOT NULL DEFAULT 'DYNAMIC',
    amount DECIMAL(20,2),
    supported_methods JSONB NOT NULL,
    tag qr_link_tag NOT NULL DEFAULT 'SHOP',
    title VARCHAR(255),
    description TEXT,
    image_url TEXT,
    is_tip_enabled BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON public.transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_merchant_id ON public.transactions(merchant_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON public.transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON public.transactions(type);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON public.transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_reference ON public.transactions(reference);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_number ON public.transactions(reference_number);
CREATE INDEX IF NOT EXISTS idx_transactions_qr_link_id ON public.transactions(qr_link_id);
CREATE INDEX IF NOT EXISTS idx_transactions_hosted_checkout_id ON public.transactions(hosted_checkout_id);
CREATE INDEX IF NOT EXISTS idx_transactions_tip_processing ON public.transactions(has_tip, tip_processed, status);
CREATE INDEX IF NOT EXISTS idx_transactions_source ON public.transactions(transaction_source);
CREATE INDEX IF NOT EXISTS idx_transactions_qr_tag ON public.transactions(qr_tag);

-- PERFORMANCE OPTIMIZATION: Composite indexes for analytics queries
-- These indexes are specifically designed for high-performance analytics on billions of transactions

-- Primary analytics index: user_id + created_at (covers most analytics queries)
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_primary ON public.transactions(user_id, created_at DESC, base_amount, merchant_net, type, status, has_tip, tip_amount);

-- Analytics by status: user_id + status + created_at
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_status ON public.transactions(user_id, status, created_at DESC, base_amount, merchant_net, type);

-- Analytics by type: user_id + type + created_at  
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_type ON public.transactions(user_id, type, created_at DESC, base_amount, merchant_net, status);

-- Analytics by medium: user_id + medium + created_at
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_medium ON public.transactions(user_id, medium, created_at DESC, base_amount, merchant_net, type, status);

-- Analytics by source: user_id + transaction_source + created_at
CREATE INDEX IF NOT EXISTS idx_transactions_analytics_source ON public.transactions(user_id, transaction_source, created_at DESC, base_amount, merchant_net, type, status);

-- Chart data optimization: user_id + created_at (date truncation friendly)
CREATE INDEX IF NOT EXISTS idx_transactions_chart_data ON public.transactions(user_id, date_trunc('day', created_at), base_amount, type, status);

-- Tip analytics: user_id + has_tip + created_at
CREATE INDEX IF NOT EXISTS idx_transactions_tip_analytics ON public.transactions(user_id, has_tip, created_at DESC, tip_amount)
WHERE has_tip = true;

-- Merchant analytics (for admin/merchant-specific queries): merchant_id + created_at
CREATE INDEX IF NOT EXISTS idx_transactions_merchant_analytics ON public.transactions(merchant_id, created_at DESC, base_amount, merchant_net, type, status, user_id)
WHERE merchant_id IS NOT NULL;

-- Amount range queries: user_id + base_amount + created_at
CREATE INDEX IF NOT EXISTS idx_transactions_amount_range ON public.transactions(user_id, base_amount, created_at DESC, merchant_net, type, status);

-- QR tag analytics: user_id + qr_tag + created_at
CREATE INDEX IF NOT EXISTS idx_transactions_qr_analytics ON public.transactions(user_id, qr_tag, created_at DESC, base_amount, merchant_net, type, status)
WHERE qr_tag IS NOT NULL;

-- Index for provider transaction ID
CREATE INDEX IF NOT EXISTS idx_transactions_provider_tx_id ON public.transactions(provider_tx_id);

-- Indexes for hosted payments
CREATE INDEX IF NOT EXISTS idx_hosted_payments_user_id ON public.hosted_payments(user_id);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_merchant_id ON public.hosted_payments(merchant_id);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_status ON public.hosted_payments(status);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_created_at ON public.hosted_payments(created_at);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_expires_at ON public.hosted_payments(expires_at);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_reference ON public.hosted_payments(reference);

-- Indexes for QR links
CREATE INDEX IF NOT EXISTS idx_qr_links_user_id ON public.qr_links(user_id);
CREATE INDEX IF NOT EXISTS idx_qr_links_merchant_id ON public.qr_links(merchant_id);
CREATE INDEX IF NOT EXISTS idx_qr_links_type ON public.qr_links(type);
CREATE INDEX IF NOT EXISTS idx_qr_links_tag ON public.qr_links(tag);
CREATE INDEX IF NOT EXISTS idx_qr_links_is_active ON public.qr_links(is_active);
CREATE INDEX IF NOT EXISTS idx_qr_links_created_at ON public.qr_links(created_at);

-- name: OverrideTransactionStatus :exec
UPDATE public.transactions
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- Create a new pivot table for transaction status overrides
CREATE TABLE IF NOT EXISTS public.transaction_status_overrides (
    id UUID PRIMARY KEY,
    transaction_id UUID NOT NULL,
    reason TEXT NOT NULL,
    admin_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES public.transactions(id) ON DELETE CASCADE
); 