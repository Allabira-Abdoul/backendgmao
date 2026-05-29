package http

import (
	"net/http"

	"backend-gmao/apps/audit-service/internal/core/domain"
	"backend-gmao/apps/audit-service/internal/core/ports/primary"
	"backend-gmao/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuditHandler manages HTTP routes for audit logs.
type AuditHandler struct {
	auditUseCase primary.AuditUseCase
}

// NewAuditHandler creates a new AuditHandler.
func NewAuditHandler(auditUseCase primary.AuditUseCase) *AuditHandler {
	return &AuditHandler{auditUseCase: auditUseCase}
}

// ListLogs lists all recorded audit events. Only allowed for auditors.
func (h *AuditHandler) ListLogs(c *gin.Context) {
	pagination := response.GetPagination(c, 1, 100)

	filter := domain.AuditFilter{
		ServiceName:  c.Query("service_name"),
		Action:       c.Query("action"),
		ResourceType: c.Query("resource_type"),
		ResourceID:   c.Query("resource_id"),
		Limit:        pagination.Limit,
		Offset:       pagination.Offset,
	}

	resp, total, err := h.auditUseCase.GetLogs(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, resp, response.NewMeta(pagination.Page, pagination.PerPage, total))
}
