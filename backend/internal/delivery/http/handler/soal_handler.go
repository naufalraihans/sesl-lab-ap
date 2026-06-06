package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type SoalHandler struct {
	uc *usecase.SoalUsecase
}

func NewSoalHandler(uc *usecase.SoalUsecase) *SoalHandler { return &SoalHandler{uc: uc} }

// ListByCourse GET /api/admin/soal?course_id=
func (h *SoalHandler) ListByCourse(c *gin.Context) {
	cid := queryIntPtr(c, "course_id")
	if cid == nil {
		response.Fail(c, http.StatusBadRequest, "course_id wajib", nil)
		return
	}
	res, err := h.uc.ListByCourse(*cid)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar soal", res)
}

func (h *SoalHandler) Create(c *gin.Context) {
	var req dto.SoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.Create(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Soal dibuat", res)
}

func (h *SoalHandler) Update(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	var req dto.SoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.Update(id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Soal diperbarui", res)
}

func (h *SoalHandler) Delete(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Delete(id); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Soal dihapus", nil)
}
