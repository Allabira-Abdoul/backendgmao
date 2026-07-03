package secondary

import (
	"context"

	"backend-gmao/apps/prediction-service/internal/core/domain"
	"github.com/google/uuid"
)

// PredictionRepository defines secondary adapter database actions.
type PredictionRepository interface {
	Save(ctx context.Context, prediction *domain.Prediction) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Prediction, error)
	FindByAssetID(ctx context.Context, assetID uuid.UUID) ([]domain.Prediction, error)
	FindAll(ctx context.Context) ([]domain.Prediction, error)
}
