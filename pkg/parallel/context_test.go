package parallel

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// ============================================================================
// RequestContext Tests
// ============================================================================

func TestNewRequestContext(t *testing.T) {
	ctx := NewRequestContext("req-1")

	if ctx.ID != "req-1" {
		t.Errorf("Expected ID 'req-1', got '%s'", ctx.ID)
	}

	if ctx.IsCancelled() {
		t.Error("New context should not be cancelled")
	}

	if ctx.GlobalCount() != 0 {
		t.Error("New context should have no globals")
	}

	if ctx.IsTimedOut() {
		t.Error("Context without timeout should never time out")
	}
}

func TestNewRequestContextWithTimeout(t *testing.T) {
	ctx := NewRequestContextWithTimeout("req-1", 100*time.Millisecond)

	if ctx.timeout != 100*time.Millisecond {
		t.Errorf("Expected timeout 100ms, got %v", ctx.timeout)
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	if !ctx.IsTimedOut() {
		t.Error("Context should be timed out")
	}
}

func TestRequestContextCancel(t *testing.T) {
	ctx := NewRequestContext("req-1")

	if ctx.IsCancelled() {
		t.Error("Context should not be cancelled initially")
	}

	ctx.Cancel()

	if !ctx.IsCancelled() {
		t.Error("Context should be cancelled after Cancel()")
	}
}

func TestRequestContextElapsed(t *testing.T) {
	ctx := NewRequestContext("req-1")

	time.Sleep(50 * time.Millisecond)

	elapsed := ctx.Elapsed()
	if elapsed < 50*time.Millisecond {
		t.Errorf("Expected elapsed >= 50ms, got %v", elapsed)
	}
}

// ============================================================================
// Global Variable Tests
// ============================================================================

func TestRequestContextGlobals(t *testing.T) {
	ctx := NewRequestContext("req-1")

	// Set globals
	ctx.SetGlobal("user", "john")
	ctx.SetGlobal("count", 42)

	// Get globals
	if value, ok := ctx.GetGlobal("user"); !ok || value != "john" {
		t.Errorf("Expected 'john', got %v", value)
	}

	if value, ok := ctx.GetGlobal("count"); !ok || value != 42 {
		t.Errorf("Expected 42, got %v", value)
	}

	// Has global
	if !ctx.HasGlobal("user") {
		t.Error("Should have 'user' global")
	}

	if ctx.HasGlobal("nonexistent") {
		t.Error("Should not have 'nonexistent' global")
	}

	// Count
	if ctx.GlobalCount() != 2 {
		t.Errorf("Expected 2 globals, got %d", ctx.GlobalCount())
	}
}

func TestRequestContextDeleteGlobal(t *testing.T) {
	ctx := NewRequestContext("req-1")

	ctx.SetGlobal("test", "value")

	if !ctx.HasGlobal("test") {
		t.Error("Should have 'test' global")
	}

	ctx.DeleteGlobal("test")

	if ctx.HasGlobal("test") {
		t.Error("Should not have 'test' global after deletion")
	}
}

func TestRequestContextClearGlobals(t *testing.T) {
	ctx := NewRequestContext("req-1")

	ctx.SetGlobal("a", 1)
	ctx.SetGlobal("b", 2)
	ctx.SetGlobal("c", 3)

	if ctx.GlobalCount() != 3 {
		t.Error("Should have 3 globals")
	}

	ctx.ClearGlobals()

	if ctx.GlobalCount() != 0 {
		t.Error("Should have no globals after clear")
	}
}

// ============================================================================
// Output Tests
// ============================================================================

func TestRequestContextOutput(t *testing.T) {
	ctx := NewRequestContext("req-1")

	// Write bytes
	n, err := ctx.Write([]byte("Hello "))
	if err != nil || n != 6 {
		t.Errorf("Write failed: %v, wrote %d bytes", err, n)
	}

	// Write string
	n, err = ctx.WriteString("World!")
	if err != nil || n != 6 {
		t.Errorf("WriteString failed: %v, wrote %d bytes", err, n)
	}

	// Get output
	output := ctx.GetOutput()
	if string(output) != "Hello World!" {
		t.Errorf("Expected 'Hello World!', got '%s'", string(output))
	}

	// Check size
	if ctx.OutputSize() != 12 {
		t.Errorf("Expected output size 12, got %d", ctx.OutputSize())
	}
}

func TestRequestContextClearOutput(t *testing.T) {
	ctx := NewRequestContext("req-1")

	ctx.WriteString("test output")

	if ctx.OutputSize() != 11 {
		t.Error("Should have output")
	}

	ctx.ClearOutput()

	if ctx.OutputSize() != 0 {
		t.Error("Output should be cleared")
	}
}

// ============================================================================
// Error Tests
// ============================================================================

func TestRequestContextErrors(t *testing.T) {
	ctx := NewRequestContext("req-1")

	if ctx.HasErrors() {
		t.Error("New context should have no errors")
	}

	// Add errors
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	ctx.AddError(err1)
	ctx.AddError(err2)

	if !ctx.HasErrors() {
		t.Error("Context should have errors")
	}

	if ctx.ErrorCount() != 2 {
		t.Errorf("Expected 2 errors, got %d", ctx.ErrorCount())
	}

	// Get errors
	errs := ctx.GetErrors()
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errs))
	}

	if errs[0] != err1 || errs[1] != err2 {
		t.Error("Errors don't match")
	}
}

