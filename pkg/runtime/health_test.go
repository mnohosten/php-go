package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestHealthStatus tests health status string representation
func TestHealthStatus(t *testing.T) {
	tests := []struct {
		status   HealthStatus
		expected string
	}{
		{HealthStatusHealthy, "healthy"},
		{HealthStatusDegraded, "degraded"},
		{HealthStatusUnhealthy, "unhealthy"},
		{HealthStatusUnknown, "unknown"},
		{HealthStatus(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("HealthStatus.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestHealthCheckType tests health check type string representation
func TestHealthCheckType(t *testing.T) {
	tests := []struct {
		checkType HealthCheckType
		expected  string
	}{
		{HealthCheckTypeLiveness, "liveness"},
		{HealthCheckTypeReadiness, "readiness"},
		{HealthCheckTypeStartup, "startup"},
		{HealthCheckType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.checkType.String(); got != tt.expected {
				t.Errorf("HealthCheckType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestNewHealthChecker tests creating a new health checker
func TestNewHealthChecker(t *testing.T) {
	hc := NewHealthChecker()
	if hc == nil {
		t.Fatal("NewHealthChecker() returned nil")
	}
	if hc.checks == nil {
		t.Error("Health checker checks map is nil")
	}
}

// TestRegisterHealthCheck tests registering a health check
func TestRegisterHealthCheck(t *testing.T) {
	hc := NewHealthChecker()

	checkFunc := func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}

	// Test successful registration
	err := hc.Register("test_check", HealthCheckTypeLiveness, checkFunc, 5*time.Second)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Verify check is registered
	checks := hc.List()
	if len(checks) != 1 || checks[0] != "test_check" {
		t.Errorf("List() = %v, want [test_check]", checks)
	}

	// Test duplicate registration
	err = hc.Register("test_check", HealthCheckTypeLiveness, checkFunc, 5*time.Second)
	if err == nil {
		t.Error("Expected error for duplicate registration")
	}
}

// TestUnregisterHealthCheck tests unregistering a health check
func TestUnregisterHealthCheck(t *testing.T) {
	hc := NewHealthChecker()

	checkFunc := func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}

	hc.Register("test_check", HealthCheckTypeLiveness, checkFunc, 5*time.Second)

	// Test successful unregistration
	err := hc.Unregister("test_check")
	if err != nil {
		t.Fatalf("Unregister() error = %v", err)
	}

	// Verify check is unregistered
	checks := hc.List()
	if len(checks) != 0 {
		t.Errorf("List() = %v, want []", checks)
	}

	// Test unregistering non-existent check
	err = hc.Unregister("non_existent")
	if err == nil {
		t.Error("Expected error for unregistering non-existent check")
	}
}

// TestHealthCheckRun tests running a single health check
func TestHealthCheckRun(t *testing.T) {
	tests := []struct {
		name           string
		checkFunc      HealthCheckFunc
		expectedStatus HealthStatus
		expectError    bool
	}{
		{
			name: "healthy check",
			checkFunc: func(ctx context.Context) (HealthStatus, string, error) {
				return HealthStatusHealthy, "All systems operational", nil
			},
			expectedStatus: HealthStatusHealthy,
			expectError:    false,
		},
		{
			name: "degraded check",
			checkFunc: func(ctx context.Context) (HealthStatus, string, error) {
				return HealthStatusDegraded, "Performance degraded", nil
			},
			expectedStatus: HealthStatusDegraded,
			expectError:    false,
		},
		{
			name: "unhealthy check",
			checkFunc: func(ctx context.Context) (HealthStatus, string, error) {
				return HealthStatusUnhealthy, "Service unavailable", nil
			},
			expectedStatus: HealthStatusUnhealthy,
			expectError:    false,
		},
		{
			name: "check with error",
			checkFunc: func(ctx context.Context) (HealthStatus, string, error) {
				return HealthStatusHealthy, "Check failed", errors.New("connection error")
			},
			expectedStatus: HealthStatusUnhealthy, // Should be unhealthy when error occurs
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := &HealthCheck{
				Name:    tt.name,
				Type:    HealthCheckTypeLiveness,
				Check:   tt.checkFunc,
				Timeout: 5 * time.Second,
			}

			ctx := context.Background()
			result := check.Run(ctx)

			if result.Status != tt.expectedStatus {
				t.Errorf("Run() status = %v, want %v", result.Status, tt.expectedStatus)
			}

			if tt.expectError && result.Error == "" {
				t.Error("Expected error in result but got none")
			}

			if !tt.expectError && result.Error != "" {
				t.Errorf("Unexpected error in result: %s", result.Error)
			}

			if result.Name != tt.name {
				t.Errorf("Result name = %v, want %v", result.Name, tt.name)
			}

			if result.Duration == 0 {
				t.Error("Result duration is 0")
			}

			if result.Timestamp.IsZero() {
				t.Error("Result timestamp is zero")
			}
		})
	}
}

// TestHealthCheckTimeout tests health check timeout
func TestHealthCheckTimeout(t *testing.T) {
	checkFunc := func(ctx context.Context) (HealthStatus, string, error) {
		// Simulate a slow check
		select {
		case <-time.After(2 * time.Second):
			return HealthStatusHealthy, "OK", nil
		case <-ctx.Done():
			return HealthStatusUnhealthy, "Timeout", ctx.Err()
		}
	}

	check := &HealthCheck{
		Name:    "slow_check",
		Type:    HealthCheckTypeLiveness,
		Check:   checkFunc,
		Timeout: 100 * time.Millisecond, // Short timeout
	}

	ctx := context.Background()
	result := check.Run(ctx)

	// Check should timeout and return unhealthy
	if result.Status == HealthStatusHealthy {
		t.Error("Expected unhealthy status due to timeout")
	}

	if result.Error == "" {
		t.Error("Expected error due to timeout")
	}
}

// TestHealthCheckerCheck tests checking a specific health check
func TestHealthCheckerCheck(t *testing.T) {
	hc := NewHealthChecker()

	checkFunc := func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}

	hc.Register("test_check", HealthCheckTypeLiveness, checkFunc, 5*time.Second)

	ctx := context.Background()
	result, err := hc.Check(ctx, "test_check")

	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}

	if result.Status != HealthStatusHealthy {
		t.Errorf("Check() status = %v, want %v", result.Status, HealthStatusHealthy)
	}

	// Test checking non-existent check
	_, err = hc.Check(ctx, "non_existent")
	if err == nil {
		t.Error("Expected error for non-existent check")
	}
}

// TestHealthCheckerCheckAll tests running all health checks
func TestHealthCheckerCheckAll(t *testing.T) {
	hc := NewHealthChecker()

	// Register multiple checks
	hc.Register("check1", HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	hc.Register("check2", HealthCheckTypeReadiness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	hc.Register("check3", HealthCheckTypeStartup, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusDegraded, "Slow", nil
	}, 5*time.Second)

	ctx := context.Background()
	results := hc.CheckAll(ctx)

	if len(results) != 3 {
		t.Errorf("CheckAll() returned %d results, want 3", len(results))
	}

	if _, ok := results["check1"]; !ok {
		t.Error("CheckAll() missing check1")
	}

	if _, ok := results["check2"]; !ok {
		t.Error("CheckAll() missing check2")
	}

	if _, ok := results["check3"]; !ok {
		t.Error("CheckAll() missing check3")
	}
}

// TestHealthCheckerCheckType tests running health checks by type
func TestHealthCheckerCheckType(t *testing.T) {
	hc := NewHealthChecker()

	// Register checks of different types
	hc.Register("liveness1", HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	hc.Register("liveness2", HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	hc.Register("readiness1", HealthCheckTypeReadiness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	ctx := context.Background()

	// Check liveness type
	livenessResults := hc.CheckType(ctx, HealthCheckTypeLiveness)
	if len(livenessResults) != 2 {
		t.Errorf("CheckType(Liveness) returned %d results, want 2", len(livenessResults))
	}

	// Check readiness type
	readinessResults := hc.CheckType(ctx, HealthCheckTypeReadiness)
	if len(readinessResults) != 1 {
		t.Errorf("CheckType(Readiness) returned %d results, want 1", len(readinessResults))
	}

	// Check startup type (none registered)
	startupResults := hc.CheckType(ctx, HealthCheckTypeStartup)
	if len(startupResults) != 0 {
		t.Errorf("CheckType(Startup) returned %d results, want 0", len(startupResults))
	}
}

// TestHealthCheckerStatus tests overall health status
func TestHealthCheckerStatus(t *testing.T) {
	tests := []struct {
		name           string
		checks         map[string]HealthStatus
		expectedStatus HealthStatus
	}{
		{
			name: "all healthy",
			checks: map[string]HealthStatus{
				"check1": HealthStatusHealthy,
				"check2": HealthStatusHealthy,
			},
			expectedStatus: HealthStatusHealthy,
		},
		{
			name: "one degraded",
			checks: map[string]HealthStatus{
				"check1": HealthStatusHealthy,
				"check2": HealthStatusDegraded,
			},
			expectedStatus: HealthStatusDegraded,
		},
		{
			name: "one unhealthy",
			checks: map[string]HealthStatus{
				"check1": HealthStatusHealthy,
				"check2": HealthStatusUnhealthy,
			},
			expectedStatus: HealthStatusUnhealthy,
		},
		{
			name: "multiple unhealthy",
			checks: map[string]HealthStatus{
				"check1": HealthStatusUnhealthy,
				"check2": HealthStatusUnhealthy,
			},
			expectedStatus: HealthStatusUnhealthy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := NewHealthChecker()

			for name, status := range tt.checks {
				capturedStatus := status // Capture for closure
				hc.Register(name, HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
					return capturedStatus, "Test", nil
				}, 5*time.Second)
			}

			ctx := context.Background()
			status := hc.Status(ctx)

			if status != tt.expectedStatus {
				t.Errorf("Status() = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}

// TestHealthCheckerStatusByType tests status by check type
func TestHealthCheckerStatusByType(t *testing.T) {
	hc := NewHealthChecker()

	// Register checks of different types with different statuses
	hc.Register("liveness_healthy", HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	hc.Register("readiness_unhealthy", HealthCheckTypeReadiness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusUnhealthy, "Not ready", nil
	}, 5*time.Second)

	ctx := context.Background()

	// Liveness should be healthy
	livenessStatus := hc.StatusByType(ctx, HealthCheckTypeLiveness)
	if livenessStatus != HealthStatusHealthy {
		t.Errorf("StatusByType(Liveness) = %v, want %v", livenessStatus, HealthStatusHealthy)
	}

	// Readiness should be unhealthy
	readinessStatus := hc.StatusByType(ctx, HealthCheckTypeReadiness)
	if readinessStatus != HealthStatusUnhealthy {
		t.Errorf("StatusByType(Readiness) = %v, want %v", readinessStatus, HealthStatusUnhealthy)
	}

	// Startup should be unknown (no checks registered)
	startupStatus := hc.StatusByType(ctx, HealthCheckTypeStartup)
	if startupStatus != HealthStatusUnknown {
		t.Errorf("StatusByType(Startup) = %v, want %v", startupStatus, HealthStatusUnknown)
	}
}

// TestHealthCheckerGetResult tests getting the last result
func TestHealthCheckerGetResult(t *testing.T) {
	hc := NewHealthChecker()

	checkFunc := func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}

	hc.Register("test_check", HealthCheckTypeLiveness, checkFunc, 5*time.Second)

	// Test getting result before check has run
	_, err := hc.GetResult("test_check")
	if err == nil {
		t.Error("Expected error when getting result before check has run")
	}

	// Run the check
	ctx := context.Background()
	hc.Check(ctx, "test_check")

	// Test getting result after check has run
	result, err := hc.GetResult("test_check")
	if err != nil {
		t.Fatalf("GetResult() error = %v", err)
	}

	if result.Status != HealthStatusHealthy {
		t.Errorf("GetResult() status = %v, want %v", result.Status, HealthStatusHealthy)
	}

	// Test getting result for non-existent check
	_, err = hc.GetResult("non_existent")
	if err == nil {
		t.Error("Expected error for non-existent check")
	}
}

// TestHealthCheckerList tests listing health checks
func TestHealthCheckerList(t *testing.T) {
	hc := NewHealthChecker()

	checkFunc := func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}

	// Empty list
	list := hc.List()
	if len(list) != 0 {
		t.Errorf("List() = %v, want []", list)
	}

	// Add some checks
	hc.Register("check1", HealthCheckTypeLiveness, checkFunc, 5*time.Second)
	hc.Register("check2", HealthCheckTypeReadiness, checkFunc, 5*time.Second)

	list = hc.List()
	if len(list) != 2 {
		t.Errorf("List() returned %d checks, want 2", len(list))
	}

	// Verify all checks are in the list
	found := make(map[string]bool)
	for _, name := range list {
		found[name] = true
	}

	if !found["check1"] || !found["check2"] {
		t.Errorf("List() = %v, missing expected checks", list)
	}
}

// TestHealthReport tests generating a health report
func TestHealthReport(t *testing.T) {
	hc := NewHealthChecker()

	hc.Register("check1", HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	hc.Register("check2", HealthCheckTypeReadiness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusDegraded, "Slow", nil
	}, 5*time.Second)

	ctx := context.Background()
	report := hc.Report(ctx)

	if report.Status != HealthStatusDegraded {
		t.Errorf("Report() status = %v, want %v", report.Status, HealthStatusDegraded)
	}

	if len(report.Checks) != 2 {
		t.Errorf("Report() has %d checks, want 2", len(report.Checks))
	}

	if report.Timestamp.IsZero() {
		t.Error("Report() timestamp is zero")
	}
}

// TestHealthReportByType tests generating a health report by type
func TestHealthReportByType(t *testing.T) {
	hc := NewHealthChecker()

	hc.Register("liveness1", HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	hc.Register("readiness1", HealthCheckTypeReadiness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	ctx := context.Background()

	// Test liveness report
	livenessReport := hc.ReportByType(ctx, HealthCheckTypeLiveness)
	if len(livenessReport.Checks) != 1 {
		t.Errorf("ReportByType(Liveness) has %d checks, want 1", len(livenessReport.Checks))
	}

	// Test readiness report
	readinessReport := hc.ReportByType(ctx, HealthCheckTypeReadiness)
	if len(readinessReport.Checks) != 1 {
		t.Errorf("ReportByType(Readiness) has %d checks, want 1", len(readinessReport.Checks))
	}
}

// TestHealthCheckResultJSON tests JSON marshaling of health check results
func TestHealthCheckResultJSON(t *testing.T) {
	result := &HealthCheckResult{
		Name:      "test_check",
		Status:    HealthStatusHealthy,
		Message:   "All systems operational",
		Timestamp: time.Now(),
		Duration:  100 * time.Millisecond,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(data)

	// Verify JSON contains expected fields
	if !strings.Contains(jsonStr, `"name":"test_check"`) {
		t.Error("JSON missing name field")
	}

	if !strings.Contains(jsonStr, `"status":"healthy"`) {
		t.Error("JSON missing or incorrect status field")
	}

	if !strings.Contains(jsonStr, `"message":"All systems operational"`) {
		t.Error("JSON missing message field")
	}

	if !strings.Contains(jsonStr, `"duration_ms"`) {
		t.Error("JSON missing duration_ms field")
	}
}

// TestHealthReportJSON tests JSON marshaling of health reports
func TestHealthReportJSON(t *testing.T) {
	hc := NewHealthChecker()

	hc.Register("check1", HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}, 5*time.Second)

	ctx := context.Background()
	report := hc.Report(ctx)

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(data)

	// Verify JSON contains expected fields
	if !strings.Contains(jsonStr, `"status":"healthy"`) {
		t.Error("JSON missing or incorrect status field")
	}

	if !strings.Contains(jsonStr, `"checks"`) {
		t.Error("JSON missing checks field")
	}

	if !strings.Contains(jsonStr, `"timestamp"`) {
		t.Error("JSON missing timestamp field")
	}
}

// TestGlobalHealthChecker tests the global health checker
func TestGlobalHealthChecker(t *testing.T) {
	hc := GetGlobalHealthChecker()
	if hc == nil {
		t.Fatal("GetGlobalHealthChecker() returned nil")
	}

	// Global health checker should have default checks registered
	checks := hc.List()
	if len(checks) < 2 {
		t.Errorf("Global health checker has %d checks, expected at least 2 default checks", len(checks))
	}

	// Verify default liveness check
	ctx := context.Background()
	result, err := hc.Check(ctx, "liveness")
	if err != nil {
		t.Errorf("Default liveness check not found: %v", err)
	} else if result.Status != HealthStatusHealthy {
		t.Errorf("Default liveness check status = %v, want %v", result.Status, HealthStatusHealthy)
	}

	// Verify default readiness check
	result, err = hc.Check(ctx, "readiness")
	if err != nil {
		t.Errorf("Default readiness check not found: %v", err)
	} else if result.Status != HealthStatusHealthy {
		t.Errorf("Default readiness check status = %v, want %v", result.Status, HealthStatusHealthy)
	}
}

// TestPackageLevelFunctions tests package-level convenience functions
func TestPackageLevelFunctions(t *testing.T) {
	// Clean up any existing custom checks from other tests
	hc := GetGlobalHealthChecker()
	for _, name := range hc.List() {
		if name != "liveness" && name != "readiness" {
			hc.Unregister(name)
		}
	}

	checkFunc := func(ctx context.Context) (HealthStatus, string, error) {
		return HealthStatusHealthy, "OK", nil
	}

	// Test RegisterHealthCheck
	err := RegisterHealthCheck("package_test", HealthCheckTypeLiveness, checkFunc, 5*time.Second)
	if err != nil {
		t.Fatalf("RegisterHealthCheck() error = %v", err)
	}

	// Test HealthCheckStatus
	ctx := context.Background()
	result, err := HealthCheckStatus(ctx, "package_test")
	if err != nil {
		t.Fatalf("HealthCheckStatus() error = %v", err)
	}
	if result.Status != HealthStatusHealthy {
		t.Errorf("HealthCheckStatus() status = %v, want %v", result.Status, HealthStatusHealthy)
	}

	// Test GetHealthStatus
	status := GetHealthStatus(ctx)
	if status != HealthStatusHealthy {
		t.Errorf("GetHealthStatus() = %v, want %v", status, HealthStatusHealthy)
	}

	// Test GetHealthReport
	report := GetHealthReport(ctx)
	if report == nil {
		t.Fatal("GetHealthReport() returned nil")
	}
	if len(report.Checks) < 3 {
		t.Errorf("HealthReport() has %d checks, expected at least 3", len(report.Checks))
	}

	// Test UnregisterHealthCheck
	err = UnregisterHealthCheck("package_test")
	if err != nil {
		t.Fatalf("UnregisterHealthCheck() error = %v", err)
	}
}

// TestConcurrentHealthChecks tests concurrent health check execution
func TestConcurrentHealthChecks(t *testing.T) {
	hc := NewHealthChecker()

	// Register multiple checks that take some time
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("check%d", i)
		hc.Register(name, HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
			time.Sleep(10 * time.Millisecond)
			return HealthStatusHealthy, "OK", nil
		}, 5*time.Second)
	}

	ctx := context.Background()
	start := time.Now()
	results := hc.CheckAll(ctx)
	duration := time.Since(start)

	if len(results) != 10 {
		t.Errorf("CheckAll() returned %d results, want 10", len(results))
	}

	// If checks ran concurrently, total duration should be much less than 10 * 10ms = 100ms
	// Allow some overhead but it should definitely be under 50ms
	if duration > 50*time.Millisecond {
		t.Errorf("CheckAll() took %v, expected concurrent execution to be faster", duration)
	}
}

// BenchmarkHealthCheckRun benchmarks running a single health check
func BenchmarkHealthCheckRun(b *testing.B) {
	check := &HealthCheck{
		Name: "bench_check",
		Type: HealthCheckTypeLiveness,
		Check: func(ctx context.Context) (HealthStatus, string, error) {
			return HealthStatusHealthy, "OK", nil
		},
		Timeout: 5 * time.Second,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		check.Run(ctx)
	}
}

// BenchmarkHealthCheckerCheckAll benchmarks running all health checks
func BenchmarkHealthCheckerCheckAll(b *testing.B) {
	hc := NewHealthChecker()

	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("check%d", i)
		hc.Register(name, HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
			return HealthStatusHealthy, "OK", nil
		}, 5*time.Second)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hc.CheckAll(ctx)
	}
}

// BenchmarkHealthReport benchmarks generating a health report
func BenchmarkHealthReport(b *testing.B) {
	hc := NewHealthChecker()

	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("check%d", i)
		hc.Register(name, HealthCheckTypeLiveness, func(ctx context.Context) (HealthStatus, string, error) {
			return HealthStatusHealthy, "OK", nil
		}, 5*time.Second)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hc.Report(ctx)
	}
}
