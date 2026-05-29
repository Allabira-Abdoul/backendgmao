package domain

import (
	"time"

	"github.com/google/uuid"
)

// AssetKpiState stores the running totals needed to calculate MTTR and MTBF for an asset.
type AssetKpiState struct {
	AssetID         uuid.UUID `gorm:"column:asset_id;type:uuid;primaryKey" json:"asset_id"`
	AssetCategory   string    `gorm:"column:asset_category;not null" json:"asset_category"`
	PurchaseDate    time.Time `gorm:"column:purchase_date;not null" json:"purchase_date"`
	TotalRepairTime float64   `gorm:"column:total_repair_time;not null;default:0" json:"total_repair_time"` // in hours
	TotalBreakdowns int       `gorm:"column:total_breakdowns;not null;default:0" json:"total_breakdowns"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:current_timestamp" json:"updated_at"`
}

// TableName overrides GORM's default table name.
func (AssetKpiState) TableName() string {
	return "asset_kpi_states"
}

// MaintenanceEvent represents the payload sent from maintenance-service when an intervention is recorded.
type MaintenanceEvent struct {
	AssetID             uuid.UUID `json:"asset_id" binding:"required"`
	MaintenanceCategory string    `json:"maintenance_category" binding:"required"` // e.g. CORRECTIVE
	DurationMinutes     float64   `json:"duration_minutes"`
}

// KpiResponse is the API DTO for returning aggregated KPIs.
type KpiResponse struct {
	Level           string  `json:"level"`            // "global", "category", "asset"
	Identifier      string  `json:"identifier"`       // category name, asset id, or "all"
	MTTR            float64 `json:"mttr"`             // Mean Time To Repair (hours)
	MTBF            float64 `json:"mtbf"`             // Mean Time Between Failures (hours)
	Availability    float64 `json:"availability"`     // Global Availability (%)
	TotalBreakdowns int     `json:"total_breakdowns"` // Number of corrective actions
}
