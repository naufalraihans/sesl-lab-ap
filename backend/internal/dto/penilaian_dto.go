package dto

// NilaiRequest: set nilai + feedback untuk satu jawaban (0..poin).
type NilaiRequest struct {
	JawabanID int     `json:"jawaban_id" binding:"required"`
	Nilai     float64 `json:"nilai" binding:"min=0"`
	Feedback  *string `json:"feedback"`
}

// KeaktifanItem: set keaktifan satu pengerjaan (pretest/posttest), 0..100.
type KeaktifanItem struct {
	PengerjaanID int     `json:"pengerjaan_id" binding:"required"`
	Nilai        float64 `json:"nilai" binding:"min=0,max=100"`
}

// KeaktifanRequest: input keaktifan massal dari Rekap Nilai.
type KeaktifanRequest struct {
	Items []KeaktifanItem `json:"items" binding:"required,dive"`
}

// RekapItem: satu baris rekap jawaban mahasiswa per soal.
type RekapItem struct {
	JawabanID      int      `json:"jawaban_id"`
	MahasiswaID    int      `json:"mahasiswa_id"`
	SoalTerpilihID int      `json:"soal_terpilih_id"`
	NamaMhs        string   `json:"nama_mahasiswa"`
	NIM            string   `json:"nim"`
	SoalID         int      `json:"soal_id"`
	TeksSoal       string   `json:"teks_soal"`
	JenisSoal      string   `json:"jenis_soal"`
	Poin           float64  `json:"poin"`
	Urutan         int      `json:"urutan"`
	JawabanTeks    string   `json:"jawaban_teks"`
	IsSubmitted    bool     `json:"is_submitted"`
	Nilai          *float64 `json:"nilai"`
	Feedback       *string  `json:"feedback"`
}

// RekapResponse: rekap jawaban satu aktivasi+course.
type RekapResponse struct {
	AktivasiSesiID int         `json:"aktivasi_sesi_id"`
	CourseID       int         `json:"course_id"`
	Items          []RekapItem `json:"items"`
}

// AIGradeOneRequest: nilai SATU jawaban dengan AI (dipanggil berulang oleh frontend).
type AIGradeOneRequest struct {
	JawabanID int `json:"jawaban_id" binding:"required"`
}

// AIGradeOneResponse: hasil penilaian AI untuk satu jawaban.
type AIGradeOneResponse struct {
	JawabanID int     `json:"jawaban_id"`
	Nilai     float64 `json:"nilai"`
	Feedback  string  `json:"feedback"`
}

// AIGradeTargetsResponse: daftar jawaban_id yang perlu dinilai AI untuk satu course.
type AIGradeTargetsResponse struct {
	JawabanIDs []int `json:"jawaban_ids"`
	Total      int   `json:"total"`
}
