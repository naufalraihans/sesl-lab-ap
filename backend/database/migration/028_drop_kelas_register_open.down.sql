-- Rollback 028: kembalikan kolom akses register per kelas.
ALTER TABLE kelas ADD COLUMN IF NOT EXISTS is_register_open BOOLEAN NOT NULL DEFAULT FALSE;