func TestRequestContextAddNilError(t *testing.T) {
	ctx := NewRequestContext("req-1")

	ctx.AddError(nil)

	if ctx.HasErrors() {
		t.Error("Adding nil error should not create an error entry")
	}
}

func TestRequestContextClearErrors(t *testing.T) {
	ctx := NewRequestContext("req-1")

	ctx.AddError(errors.New("error 1"))
	ctx.AddError(errors.New("error 2"))

	if ctx.ErrorCount() != 2 {
		t.Error("Should have 2 errors")
	}

	ctx.ClearErrors()

	if ctx.ErrorCount() != 0 {
		t.Error("Should have no errors after clear")
	}
}

// ============================================================================
// Metadata Tests
// ============================================================================

func TestRequestContextMetadata(t *testing.T) {
	ctx := NewRequestContext("req-1")

	ctx.SetMetadata("ip", "192.168.1.1")
	ctx.SetMetadata("method", "GET")

	if value, ok := ctx.GetMetadata("ip"); !ok || value != "192.168.1.1" {
		t.Error("Metadata not stored correctly")
	}

	if !ctx.HasMetadata("method") {
		t.Error("Should have 'method' metadata")
	}

	if ctx.HasMetadata("nonexistent") {
		t.Error("Should not have 'nonexistent' metadata")
	}
}

// ============================================================================
// Cleanup Tests
// ============================================================================

func TestRequestContextCleanup(t *testing.T) {
	ctx := NewRequestContext("req-1")

	ctx.SetGlobal("test", "value")
	ctx.WriteString("output")
	ctx.AddError(errors.New("test error"))
	ctx.SetMetadata("key", "value")

	ctx.Cleanup()

	// Context should be cancelled
	if !ctx.IsCancelled() {
		t.Error("Context should be cancelled after cleanup")
	}

	// Note: Internal maps are set to nil, so we can't test them directly
	// but the cleanup should prevent further use
}

// ============================================================================
// RequestManager Tests
// ============================================================================

func TestNewRequestManager(t *testing.T) {
	rm := NewRequestManager(100, 5*time.Second)

	if rm.maxActive != 100 {
		t.Errorf("Expected max active 100, got %d", rm.maxActive)
	}

	if rm.timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", rm.timeout)
	}
}

func TestRequestManagerCreateRequest(t *testing.T) {
	rm := NewRequestManager(10, 0)

	ctx, err := rm.CreateRequest("req-1")
	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	if ctx.ID != "req-1" {
		t.Errorf("Expected ID 'req-1', got '%s'", ctx.ID)
	}

	if rm.ActiveCount() != 1 {
		t.Errorf("Expected 1 active request, got %d", rm.ActiveCount())
	}
}

