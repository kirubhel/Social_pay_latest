-- Add phone_prefix and phone_number columns to auth.users table for Auth v2 compatibility

-- Add phone_prefix column if it doesn't exist
DO $$ BEGIN
    ALTER TABLE auth.users ADD COLUMN phone_prefix VARCHAR(10);
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Add phone_number column if it doesn't exist
DO $$ BEGIN
    ALTER TABLE auth.users ADD COLUMN phone_number VARCHAR(50);
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Add composite unique constraint for phone_prefix + phone_number combination
-- This ensures no duplicate phone numbers with same prefix
DO $$ BEGIN
    ALTER TABLE auth.users ADD CONSTRAINT unique_users_phone_prefix_number 
    UNIQUE (phone_prefix, phone_number);
EXCEPTION
    WHEN duplicate_table THEN NULL;
END $$;

-- Add index for better performance on phone number lookups
CREATE INDEX IF NOT EXISTS idx_users_phone_prefix_number ON auth.users(phone_prefix, phone_number);

-- Add index on phone_prefix for better performance
CREATE INDEX IF NOT EXISTS idx_users_phone_prefix ON auth.users(phone_prefix);

-- Add index on phone_number for better performance  
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON auth.users(phone_number);

-- Add check constraint to ensure phone_prefix is numeric and reasonable length
DO $$ BEGIN
    ALTER TABLE auth.users ADD CONSTRAINT chk_users_phone_prefix_format 
    CHECK (phone_prefix IS NULL OR (phone_prefix ~ '^[0-9]{1,4}$'));
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Add check constraint to ensure phone_number is numeric and reasonable length
DO $$ BEGIN
    ALTER TABLE auth.users ADD CONSTRAINT chk_users_phone_number_format 
    CHECK (phone_number IS NULL OR (phone_number ~ '^[0-9]{6,15}$'));
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$; 