-- Add id column to auth.user_groups table for Auth v2 compatibility

-- Add id column as UUID primary key if it doesn't exist
DO $$ BEGIN
    -- First check if id column already exists
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'auth' 
        AND table_name = 'user_groups' 
        AND column_name = 'id'
    ) THEN
        -- Add id column
        ALTER TABLE auth.user_groups ADD COLUMN id UUID DEFAULT gen_random_uuid();
        
        -- Make it NOT NULL after adding default values
        ALTER TABLE auth.user_groups ALTER COLUMN id SET NOT NULL;
        
        -- Drop existing primary key if it exists (usually composite key)
        BEGIN
            ALTER TABLE auth.user_groups DROP CONSTRAINT IF EXISTS user_groups_pkey;
        EXCEPTION
            WHEN OTHERS THEN NULL;
        END;
        
        -- Add new primary key on id column
        ALTER TABLE auth.user_groups ADD CONSTRAINT user_groups_pkey PRIMARY KEY (id);
        
        -- Ensure the composite unique constraint exists for user_id + group_id
        -- This prevents duplicate user-group assignments
        ALTER TABLE auth.user_groups ADD CONSTRAINT unique_user_groups_user_group 
        UNIQUE (user_id, group_id);
    END IF;
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Add created_at column if it doesn't exist  
DO $$ BEGIN
    ALTER TABLE auth.user_groups ADD COLUMN created_at TIMESTAMP DEFAULT NOW();
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Add updated_at column if it doesn't exist
DO $$ BEGIN
    ALTER TABLE auth.user_groups ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Update existing rows to have timestamps if they're null
UPDATE auth.user_groups 
SET created_at = NOW() 
WHERE created_at IS NULL;

UPDATE auth.user_groups 
SET updated_at = NOW() 
WHERE updated_at IS NULL;

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_groups_id ON auth.user_groups(id);
CREATE INDEX IF NOT EXISTS idx_user_groups_user_id ON auth.user_groups(user_id);  
CREATE INDEX IF NOT EXISTS idx_user_groups_group_id ON auth.user_groups(group_id);
CREATE INDEX IF NOT EXISTS idx_user_groups_created_at ON auth.user_groups(created_at); 