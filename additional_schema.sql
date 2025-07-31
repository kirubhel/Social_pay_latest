-- Add 2FA support to users table
ALTER TABLE auth.users ADD COLUMN IF NOT EXISTS two_factor_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE auth.users ADD COLUMN IF NOT EXISTS two_factor_verified_at TIMESTAMP;

-- Create 2FA verification codes table
CREATE TABLE IF NOT EXISTS auth.two_factor_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_two_factor_codes_user_id ON auth.two_factor_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_two_factor_codes_expires_at ON auth.two_factor_codes(expires_at);
CREATE INDEX IF NOT EXISTS idx_two_factor_codes_used ON auth.two_factor_codes(used);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_two_factor_codes_updated_at 
    BEFORE UPDATE ON auth.two_factor_codes 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();