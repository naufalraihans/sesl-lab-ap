package usecase

import (
	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
)

// toUserResponse memetakan entity.User ke DTO aman (tanpa password hash).
func toUserResponse(u *entity.User) dto.UserResponse {
	resp := dto.UserResponse{
		ID:         u.ID,
		Role:       string(u.Role),
		NIM:        u.NIM,
		Nama:       u.Nama,
		KelasID:    u.KelasID,
		Shift:      u.Shift,
		FotoURL:    u.FotoURL,
		NomorHP:    u.NomorHP,
		MedsosLink: u.MedsosLink,
	}
	if u.Kelas != nil {
		resp.NamaKelas = u.Kelas.NamaKelas
	}
	return resp
}
