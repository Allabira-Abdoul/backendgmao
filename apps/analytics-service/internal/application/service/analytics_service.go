package service

import (
	"context"
	"time"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"backend-gmao/apps/analytics-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

// AnalyticsService implements primary.AnalyticsService.
type AnalyticsService struct {
	metricRepo  secondary.MetricRepository
	kpiRepo     secondary.KpiRepository
	assetClient secondary.AssetClient
}

// NewAnalyticsService initializes a new AnalyticsService instance.
func NewAnalyticsService(metricRepo secondary.MetricRepository, kpiRepo secondary.KpiRepository, assetClient secondary.AssetClient) *AnalyticsService {
	return &AnalyticsService{
		metricRepo:  metricRepo,
		kpiRepo:     kpiRepo,
		assetClient: assetClient,
	}
}

func (s *AnalyticsService) RecordMetric(ctx context.Context, req domain.CreateMetricRequest) (*domain.MetricResponse, error) {
	metric := &domain.Metric{
		ID:        uuid.New(),
		Name:      req.Name,
		Value:     req.Value,
		Category:  req.Category,
		Timestamp: time.Now(),
	}

	if err := s.metricRepo.Save(ctx, metric); err != nil {
		return nil, err
	}

	resp := metric.ToResponse()
	return &resp, nil
}

func (s *AnalyticsService) GetMetric(ctx context.Context, id uuid.UUID) (*domain.MetricResponse, error) {
	metric, err := s.metricRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := metric.ToResponse()
	return &resp, nil
}

func (s *AnalyticsService) GetAllMetrics(ctx context.Context) ([]domain.MetricResponse, error) {
	metrics, err := s.metricRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.MetricResponse, len(metrics))
	for i, m := range metrics {
		responses[i] = m.ToResponse()
	}
	return responses, nil
}

func (s *AnalyticsService) GetMetricsByCategory(ctx context.Context, category string) ([]domain.MetricResponse, error) {
	metrics, err := s.metricRepo.FindByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.MetricResponse, len(metrics))
	for i, m := range metrics {
		responses[i] = m.ToResponse()
	}
	return responses, nil
}

func (s *AnalyticsService) ProcessMaintenanceEvent(ctx context.Context, event domain.MaintenanceEvent) error {
	// Only care about corrective maintenance for MTTR/MTBF
	if event.MaintenanceCategory != "CORRECTIVE" {
		return nil
	}

	state, err := s.kpiRepo.GetByAssetID(ctx, event.AssetID)
	if err != nil {
		return err
	}

	if state == nil {
		// Fetch asset info from asset-service
		assetInfo, err := s.assetClient.GetAssetInfo(ctx, event.AssetID)
		if err != nil {
			return err
		}

		state = &domain.AssetKpiState{
			AssetID:         assetInfo.ID,
			AssetCategory:   assetInfo.Category,
			PurchaseDate:    assetInfo.PurchaseDate,
			TotalRepairTime: 0,
			TotalBreakdowns: 0,
		}
	}

	state.TotalBreakdowns += 1
	state.TotalRepairTime += event.DurationMinutes / 60.0 // Convert to hours
	state.UpdatedAt = time.Now()

	return s.kpiRepo.Save(ctx, state)
}

func calculateMTTR(totalRepairTime float64, totalBreakdowns int) float64 {
	if totalBreakdowns == 0 {
		return 0
	}
	return totalRepairTime / float64(totalBreakdowns)
}

func calculateMTBF(totalOperatingTime float64, totalRepairTime float64, totalBreakdowns int) float64 {
	if totalBreakdowns == 0 {
		return 0
	}
	activeTime := totalOperatingTime - totalRepairTime
	if activeTime < 0 {
		activeTime = 0
	}
	return activeTime / float64(totalBreakdowns)
}

func calculateAvailability(totalOperatingTime float64, totalRepairTime float64) float64 {
	if totalOperatingTime <= 0 {
		return 100.0 // Default to 100% if no operating time
	}
	activeTime := totalOperatingTime - totalRepairTime
	if activeTime < 0 {
		activeTime = 0
	}
	return (activeTime / totalOperatingTime) * 100.0
}

func (s *AnalyticsService) GetGlobalKpi(ctx context.Context) (*domain.KpiResponse, error) {
	totalBreakdowns, totalRepairTime, err := s.kpiRepo.GetGlobalAggregates(ctx)
	if err != nil {
		return nil, err
	}

	totalOperatingTime, err := s.kpiRepo.GetTotalOperatingTimeGlobal(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.KpiResponse{
		Level:           "global",
		Identifier:      "all",
		MTTR:            calculateMTTR(totalRepairTime, totalBreakdowns),
		MTBF:            calculateMTBF(totalOperatingTime, totalRepairTime, totalBreakdowns),
		Availability:    calculateAvailability(totalOperatingTime, totalRepairTime),
		TotalBreakdowns: totalBreakdowns,
	}, nil
}

func (s *AnalyticsService) GetCategoryKpi(ctx context.Context, category string) (*domain.KpiResponse, error) {
	totalBreakdowns, totalRepairTime, err := s.kpiRepo.GetCategoryAggregates(ctx, category)
	if err != nil {
		return nil, err
	}

	totalOperatingTime, err := s.kpiRepo.GetTotalOperatingTimeByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return &domain.KpiResponse{
		Level:           "category",
		Identifier:      category,
		MTTR:            calculateMTTR(totalRepairTime, totalBreakdowns),
		MTBF:            calculateMTBF(totalOperatingTime, totalRepairTime, totalBreakdowns),
		Availability:    calculateAvailability(totalOperatingTime, totalRepairTime),
		TotalBreakdowns: totalBreakdowns,
	}, nil
}

func (s *AnalyticsService) GetAssetKpi(ctx context.Context, assetID uuid.UUID) (*domain.KpiResponse, error) {
	state, err := s.kpiRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return &domain.KpiResponse{
			Level:           "asset",
			Identifier:      assetID.String(),
			MTTR:            0,
			MTBF:            0,
			Availability:    100.0,
			TotalBreakdowns: 0,
		}, nil
	}

	totalOperatingTime := time.Since(state.PurchaseDate).Hours()

	return &domain.KpiResponse{
		Level:           "asset",
		Identifier:      assetID.String(),
		MTTR:            calculateMTTR(state.TotalRepairTime, state.TotalBreakdowns),
		MTBF:            calculateMTBF(totalOperatingTime, state.TotalRepairTime, state.TotalBreakdowns),
		Availability:    calculateAvailability(totalOperatingTime, state.TotalRepairTime),
		TotalBreakdowns: state.TotalBreakdowns,
	}, nil
}
