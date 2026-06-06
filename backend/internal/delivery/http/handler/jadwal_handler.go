package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type JadwalHandler struct {
	uc *usecase.JadwalUsecase
}

func NewJadwalHandler(uc *usecase.JadwalUsecase) *JadwalHandler { return &JadwalHandler{uc: uc} }

// List GET /api/info/jadwal (publik) & /api/admin/jadwal
func (h *JadwalHandler) List(c *gin.Context) {
	res, err := h.uc.List()
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar jadwal", res)
}

func (h *JadwalHandler) Create(c *gin.Context) {
	var req dto.JadwalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.Create(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Jadwal dibuat", res)
}

func (h *JadwalHandler) Update(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	var req dto.JadwalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.Update(id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Jadwal diperbarui", res)
}

func (h *JadwalHandler) Delete(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Delete(id); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Jadwal dihapus", nil)
}
