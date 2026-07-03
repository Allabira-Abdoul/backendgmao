package secondary

import (
	"context"
	"github.com/google/uuid"
)

// MaintenanceEvent represents the payload sent to analytics-service.
type MaintenanceEvent struct {
	AssetID             uuid.UUID `json:"asset_id"`
	MaintenanceCategory string    `json:"maintenance_category"`
	DurationMinutes     float64   `json:"duration_minutes"`
}

// AnalyticsClient defines the interface to trigger analytics calculations.
type AnalyticsClient interface {
	PublishMaintenanceEvent(ctx context.Context, event MaintenanceEvent) error
}
