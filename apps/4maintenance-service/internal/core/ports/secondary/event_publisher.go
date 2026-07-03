package secondary

import (
	"context"
	"github.com/google/uuid"
)

// EventPublisher defines the port for publishing domain events to the event bus.
type EventPublisher interface {
	PublishWorkOrderStarted(ctx context.Context, workOrderID uuid.UUID, assetID uuid.UUID, woType string) error
	PublishWorkOrderCompleted(ctx context.Context, workOrderID uuid.UUID, assetID uuid.UUID, woType string, maintenanceType string) error
	PublishAuditLog(ctx context.Context, action string, resourceType string, resourceID string, actorID *uuid.UUID, changes map[string]interface{}) error
}
