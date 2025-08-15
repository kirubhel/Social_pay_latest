CREATE SCHEMA IF NOT EXISTS admin;

CREATE SCHEMA IF NOT EXISTS merchants;

CREATE TABLE IF NOT EXISTS merchants.merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    trading_name VARCHAR(255),
    business_registration_number VARCHAR(100) NOT NULL,
    tax_identification_number VARCHAR(100) NOT NULL,
    business_type VARCHAR(100) NOT NULL,
    industry_category VARCHAR(100),
    is_betting_company BOOLEAN DEFAULT false,
    lottery_certificate_number VARCHAR(100),
    website_url VARCHAR(255),
    established_date DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'pending_verification',
    commission_active BOOLEAN NOT NULL DEFAULT false,
    commission_percent DECIMAL(5,2),
    commission_cent DECIMAL(19,4)
);

CREATE TABLE IF NOT EXISTS admin.settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL UNIQUE,
    value JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_schema = 'merchants' 
        AND table_name = 'merchants' 
        AND column_name = 'commission_active'
    ) THEN
        ALTER TABLE merchants.merchants 
        ADD COLUMN commission_active BOOLEAN NOT NULL DEFAULT false,
        ADD COLUMN commission_percent DECIMAL(5,2),
        ADD COLUMN commission_cent DECIMAL(19,4);
    END IF;
END $$; 