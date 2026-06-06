package usecase

import (
	"errors"
	"time"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"

	"gorm.io/gorm"
)

type JawabanUsecase struct {
	aktivasi   repository.AktivasiRepository
	terpilih   repository.SoalTerpilihRepository
	jawaban    repository.JawabanRepository
	pengerjaan repository.PengerjaanRepository
	users      repository.UserRepository
	course     repository.CourseRepository
}

func NewJawabanUsecase(
	a repository.AktivasiRepository,
	t repository.SoalTerpilihRepository,
	j repository.JawabanRepository,
	p repository.PengerjaanRepository,
	u repository.UserRepository,
	c repository.CourseRepository,
) *JawabanUsecase {
	return &JawabanUsecase{aktivasi: a, terpilih: t, jawaban: j, pengerjaan: p, users: u, course: c}
}

// cekAkses memverifikasi mahasiswa boleh mengakses aktivasi (kelas+shift cocok ATAU susulan).
func (uc *JawabanUsecase) cekAkses(user *entity.User, aktivasi *entity.AktivasiSesi) (bool, error) {
	cocok := user.KelasID != nil && user.Shift != nil &&
		*user.KelasID == aktivasi.KelasID && *user.Shift == aktivasi.Shift
	if cocok {
		return true, nil
	}
	return uc.aktivasi.IsSusulan(aktivasi.ID, user.ID)
}

// loadKonteks mengambil user, aktivasi, aktivasi_course, dan course.
func (uc *JawabanUsecase) loadKonteks(userID, aktivasiSesiID, courseID int) (*entity.User, *entity.AktivasiSesi, *entity.AktivasiCourse, *entity.Course, error) {
	user, err := uc.users.FindByID(userID)
	if err != nil {
		return nil, nil, nil, nil, ErrUnauthorized
	}
	aktivasi, err := uc.aktivasi.FindSesiByID(aktivasiSesiID)
	if err != nil {
		return nil, nil, nil, nil, ErrNotFound
	}
	ac, err := uc.aktivasi.FindCourse(aktivasiSesiID, courseID)
	if err != nil {
		return nil, nil, nil, nil, ErrForbidden // course tidak diaktifkan untuk aktivasi ini
	}
	course, err := uc.course.FindByID(courseID)
	if err != nil {
		return nil, nil, nil, nil, ErrNotFound
	}
	ok, err := uc.cekAkses(user, aktivasi)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if !ok {
		return nil, nil, nil, nil, ErrForbidden
	}
	return user, aktivasi, ac, course, nil
}

func deadlineFrom(mulai *time.Time, durasiMenit int) *time.Time {
	if mulai == nil {
		return nil
	}
	d := mulai.Add(time.Duration(durasiMenit) * time.Minute)
	return &d
}

// GetRuang mengembalikan data ruang pengerjaan (soal + jawaban tersimpan + timer).
// Read-only: tidak memulai timer.
func (uc *JawabanUsecase) GetRuang(userID, aktivasiSesiID, courseID int) (*dto.RuangCourseResponse, error) {
	_, _, ac, course, err := uc.loadKonteks(userID, aktivasiSesiID, courseID)
	if err != nil {
		return nil, err
	}
	return uc.buildRuang(userID, aktivasiSesiID, ac, course)
}

// Mulai menandai mahasiswa mulai mengerjakan (set waktu_mulai sekali).
func (uc *JawabanUsecase) Mulai(userID, aktivasiSesiID, courseID int) (*dto.RuangCourseResponse, error) {
	_, _, ac, course, err := uc.loadKonteks(userID, aktivasiSesiID, courseID)
	if err != nil {
		return nil, err
	}
	if !ac.IsOpen {
		return nil, ErrAlreadyDone
	}
	p, err := uc.pengerjaan.FindOrCreate(userID, aktivasiSesiID, courseID)
	if err != nil {
		return nil, err
	}
	if p.Status == entity.StatusSelesai {
		return nil, ErrAlreadyDone
	}
	if p.WaktuMulai == nil {
		now := time.Now()
		p.WaktuMulai = &now
		p.Status = entity.StatusSedang
		if err := uc.pengerjaan.Update(p); err != nil {
			return nil, err
		}
	}
	return uc.buildRuang(userID, aktivasiSesiID, ac, course)
}

