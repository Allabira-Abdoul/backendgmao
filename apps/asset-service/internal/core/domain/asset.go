package domain

import (
	"time"

	"github.com/google/uuid"
)

var ValidLocations = []string{"DLA", "NSI", "GOU", "MVR", "BPC", "NGE"}

func IsValidLocation(loc string) bool {
	for _, valid := range ValidLocations {
		if loc == valid {
			return true
		}
	}
	return false
}

// EquipmentModel represents the catalog definition (Blueprint) of an equipment.
type EquipmentModel struct {
	ID          uuid.UUID         `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string            `gorm:"column:name;not null;uniqueIndex" json:"name"`
	Category    string            `gorm:"column:category;not null" json:"category"`
	Description string            `gorm:"column:description" json:"description"`
	Suppliers        []ModelSupplier               `gorm:"foreignKey:EquipmentModelID" json:"suppliers,omitempty"`
	PartRequirements []EquipmentModelPartRequirement `gorm:"foreignKey:EquipmentModelID" json:"part_requirements,omitempty"`
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
	Suppliers     []ModelSupplier   `gorm:"foreignKey:PartModelID" json:"suppliers,omitempty"`
	CreatedAt     time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"column:updated_at" json:"updated_at"`
}

func (PartModel) TableName() string { return "part_models" }

// EquipmentModelPartRequirement defines the required parts for an equipment model.
type EquipmentModelPartRequirement struct {
	ID               uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EquipmentModelID uuid.UUID `gorm:"column:equipment_model_id;type:uuid;not null;uniqueIndex:idx_eq_part_req" json:"equipment_model_id"`
	PartModelID      uuid.UUID `gorm:"column:part_model_id;type:uuid;not null;uniqueIndex:idx_eq_part_req" json:"part_model_id"`
	PartModel        PartModel `gorm:"foreignKey:PartModelID" json:"part_model,omitempty"`
	Quantity         int       `gorm:"column:quantity;not null;default:1" json:"quantity"`
}

func (EquipmentModelPartRequirement) TableName() string { return "equipment_model_part_requirements" }

