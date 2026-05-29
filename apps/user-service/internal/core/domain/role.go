package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a role in the GMAO system (e.g., Administrator, Technician, Manager).
type Role struct {
	ID          uuid.UUID       `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string          `gorm:"column:name;uniqueIndex;not null" json:"name"`
	Description string          `gorm:"column:description" json:"description"`
	Privileges  []RolePrivilege `gorm:"foreignKey:RoleID;references:ID" json:"privileges,omitempty"`
	CreatedAt   time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (Role) TableName() string {
	return "roles"
}

// RolePrivilege represents the many-to-many relationship between roles and privileges.
type RolePrivilege struct {
	RoleID    uuid.UUID `gorm:"column:role_id;type:uuid;primaryKey" json:"role_id"`
	Privilege string    `gorm:"column:privilege;primaryKey" json:"privilege"`
}

// TableName overrides the default table name.
func (RolePrivilege) TableName() string {
	return "role_privileges"
}

// RoleResponse is the DTO returned by API endpoints.
type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Privileges  []string  `json:"privileges"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a Role to a RoleResponse.
func (r *Role) ToResponse() RoleResponse {
	privileges := make([]string, 0, len(r.Privileges))
	for _, rp := range r.Privileges {
		privileges = append(privileges, rp.Privilege)
	}

	return RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Privileges:  privileges,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// GetPrivilegeStrings extracts the privilege names from the role.
func (r *Role) GetPrivilegeStrings() []string {
	privileges := make([]string, 0, len(r.Privileges))
	for _, rp := range r.Privileges {
		privileges = append(privileges, rp.Privilege)
	}
	return privileges
}

// CreateRoleRequest is the DTO for creating a new role.
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=100"`
	Description string   `json:"description" binding:"max=500"`
	Privileges  []string `json:"privileges" binding:"required,min=1"`
}

// UpdateRoleRequest is the DTO for updating an existing role.
type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=500"`
	Privileges  *[]string `json:"privileges,omitempty" binding:"omitempty,min=1"`
}

// SetPrivilegesRequest is the DTO for setting a role's privileges.
type SetPrivilegesRequest struct {
	Privileges []string `json:"privileges" binding:"required,min=1"`
}
