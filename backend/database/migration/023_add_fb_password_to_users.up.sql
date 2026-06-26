-- Hash+salt password lama dari Firebase Auth, untuk migrasi password bertahap.
-- Login pertama verifikasi scrypt Firebase -> rehash ke bcrypt -> kolom ini dikosongkan.
ALTER TABLE users ADD COLUMN IF NOT EXISTS fb_password_hash VARCHAR(255) NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS fb_password_salt VARCHAR(255) NULL;
