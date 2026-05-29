package secondary

import (
	"context"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
)

type AssetRepository interface {
	// Models
	CreateEquipmentModel(ctx context.Context, model *domain.EquipmentModel) error
	CreatePartModel(ctx context.Context, model *domain.PartModel) error
	GetEquipmentModels(ctx context.Context) ([]domain.EquipmentModel, error)
	GetPartModels(ctx context.Context) ([]domain.PartModel, error)
	GetEquipmentModelByID(ctx context.Context, id uuid.UUID) (*domain.EquipmentModel, error)
	GetPartModelByID(ctx context.Context, id uuid.UUID) (*domain.PartModel, error)
	UpdatePartModel(ctx context.Context, model *domain.PartModel) error

	// Instances
	CreateEquipmentInstance(ctx context.Context, instance *domain.EquipmentInstance) error
	CreatePartInstance(ctx context.Context, instance *domain.PartInstance) error
	GetEquipmentInstances(ctx context.Context) ([]domain.EquipmentInstance, error)
	GetEquipmentInstanceByCode(ctx context.Context, code string) (*domain.EquipmentInstance, error)
	GetEquipmentInstanceByID(ctx context.Context, id uuid.UUID) (*domain.EquipmentInstance, error)
	GetPartInstanceByID(ctx context.Context, id uuid.UUID) (*domain.PartInstance, error)
	UpdatePartInstance(ctx context.Context, instance *domain.PartInstance) error

	// Logs
	CreatePartConsumptionLog(ctx context.Context, log *domain.PartConsumptionLog) error

	// Thresholds
	GetMetricThresholds(ctx context.Context, metricName string, eqID *uuid.UUID, partID *uuid.UUID) ([]domain.MetricThreshold, error)
}
