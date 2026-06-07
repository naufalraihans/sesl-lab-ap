package dto

// RekapKolom merepresentasikan kolom dinamis dari hasil pivot (misal "course_1").
type RekapKolom struct {
	Key   string `json:"key"`   // identifier unik course, misal "course_1"
	Label string `json:"label"` // label human-readable, misal "Modul 1 - Pre-test"
}

// RekapMahasiswa merepresentasikan baris pada tabel pivot.
type RekapMahasiswa struct {
	NIM        string             `json:"nim"`
	Nama       string             `json:"nama"`
	Scores     map[string]float64 `json:"scores"` // Key: course_id, Value: total nilai
	TotalScore float64            `json:"total_score"`
}

// RekapKelasResponse adalah response lengkap untuk endpoint rekapitulasi.
type RekapKelasResponse struct {
	Columns []RekapKolom     `json:"columns"`
	Data    []RekapMahasiswa `json:"data"`
}
