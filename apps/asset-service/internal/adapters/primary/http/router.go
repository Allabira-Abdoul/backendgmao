package http

import (
	"backend-gmao/apps/asset-service/internal/core/ports/primary"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the asset service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	assetService primary.AssetService,
) {
	assetHandler := NewAssetHandler(assetService)

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		// Legacy endpoint for backward compatibility
		assets := authenticated.Group("/assets")
		{
			assets.GET("", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetLegacyAssets)
		}

		models := authenticated.Group("/models")
		{
			models.POST("/equipment", middleware.RequirePrivilege("ASSET_CREATE"), assetHandler.CreateEquipmentModel)
			models.POST("/parts", middleware.RequirePrivilege("ASSET_CREATE"), assetHandler.CreatePartModel)
			models.GET("/equipment", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetEquipmentModels)
			models.GET("/parts", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetPartModels)

			models.PUT("/equipment/:id", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.UpdateEquipmentModel)
			models.PUT("/parts/:id", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.UpdatePartModel)
			models.POST("/equipment/:id/suppliers", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.AddSupplierToEquipmentModel)
			models.POST("/parts/:id/suppliers", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.AddSupplierToPartModel)
		}

		suppliers := authenticated.Group("/suppliers")
		{
			suppliers.POST("", middleware.RequirePrivilege("ASSET_CREATE"), assetHandler.CreateSupplier)
			suppliers.GET("", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetSuppliers)
			suppliers.PUT("/:id", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.UpdateSupplier)
			suppliers.DELETE("/:id", middleware.RequirePrivilege("ASSET_DELETE"), assetHandler.DeleteSupplier)
		}

		instances := authenticated.Group("/instances")
		{
			instances.POST("/equipment", middleware.RequirePrivilege("ASSET_CREATE"), assetHandler.CreateEquipmentInstance)
			instances.POST("/parts", middleware.RequirePrivilege("ASSET_CREATE"), assetHandler.CreatePartInstance)
			instances.GET("/equipment", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetEquipmentInstances)
			instances.GET("/equipment/code/:code", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetEquipmentInstanceByCode)
			instances.GET("/equipment/:id", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetEquipmentInstanceByID)

			instances.POST("/parts/:id/move", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.MovePartInstance)
		}

		actions := authenticated.Group("/actions")
		{
			actions.POST("/consume-part", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.ConsumePart)
		}

		measurements := authenticated.Group("/measurements")
		{
			measurements.POST("", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.IngestMeasurement)
			measurements.GET("/:targetType/:targetID", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetMeasurements)
		}

		thresholds := authenticated.Group("/thresholds")
		{
			thresholds.POST("", middleware.RequirePrivilege("ASSET_CREATE"), assetHandler.CreateMetricThreshold)
			thresholds.PUT("/:id", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.UpdateMetricThreshold)
			thresholds.DELETE("/:id", middleware.RequirePrivilege("ASSET_DELETE"), assetHandler.DeleteMetricThreshold)
		}
	}
}
