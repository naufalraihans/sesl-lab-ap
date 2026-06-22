package repository

import (
	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

type SoalRepository interface {
	Create(s *entity.Soal) error
	Update(s *entity.Soal) error
	Delete(id int) error
	FindByID(id int) (*entity.Soal, error)
	ListByCourse(courseID int) ([]entity.Soal, error)
	// PoolByDifficulty: pool soal per course difilter difficulty (untuk acak pretest/posttest).
	PoolByDifficulty(courseID int, diff entity.Difficulty) ([]entity.Soal, error)
	// PoolAll: seluruh pool course (untuk keterampilan).
	PoolAll(courseID int) ([]entity.Soal, error)
	// PoolByKategori: pool soal ujian praktik per kategori.
	PoolByKategori(courseID int, kat entity.KategoriUjian) ([]entity.Soal, error)
}

type soalRepository struct{ db *gorm.DB }

func NewSoalRepository(db *gorm.DB) SoalRepository { return &soalRepository{db: db} }

func (r *soalRepository) Create(s *entity.Soal) error { return r.db.Create(s).Error }
func (r *soalRepository) Update(s *entity.Soal) error { return r.db.Save(s).Error }
func (r *soalRepository) Delete(id int) error {
	return r.db.Delete(&entity.Soal{}, id).Error
}

func (r *soalRepository) FindByID(id int) (*entity.Soal, error) {
	var s entity.Soal
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *soalRepository) ListByCourse(courseID int) ([]entity.Soal, error) {
	// Sembunyikan soal arsip dari editor pool admin (hanya untuk Rekap Jawaban).
	var ss []entity.Soal
	return ss, r.db.Where("course_id = ? AND is_arsip = false", courseID).Order("id asc").Find(&ss).Error
}

// Pool* mengecualikan soal arsip (is_arsip) agar gacha hanya memilih soal aktif.
func (r *soalRepository) PoolByDifficulty(courseID int, diff entity.Difficulty) ([]entity.Soal, error) {
	var ss []entity.Soal
	err := r.db.Where("course_id = ? AND difficulty = ? AND is_arsip = false", courseID, diff).Find(&ss).Error
	return ss, err
}

func (r *soalRepository) PoolAll(courseID int) ([]entity.Soal, error) {
	var ss []entity.Soal
	err := r.db.Where("course_id = ? AND is_arsip = false", courseID).Find(&ss).Error
	return ss, err
}

func (r *soalRepository) PoolByKategori(courseID int, kat entity.KategoriUjian) ([]entity.Soal, error) {
	var ss []entity.Soal
	err := r.db.Where("course_id = ? AND kategori_ujian = ? AND is_arsip = false", courseID, kat).Find(&ss).Error
	return ss, err
}
