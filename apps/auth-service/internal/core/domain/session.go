package domain

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user session in the auth system.
type Session struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	AccessToken      string    `gorm:"column:access_token;uniqueIndex;not null" json:"access_token"`
	RefreshToken     string    `gorm:"column:refresh_token;uniqueIndex;not null" json:"refresh_token"`
	AccessExpiredAt  time.Time `gorm:"column:access_expired_at;not null" json:"access_expired_at"`
	RefreshExpiredAt time.Time `gorm:"column:refresh_expired_at;not null" json:"refresh_expired_at"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName overrides GORM's default table name.
func (Session) TableName() string {
	return "sessions"
}

// SessionResponse is the DTO returned by API endpoints.
type SessionResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiredAt  time.Time `json:"access_expired_at"`
	RefreshExpiredAt time.Time `json:"refresh_expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a Session to SessionResponse DTO.
func (s *Session) ToResponse() SessionResponse {
	return SessionResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		AccessToken:      s.AccessToken,
		RefreshToken:     s.RefreshToken,
		AccessExpiredAt:  s.AccessExpiredAt,
		RefreshExpiredAt: s.RefreshExpiredAt,
		CreatedAt: s.CreatedAt,
	}
}

// CreateSessionRequest is the DTO used to submit a new session.
type CreateSessionRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshSessionRequest is the DTO used to request a new session via a refresh token.
type RefreshSessionRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccountStatus string

const (
	StatusActive   AccountStatus = "ACTIVE"
	StatusInactive AccountStatus = "INACTIVE"
	StatusLocked   AccountStatus = "LOCKED"
)

type User struct {
	ID            uuid.UUID     `json:"id"`
	FullName      string        `json:"full_name"`
	Email         string        `json:"email"`
	Password      string        `json:"password"`
	Status        AccountStatus `json:"status"`
	RoleName      string        `json:"role_name"`
	Privileges    []string      `json:"privileges"`
}
