package domain

import (
	"time"

	"backend-gmao/pkg/common"
	"github.com/google/uuid"
)

// MetricMeasurement captures actual readings during an intervention.
type MetricMeasurement struct {
	ID                  uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	InterventionID      *uuid.UUID `gorm:"column:intervention_id;type:uuid" json:"intervention_id"`
	InspectionID        *uuid.UUID `gorm:"column:inspection_id;type:uuid" json:"inspection_id"`
	ComponentID         *uuid.UUID `gorm:"column:component_id;type:uuid" json:"component_id"`
	MetricName          string     `gorm:"column:metric_name;not null" json:"metric_name"`
	Value               float64    `gorm:"column:value;not null" json:"value"`
	Unit                string     `gorm:"column:unit;not null" json:"unit"`
	IsThresholdBreached bool       `gorm:"column:is_threshold_breached;default:false" json:"is_threshold_breached"`
	CreatedAt           time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (MetricMeasurement) TableName() string {
	return "metric_measurements"
}

// MetricMeasurementResponse is the DTO for MetricMeasurement.
type MetricMeasurementResponse struct {
	ID                  uuid.UUID            `json:"id"`
	Intervention        *common.ResourceRef  `json:"intervention,omitempty"`
	Inspection          *common.ResourceRef  `json:"inspection,omitempty"`
	Component           *common.ResourceRef  `json:"component,omitempty"`
	MetricName          string     `json:"metric_name"`
	Value               float64    `json:"value"`
	Unit                string     `json:"unit"`
	IsThresholdBreached bool       `json:"is_threshold_breached"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (m *MetricMeasurement) ToResponse(parentID string, componentName string) MetricMeasurementResponse {
	var compRef *common.ResourceRef
	if m.ComponentID != nil {
		compRef = &common.ResourceRef{ID: *m.ComponentID, Name: componentName}
	}
	
	var invRef *common.ResourceRef
	if m.InterventionID != nil {
		invRef = &common.ResourceRef{ID: *m.InterventionID, Name: "Intervention " + parentID}
	}
	
	var insRef *common.ResourceRef
	if m.InspectionID != nil {
		insRef = &common.ResourceRef{ID: *m.InspectionID, Name: "Inspection " + parentID}
	}

	return MetricMeasurementResponse{
		ID:                  m.ID,
		Intervention:        invRef,
		Inspection:          insRef,
		Component:           compRef,
		MetricName:          m.MetricName,
		Value:               m.Value,
		Unit:                m.Unit,
		IsThresholdBreached: m.IsThresholdBreached,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}
