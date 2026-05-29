package secondary

import (
	"context"

	"backend-gmao/apps/user-service/internal/core/domain"
	"github.com/google/uuid"
)

// RoleRepository defines the secondary port for role persistence operations.
type RoleRepository interface {
	Create(ctx context.Context, role *domain.Role) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	FindAll(ctx context.Context) ([]domain.Role, error)
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetPrivileges(ctx context.Context, roleID uuid.UUID, privileges []string) error
}
