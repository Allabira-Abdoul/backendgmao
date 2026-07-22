package audit

import (
	"context"

	"backend-gmao/apps/maintenance-service/internal/core/ports/secondary"
	pkgAudit "backend-gmao/pkg/audit"
	"backend-gmao/pkg/middleware"
	"github.com/google/uuid"
)

type compositeLogger struct {
	httpClient pkgAudit.Client
	eventPub   secondary.EventPublisher
}

// NewCompositeLogger creates a new secondary.AuditLogger that delegates to an HTTP client and/or an EventPublisher.
func NewCompositeLogger(httpClient pkgAudit.Client, eventPub secondary.EventPublisher) secondary.AuditLogger {
	return &compositeLogger{
		httpClient: httpClient,
		eventPub:   eventPub,
	}
}

// LogAction logs the given action and details, optionally extracting User metadata from the context.
func (c *compositeLogger) LogAction(ctx context.Context, action, details string) {
	userID, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	var uidPtr *string
	if ok && userID != "" {
		uidPtr = &userID
	}

	var userName string
	if name, ok := ctx.Value(middleware.ContextKeyFullName).(string); ok {
		userName = name
	}

	go func() {
		bgCtx := context.Background()

		// Publish to HTTP Audit Service
		if c.httpClient != nil {
			_ = c.httpClient.LogEvent(bgCtx, pkgAudit.AuditEvent{
				ServiceName: "maintenance-service",
				Action:      action,
				Details:     details,
				UserID:      uidPtr,
				UserName:    userName,
			})
		}

		// Publish to RabbitMQ Audit Exchange
		if c.eventPub != nil {
			var actorID *uuid.UUID
			if uidPtr != nil {
				parsed, err := uuid.Parse(*uidPtr)
				if err == nil {
					actorID = &parsed
				}
			}
			changes := map[string]interface{}{
				"details": details,
			}
			if userName != "" {
				changes["user_name"] = userName
			}
			_ = c.eventPub.PublishAuditLog(bgCtx, action, "MAINTENANCE", "", actorID, changes)
		}
	}()
}
