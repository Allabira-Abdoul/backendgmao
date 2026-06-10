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

			// Inspections under work order
			workorders.POST("/:id/inspections", maintenanceHandler.CreateInspection)
			workorders.POST("/:id/inspections/:ins_id/start", maintenanceHandler.StartInspection)
			workorders.PUT("/:id/inspections/:ins_id", maintenanceHandler.UpdateInspection)
			workorders.POST("/:id/inspections/:ins_id/end", maintenanceHandler.EndInspection)
		}
	}
}
