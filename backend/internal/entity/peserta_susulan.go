package entity

import "time"

// PesertaSusulan: mahasiswa yang diberi akses ke aktivasi yang BUKAN kelas/shift aslinya.
type PesertaSusulan struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AktivasiSesiID int       `gorm:"not null;uniqueIndex:idx_susulan" json:"aktivasi_sesi_id"`
	MahasiswaID    int       `gorm:"not null;uniqueIndex:idx_susulan" json:"mahasiswa_id"`
	Alasan         string    `gorm:"type:varchar(255)" json:"alasan"`
	CreatedAt      time.Time `json:"created_at"`

	Mahasiswa *User `gorm:"foreignKey:MahasiswaID" json:"mahasiswa,omitempty"`
}

func (PesertaSusulan) TableName() string { return "peserta_susulan" }
