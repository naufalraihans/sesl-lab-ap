package repository

import (
	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

type KelasRepository interface {
	Create(k *entity.Kelas) error
	Update(k *entity.Kelas) error
	Delete(id int) error
	FindByID(id int) (*entity.Kelas, error)
	List() ([]entity.Kelas, error)
	SetRegisterOpen(id int, open bool) error
}

type kelasRepository struct{ db *gorm.DB }

func NewKelasRepository(db *gorm.DB) KelasRepository { return &kelasRepository{db: db} }

func (r *kelasRepository) Create(k *entity.Kelas) error { return r.db.Create(k).Error }
func (r *kelasRepository) Update(k *entity.Kelas) error { return r.db.Save(k).Error }
func (r *kelasRepository) Delete(id int) error {
	return r.db.Delete(&entity.Kelas{}, id).Error
}

func (r *kelasRepository) FindByID(id int) (*entity.Kelas, error) {
	var k entity.Kelas
	if err := r.db.First(&k, id).Error; err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *kelasRepository) List() ([]entity.Kelas, error) {
	var ks []entity.Kelas
	return ks, r.db.Order("nama_kelas asc").Find(&ks).Error
}

func (r *kelasRepository) SetRegisterOpen(id int, open bool) error {
	return r.db.Model(&entity.Kelas{}).Where("id = ?", id).
		Update("is_register_open", open).Error
}
