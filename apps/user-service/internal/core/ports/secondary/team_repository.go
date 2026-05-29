package secondary

import (
	"context"

	"backend-gmao/apps/user-service/internal/core/domain"
	"github.com/google/uuid"
)

// TeamRepository defines the interface for data access of teams.
type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Team, error)
	FindByName(ctx context.Context, name string) (*domain.Team, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Team, error)
	FindAll(ctx context.Context, limit, offset int) ([]domain.Team, int64, error)
	Update(ctx context.Context, team *domain.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
}
