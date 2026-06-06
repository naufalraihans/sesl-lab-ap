package usecase

import (
	"time"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
)

type PedomanUsecase struct {
	repo repository.PedomanRepository
}

func NewPedomanUsecase(r repository.PedomanRepository) *PedomanUsecase {
	return &PedomanUsecase{repo: r}
}

func (uc *PedomanUsecase) List() ([]entity.PedomanLaporan, error) { return uc.repo.List() }

func (uc *PedomanUsecase) Create(req dto.PedomanRequest) (*entity.PedomanLaporan, error) {
	p := &entity.PedomanLaporan{
		NamaDokumen:  req.NamaDokumen,
		FileURL:      req.FileURL,
		DiunggahPada: time.Now(),
	}
	if err := uc.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *PedomanUsecase) Update(id int, req dto.PedomanRequest) (*entity.PedomanLaporan, error) {
	p, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	p.NamaDokumen = req.NamaDokumen
	p.FileURL = req.FileURL
	if err := uc.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *PedomanUsecase) Delete(id int) error { return uc.repo.Delete(id) }
