package repository

import (
	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

type PedomanRepository interface {
	Create(p *entity.PedomanLaporan) error
	Update(p *entity.PedomanLaporan) error
	Delete(id int) error
	FindByID(id int) (*entity.PedomanLaporan, error)
	List() ([]entity.PedomanLaporan, error)
}

type pedomanRepository struct{ db *gorm.DB }

func NewPedomanRepository(db *gorm.DB) PedomanRepository { return &pedomanRepository{db: db} }

func (r *pedomanRepository) Create(p *entity.PedomanLaporan) error { return r.db.Create(p).Error }
func (r *pedomanRepository) Update(p *entity.PedomanLaporan) error { return r.db.Save(p).Error }
func (r *pedomanRepository) Delete(id int) error {
	return r.db.Delete(&entity.PedomanLaporan{}, id).Error
}

func (r *pedomanRepository) FindByID(id int) (*entity.PedomanLaporan, error) {
	var p entity.PedomanLaporan
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pedomanRepository) List() ([]entity.PedomanLaporan, error) {
	var ps []entity.PedomanLaporan
	return ps, r.db.Order("diunggah_pada desc").Find(&ps).Error
}
