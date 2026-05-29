package http

import (
	"net/http"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"backend-gmao/apps/analytics-service/internal/core/ports/primary"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KpiHandler struct {
	analyticsService primary.AnalyticsService
}

func NewKpiHandler(analyticsService primary.AnalyticsService) *KpiHandler {
	return &KpiHandler{analyticsService: analyticsService}
}

func (h *KpiHandler) ProcessMaintenanceEvent(c *gin.Context) {
	var event domain.MaintenanceEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	if err := h.analyticsService.ProcessMaintenanceEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Event processed successfully"})
}

func (h *KpiHandler) GetGlobalKpi(c *gin.Context) {
	resp, err := h.analyticsService.GetGlobalKpi(c.Request.Context())
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

func (h *KpiHandler) GetCategoryKpi(c *gin.Context) {
	category := c.Param("category")
	resp, err := h.analyticsService.GetCategoryKpi(c.Request.Context(), category)
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

func (h *KpiHandler) GetAssetKpi(c *gin.Context) {
	idStr := c.Param("asset_id")
	assetID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid UUID format",
			},
		})
		return
	}

	resp, err := h.analyticsService.GetAssetKpi(c.Request.Context(), assetID)
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
