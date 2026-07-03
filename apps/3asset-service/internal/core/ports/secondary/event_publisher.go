package secondary

import (
	"context"

	"github.com/google/uuid"
)

// EventPublisher defines the port for publishing domain events (like Audit logs).
type EventPublisher interface {
	PublishAuditLog(ctx context.Context, action string, resourceType string, resourceID string, actorID *uuid.UUID, changes map[string]interface{}) error
	PublishAssetCreated(ctx context.Context, assetID uuid.UUID, modelID uuid.UUID, categoryName string, tags []string, userID *uuid.UUID) error
	PublishAssetUpdated(ctx context.Context, assetID uuid.UUID, modelID uuid.UUID, categoryName string, tags []string, userID *uuid.UUID) error
	PublishAssetStateChanged(ctx context.Context, assetID uuid.UUID, oldState string, newState string, userID *uuid.UUID) error
}
