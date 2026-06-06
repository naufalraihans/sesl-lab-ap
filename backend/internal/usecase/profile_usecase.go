package usecase

import (
	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
	"lab-ap/pkg/hash"
)

type ProfileUsecase struct {
	users repository.UserRepository
}

func NewProfileUsecase(u repository.UserRepository) *ProfileUsecase {
	return &ProfileUsecase{users: u}
}

func (uc *ProfileUsecase) Get(userID int) (*entity.User, error) {
	u, err := uc.users.FindByID(userID)
	if err != nil {
		return nil, ErrNotFound
	}
	return u, nil
}

// Update memperbarui profil (nama, kontak, foto, medsos) dan opsional ganti password.
func (uc *ProfileUsecase) Update(userID int, req dto.UpdateProfileRequest) (*entity.User, error) {
	u, err := uc.users.FindByID(userID)
	if err != nil {
		return nil, ErrNotFound
	}
	if req.Nama != nil {
		u.Nama = *req.Nama
	}
	if req.NomorHP != nil {
		u.NomorHP = req.NomorHP
	}
	if req.MedsosLink != nil {
		u.MedsosLink = req.MedsosLink
	}
	if req.FotoURL != nil {
		u.FotoURL = req.FotoURL
	}
	if req.PasswordBaru != nil && *req.PasswordBaru != "" {
		// Verifikasi password lama jika user sudah punya hash.
		if u.PasswordHash != nil {
			if req.PasswordLama == nil || !hash.Verify(*u.PasswordHash, *req.PasswordLama) {
				return nil, ErrUnauthorized
			}
		}
		hashed, err := hash.Password(*req.PasswordBaru)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = &hashed
	}
	if err := uc.users.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}
