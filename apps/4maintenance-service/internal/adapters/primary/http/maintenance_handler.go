package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/apps/maintenance-service/internal/core/ports/primary"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MaintenanceHandler manages HTTP routes for maintenance.
type MaintenanceHandler struct {
	maintenanceService primary.MaintenanceService
}

// NewMaintenanceHandler creates a new MaintenanceHandler.
func NewMaintenanceHandler(maintenanceService primary.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{maintenanceService: maintenanceService}
}

// CreateWorkOrder handles work order creation.
func (h *MaintenanceHandler) CreateWorkOrder(c *gin.Context) {
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("[DEBUG] Raw Request Body: %s\n", string(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	
	var req domain.CreateOrdreTravailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[DEBUG] CreateWorkOrder binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.CreateWorkOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateWorkOrder handles work order updating.
func (h *MaintenanceHandler) UpdateWorkOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var req domain.UpdateOrdreTravailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.UpdateWorkOrder(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetWorkOrder gets a specific work order with its interventions.
func (h *MaintenanceHandler) GetWorkOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	resp, err := h.maintenanceService.GetWorkOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work order not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListWorkOrders lists all work orders.
func (h *MaintenanceHandler) ListWorkOrders(c *gin.Context) {
	resp, err := h.maintenanceService.GetAllWorkOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteWorkOrder deletes a specific work order.
func (h *MaintenanceHandler) DeleteWorkOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	err = h.maintenanceService.DeleteWorkOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work order deleted successfully"})
}

// StartWorkOrder starts a work order.
func (h *MaintenanceHandler) StartWorkOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	resp, err := h.maintenanceService.StartWorkOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateIntervention records a new intervention for a work order.
func (h *MaintenanceHandler) CreateIntervention(c *gin.Context) {
	idStr := c.Param("id")
	workOrderID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var req domain.CreateInterventionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.RecordIntervention(c.Request.Context(), workOrderID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetInterventions lists interventions recorded under a specific work order.
func (h *MaintenanceHandler) GetInterventions(c *gin.Context) {
	idStr := c.Param("id")
	workOrderID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	resp, err := h.maintenanceService.GetInterventionsForWorkOrder(c.Request.Context(), workOrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// StartIntervention starts an intervention.
func (h *MaintenanceHandler) StartIntervention(c *gin.Context) {
	woIDStr := c.Param("id")
	woID, err := uuid.Parse(woIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work order UUID"})
		return
	}

	invIDStr := c.Param("inv_id")
	invID, err := uuid.Parse(invIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid intervention UUID"})
		return
	}

	resp, err := h.maintenanceService.StartIntervention(c.Request.Context(), woID, invID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateIntervention updates an intervention.
func (h *MaintenanceHandler) UpdateIntervention(c *gin.Context) {
	woIDStr := c.Param("id")
	woID, err := uuid.Parse(woIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work order UUID"})
		return
	}

	invIDStr := c.Param("inv_id")
	invID, err := uuid.Parse(invIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid intervention UUID"})
		return
	}

	var req domain.UpdateInterventionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.UpdateIntervention(c.Request.Context(), woID, invID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// EndIntervention ends an intervention.
func (h *MaintenanceHandler) EndIntervention(c *gin.Context) {
	woIDStr := c.Param("id")
	woID, err := uuid.Parse(woIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work order UUID"})
		return
	}

	invIDStr := c.Param("inv_id")
	invID, err := uuid.Parse(invIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid intervention UUID"})
		return
	}

	resp, err := h.maintenanceService.EndIntervention(c.Request.Context(), woID, invID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CreateInspection records a new inspection.
func (h *MaintenanceHandler) CreateInspection(c *gin.Context) {

	var req domain.CreateInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.CreateInspection(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateInspection updates an inspection.
func (h *MaintenanceHandler) UpdateInspection(c *gin.Context) {
	insIDStr := c.Param("id")
	insID, err := uuid.Parse(insIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inspection UUID"})
		return
	}

	var req domain.UpdateInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.UpdateInspection(c.Request.Context(), insID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// StartInspection starts an inspection.
func (h *MaintenanceHandler) StartInspection(c *gin.Context) {
	insIDStr := c.Param("id")
	insID, err := uuid.Parse(insIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inspection UUID"})
		return
	}

	resp, err := h.maintenanceService.StartInspection(c.Request.Context(), insID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// EndInspection ends an inspection.
func (h *MaintenanceHandler) EndInspection(c *gin.Context) {
	insIDStr := c.Param("id")
	insID, err := uuid.Parse(insIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inspection UUID"})
		return
	}

	resp, err := h.maintenanceService.EndInspection(c.Request.Context(), insID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// HandleGetAssetInspections handles fetching inspections for an asset.
func (h *MaintenanceHandler) HandleGetAssetInspections(c *gin.Context) {
	assetIDStr := c.Param("asset_id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset UUID"})
		return
	}

	resp, err := h.maintenanceService.GetInspectionsForAsset(c.Request.Context(), assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// HandleCreateSchedule handles creating a new maintenance schedule.
func (h *MaintenanceHandler) HandleCreateSchedule(c *gin.Context) {
	var req domain.CreateMaintenanceScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.CreateMaintenanceSchedule(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// HandleUpdateSchedule handles updating a maintenance schedule.
func (h *MaintenanceHandler) HandleUpdateSchedule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule UUID"})
		return
	}

	var req domain.UpdateMaintenanceScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.UpdateMaintenanceSchedule(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetAssetSchedules lists all schedules for an asset.
func (h *MaintenanceHandler) HandleGetAssetSchedules(c *gin.Context) {
	assetIDStr := c.Param("asset_id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset UUID"})
		return
	}

	resp, err := h.maintenanceService.GetMaintenanceSchedulesForAsset(c.Request.Context(), assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetAllSchedules returns all maintenance schedules.
func (h *MaintenanceHandler) HandleGetAllSchedules(c *gin.Context) {
	resp, err := h.maintenanceService.GetAllMaintenanceSchedules(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ----------------------------------------------------------------------------
// Counter Readings
func (h *MaintenanceHandler) HandleRecordCounterReading(c *gin.Context) {
	var req domain.CreateCounterReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.RecordCounterReading(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// HandleGetAssetReadings lists historical counter readings for an asset.
func (h *MaintenanceHandler) HandleGetAssetReadings(c *gin.Context) {
	assetIDStr := c.Param("asset_id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset UUID"})
		return
	}

	resp, err := h.maintenanceService.GetCounterReadingsForAsset(c.Request.Context(), assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ----------------------------------------------------------------------------
// Defect Alerts
// ----------------------------------------------------------------------------

func (h *MaintenanceHandler) HandleCreateDefectAlert(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large or invalid multipart form"})
		return
	}

	assetIDStr := c.PostForm("asset_id")
	title := c.PostForm("title")
	description := c.PostForm("description")

	// Get logged-in user from context (set by auth middleware)
	userIDStr, exists := c.Get("auth_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}
	reportedBy, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in token"})
		return
	}

	var imageURL string
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		// Save file
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		savePath := filepath.Join("uploads", filename)

		out, err := os.Create(savePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write image"})
			return
		}
		
		// For MVP, serve through the API gateway so CORS headers are applied
		imageURL = fmt.Sprintf("http://127.0.0.1:8200/api/maintenance/uploads/%s", filename) // Hardcoded for local testing
	}

	resp, err := h.maintenanceService.CreateDefectAlert(c.Request.Context(), assetID, reportedBy, title, description, imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *MaintenanceHandler) HandleGetAllDefectAlerts(c *gin.Context) {
	resp, err := h.maintenanceService.GetAllDefectAlerts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MaintenanceHandler) HandleReviewDefectAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	var req domain.ReviewDefectAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.ReviewDefectAlert(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
