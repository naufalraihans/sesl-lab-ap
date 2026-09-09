package usecase_test

import (
	"errors"
	"testing"

	"lab-ap/config"
	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository/mocks"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/hash"
	"lab-ap/pkg/jwt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupAuthUsecase(t *testing.T) (*usecase.AuthUsecase, *mocks.UserRepository, *mocks.KelasRepository) {
	mockUserRepo := mocks.NewUserRepository(t)
	mockKelasRepo := mocks.NewKelasRepository(t)
	jwtManager := jwt.NewManager("secret", 24)

	uc := usecase.NewAuthUsecase(mockUserRepo, mockKelasRepo, jwtManager, &config.Config{}, hash.FbScryptConfig{})
	return uc, mockUserRepo, mockKelasRepo
}

// rosterUser: baris hasil import asisten — belum punya password, belum diklaim.
func rosterUser() *entity.User {
	kelasID := 1
	return &entity.User{
		ID:       1,
		NIM:      "123456",
		Nama:     "Budi",
		Role:     entity.RoleUser,
		KelasID:  &kelasID,
		Email:    nil,
	}
}

// registeredUser: sudah klaim roster — punya password.
func registeredUser(role entity.RoleType) *entity.User {
	hashed, _ := hash.Password("password123")
	supa := "supa-uid-1"
	return &entity.User{
		ID:              1,
		NIM:             "123456",
		Nama:            "Budi",
		Role:            role,
		KelasID:         nil,
		IsRegistered:    true,
		PasswordHash:    &hashed,
		Email:           &[]string{"budi@gmail.com"}[0],
		SupabaseUserID:  &supa,
	}
}

// ============================= CekNIM =============================

func TestCekNIM_TidakDiRoster(t *testing.T) {
	uc, mockUserRepo, mockKelasRepo := setupAuthUsecase(t)
	mockUserRepo.On("FindByNIM", "999999").Return(nil, gorm.ErrRecordNotFound)

	resp, err := uc.CekNIM("999999")

	assert.NoError(t, err)
	assert.False(t, resp.Ditemukan)
	assert.Equal(t, "Akun tidak ada.", resp.Pesan)
	mockUserRepo.AssertExpectations(t)
	mockKelasRepo.AssertExpectations(t)
}

func TestCekNIM_SudahRegister(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	u := registeredUser(entity.RoleUser)
	mockUserRepo.On("FindByNIM", "123456").Return(u, nil)

	resp, err := uc.CekNIM("123456")

	assert.NoError(t, err)
	assert.True(t, resp.Ditemukan)
	assert.True(t, resp.IsRegistered)
	assert.Equal(t, "Silakan masukkan password Anda.", resp.Pesan)
	mockUserRepo.AssertExpectations(t)
}

func TestCekNIM_RosterBelumKlaim(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	mockUserRepo.On("FindByNIM", "123456").Return(rosterUser(), nil)

	resp, err := uc.CekNIM("123456")

	assert.NoError(t, err)
	assert.True(t, resp.Ditemukan)
	assert.False(t, resp.IsRegistered)
	assert.Equal(t, "Akun belum terdaftar. Silakan buat password.", resp.Pesan)
	mockUserRepo.AssertExpectations(t)
}

// ============================= Login =============================

