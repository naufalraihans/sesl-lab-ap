package repository

import (
	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

// PenilaianTxRepo menyediakan operasi penilaian yang ATOMIK (dalam transaksi)
// dan bebas N+1. Sengaja terpisah dari JawabanRepository agar tidak mengubah
// interface/mocks yang sudah ada — komponen ini memakai *gorm.DB langsung.
type PenilaianTxRepo struct{ db *gorm.DB }

func NewPenilaianTxRepo(db *gorm.DB) *PenilaianTxRepo { return &PenilaianTxRepo{db: db} }

// PengerjaanKey mengidentifikasi satu unit pengerjaan (mahasiswa+aktivasi+course).
type PengerjaanKey struct {
	MahasiswaID    int
	AktivasiSesiID int
	CourseID       int
}

// FindJawabanByIDs mengambil banyak jawaban sekaligus (1 query, anti N+1)
// lengkap dengan relasi SoalTerpilih + Soal.
func (r *PenilaianTxRepo) FindJawabanByIDs(ids []int) ([]entity.JawabanMahasiswa, error) {
	if len(ids) == 0 {
		return []entity.JawabanMahasiswa{}, nil
	}
	var js []entity.JawabanMahasiswa
	err := r.db.Preload("SoalTerpilih").Preload("SoalTerpilih.Soal").
		Where("id IN ?", ids).Find(&js).Error
	return js, err
}

// SetNilaiAndRecalc menyimpan nilai+feedback satu jawaban lalu menghitung ulang
// total_nilai pengerjaan_course, seluruhnya dalam SATU transaksi.
func (r *PenilaianTxRepo) SetNilaiAndRecalc(jawabanID int, nilai float64, feedback *string, key PengerjaanKey) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.JawabanMahasiswa{}).
			Where("id = ?", jawabanID).
			Updates(map[string]interface{}{"nilai": nilai, "feedback": feedback}).Error; err != nil {
			return err
		}
		return recalcTotalTx(tx, key)
	})
}

// BulkResetAndRecalc mereset nilai+feedback sekumpulan jawaban lalu recalc total
// untuk tiap pengerjaan terdampak, dalam satu transaksi.
func (r *PenilaianTxRepo) BulkResetAndRecalc(jawabanIDs []int) error {
	if len(jawabanIDs) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		keys, err := affectedKeysTx(tx, jawabanIDs)
		if err != nil {
			return err
		}
		if err := tx.Model(&entity.JawabanMahasiswa{}).
			Where("id IN ?", jawabanIDs).
			Updates(map[string]interface{}{"nilai": nil, "feedback": nil}).Error; err != nil {
			return err
		}
		return recalcKeysTx(tx, keys)
	})
}

// BulkDeleteAndRecalc menghapus sekumpulan jawaban lalu recalc total untuk tiap
// pengerjaan terdampak, dalam satu transaksi. Key diambil SEBELUM delete.
func (r *PenilaianTxRepo) BulkDeleteAndRecalc(jawabanIDs []int) error {
	if len(jawabanIDs) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		keys, err := affectedKeysTx(tx, jawabanIDs)
		if err != nil {
			return err
		}
		if err := tx.Where("id IN ?", jawabanIDs).Delete(&entity.JawabanMahasiswa{}).Error; err != nil {
			return err
		}
		return recalcKeysTx(tx, keys)
	})
}

// affectedKeysTx mengambil daftar unik (mahasiswa, aktivasi, course) untuk
// sekumpulan jawaban dalam SATU query join (anti N+1).
func affectedKeysTx(tx *gorm.DB, jawabanIDs []int) ([]PengerjaanKey, error) {
	var keys []PengerjaanKey
	err := tx.Table("jawaban_mahasiswa AS j").
		Select("DISTINCT j.mahasiswa_id AS mahasiswa_id, st.aktivasi_sesi_id AS aktivasi_sesi_id, st.course_id AS course_id").
		Joins("JOIN soal_terpilih st ON st.id = j.soal_terpilih_id").
		Where("j.id IN ?", jawabanIDs).
		Scan(&keys).Error
	return keys, err
}

func recalcKeysTx(tx *gorm.DB, keys []PengerjaanKey) error {
	for _, k := range keys {
		if err := recalcTotalTx(tx, k); err != nil {
			return err
		}
	}
	return nil
}

// recalcTotalTx menghitung ulang total_nilai (SUM nilai jawaban) untuk satu
// pengerjaan_course dan mengupsert nilainya, dalam transaksi tx.
func recalcTotalTx(tx *gorm.DB, key PengerjaanKey) error {
	var stIDs []int
	if err := tx.Model(&entity.SoalTerpilih{}).
		Where("aktivasi_sesi_id = ? AND course_id = ?", key.AktivasiSesiID, key.CourseID).
		Pluck("id", &stIDs).Error; err != nil {
		return err
	}

	var total float64
	if len(stIDs) > 0 {
		var t *float64
		if err := tx.Model(&entity.JawabanMahasiswa{}).
			Where("mahasiswa_id = ? AND soal_terpilih_id IN ?", key.MahasiswaID, stIDs).
			Select("COALESCE(SUM(nilai),0)").Scan(&t).Error; err != nil {
			return err
		}
		if t != nil {
			total = *t
		}
	}

	// Upsert total ke pengerjaan_course (update jika ada, create jika belum).
	res := tx.Model(&entity.PengerjaanCourse{}).
		Where("mahasiswa_id = ? AND aktivasi_sesi_id = ? AND course_id = ?",
			key.MahasiswaID, key.AktivasiSesiID, key.CourseID).
		Update("total_nilai", total)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return tx.Create(&entity.PengerjaanCourse{
			MahasiswaID:    key.MahasiswaID,
			AktivasiSesiID: key.AktivasiSesiID,
			CourseID:       key.CourseID,
			Status:         entity.StatusBelum,
			TotalNilai:     &total,
		}).Error
	}
	return nil
}
