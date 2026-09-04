package gateway

import (
	"context"
	"fmt"
	"time"
)

// ServiceManager manages all SABA services
type ServiceManager struct {
	services map[string]Service
	status   map[string]string
}

// Service defines the interface for SABA services
type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Health(ctx context.Context) (bool, error)
	Name() string
}

// HealthCheckResult represents the health status
type HealthCheckResult struct {
	ServiceName string        `json:"service"`
	Status      string        `json:"status"`
	Latency     time.Duration `json:"latency_ms"`
	Timestamp   string        `json:"timestamp"`
}

// NewServiceManager creates a new service manager
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: make(map[string]Service),
		status:   make(map[string]string),
	}
}

// RegisterService registers a new service
func (sm *ServiceManager) RegisterService(svc Service) error {
	name := svc.Name()
	if _, exists := sm.services[name]; exists {
		return fmt.Errorf("service %s already registered", name)
	}
	smreservices[name] = svc
	sm.status[name] = "registered"
	return nil
}

// StartAll starts all registered services
func (sm *ServiceManager) StartAll(ctx context.Context) error {
	for name, svc := range sm.services {
		if err := svc.Start(ctx); err != nil {
			sm.status[name] = "error"
			return fmt.Errorf("failed to start %s: %w", name, err)
		}
		sm.status[name] = "running"
	}
	return nil
}

// StopAll stops all registered services
func (sm *ServiceManager) StopAll(ctx context.Context) error {
	for name, svc := range sm.services {
		if err := svc.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop %s: %w", name, err)
		}
		sm.status[name] = "stopped"
	}
	return nil
}

// HealthCheck performs health checks on all services
func (sm *ServiceManager) HealthCheck(ctx context.Context) []HealthCheckResult {
	results := make([]HealthCheckResult, 0)

	for name, svc := range sm.services {
		start := time.Now()
		healthy, err := svc.Health(ctx)
		latency := time.Since(start)

		status := "healthy"
		if err != nil || !healthy {
			status = "unhealthy"
		}

		results = append(results, HealthCheckResult{
			ServiceName: name,
			Status:      status,
			Latency:     latency,
			Timestamp:   time.Now().Format(time.RFC3339),
		})
	}

	return results
}

// GetStatus returns the status of all services
func (sm *ServiceManager) GetStatus() map[string]string {
	return sm.status
}
