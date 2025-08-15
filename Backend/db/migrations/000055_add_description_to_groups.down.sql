-- Remove unique constraint from auth.group_permissions
ALTER TABLE auth.group_permissions DROP CONSTRAINT IF EXISTS unique_group_permission;

-- Remove description column from auth.groups table
ALTER TABLE auth.groups DROP COLUMN IF EXISTS description; 