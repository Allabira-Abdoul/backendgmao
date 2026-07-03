package primary

import (
	"context"

	"backend-gmao/apps/prediction-service/internal/core/domain"
	"github.com/google/uuid"
)

// PredictionService defines primary application business operations.
type PredictionService interface {
	RecordPrediction(ctx context.Context, req domain.CreatePredictionRequest) (*domain.PredictionResponse, error)
	GetPrediction(ctx context.Context, id uuid.UUID) (*domain.PredictionResponse, error)
	GetPredictionsForAsset(ctx context.Context, assetID uuid.UUID) ([]domain.PredictionResponse, error)
	GetAllPredictions(ctx context.Context) ([]domain.PredictionResponse, error)
}
