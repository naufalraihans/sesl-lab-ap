package usecase

import (
	"context"

	"lab-ap/internal/dto"
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
}

func NewAIGradingUsecase(j repository.JawabanRepository, pu *PenilaianUsecase, oc *ollama.Client) AIGradingUsecase {
	return &aiGradingUsecase{jawabanRepo: j, penilaianUsecase: pu, ollamaClient: oc}
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

	res, err := uc.ollamaClient.GradeAnswer(
		context.Background(),
		j.SoalTerpilih.Soal.TeksSoal, kunci, j.JawabanTeks, j.SoalTerpilih.Soal.Poin,
	)
	if err != nil {
		return nil, err
	}

	feedback := "[AI] " + res.Feedback
	if _, err := uc.penilaianUsecase.SetNilai(dto.NilaiRequest{
		JawabanID: jawabanID,
		Nilai:     res.Nilai,
		Feedback:  &feedback,
	}); err != nil {
		return nil, err
	}

	return &dto.AIGradeOneResponse{JawabanID: jawabanID, Nilai: res.Nilai, Feedback: feedback}, nil
}
