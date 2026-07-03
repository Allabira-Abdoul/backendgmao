package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"backend-gmao/apps/identity-service/internal/core/domain"
	"backend-gmao/apps/identity-service/internal/core/ports/primary"
	"backend-gmao/apps/identity-service/internal/core/ports/secondary"
	"backend-gmao/pkg/auth"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailExists      = errors.New("a user with this email already exists")
	ErrInvalidAccount   = errors.New("user account is not active")
	ErrCannotDeleteSelf = errors.New("you cannot delete your own account")
)

// UserService implements the primary.UserUseCase primary port.
type UserService struct {
	userRepo       secondary.UserRepository
	roleRepo       secondary.RoleRepository
	teamRepo       secondary.TeamRepository
	eventPublisher secondary.EventPublisher
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepo secondary.UserRepository, roleRepo secondary.RoleRepository, teamRepo secondary.TeamRepository, eventPublisher secondary.EventPublisher) primary.UserUseCase {
	return &UserService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		teamRepo:       teamRepo,
		eventPublisher: eventPublisher,
	}
}

// CreateUser creates a new user after validating business rules.
func (s *UserService) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
	// Check if email already exists
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailExists
	}

	// Validate role exists
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID: %w", err)
	}

	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil || role == nil {
		return nil, ErrRoleNotFound
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashedPassword,
		Status:   domain.StatusActive,
		MustChangePassword: true,
		RoleID:   roleID,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Publish audit event
	s.eventPublisher.PublishAuditLog(ctx, "USER_CREATED", "User", user.ID.String(), nil, map[string]interface{}{
		"email": user.Email,
	})

	// Reload the user with role preloaded
	created, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload user: %w", err)
	}

	resp := created.ToResponse()
	return &resp, nil
}

// GetUserByID retrieves a user by their UUID.
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	resp := user.ToResponse()
	return &resp, nil
}

// GetUserByIDInternal retrieves a user by UUID for internal authentication use.
func (s *UserService) GetUserByIDInternal(ctx context.Context, id uuid.UUID) (*domain.InternalUserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	resp := user.ToInternalResponse()
	return &resp, nil
}

// GetUserByEmail retrieves a user by email for internal authentication use.
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.InternalUserResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	resp := user.ToInternalResponse()
	return &resp, nil
}

// ListUsers returns a paginated list of users.
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]domain.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAll(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	responses := make([]domain.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}

	return responses, total, nil
}

// GetCompactUsers returns a lightweight list of users.
func (s *UserService) GetCompactUsers(ctx context.Context) ([]domain.CompactUserResponse, error) {
	users, _, err := s.userRepo.FindAll(ctx, 0, 1000) // 1000 max for dropdowns
	if err != nil {
		return nil, fmt.Errorf("failed to list compact users: %w", err)
	}

	responses := make([]domain.CompactUserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, u.ToCompactResponse())
	}

	return responses, nil
}

// UpdateUser updates an existing user's fields.
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	changes := make(map[string]interface{})

	if req.FullName != nil {
		changes["full_name"] = map[string]interface{}{"old": user.FullName, "new": *req.FullName}
		user.FullName = *req.FullName
	}

	if req.Email != nil {
		// Check email uniqueness
		existing, _ := s.userRepo.FindByEmail(ctx, *req.Email)
		if existing != nil && existing.ID != id {
			return nil, ErrEmailExists
		}
		changes["email"] = map[string]interface{}{"old": user.Email, "new": *req.Email}
		user.Email = *req.Email
	}

	if req.Status != nil {
		changes["status"] = map[string]interface{}{"old": string(user.Status), "new": *req.Status}
		user.Status = domain.AccountStatus(*req.Status)
	}

	if req.RoleID != nil {
		roleID, err := uuid.Parse(*req.RoleID)
		if err != nil {
			return nil, fmt.Errorf("invalid role ID: %w", err)
		}
		role, err := s.roleRepo.FindByID(ctx, roleID)
		if err != nil || role == nil {
			return nil, ErrRoleNotFound
		}
		changes["role_id"] = map[string]interface{}{"old": user.RoleID.String(), "new": roleID.String()}
		user.RoleID = roleID
	}

	if req.TeamID != nil {
		if *req.TeamID == "" {
			changes["team_id"] = map[string]interface{}{"old": user.TeamID, "new": nil}
			user.TeamID = nil
		} else {
			teamID, err := uuid.Parse(*req.TeamID)
			if err != nil {
				return nil, fmt.Errorf("invalid team ID: %w", err)
			}
			changes["team_id"] = map[string]interface{}{"old": user.TeamID, "new": teamID.String()}
			user.TeamID = &teamID
		}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if len(changes) > 0 {
		s.eventPublisher.PublishAuditLog(ctx, "USER_UPDATE", "User", user.ID.String(), nil, changes)
	}

	// Reload with role
	updated, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload user: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// DeleteUser removes a user by their UUID.
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.eventPublisher.PublishAuditLog(ctx, "USER_DELETE", "User", id.String(), nil, nil)
	return nil
}

// SuspendUser sets a user's status to INACTIVE.
func (s *UserService) SuspendUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	oldStatus := user.Status
	user.Suspend()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to suspend user: %w", err)
	}

	s.eventPublisher.PublishAuditLog(ctx, "USER_SUSPEND", "User", user.ID.String(), nil, map[string]interface{}{
		"status": map[string]interface{}{"old": string(oldStatus), "new": string(user.Status)},
	})

	return nil
}

// AssignRoleToUser assigns a new role to the user and publishes an audit log.
func (s *UserService) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil || role == nil {
		return ErrRoleNotFound
	}

	oldRoleID := user.RoleID
	user.RoleID = roleID

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Using the EventPublisher to create the Audit Log
	s.eventPublisher.PublishAuditLog(ctx, "USER_ASSIGN_ROLE", "User", user.ID.String(), nil, map[string]interface{}{
		"role_id": map[string]interface{}{"old": oldRoleID.String(), "new": roleID.String()},
	})

	return nil
}

// generateRandom6DigitCode generates a secure 6-digit numeric string.
func generateRandom6DigitCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

// AdminResetPassword generates a 6-digit code, updates the user's password, and forces them to change it.
func (s *UserService) AdminResetPassword(ctx context.Context, id uuid.UUID) (string, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return "", ErrUserNotFound
	}

	code := generateRandom6DigitCode()
	hashedPassword, err := auth.HashPassword(code)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	user.MustChangePassword = true

	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", fmt.Errorf("failed to update user password: %w", err)
	}

	s.eventPublisher.PublishAuditLog(ctx, "USER_RESET_PASSWORD", "User", user.ID.String(), nil, nil)
	return code, nil
}

// ChangePassword allows a user to change their own password, lifting the MustChangePassword flag.
func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	user.MustChangePassword = false

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	s.eventPublisher.PublishAuditLog(ctx, "USER_CHANGE_PASSWORD", "User", user.ID.String(), nil, nil)
	return nil
}
