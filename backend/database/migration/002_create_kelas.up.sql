CREATE TABLE IF NOT EXISTS kelas (
    id                SERIAL PRIMARY KEY,
    nama_kelas        VARCHAR(100) NOT NULL,
    is_register_open  BOOLEAN      NOT NULL DEFAULT FALSE
);
