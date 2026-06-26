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

// SetKeaktifanBulk sukses utk pretest/posttest (RowsAffected>0), tolak bila 0.
func TestSetKeaktifanBulk(t *testing.T) {
	pm := &mocks.PengerjaanRepository{}
	pm.On("SetKeaktifanIfTest", 5, 10.0).Return(int64(1), nil) // pretest/posttest
	pm.On("SetKeaktifanIfTest", 9, 7.0).Return(int64(0), nil)  // bukan / tak ada
	uc := &PenilaianUsecase{pengerjaan: pm}

	if err := uc.SetKeaktifanBulk([]dto.KeaktifanItem{{PengerjaanID: 5, Nilai: 10}}); err != nil {
		t.Fatalf("valid harus sukses: %v", err)
	}
	if err := uc.SetKeaktifanBulk([]dto.KeaktifanItem{{PengerjaanID: 9, Nilai: 7}}); err == nil {
		t.Fatal("RowsAffected 0 harus error")
	}
}
