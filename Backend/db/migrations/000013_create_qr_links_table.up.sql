-- Create QR link type enum
CREATE TYPE qr_link_type AS ENUM (
    'DYNAMIC',
    'STATIC'
);

-- Create QR link tag enum
CREATE TYPE qr_link_tag AS ENUM (
    'RESTAURANT',
    'DONATION',
    'SHOP'
);

-- Create QR links table
CREATE TABLE IF NOT EXISTS public.qr_links (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- User and merchant information
    user_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    
    -- QR link configuration
    type qr_link_type NOT NULL DEFAULT 'DYNAMIC',
    amount DECIMAL(20,2), -- Only required for STATIC type
    
    -- Supported payment mediums (JSON array)
    supported_methods JSONB NOT NULL,
    
    -- QR link metadata
    tag qr_link_tag NOT NULL DEFAULT 'SHOP',
    title VARCHAR(255),
    description TEXT,
    image_url TEXT,
    
    -- Features
    is_tip_enabled BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key constraints
ALTER TABLE public.qr_links 
ADD CONSTRAINT fk_qr_links_user 
FOREIGN KEY (user_id) REFERENCES auth.users(id);

ALTER TABLE public.qr_links 
ADD CONSTRAINT fk_qr_links_merchant 
FOREIGN KEY (merchant_id) REFERENCES merchants.merchants(id) ON DELETE CASCADE;

-- Create indexes for QR links
CREATE INDEX IF NOT EXISTS idx_qr_links_user_id ON public.qr_links(user_id);
CREATE INDEX IF NOT EXISTS idx_qr_links_merchant_id ON public.qr_links(merchant_id);
CREATE INDEX IF NOT EXISTS idx_qr_links_type ON public.qr_links(type);
CREATE INDEX IF NOT EXISTS idx_qr_links_tag ON public.qr_links(tag);
CREATE INDEX IF NOT EXISTS idx_qr_links_is_active ON public.qr_links(is_active);
CREATE INDEX IF NOT EXISTS idx_qr_links_created_at ON public.qr_links(created_at); 