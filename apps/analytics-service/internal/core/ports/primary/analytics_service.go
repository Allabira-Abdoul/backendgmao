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
	
	ProcessMaintenanceEvent(ctx context.Context, event domain.MaintenanceEvent) error
	GetGlobalKpi(ctx context.Context) (*domain.KpiResponse, error)
	GetCategoryKpi(ctx context.Context, category string) (*domain.KpiResponse, error)
	GetAssetKpi(ctx context.Context, assetID uuid.UUID) (*domain.KpiResponse, error)
}
