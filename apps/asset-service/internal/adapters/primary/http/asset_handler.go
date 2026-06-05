package http

import (
	"net/http"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/apps/asset-service/internal/core/ports/primary"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	service primary.AssetService
}

func NewAssetHandler(service primary.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

// --- Models ---

func (h *AssetHandler) CreateEquipmentModel(c *gin.Context) {
	var req domain.CreateEquipmentModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateEquipmentModel(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AssetHandler) CreatePartModel(c *gin.Context) {
	var req domain.CreatePartModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreatePartModel(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AssetHandler) GetEquipmentModels(c *gin.Context) {
	res, err := h.service.GetEquipmentModels(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AssetHandler) GetPartModels(c *gin.Context) {
	res, err := h.service.GetPartModels(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// --- Instances ---

func (h *AssetHandler) CreateEquipmentInstance(c *gin.Context) {
	var req domain.CreateEquipmentInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateEquipmentInstance(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AssetHandler) CreatePartInstance(c *gin.Context) {
	var req domain.CreatePartInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreatePartInstance(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AssetHandler) GetEquipmentInstances(c *gin.Context) {
	res, err := h.service.GetEquipmentInstances(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AssetHandler) GetEquipmentInstanceByCode(c *gin.Context) {
	code := c.Param("code")
	res, err := h.service.GetEquipmentInstanceByCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AssetHandler) GetEquipmentInstanceByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	res, err := h.service.GetEquipmentInstanceByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// --- Actions ---

func (h *AssetHandler) ConsumePart(c *gin.Context) {
	var req domain.ConsumePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract userID from context (set by Auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.service.ConsumePart(c.Request.Context(), req, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "part consumed successfully"})
}

func (h *AssetHandler) MovePartInstance(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	var req domain.MovePartInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.MovePartInstance(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// --- Measurements ---

func (h *AssetHandler) IngestMeasurement(c *gin.Context) {
	var req domain.IngestMeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		parsed, err := uuid.Parse(userIDStr.(string))
		if err == nil {
			userID = &parsed
		}
	}

	res, err := h.service.IngestMeasurement(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AssetHandler) GetMeasurements(c *gin.Context) {
	targetType := c.Param("targetType") // "equipment" or "part"
	targetIDParam := c.Param("targetID")

	targetID, err := uuid.Parse(targetIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target ID"})
		return
	}

	since := c.Query("since")

	res, err := h.service.GetMeasurements(c.Request.Context(), targetType, targetID, since)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// --- Backward compatibility wrapper for old clients ---
func (h *AssetHandler) GetLegacyAssets(c *gin.Context) {
	res, err := h.service.GetEquipmentInstances(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Map new response to old response structure
	legacy := make([]domain.AssetResponse, len(res))
	for i, inst := range res {
		legacy[i] = domain.AssetResponse{
			ID:            inst.ID,
			Type:          "EQUIPMENT",
			Name:          inst.EquipmentModel.Name,
			Code:          inst.Code,
			Status:        inst.Status,
			Category:      inst.EquipmentModel.Category,
			Location:      inst.Location,
			StockQuantity: 1, // Equipment instances are inherently 1
			CreatedAt:     inst.CreatedAt,
			UpdatedAt:     inst.UpdatedAt,
		}
	}
	c.JSON(http.StatusOK, legacy)
}

// --- Suppliers ---

func (h *AssetHandler) CreateSupplier(c *gin.Context) {
	var req domain.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateSupplier(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AssetHandler) GetSuppliers(c *gin.Context) {
	res, err := h.service.GetSuppliers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AssetHandler) AddSupplierToEquipmentModel(c *gin.Context) {
	idParam := c.Param("id")
	modelID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid equipment model id"})
		return
	}

	var req domain.AddModelSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.AddSupplierToEquipmentModel(c.Request.Context(), modelID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AssetHandler) AddSupplierToPartModel(c *gin.Context) {
	idParam := c.Param("id")
	modelID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid part model id"})
		return
	}

	var req domain.AddModelSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.AddSupplierToPartModel(c.Request.Context(), modelID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
