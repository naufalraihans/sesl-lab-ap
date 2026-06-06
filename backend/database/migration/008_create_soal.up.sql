CREATE TABLE IF NOT EXISTS soal (
    id             SERIAL PRIMARY KEY,
    course_id      INT              NOT NULL,
    jenis_soal     VARCHAR(10)      NOT NULL,
    difficulty     VARCHAR(10)      NULL,
    kategori_ujian VARCHAR(15)      NULL,
    teks_soal      TEXT             NOT NULL,
    gambar_url     VARCHAR(500)     NULL,
    poin           DOUBLE PRECISION NOT NULL DEFAULT 0,
    kunci_jawaban  TEXT             NULL,
    created_at     TIMESTAMP        NULL,
    CONSTRAINT fk_soal_course FOREIGN KEY (course_id) REFERENCES course(id) ON DELETE CASCADE
);
CREATE INDEX idx_soal_course ON soal(course_id);
