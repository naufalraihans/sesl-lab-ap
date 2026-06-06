package repository

import (
	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

type AmpuanRepository interface {
	List() ([]entity.AmpuanKelompok, error)
	ListByKelas(kelasID int) ([]entity.AmpuanKelompok, error)
	Create(a *entity.AmpuanKelompok) error
	Delete(id int) error
}

type ampuanRepository struct{ db *gorm.DB }

func NewAmpuanRepository(db *gorm.DB) AmpuanRepository { return &ampuanRepository{db: db} }

func (r *ampuanRepository) List() ([]entity.AmpuanKelompok, error) {
	var list []entity.AmpuanKelompok
	err := r.db.Preload("Asisten").Preload("Kelas").Order("kelas_id, kelompok").Find(&list).Error
	return list, err
}

func (r *ampuanRepository) ListByKelas(kelasID int) ([]entity.AmpuanKelompok, error) {
	var list []entity.AmpuanKelompok
	err := r.db.Preload("Asisten").Where("kelas_id = ?", kelasID).Order("kelompok").Find(&list).Error
	return list, err
}

func (r *ampuanRepository) Create(a *entity.AmpuanKelompok) error { return r.db.Create(a).Error }

func (r *ampuanRepository) Delete(id int) error {
	return r.db.Delete(&entity.AmpuanKelompok{}, id).Error
}
