package domain

import (
	"time"

	"github.com/google/uuid"
)

// EquipmentModel represents the catalog definition (Blueprint) of an equipment.
type EquipmentModel struct {
	ID          uuid.UUID         `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string            `gorm:"column:name;not null;uniqueIndex" json:"name"`
	Category    string            `gorm:"column:category;not null" json:"category"`
	Description string            `gorm:"column:description" json:"description"`
	Thresholds  []MetricThreshold `gorm:"foreignKey:EquipmentModelID" json:"thresholds,omitempty"`
	Suppliers   []ModelSupplier   `gorm:"foreignKey:EquipmentModelID" json:"suppliers,omitempty"`
	CreatedAt   time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"column:updated_at" json:"updated_at"`
}

func (EquipmentModel) TableName() string { return "equipment_models" }

// PartModel represents the catalog definition of a standard part, managing global spare inventory.
type PartModel struct {
	ID            uuid.UUID         `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string            `gorm:"column:name;not null;uniqueIndex" json:"name"`
	Category      string            `gorm:"column:category;not null" json:"category"`
	SpareQuantity int               `gorm:"column:spare_quantity;default:0" json:"spare_quantity"`
	IsSerialized  bool              `gorm:"column:is_serialized;default:false" json:"is_serialized"`
	Thresholds    []MetricThreshold `gorm:"foreignKey:PartModelID" json:"thresholds,omitempty"`
	Suppliers     []ModelSupplier   `gorm:"foreignKey:PartModelID" json:"suppliers,omitempty"`
	CreatedAt     time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"column:updated_at" json:"updated_at"`
}

func (PartModel) TableName() string { return "part_models" }

// EquipmentInstance represents the actual physical equipment.
type EquipmentInstance struct {
	ID               uuid.UUID         `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code             string            `gorm:"column:code;uniqueIndex;not null" json:"code"`
	EquipmentModelID uuid.UUID         `gorm:"column:equipment_model_id;type:uuid;not null" json:"equipment_model_id"`
	EquipmentModel   EquipmentModel    `gorm:"foreignKey:EquipmentModelID" json:"equipment_model,omitempty"`
	Status           string            `gorm:"column:status;not null;default:'OPERATIONAL'" json:"status"`
	Location         string            `gorm:"column:location;not null" json:"location"`
	PurchaseDate     time.Time         `gorm:"column:purchase_date" json:"purchase_date"`
	PurchaseValue    float64           `gorm:"column:purchase_value" json:"purchase_value"`
	Parts            []PartInstance    `gorm:"foreignKey:EquipmentInstanceID" json:"parts,omitempty"`
	Thresholds       []MetricThreshold `gorm:"foreignKey:EquipmentInstanceID" json:"thresholds,omitempty"` // Instance level overrides
	CreatedAt        time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time         `gorm:"column:updated_at" json:"updated_at"`
}

func (EquipmentInstance) TableName() string { return "equipment_instances" }

// PartInstance represents the actual physical part installed on an equipment.
type PartInstance struct {
	ID                  uuid.UUID         `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EquipmentInstanceID *uuid.UUID        `gorm:"column:equipment_instance_id;type:uuid" json:"equipment_instance_id"`
	PartModelID         uuid.UUID         `gorm:"column:part_model_id;type:uuid;not null" json:"part_model_id"`
	PartModel           PartModel         `gorm:"foreignKey:PartModelID" json:"part_model,omitempty"`
	SerialNumber        string            `gorm:"column:serial_number" json:"serial_number"`
	Status              string            `gorm:"column:status;not null;default:'OPERATIONAL'" json:"status"`
	CurrentLocation     string            `gorm:"column:current_location;not null;default:'Warehouse'" json:"current_location"`
	Thresholds          []MetricThreshold `gorm:"foreignKey:PartInstanceID" json:"thresholds,omitempty"` // Instance level overrides
	CreatedAt           time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time         `gorm:"column:updated_at" json:"updated_at"`
}

func (PartInstance) TableName() string { return "part_instances" }

