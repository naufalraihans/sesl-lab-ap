package repository

import (
	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

type SoalTerpilihRepository interface {
	BulkCreate(items []entity.SoalTerpilih) error
	ListByAktivasiCourse(aktivasiSesiID, courseID int) ([]entity.SoalTerpilih, error)
	ExistsForAktivasiCourse(aktivasiSesiID, courseID int) (bool, error)
	FindByID(id int) (*entity.SoalTerpilih, error)
}

type soalTerpilihRepository struct{ db *gorm.DB }

func NewSoalTerpilihRepository(db *gorm.DB) SoalTerpilihRepository {
	return &soalTerpilihRepository{db: db}
}

func (r *soalTerpilihRepository) BulkCreate(items []entity.SoalTerpilih) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}

func (r *soalTerpilihRepository) ListByAktivasiCourse(aktivasiSesiID, courseID int) ([]entity.SoalTerpilih, error) {
	var sts []entity.SoalTerpilih
	err := r.db.Preload("Soal").
		Where("aktivasi_sesi_id = ? AND course_id = ?", aktivasiSesiID, courseID).
		Order("urutan asc").Find(&sts).Error
	return sts, err
}

func (r *soalTerpilihRepository) ExistsForAktivasiCourse(aktivasiSesiID, courseID int) (bool, error) {
	var n int64
	err := r.db.Model(&entity.SoalTerpilih{}).
		Where("aktivasi_sesi_id = ? AND course_id = ?", aktivasiSesiID, courseID).
		Count(&n).Error
	return n > 0, err
}

func (r *soalTerpilihRepository) FindByID(id int) (*entity.SoalTerpilih, error) {
	var st entity.SoalTerpilih
	if err := r.db.Preload("Soal").First(&st, id).Error; err != nil {
		return nil, err
	}
	return &st, nil
}
