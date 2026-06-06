CREATE TABLE IF NOT EXISTS peserta_susulan (
    id               SERIAL PRIMARY KEY,
    aktivasi_sesi_id INT          NOT NULL,
    mahasiswa_id     INT          NOT NULL,
    alasan           VARCHAR(255) NULL,
    created_at       TIMESTAMP    NULL,
    CONSTRAINT uq_susulan UNIQUE (aktivasi_sesi_id, mahasiswa_id),
    CONSTRAINT fk_susulan_aktivasi FOREIGN KEY (aktivasi_sesi_id) REFERENCES aktivasi_sesi(id) ON DELETE CASCADE,
    CONSTRAINT fk_susulan_mhs FOREIGN KEY (mahasiswa_id) REFERENCES users(id) ON DELETE CASCADE
);
