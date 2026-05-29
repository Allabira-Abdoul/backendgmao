package secondary

import (
	"context"
	"backend-gmao/apps/analytics-service/internal/core/domain"
	"github.com/google/uuid"
)

// KpiRepository defines the persistence interface for Asset KPI states.
type KpiRepository interface {
	GetByAssetID(ctx context.Context, assetID uuid.UUID) (*domain.AssetKpiState, error)
	Save(ctx context.Context, state *domain.AssetKpiState) error
	GetCategoryAggregates(ctx context.Context, category string) (totalBreakdowns int, totalRepairTime float64, err error)
	GetGlobalAggregates(ctx context.Context) (totalBreakdowns int, totalRepairTime float64, err error)
	GetTotalOperatingTimeByCategory(ctx context.Context, category string) (totalHours float64, err error)
	GetTotalOperatingTimeGlobal(ctx context.Context) (totalHours float64, err error)
}
