package dto

import "time"

// AutoSaveRequest: simpan jawaban berkala (per soal).
type AutoSaveRequest struct {
	SoalTerpilihID int    `json:"soal_terpilih_id" binding:"required"`
	JawabanTeks    string `json:"jawaban_teks"`
}

// SubmitRequest: submit manual seluruh course.
type SubmitRequest struct {
	AktivasiSesiID int `json:"aktivasi_sesi_id" binding:"required"`
	CourseID       int `json:"course_id" binding:"required"`
}

// MulaiCourseRequest: tandai mulai mengerjakan course (set waktu_mulai sekali).
type MulaiCourseRequest struct {
	AktivasiSesiID int     `json:"aktivasi_sesi_id" binding:"required"`
	CourseID       int     `json:"course_id" binding:"required"`
	Token          *string `json:"token"`
}

// SoalTampilResponse: soal yang ditampilkan ke mahasiswa (tanpa kunci jawaban).
type SoalTampilResponse struct {
	SoalTerpilihID int     `json:"soal_terpilih_id"`
	Urutan         int     `json:"urutan"`
	JenisSoal      string  `json:"jenis_soal"`
	KategoriUjian  *string `json:"kategori_ujian,omitempty"`
	TeksSoal       string  `json:"teks_soal"`
	GambarURL      *string `json:"gambar_url,omitempty"`
	Poin           float64 `json:"poin"`
	JawabanTeks    string  `json:"jawaban_teks"`
	IsSubmitted    bool    `json:"is_submitted"`
}

// RuangCourseResponse: data lengkap ruang pengerjaan course.
type RuangCourseResponse struct {
	AktivasiSesiID int                  `json:"aktivasi_sesi_id"`
	CourseID       int                  `json:"course_id"`
	Jenis          string               `json:"jenis"`
	DurasiMenit    int                  `json:"durasi_menit"`
	WaktuMulai     *time.Time           `json:"waktu_mulai"`
	Deadline       *time.Time           `json:"deadline"`
	Status         string               `json:"status"`
	IsOpen         bool                 `json:"is_open"`
	RequireToken   bool                 `json:"require_token"`
	Soal           []SoalTampilResponse `json:"soal"`
}
