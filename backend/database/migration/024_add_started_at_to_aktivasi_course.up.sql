-- Anchor timer GLOBAL ujian: di-set sekali oleh peserta PERTAMA yang mulai.
-- Deadline semua peserta course = started_at + course.durasi_menit (kecuali susulan).
ALTER TABLE aktivasi_course ADD COLUMN IF NOT EXISTS started_at TIMESTAMP NULL;
