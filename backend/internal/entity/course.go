package entity

// Course: bagian dari sesi (pretest/posttest/keterampilan/ujian_praktik).
// Status buka/tutup TIDAK di sini, melainkan per kelas+shift di aktivasi_course.
type Course struct {
	ID              int         `gorm:"primaryKey;autoIncrement" json:"id"`
	SesiPraktikumID int         `gorm:"not null;uniqueIndex:idx_course_sesi_jenis" json:"sesi_praktikum_id"`
	Jenis           JenisCourse `gorm:"type:varchar(20);not null;uniqueIndex:idx_course_sesi_jenis" json:"jenis"`
	Judul           string      `gorm:"type:varchar(200)" json:"judul"`
	Deskripsi       string      `gorm:"type:text" json:"deskripsi"`
	DurasiMenit     int         `gorm:"not null;default:30" json:"durasi_menit"`

	Soal []Soal `gorm:"foreignKey:CourseID" json:"soal,omitempty"`
}

func (Course) TableName() string { return "course" }
