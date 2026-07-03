package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"backend-gmao/pkg/discovery"
	"backend-gmao/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Consul Config
	consulHost := os.Getenv("CONSUL_HOST")
	if consulHost == "" {
		consulHost = "127.0.0.1"
	}
	consulPort := os.Getenv("CONSUL_PORT")
	if consulPort == "" {
		consulPort = "8500"
	}
	consulURL := net.JoinHostPort(consulHost, consulPort)

	// Initialize Consul Registry
	registry, err := discovery.NewConsulRegistry(consulURL)
	if err != nil {
		log.Fatalf("Failed to create consul registry: %v", err)
	}
	//
	// Service Config
	serviceID := fmt.Sprintf("api-gateway-%s", os.Getenv("INSTANCE_ID"))
	if os.Getenv("INSTANCE_ID") == "" {
		serviceID = "api-gateway-1"
	}
	serviceName := "api-gateway"
	host := os.Getenv("SERVICE_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := 8080
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &port)
	}

	// Register the gateway itself with Consul
	err = registry.Register(serviceID, serviceName, host, port)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(middleware.Cors())

	// 🛡️ Security: Add global rate limiting middleware (defense-in-depth)
	router.Use(middleware.RateLimit(100, 1*time.Minute))

	// Add global security headers middleware
	router.Use(func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		// Enforce HTTPS (HTTP Strict Transport Security)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	})

	// Gateway health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"service": "api-gateway",
		})
	})

	// Allowed public services whitelist to prevent exposing internal services
	allowedServices := map[string]bool{
		"analytics":      true,
		"asset":          true,
		"auth":           true,
		"maintenance":    true,
		"prediction":     true,
		"user":           true,
		"identity":       true,
		"audit":          true,
	}

	// Dynamic reverse proxy for all service routes
	// Pattern: /api/{service-name}/* -> {service-name}-service/
	// Example: GET /api/analytics/health -> analytics-service /health
	router.Any("/api/:service/*path", func(c *gin.Context) {
		serviceSuffix := c.Param("service")
		targetPath := c.Param("path")

		// 🛡️ Security: Clean the path to prevent path traversal (SSRF)
		cleanedPath := path.Clean("/" + targetPath)

		// 🛡️ Security: Block external access to internal endpoints
		if strings.HasPrefix(cleanedPath, "/internal/") || cleanedPath == "/internal" {
			log.Printf("Security alert: Blocked attempt to access internal endpoint via gateway: %s", cleanedPath)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Access to internal endpoints is not allowed",
			})
			return
		}

		// 🛡️ Security: Enforce service whitelist to prevent SSRF / unauthorized access to internal services
		if !allowedServices[serviceSuffix] {
			log.Printf("Security alert: Blocked attempt to access unauthorized service: %s", serviceSuffix)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Access to this service is not allowed",
			})
			return
		}

		// 🛡️ Security: Block external access to internal endpoints
		if strings.HasPrefix(targetPath, "/internal/") || targetPath == "/internal" || targetPath == "internal" || strings.HasPrefix(targetPath, "internal/") {
			log.Printf("Security alert: Blocked attempt to access internal endpoint: %s", targetPath)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Access to internal endpoints is not allowed",
			})
			return
		}

		// Build the Consul service name
		targetServiceName := serviceSuffix + "-service"
		if serviceSuffix == "user" || serviceSuffix == "auth" {
			targetServiceName = "identity-service"
		}

		// Discover the service via Consul
		addr, err := registry.Discover(targetServiceName)
		if err != nil {
			log.Printf("Service discovery failed for %s: %v", targetServiceName, err)
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   "service_unavailable",
				"message": fmt.Sprintf("Could not discover service: %s", targetServiceName),
			})
			return
		}

		// Build target URL
		// ⚡ Bolt Optimization: Replace fmt.Sprintf + url.Parse with direct struct instantiation
		// to avoid string allocation and parsing overhead in the routing hot path.
		targetURL := &url.URL{
			Scheme: "http",
			Host:   addr,
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Modify the request
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			// Use the previously cleaned path
			req.URL.Path = cleanedPath
			req.URL.RawQuery = c.Request.URL.RawQuery
			req.Host = targetURL.Host

			// ⚡ Bolt Optimization: Removed redundant O(N) header copying loop.
			// The proxy's incoming request clone already contains all original headers.

			// 🛡️ Security: Drop internal headers to prevent external spoofing
			req.Header.Del("X-Internal-Service")

			// Add gateway-specific headers
			req.Header.Set("X-Forwarded-For", c.ClientIP())
			req.Header.Set("X-Forwarded-Host", c.Request.Host)
			req.Header.Set("X-Gateway-Service", "api-gateway")
		}

		// Custom error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error for %s: %v", targetServiceName, err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf(`{"error":"proxy_error","message":"Failed to reach %s"}`, targetServiceName)))
		}

		log.Printf("Proxying request: %s %s -> %s%s", c.Request.Method, c.Request.URL.Path, addr, targetPath)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Also support /api/{service} without trailing path (e.g. /api/analytics)
	router.Any("/api/:service", func(c *gin.Context) {
		// Redirect to include a trailing slash
		c.Request.URL.Path = strings.TrimRight(c.Request.URL.Path, "/") + "/"
		router.HandleContext(c)
	})

	// Start Server
	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Printf("Starting %s on %s\n", serviceName, addr)
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API Gateway...")

	if err := registry.Deregister(serviceID); err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}
	log.Println("API Gateway stopped")
}
