package secondary

import (
	"context"

	"github.com/google/uuid"
)

// EventPublisher defines the port for publishing domain events (like Audit logs).
type EventPublisher interface {
	PublishAuditLog(ctx context.Context, action string, resourceType string, resourceID string, actorID *uuid.UUID, changes map[string]interface{}) error
}
