package domain

import (
	"time"

	"github.com/google/uuid"
)

// Team represents a team or group of users in the GMAO system.
type Team struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"column:name;uniqueIndex;not null" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (Team) TableName() string {
	return "teams"
}

// TeamResponse is the DTO returned by API endpoints.
type TeamResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a Team to a TeamResponse.
func (e *Team) ToResponse() TeamResponse {
	return TeamResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// CreateTeamRequest is the DTO for creating a new team.
type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"omitempty"`
}

// UpdateTeamRequest is the DTO for updating an existing team.
type UpdateTeamRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty"`
}
