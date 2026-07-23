package secondary

import (
	"context"

	"github.com/google/uuid"
	"backend-gmao/apps/asset-service/internal/core/domain"
)

type AssetRepository interface {
	// Hierarchy Operations
	CreateSite(ctx context.Context, site *domain.Site) error
	GetSiteHierarchy(ctx context.Context, siteID uuid.UUID) (*domain.Site, error)
	GetAllSites(ctx context.Context) ([]domain.Site, error)

	CreateSystem(ctx context.Context, system *domain.System) error
	CreateAsset(ctx context.Context, asset *domain.Asset) error
	UpdateAssetStatus(ctx context.Context, assetID uuid.UUID, status string) error
	CreateSubsystem(ctx context.Context, subsystem *domain.Subsystem) error

	// Inventory Operations
	CreateInventoryItem(ctx context.Context, item *domain.InventoryItem) error
	GetInventoryItem(ctx context.Context, itemID uuid.UUID) (*domain.InventoryItem, error)
	UpdateInventoryStock(ctx context.Context, itemID uuid.UUID, delta int) error

	// Component Operations
	CreateComponent(ctx context.Context, component *domain.Component) error
	GetComponent(ctx context.Context, componentID uuid.UUID) (*domain.Component, error)
}
