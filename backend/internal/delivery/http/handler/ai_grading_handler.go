package handler

import (
	"net/http"

	"lab-ap/internal/dto"
	"lab-ap/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AIGradingHandler struct {
	aiUsecase usecase.AIGradingUsecase
}

func NewAIGradingHandler(ai usecase.AIGradingUsecase) *AIGradingHandler {
	return &AIGradingHandler{aiUsecase: ai}
}

// BulkGrade POST /api/admin/penilaian/ai-grade/bulk
// @Summary AI Grading Massal
// @Description Memulai proses penilaian otomatis menggunakan AI di latar belakang (Job Queue)
// @Tags Admin - Penilaian
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param request body dto.AIGradingBulkRequest true "Daftar Mahasiswa ID dan Course ID"
// @Success 202 {object} response.Envelope{data=dto.AIGradingJobResponse}
// @Router /admin/penilaian/ai-grade/bulk [post]
func (h *AIGradingHandler) BulkGrade(c *gin.Context) {
	var req dto.AIGradingBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error(), "message": err.Error()})
		return
	}

	resp, err := h.aiUsecase.QueueJob(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error(), "message": err.Error()})
		return
	}

	// Sesuai standar asynchronous API, kembalikan 202 Accepted
	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Permintaan AI Grading sedang diproses di latar belakang.",
		"data":    resp,
	})
}

// GetJobStatus GET /api/admin/jobs/:id
// @Summary Status AI Grading Job
// @Description Mengambil status pengerjaan job AI Grading yang berjalan di latar belakang
// @Tags Admin - Penilaian
// @Security bearerAuth
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} response.Envelope{data=dto.AIGradingJobResponse}
// @Router /admin/jobs/{id} [get]
func (h *AIGradingHandler) GetJobStatus(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "job ID dibutuhkan", "message": "job ID dibutuhkan"})
		return
	}

	resp, err := h.aiUsecase.GetJobStatus(jobID)
	if err != nil {
		if err == usecase.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Job tidak ditemukan", "message": "Job tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error(), "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp, "message": "OK"})
}
