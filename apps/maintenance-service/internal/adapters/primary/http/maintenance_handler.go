package http

import (
	"net/http"

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
	var req domain.CreateOrdreTravailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

// CreateInspection records a new inspection for a work order.
func (h *MaintenanceHandler) CreateInspection(c *gin.Context) {
	idStr := c.Param("id")
	workOrderID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var req domain.CreateInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.maintenanceService.CreateInspection(c.Request.Context(), workOrderID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// StartInspection starts an inspection.
func (h *MaintenanceHandler) StartInspection(c *gin.Context) {
	woIDStr := c.Param("id")
	woID, err := uuid.Parse(woIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work order UUID"})
		return
	}

	insIDStr := c.Param("ins_id")
	insID, err := uuid.Parse(insIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inspection UUID"})
		return
	}

	resp, err := h.maintenanceService.StartInspection(c.Request.Context(), woID, insID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// EndInspection ends an inspection.
func (h *MaintenanceHandler) EndInspection(c *gin.Context) {
	woIDStr := c.Param("id")
	woID, err := uuid.Parse(woIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work order UUID"})
		return
	}

	insIDStr := c.Param("ins_id")
	insID, err := uuid.Parse(insIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inspection UUID"})
		return
	}

	resp, err := h.maintenanceService.EndInspection(c.Request.Context(), woID, insID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