// AutoSave menyimpan jawaban satu soal (validasi server: open, belum lewat deadline, belum submit).
func (uc *JawabanUsecase) AutoSave(userID int, req dto.AutoSaveRequest) error {
	st, err := uc.terpilih.FindByID(req.SoalTerpilihID)
	if err != nil {
		return ErrNotFound
	}
	user, _, ac, course, err := uc.loadKonteks(userID, st.AktivasiSesiID, st.CourseID)
	if err != nil {
		return err
	}
	_ = user
	if !ac.IsOpen {
		return ErrAlreadyDone
	}
	// Validasi deadline (server-authoritative).
	p, err := uc.pengerjaan.Find(userID, st.AktivasiSesiID, st.CourseID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		// belum mulai → mulai otomatis saat save pertama
		p, err = uc.pengerjaan.FindOrCreate(userID, st.AktivasiSesiID, st.CourseID)
		if err != nil {
			return err
		}
		now := time.Now()
		p.WaktuMulai = &now
		p.Status = entity.StatusSedang
		_ = uc.pengerjaan.Update(p)
	}
	if p.Status == entity.StatusSelesai {
		return ErrAlreadyDone
	}
	if dl := deadlineFrom(p.WaktuMulai, course.DurasiMenit); dl != nil && time.Now().After(*dl) {
		// Lewat deadline → auto-submit & tolak.
		_ = uc.finalize(userID, st.AktivasiSesiID, st.CourseID)
		return ErrTimeUp
	}

	// Upsert jawaban.
	j, err := uc.jawaban.FindByMahasiswaSoal(userID, st.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		j = &entity.JawabanMahasiswa{
			MahasiswaID:    userID,
			SoalTerpilihID: st.ID,
			JawabanTeks:    req.JawabanTeks,
			UpdatedAt:      time.Now(),
		}
		return uc.jawaban.Create(j)
	}
	if j.IsSubmitted {
		return ErrAlreadyDone
	}
	j.JawabanTeks = req.JawabanTeks
	j.UpdatedAt = time.Now()
	return uc.jawaban.Update(j)
}

// Submit melakukan submit manual seluruh course untuk mahasiswa.
func (uc *JawabanUsecase) Submit(userID int, req dto.SubmitRequest) error {
	_, _, _, _, err := uc.loadKonteks(userID, req.AktivasiSesiID, req.CourseID)
	if err != nil {
		return err
	}
	return uc.finalize(userID, req.AktivasiSesiID, req.CourseID)
}

// finalize menandai semua jawaban mahasiswa submitted & pengerjaan selesai.
func (uc *JawabanUsecase) finalize(userID, aktivasiSesiID, courseID int) error {
	if _, err := uc.jawaban.MarkSubmittedForMahasiswaCourse(userID, aktivasiSesiID, courseID); err != nil {
		return err
	}
	p, err := uc.pengerjaan.FindOrCreate(userID, aktivasiSesiID, courseID)
	if err != nil {
		return err
	}
	now := time.Now()
	p.Status = entity.StatusSelesai
	p.WaktuSelesai = &now
	return uc.pengerjaan.Update(p)
}

// AutoSubmitExpired disapu oleh background job: auto-submit pengerjaan yang lewat deadline.
func (uc *JawabanUsecase) AutoSubmitExpired() (int, error) {
	expired, err := uc.pengerjaan.FindExpired()
	if err != nil {
		return 0, err
	}
	n := 0
	for _, e := range expired {
		if err := uc.finalize(e.MahasiswaID, e.AktivasiSesiID, e.CourseID); err != nil {
			continue
		}
		n++
	}
	return n, nil
}

// buildRuang menyusun response ruang course.
func (uc *JawabanUsecase) buildRuang(userID, aktivasiSesiID int, ac *entity.AktivasiCourse, course *entity.Course) (*dto.RuangCourseResponse, error) {
	soalTerpilih, err := uc.terpilih.ListByAktivasiCourse(aktivasiSesiID, course.ID)
	if err != nil {
		return nil, err
	}
	jawabanList, err := uc.jawaban.ListByMahasiswaCourse(userID, aktivasiSesiID, course.ID)
	if err != nil {
		return nil, err
	}
	jawabanMap := make(map[int]entity.JawabanMahasiswa, len(jawabanList))
	for _, j := range jawabanList {
		jawabanMap[j.SoalTerpilihID] = j
	}

	resp := &dto.RuangCourseResponse{
		AktivasiSesiID: aktivasiSesiID,
		CourseID:       course.ID,
		Jenis:          string(course.Jenis),
		DurasiMenit:    course.DurasiMenit,
		IsOpen:         ac.IsOpen,
		Status:         string(entity.StatusBelum),
	}

	if p, err := uc.pengerjaan.Find(userID, aktivasiSesiID, course.ID); err == nil {
		resp.Status = string(p.Status)
		resp.WaktuMulai = p.WaktuMulai
		resp.Deadline = deadlineFrom(p.WaktuMulai, course.DurasiMenit)
	}

	for _, st := range soalTerpilih {
		item := dto.SoalTampilResponse{
			SoalTerpilihID: st.ID,
			Urutan:         st.Urutan,
		}
		if st.Soal != nil {
			item.JenisSoal = string(st.Soal.JenisSoal)
			item.TeksSoal = st.Soal.TeksSoal
			item.GambarURL = st.Soal.GambarURL
			item.Poin = st.Soal.Poin
			if st.Soal.KategoriUjian != nil {
				k := string(*st.Soal.KategoriUjian)
				item.KategoriUjian = &k
			}
		}
		if j, ok := jawabanMap[st.ID]; ok {
			item.JawabanTeks = j.JawabanTeks
			item.IsSubmitted = j.IsSubmitted
		}
		resp.Soal = append(resp.Soal, item)
	}
	return resp, nil
}
