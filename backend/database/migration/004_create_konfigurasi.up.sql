CREATE TABLE IF NOT EXISTS konfigurasi (
    id         SERIAL PRIMARY KEY,
    "key"      VARCHAR(100) NOT NULL,
    value      TEXT         NULL,
    updated_at TIMESTAMP    NULL
);
CREATE UNIQUE INDEX uq_konfigurasi_key ON konfigurasi("key");
