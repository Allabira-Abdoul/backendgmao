package postgres

import (
	"context"


	"backend-gmao/apps/identity-service/internal/core/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a GORM session repository.
func NewSessionRepository(db *gorm.DB) *sessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	err := r.db.WithContext(ctx).Create(session).Error
	return session, err
}



func (r *sessionRepository) Logout(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&domain.Session{}, "access_token = ? OR refresh_token = ?", token, token).Error
}

func (r *sessionRepository) FindByAccessToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "access_token = ?", token).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) FindByRefreshToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "refresh_token = ?", token).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	var sessions []domain.Session
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}
