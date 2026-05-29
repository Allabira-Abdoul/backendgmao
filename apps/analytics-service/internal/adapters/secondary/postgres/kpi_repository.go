package postgres

import (
	"context"
	"errors"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"backend-gmao/apps/analytics-service/internal/core/ports/secondary"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type kpiRepository struct {
	db *gorm.DB
}

// NewKpiRepository creates a new KpiRepository.
func NewKpiRepository(db *gorm.DB) secondary.KpiRepository {
	return &kpiRepository{db: db}
}

func (r *kpiRepository) GetByAssetID(ctx context.Context, assetID uuid.UUID) (*domain.AssetKpiState, error) {
	var state domain.AssetKpiState
	result := r.db.WithContext(ctx).Where("asset_id = ?", assetID).First(&state)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, let service handle creation
		}
		return nil, result.Error
	}
	return &state, nil
}

func (r *kpiRepository) Save(ctx context.Context, state *domain.AssetKpiState) error {
	return r.db.WithContext(ctx).Save(state).Error
}

func (r *kpiRepository) GetCategoryAggregates(ctx context.Context, category string) (int, float64, error) {
	var result struct {
		TotalBreakdowns int
		TotalRepairTime float64
	}
	err := r.db.WithContext(ctx).
		Model(&domain.AssetKpiState{}).
		Where("asset_category = ?", category).
		Select("COALESCE(SUM(total_breakdowns), 0) as total_breakdowns, COALESCE(SUM(total_repair_time), 0) as total_repair_time").
		Scan(&result).Error
	return result.TotalBreakdowns, result.TotalRepairTime, err
}

func (r *kpiRepository) GetGlobalAggregates(ctx context.Context) (int, float64, error) {
	var result struct {
		TotalBreakdowns int
		TotalRepairTime float64
	}
	err := r.db.WithContext(ctx).
		Model(&domain.AssetKpiState{}).
		Select("COALESCE(SUM(total_breakdowns), 0) as total_breakdowns, COALESCE(SUM(total_repair_time), 0) as total_repair_time").
		Scan(&result).Error
	return result.TotalBreakdowns, result.TotalRepairTime, err
}

func (r *kpiRepository) GetTotalOperatingTimeByCategory(ctx context.Context, category string) (float64, error) {
	var totalHours float64
	err := r.db.WithContext(ctx).
		Model(&domain.AssetKpiState{}).
		Where("asset_category = ?", category).
		Select("COALESCE(SUM(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - purchase_date)) / 3600), 0)").
		Scan(&totalHours).Error
	return totalHours, err
}

func (r *kpiRepository) GetTotalOperatingTimeGlobal(ctx context.Context) (float64, error) {
	var totalHours float64
	err := r.db.WithContext(ctx).
		Model(&domain.AssetKpiState{}).
		Select("COALESCE(SUM(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - purchase_date)) / 3600), 0)").
		Scan(&totalHours).Error
	return totalHours, err
}
