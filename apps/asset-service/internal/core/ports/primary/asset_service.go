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
	UpdateEquipmentModel(ctx context.Context, id uuid.UUID, req domain.UpdateEquipmentModelRequest) (domain.EquipmentModelResponse, error)
	UpdatePartModel(ctx context.Context, id uuid.UUID, req domain.UpdatePartModelRequest) (domain.PartModelResponse, error)

	// Suppliers
	CreateSupplier(ctx context.Context, req domain.CreateSupplierRequest) (domain.SupplierResponse, error)
	GetSuppliers(ctx context.Context) ([]domain.SupplierResponse, error)
	UpdateSupplier(ctx context.Context, id uuid.UUID, req domain.UpdateSupplierRequest) (domain.SupplierResponse, error)
	DeleteSupplier(ctx context.Context, id uuid.UUID) error
	AddSupplierToEquipmentModel(ctx context.Context, modelID uuid.UUID, req domain.AddModelSupplierRequest) (domain.ModelSupplierResponse, error)
	AddSupplierToPartModel(ctx context.Context, modelID uuid.UUID, req domain.AddModelSupplierRequest) (domain.ModelSupplierResponse, error)

	// Instances
	CreateEquipmentInstance(ctx context.Context, req domain.CreateEquipmentInstanceRequest) (domain.EquipmentInstanceResponse, error)
	CreatePartInstance(ctx context.Context, req domain.CreatePartInstanceRequest) (domain.PartInstanceResponse, error)
	GetEquipmentInstances(ctx context.Context) ([]domain.EquipmentInstanceResponse, error)
	GetEquipmentInstanceByCode(ctx context.Context, code string) (domain.EquipmentInstanceResponse, error)
	GetEquipmentInstanceByID(ctx context.Context, id uuid.UUID) (domain.EquipmentInstanceResponse, error)
	UpdateEquipmentStatus(ctx context.Context, id uuid.UUID, newStatus string) error
	UpdateEquipmentLocation(ctx context.Context, id uuid.UUID, newLocation string) error

	MovePartInstance(ctx context.Context, partInstanceID uuid.UUID, req domain.MovePartInstanceRequest) (domain.PartInstanceResponse, error)
	ConsumePart(ctx context.Context, req domain.ConsumePartRequest, userID uuid.UUID) error


}
