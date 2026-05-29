package service

import (
	"context"
	"errors"
	"fmt"

	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/apps/user-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

var (
	ErrTeamNotFound      = errors.New("team not found")
	ErrTeamNameExists    = errors.New("a team with this name already exists")
	ErrTeamHasUsers      = errors.New("cannot delete a team that is assigned to users")
)

// TeamService implements the primary.TeamUseCase primary port.
type TeamService struct {
	teamRepo secondary.TeamRepository
	userRepo secondary.UserRepository
}

// NewTeamService creates a new TeamService instance.
func NewTeamService(teamRepo secondary.TeamRepository, userRepo secondary.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// CreateTeam creates a new team.
func (s *TeamService) CreateTeam(ctx context.Context, req domain.CreateTeamRequest) (*domain.TeamResponse, error) {
	// Check if name already exists
	existing, _ := s.teamRepo.FindByName(ctx, req.Name)
	if existing != nil {
		return nil, ErrTeamNameExists
	}

	var managerID *uuid.UUID
	if req.ManagerID != nil {
		id, err := uuid.Parse(*req.ManagerID)
		if err != nil {
			return nil, fmt.Errorf("invalid manager id: %w", err)
		}
		managerID = &id
	}

	team := &domain.Team{
		Name:        req.Name,
		ManagerID:   managerID,
		Description: req.Description,
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	// Reload with members
	created, err := s.teamRepo.FindByID(ctx, team.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload team: %w", err)
	}

	resp := created.ToResponse()
	return &resp, nil
}

// GetTeamByID retrieves a team by its UUID.
func (s *TeamService) GetTeamByID(ctx context.Context, id uuid.UUID) (*domain.TeamResponse, error) {
	team, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTeamNotFound
	}

	resp := team.ToResponse()
	return &resp, nil
}

// ListTeams returns all teams.
func (s *TeamService) ListTeams(ctx context.Context, limit, offset int) ([]domain.TeamResponse, int64, error) {
	teams, total, err := s.teamRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list teams: %w", err)
	}

	responses := make([]domain.TeamResponse, 0, len(teams))
	for _, t := range teams {
		responses = append(responses, t.ToResponse())
	}

	return responses, total, nil
}

// UpdateTeam updates an existing team's name and/or description.
func (s *TeamService) UpdateTeam(ctx context.Context, id uuid.UUID, req domain.UpdateTeamRequest) (*domain.TeamResponse, error) {
	team, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTeamNotFound
	}

	if req.Name != nil {
		// Check uniqueness
		existing, _ := s.teamRepo.FindByName(ctx, *req.Name)
		if existing != nil && existing.ID != id {
			return nil, ErrTeamNameExists
		}
		team.Name = *req.Name
	}

	if req.Description != nil {
		team.Description = *req.Description
	}

	if req.ManagerID != nil {
		if *req.ManagerID == "" {
			team.ManagerID = nil
		} else {
			id, err := uuid.Parse(*req.ManagerID)
			if err != nil {
				return nil, fmt.Errorf("invalid manager id: %w", err)
			}
			team.ManagerID = &id
		}
	}

	if err := s.teamRepo.Update(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	// Reload
	updated, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload team: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// DeleteTeam removes a team if it has no assigned users.
func (s *TeamService) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	_, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return ErrTeamNotFound
	}

	// Check if any users are assigned to this team
	users, err := s.userRepo.FindByTeamID(ctx, id)
	if err == nil && len(users) > 0 {
		return ErrTeamHasUsers
	}

	if err := s.teamRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}
