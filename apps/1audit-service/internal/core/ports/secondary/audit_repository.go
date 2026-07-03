package secondary

import (
	"context"

	"backend-gmao/apps/audit-service/internal/core/domain"
)

// AuditRepository defines secondary adapter database actions.
// Crucially, it excludes modification and deletion interfaces to guarantee append-only immutability.
type AuditRepository interface {
	Save(ctx context.Context, log domain.AuditLog) error
	Find(ctx context.Context, filter domain.AuditFilter) ([]domain.AuditLog, int64, error)
}
