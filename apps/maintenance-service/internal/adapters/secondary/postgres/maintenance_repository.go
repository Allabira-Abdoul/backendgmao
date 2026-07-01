package postgres

import (
	"context"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type maintenanceRepository struct {
	db *gorm.DB
}

// NewMaintenanceRepository creates a GORM maintenance repository.
func NewMaintenanceRepository(db *gorm.DB) *maintenanceRepository {
	return &maintenanceRepository{db: db}
}

func (r *maintenanceRepository) CreateWorkOrder(ctx context.Context, wo *domain.OrdreTravail) error {
	return r.db.WithContext(ctx).Create(wo).Error
}

func (r *maintenanceRepository) UpdateWorkOrder(ctx context.Context, wo *domain.OrdreTravail) error {
	return r.db.WithContext(ctx).Save(wo).Error
}

func (r *maintenanceRepository) DeleteWorkOrder(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.OrdreTravail{}, "id = ?", id).Error
}

func (r *maintenanceRepository) FindWorkOrderByID(ctx context.Context, id uuid.UUID) (*domain.OrdreTravail, error) {
	var wo domain.OrdreTravail
	if err := r.db.WithContext(ctx).Preload("Interventions").Preload("Inspections").First(&wo, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wo, nil
}

func (r *maintenanceRepository) FindAllWorkOrders(ctx context.Context) ([]domain.OrdreTravail, error) {
	var wos []domain.OrdreTravail
	if err := r.db.WithContext(ctx).Preload("Interventions").Preload("Inspections").Find(&wos).Error; err != nil {
		return nil, err
	}
	return wos, nil
}

func (r *maintenanceRepository) CreateIntervention(ctx context.Context, intervention *domain.Intervention) error {
	return r.db.WithContext(ctx).Create(intervention).Error
}

func (r *maintenanceRepository) FindInterventionsByWorkOrderID(ctx context.Context, workOrderID uuid.UUID) ([]domain.Intervention, error) {
	var interventions []domain.Intervention
	if err := r.db.WithContext(ctx).Preload("Measurements").Where("work_order_id = ?", workOrderID).Find(&interventions).Error; err != nil {
		return nil, err
	}
	return interventions, nil
}

func (r *maintenanceRepository) FindInterventionByID(ctx context.Context, id uuid.UUID) (*domain.Intervention, error) {
	var inv domain.Intervention
	if err := r.db.WithContext(ctx).First(&inv, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *maintenanceRepository) UpdateIntervention(ctx context.Context, intervention *domain.Intervention) error {
	return r.db.WithContext(ctx).Save(intervention).Error
}

func (r *maintenanceRepository) CreateInspection(ctx context.Context, inspection *domain.Inspection) error {
	return r.db.WithContext(ctx).Create(inspection).Error
}

func (r *maintenanceRepository) FindInspectionsByWorkOrderID(ctx context.Context, workOrderID uuid.UUID) ([]domain.Inspection, error) {
	var inspections []domain.Inspection
	if err := r.db.WithContext(ctx).Preload("Measurements").Where("work_order_id = ?", workOrderID).Find(&inspections).Error; err != nil {
		return nil, err
	}
	return inspections, nil
}

func (r *maintenanceRepository) FindInspectionByID(ctx context.Context, id uuid.UUID) (*domain.Inspection, error) {
	var ins domain.Inspection
	if err := r.db.WithContext(ctx).First(&ins, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ins, nil
}

func (r *maintenanceRepository) UpdateInspection(ctx context.Context, inspection *domain.Inspection) error {
	return r.db.WithContext(ctx).Save(inspection).Error
}

func (r *maintenanceRepository) CreateMaintenanceSchedule(ctx context.Context, schedule *domain.MaintenanceSchedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

func (r *maintenanceRepository) FindMaintenanceSchedulesByAssetID(ctx context.Context, assetID uuid.UUID) ([]domain.MaintenanceSchedule, error) {
	var schedules []domain.MaintenanceSchedule
	if err := r.db.WithContext(ctx).Where("asset_id = ?", assetID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *maintenanceRepository) FindAllMaintenanceSchedules(ctx context.Context) ([]domain.MaintenanceSchedule, error) {
	var schedules []domain.MaintenanceSchedule
	if err := r.db.WithContext(ctx).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *maintenanceRepository) FindMaintenanceScheduleByID(ctx context.Context, id uuid.UUID) (*domain.MaintenanceSchedule, error) {
	var schedule domain.MaintenanceSchedule
	if err := r.db.WithContext(ctx).First(&schedule, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *maintenanceRepository) UpdateMaintenanceSchedule(ctx context.Context, schedule *domain.MaintenanceSchedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

func (r *maintenanceRepository) CreateCounterReading(ctx context.Context, reading *domain.CounterReading) error {
	return r.db.WithContext(ctx).Create(reading).Error
}

func (r *maintenanceRepository) FindCounterReadingsByAssetID(ctx context.Context, assetID uuid.UUID) ([]domain.CounterReading, error) {
	var readings []domain.CounterReading
	if err := r.db.WithContext(ctx).Where("asset_id = ?", assetID).Find(&readings).Error; err != nil {
		return nil, err
	}
	return readings, nil
}

// ----------------------------------------------------------------------------
// Defect Alerts
// ----------------------------------------------------------------------------

func (r *maintenanceRepository) CreateDefectAlert(ctx context.Context, alert *domain.DefectAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

func (r *maintenanceRepository) FindAllDefectAlerts(ctx context.Context) ([]domain.DefectAlert, error) {
	var alerts []domain.DefectAlert
	if err := r.db.WithContext(ctx).Find(&alerts).Error; err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *maintenanceRepository) FindDefectAlertByID(ctx context.Context, id uuid.UUID) (*domain.DefectAlert, error) {
	var alert domain.DefectAlert
	if err := r.db.WithContext(ctx).First(&alert, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *maintenanceRepository) UpdateDefectAlert(ctx context.Context, alert *domain.DefectAlert) error {
	return r.db.WithContext(ctx).Save(alert).Error
}
