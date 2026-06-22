-- Soal arsip (hasil migrasi jawaban historis dari Firestore). Ditandai agar tidak
-- ikut diacak gacha untuk ujian baru, tapi tetap tampil di Rekap Jawaban. (update13)
ALTER TABLE soal ADD COLUMN IF NOT EXISTS is_arsip BOOLEAN NOT NULL DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_soal_arsip ON soal(is_arsip);
