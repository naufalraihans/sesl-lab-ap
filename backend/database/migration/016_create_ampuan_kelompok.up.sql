CREATE TABLE IF NOT EXISTS ampuan_kelompok (
    id         SERIAL PRIMARY KEY,
    asisten_id INT         NOT NULL,
    kelas_id   INT         NOT NULL,
    kelompok   VARCHAR(50) NOT NULL,
    CONSTRAINT uq_ampuan UNIQUE (asisten_id, kelas_id, kelompok),
    CONSTRAINT fk_ampuan_asisten FOREIGN KEY (asisten_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_ampuan_kelas FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE CASCADE
);
