-- Fix phone data issues
-- This script ensures that phone records exist for all phone identities

-- First, let's see what phone identities are missing phone records
SELECT 'Missing Phone Records' as issue,
       pi.user_id,
       pi.phone_id,
       pi.created_at
FROM auth.phone_identities pi
LEFT JOIN auth.phones p ON pi.phone_id = p.id
WHERE p.id IS NULL;

-- Create missing phone records for phone identities that don't have them
-- We'll use a default prefix and number since we don't have the actual data
INSERT INTO auth.phones (id, prefix, number, created_at)
SELECT 
    pi.phone_id,
    '251', -- Default prefix for Ethiopia
    '000000000', -- Default number
    NOW()
FROM auth.phone_identities pi
LEFT JOIN auth.phones p ON pi.phone_id = p.id
WHERE p.id IS NULL
ON CONFLICT (id) DO NOTHING;

-- Update the default phone numbers with actual data if available
-- This is a placeholder - you'll need to manually update these with real data
UPDATE auth.phones 
SET prefix = '251', number = '911234567' -- Replace with actual phone number
WHERE id = '20415e38-a3c2-4963-a49f-3480b631fec7'
AND (prefix = '000' OR number = '000000000');

-- Verify the fix worked
SELECT 'Verification' as check_type,
       pi.user_id,
       pi.phone_id,
       p.prefix,
       p.number,
       CASE 
           WHEN p.id IS NULL THEN 'Still missing'
           ELSE 'Fixed'
       END as status
FROM auth.phone_identities pi
LEFT JOIN auth.phones p ON pi.phone_id = p.id
WHERE pi.user_id = '9b29fcdb-3326-4604-8a31-8ee844ec3fef'; 