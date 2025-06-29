-- Create hosted payment status enum type
CREATE TYPE hosted_payment_status AS ENUM (
    'PENDING',
    'COMPLETED',
    'EXPIRED',
    'CANCELED'
);

-- Create hosted payments table
CREATE TABLE IF NOT EXISTS public.hosted_payments (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- User and merchant information
    user_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    
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
    transaction_id UUID,
    
    -- Selected payment details (filled when user makes payment)
    selected_medium VARCHAR(50),
    selected_phone_number VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours')
);

-- Add foreign key constraints
ALTER TABLE public.hosted_payments 
ADD CONSTRAINT fk_hosted_payments_user 
FOREIGN KEY (user_id) REFERENCES auth.users(id);

ALTER TABLE public.hosted_payments 
ADD CONSTRAINT fk_hosted_payments_merchant 
FOREIGN KEY (merchant_id) REFERENCES merchants.merchants(id) ON DELETE CASCADE;

ALTER TABLE public.hosted_payments 
ADD CONSTRAINT fk_hosted_payments_transaction 
FOREIGN KEY (transaction_id) REFERENCES public.transactions(id) ON DELETE SET NULL;

-- Create indexes for hosted payments
CREATE INDEX IF NOT EXISTS idx_hosted_payments_user_id ON public.hosted_payments(user_id);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_merchant_id ON public.hosted_payments(merchant_id);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_status ON public.hosted_payments(status);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_created_at ON public.hosted_payments(created_at);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_expires_at ON public.hosted_payments(expires_at);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_reference ON public.hosted_payments(reference);
CREATE INDEX IF NOT EXISTS idx_hosted_payments_transaction_id ON public.hosted_payments(transaction_id);
