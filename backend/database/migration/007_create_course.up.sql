CREATE TABLE IF NOT EXISTS course (
    id                SERIAL PRIMARY KEY,
    sesi_praktikum_id INT          NOT NULL,
    jenis             VARCHAR(20)  NOT NULL,
    judul             VARCHAR(200) NULL,
    deskripsi         TEXT         NULL,
    durasi_menit      INT          NOT NULL DEFAULT 30,
    CONSTRAINT uq_course_sesi_jenis UNIQUE (sesi_praktikum_id, jenis),
    CONSTRAINT fk_course_sesi FOREIGN KEY (sesi_praktikum_id) REFERENCES sesi_praktikum(id) ON DELETE CASCADE
);
