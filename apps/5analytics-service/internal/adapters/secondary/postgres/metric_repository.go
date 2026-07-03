package postgres

import (
	"context"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type metricRepository struct {
	db *gorm.DB
}

// NewMetricRepository creates a GORM metric repository.
func NewMetricRepository(db *gorm.DB) *metricRepository {
	return &metricRepository{db: db}
}

func (r *metricRepository) Save(ctx context.Context, metric *domain.Metric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

func (r *metricRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Metric, error) {
	var metric domain.Metric
	if err := r.db.WithContext(ctx).First(&metric, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *metricRepository) FindAll(ctx context.Context) ([]domain.Metric, error) {
	var metrics []domain.Metric
	if err := r.db.WithContext(ctx).Order("timestamp desc").Find(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}

func (r *metricRepository) FindByCategory(ctx context.Context, category string) ([]domain.Metric, error) {
	var metrics []domain.Metric
	if err := r.db.WithContext(ctx).Where("category = ?", category).Order("timestamp desc").Find(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}
