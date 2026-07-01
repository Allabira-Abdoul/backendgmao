package domain

import (
	"time"

	"backend-gmao/pkg/common"
	"github.com/google/uuid"
)

const (
	DefectStatusPending  = "PENDING"
	DefectStatusApproved = "APPROVED"
	DefectStatusRejected = "REJECTED"
)

// DefectAlert represents a fault reported by a technician.
type DefectAlert struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AssetID     uuid.UUID `gorm:"column:asset_id;type:uuid;not null" json:"asset_id"`
	ReportedBy  uuid.UUID `gorm:"column:reported_by;type:uuid;not null" json:"reported_by"`
	Title       string    `gorm:"column:title;not null" json:"title"`
	Description string    `gorm:"column:description;not null" json:"description"`
	ImageURL    string    `gorm:"column:image_url" json:"image_url"`
	Status      string    `gorm:"column:status;not null;default:'PENDING'" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName overrides GORM's default table name.
func (DefectAlert) TableName() string {
	return "defect_alerts"
}

// DefectAlertResponse is the DTO returned by API endpoints.
type DefectAlertResponse struct {
	ID          uuid.UUID          `json:"id"`
	Asset       common.ResourceRef `json:"asset"`
	ReportedBy  common.ResourceRef `json:"reported_by"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	ImageURL    string             `json:"image_url"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ToResponse converts a DefectAlert to a DefectAlertResponse.
func (d *DefectAlert) ToResponse(assetName string, reportedByName string) DefectAlertResponse {
	return DefectAlertResponse{
		ID:          d.ID,
		Asset:       common.ResourceRef{ID: d.AssetID, Name: assetName},
		ReportedBy:  common.ResourceRef{ID: d.ReportedBy, Name: reportedByName},
		Title:       d.Title,
		Description: d.Description,
		ImageURL:    d.ImageURL,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// ReviewDefectAlertRequest represents the payload for approving/rejecting a defect.
type ReviewDefectAlertRequest struct {
	Status      string `json:"status" binding:"required,oneof=APPROVED REJECTED"`
	ReviewNotes string `json:"review_notes,omitempty"`
}
