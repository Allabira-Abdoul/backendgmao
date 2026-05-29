package http

import (
	"backend-gmao/apps/auth-service/internal/application/service"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the auth service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	authService *service.AuthService,
) {
	authHandler := NewAuthHandler(authService)

	// Public routes
	router.POST("/sessions", authHandler.CreateSession)
	router.POST("/sessions/validate", authHandler.ValidateSession)
	router.POST("/sessions/refresh", authHandler.RefreshSession)

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		authenticated.DELETE("/sessions", authHandler.RevokeSession)
	}
}
