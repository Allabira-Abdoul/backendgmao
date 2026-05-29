package domain

import (
	"time"

	"github.com/google/uuid"
)

// MetricThreshold defines acceptable operational ranges for a model or a specific instance.
type MetricThreshold struct {
	ID                  uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EquipmentModelID    *uuid.UUID `gorm:"column:equipment_model_id;type:uuid" json:"equipment_model_id"`
	PartModelID         *uuid.UUID `gorm:"column:part_model_id;type:uuid" json:"part_model_id"`
	EquipmentInstanceID *uuid.UUID `gorm:"column:equipment_instance_id;type:uuid" json:"equipment_instance_id"`
	PartInstanceID      *uuid.UUID `gorm:"column:part_instance_id;type:uuid" json:"part_instance_id"`
	
	MetricName          string     `gorm:"column:metric_name;not null" json:"metric_name"`
	MinValue            *float64   `gorm:"column:min_value" json:"min_value"`
	MaxValue            *float64   `gorm:"column:max_value" json:"max_value"`
	Unit                string     `gorm:"column:unit;not null" json:"unit"`
	CreatedAt           time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (MetricThreshold) TableName() string {
	return "metric_thresholds"
}

// MetricThresholdResponse is the DTO for MetricThreshold.
type MetricThresholdResponse struct {
	ID         uuid.UUID  `json:"id"`
	MetricName string     `json:"metric_name"`
	MinValue   *float64   `json:"min_value"`
	MaxValue   *float64   `json:"max_value"`
	Unit       string     `json:"unit"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (t *MetricThreshold) ToResponse() MetricThresholdResponse {
	return MetricThresholdResponse{
		ID:         t.ID,
		MetricName: t.MetricName,
		MinValue:   t.MinValue,
		MaxValue:   t.MaxValue,
		Unit:       t.Unit,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}
