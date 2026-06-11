package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"lab-ap/internal/dto"
	"lab-ap/internal/repository"
	"lab-ap/pkg/ollama"

	"github.com/google/uuid"
)

// jobRetention adalah lama job disimpan di memori setelah selesai/gagal
// sebelum dibersihkan otomatis (mencegah kebocoran memori map jobs).
const jobRetention = 10 * time.Minute

// JobData menampung state job saat ini.
// Akses field dilindungi mu karena ditulis worker goroutine & dibaca handler HTTP.
type JobData struct {
	mu        sync.Mutex
	ID        string `json:"id"`
	Status    string `json:"status"` // "queued", "processing", "completed", "failed"
	Total     int    `json:"total"`
	Processed int    `json:"processed"`
	Message   string `json:"message"`
}

// set memperbarui status & message secara aman.
func (j *JobData) set(status, message string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if status != "" {
		j.Status = status
	}
	j.Message = message
}

// setTotal menyetel total target secara aman.
func (j *JobData) setTotal(total int) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Total = total
}

// incProcessed menambah counter progress secara aman.
func (j *JobData) incProcessed() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Processed++
}

// snapshot membaca seluruh field secara aman untuk response.
func (j *JobData) snapshot() dto.AIGradingJobResponse {
	j.mu.Lock()
	defer j.mu.Unlock()
	return dto.AIGradingJobResponse{
		JobID:     j.ID,
		Status:    j.Status,
		Total:     j.Total,
		Processed: j.Processed,
		Message:   j.Message,
	}
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

	snap := jobData.snapshot()
	return &snap, nil
}

func (uc *aiGradingUsecase) GetJobStatus(jobID string) (*dto.AIGradingJobResponse, error) {
	val, ok := uc.jobs.Load(jobID)
	if !ok {
		return nil, ErrNotFound
	}
	snap := val.(*JobData).snapshot()
	return &snap, nil
}

// scheduleCleanup menghapus job dari memori setelah masa retensi
// agar map jobs tidak tumbuh tanpa batas.
func (uc *aiGradingUsecase) scheduleCleanup(jobID string) {
	uc.requests.Delete(jobID)
	time.AfterFunc(jobRetention, func() {
		uc.jobs.Delete(jobID)
	})
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
		jobData.set("failed", "Request data hilang")
		uc.scheduleCleanup(jobID)
		return
	}
	req := reqVal.(dto.AIGradingBulkRequest)

	jobData.set("processing", "Mengambil data jawaban...")

	// Dapatkan semua jawaban untuk course ini
	allJawaban, err := uc.jawabanRepo.ListRekap(req.AktivasiSesiID, req.CourseID)
	if err != nil {
		jobData.set("failed", fmt.Sprintf("Gagal fetch data: %v", err))
		uc.scheduleCleanup(jobID)
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

	jobData.setTotal(len(targets))
	jobData.set("", fmt.Sprintf("Menilai %d jawaban dengan AI...", len(targets)))

	if len(targets) == 0 {
		jobData.set("completed", "Selesai. Tidak ada jawaban baru yang perlu dinilai.")
		uc.scheduleCleanup(jobID)
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

		jobData.incProcessed()
	}

	jobData.set("completed", "Selesai menilai semua jawaban.")
	// Hapus request & jadwalkan pembersihan job dari memori.
	uc.scheduleCleanup(jobID)
}
