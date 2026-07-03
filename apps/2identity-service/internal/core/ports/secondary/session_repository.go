package secondary

import (
	"context"

	"backend-gmao/apps/identity-service/internal/core/domain"
	"github.com/google/uuid"
)

// SessionRepository defines secondary adapter database actions.
type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) (*domain.Session, error)
	Logout(ctx context.Context, token string) error
	FindByAccessToken(ctx context.Context, token string) (*domain.Session, error)
	FindByRefreshToken(ctx context.Context, token string) (*domain.Session, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Session, error)
}
