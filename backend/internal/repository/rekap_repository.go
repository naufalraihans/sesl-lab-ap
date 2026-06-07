package repository

import (
	"gorm.io/gorm"
)

type RekapRow struct {
	MahasiswaID int      `gorm:"column:mahasiswa_id"`
	NIM         string   `gorm:"column:nim"`
	Nama        string   `gorm:"column:nama"`
	CourseID    *int     `gorm:"column:course_id"`
	CourseJenis *string  `gorm:"column:course_jenis"`
	SesiJudul   *string  `gorm:"column:sesi_judul"`
	SesiUrutan  *int     `gorm:"column:sesi_urutan"`
	TotalNilai  *float64 `gorm:"column:total_nilai"`
}

type RekapRepository interface {
	GetRekapByKelas(kelasID int) ([]RekapRow, error)
}

type rekapRepository struct {
	db *gorm.DB
}

func NewRekapRepository(db *gorm.DB) RekapRepository {
	return &rekapRepository{db: db}
}

func (r *rekapRepository) GetRekapByKelas(kelasID int) ([]RekapRow, error) {
	var rows []RekapRow
	query := `
		SELECT 
			u.id as mahasiswa_id, u.nim, u.nama,
			c.id as course_id, c.jenis as course_jenis,
			s.judul_sesi as sesi_judul, s.urutan as sesi_urutan,
			p.total_nilai
		FROM users u
		LEFT JOIN pengerjaan_course p ON p.mahasiswa_id = u.id
		LEFT JOIN course c ON p.course_id = c.id
		LEFT JOIN sesi_praktikum s ON c.sesi_praktikum_id = s.id
		WHERE u.kelas_id = ? AND u.role = 'user'
		ORDER BY u.nim ASC, s.urutan ASC, c.id ASC
	`
	if err := r.db.Raw(query, kelasID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
