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
