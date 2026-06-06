package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type AktivasiHandler struct {
	uc *usecase.AktivasiUsecase
}

func NewAktivasiHandler(uc *usecase.AktivasiUsecase) *AktivasiHandler { return &AktivasiHandler{uc: uc} }

func (h *AktivasiHandler) List(c *gin.Context) {
	res, err := h.uc.List()
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar aktivasi", res)
}

func (h *AktivasiHandler) Get(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	res, err := h.uc.Get(id)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Detail aktivasi", res)
}

// Aktivasi POST /api/admin/aktivasi (aktifkan sesi + gacha)
func (h *AktivasiHandler) Aktivasi(c *gin.Context) {
	var req dto.AktivasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.Aktivasi(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Sesi diaktifkan", res)
}

// BukaTutupCourse POST /api/admin/aktivasi/course/buka-tutup
func (h *AktivasiHandler) BukaTutupCourse(c *gin.Context) {
	var req dto.BukaTutupCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.BukaTutupCourse(req)
	if err != nil {
		mapError(c, err)
		return
	}
	msg := "Course ditutup (auto-submit massal dijalankan)"
	if req.IsOpen {
		msg = "Course dibuka"
	}
	response.OK(c, http.StatusOK, msg, res)
}

// AddSusulan POST /api/admin/aktivasi/:id/susulan
func (h *AktivasiHandler) AddSusulan(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	var req dto.SusulanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	req.AktivasiSesiID = id
	if err := h.uc.AddSusulan(req); err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Mahasiswa susulan didaftarkan", nil)
}

// ListSusulan GET /api/admin/aktivasi/:id/susulan
func (h *AktivasiHandler) ListSusulan(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	res, err := h.uc.ListSusulan(id)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar susulan", res)
}

// RemoveSusulan DELETE /api/admin/aktivasi/:id/susulan/:mahasiswaId
func (h *AktivasiHandler) RemoveSusulan(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	mid, ok := idParam(c, "mahasiswaId")
	if !ok {
		return
	}
	if err := h.uc.RemoveSusulan(id, mid); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Susulan dihapus", nil)
}
