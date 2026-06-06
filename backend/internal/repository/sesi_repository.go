package repository

import (
	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

type SesiRepository interface {
	Create(s *entity.SesiPraktikum) error
	Update(s *entity.SesiPraktikum) error
	Delete(id int) error
	FindByID(id int) (*entity.SesiPraktikum, error)
	List() ([]entity.SesiPraktikum, error)
}

type sesiRepository struct{ db *gorm.DB }

func NewSesiRepository(db *gorm.DB) SesiRepository { return &sesiRepository{db: db} }

func (r *sesiRepository) Create(s *entity.SesiPraktikum) error { return r.db.Create(s).Error }
func (r *sesiRepository) Update(s *entity.SesiPraktikum) error { return r.db.Save(s).Error }
func (r *sesiRepository) Delete(id int) error {
	return r.db.Delete(&entity.SesiPraktikum{}, id).Error
}

func (r *sesiRepository) FindByID(id int) (*entity.SesiPraktikum, error) {
	var s entity.SesiPraktikum
	if err := r.db.Preload("Courses").First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *sesiRepository) List() ([]entity.SesiPraktikum, error) {
	var ss []entity.SesiPraktikum
	return ss, r.db.Preload("Courses").Order("urutan asc").Find(&ss).Error
}
