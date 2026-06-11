// Package glot membungkus Glot.io Run API untuk menjalankan kode (C/Python)
// di sandbox eksternal. Backend hanya proxy; token disimpan di server.
package glot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

// New membuat client. baseURL default "https://glot.io/api/run".
func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http:    &http.Client{Timeout: 25 * time.Second},
	}
}

// Enabled menandakan token tersedia.
func (c *Client) Enabled() bool { return c.token != "" }

// Result adalah hasil eksekusi dari Glot.io.
type Result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}

type glotFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type glotRequest struct {
	Files   []glotFile `json:"files"`
	Stdin   string     `json:"stdin"`
	Command string     `json:"command"`
}

// fileName memetakan bahasa ke nama file & validitasnya.
func fileName(lang string) (string, bool) {
	switch lang {
	case "c":
		return "main.c", true
	case "python":
		return "main.py", true
	}
	return "", false
}

// Run mengeksekusi source untuk bahasa tertentu (c|python) dengan stdin opsional.
func (c *Client) Run(lang, source, stdin string) (*Result, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("glot belum dikonfigurasi")
	}
	name, ok := fileName(lang)
	if !ok {
		return nil, fmt.Errorf("bahasa tidak didukung: %s", lang)
	}

	payload, err := json.Marshal(glotRequest{
		Files: []glotFile{{Name: name, Content: source}},
		Stdin: stdin,
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/latest", c.baseURL, lang)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("glot status %d: %s", resp.StatusCode, string(body))
	}

	var r Result
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
