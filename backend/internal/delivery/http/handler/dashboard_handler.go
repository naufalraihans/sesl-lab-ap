package handler

import (
	"net/http"

	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	uc *usecase.DashboardUsecase
}

func NewDashboardHandler(uc *usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{uc: uc}
}

// Statistik GET /api/admin/dashboard
func (h *DashboardHandler) Statistik(c *gin.Context) {
	res, err := h.uc.Statistik()
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Statistik dashboard", res)
}

// Online GET /api/admin/dashboard/online
func (h *DashboardHandler) Online(c *gin.Context) {
	response.OK(c, http.StatusOK, "Jumlah online", h.uc.OnlineCount())
}
