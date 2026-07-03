package service

import (
	"context"
	"log"
	"time"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"backend-gmao/apps/analytics-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

// AnalyticsService implements primary.AnalyticsService.
type AnalyticsService struct {
	metricRepo     secondary.MetricRepository
	kpiRepo        secondary.KpiRepository
	assetClient    secondary.AssetClient
	eventPublisher secondary.EventPublisher
}

// NewAnalyticsService initializes a new AnalyticsService instance.
func NewAnalyticsService(metricRepo secondary.MetricRepository, kpiRepo secondary.KpiRepository, assetClient secondary.AssetClient, eventPublisher secondary.EventPublisher) *AnalyticsService {
	return &AnalyticsService{
		metricRepo:     metricRepo,
		kpiRepo:        kpiRepo,
		assetClient:    assetClient,
		eventPublisher: eventPublisher,
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

	if s.eventPublisher != nil {
		s.eventPublisher.PublishAuditLog(ctx, "RECORD", "METRIC", metric.ID.String(), nil, map[string]interface{}{
			"name":     metric.Name,
			"value":    metric.Value,
			"category": metric.Category,
		})
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

func (s *AnalyticsService) GetCategoryHealthMetrics(ctx context.Context) ([]domain.CategoryHealthMetrics, error) {
	return s.kpiRepo.GetCategoryHealthMetrics(ctx)
}

func (s *AnalyticsService) StartBackgroundRefresher(ctx context.Context) {
	// Refreshes the Materialized Views once per day, as requested by user.
	ticker := time.NewTicker(24 * time.Hour)
	
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("Running daily materialized view refresh...")
				if err := s.kpiRepo.RefreshMaterializedViews(context.Background()); err != nil {
					log.Printf("Error refreshing materialized views: %v\n", err)
				}
			}
		}
	}()
}