func TestLogin_Sukses_NIM(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	u := registeredUser(entity.RoleUser)
	mockUserRepo.On("FindByNIM", "123456").Return(u, nil)
	mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)

	resp, err := uc.Login(dto.LoginRequest{Identifier: "123456", Password: "password123"})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "123456", resp.User.NIM)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_Sukses_Email(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	u := registeredUser(entity.RoleUser)
	email := "budi@gmail.com"
	u.Email = &email
	mockUserRepo.On("FindByEmail", "budi@gmail.com").Return(u, nil)
	mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)

	resp, err := uc.Login(dto.LoginRequest{Identifier: "Budi@Gmail.com", Password: "password123"})

	assert.NoError(t, err) // email dinormalisasi lowercase sebelum lookup
	assert.NotEmpty(t, resp.Token)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_PasswordSalah(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	u := registeredUser(entity.RoleUser)
	mockUserRepo.On("FindByNIM", "123456").Return(u, nil)

	resp, err := uc.Login(dto.LoginRequest{Identifier: "123456", Password: "salah"})

	assert.ErrorIs(t, err, usecase.ErrUnauthorized)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_NIMTidakAda(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	mockUserRepo.On("FindByNIM", "000000").Return(nil, gorm.ErrRecordNotFound)

	resp, err := uc.Login(dto.LoginRequest{Identifier: "000000", Password: "x"})

	assert.ErrorIs(t, err, usecase.ErrNotFound)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

// Anti-enumeration: user roster yang belum register = "akun tidak ada", bukan "belum aktivasi".
func TestLogin_RosterBelumKlaim_Ditolak(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	mockUserRepo.On("FindByNIM", "123456").Return(rosterUser(), nil)

	resp, err := uc.Login(dto.LoginRequest{Identifier: "123456", Password: "apapun"})

	assert.ErrorIs(t, err, usecase.ErrNotFound)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

// Admin tanpa IsRegistered tetap bisa login (guard cuma untuk role user).
func TestLogin_AdminBelumRegisterFlag_TetapMasuk(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	u := registeredUser(entity.RoleAdmin)
	u.IsRegistered = false
	mockUserRepo.On("FindByNIM", "123456").Return(u, nil)
	mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)

	resp, err := uc.Login(dto.LoginRequest{Identifier: "123456", Password: "password123"})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	mockUserRepo.AssertExpectations(t)
}

// ============================= Register =============================

func TestRegister_Sukses_LocalFallback(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	// cfg tanpa Supabase -> jalur fallback lokal (tanpa network).
	mockUserRepo.On("FindByNIM", "123456").Return(rosterUser(), nil)
	mockUserRepo.On("FindByEmail", "budi@gmail.com").Return(nil, gorm.ErrRecordNotFound)
	mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)
	mockUserRepo.On("FindByID", 1).Return(rosterUser(), nil)

	resp, err := uc.Register(dto.RegisterRequest{NIM: "123456", Email: "budi@gmail.com", Password: "password123"})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "123456", resp.User.NIM)
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_NIMTidakDiRoster(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	mockUserRepo.On("FindByNIM", "000000").Return(nil, gorm.ErrRecordNotFound)

	resp, err := uc.Register(dto.RegisterRequest{NIM: "000000", Email: "budi@gmail.com", Password: "password123"})

	assert.ErrorIs(t, err, usecase.ErrNotFound)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_SudahDiklaim(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	mockUserRepo.On("FindByNIM", "123456").Return(registeredUser(entity.RoleUser), nil)

	resp, err := uc.Register(dto.RegisterRequest{NIM: "123456", Email: "budi@gmail.com", Password: "password123"})

	assert.ErrorIs(t, err, usecase.ErrConflict)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_KelasBelumLengkap(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	u := rosterUser()
	u.KelasID = nil
	mockUserRepo.On("FindByNIM", "123456").Return(u, nil)

	resp, err := uc.Register(dto.RegisterRequest{NIM: "123456", Email: "budi@gmail.com", Password: "password123"})

	assert.ErrorIs(t, err, usecase.ErrBadRequest)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "kelas belum lengkap")
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_EmailSudahDipakai(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	mockUserRepo.On("FindByNIM", "123456").Return(rosterUser(), nil)
	mockUserRepo.On("FindByEmail", "budi@gmail.com").Return(registeredUser(entity.RoleUser), nil)

	resp, err := uc.Register(dto.RegisterRequest{NIM: "123456", Email: "budi@gmail.com", Password: "password123"})

	assert.ErrorIs(t, err, usecase.ErrConflict)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

// ---- Validasi email (indirect lewat Register; allowlist + anti-alias) ----

func TestRegister_EmailDomainTidakDiizinkan(t *testing.T) {
	uc, _, _ := setupAuthUsecase(t)

	resp, err := uc.Register(dto.RegisterRequest{NIM: "123456", Email: "budi@bukanmail.xyz", Password: "password123"})

	assert.ErrorIs(t, err, usecase.ErrBadRequest)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "tidak diizinkan")
}

func TestRegister_EmailFormatBuruk(t *testing.T) {
	uc, _, _ := setupAuthUsecase(t)

	_, err := uc.Register(dto.RegisterRequest{NIM: "123456", Email: "tanpa-at", Password: "password123"})

	assert.ErrorIs(t, err, usecase.ErrBadRequest)
	assert.Contains(t, err.Error(), "format email tidak valid")
}

func TestRegister_GmailDotTrick(t *testing.T) {
	uc, _, _ := setupAuthUsecase(t)

	_, err := uc.Register(dto.RegisterRequest{NIM: "123456", Email: "b.u.d.i@gmail.com", Password: "password123"})

	assert.ErrorIs(t, err, usecase.ErrBadRequest)
	assert.Contains(t, err.Error(), "'.'")
}

func TestRegister_PlusAlias(t *testing.T) {
	uc, _, _ := setupAuthUsecase(t)

	_, err := uc.Register(dto.RegisterRequest{NIM: "123456", Email: "budi+spam@gmail.com", Password: "password123"})

	assert.ErrorIs(t, err, usecase.ErrBadRequest)
	assert.Contains(t, err.Error(), "'+'")
}

func TestRegister_ItplnTanpaFlag_Ditolak(t *testing.T) {
	uc, _, _ := setupAuthUsecase(t)
	// cfg default: EmailAllowItpln = false -> ditolak validasi email, repo tidak disentuh.

	_, err := uc.Register(dto.RegisterRequest{NIM: "123456", Email: "budi@itpln.ac.id", Password: "password123"})

	assert.ErrorIs(t, err, usecase.ErrBadRequest)
	assert.Contains(t, err.Error(), "tidak diizinkan")
}

// ============================= Forgot/Reset =============================

// Anti-enumeration: email tak terdaftar tetap 200 (err nil), tidak panic tanpa Supabase.
func TestForgotPassword_EmailTidakTerdaftar_TetapSukses(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	mockUserRepo.On("FindByEmail", "hantu@gmail.com").Return(nil, gorm.ErrRecordNotFound)

	err := uc.ForgotPassword("hantu@gmail.com")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

// Terdaftar tapi tanpa konfigurasi Supabase -> diam-diam no-op (tidak panic).
func TestForgotPassword_TanpaSupabase_NoOp(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)
	mockUserRepo.On("FindByEmail", "budi@gmail.com").Return(registeredUser(entity.RoleUser), nil)

	err := uc.ForgotPassword("budi@gmail.com")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

var _ = errors.New // jaga import errors jika suite dikembangkan
