package entity

import "time"

// AktivasiSesi: admin mengaktifkan sesi untuk kelas + shift tertentu.
type AktivasiSesi struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SesiPraktikumID int       `gorm:"not null;uniqueIndex:idx_aktivasi_sesi_kelas_shift" json:"sesi_praktikum_id"`
	KelasID         int       `gorm:"not null;uniqueIndex:idx_aktivasi_sesi_kelas_shift" json:"kelas_id"`
	Shift           int       `gorm:"not null;uniqueIndex:idx_aktivasi_sesi_kelas_shift" json:"shift"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	ActivatedAt     time.Time `json:"activated_at"`

	Sesi            *SesiPraktikum   `gorm:"foreignKey:SesiPraktikumID" json:"sesi,omitempty"`
	Kelas           *Kelas           `gorm:"foreignKey:KelasID" json:"kelas,omitempty"`
	AktivasiCourses []AktivasiCourse `gorm:"foreignKey:AktivasiSesiID" json:"aktivasi_courses,omitempty"`
}

func (AktivasiSesi) TableName() string { return "aktivasi_sesi" }
