package http

import (
	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/apps/user-service/internal/core/ports/primary"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the user service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	userService primary.UserUseCase,
	roleService primary.RoleUseCase,
	teamService primary.TeamUseCase,
) {
	userHandler := NewUserHandler(userService)
	roleHandler := NewRoleHandler(roleService)
	teamHandler := NewTeamHandler(teamService)
	internalHandler := NewInternalHandler(userService)

	// --- Internal endpoints (service-to-service only) ---
	internal := router.Group("/internal")
	internal.Use(middleware.RequireInternalService())
	{
		internal.GET("/by-email", internalHandler.GetUserByEmail)
		internal.GET("/by-id", internalHandler.GetUserByID)
		internal.GET("/user-name-by-id", internalHandler.GetUserNameByID)
	}

	// --- Authenticated endpoints ---
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		// Current user profile (any authenticated user)
		authenticated.GET("/users/me", userHandler.GetCurrentUser)
		authenticated.POST("/users/me/change-password", userHandler.ChangePassword)

		// User CRUD (privilege-protected)
		users := authenticated.Group("/users")
		{
			users.GET("", middleware.RequirePrivilege(domain.PrivilegeUserView), userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", middleware.RequirePrivilege(domain.PrivilegeUserCreate), userHandler.CreateUser)
			users.PUT("/:id", middleware.RequirePrivilege(domain.PrivilegeUserUpdate), userHandler.UpdateUser)
			users.DELETE("/:id", middleware.RequirePrivilege(domain.PrivilegeUserDelete), userHandler.DeleteUser)
			users.POST("/:id/reset-password", middleware.RequirePrivilege(domain.PrivilegeUserUpdate), userHandler.AdminResetPassword)
		}

		// Role CRUD (privilege-protected)
		roles := authenticated.Group("/roles")
		{
			roles.GET("", middleware.RequirePrivilege(domain.PrivilegeRoleView), roleHandler.ListRoles)
			roles.GET("/:id", middleware.RequirePrivilege(domain.PrivilegeRoleView), roleHandler.GetRole)
			roles.POST("", middleware.RequirePrivilege(domain.PrivilegeRoleCreate), roleHandler.CreateRole)
			roles.PUT("/:id", middleware.RequirePrivilege(domain.PrivilegeRoleUpdate), roleHandler.UpdateRole)
			roles.DELETE("/:id", middleware.RequirePrivilege(domain.PrivilegeRoleDelete), roleHandler.DeleteRole)
			roles.PUT("/:id/privileges", middleware.RequirePrivilege(domain.PrivilegeRoleUpdate), roleHandler.SetRolePrivileges)
			roles.GET("/privileges", middleware.RequireAnyPrivilege(domain.PrivilegeSystemConfig, domain.PrivilegeSystemAdmin), roleHandler.ListPrivileges)
		}

		// Team CRUD (privilege-protected)
		teams := authenticated.Group("/teams")
		{
			teams.GET("", middleware.RequirePrivilege(domain.PrivilegeTeamView), teamHandler.ListTeams)
			teams.GET("/:id", middleware.RequirePrivilege(domain.PrivilegeTeamView), teamHandler.GetTeam)
			teams.POST("", middleware.RequirePrivilege(domain.PrivilegeTeamCreate), teamHandler.CreateTeam)
			teams.PUT("/:id", middleware.RequirePrivilege(domain.PrivilegeTeamUpdate), teamHandler.UpdateTeam)
			teams.DELETE("/:id", middleware.RequirePrivilege(domain.PrivilegeTeamDelete), teamHandler.DeleteTeam)
		}
	}
}
