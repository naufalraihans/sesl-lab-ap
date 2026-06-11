package dto

type AmpuanRequest struct {
	AsistenID int    `json:"asisten_id" binding:"required"`
	KelasID   int    `json:"kelas_id" binding:"required"`
	Kelompok  string `json:"kelompok" binding:"required"`
}

// PublicMahasiswaItem: data mahasiswa MINIMAL untuk endpoint publik
// (tanpa PII seperti nomor_hp, medsos, last_login, status registrasi).
type PublicMahasiswaItem struct {
	ID       int     `json:"id"`
	NIM      string  `json:"nim"`
	Nama     string  `json:"nama"`
	Shift    *int    `json:"shift"`
	Kelompok *string `json:"kelompok"`
}
