package http

import (
	"errors"
	"net/http"

	"backend-gmao/apps/user-service/internal/application/service"
	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RoleHandler handles HTTP requests for role operations.
type RoleHandler struct {
	service *service.RoleService
}

// NewRoleHandler creates a new RoleHandler.
func NewRoleHandler(service *service.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

// ListRoles handles GET /roles
func (h *RoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list roles")
		return
	}

	response.Success(c, http.StatusOK, roles)
}

// GetRole handles GET /roles/:id
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid role ID format")
		return
	}

	role, err := h.service.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrRoleNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Role not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get role")
		return
	}

	response.Success(c, http.StatusOK, role)
}

// ListPrivileges handles GET /privileges — returns all system-defined privileges
func (h *RoleHandler) ListPrivileges(c *gin.Context) {
	privileges, err := h.service.ListPrivileges(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list privileges")
		return
	}

	privilegesByDomain, err := h.service.PrivilegesByDomain(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list privileges by domain")
		return
	}
	response.Success(c, http.StatusOK, gin.H{
		"privileges":           privileges,
		"privileges_by_domain": privilegesByDomain,
	})
}

// CreateRole handles POST /roles
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req domain.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	role, err := h.service.CreateRole(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrRoleNameExists) {
			response.Error(c, http.StatusConflict, "ROLE_EXISTS", err.Error())
			return
		}
		if errors.Is(err, service.ErrInvalidPrivileges) {
			response.Error(c, http.StatusBadRequest, "INVALID_PRIVILEGES", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create role")
		return
	}

	response.Success(c, http.StatusCreated, role)
}

// UpdateRole handles PUT /roles/:id
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid role ID format")
		return
	}

	var req domain.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	role, err := h.service.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrRoleNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Role not found")
			return
		}
		if errors.Is(err, service.ErrRoleNameExists) {
			response.Error(c, http.StatusConflict, "ROLE_EXISTS", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update role")
		return
	}

	response.Success(c, http.StatusOK, role)
}

// DeleteRole handles DELETE /roles/:id
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid role ID format")
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrRoleNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Role not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete role")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// SetRolePrivileges handles PUT /roles/:id/privileges
func (h *RoleHandler) SetRolePrivileges(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid role ID format")
		return
	}

	var req domain.SetPrivilegesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	role, err := h.service.SetRolePrivileges(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrRoleNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Role not found")
			return
		}
		if errors.Is(err, service.ErrInvalidPrivileges) {
			response.Error(c, http.StatusBadRequest, "INVALID_PRIVILEGES", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to set privileges")
		return
	}

	response.Success(c, http.StatusOK, role)
}