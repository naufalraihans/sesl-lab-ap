package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type SesiHandler struct {
	uc *usecase.SesiUsecase
}

func NewSesiHandler(uc *usecase.SesiUsecase) *SesiHandler { return &SesiHandler{uc: uc} }

func (h *SesiHandler) List(c *gin.Context) {
	res, err := h.uc.List()
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar sesi", res)
}

func (h *SesiHandler) Get(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	res, err := h.uc.Get(id)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Detail sesi", res)
}

func (h *SesiHandler) Create(c *gin.Context) {
	var req dto.SesiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.Create(req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Sesi dibuat", res)
}

func (h *SesiHandler) Update(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	var req dto.SesiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.Update(id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Sesi diperbarui", res)
}

func (h *SesiHandler) Delete(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Delete(id); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Sesi dihapus", nil)
}

// ---- Course ----

func (h *SesiHandler) ListCourse(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	res, err := h.uc.ListCourse(id)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar course", res)
}

func (h *SesiHandler) CreateCourse(c *gin.Context) {
	id, ok := idParam(c, "id")
	if !ok {
		return
	}
	var req dto.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.CreateCourse(id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.Created(c, "Course dibuat", res)
}

func (h *SesiHandler) UpdateCourse(c *gin.Context) {
	cid, ok := idParam(c, "courseId")
	if !ok {
		return
	}
	var req dto.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}
	res, err := h.uc.UpdateCourse(cid, req)
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Course diperbarui", res)
}

func (h *SesiHandler) DeleteCourse(c *gin.Context) {
	cid, ok := idParam(c, "courseId")
	if !ok {
		return
	}
	if err := h.uc.DeleteCourse(cid); err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Course dihapus", nil)
}
