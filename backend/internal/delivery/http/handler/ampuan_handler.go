package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type AmpuanHandler struct {
	uc   *usecase.AmpuanUsecase
	user *usecase.UserUsecase
}

func NewAmpuanHandler(uc *usecase.AmpuanUsecase, user *usecase.UserUsecase) *AmpuanHandler {
	return &AmpuanHandler{uc: uc, user: user}
}

func (h *AmpuanHandler) List(c *gin.Context) {
	res, err := h.uc.List()
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar ampuan", res)
}

func (h *AmpuanHandler) Create(c *gin.Context) {
	var req dto.AmpuanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.Create(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Ampuan ditambahkan", res)
}

func (h *AmpuanHandler) Delete(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Delete(id); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Ampuan dihapus", nil)
}

func (h *AmpuanHandler) PublicKelasMahasiswa(c *gin.Context) {
	kelasID, ok := idParam(c, "id")
	if !ok {
		return
	}
	mhs, err := h.user.ListMahasiswa(&kelasID, nil)
	if err != nil {
		mapError(c, err)
		return
	}
	ampuan, err := h.uc.ListByKelas(kelasID)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Mahasiswa & ampuan kelas", gin.H{
		"mahasiswa": mhs,
		"ampuan":    ampuan,
	})
}
