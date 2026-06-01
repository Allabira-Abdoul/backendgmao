package domain

import (
	"time"

	"github.com/google/uuid"
)

// TimeRange standardizes the reporting windows.
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// CoreMetrics holds the mathematical results regardless of the entity.
type CoreMetrics struct {
	TotalUptime   float64 `json:"total_uptime_seconds"`
	TotalDowntime float64 `json:"total_downtime_seconds"`
	Availability  float64 `json:"availability_percentage"`
	MTTR          float64 `json:"mttr_hours"`
	MTBF          float64 `json:"mtbf_hours"`
}

// AssetHealthMetrics represents the aggregated KPI for a specific asset.
type AssetHealthMetrics struct {
	AssetID uuid.UUID   `json:"asset_id"`
	Period  TimeRange   `json:"period"`
	Metrics CoreMetrics `json:"metrics"`
}

// CategoryHealthMetrics represents the aggregated KPI for an entire category.
type CategoryHealthMetrics struct {
	CategoryName string      `json:"category_name"` // e.g., "Ground Support Equipment"
	AssetCount   int         `json:"asset_count"`   // Number of assets in this category
	Period       TimeRange   `json:"period"`
	Metrics      CoreMetrics `json:"metrics"`
}

// AnalyticsAssetDim represents the dimension table for assets.
type AnalyticsAssetDim struct {
	AssetID      uuid.UUID `gorm:"column:asset_id;type:uuid;primaryKey"`
	ModelID      uuid.UUID `gorm:"column:model_id;type:uuid"`
	CategoryName string    `gorm:"column:category_name;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:current_timestamp"`
}

// TableName overrides GORM's default table name.
func (AnalyticsAssetDim) TableName() string {
	return "analytics_asset_dim"
}

// AnalyticsStateEvent represents a raw state change event.
type AnalyticsStateEvent struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	AssetID   uuid.UUID `gorm:"column:asset_id;type:uuid;not null"`
	OldState  string    `gorm:"column:old_state;not null"`
	NewState  string    `gorm:"column:new_state;not null"`
	Timestamp time.Time `gorm:"column:timestamp;not null"`
}

// TableName overrides GORM's default table name.
func (AnalyticsStateEvent) TableName() string {
	return "analytics_state_events"
}

// AnalyticsMaintenanceEvent represents a maintenance work order event.
type AnalyticsMaintenanceEvent struct {
	WorkOrderID     uuid.UUID  `gorm:"column:work_order_id;type:uuid;primaryKey"`
	AssetID         uuid.UUID  `gorm:"column:asset_id;type:uuid;not null"`
	Type            string     `gorm:"column:type;not null"`
	MaintenanceType string     `gorm:"column:maintenance_type"`
	StartedAt       *time.Time `gorm:"column:started_at"`
	CompletedAt     *time.Time `gorm:"column:completed_at"`
	// Denormalized fields to simplify the materialized view query
	UptimeSeconds   float64    `gorm:"column:uptime_seconds;default:0"`
	DowntimeSeconds float64    `gorm:"column:downtime_seconds;default:0"`
}

// TableName overrides GORM's default table name.
func (AnalyticsMaintenanceEvent) TableName() string {
	return "analytics_maintenance_events"
}
