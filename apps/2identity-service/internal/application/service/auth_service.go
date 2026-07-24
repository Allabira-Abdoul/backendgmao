package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend-gmao/apps/identity-service/internal/core/domain"
	"backend-gmao/apps/identity-service/internal/core/ports/secondary"
	"backend-gmao/pkg/auth"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	sessionRepo secondary.SessionRepository
	userRepo    secondary.UserRepository
	eventPub    secondary.EventPublisher
	jwtManager  *auth.JWTManager
}

// NewAuthService creates a new authentication service.
func NewAuthService(
	sessionRepo secondary.SessionRepository,
	userRepo secondary.UserRepository,
	eventPub secondary.EventPublisher,
	jwtManager *auth.JWTManager,
) *AuthService {
	return &AuthService{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		eventPub:    eventPub,
		jwtManager:  jwtManager,
	}
}

func (s *AuthService) CreateSession(ctx context.Context, req domain.CreateSessionRequest) (*domain.SessionResponse, error) {
	// 1. Fetch user via UserRepository
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Event logging could happen here, though we might not have a reliable UserID if they don't exist
		_ = s.eventPub.PublishAuditLog(ctx, "LOGIN_FAILED", "USER", req.Email, nil, map[string]interface{}{"reason": "invalid credentials"})
		return nil, ErrInvalidCredentials
	}

	// 2. Verify user status
	if user.Status != domain.StatusActive {
		_ = s.eventPub.PublishAuditLog(ctx, "LOGIN_FAILED", "USER", user.ID.String(), &user.ID, map[string]interface{}{"reason": "account inactive"})
		return nil, fmt.Errorf("account is %s", strings.ToLower(string(user.Status)))
	}

	// 3. Verify password
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		_ = s.eventPub.PublishAuditLog(ctx, "LOGIN_FAILED", "USER", user.ID.String(), &user.ID, map[string]interface{}{"reason": "invalid password"})
		return nil, ErrInvalidCredentials
	}

	siteIDStr := ""
	if user.SiteID != nil {
		siteIDStr = user.SiteID.String()
	}

	// 4. Generate signed JWT access token and refresh token
	accessToken, accessExpiredAt, err := s.jwtManager.GenerateAccessToken(
		user.ID.String(),
		siteIDStr,
		user.Email,
		user.FullName,
		user.Role.Name,
		user.Role.GetPrivilegeStrings(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshExpiredAt, err := s.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 5. Save session in DB
	session := &domain.Session{
		ID:               uuid.New(),
		UserID:           user.ID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiredAt:  accessExpiredAt,
		RefreshExpiredAt: refreshExpiredAt,
	}

	if _, err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	// Trigger audit event asynchronously
	go func() {
		bgCtx := context.Background()
		_ = s.eventPub.PublishAuditLog(bgCtx, "USER_LOGIN", "SESSION", session.ID.String(), &user.ID, nil)
	}()

	sessionResp := session.ToResponse()
	return &sessionResp, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*domain.SessionResponse, error) {
	session, err := s.sessionRepo.FindByAccessToken(ctx, token)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.AccessExpiredAt) {
		s.sessionRepo.Logout(ctx, token)
		return nil, ErrSessionExpired
	}

	resp := session.ToResponse()
	return &resp, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, token string) error {
	return s.sessionRepo.Logout(ctx, token)
}

func (s *AuthService) RefreshSession(ctx context.Context, req domain.RefreshSessionRequest) (*domain.SessionResponse, error) {
	// 1. Validate the refresh token string
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 2. Find the session in the DB using the refresh token
	session, err := s.sessionRepo.FindByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.RefreshExpiredAt) {
		s.sessionRepo.Logout(ctx, req.RefreshToken)
		return nil, ErrSessionExpired
	}

	// 3. Delete old session to enforce rotation
	_ = s.sessionRepo.Logout(ctx, req.RefreshToken)

	// 4. Fetch freshest user data via UserRepository
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in claims: %w", err)
	}
	
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fresh user data: %w", err)
	}

	if user.Status != domain.StatusActive {
		return nil, fmt.Errorf("account is %s", strings.ToLower(string(user.Status)))
	}

	siteIDStr := ""
	if user.SiteID != nil {
		siteIDStr = user.SiteID.String()
	}

	// 5. Generate new token pair
	newAccessToken, newAccessExpiredAt, err := s.jwtManager.GenerateAccessToken(
		user.ID.String(),
		siteIDStr,
		user.Email,
		user.FullName,
		user.Role.Name,
		user.Role.GetPrivilegeStrings(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	newRefreshToken, newRefreshExpiredAt, err := s.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// 6. Save new session
	newSession := &domain.Session{
		ID:               uuid.New(),
		UserID:           user.ID,
		AccessToken:      newAccessToken,
		RefreshToken:     newRefreshToken,
		AccessExpiredAt:  newAccessExpiredAt,
		RefreshExpiredAt: newRefreshExpiredAt,
	}

	if _, err := s.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, err
	}

	sessionResp := newSession.ToResponse()
	return &sessionResp, nil
}
