-- Lindungi nilai (pengerjaan_course) dari terhapus tak sengaja via CASCADE saat
-- admin menghapus sesi/course/kelas/mahasiswa. RESTRICT menolak hapus induk
-- selama masih ada nilai. (lihat bugAudit B8)
ALTER TABLE pengerjaan_course DROP CONSTRAINT IF EXISTS fk_pengerjaan_mhs;
ALTER TABLE pengerjaan_course ADD CONSTRAINT fk_pengerjaan_mhs FOREIGN KEY (mahasiswa_id) REFERENCES users(id) ON DELETE RESTRICT;
ALTER TABLE pengerjaan_course DROP CONSTRAINT IF EXISTS fk_pengerjaan_aktivasi;
ALTER TABLE pengerjaan_course ADD CONSTRAINT fk_pengerjaan_aktivasi FOREIGN KEY (aktivasi_sesi_id) REFERENCES aktivasi_sesi(id) ON DELETE RESTRICT;
ALTER TABLE pengerjaan_course DROP CONSTRAINT IF EXISTS fk_pengerjaan_course;
ALTER TABLE pengerjaan_course ADD CONSTRAINT fk_pengerjaan_course FOREIGN KEY (course_id) REFERENCES course(id) ON DELETE RESTRICT;
