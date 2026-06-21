package usecase

import (
	"testing"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository/mocks"
)

// SetNilai harus menolak nilai di luar rentang [0, poin] SEBELUM menyentuh DB.
func TestSetNilaiValidasi(t *testing.T) {
	jm := &mocks.JawabanRepository{}
	j := &entity.JawabanMahasiswa{
		ID:          1,
		MahasiswaID: 5,
		SoalTerpilih: &entity.SoalTerpilih{
			AktivasiSesiID: 2, CourseID: 3,
			Soal: &entity.Soal{Poin: 10},
		},
	}
	jm.On("FindByID", 1).Return(j, nil)

	uc := &PenilaianUsecase{jawaban: jm} // tx tidak terpakai untuk input invalid

	if _, err := uc.SetNilai(dto.NilaiRequest{JawabanID: 1, Nilai: 11}); err == nil {
		t.Fatal("nilai > poin harus ditolak")
	}
	if _, err := uc.SetNilai(dto.NilaiRequest{JawabanID: 1, Nilai: -1}); err == nil {
		t.Fatal("nilai < 0 harus ditolak")
	}
}
