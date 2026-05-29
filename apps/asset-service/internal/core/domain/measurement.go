package domain

import (
	"time"

	"github.com/google/uuid"
)

// Measurement represents a single telemetry reading from an asset.
// This table will be partitioned by RecordedAt.
type Measurement struct {
	ID                  uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	EquipmentInstanceID *uuid.UUID `gorm:"column:equipment_instance_id;type:uuid" json:"equipment_instance_id,omitempty"`
	PartInstanceID      *uuid.UUID `gorm:"column:part_instance_id;type:uuid" json:"part_instance_id,omitempty"`
	MetricName          string     `gorm:"column:metric_name;not null" json:"metric_name"`
	Value               float64    `gorm:"column:value;not null" json:"value"`
	Unit                string     `gorm:"column:unit;not null" json:"unit"`
	RecordedAt          time.Time  `gorm:"column:recorded_at;primaryKey;not null" json:"recorded_at"`
	RecordedBy          *uuid.UUID `gorm:"column:recorded_by;type:uuid" json:"recorded_by,omitempty"`
}

func (Measurement) TableName() string {
	return "measurements"
}

// IngestMeasurementRequest is the payload received from a sensor or technician.
type IngestMeasurementRequest struct {
	EquipmentInstanceID *uuid.UUID `json:"equipment_instance_id,omitempty"`
	PartInstanceID      *uuid.UUID `json:"part_instance_id,omitempty"`
	MetricName          string     `json:"metric_name" binding:"required"`
	Value               float64    `json:"value" binding:"required"`
	Unit                string     `json:"unit" binding:"required"`
	RecordedAt          *time.Time `json:"recorded_at,omitempty"`
}

// MeasurementResponse is the DTO for responding to clients.
type MeasurementResponse struct {
	ID                  uuid.UUID  `json:"id"`
	EquipmentInstanceID *uuid.UUID `json:"equipment_instance_id,omitempty"`
	PartInstanceID      *uuid.UUID `json:"part_instance_id,omitempty"`
	MetricName          string     `json:"metric_name"`
	Value               float64    `json:"value"`
	Unit                string     `json:"unit"`
	RecordedAt          time.Time  `json:"recorded_at"`
	RecordedBy          *uuid.UUID `json:"recorded_by,omitempty"`
}

func (m *Measurement) ToResponse() MeasurementResponse {
	return MeasurementResponse{
		ID:                  m.ID,
		EquipmentInstanceID: m.EquipmentInstanceID,
		PartInstanceID:      m.PartInstanceID,
		MetricName:          m.MetricName,
		Value:               m.Value,
		Unit:                m.Unit,
		RecordedAt:          m.RecordedAt,
		RecordedBy:          m.RecordedBy,
	}
}
