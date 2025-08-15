-- Add is_admin_operation column to auth.operations table
DO $$ BEGIN
    ALTER TABLE auth.operations ADD COLUMN is_admin_operation BOOLEAN DEFAULT FALSE;
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Update existing admin operations to set is_admin_operation = TRUE
UPDATE auth.operations SET is_admin_operation = TRUE 
WHERE name IN ('ADMIN_READ', 'ADMIN_CREATE', 'ADMIN_UPDATE', 'ADMIN_DELETE', 'ADMIN_ALL');

-- Add email column to auth.users table
DO $$ BEGIN
    ALTER TABLE auth.users ADD COLUMN email VARCHAR(255) UNIQUE;
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Create index on email for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON auth.users(email);

-- Add not null constraint for email after adding the column
-- First we'll need to update existing users with placeholder emails
DO $$ 
BEGIN
    -- Update existing users with placeholder emails based on phone
    UPDATE auth.users 
    SET email = CONCAT(phone_prefix, phone_number, '@socialpay.placeholder') 
    WHERE email IS NULL;
    
    -- Now make email NOT NULL
    ALTER TABLE auth.users ALTER COLUMN email SET NOT NULL;
EXCEPTION
    WHEN OTHERS THEN 
        RAISE NOTICE 'Could not set email to NOT NULL, please ensure all users have email addresses';
END $$; 