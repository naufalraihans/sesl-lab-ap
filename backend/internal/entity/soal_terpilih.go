package entity

// SoalTerpilih: soal hasil acak untuk setiap aktivasi sesi + course.
type SoalTerpilih struct {
	ID             int `gorm:"primaryKey;autoIncrement" json:"id"`
	AktivasiSesiID int `gorm:"not null;uniqueIndex:idx_soal_terpilih" json:"aktivasi_sesi_id"`
	CourseID       int `gorm:"not null;uniqueIndex:idx_soal_terpilih" json:"course_id"`
	SoalID         int `gorm:"not null;uniqueIndex:idx_soal_terpilih" json:"soal_id"`
	Urutan         int `json:"urutan"`

	Soal *Soal `gorm:"foreignKey:SoalID" json:"soal,omitempty"`
}

func (SoalTerpilih) TableName() string { return "soal_terpilih" }
