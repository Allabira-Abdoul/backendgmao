package http

import (
	"backend-gmao/apps/analytics-service/internal/application/service"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the analytics service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	analyticsService *service.AnalyticsService,
) {
	metricHandler := NewMetricHandler(analyticsService)

	kpiHandler := NewKpiHandler(analyticsService)

	// Internal authenticated routes
	internal := router.Group("/internal")
	internal.Use(middleware.RequireInternalService())
	{
		// No longer exposing internal HTTP endpoints for maintenance events (using RabbitMQ instead)
	}

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		metrics := authenticated.Group("/metrics")
		{
			metrics.POST("", middleware.RequirePrivilege("ANALYTICS_WRITE"), metricHandler.RecordMetric)
			metrics.GET("", middleware.RequirePrivilege("ANALYTICS_VIEW"), metricHandler.ListMetrics)
			metrics.GET("/:id", middleware.RequirePrivilege("ANALYTICS_VIEW"), metricHandler.GetMetric)
			metrics.GET("/category/:category", middleware.RequirePrivilege("ANALYTICS_VIEW"), metricHandler.ListMetricsByCategory)
		}
		kpis := authenticated.Group("/kpis")
		{
			kpis.GET("/categories/health", middleware.RequirePrivilege("ANALYTICS_VIEW"), kpiHandler.GetCategoryHealthMetrics)
		}
	}
}
