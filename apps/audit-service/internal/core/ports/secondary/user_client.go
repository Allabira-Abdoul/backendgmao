package secondary

import (
	"context"

	"github.com/google/uuid"
)

type UserClient interface {
	GetUserNameByID(ctx context.Context, id uuid.UUID) (string, error)
	GetUserNamesByID(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error)
}	