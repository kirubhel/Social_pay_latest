-- Rollback: Remove id column from auth.user_groups table and restore original structure

-- Drop indexes that were added
DROP INDEX IF EXISTS idx_user_groups_id;
DROP INDEX IF EXISTS idx_user_groups_user_id;
DROP INDEX IF EXISTS idx_user_groups_group_id;
DROP INDEX IF EXISTS idx_user_groups_created_at;

-- Drop unique constraint that was added
ALTER TABLE auth.user_groups DROP CONSTRAINT IF EXISTS unique_user_groups_user_group;

-- Drop primary key constraint on id
ALTER TABLE auth.user_groups DROP CONSTRAINT IF EXISTS user_groups_pkey;

-- Restore composite primary key on user_id + group_id (original structure)
-- Only if both columns exist
DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'auth' 
        AND table_name = 'user_groups' 
        AND column_name = 'user_id'
    ) AND EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'auth' 
        AND table_name = 'user_groups' 
        AND column_name = 'group_id'
    ) THEN
        ALTER TABLE auth.user_groups ADD CONSTRAINT user_groups_pkey PRIMARY KEY (user_id, group_id);
    END IF;
EXCEPTION
    WHEN OTHERS THEN NULL;
END $$;

-- Drop timestamp columns if they were added by this migration
ALTER TABLE auth.user_groups DROP COLUMN IF EXISTS created_at;
ALTER TABLE auth.user_groups DROP COLUMN IF EXISTS updated_at;

-- Drop id column
ALTER TABLE auth.user_groups DROP COLUMN IF EXISTS id; 