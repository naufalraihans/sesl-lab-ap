package repository

import (
	"errors"
	"time"

	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

// AktivasiRepository mengelola aktivasi_sesi + aktivasi_course + peserta_susulan.
type AktivasiRepository interface {
	// aktivasi_sesi
	CreateSesi(a *entity.AktivasiSesi) error
	UpdateSesi(a *entity.AktivasiSesi) error
	FindSesiByID(id int) (*entity.AktivasiSesi, error)
	FindSesiByComposite(sesiID, kelasID, shift int, gelombang *int) (*entity.AktivasiSesi, error)
	ListSesi() ([]entity.AktivasiSesi, error)
	ListActiveSesi() ([]entity.AktivasiSesi, error)
	DeleteSesi(id int) error

	// aktivasi_course
	CreateCourse(ac *entity.AktivasiCourse) error
	UpdateCourse(ac *entity.AktivasiCourse) error
	FindCourseByID(id int) (*entity.AktivasiCourse, error)
	FindCourse(aktivasiSesiID, courseID int) (*entity.AktivasiCourse, error)
	ListCoursesByAktivasi(aktivasiSesiID int) ([]entity.AktivasiCourse, error)
	// SetCourseStartedAtIfNull menetapkan anchor timer global HANYA bila belum di-set
	// (atomic: peserta pertama menang). No-op bila started_at sudah terisi.
	SetCourseStartedAtIfNull(aktivasiCourseID int, t time.Time) error

	// peserta_susulan
	AddSusulan(p *entity.PesertaSusulan) error
	RemoveSusulan(aktivasiSesiID, mahasiswaID int) error
	ListSusulan(aktivasiSesiID int) ([]entity.PesertaSusulan, error)
	ListSusulanByMahasiswa(mahasiswaID int) ([]entity.PesertaSusulan, error)
	IsSusulan(aktivasiSesiID, mahasiswaID int) (bool, error)
}

type aktivasiRepository struct{ db *gorm.DB }

func NewAktivasiRepository(db *gorm.DB) AktivasiRepository { return &aktivasiRepository{db: db} }

// ---- aktivasi_sesi ----

func (r *aktivasiRepository) CreateSesi(a *entity.AktivasiSesi) error { return r.db.Create(a).Error }

func (r *aktivasiRepository) UpdateSesi(a *entity.AktivasiSesi) error { return r.db.Save(a).Error }

func (r *aktivasiRepository) FindSesiByID(id int) (*entity.AktivasiSesi, error) {
	var a entity.AktivasiSesi
	err := r.db.Preload("Sesi").Preload("Kelas").
		Preload("AktivasiCourses").Preload("AktivasiCourses.Course").
		First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *aktivasiRepository) FindSesiByComposite(sesiID, kelasID, shift int, gelombang *int) (*entity.AktivasiSesi, error) {
	g := 0 // gelombang NULL diperlakukan sebagai 0 (sesi normal), selaras dgn unique index
	if gelombang != nil {
		g = *gelombang
	}
	var a entity.AktivasiSesi
	err := r.db.Preload("AktivasiCourses").
		Where("sesi_praktikum_id = ? AND kelas_id = ? AND shift = ? AND COALESCE(gelombang, 0) = ?", sesiID, kelasID, shift, g).
		First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *aktivasiRepository) ListSesi() ([]entity.AktivasiSesi, error) {
	var as []entity.AktivasiSesi
	err := r.db.Preload("Sesi").Preload("Kelas").Preload("AktivasiCourses").
		Order("activated_at desc").Find(&as).Error
	return as, err
}

func (r *aktivasiRepository) ListActiveSesi() ([]entity.AktivasiSesi, error) {
	var as []entity.AktivasiSesi
	err := r.db.Preload("Sesi").Preload("Kelas").Preload("AktivasiCourses").Preload("AktivasiCourses.Course").
		Where("is_active = ?", true).Order("activated_at desc").Find(&as).Error
	return as, err
}

// DeleteSesi menghapus aktivasi_sesi BESERTA seluruh data anaknya (cascade manual,
// satu transaksi): jawaban mahasiswa, soal terpilih, pengerjaan (nilai), susulan,
// dan aktivasi_course. Urut anak->induk agar lolos FK RESTRICT (migrasi 021).
func (r *aktivasiRepository) DeleteSesi(id int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// jawaban mahasiswa (anak dari soal_terpilih milik aktivasi ini)
		if err := tx.Where("soal_terpilih_id IN (?)",
			tx.Model(&entity.SoalTerpilih{}).Select("id").Where("aktivasi_sesi_id = ?", id)).
			Delete(&entity.JawabanMahasiswa{}).Error; err != nil {
			return err
		}
		if err := tx.Where("aktivasi_sesi_id = ?", id).Delete(&entity.SoalTerpilih{}).Error; err != nil {
			return err
		}
		if err := tx.Where("aktivasi_sesi_id = ?", id).Delete(&entity.PengerjaanCourse{}).Error; err != nil {
			return err
		}
		if err := tx.Where("aktivasi_sesi_id = ?", id).Delete(&entity.PesertaSusulan{}).Error; err != nil {
			return err
		}
		if err := tx.Where("aktivasi_sesi_id = ?", id).Delete(&entity.AktivasiCourse{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.AktivasiSesi{}, id).Error
	})
}

