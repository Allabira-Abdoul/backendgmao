package secondary

import (
	"context"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"github.com/google/uuid"
)

// MetricRepository defines secondary adapter database actions.
type MetricRepository interface {
	Save(ctx context.Context, metric *domain.Metric) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Metric, error)
	FindAll(ctx context.Context) ([]domain.Metric, error)
	FindByCategory(ctx context.Context, category string) ([]domain.Metric, error)
}
