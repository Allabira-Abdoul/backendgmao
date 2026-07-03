package domain

import (
	"time"

	"github.com/google/uuid"
)

// MaintenanceSchedule defines a recurring maintenance rule for an asset.
type MaintenanceSchedule struct {
	ID                    uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AssetID               uuid.UUID  `gorm:"column:asset_id;type:uuid;not null" json:"asset_id"`
	Title                 string     `gorm:"column:title;not null" json:"title"`
	Description           string     `gorm:"column:description" json:"description"`
	Frequency             string     `gorm:"column:frequency;not null" json:"frequency"` // ONCE, DAILY, WEEKLY, MONTHLY, YEARLY, USAGE_BASED
	IntervalMonths        *int       `gorm:"column:interval_months" json:"interval_months"`
	IntervalHours         *float64   `gorm:"column:interval_hours" json:"interval_hours"`
	StartDate             *time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate               *time.Time `gorm:"column:end_date" json:"end_date"`
	NextScheduledDate     *time.Time `gorm:"column:next_scheduled_date" json:"next_scheduled_date"`
	NextScheduledUsage    *float64   `gorm:"column:next_scheduled_usage" json:"next_scheduled_usage"`
	MaintenanceCategory   string     `gorm:"column:maintenance_category;not null" json:"maintenance_category"` // PREVENTIVE
	MaintenanceType       string     `gorm:"column:maintenance_type;not null" json:"maintenance_type"`         // SYSTEMATIC, CONDITIONAL
	IsActive              bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	RequireCounterReading bool       `gorm:"column:require_counter_reading;not null;default:false" json:"require_counter_reading"` // for daily inspections
	CreatedAt             time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (MaintenanceSchedule) TableName() string { return "maintenance_schedules" }

// MaintenanceScheduleResponse is the DTO for returning a schedule.
type MaintenanceScheduleResponse struct {
	ID                    uuid.UUID  `json:"id"`
	AssetID               uuid.UUID  `json:"asset_id"`
	Title                 string     `json:"title"`
	Description           string     `json:"description"`
	Frequency             string     `json:"frequency"`
	IntervalMonths        *int       `json:"interval_months,omitempty"`
	IntervalHours         *float64   `json:"interval_hours,omitempty"`
	StartDate             *time.Time `json:"start_date,omitempty"`
	EndDate               *time.Time `json:"end_date,omitempty"`
	NextScheduledDate     *time.Time `json:"next_scheduled_date,omitempty"`
	NextScheduledUsage    *float64   `json:"next_scheduled_usage,omitempty"`
	MaintenanceCategory   string     `json:"maintenance_category"`
	MaintenanceType       string     `json:"maintenance_type"`
	IsActive              bool       `json:"is_active"`
	RequireCounterReading bool       `json:"require_counter_reading"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// ToResponse converts a domain object to a response DTO.
func (m *MaintenanceSchedule) ToResponse() MaintenanceScheduleResponse {
	return MaintenanceScheduleResponse{
		ID:                    m.ID,
		AssetID:               m.AssetID,
		Title:                 m.Title,
		Description:           m.Description,
		Frequency:             m.Frequency,
		IntervalMonths:        m.IntervalMonths,
		IntervalHours:         m.IntervalHours,
		StartDate:             m.StartDate,
		EndDate:               m.EndDate,
		NextScheduledDate:     m.NextScheduledDate,
		NextScheduledUsage:    m.NextScheduledUsage,
		MaintenanceCategory:   m.MaintenanceCategory,
		MaintenanceType:       m.MaintenanceType,
		IsActive:              m.IsActive,
		RequireCounterReading: m.RequireCounterReading,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
}

// CreateMaintenanceScheduleRequest represents the request to create a schedule.
type CreateMaintenanceScheduleRequest struct {
	AssetID               string     `json:"asset_id" binding:"required,uuid"`
	Title                 string     `json:"title" binding:"required"`
	Description           string     `json:"description"`
	Frequency             string     `json:"frequency" binding:"required,oneof=ONCE DAILY WEEKLY MONTHLY YEARLY USAGE_BASED"`
	IntervalMonths        *int       `json:"interval_months"`
	IntervalHours         *float64   `json:"interval_hours"`
	StartDate             *time.Time `json:"start_date"`
	EndDate               *time.Time `json:"end_date"`
	NextScheduledDate     *time.Time `json:"next_scheduled_date"`
	NextScheduledUsage    *float64   `json:"next_scheduled_usage"`
	MaintenanceCategory   string     `json:"maintenance_category" binding:"required"`
	MaintenanceType       string     `json:"maintenance_type" binding:"required"`
	IsActive              *bool      `json:"is_active"`
	RequireCounterReading *bool      `json:"require_counter_reading"`
}

// UpdateMaintenanceScheduleRequest represents the request to update a schedule.
type UpdateMaintenanceScheduleRequest struct {
	Title                 *string    `json:"title"`
	Description           *string    `json:"description"`
	Frequency             *string    `json:"frequency" binding:"omitempty,oneof=ONCE DAILY WEEKLY MONTHLY YEARLY USAGE_BASED"`
	IntervalMonths        *int       `json:"interval_months"`
	IntervalHours         *float64   `json:"interval_hours"`
	StartDate             *time.Time `json:"start_date"`
	EndDate               *time.Time `json:"end_date"`
	NextScheduledDate     *time.Time `json:"next_scheduled_date"`
	NextScheduledUsage    *float64   `json:"next_scheduled_usage"`
	MaintenanceCategory   *string    `json:"maintenance_category"`
	MaintenanceType       *string    `json:"maintenance_type"`
	IsActive              *bool      `json:"is_active"`
	RequireCounterReading *bool      `json:"require_counter_reading"`
}
