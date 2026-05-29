package eventbus

import (
	"context"

	"backend-gmao/apps/auth-service/internal/core/ports/secondary"
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
		"service_name":  "auth-service",
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

	return p.bus.Publish(ctx, "audit.logs", "audit.log.auth", event)
}
