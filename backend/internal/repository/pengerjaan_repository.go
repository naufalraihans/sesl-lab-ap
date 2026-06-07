package repository

import (
	"errors"

	"lab-ap/internal/entity"

	"gorm.io/gorm"
)

type PengerjaanRepository interface {
	Create(p *entity.PengerjaanCourse) error
	Update(p *entity.PengerjaanCourse) error
	Find(mahasiswaID, aktivasiSesiID, courseID int) (*entity.PengerjaanCourse, error)
	// FindOrCreate mengembalikan record yang ada atau membuat baru (status belum_dikerjakan).
	FindOrCreate(mahasiswaID, aktivasiSesiID, courseID int) (*entity.PengerjaanCourse, error)
	ListByMahasiswa(mahasiswaID int) ([]entity.PengerjaanCourse, error)
	ListByAktivasiCourse(aktivasiSesiID, courseID int) ([]entity.PengerjaanCourse, error)
	// MarkSelesaiForCourse: set status selesai utk semua yg belum selesai (mass auto-submit).
	MarkSelesaiForCourse(aktivasiSesiID, courseID int) error
	ProgressSummary(aktivasiSesiID, courseID int) (ProgressSummary, error)
	// FindExpired mengambil pengerjaan yang sedang berjalan & sudah lewat deadline (untuk sweeper).
	FindExpired() ([]ExpiredPengerjaan, error)
}

// ProgressSummary: ringkasan progress per course aktivasi.
type ProgressSummary struct {
	Selesai int64 `json:"selesai"`
	Sedang  int64 `json:"sedang"`
	Belum   int64 `json:"belum"`
}

// ExpiredPengerjaan: pengerjaan yang sudah lewat deadline (waktu_mulai + durasi_menit).
type ExpiredPengerjaan struct {
	MahasiswaID    int
	AktivasiSesiID int
	CourseID       int
}

type pengerjaanRepository struct{ db *gorm.DB }

func NewPengerjaanRepository(db *gorm.DB) PengerjaanRepository {
	return &pengerjaanRepository{db: db}
}

func (r *pengerjaanRepository) Create(p *entity.PengerjaanCourse) error { return r.db.Create(p).Error }
func (r *pengerjaanRepository) Update(p *entity.PengerjaanCourse) error { return r.db.Save(p).Error }

func (r *pengerjaanRepository) Find(mahasiswaID, aktivasiSesiID, courseID int) (*entity.PengerjaanCourse, error) {
	var p entity.PengerjaanCourse
	err := r.db.Where("mahasiswa_id = ? AND aktivasi_sesi_id = ? AND course_id = ?",
		mahasiswaID, aktivasiSesiID, courseID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pengerjaanRepository) FindOrCreate(mahasiswaID, aktivasiSesiID, courseID int) (*entity.PengerjaanCourse, error) {
	p, err := r.Find(mahasiswaID, aktivasiSesiID, courseID)
	if err == nil {
		return p, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	np := &entity.PengerjaanCourse{
		MahasiswaID:    mahasiswaID,
		AktivasiSesiID: aktivasiSesiID,
		CourseID:       courseID,
		Status:         entity.StatusBelum,
	}
	if err := r.db.Create(np).Error; err != nil {
		return nil, err
	}
	return np, nil
}

func (r *pengerjaanRepository) ListByMahasiswa(mahasiswaID int) ([]entity.PengerjaanCourse, error) {
	var ps []entity.PengerjaanCourse
	err := r.db.Where("mahasiswa_id = ?", mahasiswaID).Find(&ps).Error
	return ps, err
}

func (r *pengerjaanRepository) ListByAktivasiCourse(aktivasiSesiID, courseID int) ([]entity.PengerjaanCourse, error) {
	var ps []entity.PengerjaanCourse
	err := r.db.Where("aktivasi_sesi_id = ? AND course_id = ?", aktivasiSesiID, courseID).Find(&ps).Error
	return ps, err
}

func (r *pengerjaanRepository) MarkSelesaiForCourse(aktivasiSesiID, courseID int) error {
	return r.db.Model(&entity.PengerjaanCourse{}).
		Where("aktivasi_sesi_id = ? AND course_id = ? AND status <> ?", aktivasiSesiID, courseID, entity.StatusSelesai).
		Updates(map[string]interface{}{"status": entity.StatusSelesai, "waktu_selesai": gorm.Expr("NOW()")}).Error
}

func (r *pengerjaanRepository) ProgressSummary(aktivasiSesiID, courseID int) (ProgressSummary, error) {
	var s ProgressSummary
	base := r.db.Model(&entity.PengerjaanCourse{}).
		Where("aktivasi_sesi_id = ? AND course_id = ?", aktivasiSesiID, courseID)
	if err := base.Session(&gorm.Session{}).Where("status = ?", entity.StatusSelesai).Count(&s.Selesai).Error; err != nil {
		return s, err
	}
	if err := base.Session(&gorm.Session{}).Where("status = ?", entity.StatusSedang).Count(&s.Sedang).Error; err != nil {
		return s, err
	}
	if err := base.Session(&gorm.Session{}).Where("status = ?", entity.StatusBelum).Count(&s.Belum).Error; err != nil {
		return s, err
	}
	return s, nil
}

func (r *pengerjaanRepository) FindExpired() ([]ExpiredPengerjaan, error) {
	var out []ExpiredPengerjaan
	// Deadline = waktu_mulai + course.durasi_menit menit. Hanya yg status sedang_dikerjakan.
	err := r.db.Table("pengerjaan_course AS pc").
		Select("pc.mahasiswa_id, pc.aktivasi_sesi_id, pc.course_id").
		Joins("JOIN course c ON c.id = pc.course_id").
		Where("pc.status = ?", entity.StatusSedang).
		Where("pc.waktu_mulai IS NOT NULL").
		Where("NOW() > pc.waktu_mulai + (c.durasi_menit * interval '1 minute')").
		Scan(&out).Error
	return out, err
}
