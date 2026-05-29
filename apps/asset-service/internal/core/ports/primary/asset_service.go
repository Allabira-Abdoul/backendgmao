package primary

import (
	"context"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
)

type AssetService interface {
	// Models
	CreateEquipmentModel(ctx context.Context, req domain.CreateEquipmentModelRequest) (domain.EquipmentModelResponse, error)
	CreatePartModel(ctx context.Context, req domain.CreatePartModelRequest) (domain.PartModelResponse, error)
	GetEquipmentModels(ctx context.Context) ([]domain.EquipmentModelResponse, error)
	GetPartModels(ctx context.Context) ([]domain.PartModelResponse, error)

	// Instances
	CreateEquipmentInstance(ctx context.Context, req domain.CreateEquipmentInstanceRequest) (domain.EquipmentInstanceResponse, error)
	CreatePartInstance(ctx context.Context, req domain.CreatePartInstanceRequest) (domain.PartInstanceResponse, error)
	GetEquipmentInstances(ctx context.Context) ([]domain.EquipmentInstanceResponse, error)
	GetEquipmentInstanceByCode(ctx context.Context, code string) (domain.EquipmentInstanceResponse, error)
	GetEquipmentInstanceByID(ctx context.Context, id uuid.UUID) (domain.EquipmentInstanceResponse, error)
}
