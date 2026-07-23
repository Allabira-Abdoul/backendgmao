package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/apps/asset-service/internal/core/ports/primary"
)

type AssetHandler struct {
	service primary.AssetService
}

func NewAssetHandler(service primary.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

func (h *AssetHandler) CreateSite(c *gin.Context) {
	var req domain.CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateSite(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AssetHandler) GetAllSites(c *gin.Context) {
	resp, err := h.service.GetAllSites(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AssetHandler) GetSiteHierarchy(c *gin.Context) {
	siteIDStr := c.Param("id")
	siteID, err := uuid.Parse(siteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid site ID"})
		return
	}

	resp, err := h.service.GetSiteHierarchy(c.Request.Context(), siteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "site not found or hierarchy error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AssetHandler) CreateSystem(c *gin.Context) {
	var req domain.CreateSystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateSystem(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var req domain.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateAsset(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AssetHandler) CreateSubsystem(c *gin.Context) {
	var req domain.CreateSubsystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateSubsystem(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AssetHandler) CreateInventoryItem(c *gin.Context) {
	var req domain.CreateInventoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateInventoryItem(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AssetHandler) CreateComponent(c *gin.Context) {
	var req domain.CreateComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateComponent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}
