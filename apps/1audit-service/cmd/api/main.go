package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	eventbusadapter "backend-gmao/apps/audit-service/internal/adapters/primary/eventbus"
	httphandler "backend-gmao/apps/audit-service/internal/adapters/primary/http"
	pgadapter "backend-gmao/apps/audit-service/internal/adapters/secondary/postgres"
	"backend-gmao/apps/audit-service/internal/application/service"
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
	serviceID := fmt.Sprintf("audit-service-%s", getEnv("INSTANCE_ID", "1"))
	serviceName := "audit-service"
	host := getEnv("SERVICE_HOST", "127.0.0.1")

	port := 8086
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

	// --- Initialize Partitioned Table ---
	log.Println("Initializing Postgres partitioned table...")
	if err := pgadapter.InitPartitionedTable(database); err != nil {
		log.Fatalf("Failed to init partitioned table: %v", err)
	}
	log.Println("Postgres partitioning setup completed")

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
	auditRepo := pgadapter.NewAuditRepository(database)

	// --- Application Services ---
	auditService := service.NewAuditService(auditRepo)

	// --- EventBus Config & Connection ---
	rabbitMQUrl := getEnv("RABBITMQ_URL", "amqp://guest:guest@127.0.0.1:5672/")
	bus, err := eventbus.NewRabbitMQBus(rabbitMQUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer bus.Close()

	// --- Start RabbitMQ Consumer ---
	rabbitMQConsumer := eventbusadapter.NewRabbitMQConsumer(bus, auditService)
	if err := rabbitMQConsumer.Start(); err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
	}

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
	httphandler.RegisterRoutes(router, jwtManager, auditService)

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
