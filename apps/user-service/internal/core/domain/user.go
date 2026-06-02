package domain

import (
	"time"

	"github.com/google/uuid"
)

// AccountStatus represents the status of a user account.
type AccountStatus string

const (
	StatusActive   AccountStatus = "ACTIVE"
	StatusInactive AccountStatus = "INACTIVE"
	StatusLocked   AccountStatus = "LOCKED"
)

// User represents the User entity in the GMAO system.
type User struct {
	ID            uuid.UUID     `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FullName      string        `gorm:"column:full_name;not null" json:"full_name"`
	Email         string        `gorm:"column:email;uniqueIndex;not null" json:"email"`
	Password           string        `gorm:"column:password;not null" json:"-"`
	MustChangePassword bool          `gorm:"column:must_change_password;default:false" json:"must_change_password"`
	Status             AccountStatus `gorm:"column:status;type:varchar(20);default:'ACTIVE'" json:"status"`
	RoleID        uuid.UUID     `gorm:"column:role_id;type:uuid;not null" json:"role_id"`
	Role          Role          `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"role,omitempty"`
	TeamID        *uuid.UUID    `gorm:"column:team_id;type:uuid" json:"team_id"`
	Team          *Team         `gorm:"foreignKey:TeamID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"team,omitempty"`
	CreatedAt     time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "users"
}

// Suspend sets the user's status to INACTIVE.
func (u *User) Suspend() {
	u.Status = StatusInactive
}

// Activate sets the user's status to ACTIVE.
func (u *User) Activate() {
	u.Status = StatusActive
}

// UserResponse is the DTO returned by API endpoints (excludes password).
type UserResponse struct {
	ID            uuid.UUID     `json:"id"`
	FullName      string        `json:"full_name"`
	Email              string        `json:"email"`
	Status             AccountStatus `json:"status"`
	MustChangePassword bool          `json:"must_change_password"`
	Role          *RoleResponse `json:"role,omitempty"`
	Team          *TeamResponse `json:"team,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// ToResponse converts a User to a UserResponse (safe for API output).
func (u *User) ToResponse() UserResponse {
	resp := UserResponse{
		ID:            u.ID,
		FullName:      u.FullName,
		Email:              u.Email,
		Status:             u.Status,
		MustChangePassword: u.MustChangePassword,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}

	if u.Role.ID != uuid.Nil {
		roleResp := u.Role.ToResponse()
		resp.Role = &roleResp
	}

	if u.Team != nil && u.Team.ID != uuid.Nil {
		teamResp := u.Team.ToResponse()
		resp.Team = &teamResp
	}

	return resp
}

// CompactUserResponse is a lightweight DTO for dropdowns.
type CompactUserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ToCompactResponse converts a User to a CompactUserResponse.
func (u *User) ToCompactResponse() CompactUserResponse {
	return CompactUserResponse{
		ID:   u.ID,
		Name: u.FullName,
	}
}

// InternalUserResponse is the DTO used for inter-service communication.
// It includes the hashed password for authentication verification.
type InternalUserResponse struct {
	ID            uuid.UUID     `json:"id"`
	FullName      string        `json:"full_name"`
	Email         string        `json:"email"`
	Password           string        `json:"password"`
	Status             AccountStatus `json:"status"`
	MustChangePassword bool          `json:"must_change_password"`
	RoleName      string        `json:"role_name"`
	Privileges    []string      `json:"privileges"`
}

// ToInternalResponse converts a User to an InternalUserResponse.
func (u *User) ToInternalResponse() InternalUserResponse {
	return InternalUserResponse{
		ID:                 u.ID,
		FullName:           u.FullName,
		Email:              u.Email,
		Password:           u.Password,
		Status:             u.Status,
		MustChangePassword: u.MustChangePassword,
		RoleName:           u.Role.Name,
		Privileges:         u.Role.GetPrivilegeStrings(),
	}
}

// CreateUserRequest is the DTO for creating a new user.
type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	RoleID   string `json:"role_id" binding:"required,uuid"`
}

// UpdateUserRequest is the DTO for updating an existing user.
type UpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty" binding:"omitempty,min=2,max=255"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Status   *string `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE INACTIVE LOCKED"`
	RoleID   *string `json:"role_id,omitempty" binding:"omitempty,uuid"`
	TeamID   *string `json:"team_id,omitempty" binding:"omitempty,uuid"`
}

// ChangePasswordRequest is the DTO for a user changing their own password.
type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
