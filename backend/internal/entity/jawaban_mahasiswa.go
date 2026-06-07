package entity

import "time"

// JawabanMahasiswa: jawaban per soal (auto-save berkala; auto-submit saat timer habis / akses ditutup).
type JawabanMahasiswa struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	MahasiswaID    int        `gorm:"not null;uniqueIndex:idx_jawaban_unik" json:"mahasiswa_id"`
	SoalTerpilihID int        `gorm:"not null;uniqueIndex:idx_jawaban_unik" json:"soal_terpilih_id"`
	JawabanTeks    string     `gorm:"type:longtext" json:"jawaban_teks"`
	IsSubmitted    bool       `gorm:"default:false" json:"is_submitted"`
	Nilai          *float64   `json:"nilai"`
	Feedback       *string    `gorm:"type:text" json:"feedback"`
	WaktuSubmit    *time.Time `json:"waktu_submit"`
	UpdatedAt      time.Time  `json:"updated_at"`

	SoalTerpilih *SoalTerpilih `gorm:"foreignKey:SoalTerpilihID" json:"soal_terpilih,omitempty"`
	Mahasiswa    *User         `gorm:"foreignKey:MahasiswaID" json:"mahasiswa,omitempty"`
}

func (JawabanMahasiswa) TableName() string { return "jawaban_mahasiswa" }
