package domain

import (
	"time"

	"github.com/google/uuid"
)

// Site represents the physical location.
type Site struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"column:name;not null;uniqueIndex" json:"name"`
	Location    string    `gorm:"column:location" json:"location"`
	Description string    `gorm:"column:description" json:"description"`
	Systems     []System  `gorm:"foreignKey:SiteID" json:"systems,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Site) TableName() string { return "sites" }

// System represents a major functional group within a site.
type System struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SiteID      uuid.UUID `gorm:"column:site_id;type:uuid;not null;uniqueIndex:idx_sys_site_name" json:"site_id"`
	Name        string    `gorm:"column:name;not null;uniqueIndex:idx_sys_site_name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Status      string    `gorm:"column:status;not null;default:'OPERATIONAL'" json:"status"`
	Assets      []Asset   `gorm:"foreignKey:SystemID" json:"assets,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (System) TableName() string { return "systems" }

// Asset represents a specific piece of equipment within a system.
type Asset struct {
	ID            uuid.UUID   `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SystemID      uuid.UUID   `gorm:"column:system_id;type:uuid;not null;uniqueIndex:idx_ast_sys_name" json:"system_id"`
	Name          string      `gorm:"column:name;not null;uniqueIndex:idx_ast_sys_name" json:"name"`
	Code          string      `gorm:"column:code;uniqueIndex" json:"code"`
	Model         string      `gorm:"column:model" json:"model"`
	Manufacturer  string      `gorm:"column:manufacturer" json:"manufacturer"`
	Status        string      `gorm:"column:status;not null;default:'OPERATIONAL'" json:"status"`
	RulPercentage float64     `gorm:"column:rul_percentage;not null;default:100.0" json:"rul_percentage"`
	Subsystems    []Subsystem `gorm:"foreignKey:AssetID" json:"subsystems,omitempty"`
	CreatedAt     time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

func (Asset) TableName() string { return "assets" }

// Subsystem represents a functional sub-part of an asset.
type Subsystem struct {
	ID          uuid.UUID   `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AssetID     uuid.UUID   `gorm:"column:asset_id;type:uuid;not null;uniqueIndex:idx_subast_name" json:"asset_id"`
	Name        string      `gorm:"column:name;not null;uniqueIndex:idx_subast_name" json:"name"`
	Description string      `gorm:"column:description" json:"description"`
	Criticality string      `gorm:"column:criticality;default:'MEDIUM'" json:"criticality"`
	Components  []Component `gorm:"foreignKey:SubsystemID" json:"components,omitempty"`
	CreatedAt   time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

func (Subsystem) TableName() string { return "subsystems" }

// InventoryItem represents a catalog item (SPARE_PART or CONSUMABLE).
type InventoryItem struct {
	ID                   uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ItemType             string    `gorm:"column:item_type;not null" json:"item_type"` // 'SPARE_PART' or 'CONSUMABLE'
	PartNumber           string    `gorm:"column:part_number;uniqueIndex;not null" json:"part_number"`
	Name                 string    `gorm:"column:name;not null" json:"name"`
	Category             string    `gorm:"column:category" json:"category"`
	StockQuantity        int       `gorm:"column:stock_quantity;default:0" json:"stock_quantity"`
	ReorderPoint         int       `gorm:"column:reorder_point;default:0" json:"reorder_point"`
	SupplierLeadTimeDays int       `gorm:"column:supplier_lead_time_days;default:0" json:"supplier_lead_time_days"`
	UnitOfMeasure        string    `gorm:"column:unit_of_measure;default:'unit'" json:"unit_of_measure"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (InventoryItem) TableName() string { return "inventory_items" }

// Component represents an actual physical instance installed on a subsystem.
type Component struct {
	ID              uuid.UUID     `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SubsystemID     uuid.UUID     `gorm:"column:subsystem_id;type:uuid;not null;uniqueIndex:idx_comp_subsys_name" json:"subsystem_id"`
	InventoryItemID uuid.UUID     `gorm:"column:inventory_item_id;type:uuid;not null" json:"inventory_item_id"`
	InventoryItem   InventoryItem `gorm:"foreignKey:InventoryItemID" json:"inventory_item,omitempty"`
	Name            string        `gorm:"column:name;not null;uniqueIndex:idx_comp_subsys_name" json:"name"`
	SerialNumber    string        `gorm:"column:serial_number" json:"serial_number"`
	Status          string        `gorm:"column:status;not null;default:'OPERATIONAL'" json:"status"`
	InstalledAt     time.Time     `gorm:"column:installed_at;default:CURRENT_TIMESTAMP" json:"installed_at"`
	CreatedAt       time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

func (Component) TableName() string { return "components" }

// -----------------------------------------------------------------------------
// DTOs
// -----------------------------------------------------------------------------

type SiteResponse struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Location    string           `json:"location"`
	Description string           `json:"description"`
	Systems     []SystemResponse `json:"systems,omitempty"`
}

type SystemResponse struct {
	ID          uuid.UUID       `json:"id"`
	SiteID      uuid.UUID       `json:"site_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Assets      []AssetResponse `json:"assets,omitempty"`
}

type AssetResponse struct {
	ID            uuid.UUID           `json:"id"`
	SystemID      uuid.UUID           `json:"system_id"`
	Name          string              `json:"name"`
	Code          string              `json:"code"`
	Model         string              `json:"model"`
	Manufacturer  string              `json:"manufacturer"`
	Status        string              `json:"status"`
	RulPercentage float64             `json:"rul_percentage"`
	Subsystems    []SubsystemResponse `json:"subsystems,omitempty"`
}

type SubsystemResponse struct {
	ID          uuid.UUID           `json:"id"`
	AssetID     uuid.UUID           `json:"asset_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Criticality string              `json:"criticality"`
	Components  []ComponentResponse `json:"components,omitempty"`
}

type ComponentResponse struct {
	ID              uuid.UUID              `json:"id"`
	SubsystemID     uuid.UUID              `json:"subsystem_id"`
	InventoryItemID uuid.UUID              `json:"inventory_item_id"`
	InventoryItem   *InventoryItemResponse `json:"inventory_item,omitempty"`
	Name            string                 `json:"name"`
	SerialNumber    string                 `json:"serial_number"`
	Status          string                 `json:"status"`
	InstalledAt     time.Time              `json:"installed_at"`
}

type InventoryItemResponse struct {
	ID                   uuid.UUID `json:"id"`
	ItemType             string    `json:"item_type"`
	PartNumber           string    `json:"part_number"`
	Name                 string    `json:"name"`
	Category             string    `json:"category"`
	StockQuantity        int       `json:"stock_quantity"`
	ReorderPoint         int       `json:"reorder_point"`
	SupplierLeadTimeDays int       `json:"supplier_lead_time_days"`
	UnitOfMeasure        string    `json:"unit_of_measure"`
}

// -----------------------------------------------------------------------------
// Converters
// -----------------------------------------------------------------------------

func (s *Site) ToResponse() SiteResponse {
	systems := make([]SystemResponse, len(s.Systems))
	for i, sys := range s.Systems {
		systems[i] = sys.ToResponse()
	}
	return SiteResponse{
		ID:          s.ID,
		Name:        s.Name,
		Location:    s.Location,
		Description: s.Description,
		Systems:     systems,
	}
}

func (sys *System) ToResponse() SystemResponse {
	assets := make([]AssetResponse, len(sys.Assets))
	for i, a := range sys.Assets {
		assets[i] = a.ToResponse()
	}
	return SystemResponse{
		ID:          sys.ID,
		SiteID:      sys.SiteID,
		Name:        sys.Name,
		Description: sys.Description,
		Status:      sys.Status,
		Assets:      assets,
	}
}

func (a *Asset) ToResponse() AssetResponse {
	subsystems := make([]SubsystemResponse, len(a.Subsystems))
	for i, sub := range a.Subsystems {
		subsystems[i] = sub.ToResponse()
	}
	return AssetResponse{
		ID:            a.ID,
		SystemID:      a.SystemID,
		Name:          a.Name,
		Code:          a.Code,
		Model:         a.Model,
		Manufacturer:  a.Manufacturer,
		Status:        a.Status,
		RulPercentage: a.RulPercentage,
		Subsystems:    subsystems,
	}
}

func (sub *Subsystem) ToResponse() SubsystemResponse {
	components := make([]ComponentResponse, len(sub.Components))
	for i, c := range sub.Components {
		components[i] = c.ToResponse()
	}
	return SubsystemResponse{
		ID:          sub.ID,
		AssetID:     sub.AssetID,
		Name:        sub.Name,
		Description: sub.Description,
		Criticality: sub.Criticality,
		Components:  components,
	}
}

func (c *Component) ToResponse() ComponentResponse {
	var invResp *InventoryItemResponse
	if c.InventoryItem.PartNumber != "" {
		r := c.InventoryItem.ToResponse()
		invResp = &r
	}
	return ComponentResponse{
		ID:              c.ID,
		SubsystemID:     c.SubsystemID,
		InventoryItemID: c.InventoryItemID,
		InventoryItem:   invResp,
		Name:            c.Name,
		SerialNumber:    c.SerialNumber,
		Status:          c.Status,
		InstalledAt:     c.InstalledAt,
	}
}

func (i *InventoryItem) ToResponse() InventoryItemResponse {
	return InventoryItemResponse{
		ID:                   i.ID,
		ItemType:             i.ItemType,
		PartNumber:           i.PartNumber,
		Name:                 i.Name,
		Category:             i.Category,
		StockQuantity:        i.StockQuantity,
		ReorderPoint:         i.ReorderPoint,
		SupplierLeadTimeDays: i.SupplierLeadTimeDays,
		UnitOfMeasure:        i.UnitOfMeasure,
	}
}

// -----------------------------------------------------------------------------
// Request DTOs
// -----------------------------------------------------------------------------

type CreateSiteRequest struct {
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type CreateSystemRequest struct {
	SiteID      uuid.UUID `json:"site_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
}

type CreateAssetRequest struct {
	SystemID     uuid.UUID `json:"system_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Code         string    `json:"code"`
	Model        string    `json:"model"`
	Manufacturer string    `json:"manufacturer"`
}

type CreateSubsystemRequest struct {
	AssetID     uuid.UUID `json:"asset_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Criticality string    `json:"criticality"`
}

type CreateInventoryItemRequest struct {
	ItemType             string `json:"item_type" binding:"required,oneof=SPARE_PART CONSUMABLE"`
	PartNumber           string `json:"part_number" binding:"required"`
	Name                 string `json:"name" binding:"required"`
	Category             string `json:"category"`
	StockQuantity        int    `json:"stock_quantity"`
	ReorderPoint         int    `json:"reorder_point"`
	SupplierLeadTimeDays int    `json:"supplier_lead_time_days"`
	UnitOfMeasure        string `json:"unit_of_measure"`
}

type CreateComponentRequest struct {
	SubsystemID     uuid.UUID `json:"subsystem_id" binding:"required"`
	InventoryItemID uuid.UUID `json:"inventory_item_id" binding:"required"`
	Name            string    `json:"name" binding:"required"`
	SerialNumber    string    `json:"serial_number"`
}
