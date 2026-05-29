package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend-gmao/pkg/audit"
	"backend-gmao/pkg/auth"
)

func main() {
	jwtManager := auth.NewJWTManager("gmao-dev-secret-change-in-production", time.Minute*5, time.Minute*5)
	client := audit.NewClient("test-service", jwtManager)
	
	event := audit.AuditEvent{
		ServiceName: "test-service",
		Action:      "TEST_NIL_USER",
		Details:     "Testing audit service with nil user",
		UserID:      nil, // nil user
	}
	
	err := client.LogEvent(context.Background(), event)
	if err != nil {
		log.Fatalf("Failed to log event: %v", err)
	}
	
	fmt.Println("Successfully logged event with nil user!")
}
