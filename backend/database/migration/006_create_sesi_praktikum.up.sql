CREATE TABLE IF NOT EXISTS sesi_praktikum (
    id               SERIAL PRIMARY KEY,
    judul_sesi       VARCHAR(200) NOT NULL,
    deskripsi        TEXT         NULL,
    urutan           INT          NOT NULL DEFAULT 0,
    is_ujian_praktik BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMP    NULL,
    updated_at       TIMESTAMP    NULL
);
