package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// HealthStatus represents the status of a health check
type HealthStatus int

const (
	HealthStatusUnknown  HealthStatus = iota // Unknown - check hasn't run yet
	HealthStatusHealthy                      // Healthy - check passed
	HealthStatusDegraded                     // Degraded - check passed but performance is impacted
	HealthStatusUnhealthy                    // Unhealthy - check failed
)

// String returns the string representation of a health status
func (s HealthStatus) String() string {
	switch s {
	case HealthStatusHealthy:
		return "healthy"
	case HealthStatusDegraded:
		return "degraded"
	case HealthStatusUnhealthy:
		return "unhealthy"
	case HealthStatusUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

// HealthCheckType represents the type of health check
type HealthCheckType int

const (
	HealthCheckTypeLiveness  HealthCheckType = iota // Liveness - is the app running?
	HealthCheckTypeReadiness                        // Readiness - can the app accept traffic?
	HealthCheckTypeStartup                          // Startup - has the app finished starting up?
)

// String returns the string representation of a health check type
func (t HealthCheckType) String() string {
	switch t {
	case HealthCheckTypeLiveness:
		return "liveness"
	case HealthCheckTypeReadiness:
		return "readiness"
	case HealthCheckTypeStartup:
		return "startup"
	default:
		return "unknown"
	}
}

// HealthCheckFunc is a function that performs a health check
// It should return the health status, a message, and any error encountered
type HealthCheckFunc func(ctx context.Context) (HealthStatus, string, error)

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	Name      string                 `json:"name"`
	Status    HealthStatus           `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration_ms"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// MarshalJSON custom JSON marshaling for HealthCheckResult
func (r *HealthCheckResult) MarshalJSON() ([]byte, error) {
	type Alias HealthCheckResult
	return json.Marshal(&struct {
		Status   string  `json:"status"`
		Duration float64 `json:"duration_ms"`
		*Alias
	}{
		Status:   r.Status.String(),
		Duration: float64(r.Duration.Milliseconds()),
		Alias:    (*Alias)(r),
	})
}

// HealthCheck represents a registered health check
type HealthCheck struct {
	Name     string
	Type     HealthCheckType
	Check    HealthCheckFunc
	Timeout  time.Duration
	Interval time.Duration // For periodic checks
	LastRun  time.Time
	Result   *HealthCheckResult
	mu       sync.RWMutex
}

// Run executes the health check with timeout
func (h *HealthCheck) Run(ctx context.Context) *HealthCheckResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	start := time.Now()

	// Create context with timeout
	checkCtx := ctx
	if h.Timeout > 0 {
		var cancel context.CancelFunc
		checkCtx, cancel = context.WithTimeout(ctx, h.Timeout)
		defer cancel()
	}

	// Run the check
	status, message, err := h.Check(checkCtx)

	duration := time.Since(start)

	result := &HealthCheckResult{
		Name:      h.Name,
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
		Duration:  duration,
		Metadata:  make(map[string]interface{}),
	}

	if err != nil {
		result.Error = err.Error()
		// If check returned an error, force status to unhealthy
		if status == HealthStatusHealthy {
			result.Status = HealthStatusUnhealthy
		}
	}

	h.LastRun = time.Now()
	h.Result = result

	return result
}

// HealthChecker manages health checks for the application
type HealthChecker struct {
	checks map[string]*HealthCheck
	mu     sync.RWMutex
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]*HealthCheck),
	}
}

// Register registers a new health check
func (hc *HealthChecker) Register(name string, checkType HealthCheckType, checkFunc HealthCheckFunc, timeout time.Duration) error {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if _, exists := hc.checks[name]; exists {
		return fmt.Errorf("health check '%s' already registered", name)
	}

	if timeout == 0 {
		timeout = 5 * time.Second // Default timeout
	}

	hc.checks[name] = &HealthCheck{
		Name:    name,
		Type:    checkType,
		Check:   checkFunc,
		Timeout: timeout,
	}

	return nil
}

// Unregister removes a health check
func (hc *HealthChecker) Unregister(name string) error {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if _, exists := hc.checks[name]; !exists {
		return fmt.Errorf("health check '%s' not found", name)
	}

	delete(hc.checks, name)
	return nil
}

// Check runs a specific health check by name
func (hc *HealthChecker) Check(ctx context.Context, name string) (*HealthCheckResult, error) {
	hc.mu.RLock()
	check, exists := hc.checks[name]
	hc.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("health check '%s' not found", name)
	}

	return check.Run(ctx), nil
}

// CheckAll runs all registered health checks
func (hc *HealthChecker) CheckAll(ctx context.Context) map[string]*HealthCheckResult {
	hc.mu.RLock()
	checks := make([]*HealthCheck, 0, len(hc.checks))
	for _, check := range hc.checks {
		checks = append(checks, check)
	}
	hc.mu.RUnlock()

	results := make(map[string]*HealthCheckResult)
	var resultsMu sync.Mutex

	// Run checks concurrently
	var wg sync.WaitGroup
	for _, check := range checks {
		wg.Add(1)
		go func(c *HealthCheck) {
			defer wg.Done()
			result := c.Run(ctx)
			resultsMu.Lock()
			results[c.Name] = result
			resultsMu.Unlock()
		}(check)
	}

	wg.Wait()
	return results
}

// CheckType runs all health checks of a specific type
func (hc *HealthChecker) CheckType(ctx context.Context, checkType HealthCheckType) map[string]*HealthCheckResult {
	hc.mu.RLock()
	checks := make([]*HealthCheck, 0)
	for _, check := range hc.checks {
		if check.Type == checkType {
			checks = append(checks, check)
		}
	}
	hc.mu.RUnlock()

	results := make(map[string]*HealthCheckResult)
	var resultsMu sync.Mutex

	// Run checks concurrently
	var wg sync.WaitGroup
	for _, check := range checks {
		wg.Add(1)
		go func(c *HealthCheck) {
			defer wg.Done()
			result := c.Run(ctx)
			resultsMu.Lock()
			results[c.Name] = result
			resultsMu.Unlock()
		}(check)
	}

	wg.Wait()
	return results
}

// Status returns the overall health status based on all checks
func (hc *HealthChecker) Status(ctx context.Context) HealthStatus {
	results := hc.CheckAll(ctx)

	if len(results) == 0 {
		return HealthStatusUnknown
	}

	overallStatus := HealthStatusHealthy
	for _, result := range results {
		if result.Status == HealthStatusUnhealthy {
			return HealthStatusUnhealthy
		}
		if result.Status == HealthStatusDegraded && overallStatus == HealthStatusHealthy {
			overallStatus = HealthStatusDegraded
		}
	}

	return overallStatus
}

// StatusByType returns the overall health status for a specific check type
func (hc *HealthChecker) StatusByType(ctx context.Context, checkType HealthCheckType) HealthStatus {
	results := hc.CheckType(ctx, checkType)

	if len(results) == 0 {
		return HealthStatusUnknown
	}

	overallStatus := HealthStatusHealthy
	for _, result := range results {
		if result.Status == HealthStatusUnhealthy {
			return HealthStatusUnhealthy
		}
		if result.Status == HealthStatusDegraded && overallStatus == HealthStatusHealthy {
			overallStatus = HealthStatusDegraded
		}
	}

	return overallStatus
}

// GetResult returns the last result for a specific check
func (hc *HealthChecker) GetResult(name string) (*HealthCheckResult, error) {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	check, exists := hc.checks[name]
	if !exists {
		return nil, fmt.Errorf("health check '%s' not found", name)
	}

	check.mu.RLock()
	defer check.mu.RUnlock()

	if check.Result == nil {
		return nil, fmt.Errorf("health check '%s' has not run yet", name)
	}

	return check.Result, nil
}

// List returns the names of all registered health checks
func (hc *HealthChecker) List() []string {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	names := make([]string, 0, len(hc.checks))
	for name := range hc.checks {
		names = append(names, name)
	}
	return names
}

// HealthReport represents a comprehensive health report
type HealthReport struct {
	Status    HealthStatus                  `json:"status"`
	Timestamp time.Time                     `json:"timestamp"`
	Checks    map[string]*HealthCheckResult `json:"checks"`
}

// MarshalJSON custom JSON marshaling for HealthReport
func (r *HealthReport) MarshalJSON() ([]byte, error) {
	type Alias HealthReport
	return json.Marshal(&struct {
		Status string `json:"status"`
		*Alias
	}{
		Status: r.Status.String(),
		Alias:  (*Alias)(r),
	})
}

// Report generates a comprehensive health report
func (hc *HealthChecker) Report(ctx context.Context) *HealthReport {
	results := hc.CheckAll(ctx)

	report := &HealthReport{
		Status:    hc.Status(ctx),
		Timestamp: time.Now(),
		Checks:    results,
	}

	return report
}

// ReportByType generates a health report for a specific check type
func (hc *HealthChecker) ReportByType(ctx context.Context, checkType HealthCheckType) *HealthReport {
	results := hc.CheckType(ctx, checkType)

	report := &HealthReport{
		Status:    hc.StatusByType(ctx, checkType),
		Timestamp: time.Now(),
		Checks:    results,
	}

	return report
}

// Global health checker instance
var (
	globalHealthChecker     *HealthChecker
	globalHealthCheckerOnce sync.Once
)

// GetGlobalHealthChecker returns the global health checker instance
func GetGlobalHealthChecker() *HealthChecker {
	globalHealthCheckerOnce.Do(func() {
		globalHealthChecker = NewHealthChecker()
		// Register default checks
		registerDefaultHealthChecks(globalHealthChecker)
	})
	return globalHealthChecker
}

// Package-level convenience functions

// RegisterHealthCheck registers a health check with the global checker
func RegisterHealthCheck(name string, checkType HealthCheckType, checkFunc HealthCheckFunc, timeout time.Duration) error {
	return GetGlobalHealthChecker().Register(name, checkType, checkFunc, timeout)
}

// UnregisterHealthCheck removes a health check from the global checker
func UnregisterHealthCheck(name string) error {
	return GetGlobalHealthChecker().Unregister(name)
}

// HealthCheckStatus runs a specific health check and returns its status
func HealthCheckStatus(ctx context.Context, name string) (*HealthCheckResult, error) {
	return GetGlobalHealthChecker().Check(ctx, name)
}

// GetHealthStatus returns the overall health status
func GetHealthStatus(ctx context.Context) HealthStatus {
	return GetGlobalHealthChecker().Status(ctx)
}

// GetHealthReport generates a comprehensive health report
func GetHealthReport(ctx context.Context) *HealthReport {
	return GetGlobalHealthChecker().Report(ctx)
}

// registerDefaultHealthChecks registers built-in health checks
func registerDefaultHealthChecks(hc *HealthChecker) {
	// Liveness check - always healthy if the process is running
	hc.Register("liveness", HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "Process is running", nil
	}, 1*time.Second)

	// Readiness check - checks if the runtime is ready to accept requests
	hc.Register("readiness", HealthCheckTypeReadiness, func(ctx context.Context) (HealthStatus, string, error) {
		// Basic readiness check - can be extended to check runtime state
		return HealthStatusHealthy, "Runtime is ready", nil
	}, 2*time.Second)
}
