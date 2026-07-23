package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, assetHandler *AssetHandler) {
	v1 := router.Group("/api/v1")
	{
		// Hierarchy Endpoints
		v1.POST("/sites", assetHandler.CreateSite)
		v1.GET("/sites", assetHandler.GetAllSites)
		v1.GET("/sites/:id/hierarchy", assetHandler.GetSiteHierarchy)
		
		v1.POST("/systems", assetHandler.CreateSystem)
		v1.POST("/assets", assetHandler.CreateAsset)
		v1.POST("/subsystems", assetHandler.CreateSubsystem)

		// Inventory & Components
		v1.POST("/inventory", assetHandler.CreateInventoryItem)
		v1.POST("/components", assetHandler.CreateComponent)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
}
