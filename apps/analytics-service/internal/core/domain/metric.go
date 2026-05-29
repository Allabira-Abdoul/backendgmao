package domain

import (
	"time"

	"github.com/google/uuid"
)

// Metric represents a recorded performance or operational metric.
type Metric struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	Value     float64   `gorm:"column:value;not null" json:"value"`
	Category  string    `gorm:"column:category;not null" json:"category"`
	Timestamp time.Time `gorm:"column:timestamp;not null;default:current_timestamp" json:"timestamp"`
}

// TableName overrides GORM's default table name.
func (Metric) TableName() string {
	return "metrics"
}

// MetricResponse represents the API DTO for Metric.
type MetricResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Category  string    `json:"category"`
	Timestamp time.Time `json:"timestamp"`
}

// ToResponse converts a Metric to MetricResponse DTO.
func (m *Metric) ToResponse() MetricResponse {
	return MetricResponse{
		ID:        m.ID,
		Name:      m.Name,
		Value:     m.Value,
		Category:  m.Category,
		Timestamp: m.Timestamp,
	}
}

// CreateMetricRequest is the DTO used to submit a new metric.
type CreateMetricRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=255"`
	Value    float64 `json:"value" binding:"required"`
	Category string  `json:"category" binding:"required"`
}
