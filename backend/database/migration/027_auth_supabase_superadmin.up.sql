-- 027: Supabase Auth + superadmin (idempotent)
-- Menambah kolom email & supabase_user_id untuk flow register roster-gated
-- serta index unik parsial. Role 'superadmin' muat di VARCHAR(10) existing.

ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(150);
ALTER TABLE users ADD COLUMN IF NOT EXISTS supabase_user_id UUID;

-- Index unik parsial: email boleh NULL (asisten), tapi jika terisi harus unik (case-sensitive di DB, normalisasi di app)
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_email ON users(email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_supabase_user_id ON users(supabase_user_id) WHERE supabase_user_id IS NOT NULL;

-- Index bantu lookup login via email (lower)
CREATE INDEX IF NOT EXISTS idx_users_email_lower ON users(LOWER(email));
