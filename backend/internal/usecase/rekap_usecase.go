package usecase

import (
	"fmt"
	"strings"

	"lab-ap/internal/dto"
	"lab-ap/internal/repository"
)

type RekapUsecase struct {
	rekap repository.RekapRepository
	kelas repository.KelasRepository
}

func NewRekapUsecase(r repository.RekapRepository, k repository.KelasRepository) *RekapUsecase {
	return &RekapUsecase{rekap: r, kelas: k}
}

func (uc *RekapUsecase) GetRekapKelas(kelasID int) (*dto.RekapKelasResponse, error) {
	if _, err := uc.kelas.FindByID(kelasID); err != nil {
		return nil, ErrNotFound
	}

	rows, err := uc.rekap.GetRekapByKelas(kelasID)
	if err != nil {
		return nil, err
	}

	// 1. Ekstrak metadata kolom unik dan urutannya
	// Kita gunakan map untuk memastikan unique, dan slice untuk menjaga urutan penemuan.
	// Karena query sudah ORDER BY s.urutan, c.id, urutan penemuan pertama adalah urutan yang benar.
	type colMeta struct {
		key   string
		label string
	}
	var columns []colMeta
	colSeen := make(map[string]bool)

	// 2. Kelompokkan data per mahasiswa
	type mhsData struct {
		NIM        string
		Nama       string
		Scores     map[string]float64
		TotalScore float64
	}
	mhsMap := make(map[string]*mhsData) // Key by NIM
	var nimOrder []string               // Untuk menjaga urutan NIM

	for _, row := range rows {
		// Init mahasiswa jika belum ada
		if _, ok := mhsMap[row.NIM]; !ok {
			mhsMap[row.NIM] = &mhsData{
				NIM:        row.NIM,
				Nama:       row.Nama,
				Scores:     make(map[string]float64),
				TotalScore: 0,
			}
			nimOrder = append(nimOrder, row.NIM)
		}

		mhs := mhsMap[row.NIM]

		// Jika ada pengerjaan course
		if row.CourseID != nil && row.CourseJenis != nil && row.SesiJudul != nil {
			key := fmt.Sprintf("course_%d", *row.CourseID)
			
			// Buat label, misalnya "Modul 1 - pretest"
			// SesiJudul misal "Modul 1 - Pengenalan"
			// Kita ambil teks sebelum dash "-" jika ada untuk mempersingkat label.
			shortSesi := *row.SesiJudul
			if parts := strings.SplitN(*row.SesiJudul, "-", 2); len(parts) > 1 {
				shortSesi = strings.TrimSpace(parts[0])
			}
			label := fmt.Sprintf("%s - %s", shortSesi, *row.CourseJenis)

			if !colSeen[key] {
				colSeen[key] = true
				columns = append(columns, colMeta{key: key, label: label})
			}

			if row.TotalNilai != nil {
				mhs.Scores[key] = *row.TotalNilai
				mhs.TotalScore += *row.TotalNilai
			}
		}
	}

	// 3. Susun respons akhir
	resp := &dto.RekapKelasResponse{
		Columns: make([]dto.RekapKolom, 0, len(columns)),
		Data:    make([]dto.RekapMahasiswa, 0, len(nimOrder)),
	}

	for _, c := range columns {
		resp.Columns = append(resp.Columns, dto.RekapKolom{
			Key:   c.key,
			Label: c.label,
		})
	}

	for _, nim := range nimOrder {
		md := mhsMap[nim]
		resp.Data = append(resp.Data, dto.RekapMahasiswa{
			NIM:        md.NIM,
			Nama:       md.Nama,
			Scores:     md.Scores,
			TotalScore: md.TotalScore,
		})
	}

	return resp, nil
}
