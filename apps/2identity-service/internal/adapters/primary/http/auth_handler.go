package http

import (
	"net/http"

	"backend-gmao/apps/identity-service/internal/core/domain"
	"backend-gmao/apps/identity-service/internal/core/ports/primary"
	"github.com/gin-gonic/gin"
)

// AuthHandler manages HTTP routes for sessions.
type AuthHandler struct {
	authService primary.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService primary.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// CreateSession creates a new login session.
func (h *AuthHandler) CreateSession(c *gin.Context) {
	var req domain.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.CreateSession(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ValidateSession verifies if a session token is active and valid.
func (h *AuthHandler) ValidateSession(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter token is required"})
		return
	}

	resp, err := h.authService.ValidateSession(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RevokeSession revokes an active session token.
func (h *AuthHandler) RevokeSession(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter token is required"})
		return
	}

	err := h.authService.RevokeSession(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked successfully"})
}

// RefreshSession creates a new session using a valid refresh token.
func (h *AuthHandler) RefreshSession(c *gin.Context) {
	var req domain.RefreshSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.RefreshSession(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
