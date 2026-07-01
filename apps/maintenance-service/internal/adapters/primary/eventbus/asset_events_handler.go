package eventbus

import (
	"context"
	"encoding/json"
	"log"

	"backend-gmao/apps/maintenance-service/internal/core/ports/primary"
	"backend-gmao/pkg/eventbus"
	"backend-gmao/pkg/middleware"
	"github.com/google/uuid"
)

// AssetEventsHandler listens to asset events.
type AssetEventsHandler struct {
	bus                eventbus.EventBus
	maintenanceService primary.MaintenanceService
}

// NewAssetEventsHandler creates a new AssetEventsHandler.
func NewAssetEventsHandler(bus eventbus.EventBus, maintenanceService primary.MaintenanceService) *AssetEventsHandler {
	return &AssetEventsHandler{
		bus:                bus,
		maintenanceService: maintenanceService,
	}
}

// Start starts listening.
func (h *AssetEventsHandler) Start() error {
	exchange := "asset.events"
	queueName := "maintenance_service_asset_events_queue"
	routingKey := "asset.created"

	err := h.bus.Subscribe(exchange, queueName, routingKey, h.handleAssetCreated)
	if err != nil {
		log.Printf("[MaintenanceService] Failed to subscribe to asset events: %v", err)
		return err
	}

	log.Printf("[MaintenanceService] RabbitMQ Consumer started listening on exchange %s with routing key %s", exchange, routingKey)
	return nil
}

func (h *AssetEventsHandler) handleAssetCreated(event eventbus.Event) {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		log.Printf("[MaintenanceService] Error marshaling event payload: %v", err)
		return
	}

	var payload struct {
		AssetID uuid.UUID `json:"asset_id"`
		ModelID uuid.UUID `json:"model_id"`
		UserID  *string   `json:"user_id"`
	}

	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		log.Printf("[MaintenanceService] Error unmarshaling event payload: %v", err)
		return
	}

	ctx := context.Background()
	if payload.UserID != nil && *payload.UserID != "" {
		ctx = context.WithValue(ctx, middleware.ContextKeyUserID, *payload.UserID)
	}

	if err := h.maintenanceService.HandleAssetCreated(ctx, payload.AssetID, payload.ModelID); err != nil {
		log.Printf("[MaintenanceService] Failed to handle asset created event: %v", err)
	} else {
		log.Printf("[MaintenanceService] Successfully handled asset created event for asset: %s", payload.AssetID)
	}
}
