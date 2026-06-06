package usecase

import (
	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
)

type AmpuanUsecase struct {
	repo repository.AmpuanRepository
}

func NewAmpuanUsecase(r repository.AmpuanRepository) *AmpuanUsecase {
	return &AmpuanUsecase{repo: r}
}

func (uc *AmpuanUsecase) List() ([]entity.AmpuanKelompok, error) { return uc.repo.List() }

func (uc *AmpuanUsecase) ListByKelas(kelasID int) ([]entity.AmpuanKelompok, error) {
	return uc.repo.ListByKelas(kelasID)
}

func (uc *AmpuanUsecase) Create(req dto.AmpuanRequest) (*entity.AmpuanKelompok, error) {
	a := &entity.AmpuanKelompok{
		AsistenID: req.AsistenID,
		KelasID:   req.KelasID,
		Kelompok:  req.Kelompok,
	}
	if err := uc.repo.Create(a); err != nil {
		return nil, ErrConflict
	}
	return a, nil
}

func (uc *AmpuanUsecase) Delete(id int) error { return uc.repo.Delete(id) }
