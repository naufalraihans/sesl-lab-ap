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
func (h *AIGradingHandler) BulkGrade(c *gin.Context) {
	var req dto.AIGradingBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.aiUsecase.QueueJob(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Sesuai standar asynchronous API, kembalikan 202 Accepted
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Permintaan AI Grading sedang diproses di latar belakang.",
		"data":    resp,
	})
}

// GetJobStatus GET /api/admin/jobs/:id
func (h *AIGradingHandler) GetJobStatus(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job ID dibutuhkan"})
		return
	}

	resp, err := h.aiUsecase.GetJobStatus(jobID)
	if err != nil {
		if err == usecase.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}
