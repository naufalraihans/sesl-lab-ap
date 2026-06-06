package usecase

import "errors"

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
)
