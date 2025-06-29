CREATE SCHEMA IF NOT EXISTS merchants;

CREATE TABLE IF NOT EXISTS merchants.merchants (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL REFERENCES auth.users(id) ON DELETE RESTRICT,
    legal_name VARCHAR(255) NOT NULL,
    trading_name VARCHAR(255),
    business_registration_number VARCHAR(100) NOT NULL UNIQUE,
    tax_identification_number VARCHAR(100) NOT NULL UNIQUE,
    business_type VARCHAR(100) NOT NULL, -- e.g., 'retail', 'ecommerce', 'betting'
    industry_category VARCHAR(100),
    is_betting_company BOOLEAN DEFAULT FALSE,
    lottery_certificate_number VARCHAR(100), -- Only for betting companies
    website_url VARCHAR(255),
    established_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'pending_verification' -- pending_verification, active, suspended, terminated
);

COMMENT ON TABLE merchants.merchants IS 'Core merchant information table';
COMMENT ON COLUMN merchants.merchants.lottery_certificate_number IS 'Required for betting/gaming merchants'; 


CREATE TABLE IF NOT EXISTS merchants.addresses (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id uuid NOT NULL REFERENCES merchants.merchants(id) ON DELETE CASCADE,
    address_type VARCHAR(50) NOT NULL, -- 'legal', 'operational', 'billing'
    street_address_1 VARCHAR(255) NOT NULL,
    street_address_2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    postal_code VARCHAR(50),
    country VARCHAR(100) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_addresses_merchant_id ON merchants.addresses(merchant_id);

CREATE TABLE IF NOT EXISTS merchants.contacts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id uuid NOT NULL REFERENCES merchants.merchants(id) ON DELETE CASCADE,
    contact_type VARCHAR(50) NOT NULL, -- 'primary', 'technical', 'billing', 'support'
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    position VARCHAR(100),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contacts_merchant_id ON merchants.contacts(merchant_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_contacts_email_unique ON merchants.contacts(email) WHERE is_verified = TRUE;


CREATE TABLE IF NOT EXISTS merchants.documents (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id uuid NOT NULL REFERENCES merchants.merchants(id) ON DELETE CASCADE,
    document_type VARCHAR(100) NOT NULL, -- 'business_license', 'tin_certificate', 'lottery_certificate', 'bank_statement'
    document_number VARCHAR(100),
    file_url VARCHAR(255) NOT NULL,
    file_hash VARCHAR(64), -- SHA-256 hash of the file
    verified_by uuid REFERENCES auth.users(id),
    verified_at TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, approved, rejected
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_documents_merchant_id ON merchants.documents(merchant_id);


CREATE TABLE IF NOT EXISTS merchants.bank_accounts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id uuid NOT NULL REFERENCES merchants.merchants(id) ON DELETE CASCADE,
    account_holder_name VARCHAR(255) NOT NULL,
    bank_name VARCHAR(255) NOT NULL,
    bank_code VARCHAR(50) NOT NULL,
    branch_code VARCHAR(50),
    account_number VARCHAR(50) NOT NULL,
    account_type VARCHAR(50) NOT NULL, -- 'checking', 'savings'
    currency VARCHAR(3) NOT NULL DEFAULT 'ETB',
    is_primary BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_document_id uuid REFERENCES merchants.documents(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bank_accounts_merchant_id ON merchants.bank_accounts(merchant_id);


CREATE TABLE IF NOT EXISTS merchants.api_keys (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id uuid NOT NULL REFERENCES merchants.merchants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    public_key TEXT NOT NULL,
    private_key_encrypted TEXT, -- Encrypted with master key
    key_type VARCHAR(50) NOT NULL, -- 'test', 'live'
    allowed_ips CIDR[],
    scopes TEXT[] NOT NULL DEFAULT ARRAY['payments:create', 'payments:read'],
	api_key TEXT NOT NULL,
	store TEXT NOT NULL,
    expiry_date TIMESTAMPTZ,
	service TEXT NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by uuid REFERENCES auth.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_api_keys_merchant_id ON merchants.api_keys(merchant_id);


CREATE TABLE IF NOT EXISTS merchants.settings (
    merchant_id uuid PRIMARY KEY REFERENCES merchants.merchants(id) ON DELETE CASCADE,
    default_currency VARCHAR(3) NOT NULL DEFAULT 'ETB',
    default_language VARCHAR(10) NOT NULL DEFAULT 'en',
    checkout_theme VARCHAR(50),
    enable_webhooks BOOLEAN DEFAULT FALSE,
    webhook_url VARCHAR(255),
    webhook_secret VARCHAR(255),
    auto_settlement BOOLEAN DEFAULT TRUE,
    settlement_frequency VARCHAR(50) DEFAULT 'daily', -- daily, weekly, monthly
    risk_settings JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS merchants.audit_logs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id uuid REFERENCES merchants.merchants(id) ON DELETE SET NULL,
    user_id uuid REFERENCES auth.users(id),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100),
    entity_id uuid,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_merchant_id ON merchants.audit_logs(merchant_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON merchants.audit_logs(created_at);