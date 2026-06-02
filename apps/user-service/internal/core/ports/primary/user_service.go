package primary

import (
	"context"

	"backend-gmao/apps/user-service/internal/core/domain"
	"github.com/google/uuid"
)

// UserUseCase defines the primary port for user-related use cases.
type UserUseCase interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error)
	GetUserByIDInternal(ctx context.Context, id uuid.UUID) (*domain.InternalUserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.InternalUserResponse, error)
	ListUsers(ctx context.Context, limit, offset int) ([]domain.UserResponse, int64, error)
	GetCompactUsers(ctx context.Context) ([]domain.CompactUserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SuspendUser(ctx context.Context, id uuid.UUID) error
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
	AdminResetPassword(ctx context.Context, id uuid.UUID) (string, error)
	ChangePassword(ctx context.Context, id uuid.UUID, newPassword string) error
}

// RoleUseCase defines the primary port for role-related use cases.
type RoleUseCase interface {
	CreateRole(ctx context.Context, req domain.CreateRoleRequest) (*domain.RoleResponse, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.RoleResponse, error)
	ListRoles(ctx context.Context) ([]domain.RoleResponse, error)
	GetCompactRoles(ctx context.Context) ([]domain.CompactRoleResponse, error)
	UpdateRole(ctx context.Context, id uuid.UUID, req domain.UpdateRoleRequest) (*domain.RoleResponse, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error
	SetRolePrivileges(ctx context.Context, roleID uuid.UUID, req domain.SetPrivilegesRequest) (*domain.RoleResponse, error)
	ListPrivileges(ctx context.Context) ([]string, error)
	PrivilegesByDomain(ctx context.Context) (map[string][]string, error)
}

// TeamUseCase defines the primary port for team-related use cases.
type TeamUseCase interface {
	CreateTeam(ctx context.Context, req domain.CreateTeamRequest) (*domain.TeamResponse, error)
	GetTeamByID(ctx context.Context, id uuid.UUID) (*domain.TeamResponse, error)
	ListTeams(ctx context.Context, limit, offset int) ([]domain.TeamResponse, int64, error)
	GetCompactTeams(ctx context.Context) ([]domain.CompactTeamResponse, error)
	UpdateTeam(ctx context.Context, id uuid.UUID, req domain.UpdateTeamRequest) (*domain.TeamResponse, error)
	DeleteTeam(ctx context.Context, id uuid.UUID) error
}
