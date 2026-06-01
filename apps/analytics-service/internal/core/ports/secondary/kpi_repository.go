package secondary

import (
	"context"
	"backend-gmao/apps/analytics-service/internal/core/domain"
	"github.com/google/uuid"
)

// KpiRepository defines the persistence interface for CQRS Analytics.
type KpiRepository interface {
	UpsertAssetDim(ctx context.Context, dim *domain.AnalyticsAssetDim) error
	InsertStateEvent(ctx context.Context, event *domain.AnalyticsStateEvent) error
	UpsertMaintenanceEvent(ctx context.Context, event *domain.AnalyticsMaintenanceEvent) error

	RefreshMaterializedViews(ctx context.Context) error
	GetCategoryHealthMetrics(ctx context.Context) ([]domain.CategoryHealthMetrics, error)
	GetAssetHealthMetrics(ctx context.Context, assetID uuid.UUID) (*domain.AssetHealthMetrics, error)
}