// PartConsumptionLog tracks the usage of non-serialized (consumable) parts.
type PartConsumptionLog struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PartModelID  uuid.UUID `gorm:"column:part_model_id;type:uuid;not null" json:"part_model_id"`
	QuantityUsed int       `gorm:"column:quantity_used;not null" json:"quantity_used"`
	WorkOrderID  *uuid.UUID `gorm:"column:work_order_id;type:uuid" json:"work_order_id"`
	ConsumedBy   uuid.UUID `gorm:"column:consumed_by;type:uuid;not null" json:"consumed_by"`
	Notes        string    `gorm:"column:notes" json:"notes"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (PartConsumptionLog) TableName() string { return "part_consumption_logs" }

// DTOs

type EquipmentModelResponse struct {
	ID          uuid.UUID                 `json:"id"`
	Name        string                    `json:"name"`
	Category    string                    `json:"category"`
	Description string                    `json:"description"`
	Thresholds  []MetricThresholdResponse `json:"thresholds,omitempty"`
	Suppliers   []ModelSupplierResponse   `json:"suppliers,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

type PartModelResponse struct {
	ID            uuid.UUID                 `json:"id"`
	Name          string                    `json:"name"`
	Category      string                    `json:"category"`
	SpareQuantity int                       `json:"spare_quantity"`
	IsSerialized  bool                      `json:"is_serialized"`
	Thresholds    []MetricThresholdResponse `json:"thresholds,omitempty"`
	Suppliers     []ModelSupplierResponse   `json:"suppliers,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

type EquipmentInstanceResponse struct {
	ID               uuid.UUID                   `json:"id"`
	Code             string                      `json:"code"`
	EquipmentModelID uuid.UUID                   `json:"equipment_model_id"`
	EquipmentModel   *EquipmentModelResponse     `json:"equipment_model,omitempty"`
	Status           string                      `json:"status"`
	Location         string                      `json:"location"`
	PurchaseDate     time.Time                   `json:"purchase_date"`
	PurchaseValue    float64                     `json:"purchase_value"`
	Parts            []PartInstanceResponse      `json:"parts,omitempty"`
	Thresholds       []MetricThresholdResponse   `json:"thresholds,omitempty"`
	CreatedAt        time.Time                   `json:"created_at"`
	UpdatedAt        time.Time                   `json:"updated_at"`
}

type PartInstanceResponse struct {
	ID                  uuid.UUID                 `json:"id"`
	EquipmentInstanceID *uuid.UUID                `json:"equipment_instance_id,omitempty"`
	PartModelID         uuid.UUID                 `json:"part_model_id"`
	PartModel           *PartModelResponse        `json:"part_model,omitempty"`
	SerialNumber        string                    `json:"serial_number,omitempty"`
	Status              string                    `json:"status"`
	CurrentLocation     string                    `json:"current_location"`
	Thresholds          []MetricThresholdResponse `json:"thresholds,omitempty"`
	CreatedAt           time.Time                 `json:"created_at"`
	UpdatedAt           time.Time                 `json:"updated_at"`
}

// Converters

func (e *EquipmentModel) ToResponse() EquipmentModelResponse {
	thresh := make([]MetricThresholdResponse, len(e.Thresholds))
	for i, t := range e.Thresholds { thresh[i] = t.ToResponse() }
	sups := make([]ModelSupplierResponse, len(e.Suppliers))
	for i, s := range e.Suppliers { sups[i] = s.ToResponse() }
	return EquipmentModelResponse{
		ID:          e.ID,
		Name:        e.Name,
		Category:    e.Category,
		Description: e.Description,
		Thresholds:  thresh,
		Suppliers:   sups,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (p *PartModel) ToResponse() PartModelResponse {
	thresh := make([]MetricThresholdResponse, len(p.Thresholds))
	for i, t := range p.Thresholds { thresh[i] = t.ToResponse() }
	sups := make([]ModelSupplierResponse, len(p.Suppliers))
	for i, s := range p.Suppliers { sups[i] = s.ToResponse() }
	return PartModelResponse{
		ID:            p.ID,
		Name:          p.Name,
		Category:      p.Category,
		SpareQuantity: p.SpareQuantity,
		IsSerialized:  p.IsSerialized,
		Thresholds:    thresh,
		Suppliers:     sups,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func (e *EquipmentInstance) ToResponse() EquipmentInstanceResponse {
	parts := make([]PartInstanceResponse, len(e.Parts))
	for i, p := range e.Parts { parts[i] = p.ToResponse() }
	thresh := make([]MetricThresholdResponse, len(e.Thresholds))
	for i, t := range e.Thresholds { thresh[i] = t.ToResponse() }
	
	var emResp *EquipmentModelResponse
	if e.EquipmentModel.Name != "" { // simple check if loaded
		r := e.EquipmentModel.ToResponse()
		emResp = &r
	}

	return EquipmentInstanceResponse{
		ID:               e.ID,
		Code:             e.Code,
		EquipmentModelID: e.EquipmentModelID,
		EquipmentModel:   emResp,
		Status:           e.Status,
		Location:         e.Location,
		PurchaseDate:     e.PurchaseDate,
		PurchaseValue:    e.PurchaseValue,
		Parts:            parts,
		Thresholds:       thresh,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

func (p *PartInstance) ToResponse() PartInstanceResponse {
	thresh := make([]MetricThresholdResponse, len(p.Thresholds))
	for i, t := range p.Thresholds { thresh[i] = t.ToResponse() }
	
	var pmResp *PartModelResponse
	if p.PartModel.Name != "" {
		r := p.PartModel.ToResponse()
		pmResp = &r
	}

	return PartInstanceResponse{
		ID:                  p.ID,
		EquipmentInstanceID: p.EquipmentInstanceID,
		PartModelID:         p.PartModelID,
		PartModel:           pmResp,
		SerialNumber:        p.SerialNumber,
		Status:              p.Status,
		CurrentLocation:     p.CurrentLocation,
		Thresholds:          thresh,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}

// Request DTOs

type CreateEquipmentModelRequest struct {
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Description string `json:"description"`
}

type CreatePartModelRequest struct {
	Name          string `json:"name" binding:"required"`
	Category      string `json:"category" binding:"required"`
	SpareQuantity int    `json:"spare_quantity"`
	IsSerialized  bool   `json:"is_serialized"`
}

type CreateEquipmentInstanceRequest struct {
	Code             string    `json:"code" binding:"required"`
	EquipmentModelID uuid.UUID `json:"equipment_model_id" binding:"required"`
	Location         string    `json:"location" binding:"required"`
	PurchaseDate     time.Time `json:"purchase_date"`
	PurchaseValue    float64   `json:"purchase_value"`
}

type CreatePartInstanceRequest struct {
	PartModelID         uuid.UUID `json:"part_model_id" binding:"required"`
	SerialNumber        string    `json:"serial_number" binding:"required"`
	EquipmentInstanceID *string   `json:"equipment_instance_id,omitempty" binding:"omitempty,uuid"`
	CurrentLocation     string    `json:"current_location" binding:"required"`
}

type ConsumePartRequest struct {
	PartModelID  uuid.UUID  `json:"part_model_id" binding:"required"`
	Quantity     int        `json:"quantity" binding:"required,min=1"`
	WorkOrderID  *uuid.UUID `json:"work_order_id,omitempty"`
	Notes        string     `json:"notes"`
}

type MovePartInstanceRequest struct {
	EquipmentInstanceID *string `json:"equipment_instance_id,omitempty" binding:"omitempty,uuid"`
	CurrentLocation     string  `json:"current_location" binding:"required"`
}

type UpdateInstanceStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// Old Asset Response wrapper for backward compatibility, mapping the /assets endpoint 
// If required, we can assemble an AssetResponse from EquipmentInstance.
type AssetResponse struct {
	ID            uuid.UUID                 `json:"id"`
	Type          string                    `json:"type"`
	ParentAssetID *uuid.UUID                `json:"parent_asset_id,omitempty"`
	Name          string                    `json:"name"`
	Code          string                    `json:"code"`
	Status        string                    `json:"status"`
	Category      string                    `json:"category"`
	Location      string                    `json:"location"`
	StockQuantity int                       `json:"stock_quantity"`
	PurchaseDate  time.Time                 `json:"purchase_date"`
	PurchaseValue float64                   `json:"purchase_value"`
	Parts         []AssetResponse           `json:"parts,omitempty"`
	Thresholds    []MetricThresholdResponse `json:"thresholds,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}
