package secondary

import (
	"context"

	"github.com/google/uuid"
)

type AssetClient interface {
	GetAssetName(ctx context.Context, id uuid.UUID) (string, error)
}
