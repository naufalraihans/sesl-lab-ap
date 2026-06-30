ALTER TABLE users ADD COLUMN gelombang INT NULL;
ALTER TABLE aktivasi_sesi ADD COLUMN gelombang INT NULL;
ALTER TABLE aktivasi_sesi DROP CONSTRAINT IF EXISTS uq_aktivasi_sesi;
CREATE UNIQUE INDEX uq_aktivasi_sesi_gel ON aktivasi_sesi (sesi_praktikum_id, kelas_id, shift, COALESCE(gelombang, 0));
