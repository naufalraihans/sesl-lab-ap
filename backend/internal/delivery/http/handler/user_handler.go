package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler { return &UserHandler{uc: uc} }

// ---- Mahasiswa ----

// ListMahasiswa GET /api/admin/users?kelas_id=&shift=
func (h *UserHandler) ListMahasiswa(c *gin.Context) {
	res, err := h.uc.ListMahasiswa(queryIntPtr(c, "kelas_id"), queryIntPtr(c, "shift"))
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar mahasiswa", res)
}

func (h *UserHandler) CreateMahasiswa(c *gin.Context) {
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.CreateMahasiswa(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Mahasiswa dibuat", res)
}

func (h *UserHandler) UpdateMahasiswa(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.UpdateMahasiswa(id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Mahasiswa diperbarui", res)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Delete(id); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "User dihapus", nil)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	if err := h.uc.ResetPassword(id); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Password direset (mahasiswa harus register ulang)", nil)
}

// SetRegisterOpen POST /api/admin/kelas/register
func (h *UserHandler) SetRegisterOpen(c *gin.Context) {
	var req dto.RegisterOpenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	if err := h.uc.SetRegisterOpen(req.KelasID, req.Open); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Status register kelas diperbarui", nil)
}

// ---- Asisten ----

// ListAsisten GET /api/admin/asisten & /api/info/asisten (publik)
func (h *UserHandler) ListAsisten(c *gin.Context) {
	res, err := h.uc.ListAsisten()
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar asisten", res)
}

func (h *UserHandler) CreateAsisten(c *gin.Context) {
	var req dto.AsistenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.CreateAsisten(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Asisten dibuat", res)
}

func (h *UserHandler) UpdateAsisten(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	var req dto.AsistenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.UpdateAsisten(id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Asisten diperbarui", res)
}
