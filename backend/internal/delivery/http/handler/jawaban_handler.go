package handler

import (
	"net/http"

	"lab-ap/internal/delivery/http/middleware"
	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type JawabanHandler struct {
	uc *usecase.JawabanUsecase
}

func NewJawabanHandler(uc *usecase.JawabanUsecase) *JawabanHandler { return &JawabanHandler{uc: uc} }

// GetRuang GET /api/praktikum/ruang?aktivasi_sesi_id=&course_id=
func (h *JawabanHandler) GetRuang(c *gin.Context) {
	aks := queryIntPtr(c, "aktivasi_sesi_id")
	cid := queryIntPtr(c, "course_id")
	if aks == nil || cid == nil {
		response.Fail(c, http.StatusBadRequest, "aktivasi_sesi_id & course_id wajib", nil)
		return
	}
	res, err := h.uc.GetRuang(middleware.UserID(c), *aks, *cid)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Ruang course", res)
}

// Mulai POST /api/praktikum/mulai
func (h *JawabanHandler) Mulai(c *gin.Context) {
	var req dto.MulaiCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.Mulai(middleware.UserID(c), req.AktivasiSesiID, req.CourseID)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Mulai mengerjakan", res)
}

// AutoSave POST /api/praktikum/autosave
func (h *JawabanHandler) AutoSave(c *gin.Context) {
	var req dto.AutoSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	if err := h.uc.AutoSave(middleware.UserID(c), req); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Tersimpan", nil)
}

// Submit POST /api/praktikum/submit
func (h *JawabanHandler) Submit(c *gin.Context) {
	var req dto.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	if err := h.uc.Submit(middleware.UserID(c), req); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Jawaban ter-submit", nil)
}
