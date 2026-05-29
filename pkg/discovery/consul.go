package discovery

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	consul "github.com/hashicorp/consul/api"
	"golang.org/x/sync/singleflight"
)

// Registry defines the interface for service discovery operations (Secondary Port)
type Registry interface {
	Register(serviceID, serviceName, host string, port int) error
	Deregister(serviceID string) error
	Discover(serviceName string) (string, error)
}

type cacheEntry struct {
	addr      string
	expiresAt time.Time
}

// ConsulRegistry is a Consul-backed implementation of the Registry interface (Secondary Adapter)
type ConsulRegistry struct {
	client *consul.Client
	cache  map[string]cacheEntry
	mu     sync.RWMutex
	sf     singleflight.Group
}

// NewConsulRegistry creates a new ConsulRegistry connected to the given Consul address.
func NewConsulRegistry(addr string) (*ConsulRegistry, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &ConsulRegistry{
		client: client,
		cache:  make(map[string]cacheEntry),
	}, nil
}

// Register registers a service instance with Consul, including an HTTP health check.
func (r *ConsulRegistry) Register(serviceID, serviceName, host string, port int) error {
	registration := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: host,
		Port:    port,
		Check: &consul.AgentServiceCheck{
			HTTP:                           "http://" + net.JoinHostPort(host, strconv.Itoa(port)) + "/health",
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	err := r.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service %s: %w", serviceName, err)
	}

	log.Printf("Service %s registered with Consul (ID: %s, Address: %s:%d)", serviceName, serviceID, host, port)
	return nil
}

// Deregister removes a service instance from Consul.
func (r *ConsulRegistry) Deregister(serviceID string) error {
	err := r.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service %s: %w", serviceID, err)
	}

	log.Printf("Service %s deregistered from Consul", serviceID)
	return nil
}

// Discover finds a healthy instance of a service by name and returns its address.
// ⚡ Bolt Optimization: Added caching to avoid synchronous network calls to Consul on every request.
func (r *ConsulRegistry) Discover(serviceName string) (string, error) {
	// 1. Check cache first
	r.mu.RLock()
	entry, found := r.cache[serviceName]
	r.mu.RUnlock()

	if found && time.Now().Before(entry.expiresAt) {
		return entry.addr, nil
	}

	// 2. Cache miss or expired, fetch from Consul
	// Use singleflight to prevent cache stampedes (thundering herd)
	v, err, _ := r.sf.Do(serviceName, func() (interface{}, error) {
		entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
		if err != nil {
			// If Consul is down but we have a stale cache, we could optionally return it here.
			// For now, we return the error to maintain strict correctness, but if high availability
			// is preferred over strict correctness, we could return the stale entry.
			return "", fmt.Errorf("failed to discover service %s: %w", serviceName, err)
		}

		if len(entries) == 0 {
			return "", fmt.Errorf("no healthy instances found for service %s", serviceName)
		}

		// Return the first healthy instance
		consulEntry := entries[0]
		addr := net.JoinHostPort(consulEntry.Service.Address, strconv.Itoa(consulEntry.Service.Port))

		// 3. Update cache (10 seconds TTL)
		r.mu.Lock()
		r.cache[serviceName] = cacheEntry{
			addr:      addr,
			expiresAt: time.Now().Add(10 * time.Second),
		}
		r.mu.Unlock()

		return addr, nil
	})

	if err != nil {
		return "", err
	}

	return v.(string), nil
}
