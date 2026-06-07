package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	"lab-ap/internal/entity"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type KonfigurasiHandler struct {
	uc *usecase.KonfigurasiUsecase
}

func NewKonfigurasiHandler(uc *usecase.KonfigurasiUsecase) *KonfigurasiHandler {
	return &KonfigurasiHandler{uc: uc}
}

// All GET /api/admin/konfigurasi
// @Summary Seluruh Konfigurasi
// @Description Mengambil semua data konfigurasi global
// @Tags Admin - Konfigurasi
// @Security bearerAuth
// @Produce json
// @Success 200 {object} response.Envelope{data=[]entity.Konfigurasi}
// @Router /admin/konfigurasi [get]
func (h *KonfigurasiHandler) All(c *gin.Context) {
	res, err := h.uc.All()
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Konfigurasi", res)
}

// Set POST /api/admin/konfigurasi
// @Summary Set Konfigurasi
// @Description Mengubah atau menambah konfigurasi global (Upsert)
// @Tags Admin - Konfigurasi
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param request body dto.KonfigurasiRequest true "Payload Konfigurasi"
// @Success 200 {object} response.Envelope
// @Router /admin/konfigurasi [post]
func (h *KonfigurasiHandler) Set(c *gin.Context) {
	var req dto.KonfigurasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	if err := h.uc.Set(req.Key, req.Value); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Konfigurasi disimpan", nil)
}

// ---- Publik (info) ----

// PublicModul GET /api/info/modul
// @Summary File Modul Publik
// @Description Mengambil URL publik dari Modul Praktikum
// @Tags Info
// @Produce json
// @Success 200 {object} response.Envelope
// @Router /info/modul [get]
func (h *KonfigurasiHandler) PublicModul(c *gin.Context) {
	url, _ := h.uc.Get(entity.KeyModulFileURL)
	response.OK(c, http.StatusOK, "Modul praktikum", gin.H{"file_url": url})
}

// PublicJadwalConfig GET /api/info/jadwal/config
// @Summary Konfigurasi Jadwal Publik
// @Description Mengecek apakah jadwal menggunakan sistem internal atau GDrive URL
// @Tags Info
// @Produce json
// @Success 200 {object} response.Envelope
// @Router /info/jadwal/config [get]
func (h *KonfigurasiHandler) PublicJadwalConfig(c *gin.Context) {
	mode, _ := h.uc.Get(entity.KeyJadwalMode)
	gdrive, _ := h.uc.Get(entity.KeyGDriveJadwalURL)
	if mode == "" {
		mode = "internal"
	}
	response.OK(c, http.StatusOK, "Konfigurasi jadwal", gin.H{"mode": mode, "gdrive_url": gdrive})
}
