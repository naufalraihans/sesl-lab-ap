package usecase

import (
	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
)

type SesiUsecase struct {
	sesi   repository.SesiRepository
	course repository.CourseRepository
}

func NewSesiUsecase(s repository.SesiRepository, c repository.CourseRepository) *SesiUsecase {
	return &SesiUsecase{sesi: s, course: c}
}

// ---- Sesi ----

func (uc *SesiUsecase) List() ([]entity.SesiPraktikum, error) { return uc.sesi.List() }

func (uc *SesiUsecase) Get(id int) (*entity.SesiPraktikum, error) {
	s, err := uc.sesi.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return s, nil
}

func (uc *SesiUsecase) Create(req dto.SesiRequest) (*entity.SesiPraktikum, error) {
	s := &entity.SesiPraktikum{
		JudulSesi:      req.JudulSesi,
		Deskripsi:      req.Deskripsi,
		Urutan:         req.Urutan,
		IsUjianPraktik: req.IsUjianPraktik,
	}
	if err := uc.sesi.Create(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *SesiUsecase) Update(id int, req dto.SesiRequest) (*entity.SesiPraktikum, error) {
	s, err := uc.sesi.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	s.JudulSesi = req.JudulSesi
	s.Deskripsi = req.Deskripsi
	s.Urutan = req.Urutan
	s.IsUjianPraktik = req.IsUjianPraktik
	if err := uc.sesi.Update(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *SesiUsecase) Delete(id int) error { return uc.sesi.Delete(id) }

// ---- Course ----

func (uc *SesiUsecase) ListCourse(sesiID int) ([]entity.Course, error) {
	return uc.course.ListBySesi(sesiID)
}

func (uc *SesiUsecase) CreateCourse(sesiID int, req dto.CourseRequest) (*entity.Course, error) {
	if _, err := uc.sesi.FindByID(sesiID); err != nil {
		return nil, ErrNotFound
	}
	c := &entity.Course{
		SesiPraktikumID: sesiID,
		Jenis:           entity.JenisCourse(req.Jenis),
		Judul:           req.Judul,
		Deskripsi:       req.Deskripsi,
		DurasiMenit:     req.DurasiMenit,
	}
	if err := uc.course.Create(c); err != nil {
		return nil, ErrConflict // kemungkinan duplikat jenis per sesi
	}
	return c, nil
}

func (uc *SesiUsecase) UpdateCourse(id int, req dto.CourseRequest) (*entity.Course, error) {
	c, err := uc.course.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	c.Jenis = entity.JenisCourse(req.Jenis)
	c.Judul = req.Judul
	c.Deskripsi = req.Deskripsi
	c.DurasiMenit = req.DurasiMenit
	if err := uc.course.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *SesiUsecase) DeleteCourse(id int) error { return uc.course.Delete(id) }
