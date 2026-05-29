package service

import (
	"context"
	"errors"
	"fmt"

	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/apps/user-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleNameExists    = errors.New("a role with this name already exists")
	ErrInvalidPrivileges = errors.New("one or more privileges are invalid")
	ErrRoleHasUsers      = errors.New("cannot delete a role that is assigned to users")
)

// RoleService implements the primary.RoleUseCase primary port.
type RoleService struct {
	roleRepo secondary.RoleRepository
	userRepo secondary.UserRepository
}

// NewRoleService creates a new RoleService instance.
func NewRoleService(roleRepo secondary.RoleRepository, userRepo secondary.UserRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// CreateRole creates a new role with validated privileges.
func (s *RoleService) CreateRole(ctx context.Context, req domain.CreateRoleRequest) (*domain.RoleResponse, error) {
	// Check if name already exists
	existing, _ := s.roleRepo.FindByName(ctx, req.Name)
	if existing != nil {
		return nil, ErrRoleNameExists
	}

	// Validate all privileges are system-defined
	invalidPrivs := domain.ValidatePrivileges(req.Privileges)
	if len(invalidPrivs) > 0 {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPrivileges, invalidPrivs)
	}

	// Build the role with privileges
	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Set privileges
	if err := s.roleRepo.SetPrivileges(ctx, role.ID, req.Privileges); err != nil {
		return nil, fmt.Errorf("failed to set privileges: %w", err)
	}

	// Reload with privileges
	created, err := s.roleRepo.FindByID(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload role: %w", err)
	}

	resp := created.ToResponse()
	return &resp, nil
}

// GetRoleByID retrieves a role by its UUID.
func (s *RoleService) GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	resp := role.ToResponse()
	return &resp, nil
}

// ListRoles returns all roles.
func (s *RoleService) ListRoles(ctx context.Context) ([]domain.RoleResponse, error) {
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	responses := make([]domain.RoleResponse, 0, len(roles))
	for _, r := range roles {
		responses = append(responses, r.ToResponse())
	}

	return responses, nil
}

// UpdateRole updates an existing role's name and/or description.
func (s *RoleService) UpdateRole(ctx context.Context, id uuid.UUID, req domain.UpdateRoleRequest) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	if req.Name != nil {
		// Check uniqueness
		existing, _ := s.roleRepo.FindByName(ctx, *req.Name)
		if existing != nil && existing.ID != id {
			return nil, ErrRoleNameExists
		}
		role.Name = *req.Name
	}

	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Reload
	updated, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload role: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// DeleteRole removes a role if it has no assigned users.
func (s *RoleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	_, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return ErrRoleNotFound
	}

	// Check if any users are assigned to this role before deleting.
	count, err := s.userRepo.CountByRoleID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check role usage: %w", err)
	}
	if count > 0 {
		return ErrRoleHasUsers
	}

	if err := s.roleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// SetRolePrivileges replaces all privileges for a role.
func (s *RoleService) SetRolePrivileges(ctx context.Context, roleID uuid.UUID, req domain.SetPrivilegesRequest) (*domain.RoleResponse, error) {
	_, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	// Validate all privileges
	invalidPrivs := domain.ValidatePrivileges(req.Privileges)
	if len(invalidPrivs) > 0 {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPrivileges, invalidPrivs)
	}

	if err := s.roleRepo.SetPrivileges(ctx, roleID, req.Privileges); err != nil {
		return nil, fmt.Errorf("failed to set privileges: %w", err)
	}

	// Reload
	updated, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload role: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// ListPrivileges returns all system-defined privileges.
func (s *RoleService) ListPrivileges(ctx context.Context) ([]string, error) {
	p := domain.AllPrivileges()
	if p == nil {
		return nil, fmt.Errorf("failed to get privileges")
	}
	return p, nil
}

// PrivilegesByDomain returns all system-defined privileges grouped by domain.
func (s *RoleService) PrivilegesByDomain(ctx context.Context) (map[string][]string, error) {
	p := domain.PrivilegesByDomain()
	if p == nil {
		return nil, fmt.Errorf("failed to get privileges by domain")
	}
	return p, nil
}