// EquipmentInstance represents the actual physical equipment.
type EquipmentInstance struct {
	ID               uuid.UUID         `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code             string            `gorm:"column:code;uniqueIndex;not null" json:"code"`
	EquipmentModelID uuid.UUID         `gorm:"column:equipment_model_id;type:uuid;not null" json:"equipment_model_id"`
	EquipmentModel   EquipmentModel    `gorm:"foreignKey:EquipmentModelID" json:"equipment_model,omitempty"`
	SupplierID       *uuid.UUID        `gorm:"column:supplier_id;type:uuid" json:"supplier_id,omitempty"`
	Supplier         *Supplier         `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Status           string            `gorm:"column:status;not null;default:'OPERATIONAL'" json:"status"`
	Location         string            `gorm:"column:location;not null" json:"location"`
	Parts            []PartInstance    `gorm:"foreignKey:EquipmentInstanceID" json:"parts,omitempty"`
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
	SupplierID          *uuid.UUID        `gorm:"column:supplier_id;type:uuid" json:"supplier_id,omitempty"`
	Supplier            *Supplier         `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	SerialNumber        string            `gorm:"column:serial_number" json:"serial_number"`
	Status              string            `gorm:"column:status;not null;default:'OPERATIONAL'" json:"status"`
	CurrentLocation     string            `gorm:"column:current_location;not null;default:'Warehouse'" json:"current_location"`
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
	Suppliers        []ModelSupplierResponse               `json:"suppliers,omitempty"`
	PartRequirements []EquipmentModelPartRequirementResponse `json:"part_requirements,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

type PartModelResponse struct {
	ID            uuid.UUID                 `json:"id"`
	Name          string                    `json:"name"`
	Category      string                    `json:"category"`
	SpareQuantity int                       `json:"spare_quantity"`
	IsSerialized  bool                      `json:"is_serialized"`
	Suppliers     []ModelSupplierResponse   `json:"suppliers,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

type EquipmentModelPartRequirementResponse struct {
	ID               uuid.UUID          `json:"id"`
	EquipmentModelID uuid.UUID          `json:"equipment_model_id"`
	PartModelID      uuid.UUID          `json:"part_model_id"`
	PartModel        *PartModelResponse `json:"part_model,omitempty"`
	Quantity         int                `json:"quantity"`
}

type EquipmentInstanceResponse struct {
	ID               uuid.UUID                   `json:"id"`
	Code             string                      `json:"code"`
	EquipmentModelID uuid.UUID                   `json:"equipment_model_id"`
	EquipmentModel   *EquipmentModelResponse     `json:"equipment_model,omitempty"`
	SupplierID       *uuid.UUID                  `json:"supplier_id,omitempty"`
	Supplier         *SupplierResponse           `json:"supplier,omitempty"`
	Status           string                      `json:"status"`
	Location         string                      `json:"location"`
	Parts            []PartInstanceResponse      `json:"parts,omitempty"`
	CreatedAt        time.Time                   `json:"created_at"`
	UpdatedAt        time.Time                   `json:"updated_at"`
}

type PartInstanceResponse struct {
	ID                  uuid.UUID                 `json:"id"`
	EquipmentInstanceID *uuid.UUID                `json:"equipment_instance_id,omitempty"`
	PartModelID         uuid.UUID                 `json:"part_model_id"`
	PartModel           *PartModelResponse        `json:"part_model,omitempty"`
	SupplierID          *uuid.UUID                `json:"supplier_id,omitempty"`
	Supplier            *SupplierResponse         `json:"supplier,omitempty"`
	SerialNumber        string                    `json:"serial_number,omitempty"`
	Status              string                    `json:"status"`
	CurrentLocation     string                    `json:"current_location"`
	CreatedAt           time.Time                 `json:"created_at"`
	UpdatedAt           time.Time                 `json:"updated_at"`
}

// Converters

func (e *EquipmentModel) ToResponse() EquipmentModelResponse {
	sups := make([]ModelSupplierResponse, len(e.Suppliers))
	for i, s := range e.Suppliers { sups[i] = s.ToResponse() }
	reqs := make([]EquipmentModelPartRequirementResponse, len(e.PartRequirements))
	for i, r := range e.PartRequirements { reqs[i] = r.ToResponse() }
	return EquipmentModelResponse{
		ID:               e.ID,
		Name:             e.Name,
		Category:         e.Category,
		Description:      e.Description,
		Suppliers:        sups,
		PartRequirements: reqs,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (p *PartModel) ToResponse() PartModelResponse {
	sups := make([]ModelSupplierResponse, len(p.Suppliers))
	for i, s := range p.Suppliers { sups[i] = s.ToResponse() }
	return PartModelResponse{
		ID:            p.ID,
		Name:          p.Name,
		Category:      p.Category,
		SpareQuantity: p.SpareQuantity,
		IsSerialized:  p.IsSerialized,
		Suppliers:     sups,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func (r *EquipmentModelPartRequirement) ToResponse() EquipmentModelPartRequirementResponse {
	var pmResp *PartModelResponse
	if r.PartModel.Name != "" {
		pm := r.PartModel.ToResponse()
		pmResp = &pm
	}
	return EquipmentModelPartRequirementResponse{
		ID:               r.ID,
		EquipmentModelID: r.EquipmentModelID,
		PartModelID:      r.PartModelID,
		PartModel:        pmResp,
		Quantity:         r.Quantity,
	}
}

func (e *EquipmentInstance) ToResponse() EquipmentInstanceResponse {
	parts := make([]PartInstanceResponse, len(e.Parts))
	for i, p := range e.Parts { parts[i] = p.ToResponse() }
	
	var emResp *EquipmentModelResponse
	if e.EquipmentModel.Name != "" { // simple check if loaded
		r := e.EquipmentModel.ToResponse()
		emResp = &r
	}

	var supResp *SupplierResponse
	if e.Supplier != nil && e.Supplier.Name != "" {
		s := e.Supplier.ToResponse()
		supResp = &s
	}

	return EquipmentInstanceResponse{
		ID:               e.ID,
		Code:             e.Code,
		EquipmentModelID: e.EquipmentModelID,
		EquipmentModel:   emResp,
		SupplierID:       e.SupplierID,
		Supplier:         supResp,
		Status:           e.Status,
		Location:         e.Location,
		Parts:            parts,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

func (p *PartInstance) ToResponse() PartInstanceResponse {
	
	var pmResp *PartModelResponse
	if p.PartModel.Name != "" {
		r := p.PartModel.ToResponse()
		pmResp = &r
	}

	var supResp *SupplierResponse
	if p.Supplier != nil && p.Supplier.Name != "" {
		s := p.Supplier.ToResponse()
		supResp = &s
	}

	return PartInstanceResponse{
		ID:                  p.ID,
		EquipmentInstanceID: p.EquipmentInstanceID,
		PartModelID:         p.PartModelID,
		PartModel:           pmResp,
		SupplierID:          p.SupplierID,
		Supplier:            supResp,
		SerialNumber:        p.SerialNumber,
		Status:              p.Status,
		CurrentLocation:     p.CurrentLocation,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}

// Request DTOs

type PartRequirementReq struct {
	PartModelID uuid.UUID `json:"part_model_id"`
	Quantity    int       `json:"quantity"`
}

type CreateEquipmentModelRequest struct {
	Name             string               `json:"name" binding:"required"`
	Category         string               `json:"category" binding:"required"`
	Description      string               `json:"description"`
	PartRequirements []PartRequirementReq `json:"part_requirements"`
}

type CreatePartModelRequest struct {
	Name          string `json:"name" binding:"required"`
	Category      string `json:"category" binding:"required"`
	SpareQuantity int    `json:"spare_quantity"`
	IsSerialized  bool   `json:"is_serialized"`
}

type UpdateEquipmentModelRequest struct {
	Name             *string              `json:"name"`
	Category         *string              `json:"category"`
	Description      *string              `json:"description"`
	PartRequirements []PartRequirementReq `json:"part_requirements"`
}

type UpdateEquipmentLocationRequest struct {
	Location string `json:"location" binding:"required"`
}

type UpdatePartModelRequest struct {
	Name          *string `json:"name"`
	Category      *string `json:"category"`
	SpareQuantity *int    `json:"spare_quantity"`
	IsSerialized  *bool   `json:"is_serialized"`
}

type CreateEquipmentInstanceRequest struct {
	Code             string     `json:"code" binding:"required"`
	EquipmentModelID uuid.UUID  `json:"equipment_model_id" binding:"required"`
	SupplierID       *uuid.UUID `json:"supplier_id,omitempty"`
	Location         string     `json:"location" binding:"required"`
}

type CreatePartInstanceRequest struct {
	PartModelID         uuid.UUID  `json:"part_model_id" binding:"required"`
	SupplierID          *uuid.UUID `json:"supplier_id,omitempty"`
	SerialNumber        string     `json:"serial_number" binding:"required"`
	EquipmentInstanceID *string    `json:"equipment_instance_id,omitempty" binding:"omitempty,uuid"`
	CurrentLocation     string     `json:"current_location" binding:"required"`
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
	Parts         []AssetResponse           `json:"parts,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}
