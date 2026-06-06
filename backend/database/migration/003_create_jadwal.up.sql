CREATE TABLE IF NOT EXISTS jadwal (
    id          SERIAL PRIMARY KEY,
    kelas_id    INT          NOT NULL,
    shift       INT          NOT NULL,
    hari        VARCHAR(20)  NULL,
    jam_mulai   TIME         NULL,
    jam_selesai TIME         NULL,
    keterangan  VARCHAR(150) NULL,
    CONSTRAINT uq_jadwal_kelas_shift UNIQUE (kelas_id, shift),
    CONSTRAINT fk_jadwal_kelas FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE CASCADE
);
