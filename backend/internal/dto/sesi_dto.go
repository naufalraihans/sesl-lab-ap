package dto

// SesiRequest: buat/ubah sesi praktikum.
type SesiRequest struct {
	JudulSesi      string `json:"judul_sesi" binding:"required"`
	Deskripsi      string `json:"deskripsi"`
	Urutan         int    `json:"urutan"`
	IsUjianPraktik bool   `json:"is_ujian_praktik"`
}

// CourseRequest: buat/ubah course dalam sesi.
type CourseRequest struct {
	Jenis       string `json:"jenis" binding:"required,oneof=pretest posttest keterampilan ujian_praktik"`
	Judul       string `json:"judul"`
	Deskripsi   string `json:"deskripsi"`
	DurasiMenit int    `json:"durasi_menit" binding:"required,min=1"`
	// Kuota gacha per difficulty (pretest/posttest). Kosong = default per jenis.
	KuotaEasy   *int `json:"kuota_easy" binding:"omitempty,min=0"`
	KuotaMedium *int `json:"kuota_medium" binding:"omitempty,min=0"`
	KuotaHard   *int `json:"kuota_hard" binding:"omitempty,min=0"`
}
