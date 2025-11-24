package runtime

import (
	"context"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

// TestNewShutdownManager tests creating a new shutdown manager
func TestNewShutdownManager(t *testing.T) {
	timeout := 5 * time.Second
	sm := NewShutdownManager(timeout)

	if sm == nil {
		t.Fatal("NewShutdownManager returned nil")
	}

	if sm.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, sm.timeout)
	}

	if sm.shutdownInitiated {
		t.Error("Expected shutdownInitiated to be false")
	}

	if len(sm.handlers) != 0 {
		t.Errorf("Expected 0 handlers, got %d", len(sm.handlers))
	}
}

// TestGetShutdownManager tests the global shutdown manager
func TestGetShutdownManager(t *testing.T) {
	sm1 := GetShutdownManager()
	sm2 := GetShutdownManager()

	if sm1 != sm2 {
		t.Error("GetShutdownManager should return the same instance")
	}
}

// TestRegisterHandler tests registering shutdown handlers
func TestRegisterHandler(t *testing.T) {
	sm := NewShutdownManager(30 * time.Second)

	handler := func(ctx context.Context) error {
		return nil
	}

	sm.RegisterHandler("test", handler, ShutdownPriorityNormal)

	if len(sm.handlers) != 1 {
		t.Fatalf("Expected 1 handler, got %d", len(sm.handlers))
	}

	if sm.handlers[0].name != "test" {
		t.Errorf("Expected handler name 'test', got '%s'", sm.handlers[0].name)
	}

	if sm.handlers[0].priority != ShutdownPriorityNormal {
		t.Errorf("Expected priority %d, got %d", ShutdownPriorityNormal, sm.handlers[0].priority)
	}
}

// TestRegisterHandlerPriorityOrder tests that handlers are ordered by priority
func TestRegisterHandlerPriorityOrder(t *testing.T) {
	sm := NewShutdownManager(30 * time.Second)

	handler := func(ctx context.Context) error { return nil }

	sm.RegisterHandler("low", handler, ShutdownPriorityLow)
	sm.RegisterHandler("high", handler, ShutdownPriorityHigh)
	sm.RegisterHandler("normal", handler, ShutdownPriorityNormal)
	sm.RegisterHandler("highest", handler, ShutdownPriorityHighest)
	sm.RegisterHandler("lowest", handler, ShutdownPriorityLowest)

	if len(sm.handlers) != 5 {
		t.Fatalf("Expected 5 handlers, got %d", len(sm.handlers))
	}

	// Check order (highest to lowest priority)
	expected := []string{"highest", "high", "normal", "low", "lowest"}
	for i, name := range expected {
		if sm.handlers[i].name != name {
			t.Errorf("Expected handler %d to be '%s', got '%s'", i, name, sm.handlers[i].name)
		}
	}
}

// TestShutdown tests basic shutdown functionality
func TestShutdown(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	var callOrder []string
	var mu sync.Mutex

	createHandler := func(name string) ShutdownHandler {
		return func(ctx context.Context) error {
			mu.Lock()
			callOrder = append(callOrder, name)
			mu.Unlock()
			return nil
		}
	}

	sm.RegisterHandler("handler1", createHandler("handler1"), ShutdownPriorityHigh)
	sm.RegisterHandler("handler2", createHandler("handler2"), ShutdownPriorityNormal)
	sm.RegisterHandler("handler3", createHandler("handler3"), ShutdownPriorityLow)

	err := sm.Shutdown()
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	if !sm.IsShuttingDown() {
		t.Error("Expected IsShuttingDown to be true after shutdown")
	}

	// Check call order
	expectedOrder := []string{"handler1", "handler2", "handler3"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d handlers called, got %d", len(expectedOrder), len(callOrder))
	}

	for i, name := range expectedOrder {
		if callOrder[i] != name {
			t.Errorf("Expected handler %d to be '%s', got '%s'", i, name, callOrder[i])
		}
	}
}

// TestShutdownWithErrors tests shutdown when handlers return errors
func TestShutdownWithErrors(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	var callOrder []string
	var mu sync.Mutex

	sm.RegisterHandler("handler1", func(ctx context.Context) error {
		mu.Lock()
		callOrder = append(callOrder, "handler1")
		mu.Unlock()
		return errors.New("handler1 error")
	}, ShutdownPriorityHigh)

	sm.RegisterHandler("handler2", func(ctx context.Context) error {
		mu.Lock()
		callOrder = append(callOrder, "handler2")
		mu.Unlock()
		return nil
	}, ShutdownPriorityNormal)

	sm.RegisterHandler("handler3", func(ctx context.Context) error {
		mu.Lock()
		callOrder = append(callOrder, "handler3")
		mu.Unlock()
		return errors.New("handler3 error")
	}, ShutdownPriorityLow)

	err := sm.Shutdown()
	if err == nil {
		t.Fatal("Expected shutdown to return error")
	}

	// All handlers should still be called despite errors
	expectedOrder := []string{"handler1", "handler2", "handler3"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d handlers called, got %d", len(expectedOrder), len(callOrder))
	}

	for i, name := range expectedOrder {
		if callOrder[i] != name {
			t.Errorf("Expected handler %d to be '%s', got '%s'", i, name, callOrder[i])
		}
	}
}

