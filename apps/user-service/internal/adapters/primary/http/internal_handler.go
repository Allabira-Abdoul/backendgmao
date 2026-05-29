package http

import (
	"net/http"

	"backend-gmao/apps/user-service/internal/core/ports/primary"
	"backend-gmao/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InternalHandler handles service-to-service HTTP requests.
type InternalHandler struct {
	userService primary.UserUseCase
}

// NewInternalHandler creates a new InternalHandler.
func NewInternalHandler(userService primary.UserUseCase) *InternalHandler {
	return &InternalHandler{userService: userService}
}

// GetUserByEmail handles GET /internal/by-email?email=...
// This endpoint is used by authentication-service to retrieve user credentials.
func (h *InternalHandler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.Error(c, http.StatusBadRequest, "MISSING_EMAIL", "Email query parameter is required")
		return
	}

	user, err := h.userService.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	response.Success(c, http.StatusOK, user)
}

// GetUserByID handles GET /internal/by-id?id=...
func (h *InternalHandler) GetUserByID(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "MISSING_ID", "ID query parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid ID format")
		return
	}

	user, err := h.userService.GetUserByIDInternal(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	response.Success(c, http.StatusOK, user)
}

func (h *InternalHandler) GetUserNameByID(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "MISSING_ID", "ID query parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid ID format")
		return
	}

	user, err := h.userService.GetUserByIDInternal(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	response.Success(c, http.StatusOK, user.FullName)
}
