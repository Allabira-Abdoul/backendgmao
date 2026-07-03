package secondary

import (
	"context"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"github.com/google/uuid"
)

// MaintenanceRepository defines secondary adapter database actions.
type MaintenanceRepository interface {
	CreateWorkOrder(ctx context.Context, wo *domain.OrdreTravail) error
	UpdateWorkOrder(ctx context.Context, wo *domain.OrdreTravail) error
	DeleteWorkOrder(ctx context.Context, id uuid.UUID) error
	FindWorkOrderByID(ctx context.Context, id uuid.UUID) (*domain.OrdreTravail, error)
	FindAllWorkOrders(ctx context.Context) ([]domain.OrdreTravail, error)

	CreateIntervention(ctx context.Context, intervention *domain.Intervention) error
	FindInterventionsByWorkOrderID(ctx context.Context, workOrderID uuid.UUID) ([]domain.Intervention, error)
	FindInterventionByID(ctx context.Context, id uuid.UUID) (*domain.Intervention, error)
	UpdateIntervention(ctx context.Context, intervention *domain.Intervention) error

	CreateInspection(ctx context.Context, inspection *domain.Inspection) error
	FindInspectionsByAssetID(ctx context.Context, assetID uuid.UUID) ([]domain.Inspection, error)
	FindAllInspections(ctx context.Context) ([]domain.Inspection, error)
	FindInspectionByID(ctx context.Context, id uuid.UUID) (*domain.Inspection, error)
	UpdateInspection(ctx context.Context, inspection *domain.Inspection) error

	CreateMaintenanceSchedule(ctx context.Context, schedule *domain.MaintenanceSchedule) error
	FindMaintenanceSchedulesByAssetID(ctx context.Context, assetID uuid.UUID) ([]domain.MaintenanceSchedule, error)
	FindAllMaintenanceSchedules(ctx context.Context) ([]domain.MaintenanceSchedule, error)
	FindMaintenanceScheduleByID(ctx context.Context, id uuid.UUID) (*domain.MaintenanceSchedule, error)
	UpdateMaintenanceSchedule(ctx context.Context, schedule *domain.MaintenanceSchedule) error

	CreateCounterReading(ctx context.Context, reading *domain.CounterReading) error
	FindCounterReadingsByAssetID(ctx context.Context, assetID uuid.UUID) ([]domain.CounterReading, error)

	// Defect Alerts
	CreateDefectAlert(ctx context.Context, alert *domain.DefectAlert) error
	FindAllDefectAlerts(ctx context.Context) ([]domain.DefectAlert, error)
	FindDefectAlertByID(ctx context.Context, id uuid.UUID) (*domain.DefectAlert, error)
	UpdateDefectAlert(ctx context.Context, alert *domain.DefectAlert) error
}
