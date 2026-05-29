package http

import (
	"backend-gmao/apps/prediction-service/internal/application/service"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the prediction service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	predictionService *service.PredictionService,
) {
	predictionHandler := NewPredictionHandler(predictionService)

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		predictions := authenticated.Group("/predictions")
		{
			predictions.POST("", predictionHandler.CreatePrediction)
			predictions.GET("", predictionHandler.ListPredictions)
			predictions.GET("/:id", predictionHandler.GetPrediction)
			predictions.GET("/asset/:assetId", predictionHandler.ListPredictionsForAsset)
		}
	}
}