func TestRequestManagerDuplicateRequest(t *testing.T) {
	rm := NewRequestManager(10, 0)

	_, err := rm.CreateRequest("req-1")
	if err != nil {
		t.Fatalf("Failed to create first request: %v", err)
	}

	// Try to create duplicate
	_, err = rm.CreateRequest("req-1")
	if err == nil {
		t.Error("Should not allow duplicate request IDs")
	}
}

func TestRequestManagerMaxActive(t *testing.T) {
	rm := NewRequestManager(2, 0)

	// Create max requests
	_, err := rm.CreateRequest("req-1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = rm.CreateRequest("req-2")
	if err != nil {
		t.Fatal(err)
	}

	// Try to exceed max
	_, err = rm.CreateRequest("req-3")
	if err == nil {
		t.Error("Should not allow exceeding max active requests")
	}
}

func TestRequestManagerGetRequest(t *testing.T) {
	rm := NewRequestManager(10, 0)

	ctx, _ := rm.CreateRequest("req-1")

	retrieved, ok := rm.GetRequest("req-1")
	if !ok {
		t.Error("Should find created request")
	}

	if retrieved != ctx {
		t.Error("Retrieved context should match created context")
	}

	_, ok = rm.GetRequest("nonexistent")
	if ok {
		t.Error("Should not find nonexistent request")
	}
}

func TestRequestManagerCompleteRequest(t *testing.T) {
	rm := NewRequestManager(10, 0)

	rm.CreateRequest("req-1")

	if rm.ActiveCount() != 1 {
		t.Error("Should have 1 active request")
	}

	err := rm.CompleteRequest("req-1")
	if err != nil {
		t.Errorf("Failed to complete request: %v", err)
	}

	if rm.ActiveCount() != 0 {
		t.Error("Should have 0 active requests after completion")
	}

	stats := rm.Stats()
	if stats.CompletedRequests != 1 {
		t.Errorf("Expected 1 completed request, got %d", stats.CompletedRequests)
	}
}

func TestRequestManagerCompleteNonexistent(t *testing.T) {
	rm := NewRequestManager(10, 0)

	err := rm.CompleteRequest("nonexistent")
	if err == nil {
		t.Error("Should error when completing nonexistent request")
	}
}

func TestRequestManagerCancelRequest(t *testing.T) {
	rm := NewRequestManager(10, 0)

	ctx, _ := rm.CreateRequest("req-1")

	if ctx.IsCancelled() {
		t.Error("Request should not be cancelled initially")
	}

	err := rm.CancelRequest("req-1")
	if err != nil {
		t.Errorf("Failed to cancel request: %v", err)
	}

	if !ctx.IsCancelled() {
		t.Error("Request should be cancelled")
	}
}

func TestRequestManagerCleanupTimedOut(t *testing.T) {
	rm := NewRequestManager(10, 50*time.Millisecond)

	// Create requests that will timeout
	rm.CreateRequest("req-1")
	rm.CreateRequest("req-2")

	if rm.ActiveCount() != 2 {
		t.Error("Should have 2 active requests")
	}

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)

	count := rm.CleanupTimedOut()
	if count != 2 {
		t.Errorf("Expected 2 timed out requests, got %d", count)
	}

	if rm.ActiveCount() != 0 {
		t.Error("Should have 0 active requests after cleanup")
	}

	stats := rm.Stats()
	if stats.TimedOutRequests != 2 {
		t.Errorf("Expected 2 timed out, got %d", stats.TimedOutRequests)
	}
}

