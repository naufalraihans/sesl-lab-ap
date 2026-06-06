package supabase

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client membungkus operasi Supabase Storage (REST API).
// Dipakai untuk upload foto asisten, PDF modul/pedoman, dan gambar flowchart.
type Client struct {
	baseURL    string
	serviceKey string
	bucket     string
	http       *http.Client
}

func New(baseURL, serviceKey, bucket string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		serviceKey: serviceKey,
		bucket:     bucket,
		http:       &http.Client{Timeout: 30 * time.Second},
	}
}

// Enabled menandakan apakah konfigurasi Supabase tersedia.
func (c *Client) Enabled() bool {
	return c.baseURL != "" && c.serviceKey != ""
}

// Upload mengunggah file ke path tertentu dalam bucket dan mengembalikan URL publik.
// upsert=true akan menimpa file yang sudah ada.
func (c *Client) Upload(path string, content []byte, contentType string, upsert bool) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("supabase belum dikonfigurasi")
	}
	path = strings.TrimLeft(path, "/")
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.baseURL, c.bucket, path)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(content))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.serviceKey)
	req.Header.Set("Content-Type", contentType)
	if upsert {
		req.Header.Set("x-upsert", "true")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("upload gagal (%d): %s", resp.StatusCode, string(body))
	}
	return c.PublicURL(path), nil
}

// PublicURL membangun URL publik untuk objek pada bucket public.
func (c *Client) PublicURL(path string) string {
	path = strings.TrimLeft(path, "/")
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", c.baseURL, c.bucket, path)
}

// Delete menghapus objek dari bucket.
func (c *Client) Delete(path string) error {
	if !c.Enabled() {
		return fmt.Errorf("supabase belum dikonfigurasi")
	}
	path = strings.TrimLeft(path, "/")
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.baseURL, c.bucket, path)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.serviceKey)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete gagal (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}
