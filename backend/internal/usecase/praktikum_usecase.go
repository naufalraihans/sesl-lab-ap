package usecase

import (
	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
)

// PraktikumUsecase melayani sisi mahasiswa: dashboard + daftar sesi.
type PraktikumUsecase struct {
	sesi       repository.SesiRepository
	course     repository.CourseRepository
	aktivasi   repository.AktivasiRepository
	pengerjaan repository.PengerjaanRepository
	users      repository.UserRepository
	jadwal     repository.JadwalRepository
}

func NewPraktikumUsecase(
	s repository.SesiRepository,
	c repository.CourseRepository,
	a repository.AktivasiRepository,
	p repository.PengerjaanRepository,
	u repository.UserRepository,
	j repository.JadwalRepository,
) *PraktikumUsecase {
	return &PraktikumUsecase{sesi: s, course: c, aktivasi: a, pengerjaan: p, users: u, jadwal: j}
}

// ListSesi mengembalikan seluruh sesi dari sudut pandang mahasiswa (aktif/terkunci).
func (uc *PraktikumUsecase) ListSesi(userID int) ([]dto.SesiUserItem, error) {
	user, err := uc.users.FindByID(userID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	allSesi, err := uc.sesi.List()
	if err != nil {
		return nil, err
	}

	// Map sesiID -> aktivasi milik kelas+shift user.
	aktivasiBySesi := map[int]*entity.AktivasiSesi{}
	if user.KelasID != nil && user.Shift != nil {
		for _, s := range allSesi {
			if a, err := uc.aktivasi.FindSesiByComposite(s.ID, *user.KelasID, *user.Shift); err == nil {
				full, _ := uc.aktivasi.FindSesiByID(a.ID)
				aktivasiBySesi[s.ID] = full
			}
		}
	}

	// Susulan: aktivasi di kelas/shift lain.
	susulanAktivasiBySesi := map[int]*entity.AktivasiSesi{}
	if sus, err := uc.aktivasi.ListSusulanByMahasiswa(userID); err == nil {
		for _, ps := range sus {
			if a, err := uc.aktivasi.FindSesiByID(ps.AktivasiSesiID); err == nil {
				susulanAktivasiBySesi[a.SesiPraktikumID] = a
			}
		}
	}

	items := make([]dto.SesiUserItem, 0, len(allSesi))
	for _, s := range allSesi {
		item := dto.SesiUserItem{
			SesiID:         s.ID,
			Judul:          s.JudulSesi,
			Deskripsi:      s.Deskripsi,
			Urutan:         s.Urutan,
			IsUjianPraktik: s.IsUjianPraktik,
		}
		aktivasi := aktivasiBySesi[s.ID]
		susulan := false
		if aktivasi == nil {
			if a, ok := susulanAktivasiBySesi[s.ID]; ok {
				aktivasi = a
				susulan = true
			}
		}
		if aktivasi != nil && aktivasi.IsActive {
			item.Aktif = true
			item.Susulan = susulan
			id := aktivasi.ID
			item.AktivasiSesiID = &id
			item.Courses = uc.buildCourses(userID, aktivasi)
		}
		items = append(items, item)
	}
	return items, nil
}

func (uc *PraktikumUsecase) buildCourses(userID int, aktivasi *entity.AktivasiSesi) []dto.CourseUserItem {
	courses := make([]dto.CourseUserItem, 0, len(aktivasi.AktivasiCourses))
	for _, ac := range aktivasi.AktivasiCourses {
		ci := dto.CourseUserItem{
			CourseID:         ac.CourseID,
			AktivasiCourseID: ac.ID,
			IsOpen:           ac.IsOpen,
			Status:           string(entity.StatusBelum),
		}
		if ac.Course != nil {
			ci.Jenis = string(ac.Course.Jenis)
			ci.Judul = ac.Course.Judul
			ci.DurasiMenit = ac.Course.DurasiMenit
		}
		if p, err := uc.pengerjaan.Find(userID, aktivasi.ID, ac.CourseID); err == nil {
			ci.Status = string(p.Status)
			ci.TotalNilai = p.TotalNilai
		}
		courses = append(courses, ci)
	}
	return courses
}

// Dashboard menyusun data dashboard mahasiswa.
func (uc *PraktikumUsecase) Dashboard(userID int) (*dto.DashboardUserResponse, error) {
	user, err := uc.users.FindByID(userID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	resp := &dto.DashboardUserResponse{
		Profil:       toUserResponse(user),
		SesiAktif:    make([]dto.SesiUserItem, 0),
		RiwayatNilai: make([]dto.NilaiCourseItem, 0),
	}
	if user.KelasID != nil && user.Shift != nil {
		if j, err := uc.jadwal.FindByKelasShift(*user.KelasID, *user.Shift); err == nil {
			resp.Jadwal = &dto.JadwalInfo{
				Hari:       j.Hari,
				JamMulai:   j.JamMulai,
				JamSelesai: j.JamSelesai,
				Keterangan: j.Keterangan,
			}
		}
	}

	allSesi, _ := uc.ListSesi(userID)
	for _, s := range allSesi {
		if s.Aktif {
			resp.SesiAktif = append(resp.SesiAktif, s)
		}
	}

	// Riwayat nilai dari pengerjaan yang sudah dinilai.
	if pengList, err := uc.pengerjaan.ListByMahasiswa(userID); err == nil {
		for _, p := range pengList {
			if p.TotalNilai == nil {
				continue
			}
			item := dto.NilaiCourseItem{Status: string(p.Status), TotalNilai: p.TotalNilai}
			if c, err := uc.course.FindByID(p.CourseID); err == nil {
				item.Jenis = string(c.Jenis)
				if sesi, err := uc.sesi.FindByID(c.SesiPraktikumID); err == nil {
					item.SesiJudul = sesi.JudulSesi
				}
			}
			resp.RiwayatNilai = append(resp.RiwayatNilai, item)
		}
	}

	return resp, nil
}
