package postgres

import (
	"context"
	"fmt"

	"backend-gmao/apps/user-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleRepository is the GORM-based implementation of the RoleRepository port.
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new RoleRepository instance.
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create persists a new role to the database.
func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	result := r.db.WithContext(ctx).Create(role)
	if result.Error != nil {
		return fmt.Errorf("postgres create role: %w", result.Error)
	}
	return nil
}

// FindByID retrieves a role by UUID, preloading its privileges.
func (r *RoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	var role domain.Role
	result := r.db.WithContext(ctx).
		Preload("Privileges").
		Where("id = ?", id).
		First(&role)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find role by id: %w", result.Error)
	}
	return &role, nil
}

// FindByName retrieves a role by its name, preloading its privileges.
func (r *RoleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role
	result := r.db.WithContext(ctx).
		Preload("Privileges").
		Where("name = ?", name).
		First(&role)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find role by name: %w", result.Error)
	}
	return &role, nil
}

// FindAll retrieves all roles with their privileges.
func (r *RoleRepository) FindAll(ctx context.Context) ([]domain.Role, error) {
	var roles []domain.Role
	result := r.db.WithContext(ctx).
		Preload("Privileges").
		Order("name ASC").
		Find(&roles)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find all roles: %w", result.Error)
	}
	return roles, nil
}

// Update updates an existing role in the database.
func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	result := r.db.WithContext(ctx).Save(role)
	if result.Error != nil {
		return fmt.Errorf("postgres update role: %w", result.Error)
	}
	return nil
}

// Delete removes a role from the database by UUID.
func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete associated privileges first
		if err := tx.Where("role_id = ?", id).Delete(&domain.RolePrivilege{}).Error; err != nil {
			return fmt.Errorf("postgres delete role privileges: %w", err)
		}
		// Delete the role
		if err := tx.Where("id = ?", id).Delete(&domain.Role{}).Error; err != nil {
			return fmt.Errorf("postgres delete role: %w", err)
		}
		return nil
	})
}

// SetPrivileges replaces all privileges for a role within a transaction.
func (r *RoleRepository) SetPrivileges(ctx context.Context, roleID uuid.UUID, privileges []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing privileges
		if err := tx.Where("role_id = ?", roleID).Delete(&domain.RolePrivilege{}).Error; err != nil {
			return fmt.Errorf("postgres clear role privileges: %w", err)
		}

		// Insert new privileges
		rolePrivileges := make([]domain.RolePrivilege, 0, len(privileges))
		for _, p := range privileges {
			rolePrivileges = append(rolePrivileges, domain.RolePrivilege{
				RoleID:    roleID,
				Privilege: p,
			})
		}

		if len(rolePrivileges) > 0 {
			if err := tx.Create(&rolePrivileges).Error; err != nil {
				return fmt.Errorf("postgres set role privileges: %w", err)
			}
		}

		return nil
	})
}
