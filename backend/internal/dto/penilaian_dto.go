package dto

// NilaiRequest: set nilai + feedback untuk satu jawaban (0..poin).
type NilaiRequest struct {
	JawabanID int      `json:"jawaban_id" binding:"required"`
	Nilai     float64  `json:"nilai" binding:"min=0"`
	Feedback  *string  `json:"feedback"`
}

// RekapItem: satu baris rekap jawaban mahasiswa per soal.
type RekapItem struct {
	JawabanID    int      `json:"jawaban_id"`
	MahasiswaID  int      `json:"mahasiswa_id"`
	NamaMhs      string   `json:"nama_mahasiswa"`
	NIM          string   `json:"nim"`
	SoalID       int      `json:"soal_id"`
	TeksSoal     string   `json:"teks_soal"`
	Poin         float64  `json:"poin"`
	JawabanTeks  string   `json:"jawaban_teks"`
	IsSubmitted  bool     `json:"is_submitted"`
	Nilai        *float64 `json:"nilai"`
	Feedback     *string  `json:"feedback"`
}

// RekapResponse: rekap jawaban satu aktivasi+course.
type RekapResponse struct {
	AktivasiSesiID int         `json:"aktivasi_sesi_id"`
	CourseID       int         `json:"course_id"`
	Items          []RekapItem `json:"items"`
}
