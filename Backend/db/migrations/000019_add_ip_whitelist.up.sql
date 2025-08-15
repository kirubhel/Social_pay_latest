CREATE SCHEMA IF NOT EXISTS ip_whitelist;

CREATE TABLE IF NOT EXISTS ip_whitelist.whitelisted_ips (
    id UUID PRIMARY KEY,
    merchant_id UUID NOT NULL,
    ip_address CIDR NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (merchant_id) REFERENCES merchants.merchants(id) ON DELETE CASCADE,
    UNIQUE(merchant_id, ip_address)
);

CREATE INDEX idx_whitelisted_ips_merchant_id ON ip_whitelist.whitelisted_ips(merchant_id);
CREATE INDEX idx_whitelisted_ips_ip_address ON ip_whitelist.whitelisted_ips USING gist (ip_address inet_ops);
