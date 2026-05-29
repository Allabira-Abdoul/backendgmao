package primary

import (
	"context"

	"backend-gmao/apps/auth-service/internal/core/domain"
)

// AuthService defines primary application business operations.
type AuthService interface {
	CreateSession(ctx context.Context, req domain.CreateSessionRequest) (*domain.SessionResponse, error)
	ValidateSession(ctx context.Context, token string) (*domain.SessionResponse, error)
	RevokeSession(ctx context.Context, token string) error
	RefreshSession(ctx context.Context, req domain.RefreshSessionRequest) (*domain.SessionResponse, error)
}
