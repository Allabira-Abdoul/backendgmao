package eventbus

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"backend-gmao/apps/analytics-service/internal/core/ports/secondary"
	"backend-gmao/pkg/eventbus"
	"github.com/google/uuid"
)

type AnalyticsConsumer struct {
	bus     eventbus.EventBus
	kpiRepo secondary.KpiRepository
}

func NewAnalyticsConsumer(bus eventbus.EventBus, kpiRepo secondary.KpiRepository) *AnalyticsConsumer {
	return &AnalyticsConsumer{
		bus:     bus,
		kpiRepo: kpiRepo,
	}
}

func (c *AnalyticsConsumer) Start(ctx context.Context) error {
	handleAsset := func(event eventbus.Event) {
		if err := c.handleAssetEvent(context.Background(), event); err != nil {
			log.Printf("Error handling asset event: %v", err)
		}
	}

	if err := c.bus.Subscribe("asset.events", "analytics-asset-q", "asset.created", handleAsset); err != nil {
		return err
	}
	if err := c.bus.Subscribe("asset.events", "analytics-asset-q", "asset.updated", handleAsset); err != nil {
		return err
	}
	if err := c.bus.Subscribe("asset.events", "analytics-asset-q", "asset.state.changed", handleAsset); err != nil {
		return err
	}

	handleMaintenance := func(event eventbus.Event) {
		if err := c.handleMaintenanceEvent(context.Background(), event); err != nil {
			log.Printf("Error handling maintenance event: %v", err)
		}
	}

	if err := c.bus.Subscribe("maintenance.events", "analytics-maintenance-q", "maintenance.workorder.started", handleMaintenance); err != nil {
		return err
	}
	if err := c.bus.Subscribe("maintenance.events", "analytics-maintenance-q", "maintenance.workorder.completed", handleMaintenance); err != nil {
		return err
	}

	return nil
}

func (c *AnalyticsConsumer) handleAssetEvent(ctx context.Context, event eventbus.Event) error {
	b, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	switch event.Type {
	case "ASSET_CREATED", "ASSET_UPDATED":
		var payload struct {
			AssetID      string   `json:"asset_id"`
			ModelID      string   `json:"model_id"`
			CategoryName string   `json:"category_name"`
			Tags         []string `json:"tags"`
		}
		if err := json.Unmarshal(b, &payload); err != nil {
			return err
		}

		assetID, _ := uuid.Parse(payload.AssetID)
		modelID, _ := uuid.Parse(payload.ModelID)

		dim := &domain.AnalyticsAssetDim{
			AssetID:      assetID,
			ModelID:      modelID,
			CategoryName: payload.CategoryName,
			CreatedAt:    time.Now(),
		}
		return c.kpiRepo.UpsertAssetDim(ctx, dim)

	case "ASSET_STATE_CHANGED":
		var payload struct {
			AssetID  string `json:"asset_id"`
			OldState string `json:"old_state"`
			NewState string `json:"new_state"`
		}
		if err := json.Unmarshal(b, &payload); err != nil {
			return err
		}
		assetID, _ := uuid.Parse(payload.AssetID)
		stateEvent := &domain.AnalyticsStateEvent{
			AssetID:   assetID,
			OldState:  payload.OldState,
			NewState:  payload.NewState,
			Timestamp: time.Now(),
		}
		return c.kpiRepo.InsertStateEvent(ctx, stateEvent)
	}

	return nil
}

func (c *AnalyticsConsumer) handleMaintenanceEvent(ctx context.Context, event eventbus.Event) error {
	b, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	switch event.Type {
	case "WORK_ORDER_STARTED":
		var payload struct {
			WorkOrderID string `json:"work_order_id"`
			AssetID     string `json:"asset_id"`
			Type        string `json:"type"`
		}
		if err := json.Unmarshal(b, &payload); err != nil {
			return err
		}

		woID, _ := uuid.Parse(payload.WorkOrderID)
		assetID, _ := uuid.Parse(payload.AssetID)
		now := time.Now()

		mEvent := &domain.AnalyticsMaintenanceEvent{
			WorkOrderID: woID,
			AssetID:     assetID,
			Type:        payload.Type,
			StartedAt:   &now,
		}
		return c.kpiRepo.UpsertMaintenanceEvent(ctx, mEvent)

	case "WORK_ORDER_COMPLETED":
		var payload struct {
			WorkOrderID     string `json:"work_order_id"`
			AssetID         string `json:"asset_id"`
			Type            string `json:"type"`
			MaintenanceType string `json:"maintenance_type"`
		}
		if err := json.Unmarshal(b, &payload); err != nil {
			return err
		}

		woID, _ := uuid.Parse(payload.WorkOrderID)
		assetID, _ := uuid.Parse(payload.AssetID)
		now := time.Now()

		mEvent := &domain.AnalyticsMaintenanceEvent{
			WorkOrderID:     woID,
			AssetID:         assetID,
			Type:            payload.Type,
			MaintenanceType: payload.MaintenanceType,
			CompletedAt:     &now,
		}
		return c.kpiRepo.UpsertMaintenanceEvent(ctx, mEvent)
	}

	return nil
}
