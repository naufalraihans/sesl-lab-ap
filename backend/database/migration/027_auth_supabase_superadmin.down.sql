DROP INDEX IF EXISTS idx_users_email_lower;
DROP INDEX IF EXISTS uq_users_supabase_user_id;
DROP INDEX IF EXISTS uq_users_email;
ALTER TABLE users DROP COLUMN IF EXISTS supabase_user_id;
ALTER TABLE users DROP COLUMN IF EXISTS email;
