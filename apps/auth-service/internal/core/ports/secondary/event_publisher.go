package secondary

import (
	"context"

	"github.com/google/uuid"
)

// EventPublisher defines the port for publishing audit events.
type EventPublisher interface {
	PublishAuditLog(ctx context.Context, action string, resourceType string, resourceID string, actorID *uuid.UUID, changes map[string]interface{}) error
}
