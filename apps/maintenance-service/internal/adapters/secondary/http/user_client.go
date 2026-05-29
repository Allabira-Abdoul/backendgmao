package http

import (
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

type userClient struct {
	gatewayURL string
	jwtManager *auth.JWTManager
	httpClient *http.Client
}

func NewUserClient(jwtManager *auth.JWTManager) secondary.UserClient {
	url := os.Getenv("USER_SERVICE_URL")
	if url == "" {
		url = "http://127.0.0.1:8081"
	}
	return &userClient{
		gatewayURL: url,
		jwtManager: jwtManager,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *userClient) GetUserName(ctx context.Context, id uuid.UUID) (string, error) {
	reqURL := fmt.Sprintf("%s/internal/by-id?id=%s", c.gatewayURL, id.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}

	token, _ := c.jwtManager.GenerateInternalServiceToken("maintenance-service")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Service", "maintenance-service")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("user service returned %d", resp.StatusCode)
	}

	var envelope struct {
		Data struct {
			FullName string `json:"full_name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return "", err
	}
	return envelope.Data.FullName, nil
}

func (c *userClient) GetUserNames(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error) {
	// A batch endpoint would be better, but we fall back to sequential for now.
	// Since performance is not a strict requirement for this POC, sequential is fine.
	result := make(map[uuid.UUID]string)
	for _, id := range ids {
		name, err := c.GetUserName(ctx, id)
		if err == nil {
			result[id] = name
		}
	}
	return result, nil
}
