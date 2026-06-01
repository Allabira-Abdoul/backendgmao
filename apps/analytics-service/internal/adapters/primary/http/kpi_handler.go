package http

import (
	"net/http"

	"backend-gmao/apps/analytics-service/internal/core/ports/primary"
	"github.com/gin-gonic/gin"
)

type KpiHandler struct {
	analyticsService primary.AnalyticsService
}

func NewKpiHandler(analyticsService primary.AnalyticsService) *KpiHandler {
	return &KpiHandler{analyticsService: analyticsService}
}

func (h *KpiHandler) GetCategoryHealthMetrics(c *gin.Context) {
	resp, err := h.analyticsService.GetCategoryHealthMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": resp})
}
