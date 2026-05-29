package http

import (
	"errors"
	"net/http"

	"backend-gmao/apps/user-service/internal/application/service"
	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/apps/user-service/internal/core/ports/primary"
	"backend-gmao/pkg/middleware"
	"backend-gmao/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	userService primary.UserUseCase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService primary.UserUseCase) *UserHandler {
	return &UserHandler{userService: userService}
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Extracting pagination validation into response package aligns with Single Responsibility Principle (SRP)
	// Handlers should only orchestrate request data, not validate pagination logic.
	pagination := response.GetPagination(c, 1, 20)

	users, total, err := h.userService.ListUsers(c.Request.Context(), pagination.Limit, pagination.Offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list users")
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, users, response.NewMeta(pagination.Page, pagination.PerPage, total))
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
		return
	}

	// Visibility rule: User can see themselves, but needs USER_VIEW privilege for others.
	currentUserIDStr, exists := c.Get(string(middleware.ContextKeyUserID))
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user ID in token")
		return
	}
	if id.String() != currentUserIDStr.(string) {
		hasPrivilege := false
		if privs, ok := c.Get(string(middleware.ContextKeyPrivileges)); ok {
			if userPrivs, ok := privs.([]string); ok {
				for _, p := range userPrivs {
					if p == "SYSTEM_ADMIN" || p == domain.PrivilegeUserView {
						hasPrivilege = true
						break
					}
				}
			}
		}
		if !hasPrivilege {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "Insufficient privileges to view other users")
			return
		}
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get user")
		return
	}

	response.Success(c, http.StatusOK, user)
}

// GetCurrentUser handles GET /users/me
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userIDStr, _ := c.Get(middleware.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_TOKEN", "Invalid user ID in token")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	response.Success(c, http.StatusOK, user)
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			response.Error(c, http.StatusConflict, "EMAIL_EXISTS", err.Error())
			return
		}
		if errors.Is(err, service.ErrRoleNotFound) {
			response.Error(c, http.StatusBadRequest, "ROLE_NOT_FOUND", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create user")
		return
	}

	response.Success(c, http.StatusCreated, user)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		if errors.Is(err, service.ErrEmailExists) {
			response.Error(c, http.StatusConflict, "EMAIL_EXISTS", err.Error())
			return
		}
		if errors.Is(err, service.ErrRoleNotFound) {
			response.Error(c, http.StatusBadRequest, "ROLE_NOT_FOUND", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update user")
		return
	}

	response.Success(c, http.StatusOK, user)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete user")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// AdminResetPassword handles POST /users/:id/reset-password
func (h *UserHandler) AdminResetPassword(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
		return
	}

	code, err := h.userService.AdminResetPassword(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to reset password")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "Password reset successfully. Please communicate this code to the user.",
		"code":    code,
	})
}

// ChangePassword handles POST /users/me/change-password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDStr, exists := c.Get(string(middleware.ContextKeyUserID))
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user ID in token")
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_TOKEN", "Invalid user ID in token")
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.userService.ChangePassword(c.Request.Context(), userID, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to change password")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Password changed successfully"})
}
