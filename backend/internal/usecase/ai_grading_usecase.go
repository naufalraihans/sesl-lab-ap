package usecase

import (
	"context"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/repository"
	"lab-ap/pkg/ollama"
)

// AIGradingUsecase: penilaian AI SINKRON satu-per-satu (serverless-friendly).
// Frontend memanggil GradeOne berulang per jawaban; tidak ada worker/job background.
type AIGradingUsecase interface {
	ListTargets(aktivasiSesiID, courseID int) (*dto.AIGradeTargetsResponse, error)
	GradeOne(jawabanID int) (*dto.AIGradeOneResponse, error)
}

type aiGradingUsecase struct {
	jawabanRepo      repository.JawabanRepository
	penilaianUsecase *PenilaianUsecase // untuk SetNilai (recalc total_nilai)
	ollamaClient     *ollama.Client
	konfRepo         repository.KonfigurasiRepository // untuk override model AI (key ai_model)
}

func NewAIGradingUsecase(j repository.JawabanRepository, pu *PenilaianUsecase, oc *ollama.Client, kr repository.KonfigurasiRepository) AIGradingUsecase {
	return &aiGradingUsecase{jawabanRepo: j, penilaianUsecase: pu, ollamaClient: oc, konfRepo: kr}
}

// ListTargets mengembalikan jawaban_id yang perlu dinilai AI untuk satu course:
// sudah submit, belum dinilai (nilai null), dan ada teks jawaban.
func (uc *aiGradingUsecase) ListTargets(aktivasiSesiID, courseID int) (*dto.AIGradeTargetsResponse, error) {
	all, err := uc.jawabanRepo.ListRekap(aktivasiSesiID, courseID)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0)
	for _, j := range all {
		if j.IsSubmitted && j.Nilai == nil && j.JawabanTeks != "" &&
			j.SoalTerpilih != nil && j.SoalTerpilih.Soal != nil {
			ids = append(ids, j.ID)
		}
	}
	return &dto.AIGradeTargetsResponse{JawabanIDs: ids, Total: len(ids)}, nil
}

// GradeOne menilai SATU jawaban dengan AI lalu menyimpannya (sinkron).
func (uc *aiGradingUsecase) GradeOne(jawabanID int) (*dto.AIGradeOneResponse, error) {
	j, err := uc.jawabanRepo.FindByID(jawabanID)
	if err != nil {
		return nil, ErrNotFound
	}
	if j.SoalTerpilih == nil || j.SoalTerpilih.Soal == nil {
		return nil, ErrBadRequest
	}
	if !j.IsSubmitted || j.JawabanTeks == "" {
		return nil, ErrBadRequest
	}

	kunci := ""
	if j.SoalTerpilih.Soal.KunciJawaban != nil {
		kunci = *j.SoalTerpilih.Soal.KunciJawaban
	}
	difficulty := ""
	if j.SoalTerpilih.Soal.Difficulty != nil {
		difficulty = string(*j.SoalTerpilih.Soal.Difficulty)
	}
	jenisCourse := ""
	if j.SoalTerpilih.Course != nil {
		jenisCourse = string(j.SoalTerpilih.Course.Jenis)
	}

	// Override model AI dari konfigurasi (key ai_model); kosong = default env.
	model := ""
	if k, err := uc.konfRepo.Get(entity.KeyAIModel); err == nil && k != nil {
		model = k.Value
	}

	res, err := uc.ollamaClient.GradeAnswer(context.Background(), ollama.GradeParams{
		Soal:        j.SoalTerpilih.Soal.TeksSoal,
		Kunci:       kunci,
		Jawaban:     j.JawabanTeks,
		Poin:        j.SoalTerpilih.Soal.Poin,
		JenisSoal:   string(j.SoalTerpilih.Soal.JenisSoal),
		Difficulty:  difficulty,
		JenisCourse: jenisCourse,
		Model:       model,
	})
	if err != nil {
		return nil, err
	}

	// Clamp nilai AI ke [0, poin] supaya tidak ditolak SetNilai bila model meleset.
	nilai := res.Nilai
	if nilai < 0 {
		nilai = 0
	} else if poin := j.SoalTerpilih.Soal.Poin; nilai > poin {
		nilai = poin
	}

	feedback := "[AI] " + res.Feedback
	if _, err := uc.penilaianUsecase.SetNilai(dto.NilaiRequest{
		JawabanID: jawabanID,
		Nilai:     nilai,
		Feedback:  &feedback,
	}); err != nil {
		return nil, err
	}

	return &dto.AIGradeOneResponse{JawabanID: jawabanID, Nilai: nilai, Feedback: feedback}, nil
}
