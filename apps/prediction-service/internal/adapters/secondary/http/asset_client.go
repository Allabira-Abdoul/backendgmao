package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"backend-gmao/apps/prediction-service/internal/core/ports/secondary"
	"backend-gmao/pkg/auth"
	"github.com/google/uuid"
)

type assetClient struct {
	gatewayURL string
	jwtManager *auth.JWTManager
	httpClient *http.Client
}

func NewAssetClient(jwtManager *auth.JWTManager) secondary.AssetClient {
	url := os.Getenv("ASSET_SERVICE_URL")
	if url == "" {
		url = "http://127.0.0.1:8083"
	}
	return &assetClient{
		gatewayURL: url,
		jwtManager: jwtManager,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *assetClient) GetAssetName(ctx context.Context, id uuid.UUID) (string, error) {
	reqURL := fmt.Sprintf("%s/api/asset/assets/%s", c.gatewayURL, id.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}

	token, _ := c.jwtManager.GenerateInternalServiceToken("prediction-service")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("asset service returned %d", resp.StatusCode)
	}

	var envelope struct {
		Data struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return "", err
	}
	return envelope.Data.Name, nil
}
