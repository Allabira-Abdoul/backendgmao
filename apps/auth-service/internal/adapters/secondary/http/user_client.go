package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"backend-gmao/apps/auth-service/internal/core/domain"
	"backend-gmao/apps/auth-service/internal/core/ports/secondary"
	"backend-gmao/pkg/discovery"
	"github.com/google/uuid"
)

type UserClient struct {
	registry discovery.Registry
	client   *http.Client
}

// NewUserClient creates a new HTTP client for user-service interactions.
func NewUserClient(registry discovery.Registry) secondary.UserProvider {
	return &UserClient{
		registry: registry,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *UserClient) FetchUserForAuth(ctx context.Context, email string) (*domain.User, error) {
	addr, err := c.registry.Discover("user-service")
	if err != nil {
		return nil, fmt.Errorf("failed to discover user-service: %w", err)
	}

	targetURL := fmt.Sprintf("http://%s/internal/by-email?email=%s", addr, url.QueryEscape(email))
	httpReq, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("X-Internal-Service", "auth-service")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status: %d", resp.StatusCode)
	}

	var envelope struct {
		Status string      `json:"status"`
		Data   domain.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &envelope.Data, nil
}

func (c *UserClient) FetchUserByIDForAuth(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	addr, err := c.registry.Discover("user-service")
	if err != nil {
		return nil, fmt.Errorf("failed to discover user-service: %w", err)
	}

	targetURL := fmt.Sprintf("http://%s/internal/by-id?id=%s", addr, url.QueryEscape(id.String()))
	httpReq, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("X-Internal-Service", "auth-service")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status: %d", resp.StatusCode)
	}

	var envelope struct {
		Status string      `json:"status"`
		Data   domain.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &envelope.Data, nil
}
