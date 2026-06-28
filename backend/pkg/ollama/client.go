package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"lab-ap/config"
)

// Client mengatur koneksi ke API Ollama.
type Client struct {
	cfg        *config.Config
	httpClient *http.Client
}

// NewClient membuat instance Ollama client baru.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			// 50s: di bawah batas maxDuration Vercel Hobby (60s) agar fungsi sempat
			// membalas error rapi sebelum dimatikan paksa (hindari 504 mentah).
			Timeout: 50 * time.Second,
		},
	}
}

// ChatMessage merepresentasikan pesan dalam percakapan chat.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest adalah payload request ke endpoint /api/chat Ollama.
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Format   string        `json:"format,omitempty"` // "json" untuk memastikan output JSON
	Stream   bool          `json:"stream"`
}

// ChatResponse adalah response dari endpoint /api/chat.
type ChatResponse struct {
	Model     string      `json:"model"`
	CreatedAt time.Time   `json:"created_at"`
	Message   ChatMessage `json:"message"`
	Done      bool        `json:"done"`
}

// AIResult adalah struktur yang diharapkan dari balasan JSON AI.
type AIResult struct {
	Nilai    float64 `json:"nilai"`
	Feedback string  `json:"feedback"`
}

// GradeParams: konteks penilaian satu jawaban (untuk pemilihan rubrik & prompt).
type GradeParams struct {
	Soal        string
	Kunci       string  // kunci jawaban / instruksi
	Jawaban     string  // jawaban / kode mahasiswa
	Poin        float64 // poin maksimal soal
	JenisSoal   string  // "essay" / "coding"
	Difficulty  string  // "easy" / "medium" / "hard" / ""
	JenisCourse string  // "pretest" / "posttest" / "keterampilan" / "ujian_praktik"
}

type essayResp struct {
	Kategori string  `json:"kategori"`
	Poin     float64 `json:"poin"`
	Feedback string  `json:"feedback"`
}
type codingResp struct {
	SesuaiPetunjuk     float64 `json:"sesuai_petunjuk"`
	BerjalanTanpaError float64 `json:"berjalan_tanpa_error"`
	TepatWaktuSelesai  float64 `json:"tepat_waktu_selesai"`
	Feedback           string  `json:"feedback"`
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// GradeAnswer menilai satu jawaban memakai rubrik (default: ketentuan project lama).
// Hasil di-skala ke poin soal agar adaptif terhadap bobot soal yang ditetapkan admin.
func (c *Client) GradeAnswer(ctx context.Context, p GradeParams) (*AIResult, error) {
	// Soal & kunci berasal dari editor HTML (edra) -> bersihkan agar AI baca teks/kode bersih.
	soal := stripHTML(p.Soal)
	kunci := stripHTML(p.Kunci)
	coding := p.JenisSoal == "coding"
	var prompt string
	if coding {
		prompt = promptCoding(soal, kunci, p.Jawaban, subForCourse(p.JenisCourse))
	} else {
		prompt = promptEssay(soal, kunci, p.Jawaban, p.Difficulty, rubrikEssay(p.JenisCourse, p.Difficulty))
	}

	content, err := c.chat(ctx, prompt)
	if err != nil {
		return nil, err
	}

	res := &AIResult{}
	if coding {
		var r codingResp
		if err := json.Unmarshal([]byte(content), &r); err != nil {
			return nil, fmt.Errorf("parse AI coding JSON '%s': %w", content, err)
		}
		sc := subForCourse(p.JenisCourse)
		max := float64(sc.SesuaiPetunjukMax + sc.BteMax + sc.TwMax)
		sp := clamp(r.SesuaiPetunjuk, 0, float64(sc.SesuaiPetunjukMax))
		bte := clamp(r.BerjalanTanpaError, 0, float64(sc.BteMax))
		tw := clamp(r.TepatWaktuSelesai, float64(sc.TwMin), float64(sc.TwMax))
		res.Nilai = math.Round(clamp(p.Poin*(sp+bte+tw)/max, 0, p.Poin)) // nilai bulat (tanpa koma)
		res.Feedback = r.Feedback
	} else {
		var r essayResp
		if err := json.Unmarshal([]byte(content), &r); err != nil {
			return nil, fmt.Errorf("parse AI essay JSON '%s': %w", content, err)
		}
		kat := rubrikEssay(p.JenisCourse, p.Difficulty)
		maxCat := 0
		for _, v := range kat {
			if v > maxCat {
				maxCat = v
			}
		}
		catPoin, ok := kat[r.Kategori]
		if !ok {
			catPoin = int(r.Poin) // fallback ke poin AI bila kategori tak dikenal
		}
		if maxCat == 0 {
			maxCat = 1
		}
		res.Nilai = math.Round(clamp(p.Poin*float64(catPoin)/float64(maxCat), 0, p.Poin)) // nilai bulat (tanpa koma)
		res.Feedback = r.Feedback
	}
	return res, nil
}

// chat memanggil Ollama /api/chat dengan satu prompt dan mengembalikan isi balasan.
func (c *Client) chat(ctx context.Context, prompt string) (string, error) {
	reqPayload := ChatRequest{
		Model:    c.cfg.OllamaModel,
		Messages: []ChatMessage{{Role: "user", Content: prompt}},
		Format:   "json",
		Stream:   false,
	}
	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	url := strings.TrimSuffix(c.cfg.OllamaURL, "/") + "/api/chat"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.cfg.OllamaAPIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.cfg.OllamaAPIKey)
	}
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyErr, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(bodyErr))
	}
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return chatResp.Message.Content, nil
}
