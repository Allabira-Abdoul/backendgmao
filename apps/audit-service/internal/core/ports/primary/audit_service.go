package primary

import (
	"context"

	"backend-gmao/apps/audit-service/internal/core/domain"
)

// AuditUseCase defines primary application business operations for Audit Logs.
type AuditUseCase interface {
	RecordAction(ctx context.Context, log domain.AuditLog) error
	GetLogs(ctx context.Context, filter domain.AuditFilter) ([]domain.AuditLogResponse, int64, error)
}
