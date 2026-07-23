package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/apps/asset-service/internal/core/ports/secondary"
)

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) secondary.AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) CreateSite(ctx context.Context, site *domain.Site) error {
	return r.db.WithContext(ctx).Create(site).Error
}

func (r *assetRepository) GetSiteHierarchy(ctx context.Context, siteID uuid.UUID) (*domain.Site, error) {
	var site domain.Site
	err := r.db.WithContext(ctx).
		Preload("Systems.Assets.Subsystems.Components.InventoryItem").
		Where("id = ?", siteID).
		First(&site).Error
	if err != nil {
		return nil, err
	}
	return &site, nil
}

func (r *assetRepository) GetAllSites(ctx context.Context) ([]domain.Site, error) {
	var sites []domain.Site
	err := r.db.WithContext(ctx).Find(&sites).Error
	return sites, err
}

func (r *assetRepository) CreateSystem(ctx context.Context, system *domain.System) error {
	return r.db.WithContext(ctx).Create(system).Error
}

func (r *assetRepository) CreateAsset(ctx context.Context, asset *domain.Asset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *assetRepository) UpdateAssetStatus(ctx context.Context, assetID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Asset{}).
		Where("id = ?", assetID).
		Update("status", status).Error
}

func (r *assetRepository) CreateSubsystem(ctx context.Context, subsystem *domain.Subsystem) error {
	return r.db.WithContext(ctx).Create(subsystem).Error
}

func (r *assetRepository) CreateInventoryItem(ctx context.Context, item *domain.InventoryItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *assetRepository) GetInventoryItem(ctx context.Context, itemID uuid.UUID) (*domain.InventoryItem, error) {
	var item domain.InventoryItem
	err := r.db.WithContext(ctx).Where("id = ?", itemID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *assetRepository) UpdateInventoryStock(ctx context.Context, itemID uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).
		Model(&domain.InventoryItem{}).
		Where("id = ?", itemID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", delta)).Error
}

func (r *assetRepository) CreateComponent(ctx context.Context, component *domain.Component) error {
	return r.db.WithContext(ctx).Create(component).Error
}

func (r *assetRepository) GetComponent(ctx context.Context, componentID uuid.UUID) (*domain.Component, error) {
	var component domain.Component
	err := r.db.WithContext(ctx).
		Preload("InventoryItem").
		Where("id = ?", componentID).
		First(&component).Error
	if err != nil {
		return nil, err
	}
	return &component, nil
}
