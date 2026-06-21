package usecase

import (
	"testing"
	"time"

	"lab-ap/internal/entity"
	"lab-ap/internal/repository/mocks"
)

func ptrInt(n int) *int { return &n }

func TestDeadlineFrom(t *testing.T) {
	if deadlineFrom(nil, 30) != nil {
		t.Fatal("waktu_mulai nil harus menghasilkan deadline nil")
	}
	mulai := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	dl := deadlineFrom(&mulai, 30)
	if dl == nil || !dl.Equal(mulai.Add(30*time.Minute)) {
		t.Fatalf("deadline harus mulai+30m, dapat %v", dl)
	}
}

func TestCekAkses(t *testing.T) {
	akt := &entity.AktivasiSesi{ID: 9, KelasID: 2, Shift: 1}

	// Kelas + shift cocok → boleh, tanpa menyentuh repo.
	uc := &JawabanUsecase{}
	user := &entity.User{ID: 7, KelasID: ptrInt(2), Shift: ptrInt(1)}
	if ok, err := uc.cekAkses(user, akt); err != nil || !ok {
		t.Fatalf("kelas+shift cocok harus boleh; ok=%v err=%v", ok, err)
	}

	// Tidak cocok tapi terdaftar susulan → boleh.
	mYa := &mocks.AktivasiRepository{}
	mYa.On("IsSusulan", 9, 7).Return(true, nil)
	ucSus := &JawabanUsecase{aktivasi: mYa}
	beda := &entity.User{ID: 7, KelasID: ptrInt(3), Shift: ptrInt(2)}
	if ok, err := ucSus.cekAkses(beda, akt); err != nil || !ok {
		t.Fatalf("susulan harus boleh; ok=%v err=%v", ok, err)
	}

	// Tidak cocok & bukan susulan → ditolak.
	mTidak := &mocks.AktivasiRepository{}
	mTidak.On("IsSusulan", 9, 7).Return(false, nil)
	ucTolak := &JawabanUsecase{aktivasi: mTidak}
	if ok, err := ucTolak.cekAkses(beda, akt); err != nil || ok {
		t.Fatalf("bukan kelas/shift & bukan susulan harus ditolak; ok=%v err=%v", ok, err)
	}

	// Shift nil (kasus data migrasi) → tidak cocok, jatuh ke cek susulan.
	mNil := &mocks.AktivasiRepository{}
	mNil.On("IsSusulan", 9, 7).Return(false, nil)
	ucNil := &JawabanUsecase{aktivasi: mNil}
	noShift := &entity.User{ID: 7, KelasID: ptrInt(2), Shift: nil}
	if ok, _ := ucNil.cekAkses(noShift, akt); ok {
		t.Fatal("shift nil tidak boleh otomatis lolos")
	}
}