// TestShutdownWithPanic tests shutdown when a handler panics
func TestShutdownWithPanic(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	var callOrder []string
	var mu sync.Mutex

	sm.RegisterHandler("handler1", func(ctx context.Context) error {
		mu.Lock()
		callOrder = append(callOrder, "handler1")
		mu.Unlock()
		panic("handler1 panic")
	}, ShutdownPriorityHigh)

	sm.RegisterHandler("handler2", func(ctx context.Context) error {
		mu.Lock()
		callOrder = append(callOrder, "handler2")
		mu.Unlock()
		return nil
	}, ShutdownPriorityNormal)

	err := sm.Shutdown()
	if err == nil {
		t.Fatal("Expected shutdown to return error when handler panics")
	}

	// handler2 should still be called after handler1 panics
	expectedOrder := []string{"handler1", "handler2"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d handlers called, got %d", len(expectedOrder), len(callOrder))
	}

	for i, name := range expectedOrder {
		if callOrder[i] != name {
			t.Errorf("Expected handler %d to be '%s', got '%s'", i, name, callOrder[i])
		}
	}
}

// TestShutdownTimeout tests shutdown timeout
func TestShutdownTimeout(t *testing.T) {
	sm := NewShutdownManager(100 * time.Millisecond)

	var handler1Done, handler2Started atomic.Bool

	sm.RegisterHandler("handler1", func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		handler1Done.Store(true)
		return nil
	}, ShutdownPriorityHigh)

	sm.RegisterHandler("handler2", func(ctx context.Context) error {
		handler2Started.Store(true)
		time.Sleep(200 * time.Millisecond) // This will timeout
		return nil
	}, ShutdownPriorityNormal)

	sm.RegisterHandler("handler3", func(ctx context.Context) error {
		t.Error("handler3 should not be called due to timeout")
		return nil
	}, ShutdownPriorityLow)

	err := sm.Shutdown()
	if err == nil {
		t.Fatal("Expected shutdown to return error due to timeout")
	}

	if !handler1Done.Load() {
		t.Error("handler1 should have completed")
	}

	if !handler2Started.Load() {
		t.Error("handler2 should have started")
	}
}

// TestShutdownIdempotent tests that shutdown can only be called once
func TestShutdownIdempotent(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	var callCount atomic.Int32
	sm.RegisterHandler("handler", func(ctx context.Context) error {
		callCount.Add(1)
		return nil
	}, ShutdownPriorityNormal)

	err1 := sm.Shutdown()
	if err1 != nil {
		t.Fatalf("First shutdown failed: %v", err1)
	}

	err2 := sm.Shutdown()
	if err2 == nil {
		t.Error("Second shutdown should return error")
	}

	if callCount.Load() != 1 {
		t.Errorf("Expected handler to be called once, got %d", callCount.Load())
	}
}

// TestWaitForShutdown tests waiting for shutdown
func TestWaitForShutdown(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	done := make(chan struct{})
	go func() {
		sm.WaitForShutdown()
		close(done)
	}()

	// Wait a bit to ensure goroutine is waiting
	time.Sleep(10 * time.Millisecond)

	select {
	case <-done:
		t.Fatal("WaitForShutdown should not return before shutdown")
	default:
	}

	// Initiate shutdown
	sm.Shutdown()

	// Now it should return
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitForShutdown did not return after shutdown")
	}
}

// TestWaitForShutdownComplete tests waiting for shutdown to complete
func TestWaitForShutdownComplete(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	sm.RegisterHandler("slow-handler", func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}, ShutdownPriorityNormal)

	done := make(chan struct{})
	go func() {
		sm.WaitForShutdownComplete()
		close(done)
	}()

	// Wait a bit to ensure goroutine is waiting
	time.Sleep(10 * time.Millisecond)

	select {
	case <-done:
		t.Fatal("WaitForShutdownComplete should not return before shutdown")
	default:
	}

	// Start shutdown in goroutine
	go sm.Shutdown()

	// Should return after shutdown completes
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("WaitForShutdownComplete did not return after shutdown")
	}
}

// TestSetTimeout tests setting shutdown timeout
func TestSetTimeout(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	initialTimeout := sm.GetTimeout()
	if initialTimeout != 5*time.Second {
		t.Errorf("Expected initial timeout 5s, got %v", initialTimeout)
	}

	newTimeout := 10 * time.Second
	sm.SetTimeout(newTimeout)

	if sm.GetTimeout() != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, sm.GetTimeout())
	}
}

