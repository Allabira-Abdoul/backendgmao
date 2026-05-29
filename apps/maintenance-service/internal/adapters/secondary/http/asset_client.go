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
	// Using the public endpoint since asset-service doesn't have an internal endpoint yet
	reqURL := fmt.Sprintf("%s/api/asset/assets/%s", c.gatewayURL, id.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}

	token, _ := c.jwtManager.GenerateInternalServiceToken("maintenance-service")
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

func (c *assetClient) GetAssetNames(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error) {
	result := make(map[uuid.UUID]string)
	for _, id := range ids {
		name, err := c.GetAssetName(ctx, id)
		if err == nil {
			result[id] = name
		}
	}
	return result, nil
}

func (c *assetClient) UpdateAssetStatus(ctx context.Context, id uuid.UUID, status string) error {
	reqURL := fmt.Sprintf("%s/api/asset/assets/%s", c.gatewayURL, id.String())
	
	payload := map[string]string{"status": status}
	body, _ := json.Marshal(payload)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	token, _ := c.jwtManager.GenerateInternalServiceToken("maintenance-service")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("asset service returned %d when updating status", resp.StatusCode)
	}
	return nil
}
