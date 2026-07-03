package domain

import (
	"time"

	"backend-gmao/pkg/common"
	"github.com/google/uuid"
)

// Prediction represents an AI/ML health prediction for a machine asset.
type Prediction struct {
	ID                   uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AssetID              uuid.UUID `gorm:"column:asset_id;type:uuid;not null" json:"asset_id"`
	FailureProbability   float64   `gorm:"column:failure_probability;not null" json:"failure_probability"`
	PredictedFailureDate time.Time `gorm:"column:predicted_failure_date" json:"predicted_failure_date"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName overrides GORM's default table name.
func (Prediction) TableName() string {
	return "predictions"
}

// PredictionResponse represents the API DTO for Prediction.
type PredictionResponse struct {
	ID                   uuid.UUID          `json:"id"`
	Asset                common.ResourceRef `json:"asset"`
	FailureProbability   float64   `json:"failure_probability"`
	PredictedFailureDate time.Time `json:"predicted_failure_date"`
	CreatedAt            time.Time `json:"created_at"`
}

// ToResponse converts a Prediction to PredictionResponse DTO.
func (p *Prediction) ToResponse(assetName string) PredictionResponse {
	return PredictionResponse{
		ID:                   p.ID,
		Asset:                common.ResourceRef{ID: p.AssetID, Name: assetName},
		FailureProbability:   p.FailureProbability,
		PredictedFailureDate: p.PredictedFailureDate,
		CreatedAt:            p.CreatedAt,
	}
}

// CreatePredictionRequest is the DTO used to submit a new prediction.
type CreatePredictionRequest struct {
	AssetID              string    `json:"asset_id" binding:"required,uuid"`
	FailureProbability   float64   `json:"failure_probability" binding:"required,gte=0,lte=100"`
	PredictedFailureDate time.Time `json:"predicted_failure_date" binding:"required"`
}
