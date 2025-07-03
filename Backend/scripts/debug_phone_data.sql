-- Debug script to check phone data relationship
-- Run this to see what's happening with the phone data

-- Check the specific user's phone identity
SELECT 'Phone Identity' as table_name, 
       pi.user_id, 
       pi.phone_id, 
       pi.created_at
FROM auth.phone_identities pi
WHERE pi.user_id = '9b29fcdb-3326-4604-8a31-8ee844ec3fef';

-- Check the phone record
SELECT 'Phone Record' as table_name,
       p.id,
       p.prefix,
       p.number,
       p.created_at
FROM auth.phones p
WHERE p.id = '20415e38-a3c2-4963-a49f-3480b631fec7';

-- Check if the JOIN works
SELECT 'JOIN Test' as table_name,
       pi.user_id,
       pi.phone_id,
       p.id as phone_table_id,
       p.prefix,
       p.number,
       CASE 
           WHEN p.id IS NULL THEN 'Phone record missing'
           WHEN pi.phone_id != p.id THEN 'ID mismatch'
           ELSE 'OK'
       END as status
FROM auth.phone_identities pi
LEFT JOIN auth.phones p ON pi.phone_id = p.id
WHERE pi.user_id = '9b29fcdb-3326-4604-8a31-8ee844ec3fef';

-- Check all phone identities for this user
SELECT 'All Phone Identities' as table_name,
       pi.id,
       pi.user_id,
       pi.phone_id,
       pi.created_at
FROM auth.phone_identities pi
WHERE pi.user_id = '9b29fcdb-3326-4604-8a31-8ee844ec3fef';

-- Check all phones that might be related
SELECT 'All Related Phones' as table_name,
       p.id,
       p.prefix,
       p.number,
       p.created_at
FROM auth.phones p
WHERE p.id IN (
    SELECT phone_id 
    FROM auth.phone_identities 
    WHERE user_id = '9b29fcdb-3326-4604-8a31-8ee844ec3fef'
); 