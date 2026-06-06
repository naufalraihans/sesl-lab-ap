package entity

import "time"

// Konfigurasi key-value global (mis. gdrive_jadwal_url, modul_file_url).
type Konfigurasi struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string    `gorm:"column:key;type:varchar(100);uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Konfigurasi) TableName() string { return "konfigurasi" }

// Key konfigurasi yang dikenal sistem.
const (
	KeyGDriveJadwalURL = "gdrive_jadwal_url"
	KeyJadwalMode      = "jadwal_mode" // "gdrive" | "internal"
	KeyModulFileURL    = "modul_file_url"
)
