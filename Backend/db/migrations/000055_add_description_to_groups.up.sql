-- Add description column to auth.groups table if it doesn't exist
DO $$ BEGIN
    ALTER TABLE auth.groups ADD COLUMN description TEXT;
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Add unique constraint to auth.group_permissions to prevent duplicate assignments
DO $$ BEGIN
    ALTER TABLE auth.group_permissions ADD CONSTRAINT unique_group_permission UNIQUE (group_id, permission_id);
EXCEPTION
    WHEN duplicate_table THEN NULL;
    WHEN duplicate_object THEN NULL;
END $$; 