// ---- aktivasi_course ----

func (r *aktivasiRepository) CreateCourse(ac *entity.AktivasiCourse) error {
	return r.db.Create(ac).Error
}

func (r *aktivasiRepository) UpdateCourse(ac *entity.AktivasiCourse) error {
	return r.db.Save(ac).Error
}

func (r *aktivasiRepository) FindCourseByID(id int) (*entity.AktivasiCourse, error) {
	var ac entity.AktivasiCourse
	if err := r.db.Preload("Course").First(&ac, id).Error; err != nil {
		return nil, err
	}
	return &ac, nil
}

func (r *aktivasiRepository) FindCourse(aktivasiSesiID, courseID int) (*entity.AktivasiCourse, error) {
	var ac entity.AktivasiCourse
	err := r.db.Preload("Course").
		Where("aktivasi_sesi_id = ? AND course_id = ?", aktivasiSesiID, courseID).
		First(&ac).Error
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

func (r *aktivasiRepository) SetCourseStartedAtIfNull(aktivasiCourseID int, t time.Time) error {
	return r.db.Model(&entity.AktivasiCourse{}).
		Where("id = ? AND started_at IS NULL", aktivasiCourseID).
		Update("started_at", t).Error
}

func (r *aktivasiRepository) ListCoursesByAktivasi(aktivasiSesiID int) ([]entity.AktivasiCourse, error) {
	var acs []entity.AktivasiCourse
	err := r.db.Preload("Course").
		Where("aktivasi_sesi_id = ?", aktivasiSesiID).Order("urutan asc").Find(&acs).Error
	return acs, err
}

// ---- peserta_susulan ----

func (r *aktivasiRepository) AddSusulan(p *entity.PesertaSusulan) error {
	return r.db.Create(p).Error
}

func (r *aktivasiRepository) RemoveSusulan(aktivasiSesiID, mahasiswaID int) error {
	return r.db.Where("aktivasi_sesi_id = ? AND mahasiswa_id = ?", aktivasiSesiID, mahasiswaID).
		Delete(&entity.PesertaSusulan{}).Error
}

func (r *aktivasiRepository) ListSusulan(aktivasiSesiID int) ([]entity.PesertaSusulan, error) {
	var ps []entity.PesertaSusulan
	err := r.db.Preload("Mahasiswa").
		Where("aktivasi_sesi_id = ?", aktivasiSesiID).Find(&ps).Error
	return ps, err
}

func (r *aktivasiRepository) ListSusulanByMahasiswa(mahasiswaID int) ([]entity.PesertaSusulan, error) {
	var ps []entity.PesertaSusulan
	err := r.db.Where("mahasiswa_id = ?", mahasiswaID).Find(&ps).Error
	return ps, err
}

func (r *aktivasiRepository) IsSusulan(aktivasiSesiID, mahasiswaID int) (bool, error) {
	var n int64
	err := r.db.Model(&entity.PesertaSusulan{}).
		Where("aktivasi_sesi_id = ? AND mahasiswa_id = ?", aktivasiSesiID, mahasiswaID).
		Count(&n).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return n > 0, nil
}
