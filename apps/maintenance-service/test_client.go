package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"backend-gmao/pkg/auth"
)

func main() {
	jwtManager := auth.NewJWTManager("gmao-dev-secret-change-in-production", 15*time.Minute, 24*time.Hour)
	token, _ := jwtManager.GenerateInternalServiceToken("maintenance-service")
	
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8102/instances/equipment/b5f2f1ff-df3c-4cc5-9915-d25be8066de2", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\nBody: %s\n", resp.StatusCode, string(body))
}
