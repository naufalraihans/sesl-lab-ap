package repository

import (
	"sync"
	"time"

	"lab-ap/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KonfigurasiRepository interface {
	Get(key string) (*entity.Konfigurasi, error)
	Set(key, value string) error
	All() ([]entity.Konfigurasi, error)
}

type konfigurasiRepository struct {
	db    *gorm.DB
	cache sync.Map
}

type cacheEntry struct {
	value     *entity.Konfigurasi
	expiresAt time.Time
}

func NewKonfigurasiRepository(db *gorm.DB) KonfigurasiRepository {
	return &konfigurasiRepository{db: db}
}

func (r *konfigurasiRepository) Get(key string) (*entity.Konfigurasi, error) {
	if val, ok := r.cache.Load(key); ok {
		entry := val.(cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.value, nil
		}
		r.cache.Delete(key) // Clean up expired cache
	}

	var k entity.Konfigurasi
	if err := r.db.Where("\"key\" = ?", key).First(&k).Error; err != nil {
		return nil, err
	}
	
	// Cache the result for 5 minutes
	r.cache.Store(key, cacheEntry{value: &k, expiresAt: time.Now().Add(5 * time.Minute)})
	return &k, nil
}

// Set melakukan upsert berdasarkan key.
func (r *konfigurasiRepository) Set(key, value string) error {
	r.cache.Delete(key) // Invalidate cache on update
	
	k := entity.Konfigurasi{Key: key, Value: value, UpdatedAt: time.Now()}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&k).Error
}

func (r *konfigurasiRepository) All() ([]entity.Konfigurasi, error) {
	var ks []entity.Konfigurasi
	return ks, r.db.Find(&ks).Error
}
