package http

import (
	"backend-gmao/apps/audit-service/internal/application/service"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the audit service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	auditService *service.AuditService,
) {
	auditHandler := NewAuditHandler(auditService)

	// 🔒 Internal Route: Removed. Log ingestion is now done via EventBus.
	// internal := router.Group("/internal")
	// internal.Use(middleware.RequireInternalService())
	// {
	// 	internal.POST("/audit-logs", auditHandler.WriteLog)
	// }

	// 🔒 Public/Gateway Auditable Route: Read logs (accessible only by those with AUDITOR privilege)
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	authenticated.Use(middleware.RequirePrivilege("AUDIT_LOG_VIEW"))
	{
		authenticated.GET("/audit-logs", auditHandler.ListLogs)
	}
}
