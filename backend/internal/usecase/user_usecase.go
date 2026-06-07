package usecase

import (
	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
	"lab-ap/pkg/hash"
)

type UserUsecase struct {
	users repository.UserRepository
	kelas repository.KelasRepository
}

func NewUserUsecase(u repository.UserRepository, k repository.KelasRepository) *UserUsecase {
	return &UserUsecase{users: u, kelas: k}
}

func (uc *UserUsecase) ListMahasiswa(kelasID, shift *int) ([]entity.User, error) {
	return uc.users.List(string(entity.RoleUser), kelasID, shift)
}

func (uc *UserUsecase) CreateMahasiswa(req dto.UserRequest) (*entity.User, error) {
	u := &entity.User{
		Role:         entity.RoleUser,
		NIM:          req.NIM,
		Nama:         req.Nama,
		KelasID:      req.KelasID,
		Shift:        req.Shift,
		Kelompok:     req.Kelompok,
		IsRegistered: false,
	}
	if err := uc.users.Create(u); err != nil {
		return nil, ErrConflict
	}
	return u, nil
}

// BulkUpsertMahasiswa memproses array UserRequest dan melakukan upsert.
func (uc *UserUsecase) BulkUpsertMahasiswa(req dto.UserBulkRequest) (*dto.BulkResponse, error) {
	if len(req.Users) == 0 {
		return &dto.BulkResponse{}, nil
	}

	var entities []entity.User
	for _, r := range req.Users {
		entities = append(entities, entity.User{
			Role:         entity.RoleUser,
			NIM:          r.NIM,
			Nama:         r.Nama,
			KelasID:      r.KelasID,
			Shift:        r.Shift,
			Kelompok:     r.Kelompok,
			IsRegistered: false, // Default false, tapi OnConflict TIDAK akan menimpa is_registered
		})
	}

	err := uc.users.BulkUpsert(entities)
	if err != nil {
		return nil, err
	}

	// GORM's CreateInBatches with OnConflict doesn't easily return exact insert/update counts in MySQL/Postgres
	// without complex RETURNING clauses. For now, we return TotalProcessed.
	return &dto.BulkResponse{
		TotalProcessed: len(entities),
		TotalInserted:  0, // Can be refined if needed
		TotalUpdated:   0,
	}, nil
}

func (uc *UserUsecase) UpdateMahasiswa(id int, req dto.UserRequest) (*entity.User, error) {
	u, err := uc.users.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	u.NIM = req.NIM
	u.Nama = req.Nama
	u.KelasID = req.KelasID
	u.Shift = req.Shift
	u.Kelompok = req.Kelompok
	if err := uc.users.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UserUsecase) Delete(id int) error { return uc.users.Delete(id) }

// ResetPassword mengosongkan password & menandai belum register (mahasiswa register ulang).
func (uc *UserUsecase) ResetPassword(id int) error {
	u, err := uc.users.FindByID(id)
	if err != nil {
		return ErrNotFound
	}
	u.PasswordHash = nil
	u.IsRegistered = false
	return uc.users.Update(u)
}

// ---- Asisten (role admin) ----

func (uc *UserUsecase) ListAsisten() ([]entity.User, error) { return uc.users.ListAsisten() }

func (uc *UserUsecase) CreateAsisten(req dto.AsistenRequest) (*entity.User, error) {
	u := &entity.User{
		Role:       entity.RoleAdmin,
		NIM:        req.NIM,
		Nama:       req.Nama,
		NomorHP:    req.NomorHP,
		MedsosLink: req.MedsosLink,
		FotoURL:    req.FotoURL,
	}
	if req.Password != nil && *req.Password != "" {
		h, err := hash.Password(*req.Password)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = &h
		u.IsRegistered = true
	}
	if err := uc.users.Create(u); err != nil {
		return nil, ErrConflict
	}
	return u, nil
}

func (uc *UserUsecase) UpdateAsisten(id int, req dto.AsistenRequest) (*entity.User, error) {
	u, err := uc.users.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	u.NIM = req.NIM
	u.Nama = req.Nama
	u.NomorHP = req.NomorHP
	u.MedsosLink = req.MedsosLink
	u.FotoURL = req.FotoURL
	if req.Password != nil && *req.Password != "" {
		h, err := hash.Password(*req.Password)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = &h
		u.IsRegistered = true
	}
	if err := uc.users.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

// ---- Kelas register open ----

func (uc *UserUsecase) SetRegisterOpen(kelasID int, open bool) error {
	return uc.kelas.SetRegisterOpen(kelasID, open)
}
