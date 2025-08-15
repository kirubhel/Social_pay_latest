-- Auth v2 Migration: Add missing tables and fields for the new authentication system

-- Create auth.otp_requests table for OTP management
CREATE TABLE IF NOT EXISTS auth.otp_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_id UUID NOT NULL REFERENCES auth.phones(id) ON DELETE CASCADE,
    code VARCHAR(10) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index for faster OTP lookups
CREATE INDEX idx_otp_requests_token ON auth.otp_requests(token);
CREATE INDEX idx_otp_requests_phone_id ON auth.otp_requests(phone_id);
CREATE INDEX idx_otp_requests_expires_at ON auth.otp_requests(expires_at);

-- Create auth.auth_activities table for authentication activity logging
CREATE TABLE IF NOT EXISTS auth.auth_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    activity_type VARCHAR(50) NOT NULL, -- LOGIN_SUCCESS, LOGIN_FAILURE, REGISTER, LOGOUT, etc.
    ip_address INET,
    user_agent TEXT,
    device_name VARCHAR(255),
    success BOOLEAN NOT NULL,
    details JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for auth activities
CREATE INDEX idx_auth_activities_user_id ON auth.auth_activities(user_id);
CREATE INDEX idx_auth_activities_activity_type ON auth.auth_activities(activity_type);
CREATE INDEX idx_auth_activities_created_at ON auth.auth_activities(created_at);
CREATE INDEX idx_auth_activities_success ON auth.auth_activities(success);

-- Add refresh_token column to sessions table if it doesn't exist
DO $$ BEGIN
    ALTER TABLE auth.sessions ADD COLUMN refresh_token VARCHAR(500);
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Add expires_at column to sessions table if it doesn't exist (some versions might not have it)
DO $$ BEGIN
    ALTER TABLE auth.sessions ADD COLUMN expires_at TIMESTAMP;
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Add active column to sessions table if it doesn't exist
DO $$ BEGIN
    ALTER TABLE auth.sessions ADD COLUMN active BOOLEAN DEFAULT TRUE;
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Create indexes for sessions table optimizations
CREATE INDEX IF NOT EXISTS idx_sessions_token ON auth.sessions(token);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON auth.sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON auth.sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_active ON auth.sessions(active);

-- Ensure permissions table has the right structure for v2
-- Add effect column to permissions if it doesn't exist
DO $$ BEGIN
    ALTER TABLE auth.permissions ADD COLUMN effect VARCHAR(10) DEFAULT 'allow';
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Create index on permissions for faster lookups
CREATE INDEX IF NOT EXISTS idx_permissions_resource_id ON auth.permissions(resource_id);

-- Ensure we have proper indexes on existing tables for performance
CREATE INDEX IF NOT EXISTS idx_users_user_type ON auth.users(user_type);
CREATE INDEX IF NOT EXISTS idx_phone_identities_user_id ON auth.phone_identities(user_id);
CREATE INDEX IF NOT EXISTS idx_phone_identities_phone_id ON auth.phone_identities(phone_id);
CREATE INDEX IF NOT EXISTS idx_phones_prefix_number ON auth.phones(prefix, number);
CREATE INDEX IF NOT EXISTS idx_password_identities_user_id ON auth.password_identities(user_id);
CREATE INDEX IF NOT EXISTS idx_user_groups_user_id ON auth.user_groups(user_id);
CREATE INDEX IF NOT EXISTS idx_user_groups_group_id ON auth.user_groups(group_id);
CREATE INDEX IF NOT EXISTS idx_group_permissions_group_id ON auth.group_permissions(group_id);
CREATE INDEX IF NOT EXISTS idx_group_permissions_permission_id ON auth.group_permissions(permission_id);

-- Add constraints to ensure data integrity
-- Ensure OTP codes are numeric and have proper length
ALTER TABLE auth.otp_requests ADD CONSTRAINT chk_otp_code_format CHECK (code ~ '^[0-9]{6}$');

-- Ensure activity types are valid
ALTER TABLE auth.auth_activities ADD CONSTRAINT chk_activity_type_valid 
    CHECK (activity_type IN (
        'REGISTER', 'LOGIN_SUCCESS', 'LOGIN_FAILURE', 'LOGOUT', 'PASSWORD_CHANGE',
        'OTP_REQUEST', 'OTP_VERIFY_SUCCESS', 'OTP_VERIFY_FAILURE', 'PERMISSION_DENIED',
        'SESSION_EXPIRED', 'ACCOUNT_LOCKED', 'ACCOUNT_UNLOCKED', 'PROFILE_UPDATE'
    ));

-- Ensure sessions have proper token lengths
ALTER TABLE auth.sessions ADD CONSTRAINT chk_token_length CHECK (LENGTH(token) >= 10);

-- Update existing data to have proper defaults
UPDATE auth.sessions SET active = TRUE WHERE active IS NULL;
UPDATE auth.permissions SET effect = 'allow' WHERE effect IS NULL;

-- Add trigger to update updated_at timestamp on otp_requests
CREATE OR REPLACE FUNCTION update_otp_requests_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_otp_requests_updated_at ON auth.otp_requests;
CREATE TRIGGER trigger_update_otp_requests_updated_at
    BEFORE UPDATE ON auth.otp_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_otp_requests_updated_at();

-- Create function to clean up expired OTPs (can be called by a cron job)
CREATE OR REPLACE FUNCTION cleanup_expired_otps()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM auth.otp_requests 
    WHERE expires_at < NOW() - INTERVAL '1 hour'
    OR (used = TRUE AND created_at < NOW() - INTERVAL '24 hours');
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Create function to clean up old auth activities (keep last 90 days)
CREATE OR REPLACE FUNCTION cleanup_old_auth_activities()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM auth.auth_activities 
    WHERE created_at < NOW() - INTERVAL '90 days';
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
