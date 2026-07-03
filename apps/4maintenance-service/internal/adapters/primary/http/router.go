package http

import (
	"backend-gmao/apps/maintenance-service/internal/application/service"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the maintenance service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	maintenanceService *service.MaintenanceService,
) {
	maintenanceHandler := NewMaintenanceHandler(maintenanceService)

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		workorders := authenticated.Group("/work-orders")
		{
			workorders.POST("", maintenanceHandler.CreateWorkOrder)
			workorders.GET("", maintenanceHandler.ListWorkOrders)
			workorders.GET("/:id", maintenanceHandler.GetWorkOrder)
			workorders.PUT("/:id", maintenanceHandler.UpdateWorkOrder)
			workorders.DELETE("/:id", maintenanceHandler.DeleteWorkOrder)
			workorders.POST("/:id/start", maintenanceHandler.StartWorkOrder)

			// Interventions under work order
			workorders.POST("/:id/interventions", maintenanceHandler.CreateIntervention)
			workorders.GET("/:id/interventions", maintenanceHandler.GetInterventions)
			workorders.POST("/:id/interventions/:inv_id/start", maintenanceHandler.StartIntervention)
			workorders.PUT("/:id/interventions/:inv_id", maintenanceHandler.UpdateIntervention)
			workorders.POST("/:id/interventions/:inv_id/end", maintenanceHandler.EndIntervention)

		}

		inspections := authenticated.Group("/inspections")
		{
			inspections.POST("", maintenanceHandler.CreateInspection)
			inspections.POST("/:id/start", maintenanceHandler.StartInspection)
			inspections.PUT("/:id", maintenanceHandler.UpdateInspection)
			inspections.POST("/:id/end", maintenanceHandler.EndInspection)
			inspections.GET("/asset/:asset_id", maintenanceHandler.HandleGetAssetInspections)
		}

		schedules := authenticated.Group("/schedules")
		{
			schedules.GET("", maintenanceHandler.HandleGetAllSchedules)
			schedules.POST("", maintenanceHandler.HandleCreateSchedule)
			schedules.PUT("/:id", maintenanceHandler.HandleUpdateSchedule)
			schedules.GET("/asset/:asset_id", maintenanceHandler.HandleGetAssetSchedules)
		}

		readings := authenticated.Group("/readings")
		{
			readings.POST("", maintenanceHandler.HandleRecordCounterReading)
			readings.GET("/asset/:asset_id", maintenanceHandler.HandleGetAssetReadings)
		}

		// Defect Alerts endpoints
		alerts := authenticated.Group("/alerts")
		{
			alerts.POST("", maintenanceHandler.HandleCreateDefectAlert)
			alerts.GET("", maintenanceHandler.HandleGetAllDefectAlerts)
			alerts.PUT("/:id/review", maintenanceHandler.HandleReviewDefectAlert)
		}
	}
}
