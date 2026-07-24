package secondary

import (
	"context"

	"backend-gmao/apps/identity-service/internal/core/domain"
	"github.com/google/uuid"
)

// UserRepository defines the secondary port for user persistence operations.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	CountByRoleID(ctx context.Context, roleID uuid.UUID) (int64, error)
	FindByTeamID(ctx context.Context, teamID uuid.UUID) ([]domain.User, error)	
	FindAll(ctx context.Context, siteIDFilter *string, offset, limit int) ([]domain.User, int64, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
