-- Auth v2 Migration Rollback: Remove auth v2 tables and fields

-- Drop cleanup functions
DROP FUNCTION IF EXISTS cleanup_expired_sessions();
DROP FUNCTION IF EXISTS cleanup_old_auth_activities();
DROP FUNCTION IF EXISTS cleanup_expired_otps();

-- Drop trigger and function for otp_requests
DROP TRIGGER IF EXISTS trigger_update_otp_requests_updated_at ON auth.otp_requests;
DROP FUNCTION IF EXISTS update_otp_requests_updated_at();

-- Drop constraints that were added
ALTER TABLE auth.otp_requests DROP CONSTRAINT IF EXISTS chk_otp_code_format;
ALTER TABLE auth.auth_activities DROP CONSTRAINT IF EXISTS chk_activity_type_valid;
ALTER TABLE auth.sessions DROP CONSTRAINT IF EXISTS chk_token_length;

-- Drop indexes that were added
DROP INDEX IF EXISTS idx_otp_requests_token;
DROP INDEX IF EXISTS idx_otp_requests_phone_id;
DROP INDEX IF EXISTS idx_otp_requests_expires_at;

DROP INDEX IF EXISTS idx_auth_activities_user_id;
DROP INDEX IF EXISTS idx_auth_activities_activity_type;
DROP INDEX IF EXISTS idx_auth_activities_created_at;
DROP INDEX IF EXISTS idx_auth_activities_success;

DROP INDEX IF EXISTS idx_sessions_token;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_active;

DROP INDEX IF EXISTS idx_permissions_resource_id;

DROP INDEX IF EXISTS idx_users_user_type;
DROP INDEX IF EXISTS idx_phone_identities_user_id;
DROP INDEX IF EXISTS idx_phone_identities_phone_id;
DROP INDEX IF EXISTS idx_phones_prefix_number;
DROP INDEX IF EXISTS idx_password_identities_user_id;
DROP INDEX IF EXISTS idx_user_groups_user_id;
DROP INDEX IF EXISTS idx_user_groups_group_id;
DROP INDEX IF EXISTS idx_group_permissions_group_id;
DROP INDEX IF EXISTS idx_group_permissions_permission_id;

-- Drop new tables
DROP TABLE IF EXISTS auth.auth_activities;
DROP TABLE IF EXISTS auth.otp_requests;

-- Remove columns that were added
ALTER TABLE auth.sessions DROP COLUMN IF EXISTS refresh_token;
ALTER TABLE auth.sessions DROP COLUMN IF EXISTS expires_at;
ALTER TABLE auth.sessions DROP COLUMN IF EXISTS active;
ALTER TABLE auth.permissions DROP COLUMN IF EXISTS effect; 