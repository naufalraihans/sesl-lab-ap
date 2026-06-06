package usecase

import (
	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
)

type JadwalUsecase struct {
	repo repository.JadwalRepository
}

func NewJadwalUsecase(r repository.JadwalRepository) *JadwalUsecase { return &JadwalUsecase{repo: r} }

func (uc *JadwalUsecase) List() ([]entity.Jadwal, error) { return uc.repo.List() }

func (uc *JadwalUsecase) ForKelasShift(kelasID, shift int) (*entity.Jadwal, error) {
	j, err := uc.repo.FindByKelasShift(kelasID, shift)
	if err != nil {
		return nil, ErrNotFound
	}
	return j, nil
}

func (uc *JadwalUsecase) Create(req dto.JadwalRequest) (*entity.Jadwal, error) {
	j := &entity.Jadwal{
		KelasID:    req.KelasID,
		Shift:      req.Shift,
		Hari:       req.Hari,
		JamMulai:   req.JamMulai,
		JamSelesai: req.JamSelesai,
		Keterangan: req.Keterangan,
	}
	if err := uc.repo.Create(j); err != nil {
		return nil, ErrConflict
	}
	return j, nil
}

func (uc *JadwalUsecase) Update(id int, req dto.JadwalRequest) (*entity.Jadwal, error) {
	j, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	j.KelasID = req.KelasID
	j.Shift = req.Shift
	j.Hari = req.Hari
	j.JamMulai = req.JamMulai
	j.JamSelesai = req.JamSelesai
	j.Keterangan = req.Keterangan
	if err := uc.repo.Update(j); err != nil {
		return nil, err
	}
	return j, nil
}

func (uc *JadwalUsecase) Delete(id int) error { return uc.repo.Delete(id) }
