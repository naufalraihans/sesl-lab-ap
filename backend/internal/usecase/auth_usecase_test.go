package usecase_test

import (
	"testing"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository/mocks"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/hash"
	"lab-ap/pkg/jwt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAuthUsecase(t *testing.T) (*usecase.AuthUsecase, *mocks.UserRepository, *mocks.KelasRepository) {
	mockUserRepo := mocks.NewUserRepository(t)
	mockKelasRepo := mocks.NewKelasRepository(t)
	jwtManager := jwt.NewManager("secret", 24)

	uc := usecase.NewAuthUsecase(mockUserRepo, mockKelasRepo, jwtManager)
	return uc, mockUserRepo, mockKelasRepo
}

func TestAuthUsecase_Login_Success(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)

	hashedPassword, _ := hash.Password("password123")
	mockUser := &entity.User{
		ID:           1,
		NIM:          "123456",
		Role:         entity.RoleUser,
		IsRegistered: true,
		PasswordHash: &hashedPassword,
	}

	// Ekspektasi: ketika FindByNIM dipanggil dengan "123456", kembalikan mockUser
	mockUserRepo.On("FindByNIM", "123456").Return(mockUser, nil)
	// Ekspektasi: Update last login
	mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)

	req := dto.LoginRequest{
		NIM:      "123456",
		Password: "password123",
	}

	resp, err := uc.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "123456", resp.User.NIM)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_WrongPassword(t *testing.T) {
	uc, mockUserRepo, _ := setupAuthUsecase(t)

	hashedPassword, _ := hash.Password("password123")
	mockUser := &entity.User{
		ID:           1,
		NIM:          "123456",
		Role:         entity.RoleUser,
		IsRegistered: true,
		PasswordHash: &hashedPassword,
	}

	mockUserRepo.On("FindByNIM", "123456").Return(mockUser, nil)

	req := dto.LoginRequest{
		NIM:      "123456",
		Password: "wrongpassword",
	}

	resp, err := uc.Login(req)

	assert.Error(t, err)
	assert.Equal(t, usecase.ErrUnauthorized, err)
	assert.Nil(t, resp)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthUsecase_CekNIM_NotRegistered(t *testing.T) {
	uc, mockUserRepo, mockKelasRepo := setupAuthUsecase(t)

	kelasID := 1
	mockUser := &entity.User{
		ID:           1,
		NIM:          "123456",
		IsRegistered: false,
		KelasID:      &kelasID,
	}
	mockKelas := &entity.Kelas{
		ID:             1,
		IsRegisterOpen: true,
	}

	mockUserRepo.On("FindByNIM", "123456").Return(mockUser, nil)
	mockKelasRepo.On("FindByID", 1).Return(mockKelas, nil)

	resp, err := uc.CekNIM("123456")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsRegistered)
	assert.True(t, resp.IsRegisterOpen)
	assert.Equal(t, "Akun belum terdaftar. Silakan buat password.", resp.Pesan)

	mockUserRepo.AssertExpectations(t)
	mockKelasRepo.AssertExpectations(t)
}
