package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// GradeAnswer memanggil Ollama untuk menilai jawaban.
func (c *Client) GradeAnswer(ctx context.Context, soal string, kunciJawaban string, jawabanMahasiswa string, poinMaksimal float64) (*AIResult, error) {
	promptSystem := fmt.Sprintf(`You are an objective and strict teaching assistant.
Your task is to grade a student's answer based on the provided Question (Soal) and Reference Answer (Kunci Jawaban).
The maximum score for this question is %.2f.
If the student's answer is perfectly correct or fully matches the essence of the Reference Answer, give the maximum score.
If it is partially correct, give partial credit proportional to the correctness.
If it is completely wrong or unrelated, give 0.

You must reply strictly in valid JSON format with exactly two keys:
1. "nilai": a number between 0 and %.2f representing the student's score.
2. "feedback": a brief string explaining why this score was given, written in Indonesian.

Do not include markdown blocks, explanation text outside JSON, or any other keys.`, poinMaksimal, poinMaksimal)

	kunciStr := kunciJawaban
	if kunciStr == "" {
		kunciStr = "(Tidak ada kunci jawaban, nilai berdasarkan kebenaran umum soal)"
	}

	promptUser := fmt.Sprintf("Soal:\n%s\n\nKunci Jawaban (Referensi):\n%s\n\nJawaban Mahasiswa:\n%s", soal, kunciStr, jawabanMahasiswa)

	reqPayload := ChatRequest{
		Model: c.cfg.OllamaModel,
		Messages: []ChatMessage{
			{Role: "system", Content: promptSystem},
			{Role: "user", Content: promptUser},
		},
		Format: "json", // Paksa output JSON (hanya didukung model tertentu di Ollama, namun ini standar aman)
		Stream: false,
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := strings.TrimSuffix(c.cfg.OllamaURL, "/") + "/api/chat"
	
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.cfg.OllamaAPIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.cfg.OllamaAPIKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyErr, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(bodyErr))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var result AIResult
	// Ollama terkadang masih memberikan spasi kosong berlebih atau format JSON sedikit kotor meskipun sudah mode json.
	// Namun dengan mode json native, seharusnya aman.
	if err := json.Unmarshal([]byte(chatResp.Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON response '%s': %w", chatResp.Message.Content, err)
	}

	// Sanity check
	if result.Nilai < 0 {
		result.Nilai = 0
	}
	if result.Nilai > poinMaksimal {
		result.Nilai = poinMaksimal
	}

	return &result, nil
}
