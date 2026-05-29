package secondary

import (
	"context"
	"time"
	"github.com/google/uuid"
)

// AssetInfo contains the subset of asset data needed for analytics.
type AssetInfo struct {
	ID           uuid.UUID
	Category     string
	PurchaseDate time.Time
}

// AssetClient defines the interface to fetch asset details from the asset-service.
type AssetClient interface {
	GetAssetInfo(ctx context.Context, assetID uuid.UUID) (*AssetInfo, error)
}
