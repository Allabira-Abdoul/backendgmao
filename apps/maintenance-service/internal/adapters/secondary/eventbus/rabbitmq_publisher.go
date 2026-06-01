package eventbus

import (
	"context"

	"backend-gmao/apps/maintenance-service/internal/core/ports/secondary"
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

// PublishWorkOrderStarted publishes an event when a work order starts.
func (p *RabbitMQPublisher) PublishWorkOrderStarted(ctx context.Context, workOrderID uuid.UUID, assetID uuid.UUID, woType string) error {
	payload := map[string]interface{}{
		"work_order_id": workOrderID.String(),
		"asset_id":      assetID.String(),
		"type":          woType, // INTERVENTION or INSPECTION
	}

	event := eventbus.Event{
		Type:    "WORK_ORDER_STARTED",
		Payload: payload,
	}

	return p.bus.Publish(ctx, "maintenance.events", "maintenance.workorder.started", event)
}

// PublishWorkOrderCompleted publishes an event when a work order completes.
func (p *RabbitMQPublisher) PublishWorkOrderCompleted(ctx context.Context, workOrderID uuid.UUID, assetID uuid.UUID, woType string, maintenanceType string) error {
	payload := map[string]interface{}{
		"work_order_id":    workOrderID.String(),
		"asset_id":         assetID.String(),
		"type":             woType,
		"maintenance_type": maintenanceType, // e.g., PALLIATIVE, CURATIVE, SYSTEMATIC
	}

	event := eventbus.Event{
		Type:    "WORK_ORDER_COMPLETED",
		Payload: payload,
	}

	return p.bus.Publish(ctx, "maintenance.events", "maintenance.workorder.completed", event)
}

// PublishAuditLog publishes an audit log event to the RabbitMQ exchange.
func (p *RabbitMQPublisher) PublishAuditLog(ctx context.Context, action string, resourceType string, resourceID string, actorID *uuid.UUID, changes map[string]interface{}) error {
	payload := map[string]interface{}{
		"service_name":  "maintenance-service",
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

	return p.bus.Publish(ctx, "audit.logs", "audit.log.maintenance", event)
}
