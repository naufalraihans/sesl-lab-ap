package usecase

import (
	"fmt"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
)

type PenilaianUsecase struct {
	jawaban    repository.JawabanRepository
	pengerjaan repository.PengerjaanRepository
	users      repository.UserRepository
}

func NewPenilaianUsecase(j repository.JawabanRepository, p repository.PengerjaanRepository, u repository.UserRepository) *PenilaianUsecase {
	return &PenilaianUsecase{jawaban: j, pengerjaan: p, users: u}
}

// SetNilai memberi nilai (0..poin) + feedback pada satu jawaban, lalu recalc total_nilai course.
func (uc *PenilaianUsecase) SetNilai(req dto.NilaiRequest) (*entity.JawabanMahasiswa, error) {
	j, err := uc.jawaban.FindByID(req.JawabanID)
	if err != nil {
		return nil, ErrNotFound
	}
	if j.SoalTerpilih == nil || j.SoalTerpilih.Soal == nil {
		return nil, ErrBadRequest
	}
	poin := j.SoalTerpilih.Soal.Poin
	if req.Nilai < 0 || req.Nilai > poin {
		return nil, fmt.Errorf("%w: nilai harus 0..%.2f", ErrBadRequest, poin)
	}

	nilai := req.Nilai
	j.Nilai = &nilai
	j.Feedback = req.Feedback
	if err := uc.jawaban.Update(j); err != nil {
		return nil, err
	}

	// Recalc total_nilai pada pengerjaan_course.
	aktivasiSesiID := j.SoalTerpilih.AktivasiSesiID
	courseID := j.SoalTerpilih.CourseID
	total, err := uc.jawaban.SumNilai(j.MahasiswaID, aktivasiSesiID, courseID)
	if err != nil {
		return nil, err
	}
	p, err := uc.pengerjaan.FindOrCreate(j.MahasiswaID, aktivasiSesiID, courseID)
	if err != nil {
		return nil, err
	}
	p.TotalNilai = &total
	if err := uc.pengerjaan.Update(p); err != nil {
		return nil, err
	}
	return j, nil
}

// Rekap mengembalikan rekap jawaban semua mahasiswa untuk satu aktivasi+course.
func (uc *PenilaianUsecase) Rekap(aktivasiSesiID, courseID int) (*dto.RekapResponse, error) {
	jawaban, err := uc.jawaban.ListRekap(aktivasiSesiID, courseID)
	if err != nil {
		return nil, err
	}
	userCache := map[int]*entity.User{}
	resp := &dto.RekapResponse{AktivasiSesiID: aktivasiSesiID, CourseID: courseID}
	for _, j := range jawaban {
		u, ok := userCache[j.MahasiswaID]
		if !ok {
			u, _ = uc.users.FindByID(j.MahasiswaID)
			userCache[j.MahasiswaID] = u
		}
		item := dto.RekapItem{
			JawabanID:   j.ID,
			MahasiswaID: j.MahasiswaID,
			JawabanTeks: j.JawabanTeks,
			IsSubmitted: j.IsSubmitted,
			Nilai:       j.Nilai,
			Feedback:    j.Feedback,
		}
		if u != nil {
			item.NamaMhs = u.Nama
			item.NIM = u.NIM
		}
		if j.SoalTerpilih != nil && j.SoalTerpilih.Soal != nil {
			item.SoalID = j.SoalTerpilih.Soal.ID
			item.TeksSoal = j.SoalTerpilih.Soal.TeksSoal
			item.Poin = j.SoalTerpilih.Soal.Poin
		}
		resp.Items = append(resp.Items, item)
	}
	return resp, nil
}

