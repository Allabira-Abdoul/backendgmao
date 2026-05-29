package domain

import (
	"time"

	"backend-gmao/pkg/common"
	"github.com/google/uuid"
)

// OrdreTravail represents a work order in the GMAO system.
type OrdreTravail struct {
	ID           uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title        string         `gorm:"column:title;not null" json:"title"`
	Description  string         `gorm:"column:description" json:"description"`
	AssetID      uuid.UUID      `gorm:"column:asset_id;type:uuid;not null" json:"asset_id"`
	Priority            string         `gorm:"column:priority;not null;default:'MEDIUM'" json:"priority"` // LOW, MEDIUM, HIGH, CRITICAL
	Status              string         `gorm:"column:status;not null;default:'PENDING'" json:"status"`    // PENDING, IN_PROGRESS, COMPLETED, CANCELLED
	Type                string         `gorm:"column:type;not null;default:'INTERVENTION'" json:"type"` // INTERVENTION, INSPECTION
	ScheduledAt         *time.Time     `gorm:"column:scheduled_at" json:"scheduled_at"`
	MaintenanceCategory string         `gorm:"column:maintenance_category" json:"maintenance_category"`   // CORRECTIVE, PREVENTIVE
	MaintenanceType     string         `gorm:"column:maintenance_type" json:"maintenance_type"`           // PALLIATIVE, CURATIVE, SYSTEMATIC, CONDITIONAL, PREDICTIVE
	IsMetricMeasurement bool           `gorm:"column:is_metric_measurement;not null;default:false" json:"is_metric_measurement"`
	AssignedTo          *uuid.UUID     `gorm:"column:assigned_to;type:uuid" json:"assigned_to"`
	Interventions       []Intervention `gorm:"foreignKey:WorkOrderID" json:"interventions,omitempty"`
	Inspections         []Inspection   `gorm:"foreignKey:WorkOrderID" json:"inspections,omitempty"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides GORM's default table name.
func (OrdreTravail) TableName() string {
	return "work_orders"
}

// Intervention represents an action taken on a Work Order.
type Intervention struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WorkOrderID         uuid.UUID           `gorm:"column:work_order_id;type:uuid;not null" json:"work_order_id"`
	Description         string              `gorm:"column:description;not null" json:"description"`
	MaintenanceCategory string              `gorm:"column:maintenance_category" json:"maintenance_category"`
	MaintenanceType     string              `gorm:"column:maintenance_type" json:"maintenance_type"`
	IsMetricMeasurement bool                `gorm:"column:is_metric_measurement;not null;default:false" json:"is_metric_measurement"`
	StartedAt           *time.Time          `gorm:"column:started_at" json:"started_at"`
	EndedAt             *time.Time          `gorm:"column:ended_at" json:"ended_at"`
	PerformedBy         uuid.UUID           `gorm:"column:performed_by;type:uuid;not null" json:"performed_by"`
	Measurements        []MetricMeasurement `gorm:"foreignKey:InterventionID" json:"measurements,omitempty"`
	CreatedAt           time.Time           `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time           `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides GORM's default table name.
func (Intervention) TableName() string {
	return "interventions"
}

// Inspection represents an observation/measurement action taken on a Work Order.
type Inspection struct {
	ID           uuid.UUID           `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WorkOrderID  uuid.UUID           `gorm:"column:work_order_id;type:uuid;not null" json:"work_order_id"`
	Observations string              `gorm:"column:observations;not null" json:"observations"`
	StartedAt    *time.Time          `gorm:"column:started_at" json:"started_at"`
	EndedAt      *time.Time          `gorm:"column:ended_at" json:"ended_at"`
	PerformedBy  uuid.UUID           `gorm:"column:performed_by;type:uuid;not null" json:"performed_by"`
	Measurements []MetricMeasurement `gorm:"foreignKey:InspectionID" json:"measurements,omitempty"`
	CreatedAt    time.Time           `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time           `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides GORM's default table name.
func (Inspection) TableName() string {
	return "inspections"
}

// OrdreTravailResponse represents the DTO returned by API endpoints.
type OrdreTravailResponse struct {
	ID                  uuid.UUID              `json:"id"`
	Title               string                 `json:"title"`
	Description         string                 `json:"description"`
	Asset               common.ResourceRef     `json:"asset"`
	Type                string                 `json:"type"`
	ScheduledAt         *time.Time             `json:"scheduled_at,omitempty"`
	Priority            string                 `json:"priority"`
	Status              string                 `json:"status"`
	MaintenanceCategory string                 `json:"maintenance_category"`
	MaintenanceType     string                 `json:"maintenance_type"`
	IsMetricMeasurement bool                   `json:"is_metric_measurement"`
	AssignedTo          *common.ResourceRef    `json:"assigned_to,omitempty"`
	Interventions       []InterventionResponse `json:"interventions,omitempty"`
	Inspections         []InspectionResponse   `json:"inspections,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// InterventionResponse represents the API DTO for Intervention.
type InterventionResponse struct {
	ID                  uuid.UUID                   `json:"id"`
	WorkOrder           common.ResourceRef          `json:"work_order"`
	Description         string                      `json:"description"`
	MaintenanceCategory string                      `json:"maintenance_category"`
	MaintenanceType     string                      `json:"maintenance_type"`
	IsMetricMeasurement bool                        `json:"is_metric_measurement"`
	StartedAt           *time.Time                  `json:"started_at,omitempty"`
	EndedAt             *time.Time                  `json:"ended_at,omitempty"`
	PerformedBy         common.ResourceRef          `json:"performed_by"`
	Measurements        []MetricMeasurementResponse `json:"measurements,omitempty"`
	CreatedAt           time.Time                   `json:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at"`
}

// InspectionResponse represents the API DTO for Inspection.
type InspectionResponse struct {
	ID           uuid.UUID                   `json:"id"`
	WorkOrder    common.ResourceRef          `json:"work_order"`
	Observations string                      `json:"observations"`
	StartedAt    *time.Time                  `json:"started_at,omitempty"`
	EndedAt      *time.Time                  `json:"ended_at,omitempty"`
	PerformedBy  common.ResourceRef          `json:"performed_by"`
	Measurements []MetricMeasurementResponse `json:"measurements,omitempty"`
	CreatedAt    time.Time                   `json:"created_at"`
	UpdatedAt    time.Time                   `json:"updated_at"`
}

// ToResponse converts an Intervention to InterventionResponse.
func (i *Intervention) ToResponse(workOrderName string, performedByName string, compNames map[uuid.UUID]string) InterventionResponse {
	measResp := make([]MetricMeasurementResponse, len(i.Measurements))
	for idx, m := range i.Measurements {
		var compName string
		if m.ComponentID != nil {
			compName = compNames[*m.ComponentID]
		}
		measResp[idx] = m.ToResponse(i.ID.String(), compName) // Wait, we can just pass "" for interventionName if we don't have it here. Let's see.
	}

	return InterventionResponse{
		ID:                  i.ID,
		WorkOrder:           common.ResourceRef{ID: i.WorkOrderID, Name: workOrderName},
		Description:         i.Description,
		MaintenanceCategory: i.MaintenanceCategory,
		MaintenanceType:     i.MaintenanceType,
		IsMetricMeasurement: i.IsMetricMeasurement,
		StartedAt:           i.StartedAt,
		EndedAt:             i.EndedAt,
		PerformedBy:         common.ResourceRef{ID: i.PerformedBy, Name: performedByName},
		Measurements:        measResp,
		CreatedAt:           i.CreatedAt,
		UpdatedAt:           i.UpdatedAt,
	}
}

// ToResponse converts an Inspection to InspectionResponse.
func (i *Inspection) ToResponse(workOrderName string, performedByName string, compNames map[uuid.UUID]string) InspectionResponse {
	measResp := make([]MetricMeasurementResponse, len(i.Measurements))
	for idx, m := range i.Measurements {
		var compName string
		if m.ComponentID != nil {
			compName = compNames[*m.ComponentID]
		}
		measResp[idx] = m.ToResponse(i.ID.String(), compName)
	}

	return InspectionResponse{
		ID:           i.ID,
		WorkOrder:    common.ResourceRef{ID: i.WorkOrderID, Name: workOrderName},
		Observations: i.Observations,
		StartedAt:    i.StartedAt,
		EndedAt:      i.EndedAt,
		PerformedBy:  common.ResourceRef{ID: i.PerformedBy, Name: performedByName},
		Measurements: measResp,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
}

// ToResponse converts an OrdreTravail to OrdreTravailResponse.
func (o *OrdreTravail) ToResponse(assetName string, assignedToName string, perfNames map[uuid.UUID]string, compNames map[uuid.UUID]string) OrdreTravailResponse {
	var assignRef *common.ResourceRef
	if o.AssignedTo != nil {
		assignRef = &common.ResourceRef{ID: *o.AssignedTo, Name: assignedToName}
	}

	invResp := make([]InterventionResponse, len(o.Interventions))
	for i, inv := range o.Interventions {
		invResp[i] = inv.ToResponse(o.Title, perfNames[inv.PerformedBy], compNames)
	}

	insResp := make([]InspectionResponse, len(o.Inspections))
	for i, ins := range o.Inspections {
		insResp[i] = ins.ToResponse(o.Title, perfNames[ins.PerformedBy], compNames)
	}

	return OrdreTravailResponse{
		ID:                  o.ID,
		Title:               o.Title,
		Description:         o.Description,
		Asset:               common.ResourceRef{ID: o.AssetID, Name: assetName},
		Type:                o.Type,
		ScheduledAt:         o.ScheduledAt,
		Priority:            o.Priority,
		Status:              o.Status,
		MaintenanceCategory: o.MaintenanceCategory,
		MaintenanceType:     o.MaintenanceType,
		IsMetricMeasurement: o.IsMetricMeasurement,
		AssignedTo:          assignRef,
		Interventions:       invResp,
		Inspections:         insResp,
		CreatedAt:           o.CreatedAt,
		UpdatedAt:           o.UpdatedAt,
	}
}

// CreateOrdreTravailRequest is the DTO to create a new work order.
type CreateOrdreTravailRequest struct {
	Title       string     `json:"title" binding:"required,min=2,max=255"`
	Description string     `json:"description"`
	AssetID             string     `json:"asset_id" binding:"required,uuid"`
	Type                string     `json:"type" binding:"omitempty,oneof=INTERVENTION INSPECTION"`
	ScheduledAt         *time.Time `json:"scheduled_at,omitempty"`
	Priority            string     `json:"priority" binding:"required"`
	MaintenanceCategory string     `json:"maintenance_category,omitempty"`
	MaintenanceType     string     `json:"maintenance_type,omitempty"`
	IsMetricMeasurement bool       `json:"is_metric_measurement,omitempty"`
	AssignedTo          *string    `json:"assigned_to,omitempty" binding:"omitempty,uuid"`
}

// UpdateOrdreTravailRequest is the DTO to update an existing work order.
type UpdateOrdreTravailRequest struct {
	Title               *string    `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	Description         *string    `json:"description,omitempty"`
	Type                *string    `json:"type,omitempty" binding:"omitempty,oneof=INTERVENTION INSPECTION"`
	ScheduledAt         *time.Time `json:"scheduled_at,omitempty"`
	Status              *string    `json:"status,omitempty" binding:"omitempty,oneof=PENDING IN_PROGRESS COMPLETED CANCELLED"`
	Priority            *string    `json:"priority,omitempty" binding:"omitempty,oneof=LOW MEDIUM HIGH CRITICAL"`
	MaintenanceCategory *string    `json:"maintenance_category,omitempty"`
	MaintenanceType     *string    `json:"maintenance_type,omitempty"`
	IsMetricMeasurement *bool      `json:"is_metric_measurement,omitempty"`
	AssignedTo          *string    `json:"assigned_to,omitempty" binding:"omitempty,uuid"`
}

type CreateMetricMeasurementRequest struct {
	ComponentID         *string `json:"component_id,omitempty" binding:"omitempty,uuid"`
	MetricName          string  `json:"metric_name" binding:"required"`
	Value               float64 `json:"value" binding:"required"`
	Unit                string  `json:"unit" binding:"required"`
	IsThresholdBreached bool    `json:"is_threshold_breached"`
}

// CreateInterventionRequest is the DTO to record a new intervention.
type CreateInterventionRequest struct {
	Description         string                           `json:"description" binding:"required,min=2"`
	MaintenanceCategory string                           `json:"maintenance_category,omitempty"`
	MaintenanceType     string                           `json:"maintenance_type,omitempty"`
	IsMetricMeasurement bool                             `json:"is_metric_measurement,omitempty"`
	PerformedBy         string                           `json:"performed_by" binding:"required,uuid"`
	Measurements        []CreateMetricMeasurementRequest `json:"measurements,omitempty" binding:"omitempty,dive"`
}

// CreateInspectionRequest is the DTO to record a new inspection.
type CreateInspectionRequest struct {
	Observations string                           `json:"observations" binding:"required,min=2"`
	PerformedBy  string                           `json:"performed_by" binding:"required,uuid"`
	Measurements []CreateMetricMeasurementRequest `json:"measurements,omitempty" binding:"omitempty,dive"`
}
