package usecase

import (
	"context"
	"fmt"
	"sync"

	"lab-ap/internal/dto"
	"lab-ap/internal/repository"
	"lab-ap/pkg/ollama"

	"github.com/google/uuid"
)

// JobData menampung state job saat ini.
type JobData struct {
	ID        string `json:"id"`
	Status    string `json:"status"` // "queued", "processing", "completed", "failed"
	Total     int    `json:"total"`
	Processed int    `json:"processed"`
	Message   string `json:"message"`
}

type AIGradingUsecase interface {
	QueueJob(req dto.AIGradingBulkRequest) (*dto.AIGradingJobResponse, error)
	GetJobStatus(jobID string) (*dto.AIGradingJobResponse, error)
}

type aiGradingUsecase struct {
	jawabanRepo      repository.JawabanRepository
	penilaianUsecase *PenilaianUsecase // Digunakan untuk menyimpan nilai agar total_nilai ter-recalc
	ollamaClient     *ollama.Client

	jobQueue chan string      // Antrean Job ID
	jobs     sync.Map         // map[string]*JobData
	requests sync.Map         // map[string]dto.AIGradingBulkRequest
}

func NewAIGradingUsecase(j repository.JawabanRepository, pu *PenilaianUsecase, oc *ollama.Client) AIGradingUsecase {
	uc := &aiGradingUsecase{
		jawabanRepo:      j,
		penilaianUsecase: pu,
		ollamaClient:     oc,
		jobQueue:         make(chan string, 100), // Buffer antrean
	}

	// Mulai worker tunggal di latar belakang
	go uc.workerLoop()

	return uc
}

func (uc *aiGradingUsecase) QueueJob(req dto.AIGradingBulkRequest) (*dto.AIGradingJobResponse, error) {
	jobID := uuid.New().String()
	
	jobData := &JobData{
		ID:        jobID,
		Status:    "queued",
		Total:     0, // Akan dihitung saat processing dimulai
		Processed: 0,
		Message:   "Menunggu giliran",
	}

	uc.jobs.Store(jobID, jobData)
	uc.requests.Store(jobID, req)

	// Masukkan ke antrean
	uc.jobQueue <- jobID

	return &dto.AIGradingJobResponse{
		JobID:     jobData.ID,
		Status:    jobData.Status,
		Total:     jobData.Total,
		Processed: jobData.Processed,
		Message:   jobData.Message,
	}, nil
}

func (uc *aiGradingUsecase) GetJobStatus(jobID string) (*dto.AIGradingJobResponse, error) {
	val, ok := uc.jobs.Load(jobID)
	if !ok {
		return nil, ErrNotFound
	}
	jobData := val.(*JobData)
	return &dto.AIGradingJobResponse{
		JobID:     jobData.ID,
		Status:    jobData.Status,
		Total:     jobData.Total,
		Processed: jobData.Processed,
		Message:   jobData.Message,
	}, nil
}

func (uc *aiGradingUsecase) workerLoop() {
	for jobID := range uc.jobQueue {
		uc.processJob(jobID)
	}
}

func (uc *aiGradingUsecase) processJob(jobID string) {
	val, ok := uc.jobs.Load(jobID)
	if !ok {
		return
	}
	jobData := val.(*JobData)

	reqVal, ok := uc.requests.Load(jobID)
	if !ok {
		jobData.Status = "failed"
		jobData.Message = "Request data hilang"
		return
	}
	req := reqVal.(dto.AIGradingBulkRequest)

	jobData.Status = "processing"
	jobData.Message = "Mengambil data jawaban..."

	// Dapatkan semua jawaban untuk course ini
	allJawaban, err := uc.jawabanRepo.ListRekap(req.AktivasiSesiID, req.CourseID)
	if err != nil {
		jobData.Status = "failed"
		jobData.Message = fmt.Sprintf("Gagal fetch data: %v", err)
		return
	}

	// Filter jawaban yang is_submitted = true, teks_soal tidak kosong, dan nilai masih nil (Belum Dinilai)
	var targets []dto.RekapItem
	for _, j := range allJawaban {
		if j.IsSubmitted && j.Nilai == nil && j.JawabanTeks != "" && j.SoalTerpilih != nil && j.SoalTerpilih.Soal != nil {
			targets = append(targets, dto.RekapItem{
				JawabanID:   j.ID,
				TeksSoal:    j.SoalTerpilih.Soal.TeksSoal,
				Poin:        j.SoalTerpilih.Soal.Poin,
				JawabanTeks: j.JawabanTeks,
				// Kita masukkan KunciJawaban (jika ada) via casting/manual
			})
		}
	}

	jobData.Total = len(targets)
	jobData.Message = fmt.Sprintf("Menilai %d jawaban dengan AI...", len(targets))

	if len(targets) == 0 {
		jobData.Status = "completed"
		jobData.Message = "Selesai. Tidak ada jawaban baru yang perlu dinilai."
		return
	}

	// Proses satu persatu dengan AI
	ctx := context.Background()
	for _, target := range targets {
		// Dapatkan kunci jawaban dari DB (karena di ListRekap mungkin belum tertampung sempurna)
		jData, err := uc.jawabanRepo.FindByID(target.JawabanID)
		kunciStr := ""
		if err == nil && jData.SoalTerpilih != nil && jData.SoalTerpilih.Soal != nil && jData.SoalTerpilih.Soal.KunciJawaban != nil {
			kunciStr = *jData.SoalTerpilih.Soal.KunciJawaban
		}

		// Panggil Ollama
		aiResult, err := uc.ollamaClient.GradeAnswer(ctx, target.TeksSoal, kunciStr, target.JawabanTeks, target.Poin)
		
		if err == nil && aiResult != nil {
			// Simpan ke database via PenilaianUsecase agar total_nilai ter-recalc
			feedback := "[AI Graded] " + aiResult.Feedback
			uc.penilaianUsecase.SetNilai(dto.NilaiRequest{
				JawabanID: target.JawabanID,
				Nilai:     aiResult.Nilai,
				Feedback:  &feedback,
			})
		}

		jobData.Processed++
		// Sengaja tidak update Map terlalu agresif jika ingin menghindari blocking, namun karena kita simpan pointer, nilainya berubah otomatis
	}

	jobData.Status = "completed"
	jobData.Message = "Selesai menilai semua jawaban."
	// Hapus request memory setelah selesai
	uc.requests.Delete(jobID)
}
