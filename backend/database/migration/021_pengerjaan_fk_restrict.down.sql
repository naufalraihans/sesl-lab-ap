ALTER TABLE pengerjaan_course DROP CONSTRAINT IF EXISTS fk_pengerjaan_mhs;
ALTER TABLE pengerjaan_course ADD CONSTRAINT fk_pengerjaan_mhs FOREIGN KEY (mahasiswa_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE pengerjaan_course DROP CONSTRAINT IF EXISTS fk_pengerjaan_aktivasi;
ALTER TABLE pengerjaan_course ADD CONSTRAINT fk_pengerjaan_aktivasi FOREIGN KEY (aktivasi_sesi_id) REFERENCES aktivasi_sesi(id) ON DELETE CASCADE;
ALTER TABLE pengerjaan_course DROP CONSTRAINT IF EXISTS fk_pengerjaan_course;
ALTER TABLE pengerjaan_course ADD CONSTRAINT fk_pengerjaan_course FOREIGN KEY (course_id) REFERENCES course(id) ON DELETE CASCADE;