// GetRekapJawabanGlobal mengembalikan rekap jawaban secara flat untuk dashboard
func (uc *PenilaianUsecase) GetRekapJawabanGlobal(kelasID, sesiID int, search, jenis string) (*dto.RekapJawabanResponse, error) {
	jawabanList, err := uc.jawaban.GetAllJawabanFlat(kelasID, sesiID, search, jenis)
	if err != nil {
		return nil, err
	}

	resp := &dto.RekapJawabanResponse{
		Items: make([]dto.RekapJawabanItem, 0, len(jawabanList)),
		Total: int64(len(jawabanList)),
	}

	for _, j := range jawabanList {
		item := dto.RekapJawabanItem{
			JawabanID:   j.ID,
			NIM:         j.Mahasiswa.NIM,
			JawabanTeks: j.JawabanTeks,
			IsSubmitted: j.IsSubmitted,
			WaktuSubmit: j.WaktuSubmit,
			Nilai:       j.Nilai,
			Feedback:    j.Feedback,
		}

		if j.Mahasiswa != nil {
			item.NamaMahasiswa = j.Mahasiswa.Nama
			if j.Mahasiswa.Kelas != nil {
				item.KelasID = j.Mahasiswa.Kelas.ID
				item.NamaKelas = j.Mahasiswa.Kelas.NamaKelas
			}
		}

		if j.SoalTerpilih != nil {
			if j.SoalTerpilih.Soal != nil {
				item.JenisSoal = string(j.SoalTerpilih.Soal.JenisSoal)
				item.TeksSoal = j.SoalTerpilih.Soal.TeksSoal
				item.PoinMaksimal = j.SoalTerpilih.Soal.Poin
			}
			if j.SoalTerpilih.Course != nil {
				item.CourseID = j.SoalTerpilih.Course.ID
				item.JudulCourse = j.SoalTerpilih.Course.Judul
				item.JenisCourse = string(j.SoalTerpilih.Course.Jenis)
			}
			if j.SoalTerpilih.AktivasiSesi != nil && j.SoalTerpilih.AktivasiSesi.Sesi != nil {
				item.SesiPraktikumID = j.SoalTerpilih.AktivasiSesi.Sesi.ID
				item.JudulSesi = j.SoalTerpilih.AktivasiSesi.Sesi.JudulSesi
			}
		}
		resp.Items = append(resp.Items, item)
	}

	return resp, nil
}

func (uc *PenilaianUsecase) syncPengerjaanCourseForJawaban(jawabanIDs []int) {
	if len(jawabanIDs) == 0 {
		return
	}
	
	// Cari semua jawaban_id untuk mengetahui mahasiswa mana yang berubah
	var jList []entity.JawabanMahasiswa
	for _, jID := range jawabanIDs {
		if j, err := uc.jawaban.FindByID(jID); err == nil && j != nil && j.SoalTerpilih != nil {
			jList = append(jList, *j)
		}
	}

	// Buat map unik untuk (MahasiswaID + AktivasiSesiID + CourseID)
	type pKey struct {
		MhsID, AktSesiID, CourseID int
	}
	affected := map[pKey]bool{}

	for _, j := range jList {
		k := pKey{MhsID: j.MahasiswaID, AktSesiID: j.SoalTerpilih.AktivasiSesiID, CourseID: j.SoalTerpilih.CourseID}
		affected[k] = true
	}

	// Untuk tiap kombinasi unik, kalkulasi ulang total_nilai
	for k := range affected {
		total, err := uc.jawaban.SumNilai(k.MhsID, k.AktSesiID, k.CourseID)
		if err == nil {
			if p, err := uc.pengerjaan.FindOrCreate(k.MhsID, k.AktSesiID, k.CourseID); err == nil {
				p.TotalNilai = &total
				uc.pengerjaan.Update(p)
			}
		}
	}
}

// BulkResetNilai mereset nilai dan feedback menjadi null
func (uc *PenilaianUsecase) BulkResetNilai(jawabanIDs []int) error {
	err := uc.jawaban.BulkResetNilai(jawabanIDs)
	if err != nil {
		return err
	}
	
	// Sinkronisasi nilai ke pengerjaan_course
	uc.syncPengerjaanCourseForJawaban(jawabanIDs)
	return nil
}

// BulkDeleteJawaban menghapus jawaban secara permanen
func (uc *PenilaianUsecase) BulkDeleteJawaban(jawabanIDs []int) error {
	// Ambil data dulu sblm dihapus buat sync
	var jList []entity.JawabanMahasiswa
	for _, jID := range jawabanIDs {
		if j, err := uc.jawaban.FindByID(jID); err == nil && j != nil {
			jList = append(jList, *j)
		}
	}

	err := uc.jawaban.BulkDelete(jawabanIDs)
	if err != nil {
		return err
	}

	// Ekstrak ID yg berhasil dihapus
	type pKey struct {
		MhsID, AktSesiID, CourseID int
	}
	affected := map[pKey]bool{}
	for _, j := range jList {
		if j.SoalTerpilih != nil {
			affected[pKey{MhsID: j.MahasiswaID, AktSesiID: j.SoalTerpilih.AktivasiSesiID, CourseID: j.SoalTerpilih.CourseID}] = true
		}
	}

	// Sync
	for k := range affected {
		total, err := uc.jawaban.SumNilai(k.MhsID, k.AktSesiID, k.CourseID)
		if err == nil {
			if p, err := uc.pengerjaan.FindOrCreate(k.MhsID, k.AktSesiID, k.CourseID); err == nil {
				p.TotalNilai = &total
				uc.pengerjaan.Update(p)
			}
		}
	}
	return nil
}

