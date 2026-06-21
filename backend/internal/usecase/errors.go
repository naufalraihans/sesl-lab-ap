package usecase

import (
	"errors"
	"strings"
)

// Error domain yang dipetakan ke kode HTTP oleh handler.
var (
	ErrNotFound      = errors.New("data tidak ditemukan")
	ErrUnauthorized  = errors.New("tidak terautentikasi")
	ErrForbidden     = errors.New("akses ditolak")
	ErrBadRequest    = errors.New("permintaan tidak valid")
	ErrConflict      = errors.New("konflik data")
	ErrRegisterClosed = errors.New("akses register belum dibuka")
	ErrAlreadyDone   = errors.New("course sudah dikerjakan / ditutup")
	ErrTimeUp        = errors.New("waktu pengerjaan telah habis")
	ErrHasReferences = errors.New("tidak bisa dihapus: masih dipakai data lain (mis. ada nilai/jawaban mahasiswa)")
)

// mapDeleteErr menerjemahkan pelanggaran foreign key (SQLSTATE 23503) menjadi
// ErrHasReferences agar admin dapat pesan jelas, bukan 500. Lihat bugAudit B8.
func mapDeleteErr(err error) error {
	if err != nil && strings.Contains(err.Error(), "SQLSTATE 23503") {
		return ErrHasReferences
	}
	return err
}
