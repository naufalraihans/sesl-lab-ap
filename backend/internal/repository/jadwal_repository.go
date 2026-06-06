package repository

import (
	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

type JadwalRepository interface {
	Create(j *entity.Jadwal) error
	Update(j *entity.Jadwal) error
	Delete(id int) error
	FindByID(id int) (*entity.Jadwal, error)
	List() ([]entity.Jadwal, error)
	FindByKelasShift(kelasID, shift int) (*entity.Jadwal, error)
}

type jadwalRepository struct{ db *gorm.DB }

func NewJadwalRepository(db *gorm.DB) JadwalRepository { return &jadwalRepository{db: db} }

func (r *jadwalRepository) Create(j *entity.Jadwal) error { return r.db.Create(j).Error }
func (r *jadwalRepository) Update(j *entity.Jadwal) error { return r.db.Save(j).Error }
func (r *jadwalRepository) Delete(id int) error {
	return r.db.Delete(&entity.Jadwal{}, id).Error
}

func (r *jadwalRepository) FindByID(id int) (*entity.Jadwal, error) {
	var j entity.Jadwal
	if err := r.db.Preload("Kelas").First(&j, id).Error; err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *jadwalRepository) List() ([]entity.Jadwal, error) {
	var js []entity.Jadwal
	return js, r.db.Preload("Kelas").Order("kelas_id, shift").Find(&js).Error
}

func (r *jadwalRepository) FindByKelasShift(kelasID, shift int) (*entity.Jadwal, error) {
	var j entity.Jadwal
	err := r.db.Preload("Kelas").
		Where("kelas_id = ? AND shift = ?", kelasID, shift).First(&j).Error
	if err != nil {
		return nil, err
	}
	return &j, nil
}
