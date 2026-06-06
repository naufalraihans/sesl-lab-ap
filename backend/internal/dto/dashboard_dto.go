package dto

import "lab-ap/internal/repository"

// StatistikResponse: ringkasan dashboard admin.
type StatistikResponse struct {
	OnlineSekarang     int                          `json:"online_sekarang"`
	TotalMahasiswa     int64                        `json:"total_mahasiswa"`
	TotalAsisten       int64                        `json:"total_asisten"`
	SudahRegister      int64                        `json:"sudah_register"`
	BelumRegister      int64                        `json:"belum_register"`
	PerKelasShift      []repository.KelasShiftCount `json:"per_kelas_shift"`
	SesiAktif          []SesiAktifInfo              `json:"sesi_aktif"`
}

// SesiAktifInfo: info aktivasi yang sedang aktif + progress per course.
type SesiAktifInfo struct {
	AktivasiSesiID int                  `json:"aktivasi_sesi_id"`
	JudulSesi      string               `json:"judul_sesi"`
	NamaKelas      string               `json:"nama_kelas"`
	Shift          int                  `json:"shift"`
	Courses        []CourseProgressInfo `json:"courses"`
}

// CourseProgressInfo: progress satu course dalam aktivasi.
type CourseProgressInfo struct {
	CourseID int    `json:"course_id"`
	Jenis    string `json:"jenis"`
	IsOpen   bool   `json:"is_open"`
	Selesai  int64  `json:"selesai"`
	Sedang   int64  `json:"sedang"`
	Belum    int64  `json:"belum"`
}

// OnlineCountResponse: jumlah user online real-time.
type OnlineCountResponse struct {
	Total int `json:"total"`
	Admin int `json:"admin"`
	User  int `json:"user"`
}
