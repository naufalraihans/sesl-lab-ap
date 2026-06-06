package entity

import "time"

// AktivasiCourse: status buka/tutup tiap course dalam sebuah aktivasi (per kelas+shift).
// Hanya course yang punya baris di sini yang DIPAKAI pada aktivasi tsb (hasil gacha).
type AktivasiCourse struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	AktivasiSesiID int        `gorm:"not null;uniqueIndex:idx_aktivasi_course" json:"aktivasi_sesi_id"`
	CourseID       int        `gorm:"not null;uniqueIndex:idx_aktivasi_course" json:"course_id"`
	IsOpen         bool       `gorm:"default:false" json:"is_open"`
	Urutan         int        `json:"urutan"`
	OpenedAt       *time.Time `json:"opened_at"`
	ClosedAt       *time.Time `json:"closed_at"`

	Course *Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (AktivasiCourse) TableName() string { return "aktivasi_course" }
