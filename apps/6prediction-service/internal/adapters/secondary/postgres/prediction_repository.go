package postgres

import (
	"context"

	"backend-gmao/apps/prediction-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type predictionRepository struct {
	db *gorm.DB
}

// NewPredictionRepository creates a GORM prediction repository.
func NewPredictionRepository(db *gorm.DB) *predictionRepository {
	return &predictionRepository{db: db}
}

func (r *predictionRepository) Save(ctx context.Context, prediction *domain.Prediction) error {
	return r.db.WithContext(ctx).Save(prediction).Error
}

func (r *predictionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Prediction, error) {
	var prediction domain.Prediction
	if err := r.db.WithContext(ctx).First(&prediction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &prediction, nil
}

func (r *predictionRepository) FindByAssetID(ctx context.Context, assetID uuid.UUID) ([]domain.Prediction, error) {
	var predictions []domain.Prediction
	if err := r.db.WithContext(ctx).Where("asset_id = ?", assetID).Find(&predictions).Error; err != nil {
		return nil, err
	}
	return predictions, nil
}

func (r *predictionRepository) FindAll(ctx context.Context) ([]domain.Prediction, error) {
	var predictions []domain.Prediction
	if err := r.db.WithContext(ctx).Find(&predictions).Error; err != nil {
		return nil, err
	}
	return predictions, nil
}
