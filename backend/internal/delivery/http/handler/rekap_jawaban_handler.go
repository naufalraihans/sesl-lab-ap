package handler

import (
	"net/http"
	"strconv"

	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RekapJawabanHandler struct {
	penilaianUsecase *usecase.PenilaianUsecase
}

func NewRekapJawabanHandler(p *usecase.PenilaianUsecase) *RekapJawabanHandler {
	return &RekapJawabanHandler{penilaianUsecase: p}
}

// GetRekapJawabanGlobal GET /api/admin/rekap-jawaban
// @Summary Rekap Jawaban Global (Flat List)
// @Description Mengambil data jawaban secara flat untuk keperluan tabel rekap global
// @Tags Admin - Rekap
// @Security bearerAuth
// @Produce json
// @Param kelas_id query int false "Filter berdasarkan ID Kelas"
// @Param sesi_id query int false "Filter berdasarkan ID Sesi"
// @Param search query string false "Filter NIM atau Nama"
// @Param jenis query string false "Filter Jenis Tes (pretest, posttest, dll)"
// @Success 200 {object} response.Envelope{data=dto.RekapJawabanResponse}
// @Router /admin/rekap-jawaban [get]
func (h *RekapJawabanHandler) GetRekapJawabanGlobal(c *gin.Context) {
	kelasID, _ := strconv.Atoi(c.Query("kelas_id"))
	sesiID, _ := strconv.Atoi(c.Query("sesi_id"))
	search := c.Query("search")
	jenis := c.Query("jenis")

	resp, err := h.penilaianUsecase.GetRekapJawabanGlobal(kelasID, sesiID, search, jenis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error(), "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// BulkAction POST /api/admin/penilaian/bulk-action
// @Summary Bulk Action Penilaian (Reset / Hapus)
// @Description Mereset nilai atau menghapus jawaban secara masal
// @Tags Admin - Penilaian
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param request body dto.BulkActionRequest true "Daftar ID dan Aksi"
// @Success 200 {object} response.Envelope
// @Router /admin/penilaian/bulk-action [post]
func (h *RekapJawabanHandler) BulkAction(c *gin.Context) {
	var req dto.BulkActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error(), "message": err.Error()})
		return
	}

	var err error
	if req.Action == "delete" {
		err = h.penilaianUsecase.BulkDeleteJawaban(req.JawabanIDs)
	} else if req.Action == "reset_nilai" {
		err = h.penilaianUsecase.BulkResetNilai(req.JawabanIDs)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error(), "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Bulk action " + req.Action + " berhasil"})
}
