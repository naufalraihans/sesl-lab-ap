package entity

import "time"

// PedomanLaporan: file template/pedoman yang dinamis (muncul sebagai tombol download).
type PedomanLaporan struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaDokumen  string    `gorm:"type:varchar(200);not null" json:"nama_dokumen"`
	FileURL      string    `gorm:"type:varchar(500);not null" json:"file_url"`
	DiunggahPada time.Time `json:"diunggah_pada"`
}

func (PedomanLaporan) TableName() string { return "pedoman_laporan" }
