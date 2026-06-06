CREATE TABLE IF NOT EXISTS aktivasi_sesi (
    id                SERIAL PRIMARY KEY,
    sesi_praktikum_id INT       NOT NULL,
    kelas_id          INT       NOT NULL,
    shift             INT       NOT NULL,
    is_active         BOOLEAN   NOT NULL DEFAULT TRUE,
    activated_at      TIMESTAMP NULL,
    CONSTRAINT uq_aktivasi_sesi UNIQUE (sesi_praktikum_id, kelas_id, shift),
    CONSTRAINT fk_aktivasi_sesi_sesi FOREIGN KEY (sesi_praktikum_id) REFERENCES sesi_praktikum(id) ON DELETE CASCADE,
    CONSTRAINT fk_aktivasi_sesi_kelas FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE CASCADE
);
