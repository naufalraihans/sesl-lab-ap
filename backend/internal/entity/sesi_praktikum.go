package entity

import "time"

// SesiPraktikum: satu modul/pertemuan (mis. "Modul 1 - Pengenalan Dasar Bahasa C").
type SesiPraktikum struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	JudulSesi      string    `gorm:"type:varchar(200);not null" json:"judul_sesi"`
	Deskripsi      string    `gorm:"type:text" json:"deskripsi"`
	Urutan         int       `gorm:"not null" json:"urutan"`
	IsUjianPraktik bool      `gorm:"default:false" json:"is_ujian_praktik"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Courses []Course `gorm:"foreignKey:SesiPraktikumID" json:"courses,omitempty"`
}

func (SesiPraktikum) TableName() string { return "sesi_praktikum" }
