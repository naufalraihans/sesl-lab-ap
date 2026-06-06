package entity

type AmpuanKelompok struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	AsistenID int    `gorm:"not null;uniqueIndex:uq_ampuan" json:"asisten_id"`
	KelasID   int    `gorm:"not null;uniqueIndex:uq_ampuan" json:"kelas_id"`
	Kelompok  string `gorm:"type:varchar(50);not null;uniqueIndex:uq_ampuan" json:"kelompok"`

	Asisten *User  `gorm:"foreignKey:AsistenID" json:"asisten,omitempty"`
	Kelas   *Kelas `gorm:"foreignKey:KelasID" json:"kelas,omitempty"`
}

func (AmpuanKelompok) TableName() string { return "ampuan_kelompok" }
