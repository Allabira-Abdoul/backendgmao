package postgres

import (
	"context"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type assetRepository struct {
	db *gorm.DB
}

// NewAssetRepository creates a GORM asset repository.
func NewAssetRepository(db *gorm.DB) *assetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) CreateEquipmentModel(ctx context.Context, model *domain.EquipmentModel) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *assetRepository) CreatePartModel(ctx context.Context, model *domain.PartModel) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *assetRepository) GetEquipmentModels(ctx context.Context) ([]domain.EquipmentModel, error) {
	var models []domain.EquipmentModel
	err := r.db.WithContext(ctx).Preload("Suppliers").Preload("Suppliers.Supplier").Preload("PartRequirements").Preload("PartRequirements.PartModel").Find(&models).Error
	return models, err
}

func (r *assetRepository) GetPartModels(ctx context.Context) ([]domain.PartModel, error) {
	var models []domain.PartModel
	err := r.db.WithContext(ctx).Preload("Suppliers").Preload("Suppliers.Supplier").Find(&models).Error
	return models, err
}

func (r *assetRepository) GetEquipmentModelByID(ctx context.Context, id uuid.UUID) (*domain.EquipmentModel, error) {
	var model domain.EquipmentModel
	err := r.db.WithContext(ctx).Preload("Suppliers").Preload("Suppliers.Supplier").Preload("PartRequirements").Preload("PartRequirements.PartModel").First(&model, "id = ?", id).Error
	return &model, err
}

func (r *assetRepository) GetPartModelByID(ctx context.Context, id uuid.UUID) (*domain.PartModel, error) {
	var model domain.PartModel
	err := r.db.WithContext(ctx).Preload("Suppliers").Preload("Suppliers.Supplier").First(&model, "id = ?", id).Error
	return &model, err
}

func (r *assetRepository) UpdatePartModel(ctx context.Context, model *domain.PartModel) error {
	return r.db.WithContext(ctx).Save(model).Error
}

// --- Suppliers ---

func (r *assetRepository) CreateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	return r.db.WithContext(ctx).Create(supplier).Error
}

func (r *assetRepository) GetSuppliers(ctx context.Context) ([]domain.Supplier, error) {
	var suppliers []domain.Supplier
	err := r.db.WithContext(ctx).Find(&suppliers).Error
	return suppliers, err
}

func (r *assetRepository) GetSupplierByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	var supplier domain.Supplier
	err := r.db.WithContext(ctx).First(&supplier, "id = ?", id).Error
	return &supplier, err
}

func (r *assetRepository) AddModelSupplier(ctx context.Context, modelSupplier *domain.ModelSupplier) error {
	return r.db.WithContext(ctx).Create(modelSupplier).Error
}

func (r *assetRepository) CreateEquipmentInstance(ctx context.Context, instance *domain.EquipmentInstance) error {
	return r.db.WithContext(ctx).Create(instance).Error
}

func (r *assetRepository) CreatePartInstance(ctx context.Context, instance *domain.PartInstance) error {
	return r.db.WithContext(ctx).Create(instance).Error
}

func (r *assetRepository) GetEquipmentInstances(ctx context.Context) ([]domain.EquipmentInstance, error) {
	var instances []domain.EquipmentInstance
	err := r.db.WithContext(ctx).Preload("EquipmentModel").Preload("EquipmentModel.Suppliers").Preload("EquipmentModel.Suppliers.Supplier").Preload("Parts").Preload("Parts.PartModel").Preload("Parts.PartModel.Suppliers").Preload("Parts.PartModel.Suppliers.Supplier").Preload("Supplier").Find(&instances).Error
	return instances, err
}

func (r *assetRepository) GetEquipmentInstanceByCode(ctx context.Context, code string) (*domain.EquipmentInstance, error) {
	var instance domain.EquipmentInstance
	err := r.db.WithContext(ctx).Preload("EquipmentModel").Preload("EquipmentModel.Suppliers").Preload("EquipmentModel.Suppliers.Supplier").Preload("Supplier").Preload("Parts").Preload("Parts.PartModel").Preload("Parts.PartModel.Suppliers").Preload("Parts.PartModel.Suppliers.Supplier").Preload("Parts.Supplier").First(&instance, "code = ?", code).Error
	return &instance, err
}

func (r *assetRepository) GetEquipmentInstanceByID(ctx context.Context, id uuid.UUID) (*domain.EquipmentInstance, error) {
	var instance domain.EquipmentInstance
	err := r.db.WithContext(ctx).Preload("EquipmentModel").Preload("EquipmentModel.Suppliers").Preload("EquipmentModel.Suppliers.Supplier").Preload("Supplier").Preload("Parts").Preload("Parts.PartModel").Preload("Parts.PartModel.Suppliers").Preload("Parts.PartModel.Suppliers.Supplier").Preload("Parts.Supplier").First(&instance, "id = ?", id).Error
	return &instance, err
}

func (r *assetRepository) UpdateEquipmentInstance(ctx context.Context, instance *domain.EquipmentInstance) error {
	return r.db.WithContext(ctx).Save(instance).Error
}

func (r *assetRepository) GetPartInstanceByID(ctx context.Context, id uuid.UUID) (*domain.PartInstance, error) {
	var instance domain.PartInstance
	err := r.db.WithContext(ctx).Preload("PartModel").Preload("PartModel.Suppliers").Preload("PartModel.Suppliers.Supplier").Preload("Supplier").First(&instance, "id = ?", id).Error
	return &instance, err
}

func (r *assetRepository) UpdatePartInstance(ctx context.Context, instance *domain.PartInstance) error {
	return r.db.WithContext(ctx).Save(instance).Error
}

func (r *assetRepository) CreatePartConsumptionLog(ctx context.Context, log *domain.PartConsumptionLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}



func (r *assetRepository) UpdateEquipmentModel(ctx context.Context, model *domain.EquipmentModel) error {
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *assetRepository) CreateEquipmentModelPartRequirement(ctx context.Context, req *domain.EquipmentModelPartRequirement) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *assetRepository) DeleteEquipmentModelPartRequirements(ctx context.Context, modelID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("equipment_model_id = ?", modelID).Delete(&domain.EquipmentModelPartRequirement{}).Error
}

func (r *assetRepository) UpdateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	return r.db.WithContext(ctx).Save(supplier).Error
}

func (r *assetRepository) DeleteSupplier(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Supplier{}, "id = ?", id).Error
}




