package handler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"lab-ap/pkg/response"
	"lab-ap/pkg/supabase"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	sb *supabase.Client
}

func NewUploadHandler(sb *supabase.Client) *UploadHandler { return &UploadHandler{sb: sb} }

// Upload POST /api/admin/upload (multipart: file, folder)
// @Summary Upload File
// @Description Mengunggah file ke Supabase Storage (gambar soal, dll)
// @Tags Admin - Upload
// @Security bearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File yang diupload"
// @Param folder formData string false "Nama folder tujuan"
// @Success 201 {object} response.Envelope
// @Router /admin/upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	if !h.sb.Enabled() {
		response.Fail(c, http.StatusServiceUnavailable, "Supabase Storage belum dikonfigurasi", nil)
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "File wajib diunggah (field 'file')", err.Error())
		return
	}
	folder := strings.Trim(c.PostForm("folder"), "/")
	if folder == "" {
		folder = "uploads"
	}

	f, err := fileHeader.Open()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "Gagal membuka file", err.Error())
		return
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "Gagal membaca file", err.Error())
		return
	}

	ext := filepath.Ext(fileHeader.Filename)
	name := fmt.Sprintf("%s/%d%s", folder, time.Now().UnixNano(), ext)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	url, err := h.sb.Upload(name, content, contentType, true)
	if err != nil {
		response.Fail(c, http.StatusBadGateway, "Upload ke Supabase gagal", err.Error())
		return
	}
	response.Created(c, "File diunggah", gin.H{"url": url, "path": name})
}
