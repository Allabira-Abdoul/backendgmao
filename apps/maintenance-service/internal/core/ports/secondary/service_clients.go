package secondary

import (
	"context"

	"github.com/google/uuid"
)

type UserClient interface {
	GetUserName(ctx context.Context, id uuid.UUID) (string, error)
	GetUserNames(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error)
}

type AssetClient interface {
	GetAssetName(ctx context.Context, id uuid.UUID) (string, error)
	GetAssetNames(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error)
	UpdateAssetStatus(ctx context.Context, id uuid.UUID, status string) error
}
