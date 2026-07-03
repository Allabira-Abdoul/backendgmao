package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AuditLog represents a secure log entry of any action taken across microservices.
type AuditLog struct {
	ID           uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PerformedAt  time.Time      `gorm:"column:performed_at;not null;default:now()" json:"performed_at"`
	ActorID      *uuid.UUID     `gorm:"column:actor_id;type:uuid" json:"actor_id"`
	ServiceName  string         `gorm:"column:service_name;not null" json:"service_name"`
	Action       string         `gorm:"column:action;not null" json:"action"`
	ResourceType string         `gorm:"column:resource_type" json:"resource_type"`
	ResourceID   string         `gorm:"column:resource_id" json:"resource_id"`
	Changes      datatypes.JSON `gorm:"column:changes;type:jsonb" json:"changes"`
}

// TableName overrides GORM's default table name.
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogResponse represents the API DTO for AuditLog.
type AuditLogResponse struct {
	ID           uuid.UUID      `json:"id"`
	PerformedAt  time.Time      `json:"performed_at"`
	ActorID      *uuid.UUID     `json:"actor_id"`
	ServiceName  string         `json:"service_name"`
	Action       string         `json:"action"`
	ResourceType string         `json:"resource_type"`
	ResourceID   string         `json:"resource_id"`
	Changes      datatypes.JSON `json:"changes"`
}

// ToResponse converts an AuditLog to AuditLogResponse DTO.
func (l *AuditLog) ToResponse() AuditLogResponse {
	return AuditLogResponse{
		ID:           l.ID,
		PerformedAt:  l.PerformedAt,
		ActorID:      l.ActorID,
		ServiceName:  l.ServiceName,
		Action:       l.Action,
		ResourceType: l.ResourceType,
		ResourceID:   l.ResourceID,
		Changes:      l.Changes,
	}
}

// AuditFilter defines the parameters for querying audit logs.
type AuditFilter struct {
	ServiceName  string
	Action       string
	ResourceType string
	ResourceID   string
	ActorID      *uuid.UUID
	StartDate    *time.Time
	EndDate      *time.Time
	Limit        int
	Offset       int
}
