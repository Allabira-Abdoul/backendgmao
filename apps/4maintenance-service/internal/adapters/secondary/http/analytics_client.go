package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"backend-gmao/apps/maintenance-service/internal/core/ports/secondary"
	"backend-gmao/pkg/auth"
)

type analyticsClient struct {
	gatewayURL string
	jwtManager *auth.JWTManager
	httpClient *http.Client
}

// NewAnalyticsClient creates a new HTTP AnalyticsClient.
func NewAnalyticsClient(jwtManager *auth.JWTManager) secondary.AnalyticsClient {
	url := os.Getenv("ANALYTICS_SERVICE_URL")
	if url == "" {
		url = "http://127.0.0.1:8101" // Direct port for analytics-service
	}

	return &analyticsClient{
		gatewayURL: url,
		jwtManager: jwtManager,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *analyticsClient) PublishMaintenanceEvent(ctx context.Context, event secondary.MaintenanceEvent) error {
	reqURL := fmt.Sprintf("%s/internal/analytics/events/maintenance-completed", c.gatewayURL)

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Generate internal token
	token, err := c.jwtManager.GenerateInternalServiceToken("maintenance-service")
	if err != nil {
		return fmt.Errorf("failed to generate internal token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Service", "maintenance-service")
	req.Header.Set("Content-Type", "application/json")

	// Make it fire-and-forget but log errors, or synchronous
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("analytics service returned status: %d", resp.StatusCode)
	}

	return nil
}