// TestShutdownPackageLevelFunctions tests package-level convenience functions for shutdown
func TestShutdownPackageLevelFunctions(t *testing.T) {
	// Reset global shutdown manager for this test
	shutdownManagerOnce = sync.Once{}
	globalShutdownManager = nil

	var called atomic.Bool
	RegisterShutdownHandler("test", func(ctx context.Context) error {
		called.Store(true)
		return nil
	}, ShutdownPriorityNormal)

	if IsShuttingDown() {
		t.Error("Expected IsShuttingDown to be false initially")
	}

	originalTimeout := GetShutdownTimeout()
	newTimeout := originalTimeout + 5*time.Second
	SetShutdownTimeout(newTimeout)

	if GetShutdownTimeout() != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, GetShutdownTimeout())
	}

	err := Shutdown()
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	if !IsShuttingDown() {
		t.Error("Expected IsShuttingDown to be true after shutdown")
	}

	if !called.Load() {
		t.Error("Expected handler to be called")
	}
}

// TestConcurrentRegistration tests concurrent handler registration
func TestConcurrentRegistration(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sm.RegisterHandler(
				string(rune('A'+id)),
				func(ctx context.Context) error { return nil },
				ShutdownPriorityNormal,
			)
		}(i)
	}

	wg.Wait()

	if len(sm.handlers) != numGoroutines {
		t.Errorf("Expected %d handlers, got %d", numGoroutines, len(sm.handlers))
	}
}

// TestContextCancellation tests that handlers receive context cancellation
func TestContextCancellation(t *testing.T) {
	sm := NewShutdownManager(100 * time.Millisecond)

	var contextCancelled atomic.Bool

	sm.RegisterHandler("slow-handler", func(ctx context.Context) error {
		select {
		case <-time.After(1 * time.Second):
			return nil
		case <-ctx.Done():
			contextCancelled.Store(true)
			return ctx.Err()
		}
	}, ShutdownPriorityNormal)

	sm.Shutdown()

	if !contextCancelled.Load() {
		t.Error("Expected context to be cancelled")
	}
}

// TestShutdownWithNoHandlers tests shutdown with no registered handlers
func TestShutdownWithNoHandlers(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	err := sm.Shutdown()
	if err != nil {
		t.Fatalf("Shutdown with no handlers should succeed: %v", err)
	}

	if !sm.IsShuttingDown() {
		t.Error("Expected IsShuttingDown to be true after shutdown")
	}
}

// TestListenForSignals tests signal handling
func TestListenForSignals(t *testing.T) {
	sm := NewShutdownManager(5 * time.Second)

	var handlerCalled atomic.Bool
	sm.RegisterHandler("signal-test", func(ctx context.Context) error {
		handlerCalled.Store(true)
		return nil
	}, ShutdownPriorityNormal)

	// Start listening for signals
	sm.ListenForSignals()

	// Give it a moment to set up the signal handler
	time.Sleep(10 * time.Millisecond)

	// Send SIGTERM to ourselves
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}

	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		t.Fatalf("Failed to send signal: %v", err)
	}

	// Wait for shutdown to complete
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("Shutdown did not complete after signal")
	default:
		// Give it a moment to process
		time.Sleep(100 * time.Millisecond)
	}

	if !handlerCalled.Load() {
		t.Error("Handler should have been called after signal")
	}

	if !sm.IsShuttingDown() {
		t.Error("Expected IsShuttingDown to be true after signal")
	}
}

// TestPackageLevelWaitFunctions tests package-level wait functions
func TestPackageLevelWaitFunctions(t *testing.T) {
	// Reset global shutdown manager
	shutdownManagerOnce = sync.Once{}
	globalShutdownManager = nil

	// Start goroutine that waits for shutdown
	done := make(chan struct{})
	go func() {
		WaitForShutdown()
		close(done)
	}()

	// Give it time to start waiting
	time.Sleep(10 * time.Millisecond)

	// Trigger shutdown
	go Shutdown()

	// Should complete
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("WaitForShutdown did not return")
	}

	// Now test WaitForShutdownComplete
	// Reset again
	shutdownManagerOnce = sync.Once{}
	globalShutdownManager = nil

	done2 := make(chan struct{})
	go func() {
		WaitForShutdownComplete()
		close(done2)
	}()

	// Give it time to start waiting
	time.Sleep(10 * time.Millisecond)

	// Trigger shutdown
	go Shutdown()

	// Should complete
	select {
	case <-done2:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("WaitForShutdownComplete did not return")
	}
}

// Benchmark tests

// BenchmarkShutdown benchmarks basic shutdown with multiple handlers
func BenchmarkShutdown(b *testing.B) {
	handler := func(ctx context.Context) error { return nil }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm := NewShutdownManager(30 * time.Second)
		sm.RegisterHandler("handler1", handler, ShutdownPriorityHigh)
		sm.RegisterHandler("handler2", handler, ShutdownPriorityNormal)
		sm.RegisterHandler("handler3", handler, ShutdownPriorityLow)
		sm.Shutdown()
	}
}

// BenchmarkRegisterHandler benchmarks handler registration
func BenchmarkRegisterHandler(b *testing.B) {
	sm := NewShutdownManager(30 * time.Second)
	handler := func(ctx context.Context) error { return nil }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.RegisterHandler("handler", handler, ShutdownPriorityNormal)
	}
}

// BenchmarkIsShuttingDown benchmarks the IsShuttingDown check
func BenchmarkIsShuttingDown(b *testing.B) {
	sm := NewShutdownManager(30 * time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.IsShuttingDown()
	}
}
