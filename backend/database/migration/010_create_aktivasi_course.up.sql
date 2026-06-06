CREATE TABLE IF NOT EXISTS aktivasi_course (
    id               SERIAL PRIMARY KEY,
    aktivasi_sesi_id INT       NOT NULL,
    course_id        INT       NOT NULL,
    is_open          BOOLEAN   NOT NULL DEFAULT FALSE,
    urutan           INT       NULL,
    opened_at        TIMESTAMP NULL,
    closed_at        TIMESTAMP NULL,
    CONSTRAINT uq_aktivasi_course UNIQUE (aktivasi_sesi_id, course_id),
    CONSTRAINT fk_aktcourse_aktivasi FOREIGN KEY (aktivasi_sesi_id) REFERENCES aktivasi_sesi(id) ON DELETE CASCADE,
    CONSTRAINT fk_aktcourse_course FOREIGN KEY (course_id) REFERENCES course(id) ON DELETE CASCADE
);
