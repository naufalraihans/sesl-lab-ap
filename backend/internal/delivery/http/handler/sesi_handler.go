package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	_ "lab-ap/internal/entity"
	"lab-ap/internal/usecase"
	"lab-ap/pkg/response"

	"github.com/gin-gonic/gin"
)

type SesiHandler struct {
	uc *usecase.SesiUsecase
}

func NewSesiHandler(uc *usecase.SesiUsecase) *SesiHandler { return &SesiHandler{uc: uc} }

// List GET /api/admin/sesi
// @Summary Daftar Master Sesi
// @Description Mengambil seluruh master data sesi praktikum
// @Tags Admin - Sesi
// @Security bearerAuth
// @Produce json
// @Success 200 {object} response.Envelope{data=[]entity.SesiPraktikum}
// @Router /admin/sesi [get]
func (h *SesiHandler) List(c *gin.Context) {
	res, err := h.uc.List()
	if err != nil {
		mapError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "Daftar sesi", res)
}

// Get GET /api/admin/sesi/:id
// @Summary Detail Master Sesi
// @Description Mengambil detail master sesi praktikum
// @Tags Admin - Sesi
// @Security bearerAuth
// @Produce json
// @Param id path int true "ID Sesi"
// @Success 200 {object} response.Envelope{data=entity.SesiPraktikum}
// @Router /admin/sesi/{id} [get]
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

// Create POST /api/admin/sesi
// @Summary Tambah Master Sesi
// @Description Menambahkan master sesi praktikum baru
// @Tags Admin - Sesi
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SesiRequest true "Payload"
// @Success 201 {object} response.Envelope{data=entity.SesiPraktikum}
// @Router /admin/sesi [post]
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

// Update PUT /api/admin/sesi/:id
// @Summary Perbarui Master Sesi
// @Description Mengubah data master sesi praktikum
// @Tags Admin - Sesi
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Sesi"
// @Param request body dto.SesiRequest true "Payload"
// @Success 200 {object} response.Envelope{data=entity.SesiPraktikum}
// @Router /admin/sesi/{id} [put]
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

// Delete DELETE /api/admin/sesi/:id
// @Summary Hapus Master Sesi
// @Description Menghapus master sesi praktikum
// @Tags Admin - Sesi
// @Security bearerAuth
// @Produce json
// @Param id path int true "ID Sesi"
// @Success 200 {object} response.Envelope
// @Router /admin/sesi/{id} [delete]
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

// ListCourse GET /api/admin/sesi/:id/course
// @Summary Daftar Course
// @Description Mengambil course (pretest/posttest/dll) di dalam satu sesi
// @Tags Admin - Course
// @Security bearerAuth
// @Produce json
// @Param id path int true "ID Sesi"
// @Success 200 {object} response.Envelope{data=[]entity.Course}
// @Router /admin/sesi/{id}/course [get]
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

// CreateCourse POST /api/admin/sesi/:id/course
// @Summary Tambah Course
// @Description Menambahkan course ke dalam master sesi
// @Tags Admin - Course
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Sesi"
// @Param request body dto.CourseRequest true "Payload"
// @Success 201 {object} response.Envelope{data=entity.Course}
// @Router /admin/sesi/{id}/course [post]
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

// UpdateCourse PUT /api/admin/course/:courseId
// @Summary Perbarui Course
// @Description Mengubah course (termasuk pool soal dan durasi)
// @Tags Admin - Course
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param courseId path int true "ID Course"
// @Param request body dto.CourseRequest true "Payload"
// @Success 200 {object} response.Envelope{data=entity.Course}
// @Router /admin/course/{courseId} [put]
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

// DeleteCourse DELETE /api/admin/course/:courseId
// @Summary Hapus Course
// @Description Menghapus course dari sesi
// @Tags Admin - Course
// @Security bearerAuth
// @Produce json
// @Param courseId path int true "ID Course"
// @Success 200 {object} response.Envelope
// @Router /admin/course/{courseId} [delete]
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
