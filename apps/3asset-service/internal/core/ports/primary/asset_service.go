package primary

import (
	"context"

	"github.com/google/uuid"
	"backend-gmao/apps/asset-service/internal/core/domain"
)

type AssetService interface {
	// Site
	CreateSite(ctx context.Context, req domain.CreateSiteRequest) (*domain.SiteResponse, error)
	GetSiteHierarchy(ctx context.Context, siteID uuid.UUID) (*domain.SiteResponse, error)
	GetAllSites(ctx context.Context) ([]domain.SiteResponse, error)

	// System & Asset & Subsystem
	CreateSystem(ctx context.Context, req domain.CreateSystemRequest) (*domain.SystemResponse, error)
	CreateAsset(ctx context.Context, req domain.CreateAssetRequest) (*domain.AssetResponse, error)
	UpdateAssetStatus(ctx context.Context, assetID uuid.UUID, status string) error
	CreateSubsystem(ctx context.Context, req domain.CreateSubsystemRequest) (*domain.SubsystemResponse, error)

	// Inventory & Components
	CreateInventoryItem(ctx context.Context, req domain.CreateInventoryItemRequest) (*domain.InventoryItemResponse, error)
	CreateComponent(ctx context.Context, req domain.CreateComponentRequest) (*domain.ComponentResponse, error)
}
