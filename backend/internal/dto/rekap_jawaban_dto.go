package dto

import "time"

// RekapJawabanItem merepresentasikan satu baris jawaban di tabel Rekap Jawaban Global
type RekapJawabanItem struct {
	JawabanID       int        `json:"jawaban_id"`
	NIM             string     `json:"nim"`
	NamaMahasiswa   string     `json:"nama_mahasiswa"`
	KelasID         int        `json:"kelas_id"`
	NamaKelas       string     `json:"nama_kelas"`
	SesiPraktikumID int        `json:"sesi_praktikum_id"`
	JudulSesi       string     `json:"judul_sesi"`
	CourseID        int        `json:"course_id"`
	JudulCourse     string     `json:"judul_course"`
	JenisCourse     string     `json:"jenis_course"`
	JenisSoal       string     `json:"jenis_soal"`
	TeksSoal        string     `json:"teks_soal"`
	PoinMaksimal    float64    `json:"poin_maksimal"`
	JawabanTeks     string     `json:"jawaban_teks"`
	IsSubmitted     bool       `json:"is_submitted"`
	WaktuSubmit     *time.Time `json:"waktu_submit"`
	Nilai           *float64   `json:"nilai"`
	Feedback        *string    `json:"feedback"`
}

// RekapJawabanResponse adalah response untuk tabel Rekap Jawaban
type RekapJawabanResponse struct {
	Items []RekapJawabanItem `json:"items"`
	Total int64              `json:"total"`
}

// BulkActionRequest digunakan untuk menghapus atau mereset nilai masal
type BulkActionRequest struct {
	Action     string `json:"action" binding:"required,oneof=delete reset_nilai"`
	JawabanIDs []int  `json:"jawaban_ids" binding:"required,min=1"`
}
