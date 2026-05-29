package http

import (
	"net/http"

	"backend-gmao/apps/prediction-service/internal/core/domain"
	"backend-gmao/apps/prediction-service/internal/core/ports/primary"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PredictionHandler manages HTTP routes for predictions.
type PredictionHandler struct {
	predictionService primary.PredictionService
}

// NewPredictionHandler creates a new PredictionHandler.
func NewPredictionHandler(predictionService primary.PredictionService) *PredictionHandler {
	return &PredictionHandler{predictionService: predictionService}
}

// CreatePrediction handles prediction creation.
func (h *PredictionHandler) CreatePrediction(c *gin.Context) {
	var req domain.CreatePredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.predictionService.RecordPrediction(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetPrediction retrieves a single prediction.
func (h *PredictionHandler) GetPrediction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	resp, err := h.predictionService.GetPrediction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prediction not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListPredictionsForAsset lists predictions matching the asset identifier.
func (h *PredictionHandler) ListPredictionsForAsset(c *gin.Context) {
	assetIDStr := c.Param("assetId")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset UUID format"})
		return
	}

	resp, err := h.predictionService.GetPredictionsForAsset(c.Request.Context(), assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListPredictions lists all predictions in history.
func (h *PredictionHandler) ListPredictions(c *gin.Context) {
	resp, err := h.predictionService.GetAllPredictions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
