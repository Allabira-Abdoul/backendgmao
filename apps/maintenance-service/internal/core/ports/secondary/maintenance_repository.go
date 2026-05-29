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
	FindInspectionsByWorkOrderID(ctx context.Context, workOrderID uuid.UUID) ([]domain.Inspection, error)
	FindInspectionByID(ctx context.Context, id uuid.UUID) (*domain.Inspection, error)
	UpdateInspection(ctx context.Context, inspection *domain.Inspection) error
}
