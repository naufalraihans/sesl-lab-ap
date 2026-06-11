package repository

import (
	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

// AktivasiTxRepo menulis seluruh rangkaian aktivasi (sesi + course + soal terpilih)
// dalam SATU transaksi, sehingga kegagalan di tengah tidak meninggalkan state
// setengah jadi. Memakai *gorm.DB langsung agar tidak mengubah interface/mocks lama.
type AktivasiTxRepo struct{ db *gorm.DB }

func NewAktivasiTxRepo(db *gorm.DB) *AktivasiTxRepo { return &AktivasiTxRepo{db: db} }

// CreateActivation membuat aktivasi_sesi, lalu seluruh aktivasi_course & soal_terpilih
// (dengan aktivasi_sesi_id terisi dari sesi yang baru dibuat) secara atomik.
func (r *AktivasiTxRepo) CreateActivation(aks *entity.AktivasiSesi, courses []entity.AktivasiCourse, terpilih []entity.SoalTerpilih) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(aks).Error; err != nil {
			return err
		}
		for i := range courses {
			courses[i].AktivasiSesiID = aks.ID
		}
		if len(courses) > 0 {
			if err := tx.Create(&courses).Error; err != nil {
				return err
			}
		}
		for i := range terpilih {
			terpilih[i].AktivasiSesiID = aks.ID
		}
		if len(terpilih) > 0 {
			if err := tx.Create(&terpilih).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
