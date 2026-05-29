package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
)

// AuditEvent represents the payload to send to the audit service.
type AuditEvent struct {
	ServiceName string  `json:"service_name"`
	Action      string  `json:"action"`
	Details     string  `json:"details,omitempty"`
	UserID      *string `json:"user_id,omitempty"`
	UserName    string  `json:"user_name,omitempty"`
}

// Client defines the interface for logging audit events.
type Client interface {
	LogEvent(ctx context.Context, event AuditEvent) error
}

type auditClient struct {
	gatewayURL  string
	serviceName string
	jwtManager  *auth.JWTManager
	httpClient  *http.Client
}

// NewClient creates a new HTTP Audit Client.
func NewClient(serviceName string, jwtManager *auth.JWTManager) Client {
	url := os.Getenv("AUDIT_SERVICE_URL")
	if url == "" {
		url = "http://127.0.0.1:8106" // Direct port for audit-service
	}

	return &auditClient{
		gatewayURL:  url,
		serviceName: serviceName,
		jwtManager:  jwtManager,
		httpClient:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *auditClient) LogEvent(ctx context.Context, event AuditEvent) error {
	if event.UserName == "" {
		if name, ok := ctx.Value(middleware.ContextKeyFullName).(string); ok {
			event.UserName = name
		}
	}

	reqURL := fmt.Sprintf("%s/internal/audit-logs", c.gatewayURL)

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal audit event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Generate internal token
	token, err := c.jwtManager.GenerateInternalServiceToken(c.serviceName)
	if err != nil {
		return fmt.Errorf("failed to generate internal token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Service", c.serviceName)
	req.Header.Set("Content-Type", "application/json")

	// We'll execute synchronously, the service layer can decide to ignore the error or fire asynchronously.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("Audit request failed: %v\n", err)
		logErrorToFile(fmt.Sprintf("Audit request failed: %v\n", err))
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err := fmt.Errorf("audit service returned status: %d", resp.StatusCode)
		fmt.Printf("Audit logging failed: %v\n", err)
		logErrorToFile(fmt.Sprintf("Audit logging failed: %v\n", err))
		return err
	}

	return nil
}

func logErrorToFile(msg string) {
	f, err := os.OpenFile(`d:\1. STAGE DE FIN D'ETUDES\2.PROJET\backend\audit_errors.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.WriteString(time.Now().Format(time.RFC3339) + " - " + msg)
		f.Close()
	}
}
