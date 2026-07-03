package primary

import (
	"context"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"github.com/google/uuid"
)

// AnalyticsService defines primary application business operations.
type AnalyticsService interface {
	RecordMetric(ctx context.Context, req domain.CreateMetricRequest) (*domain.MetricResponse, error)
	GetMetric(ctx context.Context, id uuid.UUID) (*domain.MetricResponse, error)
	GetAllMetrics(ctx context.Context) ([]domain.MetricResponse, error)
	GetMetricsByCategory(ctx context.Context, category string) ([]domain.MetricResponse, error)
	GetCategoryHealthMetrics(ctx context.Context) ([]domain.CategoryHealthMetrics, error)
	StartBackgroundRefresher(ctx context.Context)
}
