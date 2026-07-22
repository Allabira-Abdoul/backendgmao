package secondary

import "context"

// AuditLogger defines the interface for logging audit events,
// decoupling the business logic from the specific auditing infrastructure.
type AuditLogger interface {
	LogAction(ctx context.Context, action, details string)
}
