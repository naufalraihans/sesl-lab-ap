CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    role            VARCHAR(10)  NOT NULL,
    nim             VARCHAR(32)  NOT NULL,
    nama            VARCHAR(150) NOT NULL,
    password_hash   VARCHAR(255) NULL,
    is_registered   BOOLEAN      NOT NULL DEFAULT FALSE,
    kelas_id        INT          NULL,
    shift           INT          NULL,
    foto_url        VARCHAR(500) NULL,
    nomor_hp        VARCHAR(30)  NULL,
    medsos_link     VARCHAR(500) NULL,
    last_login_at   TIMESTAMP    NULL,
    created_at      TIMESTAMP    NULL,
    updated_at      TIMESTAMP    NULL
);
CREATE UNIQUE INDEX uq_users_nim ON users(nim);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_kelas ON users(kelas_id);
