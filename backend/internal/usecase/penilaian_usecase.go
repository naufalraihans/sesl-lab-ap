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
	tx         *repository.PenilaianTxRepo
}

func NewPenilaianUsecase(j repository.JawabanRepository, p repository.PengerjaanRepository, u repository.UserRepository, tx *repository.PenilaianTxRepo) *PenilaianUsecase {
	return &PenilaianUsecase{jawaban: j, pengerjaan: p, users: u, tx: tx}
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

	// Simpan nilai + recalc total_nilai pengerjaan_course secara ATOMIK (1 transaksi).
	key := repository.PengerjaanKey{
		MahasiswaID:    j.MahasiswaID,
		AktivasiSesiID: j.SoalTerpilih.AktivasiSesiID,
		CourseID:       j.SoalTerpilih.CourseID,
	}
	if err := uc.tx.SetNilaiAndRecalc(j.ID, nilai, req.Feedback, key); err != nil {
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

// BulkResetNilai mereset nilai & feedback menjadi null, lalu recalc total_nilai
// pengerjaan_course terdampak — seluruhnya ATOMIK dalam satu transaksi.
func (uc *PenilaianUsecase) BulkResetNilai(jawabanIDs []int) error {
	return uc.tx.BulkResetAndRecalc(jawabanIDs)
}

// BulkDeleteJawaban menghapus jawaban permanen, lalu recalc total_nilai
// pengerjaan_course terdampak — seluruhnya ATOMIK dalam satu transaksi.
func (uc *PenilaianUsecase) BulkDeleteJawaban(jawabanIDs []int) error {
	return uc.tx.BulkDeleteAndRecalc(jawabanIDs)
}

