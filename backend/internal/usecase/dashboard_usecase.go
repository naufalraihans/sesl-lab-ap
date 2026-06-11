package usecase

import (
	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
)

type DashboardUsecase struct {
	users      repository.UserRepository
	aktivasi   repository.AktivasiRepository
	pengerjaan repository.PengerjaanRepository
}

func NewDashboardUsecase(u repository.UserRepository, a repository.AktivasiRepository, p repository.PengerjaanRepository) *DashboardUsecase {
	return &DashboardUsecase{users: u, aktivasi: a, pengerjaan: p}
}

// Statistik menyusun ringkasan dashboard admin.
func (uc *DashboardUsecase) Statistik() (*dto.StatistikResponse, error) {
	totalMhs, err := uc.users.CountByRole(entity.RoleUser)
	if err != nil {
		return nil, err
	}
	totalAsisten, err := uc.users.CountByRole(entity.RoleAdmin)
	if err != nil {
		return nil, err
	}
	sudah, err := uc.users.CountRegistered(true)
	if err != nil {
		return nil, err
	}
	belum, err := uc.users.CountRegistered(false)
	if err != nil {
		return nil, err
	}
	perKelas, err := uc.users.CountPerKelasShift()
	if err != nil {
		return nil, err
	}

	resp := &dto.StatistikResponse{
		TotalMahasiswa: totalMhs,
		TotalAsisten:   totalAsisten,
		SudahRegister:  sudah,
		BelumRegister:  belum,
		PerKelasShift:  perKelas,
	}

	aktifList, err := uc.aktivasi.ListActiveSesi()
	if err != nil {
		return nil, err
	}
	for _, a := range aktifList {
		info := dto.SesiAktifInfo{AktivasiSesiID: a.ID, Shift: a.Shift}
		if a.Sesi != nil {
			info.JudulSesi = a.Sesi.JudulSesi
		}
		if a.Kelas != nil {
			info.NamaKelas = a.Kelas.NamaKelas
		}
		for _, ac := range a.AktivasiCourses {
			prog, _ := uc.pengerjaan.ProgressSummary(a.ID, ac.CourseID)
			cp := dto.CourseProgressInfo{
				CourseID: ac.CourseID,
				IsOpen:   ac.IsOpen,
				Selesai:  prog.Selesai,
				Sedang:   prog.Sedang,
				Belum:    prog.Belum,
			}
			if ac.Course != nil {
				cp.Jenis = string(ac.Course.Jenis)
			}
			info.Courses = append(info.Courses, cp)
		}
		resp.SesiAktif = append(resp.SesiAktif, info)
	}
	return resp, nil
}
