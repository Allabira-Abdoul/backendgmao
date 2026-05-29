package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"backend-gmao/apps/analytics-service/internal/core/ports/secondary"
	"backend-gmao/pkg/auth"
	"github.com/google/uuid"
)

type assetClient struct {
	gatewayURL string
	jwtManager *auth.JWTManager
	httpClient *http.Client
}

// NewAssetClient creates a new HTTP AssetClient.
func NewAssetClient(jwtManager *auth.JWTManager) secondary.AssetClient {
	url := os.Getenv("ASSET_SERVICE_URL")
	if url == "" {
		url = "http://127.0.0.1:8102"
	}

	return &assetClient{
		gatewayURL: url,
		jwtManager: jwtManager,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Struct to match the expected JSON response from the asset service via API Gateway
type gatewayResponse struct {
	ID           uuid.UUID `json:"id"`
	Category     string    `json:"category"`
	PurchaseDate string    `json:"purchase_date"`
}

func (c *assetClient) GetAssetInfo(ctx context.Context, assetID uuid.UUID) (*secondary.AssetInfo, error) {
	reqURL := fmt.Sprintf("%s/assets/%s", c.gatewayURL, assetID.String())
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Generate internal token
	token, err := c.jwtManager.GenerateInternalServiceToken("analytics-service")
	if err != nil {
		return nil, fmt.Errorf("failed to generate internal token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Service", "analytics-service")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("asset service returned status: %d", resp.StatusCode)
	}

	var gResp gatewayResponse
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	parsedDate, err := time.Parse(time.RFC3339, gResp.PurchaseDate)
	if err != nil {
		// Fallback to current time if parsing fails
		parsedDate = time.Now()
	}

	return &secondary.AssetInfo{
		ID:           gResp.ID,
		Category:     gResp.Category,
		PurchaseDate: parsedDate,
	}, nil
}
