package entity

// Jadwal bersifat per kelas + shift (tiap shift punya periode berbeda).
type Jadwal struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	KelasID    int    `gorm:"not null;uniqueIndex:idx_jadwal_kelas_shift" json:"kelas_id"`
	Shift      int    `gorm:"not null;uniqueIndex:idx_jadwal_kelas_shift" json:"shift"`
	Hari       string `gorm:"type:varchar(20)" json:"hari"`
	JamMulai   string `gorm:"type:time" json:"jam_mulai"`
	JamSelesai string `gorm:"type:time" json:"jam_selesai"`
	Keterangan string `gorm:"type:varchar(150)" json:"keterangan"`

	Kelas *Kelas `gorm:"foreignKey:KelasID" json:"kelas,omitempty"`
}

func (Jadwal) TableName() string { return "jadwal" }
