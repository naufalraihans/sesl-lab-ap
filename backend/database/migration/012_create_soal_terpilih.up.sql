CREATE TABLE IF NOT EXISTS soal_terpilih (
    id               SERIAL PRIMARY KEY,
    aktivasi_sesi_id INT NOT NULL,
    course_id        INT NOT NULL,
    soal_id          INT NOT NULL,
    urutan           INT NULL,
    CONSTRAINT uq_soal_terpilih UNIQUE (aktivasi_sesi_id, course_id, soal_id),
    CONSTRAINT fk_st_aktivasi FOREIGN KEY (aktivasi_sesi_id) REFERENCES aktivasi_sesi(id) ON DELETE CASCADE,
    CONSTRAINT fk_st_course FOREIGN KEY (course_id) REFERENCES course(id) ON DELETE CASCADE,
    CONSTRAINT fk_st_soal FOREIGN KEY (soal_id) REFERENCES soal(id) ON DELETE CASCADE
);
