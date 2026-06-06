package usecase

import (
	"errors"
	"time"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
	"lab-ap/pkg/hash"
	"lab-ap/pkg/jwt"
	"lab-ap/pkg/online"

	"gorm.io/gorm"
)

type AuthUsecase struct {
	users  repository.UserRepository
	kelas  repository.KelasRepository
	jwt    *jwt.Manager
	online *online.Registry
}

func NewAuthUsecase(u repository.UserRepository, k repository.KelasRepository, j *jwt.Manager, o *online.Registry) *AuthUsecase {
	return &AuthUsecase{users: u, kelas: k, jwt: j, online: o}
}

// CekNIM menentukan alur first-time login (login / register / ditolak).
func (uc *AuthUsecase) CekNIM(nim string) (*dto.CekNIMResponse, error) {
	u, err := uc.users.FindByNIM(nim)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.CekNIMResponse{NIM: nim, Ditemukan: false, Pesan: "NIM tidak terdaftar. Hubungi admin."}, nil
		}
		return nil, err
	}

	resp := &dto.CekNIMResponse{
		NIM:          u.NIM,
		Ditemukan:    true,
		IsRegistered: u.IsRegistered,
		Nama:         u.Nama,
	}

	if u.IsRegistered {
		resp.Pesan = "Silakan masukkan password Anda."
		return resp, nil
	}

	// Belum register: cek apakah akses register kelasnya dibuka.
	if u.KelasID != nil {
		if k, err := uc.kelas.FindByID(*u.KelasID); err == nil {
			resp.IsRegisterOpen = k.IsRegisterOpen
		}
	}
	if resp.IsRegisterOpen {
		resp.Pesan = "Akun belum terdaftar. Silakan buat password."
	} else {
		resp.Pesan = "Akses register belum dibuka oleh admin."
	}
	return resp, nil
}

// Login memvalidasi password user yang sudah terdaftar.
func (uc *AuthUsecase) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	u, err := uc.users.FindByNIM(req.NIM)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	// Admin boleh login walau is_registered default; mahasiswa wajib sudah register.
	if u.Role == entity.RoleUser && !u.IsRegistered {
		return nil, ErrRegisterClosed
	}
	if u.PasswordHash == nil || !hash.Verify(*u.PasswordHash, req.Password) {
		return nil, ErrUnauthorized
	}
	return uc.issue(u)
}

// Register membuat password pertama kali (jika akses dibuka).
func (uc *AuthUsecase) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	u, err := uc.users.FindByNIM(req.NIM)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if u.IsRegistered {
		return nil, ErrConflict // sudah register, harus login
	}
	// Cek akses register kelas.
	if u.KelasID == nil {
		return nil, ErrRegisterClosed
	}
	k, err := uc.kelas.FindByID(*u.KelasID)
	if err != nil || !k.IsRegisterOpen {
		return nil, ErrRegisterClosed
	}

	hashed, err := hash.Password(req.Password)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = &hashed
	u.IsRegistered = true
	if err := uc.users.Update(u); err != nil {
		return nil, err
	}
	return uc.issue(u)
}

// Logout menghapus entri online.
func (uc *AuthUsecase) Logout(userID int) {
	uc.online.Remove(userID)
}

// issue membuat token + mencatat login & online.
func (uc *AuthUsecase) issue(u *entity.User) (*dto.AuthResponse, error) {
	token, err := uc.jwt.Generate(u.ID, u.NIM, string(u.Role))
	if err != nil {
		return nil, err
	}
	now := time.Now()
	u.LastLoginAt = &now
	_ = uc.users.Update(u)

	// Muat relasi kelas untuk response.
	if u.KelasID != nil && u.Kelas == nil {
		if full, err := uc.users.FindByID(u.ID); err == nil {
			u = full
		}
	}
	uc.online.Touch(u.ID, string(u.Role))
	return &dto.AuthResponse{Token: token, User: toUserResponse(u)}, nil
}
