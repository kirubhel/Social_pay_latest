-- Add merchant_id to auth.groups table if it doesn't exist
DO $$ BEGIN
    ALTER TABLE auth.groups ADD COLUMN merchant_id UUID REFERENCES merchants.merchants(id) ON DELETE CASCADE;
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Add merchant_id to auth.user_groups table if it doesn't exist  
DO $$ BEGIN
    ALTER TABLE auth.user_groups ADD COLUMN merchant_id UUID REFERENCES merchants.merchants(id) ON DELETE CASCADE;
EXCEPTION
    WHEN duplicate_column THEN NULL;
END $$;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_groups_merchant_id ON auth.groups(merchant_id);
CREATE INDEX IF NOT EXISTS idx_user_groups_merchant_id ON auth.user_groups(merchant_id);

-- Update existing user_groups to set merchant_id based on their group's merchant_id
UPDATE auth.user_groups 
SET merchant_id = g.merchant_id 
FROM auth.groups g 
WHERE auth.user_groups.group_id = g.id 
AND auth.user_groups.merchant_id IS NULL 
AND g.merchant_id IS NOT NULL; 