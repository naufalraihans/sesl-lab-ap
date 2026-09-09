-- 028: Hapus is_register_open — auth roster-gated (ala project_mikon), NIM di roster = boleh register.
ALTER TABLE kelas DROP COLUMN IF EXISTS is_register_open;
