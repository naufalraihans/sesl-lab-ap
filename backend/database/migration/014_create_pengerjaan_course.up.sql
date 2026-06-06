CREATE TABLE IF NOT EXISTS pengerjaan_course (
    id               SERIAL PRIMARY KEY,
    mahasiswa_id     INT              NOT NULL,
    aktivasi_sesi_id INT              NOT NULL,
    course_id        INT              NOT NULL,
    status           VARCHAR(20)      NOT NULL DEFAULT 'belum_dikerjakan',
    waktu_mulai      TIMESTAMP        NULL,
    waktu_selesai    TIMESTAMP        NULL,
    total_nilai      DOUBLE PRECISION NULL,
    CONSTRAINT uq_pengerjaan UNIQUE (mahasiswa_id, aktivasi_sesi_id, course_id),
    CONSTRAINT fk_pengerjaan_mhs FOREIGN KEY (mahasiswa_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_pengerjaan_aktivasi FOREIGN KEY (aktivasi_sesi_id) REFERENCES aktivasi_sesi(id) ON DELETE CASCADE,
    CONSTRAINT fk_pengerjaan_course FOREIGN KEY (course_id) REFERENCES course(id) ON DELETE CASCADE
);
