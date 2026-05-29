package secondary

import (
	"context"
	"backend-gmao/apps/auth-service/internal/core/domain"
	"github.com/google/uuid"
)

// UserProvider defines the port for fetching user identity and privileges from the user-service.
type UserProvider interface {
	FetchUserForAuth(ctx context.Context, email string) (*domain.User, error)
	FetchUserByIDForAuth(ctx context.Context, id uuid.UUID) (*domain.User, error)
}
