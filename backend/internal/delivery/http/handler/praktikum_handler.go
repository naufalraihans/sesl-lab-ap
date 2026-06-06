package handler

import (
	"net/http"

	"lab-ap/internal/delivery/http/middleware"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

// PraktikumHandler melayani sisi mahasiswa (dashboard + daftar sesi).
type PraktikumHandler struct {
	uc *usecase.PraktikumUsecase
}

func NewPraktikumHandler(uc *usecase.PraktikumUsecase) *PraktikumHandler {
	return &PraktikumHandler{uc: uc}
}

// Dashboard GET /api/praktikum/dashboard
func (h *PraktikumHandler) Dashboard(c *gin.Context) {
	res, err := h.uc.Dashboard(middleware.UserID(c))
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Dashboard mahasiswa", res)
}

// ListSesi GET /api/praktikum/sesi
func (h *PraktikumHandler) ListSesi(c *gin.Context) {
	res, err := h.uc.ListSesi(middleware.UserID(c))
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar sesi", res)
}
