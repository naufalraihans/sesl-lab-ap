DROP INDEX IF EXISTS uq_aktivasi_sesi_gel;
ALTER TABLE aktivasi_sesi ADD CONSTRAINT uq_aktivasi_sesi UNIQUE (sesi_praktikum_id, kelas_id, shift);
ALTER TABLE aktivasi_sesi DROP COLUMN gelombang;
ALTER TABLE users DROP COLUMN gelombang;
