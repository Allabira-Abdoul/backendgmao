package eventbus

import (
	"context"
	"log"

	"backend-gmao/apps/asset-service/internal/core/ports/primary"
	"backend-gmao/pkg/eventbus"
	"github.com/google/uuid"
)

// StartConsumingWorkOrderEvents subscribes to work order events and updates asset statuses.
func StartConsumingWorkOrderEvents(bus eventbus.EventBus, assetService primary.AssetService) {
	// Subscribe to WORK_ORDER_STARTED
	err := bus.Subscribe("maintenance.events", "asset_service_wo_started", "maintenance.workorder.started", func(e eventbus.Event) {
		payload, ok := e.Payload.(map[string]interface{})
		if !ok {
			log.Println("Invalid payload for WORK_ORDER_STARTED")
			return
		}

		woType, _ := payload["type"].(string)
		if woType != "INTERVENTION" {
			// Inspections don't take equipment offline
			return
		}

		assetIDStr, _ := payload["asset_id"].(string)
		assetID, err := uuid.Parse(assetIDStr)
		if err != nil {
			log.Printf("Invalid asset_id in WORK_ORDER_STARTED: %v", err)
			return
		}

		// Update to OFFLINE
		if err := assetService.UpdateEquipmentStatus(context.Background(), assetID, "OFFLINE"); err != nil {
			log.Printf("Failed to update equipment status to OFFLINE: %v", err)
		} else {
			log.Printf("Set equipment %s to OFFLINE due to intervention start", assetID)
		}
	})

	if err != nil {
		log.Printf("Failed to subscribe to WORK_ORDER_STARTED: %v", err)
	}

	// Subscribe to WORK_ORDER_COMPLETED
	err = bus.Subscribe("maintenance.events", "asset_service_wo_completed", "maintenance.workorder.completed", func(e eventbus.Event) {
		payload, ok := e.Payload.(map[string]interface{})
		if !ok {
			log.Println("Invalid payload for WORK_ORDER_COMPLETED")
			return
		}

		woType, _ := payload["type"].(string)
		if woType != "INTERVENTION" {
			return
		}

		assetIDStr, _ := payload["asset_id"].(string)
		assetID, err := uuid.Parse(assetIDStr)
		if err != nil {
			log.Printf("Invalid asset_id in WORK_ORDER_COMPLETED: %v", err)
			return
		}

		maintenanceType, _ := payload["maintenance_type"].(string)

		newStatus := "OPERATIONAL"
		if maintenanceType == "PALLIATIVE" {
			newStatus = "DEGRADED"
		}

		if err := assetService.UpdateEquipmentStatus(context.Background(), assetID, newStatus); err != nil {
			log.Printf("Failed to update equipment status to %s: %v", newStatus, err)
		} else {
			log.Printf("Set equipment %s to %s due to intervention completion", assetID, newStatus)
		}
	})

	if err != nil {
		log.Printf("Failed to subscribe to WORK_ORDER_COMPLETED: %v", err)
	}
}
