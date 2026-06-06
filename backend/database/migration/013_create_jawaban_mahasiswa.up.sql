CREATE TABLE IF NOT EXISTS jawaban_mahasiswa (
    id               SERIAL PRIMARY KEY,
    mahasiswa_id     INT              NOT NULL,
    soal_terpilih_id INT              NOT NULL,
    jawaban_teks     TEXT             NULL,
    is_submitted     BOOLEAN          NOT NULL DEFAULT FALSE,
    nilai            DOUBLE PRECISION NULL,
    feedback         TEXT             NULL,
    waktu_submit     TIMESTAMP        NULL,
    updated_at       TIMESTAMP        NULL,
    CONSTRAINT uq_jawaban_unik UNIQUE (mahasiswa_id, soal_terpilih_id),
    CONSTRAINT fk_jawaban_mhs FOREIGN KEY (mahasiswa_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_jawaban_st FOREIGN KEY (soal_terpilih_id) REFERENCES soal_terpilih(id) ON DELETE CASCADE
);