func TestRequestManagerStats(t *testing.T) {
	rm := NewRequestManager(10, 0)

	// Create and complete some requests
	rm.CreateRequest("req-1")
	rm.CreateRequest("req-2")
	rm.CreateRequest("req-3")

	rm.CompleteRequest("req-1")

	stats := rm.Stats()

	if stats.TotalRequests != 3 {
		t.Errorf("Expected 3 total requests, got %d", stats.TotalRequests)
	}

	if stats.ActiveRequests != 2 {
		t.Errorf("Expected 2 active requests, got %d", stats.ActiveRequests)
	}

	if stats.CompletedRequests != 1 {
		t.Errorf("Expected 1 completed request, got %d", stats.CompletedRequests)
	}

	if stats.MaxActive != 10 {
		t.Errorf("Expected max active 10, got %d", stats.MaxActive)
	}
}

func TestRequestManagerShutdown(t *testing.T) {
	rm := NewRequestManager(10, 0)

	// Create multiple requests
	rm.CreateRequest("req-1")
	rm.CreateRequest("req-2")
	rm.CreateRequest("req-3")

	if rm.ActiveCount() != 3 {
		t.Error("Should have 3 active requests")
	}

	rm.Shutdown()

	if rm.ActiveCount() != 0 {
		t.Error("Should have 0 active requests after shutdown")
	}
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestRequestContextConcurrentGlobals(t *testing.T) {
	ctx := NewRequestContext("req-1")

	var wg sync.WaitGroup
	iterations := 100

	// Concurrent writes
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ctx.SetGlobal(fmt.Sprintf("key-%d", n), n)
		}(i)
	}

	wg.Wait()

	if ctx.GlobalCount() != iterations {
		t.Errorf("Expected %d globals, got %d", iterations, ctx.GlobalCount())
	}
}

func TestRequestContextConcurrentOutput(t *testing.T) {
	ctx := NewRequestContext("req-1")

	var wg sync.WaitGroup
	iterations := 100

	// Concurrent writes
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ctx.WriteString(fmt.Sprintf("%d", n))
		}(i)
	}

	wg.Wait()

	// Each iteration writes at least 1 byte (single digit)
	// and at most 3 bytes (three digits for n=100)
	if ctx.OutputSize() < iterations {
		t.Errorf("Expected at least %d bytes, got %d", iterations, ctx.OutputSize())
	}
}

func TestRequestManagerConcurrentRequests(t *testing.T) {
	rm := NewRequestManager(1000, 0)

	var wg sync.WaitGroup
	numRequests := 100

	// Create requests concurrently
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("req-%d", n)
			_, err := rm.CreateRequest(id)
			if err != nil {
				t.Errorf("Failed to create request %s: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	if rm.ActiveCount() != numRequests {
		t.Errorf("Expected %d active requests, got %d", numRequests, rm.ActiveCount())
	}

	// Complete requests concurrently
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("req-%d", n)
			rm.CompleteRequest(id)
		}(i)
	}

	wg.Wait()

	if rm.ActiveCount() != 0 {
		t.Errorf("Expected 0 active requests, got %d", rm.ActiveCount())
	}

	stats := rm.Stats()
	if stats.CompletedRequests != int64(numRequests) {
		t.Errorf("Expected %d completed, got %d", numRequests, stats.CompletedRequests)
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestRequestContextLifecycle(t *testing.T) {
	rm := NewRequestManager(10, 5*time.Second)

	// Create request
	ctx, err := rm.CreateRequest("req-1")
	if err != nil {
		t.Fatal(err)
	}

	// Use the context
	ctx.SetGlobal("user_id", 123)
	ctx.WriteString("Hello from request 1")
	ctx.SetMetadata("start_time", time.Now())

	// Simulate processing
	time.Sleep(10 * time.Millisecond)

	// Verify state
	if userID, ok := ctx.GetGlobal("user_id"); !ok || userID != 123 {
		t.Error("Global variable not preserved")
	}

	if ctx.OutputSize() == 0 {
		t.Error("Output not captured")
	}

	// Complete request
	err = rm.CompleteRequest("req-1")
	if err != nil {
		t.Errorf("Failed to complete request: %v", err)
	}

	// Verify cleanup
	_, ok := rm.GetRequest("req-1")
	if ok {
		t.Error("Request should be removed after completion")
	}
}
