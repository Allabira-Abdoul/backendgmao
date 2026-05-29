package eventbus

import (
	"context"
	"encoding/json"
	"log"

	"backend-gmao/apps/audit-service/internal/core/domain"
	"backend-gmao/apps/audit-service/internal/core/ports/primary"
	"backend-gmao/pkg/eventbus"
)

// RabbitMQConsumer handles consuming events from RabbitMQ and routing them to the AuditUseCase.
type RabbitMQConsumer struct {
	bus          eventbus.EventBus
	auditUseCase primary.AuditUseCase
}

// NewRabbitMQConsumer creates a new RabbitMQConsumer.
func NewRabbitMQConsumer(bus eventbus.EventBus, auditUseCase primary.AuditUseCase) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		bus:          bus,
		auditUseCase: auditUseCase,
	}
}

// Start starts listening to the audit logs topics.
func (c *RabbitMQConsumer) Start() error {
	exchange := "audit.logs"
	queueName := "audit_service_queue"
	routingKey := "audit.log.*"

	err := c.bus.Subscribe(exchange, queueName, routingKey, c.handleEvent)
	if err != nil {
		log.Printf("Failed to subscribe to audit events: %v", err)
		return err
	}

	log.Printf("RabbitMQ Consumer started listening on exchange %s with routing key %s", exchange, routingKey)
	return nil
}

func (c *RabbitMQConsumer) handleEvent(event eventbus.Event) {
	// Re-marshal the payload to map it to our AuditLog struct.
	// Since event.Payload is unmarshalled into map[string]interface{} by default in the eventbus,
	// we need to serialize and deserialize it to strongly type it.
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		log.Printf("Error marshaling event payload: %v", err)
		return
	}

	var auditLog domain.AuditLog
	if err := json.Unmarshal(payloadBytes, &auditLog); err != nil {
		log.Printf("Error unmarshaling event payload to AuditLog: %v", err)
		return
	}

	// Make sure action is captured from the event type if not explicitly provided in payload
	if auditLog.Action == "" {
		auditLog.Action = event.Type
	}

	ctx := context.Background()
	if err := c.auditUseCase.RecordAction(ctx, auditLog); err != nil {
		log.Printf("Failed to record audit action: %v", err)
	} else {
		log.Printf("Successfully recorded audit action: %s for service: %s", auditLog.Action, auditLog.ServiceName)
	}
}
