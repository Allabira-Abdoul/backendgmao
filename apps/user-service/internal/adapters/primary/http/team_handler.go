package http

import (
	"errors"
	"net/http"

	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/apps/user-service/internal/application/service"
	"backend-gmao/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	service *service.TeamService
}

func NewTeamHandler(service *service.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

// ListTeams handles GET /teams
func (h *TeamHandler) ListTeams(c *gin.Context) {
	pagination := response.GetPagination(c, 1, 20)

	teams, total, err := h.service.ListTeams(c.Request.Context(), pagination.Limit, pagination.Offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list teams")
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, teams, response.NewMeta(pagination.Page, pagination.PerPage, total))
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req domain.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	team, err := h.service.CreateTeam(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrTeamNameExists) {
			response.Error(c, http.StatusConflict, "TEAM_EXISTS", err.Error())
			return
		}
		if errors.Is(err, service.ErrRoleNotFound) {
			response.Error(c, http.StatusBadRequest, "ROLE_NOT_FOUND", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create team")
		return
	}

	response.Success(c, http.StatusCreated, team)
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid team ID format")
		return
	}

	team, err := h.service.GetTeamByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Team not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get team")
		return
	}

	response.Success(c, http.StatusOK, team)
}

func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid team ID format")
		return
	}

	var req domain.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	team, err := h.service.UpdateTeam(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Team not found")
			return
		}
		if errors.Is(err, service.ErrTeamNameExists) {
			response.Error(c, http.StatusConflict, "TEAM_EXISTS", err.Error())
			return
		}
		if errors.Is(err, service.ErrRoleNotFound) {
			response.Error(c, http.StatusBadRequest, "ROLE_NOT_FOUND", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update team")
		return
	}

	response.Success(c, http.StatusOK, team)
}

func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid team ID format")
		return
	}

	if err := h.service.DeleteTeam(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Team not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete team")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Team deleted successfully"})
}
