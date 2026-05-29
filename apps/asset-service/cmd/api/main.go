package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	httphandler "backend-gmao/apps/asset-service/internal/adapters/primary/http"
	pgadapter "backend-gmao/apps/asset-service/internal/adapters/secondary/postgres"
	"backend-gmao/apps/asset-service/internal/application/service"
	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/pkg/audit"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/db"
	"backend-gmao/pkg/discovery"

	"github.com/gin-gonic/gin"
)

func main() {
	// --- Consul Config ---
	consulHost := getEnv("CONSUL_HOST", "127.0.0.1")
	consulPort := getEnv("CONSUL_PORT", "8500")
	consulURL := net.JoinHostPort(consulHost, consulPort)

	registry, err := discovery.NewConsulRegistry(consulURL)
	if err != nil {
		log.Fatalf("Failed to create consul registry: %v", err)
	}

	// --- Service Config ---
	serviceID := fmt.Sprintf("asset-service-%s", getEnv("INSTANCE_ID", "1"))
	serviceName := "asset-service"
	host := getEnv("SERVICE_HOST", "127.0.0.1")

	port := 8082
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &port)
	}

	// --- Database Connection ---
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	database, err := db.InitPostgres(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// --- Auto-Migrate Tables ---
	log.Println("Running database migrations...")
	if err := database.AutoMigrate(
		&domain.EquipmentModel{}, 
		&domain.PartModel{}, 
		&domain.EquipmentInstance{}, 
		&domain.PartInstance{}, 
		&domain.MetricThreshold{},
		&domain.PartConsumptionLog{},
	); err != nil {
		log.Fatalf("Failed to migrate Asset tables: %v", err)
	}
	log.Println("Database migrations completed")

	// --- Seed Default Data ---
	pgadapter.Seed(database)

	// --- JWT Manager ---
	jwtSecret := getEnv("JWT_SECRET", "gmao-dev-secret-change-in-production")
	accessExpiry := 15 * time.Minute
	refreshExpiry := 7 * 24 * time.Hour

	if exp := os.Getenv("JWT_ACCESS_EXPIRY"); exp != "" {
		if d, err := time.ParseDuration(exp); err == nil {
			accessExpiry = d
		}
	}
	if exp := os.Getenv("JWT_REFRESH_EXPIRY"); exp != "" {
		if d, err := time.ParseDuration(exp); err == nil {
			refreshExpiry = d
		}
	}

	jwtManager := auth.NewJWTManager(jwtSecret, accessExpiry, refreshExpiry)

	// --- Repositories (Secondary Adapters) ---
	assetRepo := pgadapter.NewAssetRepository(database)

	// --- Application Services ---
	jwtManagerForInternal := auth.NewJWTManager(jwtSecret, time.Minute*5, time.Minute*5)
	auditClient := audit.NewClient("asset-service", jwtManagerForInternal)
	assetService := service.NewAssetService(assetRepo, auditClient)

	// --- Register with Consul ---
	err = registry.Register(serviceID, serviceName, host, port)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// --- Initialize Gin Router ---
	router := gin.Default()

	// Health check
	healthHandler := httphandler.NewHealthHandler(database)
	router.GET("/health", healthHandler.HealthCheck)

	// Register all routes
	httphandler.RegisterRoutes(router, jwtManager, assetService)

	// --- Start Server ---
	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Printf("Starting %s on %s\n", serviceName, addr)
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down service...")

	if err := registry.Deregister(serviceID); err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}
	log.Println("Service stopped")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
