package usecase

import (
	"errors"
	"time"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"

	"gorm.io/gorm"
)

type AktivasiUsecase struct {
	aktivasi   repository.AktivasiRepository
	sesi       repository.SesiRepository
	course     repository.CourseRepository
	kelas      repository.KelasRepository
	jawaban    repository.JawabanRepository
	pengerjaan repository.PengerjaanRepository
	soalUC     *SoalUsecase
}

func NewAktivasiUsecase(
	a repository.AktivasiRepository,
	s repository.SesiRepository,
	c repository.CourseRepository,
	k repository.KelasRepository,
	j repository.JawabanRepository,
	p repository.PengerjaanRepository,
	soalUC *SoalUsecase,
) *AktivasiUsecase {
	return &AktivasiUsecase{aktivasi: a, sesi: s, course: c, kelas: k, jawaban: j, pengerjaan: p, soalUC: soalUC}
}

func (uc *AktivasiUsecase) List() ([]entity.AktivasiSesi, error)       { return uc.aktivasi.ListSesi() }
func (uc *AktivasiUsecase) ListActive() ([]entity.AktivasiSesi, error) { return uc.aktivasi.ListActiveSesi() }

func (uc *AktivasiUsecase) Get(id int) (*entity.AktivasiSesi, error) {
	a, err := uc.aktivasi.FindSesiByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return a, nil
}

// Aktivasi mengaktifkan sesi untuk kelas+shift, melakukan gacha (untuk sesi normal),
// membuat baris aktivasi_course untuk course terpilih, lalu mengacak soal-nya.
func (uc *AktivasiUsecase) Aktivasi(req dto.AktivasiRequest) (*entity.AktivasiSesi, error) {
	sesi, err := uc.sesi.FindByID(req.SesiPraktikumID)
	if err != nil {
		return nil, ErrNotFound
	}
	if _, err := uc.kelas.FindByID(req.KelasID); err != nil {
		return nil, ErrNotFound
	}

	// Cek duplikat aktivasi.
	if _, err := uc.aktivasi.FindSesiByComposite(req.SesiPraktikumID, req.KelasID, req.Shift); err == nil {
		return nil, ErrConflict
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	courses, err := uc.course.ListBySesi(sesi.ID)
	if err != nil {
		return nil, err
	}

	// Tentukan course mana yang dibuatkan aktivasi_course.
	var dipakai []entity.Course
	if sesi.IsUjianPraktik {
		for _, c := range courses {
			if c.Jenis == entity.CourseUjianPraktik {
				dipakai = append(dipakai, c)
			}
		}
		if len(dipakai) == 0 {
			return nil, errors.New("sesi ujian praktik belum punya course ujian_praktik")
		}
	} else {
		// Gacha: pretest ATAU posttest (wajib pilih), + keterampilan.
		if req.GachaPilihan == "" {
			return nil, errors.New("gacha_pilihan (pretest/posttest) wajib untuk sesi normal")
		}
		pilihan := entity.JenisCourse(req.GachaPilihan)
		for _, c := range courses {
			if c.Jenis == pilihan || c.Jenis == entity.CourseKeterampilan {
				dipakai = append(dipakai, c)
			}
		}
		if len(dipakai) == 0 {
			return nil, errors.New("course terpilih tidak tersedia di sesi ini")
		}
	}

	// Buat aktivasi_sesi.
	aks := &entity.AktivasiSesi{
		SesiPraktikumID: req.SesiPraktikumID,
		KelasID:         req.KelasID,
		Shift:           req.Shift,
		IsActive:        true,
		ActivatedAt:     time.Now(),
	}
	if err := uc.aktivasi.CreateSesi(aks); err != nil {
		return nil, err
	}

	// Buat aktivasi_course (tertutup awalnya) + acak soal untuk tiap course.
	for i, c := range dipakai {
		course := c
		ac := &entity.AktivasiCourse{
			AktivasiSesiID: aks.ID,
			CourseID:       course.ID,
			IsOpen:         false,
			Urutan:         urutanCourse(course.Jenis, i),
		}
		if err := uc.aktivasi.CreateCourse(ac); err != nil {
			return nil, err
		}
		if err := uc.soalUC.AcakUntukAktivasiCourse(aks.ID, &course); err != nil {
			return nil, err
		}
	}

	return uc.aktivasi.FindSesiByID(aks.ID)
}

// urutanCourse menentukan urutan default buka course.
func urutanCourse(jenis entity.JenisCourse, fallback int) int {
	switch jenis {
	case entity.CoursePretest:
		return 1
	case entity.CourseKeterampilan:
		return 2
	case entity.CoursePosttest:
		return 3
	case entity.CourseUjianPraktik:
		return 1
	default:
		return fallback + 1
	}
}

// BukaTutupCourse mengubah is_open sebuah aktivasi_course.
// Saat menutup (is_open=false): auto-submit massal semua jawaban belum-submit
// dan tandai pengerjaan_course = selesai.
func (uc *AktivasiUsecase) BukaTutupCourse(req dto.BukaTutupCourseRequest) (*entity.AktivasiCourse, error) {
	ac, err := uc.aktivasi.FindCourseByID(req.AktivasiCourseID)
	if err != nil {
		return nil, ErrNotFound
	}
	now := time.Now()
	ac.IsOpen = req.IsOpen
	if req.IsOpen {
		ac.OpenedAt = &now
		ac.ClosedAt = nil
	} else {
		ac.ClosedAt = &now
	}
	if err := uc.aktivasi.UpdateCourse(ac); err != nil {
		return nil, err
	}

	if !req.IsOpen {
		// Auto-submit massal.
		if _, err := uc.jawaban.MarkSubmittedForCourse(ac.AktivasiSesiID, ac.CourseID); err != nil {
			return nil, err
		}
		if err := uc.pengerjaan.MarkSelesaiForCourse(ac.AktivasiSesiID, ac.CourseID); err != nil {
			return nil, err
		}
	}
	return ac, nil
}

// ---- Susulan ----

func (uc *AktivasiUsecase) AddSusulan(req dto.SusulanRequest) error {
	if _, err := uc.aktivasi.FindSesiByID(req.AktivasiSesiID); err != nil {
		return ErrNotFound
	}
	p := &entity.PesertaSusulan{
		AktivasiSesiID: req.AktivasiSesiID,
		MahasiswaID:    req.MahasiswaID,
		Alasan:         req.Alasan,
		CreatedAt:      time.Now(),
	}
	if err := uc.aktivasi.AddSusulan(p); err != nil {
		return ErrConflict
	}
	return nil
}

func (uc *AktivasiUsecase) RemoveSusulan(aktivasiSesiID, mahasiswaID int) error {
	return uc.aktivasi.RemoveSusulan(aktivasiSesiID, mahasiswaID)
}

func (uc *AktivasiUsecase) ListSusulan(aktivasiSesiID int) ([]entity.PesertaSusulan, error) {
	return uc.aktivasi.ListSusulan(aktivasiSesiID)
}
