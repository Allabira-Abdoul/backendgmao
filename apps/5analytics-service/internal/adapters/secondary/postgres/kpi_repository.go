package postgres

import (
	"context"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"backend-gmao/apps/analytics-service/internal/core/ports/secondary"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type kpiRepository struct {
	db *gorm.DB
}

// NewKpiRepository creates a new KpiRepository.
func NewKpiRepository(db *gorm.DB) secondary.KpiRepository {
	return &kpiRepository{db: db}
}

func (r *kpiRepository) UpsertAssetDim(ctx context.Context, dim *domain.AnalyticsAssetDim) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "asset_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"model_id", "category_name"}),
	}).Create(dim).Error
}

func (r *kpiRepository) InsertStateEvent(ctx context.Context, event *domain.AnalyticsStateEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *kpiRepository) UpsertMaintenanceEvent(ctx context.Context, event *domain.AnalyticsMaintenanceEvent) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "work_order_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "completed_at"}),
	}).Create(event).Error
}

func (r *kpiRepository) RefreshMaterializedViews(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("REFRESH MATERIALIZED VIEW category_health_metrics_daily").Error
}

func (r *kpiRepository) GetCategoryHealthMetrics(ctx context.Context) ([]domain.CategoryHealthMetrics, error) {
	type ViewRow struct {
		CategoryName string
		AssetCount   int
		MttrHours    float64
		Availability float64
	}
	var rows []ViewRow
	err := r.db.WithContext(ctx).Raw("SELECT category_name, asset_count, mttr_hours, availability FROM category_health_metrics_daily").Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	var results []domain.CategoryHealthMetrics
	for _, row := range rows {
		results = append(results, domain.CategoryHealthMetrics{
			CategoryName: row.CategoryName,
			AssetCount:   row.AssetCount,
			Metrics: domain.CoreMetrics{
				Availability: row.Availability,
				MTTR:         row.MttrHours,
			},
		})
	}
	return results, nil
}

func (r *kpiRepository) GetAssetHealthMetrics(ctx context.Context, assetID uuid.UUID) (*domain.AssetHealthMetrics, error) {
	// Placeholder for individual asset metrics
	return &domain.AssetHealthMetrics{AssetID: assetID}, nil
}
