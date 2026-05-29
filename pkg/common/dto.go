package common

import "github.com/google/uuid"

// ResourceRef represents a lightweight reference to another entity,
// including both its ID and a human-readable Name, to replace raw UUIDs in API responses.
type ResourceRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
