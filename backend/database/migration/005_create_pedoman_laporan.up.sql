CREATE TABLE IF NOT EXISTS pedoman_laporan (
    id            SERIAL PRIMARY KEY,
    nama_dokumen  VARCHAR(200) NOT NULL,
    file_url      VARCHAR(500) NOT NULL,
    diunggah_pada TIMESTAMP    NULL
);
