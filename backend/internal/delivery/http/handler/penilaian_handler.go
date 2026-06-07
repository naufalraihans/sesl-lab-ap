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

// Rekap GET /api/admin/penilaian/rekap
// @Summary Rekap Nilai
// @Description Mengambil rekapitulasi nilai jawaban untuk satu course pada suatu sesi
// @Tags Admin - Penilaian
// @Security bearerAuth
// @Produce json
// @Param aktivasi_sesi_id query int true "ID Aktivasi Sesi"
// @Param course_id query int true "ID Course"
// @Success 200 {object} response.Envelope{data=dto.RekapResponse}
// @Router /admin/penilaian/rekap [get]
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
// @Summary Set Nilai Jawaban
// @Description Memberikan atau mengubah nilai (skor) pada jawaban mahasiswa secara manual
// @Tags Admin - Penilaian
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param request body dto.NilaiRequest true "Payload Nilai"
// @Success 200 {object} response.Envelope
// @Router /admin/penilaian [post]
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
