package postgres

import (
	"context"

	"backend-gmao/apps/user-service/internal/core/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, team *domain.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *TeamRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	var team domain.Team
	if err := r.db.WithContext(ctx).First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) FindByName(ctx context.Context, name string) (*domain.Team, error) {
	var team domain.Team
	if err := r.db.WithContext(ctx).First(&team, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Team, error) {
	var teams []domain.Team
	if err := r.db.WithContext(ctx).Find(&teams, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *TeamRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Team, int64, error) {
	var teams []domain.Team
	var total int64
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&teams).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Model(&domain.Team{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return teams, total, nil
}

func (r *TeamRepository) Update(ctx context.Context, team *domain.Team) error {
	return r.db.WithContext(ctx).Save(team).Error
}

func (r *TeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Team{}, "id = ?", id).Error
}
