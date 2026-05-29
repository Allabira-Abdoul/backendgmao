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
	err := r.db.WithContext(ctx).Preload("Thresholds").Find(&models).Error
	return models, err
}

func (r *assetRepository) GetPartModels(ctx context.Context) ([]domain.PartModel, error) {
	var models []domain.PartModel
	err := r.db.WithContext(ctx).Preload("Thresholds").Find(&models).Error
	return models, err
}

func (r *assetRepository) GetEquipmentModelByID(ctx context.Context, id uuid.UUID) (*domain.EquipmentModel, error) {
	var model domain.EquipmentModel
	err := r.db.WithContext(ctx).Preload("Thresholds").First(&model, "id = ?", id).Error
	return &model, err
}

func (r *assetRepository) GetPartModelByID(ctx context.Context, id uuid.UUID) (*domain.PartModel, error) {
	var model domain.PartModel
	err := r.db.WithContext(ctx).Preload("Thresholds").First(&model, "id = ?", id).Error
	return &model, err
}

func (r *assetRepository) CreateEquipmentInstance(ctx context.Context, instance *domain.EquipmentInstance) error {
	return r.db.WithContext(ctx).Create(instance).Error
}

func (r *assetRepository) CreatePartInstance(ctx context.Context, instance *domain.PartInstance) error {
	return r.db.WithContext(ctx).Create(instance).Error
}

func (r *assetRepository) GetEquipmentInstances(ctx context.Context) ([]domain.EquipmentInstance, error) {
	var instances []domain.EquipmentInstance
	err := r.db.WithContext(ctx).Preload("EquipmentModel").Preload("Parts").Preload("Parts.PartModel").Preload("Thresholds").Find(&instances).Error
	return instances, err
}

func (r *assetRepository) GetEquipmentInstanceByCode(ctx context.Context, code string) (*domain.EquipmentInstance, error) {
	var instance domain.EquipmentInstance
	err := r.db.WithContext(ctx).Preload("EquipmentModel").Preload("Parts").Preload("Parts.PartModel").Preload("Thresholds").First(&instance, "code = ?", code).Error
	return &instance, err
}

func (r *assetRepository) GetEquipmentInstanceByID(ctx context.Context, id uuid.UUID) (*domain.EquipmentInstance, error) {
	var instance domain.EquipmentInstance
	err := r.db.WithContext(ctx).Preload("EquipmentModel").Preload("Parts").Preload("Parts.PartModel").Preload("Thresholds").First(&instance, "id = ?", id).Error
	return &instance, err
}
