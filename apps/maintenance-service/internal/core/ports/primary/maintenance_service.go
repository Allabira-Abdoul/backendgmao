package primary

import (
	"context"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"github.com/google/uuid"
)

// MaintenanceService defines primary application business operations.
type MaintenanceService interface {
	CreateWorkOrder(ctx context.Context, req domain.CreateOrdreTravailRequest) (*domain.OrdreTravailResponse, error)
	UpdateWorkOrder(ctx context.Context, id uuid.UUID, req domain.UpdateOrdreTravailRequest) (*domain.OrdreTravailResponse, error)
	DeleteWorkOrder(ctx context.Context, id uuid.UUID) error
	GetWorkOrder(ctx context.Context, id uuid.UUID) (*domain.OrdreTravailResponse, error)
	GetAllWorkOrders(ctx context.Context) ([]domain.OrdreTravailResponse, error)
	StartWorkOrder(ctx context.Context, id uuid.UUID) (*domain.OrdreTravailResponse, error)

	RecordIntervention(ctx context.Context, workOrderID uuid.UUID, req domain.CreateInterventionRequest) (*domain.InterventionResponse, error)
	UpdateIntervention(ctx context.Context, workOrderID uuid.UUID, interventionID uuid.UUID, req domain.UpdateInterventionRequest) (*domain.InterventionResponse, error)
	StartIntervention(ctx context.Context, workOrderID uuid.UUID, interventionID uuid.UUID) (*domain.InterventionResponse, error)
	EndIntervention(ctx context.Context, workOrderID uuid.UUID, interventionID uuid.UUID) (*domain.InterventionResponse, error)
	GetInterventionsForWorkOrder(ctx context.Context, workOrderID uuid.UUID) ([]domain.InterventionResponse, error)

	CreateInspection(ctx context.Context, workOrderID uuid.UUID, req domain.CreateInspectionRequest) (*domain.InspectionResponse, error)
	UpdateInspection(ctx context.Context, workOrderID uuid.UUID, inspectionID uuid.UUID, req domain.UpdateInspectionRequest) (*domain.InspectionResponse, error)
	StartInspection(ctx context.Context, workOrderID uuid.UUID, inspectionID uuid.UUID) (*domain.InspectionResponse, error)
	EndInspection(ctx context.Context, workOrderID uuid.UUID, inspectionID uuid.UUID) (*domain.InspectionResponse, error)
}
