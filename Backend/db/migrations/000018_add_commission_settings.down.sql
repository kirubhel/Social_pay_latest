-- Drop commission tiers table
DROP TABLE IF EXISTS admin.commission_tiers;

-- Drop admin settings table
DROP TABLE IF EXISTS admin.settings;

-- Remove commission columns from merchants table
ALTER TABLE merchants.merchants
DROP COLUMN IF EXISTS commission_percent,
DROP COLUMN IF EXISTS commission_cent,
DROP COLUMN IF EXISTS commission_active; 