package domain

import (
	"time"

	"github.com/google/uuid"
)

// Supplier represents an entity that supplies assets (equipment or parts).
type Supplier struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"column:name;not null;uniqueIndex" json:"name"`
	ContactInfo string    `gorm:"column:contact_info" json:"contact_info"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Supplier) TableName() string { return "suppliers" }

// ModelSupplier represents the association between a Supplier and a Model (Equipment or Part).
// It tracks how a specific supplier refers to a specific asset model.
type ModelSupplier struct {
	ID                    uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SupplierID            uuid.UUID  `gorm:"column:supplier_id;type:uuid;not null" json:"supplier_id"`
	Supplier              *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	EquipmentModelID      *uuid.UUID `gorm:"column:equipment_model_id;type:uuid" json:"equipment_model_id,omitempty"`
	PartModelID           *uuid.UUID `gorm:"column:part_model_id;type:uuid" json:"part_model_id,omitempty"`
	ConsumableID          *uuid.UUID `gorm:"column:consumable_id;type:uuid" json:"consumable_id,omitempty"`
	SupplierReferenceCode string     `gorm:"column:supplier_reference_code;not null" json:"supplier_reference_code"`
	TechnicalDocReference string     `gorm:"column:technical_doc_reference" json:"technical_doc_reference"`
	CreatedAt             time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (ModelSupplier) TableName() string { return "model_suppliers" }

// DTOs

type SupplierResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	ContactInfo string    `json:"contact_info"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ModelSupplierResponse struct {
	ID                    uuid.UUID         `json:"id"`
	SupplierID            uuid.UUID         `json:"supplier_id"`
	Supplier              *SupplierResponse `json:"supplier,omitempty"`
	EquipmentModelID      *uuid.UUID        `json:"equipment_model_id,omitempty"`
	PartModelID           *uuid.UUID        `json:"part_model_id,omitempty"`
	ConsumableID          *uuid.UUID        `json:"consumable_id,omitempty"`
	SupplierReferenceCode string            `json:"supplier_reference_code"`
	TechnicalDocReference string            `json:"technical_doc_reference"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

func (s *Supplier) ToResponse() SupplierResponse {
	return SupplierResponse{
		ID:          s.ID,
		Name:        s.Name,
		ContactInfo: s.ContactInfo,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func (ms *ModelSupplier) ToResponse() ModelSupplierResponse {
	var supResp *SupplierResponse
	if ms.Supplier != nil && ms.Supplier.Name != "" {
		r := ms.Supplier.ToResponse()
		supResp = &r
	}

	return ModelSupplierResponse{
		ID:                    ms.ID,
		SupplierID:            ms.SupplierID,
		Supplier:              supResp,
		EquipmentModelID:      ms.EquipmentModelID,
		PartModelID:           ms.PartModelID,
		ConsumableID:          ms.ConsumableID,
		SupplierReferenceCode: ms.SupplierReferenceCode,
		TechnicalDocReference: ms.TechnicalDocReference,
		CreatedAt:             ms.CreatedAt,
		UpdatedAt:             ms.UpdatedAt,
	}
}

// Request DTOs

type CreateSupplierRequest struct {
	Name        string `json:"name" binding:"required"`
	ContactInfo string `json:"contact_info"`
}

type AddModelSupplierRequest struct {
	SupplierID            uuid.UUID `json:"supplier_id" binding:"required"`
	SupplierReferenceCode string    `json:"supplier_reference_code" binding:"required"`
	TechnicalDocReference string    `json:"technical_doc_reference"`
}

type UpdateSupplierRequest struct {
	Name        *string `json:"name"`
	ContactInfo *string `json:"contact_info"`
}
