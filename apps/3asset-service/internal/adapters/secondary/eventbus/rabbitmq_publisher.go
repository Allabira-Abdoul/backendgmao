package eventbus

import (
	"context"

	"backend-gmao/apps/asset-service/internal/core/ports/secondary"
	"backend-gmao/pkg/eventbus"
	"github.com/google/uuid"
)

// RabbitMQPublisher implements the EventPublisher port using RabbitMQ.
type RabbitMQPublisher struct {
	bus eventbus.EventBus
}

// NewRabbitMQPublisher creates a new EventPublisher using RabbitMQ.
func NewRabbitMQPublisher(bus eventbus.EventBus) secondary.EventPublisher {
	return &RabbitMQPublisher{
		bus: bus,
	}
}

// PublishAuditLog publishes an audit log event to the RabbitMQ exchange.
func (p *RabbitMQPublisher) PublishAuditLog(ctx context.Context, action string, resourceType string, resourceID string, actorID *uuid.UUID, changes map[string]interface{}) error {
	payload := map[string]interface{}{
		"service_name":  "asset-service",
		"action":        action,
		"resource_type": resourceType,
		"resource_id":   resourceID,
		"changes":       changes,
	}

	if actorID != nil {
		payload["actor_id"] = actorID.String()
	}

	event := eventbus.Event{
		Type:    action,
		Payload: payload,
	}

	return p.bus.Publish(ctx, "audit.logs", "audit.log.asset", event)
}

func (p *RabbitMQPublisher) PublishAssetCreated(ctx context.Context, assetID uuid.UUID, modelID uuid.UUID, categoryName string, tags []string, userID *uuid.UUID) error {
	payload := map[string]interface{}{
		"asset_id":      assetID.String(),
		"model_id":      modelID.String(),
		"category_name": categoryName,
		"tags":          tags,
	}
	if userID != nil {
		payload["user_id"] = userID.String()
	}

	event := eventbus.Event{
		Type:    "ASSET_CREATED",
		Payload: payload,
	}
	return p.bus.Publish(ctx, "asset.events", "asset.created", event)
}

func (p *RabbitMQPublisher) PublishAssetUpdated(ctx context.Context, assetID uuid.UUID, modelID uuid.UUID, categoryName string, tags []string, userID *uuid.UUID) error {
	payload := map[string]interface{}{
		"asset_id":      assetID.String(),
		"model_id":      modelID.String(),
		"category_name": categoryName,
		"tags":          tags,
	}
	if userID != nil {
		payload["user_id"] = userID.String()
	}

	event := eventbus.Event{
		Type:    "ASSET_UPDATED",
		Payload: payload,
	}
	return p.bus.Publish(ctx, "asset.events", "asset.updated", event)
}

func (p *RabbitMQPublisher) PublishAssetStateChanged(ctx context.Context, assetID uuid.UUID, oldState string, newState string, userID *uuid.UUID) error {
	payload := map[string]interface{}{
		"asset_id":  assetID.String(),
		"old_state": oldState,
		"new_state": newState,
	}
	if userID != nil {
		payload["user_id"] = userID.String()
	}

	event := eventbus.Event{
		Type:    "ASSET_STATE_CHANGED",
		Payload: payload,
	}
	return p.bus.Publish(ctx, "asset.events", "asset.state.changed", event)
}
