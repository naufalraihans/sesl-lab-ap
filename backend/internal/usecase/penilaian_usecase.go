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
