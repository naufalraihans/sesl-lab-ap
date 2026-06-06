package entity

import "time"

// PengerjaanCourse: tracking status pengerjaan mahasiswa per course per aktivasi.
type PengerjaanCourse struct {
	ID             int              `gorm:"primaryKey;autoIncrement" json:"id"`
	MahasiswaID    int              `gorm:"not null;uniqueIndex:idx_pengerjaan" json:"mahasiswa_id"`
	AktivasiSesiID int              `gorm:"not null;uniqueIndex:idx_pengerjaan" json:"aktivasi_sesi_id"`
	CourseID       int              `gorm:"not null;uniqueIndex:idx_pengerjaan" json:"course_id"`
	Status         StatusPengerjaan `gorm:"type:varchar(20);default:'belum_dikerjakan'" json:"status"`
	WaktuMulai     *time.Time       `json:"waktu_mulai"`
	WaktuSelesai   *time.Time       `json:"waktu_selesai"`
	// TotalNilai: CACHED/DERIVED akumulasi SUM(jawaban.nilai). Di-recalc tiap update nilai.
	TotalNilai *float64 `json:"total_nilai"`
}

func (PengerjaanCourse) TableName() string { return "pengerjaan_course" }
