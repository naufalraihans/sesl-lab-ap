package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type PenilaianHandler struct {
	uc *usecase.PenilaianUsecase
}

func NewPenilaianHandler(uc *usecase.PenilaianUsecase) *PenilaianHandler {
	return &PenilaianHandler{uc: uc}
}

// Rekap GET /api/admin/penilaian/rekap?aktivasi_sesi_id=&course_id=
func (h *PenilaianHandler) Rekap(c *gin.Context) {
	aks := queryIntPtr(c, "aktivasi_sesi_id")
	cid := queryIntPtr(c, "course_id")
	if aks == nil || cid == nil {
		response.Fail(c, http.StatusBadRequest, "aktivasi_sesi_id & course_id wajib", nil)
		return
	}
	res, err := h.uc.Rekap(*aks, *cid)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Rekap jawaban", res)
}

// SetNilai POST /api/admin/penilaian
func (h *PenilaianHandler) SetNilai(c *gin.Context) {
	var req dto.NilaiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.SetNilai(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Nilai disimpan", res)
}
