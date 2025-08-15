-- Create schema
CREATE SCHEMA IF NOT EXISTS ip_whitelist;

-- Create table
CREATE TABLE IF NOT EXISTS ip_whitelist.whitelisted_ips (
    id UUID PRIMARY KEY,
    merchant_id UUID NOT NULL,
    ip_address CIDR NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (merchant_id) REFERENCES merchant.merchants(id) ON DELETE CASCADE,
    UNIQUE(merchant_id, ip_address)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_whitelisted_ips_merchant_id ON ip_whitelist.whitelisted_ips(merchant_id);
CREATE INDEX IF NOT EXISTS idx_whitelisted_ips_ip_address ON ip_whitelist.whitelisted_ips USING gist (ip_address inet_ops);

-- Create trigger function
CREATE OR REPLACE FUNCTION ip_whitelist.update_whitelisted_ips_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
CREATE TRIGGER update_whitelisted_ips_updated_at
    BEFORE UPDATE ON ip_whitelist.whitelisted_ips
    FOR EACH ROW
    EXECUTE FUNCTION ip_whitelist.update_whitelisted_ips_updated_at(); 