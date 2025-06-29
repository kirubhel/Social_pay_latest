-- QR Links Schema for SQLC generation
-- This should match the migration exactly

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