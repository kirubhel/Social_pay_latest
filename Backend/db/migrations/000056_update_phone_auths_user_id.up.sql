-- Clear existing phone_auths records to avoid data integrity issues
TRUNCATE TABLE auth.phone_auths;

-- Add user_id column to auth.phone_auths table if it doesn't exist
DO $$ BEGIN
    ALTER TABLE auth.phone_auths ADD COLUMN user_id UUID NOT NULL;
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Add foreign key constraint
ALTER TABLE auth.phone_auths ADD CONSTRAINT fk_phone_auths_user_id 
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE;

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_phone_auths_user_id ON auth.phone_auths(user_id); 