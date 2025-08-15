CREATE SCHEMA IF NOT EXISTS admin;

-- Add commission columns to merchants table
ALTER TABLE merchants.merchants
ADD COLUMN commission_percent DECIMAL(5,2) DEFAULT NULL,
ADD COLUMN commission_cent DECIMAL(5,2) DEFAULT NULL,
ADD COLUMN commission_active BOOLEAN DEFAULT false;

-- Create admin settings table
CREATE TABLE IF NOT EXISTS admin.settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(50) NOT NULL UNIQUE,
    value JSONB NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Insert default commission settings
INSERT INTO admin.settings (key, value, description)
VALUES (
    'default_commission',
    '{
        "percent": 2.75,
        "cent": 0.00,
        "min_amount": 0.00,
        "max_amount": 1000000.00
    }'::jsonb,
    'Default commission settings for all transactions'
);


-- Create commission tiers table for dynamic commission rates
CREATE TABLE IF NOT EXISTS admin.commission_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    min_amount DECIMAL(19,4) NOT NULL,
    max_amount DECIMAL(19,4) NOT NULL,
    percent DECIMAL(5,2) NOT NULL,
    cent DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_amount_range CHECK (min_amount < max_amount)
);

-- Insert some default commission tiers
INSERT INTO admin.commission_tiers (min_amount, max_amount, percent, cent)
VALUES 
    (0, 1000, 3.00, 0.00),
    (1000, 5000, 2.75, 0.00),
    (5000, 10000, 2.50, 0.00),
    (10000, 1000000, 2.25, 0.00); 