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
	importPrimaryEventBus "backend-gmao/apps/asset-service/internal/adapters/primary/eventbus"
	importEventBus "backend-gmao/apps/asset-service/internal/adapters/secondary/eventbus"
	pgadapter "backend-gmao/apps/asset-service/internal/adapters/secondary/postgres"
	"backend-gmao/apps/asset-service/internal/application/service"
	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/db"
	"backend-gmao/pkg/discovery"
	"backend-gmao/pkg/eventbus"

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
		&domain.Site{}, 
		&domain.System{}, 
		&domain.Asset{}, 
		&domain.Subsystem{}, 
		&domain.InventoryItem{},
		&domain.Component{},
	); err != nil {
		log.Fatalf("Failed to migrate Asset tables: %v", err)
	}

	log.Println("Database migrations completed")

	// --- Seed Default Data ---
	pgadapter.SeedData(database)

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

	_ = auth.NewJWTManager(jwtSecret, accessExpiry, refreshExpiry)

	// --- Repositories (Secondary Adapters) ---
	assetRepo := pgadapter.NewAssetRepository(database)

	// --- EventBus (RabbitMQ) ---
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	bus, err := eventbus.NewRabbitMQBus(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	_ = importEventBus.NewRabbitMQPublisher(bus)

	// --- Application Services ---
	assetService := service.NewAssetService(assetRepo)

	// --- Start Consuming Events ---
	importPrimaryEventBus.StartConsumingWorkOrderEvents(bus, assetService)

	// --- Register with Consul ---
	err = registry.Register(serviceID, serviceName, host, port)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// --- Initialize Gin Router ---
	router := gin.Default()

	// Register all routes
	assetHandler := httphandler.NewAssetHandler(assetService)
	httphandler.RegisterRoutes(router, assetHandler)

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
